// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	aosOpsUrlPath   = "/send"
	EnvAosOpsEdgeId = "API_OPS_DATACENTER_EDGE_ID"
)

func (o *Client) talkToApiOps(ctx context.Context, in *talkToApstraIn) error {
	// create URL we'd use if we were talking to an actual Apstra server
	apstraUrl, err := o.craftUrl(in)
	if err != nil {
		return err
	}

	values := apstraUrl.Query()
	params := make(map[string]string, len(values))
	for k, v := range values {
		switch len(v) {
		case 0:
		case 1:
			params[k] = v[0]
		default:
			return fmt.Errorf("cannot format query string param %q for the proxy: only one string supported per param, got %d strings: %s", k, len(v), v)
		}
	}

	o.lock(mutexKeyHttpHeaders)
	headers := make(map[string]string, len(o.httpHeaders)+1)
	for k, v := range o.httpHeaders {
		headers[k] = v
	}
	headers["API-Ops-Datacenter-Id"] = *o.apiOpsDcId
	if in.apiInput != nil {
		headers["Content-Type"] = "application/json"
	}
	o.unlock(mutexKeyHttpHeaders)

	type proxyMessage struct {
		Method  string            `json:"method"`
		UrlPath string            `json:"urlPath"`
		Body    []byte            `json:"body,omitempty"`
		Params  map[string]string `json:"params,omitempty"`
		Headers map[string]string `json:"headers,omitempty"`
	}

	msg := proxyMessage{
		Method:  in.method,
		UrlPath: apstraUrl.Path,
		Params:  params,
		Headers: headers,
	}

	// are we sending data to the server?
	if in.apiInput != nil {
		msg.Body, err = json.Marshal(in.apiInput)
		if err != nil {
			return fmt.Errorf("error marshaling proxyMessage in talkToApiOps for url '%s' - %w", apstraUrl.String(), err)
		}
	}

	requestBody, err := json.Marshal(&struct {
		DcId         string       `json:"datacenter_edge_id"`
		ProxyMessage proxyMessage `json:"payload"`
	}{
		DcId:         *o.apiOpsDcId,
		ProxyMessage: msg,
	})
	if err != nil {
		return fmt.Errorf("error marshaling payload in talkToApiOps for url '%s' - %w", apstraUrl.String(), err)
	}

	// wrap supplied context with timeout (maybe)
	_, contextHasDeadline := ctx.Deadline()
	if !contextHasDeadline { // maybe this context already has a deadline?
		switch {
		case o.cfg.Timeout < 0: // negative Timeout is no timeout interval (infinite)
		case o.cfg.Timeout == 0: // Timeout of zero means use DefaultTimeout
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
			defer cancel()
		case o.cfg.Timeout > 0: // positive Timeout means use this value
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, o.cfg.Timeout)
			defer cancel()
		}
	}

	// create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseUrl.String()+aosOpsUrlPath, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("error creating http Request for url '%s' - %w", apstraUrl.String(), err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Dest-Fallback", "s3")
	req.Header.Set("X-Dest-Fallback-Region", "us-east-1")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do for url '%s' via api-ops proxy - %w", apstraUrl.String(), err)
	}

	// response not okay?
	if resp.StatusCode/100 != 2 {
		// noinspection GoUnhandledErrorResult
		defer resp.Body.Close()
		return newTalkToApstraErr(req, requestBody, resp, fmt.Sprintf("API-ops proxy response code: %d", resp.StatusCode))
	}

	// noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	var proxyResponse struct {
		HTTPStatusCode int    `json:"statusCode"`
		HTTPResponse   string `json:"response"`
		ErrorMsg       string `json:"errorMsg"`
	}
	err = json.NewDecoder(resp.Body).Decode(&proxyResponse)
	if err != nil {
		return fmt.Errorf("error decoding proxy response for url '%s' - %w", apstraUrl.String(), err)
	}
	if proxyResponse.ErrorMsg != "" {
		return newTalkToApstraErr(req, requestBody, resp, fmt.Sprintf("API-ops proxy error message: %s", proxyResponse.ErrorMsg))
	}

	// create a bogus http.Response so that our previously implemented logic works with it
	innerResp := new(http.Response)
	innerResp.Body = io.NopCloser(base64.NewDecoder(base64.StdEncoding, strings.NewReader(proxyResponse.HTTPResponse)))
	// noinspection GoUnhandledErrorResult
	defer innerResp.Body.Close()

	if proxyResponse.HTTPStatusCode/100 != 2 {
		return newTalkToApstraErr(req, requestBody, innerResp, fmt.Sprintf("API-ops proxy response code: %d, error message: %q", proxyResponse.HTTPStatusCode, proxyResponse.ErrorMsg))
	}

	// If the caller gave us an httpBodyWriter, copy the response body into it and return
	if in.httpBodyWriter != nil {
		_, err = io.CopyBuffer(in.httpBodyWriter, innerResp.Body, nil)
		if err != nil {
			return fmt.Errorf("error while reading http response body - %w", err)
		}
		return nil
	}

	// figure out whether Apstra responded with a task ID
	var tIdR taskIdResponse
	taskResponseFound, err := peekParseResponseBodyAsTaskId(innerResp, &tIdR)
	if err != nil {
		return newTalkToApstraErr(req, requestBody, innerResp, "error peeking response body")
	}

	// no task ID response, so no polling tomfoolery required
	if !taskResponseFound {
		if in.apiResponse == nil {
			o.Log(2, "no task ID response, and caller wants nothing back - talkToApstra done")
			// caller expects no response, so we're done here
			return nil
		}
		o.Log(2, "no task ID response, parse apstra reply for caller")
		// no task ID, decode response body into the caller-specified structure
		return json.NewDecoder(innerResp.Body).Decode(in.apiResponse)
	}

	// we got a task ID, instead of the expected response object. tasks are
	// per-blueprint, so we need to figure out the blueprint ID for task
	// progress query reasons.
	var bpId ObjectId

	// maybe the blueprintId is in the URL?
	if strings.Contains(apstraUrl.Path, apiUrlBlueprintsPrefix) {
		bpId = blueprintIdFromUrl(apstraUrl)
	}

	switch {
	case (bpId != "" && tIdR.BlueprintId != "") && (bpId != tIdR.BlueprintId):
		return fmt.Errorf("blueprint Id in URL ('%s') and returned object body ('%s') don't match", bpId, tIdR.BlueprintId)
	case bpId == "" && tIdR.BlueprintId == "":
		return newTalkToApstraErr(req, requestBody, innerResp, "blueprint id not found in url nor in response body")
	case bpId == "":
		bpId = tIdR.BlueprintId
	}
	o.Logf(2, "apstra returned task ID '%s' for blueprint '%s'", tIdR.TaskId, tIdR.BlueprintId)

	// get (wait for) full detailed response on the outstanding task ID
	taskResponse, err := waitForTaskCompletion(bpId, tIdR.TaskId, o.taskMonChan)
	if err != nil {
		return fmt.Errorf("error in task monitor - %w", err)
	}

	// there might be errors articulated in the taskResponse body
	if len(taskResponse.DetailedStatus.Errors) > 0 || taskResponse.DetailedStatus.ErrorCode != 0 {
		originalUrl, _ := url.Parse(taskResponse.RequestData.Url)
		qValues := originalUrl.Query()
		for k, v := range taskResponse.RequestData.Args {
			qValues.Add(k, v)
		}
		originalUrl.RawQuery = qValues.Encode()

		originalHdr := make(http.Header, len(taskResponse.RequestData.Headers))
		for k, v := range taskResponse.RequestData.Headers {
			originalHdr.Add(k, v)
		}

		var originalBody bytes.Buffer
		originalBody.Write(taskResponse.RequestData.Data)

		request := &http.Request{
			Method:        taskResponse.RequestData.Method,
			URL:           originalUrl,
			Header:        originalHdr,
			Body:          io.NopCloser(&originalBody),
			ContentLength: int64(len(taskResponse.RequestData.Data)),
		}

		var responseBody bytes.Buffer
		responseBody.Write(taskResponse.DetailedStatus.Errors)

		response := &http.Response{
			StatusCode:    taskResponse.DetailedStatus.ErrorCode,
			Body:          io.NopCloser(&responseBody),
			ContentLength: int64(len(taskResponse.DetailedStatus.Errors)),
		}

		dsMsg, _ := json.Marshal(&taskResponse.DetailedStatus)

		return TalkToApstraErr{
			Request:  request,
			Response: response,
			Msg:      string(dsMsg),
		}
	}

	// caller not expecting any response?
	if in.apiResponse == nil {
		return nil
	}

	// the getTaskResponse data structure is only partially unmarshaled because
	// it's impossible to know exactly what'll be in there. Extract it now into
	// whatever in.apiResponse (interface{} controlled by the caller) is.
	return json.Unmarshal(taskResponse.DetailedStatus.ApiResponse, in.apiResponse)

}

