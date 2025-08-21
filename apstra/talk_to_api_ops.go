// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"compress/gzip"
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
	envAosOpsEdgeId = "API_OPS_DATACENTER_EDGE_ID"
	envAosOpsNoGzip = "API_OPS_DISABLE_GZIP"
)

func (o *Client) talkToApiOps(ctx context.Context, in *talkToApstraIn) error {
	// create URL we'd use if we were talking to an actual Apstra server
	apstraUrl, err := o.craftUrl(in)
	if err != nil {
		return err
	}

	// we need to pack the query string params into proxy-bound http payload
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
	headers := make(map[string]string, len(o.httpHeaders)+3)
	for k, v := range o.httpHeaders {
		headers[k] = v
	}
	if !o.skipGzip {
		headers["Accept-Encoding"] = "gzip"
	}
	if in.apiInput != nil {
		headers["Content-Type"] = "application/json"
	}
	headers["X-Dest-Fallback"] = "s3"
	o.unlock(mutexKeyHttpHeaders)

	// this is the structure we send to the proxy
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
		DcId:         *o.cfg.apiOpsDcId,
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

	// create request to send to the proxy
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseUrl.String()+aosOpsUrlPath, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("error creating http Request for url '%s' - %w", apstraUrl.String(), err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return fmt.Errorf("error calling http.client.Do for url '%s' via api-ops proxy - %w", apstraUrl.String(), err)
	}

	defer func() { _ = resp.Body.Close() }() // close the response body received directly from the proxy

	// proxy's status code not okay?
	if resp.StatusCode/100 != 2 {
		req.URL = apstraUrl
		req.URL.Host = "api-ops"
		return newTalkToApstraErr(req, requestBody, resp, fmt.Sprintf("API-ops proxy response code: %d", resp.StatusCode))
	}

	var proxyResponse struct {
		Uid        string            `json:"uid"`
		Headers    map[string]string `json:"headers"`
		StatusCode int               `json:"statusCode"`
		Response   string            `json:"response"`
		ErrorMsg   string            `json:"errorMsg"`
	}
	err = json.NewDecoder(resp.Body).Decode(&proxyResponse)
	if err != nil {
		return fmt.Errorf("decoding proxy response for url '%s': %w", apstraUrl.String(), err)
	}
	if proxyResponse.ErrorMsg != "" {
		req.URL = apstraUrl
		req.URL.Host = "api-ops"
		return newTalkToApstraErr(req, requestBody, resp, fmt.Sprintf("API-ops proxy error message for transaction %s: %s", proxyResponse.Uid, proxyResponse.ErrorMsg))
	}
	// we'll check the proxied status code (the one from AOS) later, after we've assembled a stunt-double http response

	var gzResponse bool
	if ce, ok := proxyResponse.Headers["Content-Encoding"]; ok {
		if strings.Contains(ce, "gzip") {
			gzResponse = true
		}
	}

	var s3Response bool
	if xDestResp, ok := proxyResponse.Headers["X-Dest"]; ok {
		if xDestResp == "S3" {
			s3Response = true
		} else {
			return fmt.Errorf("unexpected X-Dest header: %s", xDestResp)
		}
	}

	var responseS3Url string
	if s3Response {
		b, err := base64.StdEncoding.DecodeString(proxyResponse.Response)
		if err != nil {
			return fmt.Errorf("b64 decoding s3Url for transaction %q: %w", proxyResponse.Uid, err)
		}

		responseS3Url = string(b)
	}

	// create a bogus http.Response so that our previously implemented logic works with it
	innerResp := new(http.Response)
	innerResp.StatusCode = proxyResponse.StatusCode

	defer func() { _ = innerResp.Body.Close() }()

	// Apstra may have responded with compressed data (gzResponse) and api-ops may have
	// offloaded our desired payload to S3 (s3Response). Each combination of these
	// possibilities needs to be handled slightly differently.
	switch {
	case !gzResponse && !s3Response:
		innerResp.Body = io.NopCloser(base64.NewDecoder(base64.StdEncoding, strings.NewReader(proxyResponse.Response)))

	case !gzResponse && s3Response:
		s3Resp, err := getUrl(ctx, responseS3Url, o.httpClient)
		if err != nil {
			return fmt.Errorf("fetching response from S3 URL %q: %w", responseS3Url, err)
		}

		innerResp.Body = s3Resp.Body

	case gzResponse && !s3Response:
		gzReader, err := gzip.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(proxyResponse.Response)))
		if err != nil {
			return fmt.Errorf("error creating gzip reader for transaction %s - %w", proxyResponse.Uid, err)
		}

		innerResp.Body = gzReader

	case gzResponse && s3Response:
		s3Resp, err := getUrl(ctx, responseS3Url, o.httpClient)
		if err != nil {
			return fmt.Errorf("fetching response from S3 URL %q: %w", responseS3Url, err)
		}

		r, err := gzip.NewReader(s3Resp.Body)
		if err != nil {
			return fmt.Errorf("creating gz reader for S3 response body: %w", err)
		}

		defer func() { _ = s3Resp.Body.Close() }()

		innerResp.Body = r
	}

	if innerResp.StatusCode/100 != 2 {
		req.URL = apstraUrl
		req.URL.Host = "api-ops"
		return newTalkToApstraErr(req, msg.Body, innerResp, "")
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
		req.URL = apstraUrl
		req.URL.Host = "api-ops"
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
		req.URL = apstraUrl
		req.URL.Host = "api-ops"
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

func getUrl(ctx context.Context, urlStr string, client apstraHttpClient) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request for %q: %w", urlStr, err)
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("perofrming GET for %q: %w", urlStr, err)
	}
	if resp.StatusCode/100 != 2 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("http response to GET for %q: %d", urlStr, resp.StatusCode)
	}

	return resp, nil
}