//func (o *Client) talkToApiOpsOld(apstraReq *http.Request) (*http.Response, error) {
//	// Based on:
//	// https://ssd-git.juniper.net/aide-jcloud/apstra-marvis/api-ops/-/blob/develop/proto/aospayload.pb.go#L106-125
//	type proxyMessage struct {
//		Method  string            `json:"method"`
//		UrlPath string            `json:"urlPath"`
//		Body    []byte            `json:"body,omitempty"`
//		Params  map[string]string `json:"params,omitempty"`
//		Headers map[string]string `json:"headers,omitempty"`
//	}
//
//	type proxyBody struct {
//		DcId         string       `json:"datacenter_edge_id"`
//		ProxyPayload proxyMessage `json:"payload"`
//	}
//
//	origBody, err := io.ReadAll(apstraReq.Body)
//	if err != nil {
//		return nil, fmt.Errorf("error prepping proxy request body from apstraReq: %w", err)
//	}
//
//	params := apstraReq.URL.Query()
//
//	body := proxyBody{
//		DcId: *o.apiOpsDcId,
//		ProxyPayload: proxyMessage{
//			Method:  apstraReq.Method,
//			UrlPath: apstraReq.URL.Path,
//			Body:    origBody,
//			Params:  make(map[string]string, len(params)),
//			Headers: make(map[string]string, len(apstraReq.Header)),
//		},
//	}
//
//	body.ProxyPayload.Params = make(map[string]string, len(params))
//	for k, v := range params {
//		switch len(v) {
//		case 0:
//		case 1:
//			body.ProxyPayload.Params[k] = v[0]
//		default:
//			return nil, fmt.Errorf("cannot format query string param %q for the proxy: only one string supported per param, got %d strings", k, len(v))
//		}
//	}
//
//	body.ProxyPayload.Headers = make(map[string]string, len(apstraReq.Header)+1)
//	for k, v := range apstraReq.Header {
//		switch len(v) {
//		case 0:
//		case 1:
//			body.ProxyPayload.Headers[k] = v[0]
//		default:
//			return nil, fmt.Errorf("cannot format query string header %q for the proxy: only one string supported per header, got %d strings", k, len(v))
//		}
//	}
//	body.ProxyPayload.Headers["X-Dest-Fallback"] = "s3"
//
//	bodyBytes, err := json.Marshal(body)
//	if err != nil {
//		return nil, fmt.Errorf("error prepping proxy request body for api-ops proxy: %w", err)
//	}
//
//	proxyRequest, err := http.NewRequest(http.MethodPost, o.baseUrl.String()+aosOpsUrlPath, bytes.NewReader(bodyBytes))
//	if err != nil {
//		return nil, fmt.Errorf("error prepping proxyRequest from apstraReq: %w", err)
//	}
//
//	proxyRequest = proxyRequest.WithContext(apstraReq.Context())
//
//	return o.httpClient.Do(proxyRequest)
//}
