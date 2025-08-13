// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// TalkToApstraErr implements error{} and carries around http.Request and
// http.Response object pointers. Error() method produces a string like
// "<error> - http response <status> at url <url>".
type TalkToApstraErr struct {
	Request  *http.Request
	Response *http.Response
	Msg      string
}

func (o TalkToApstraErr) Error() string {
	apstraUrl := "nil"
	if o.Request != nil {
		apstraUrl = o.Request.URL.String()
	}

	status := "nil"
	if o.Response != nil {
		status = o.Response.Status
	}

	return fmt.Sprintf("%s - http response '%s' at '%s'", o.Msg, status, apstraUrl)
}

func (o TalkToApstraErr) parseApiUrlBlueprintObjPolicyBatchApplyError() error {
	body, err := io.ReadAll(o.Response.Body)
	if err != nil {
		return fmt.Errorf("reading error response body: %w", err)
	}

	var raw struct {
		ApplicationPoints map[string]struct {
			Policies map[string]json.RawMessage `json:"policies"` // these might be strings or a structs
		} `json:"application_points"`
	}

	if json.Unmarshal(body, &raw) != nil {
		return fmt.Errorf("parsing obj-policy-batch-apply error: %w", o)
	}

	var detail ErrCtAssignmentFailedDetail // we return this

	var rawPolicyMsg struct { // we unpack each policy element here
		Policy string `json:"policy"`
	}

	policyIdRegexp := regexp.MustCompile("^Endpoint policy with node id (.*) does not exist$")

	// store collected indexes and ids in maps for de-dup reasons
	invalidApplicationPointIndexes := make(map[int]struct{})
	invalidConnectivityTemplateIds := make(map[ObjectId]struct{})

	for apKey, ap := range raw.ApplicationPoints {
		apIdx, err := strconv.Atoi(apKey)
		if err != nil {
			return fmt.Errorf("parsing obj-policy-batch-apply application point map index %q: %w': parsing obj-policy-batch-apply error: %w", apKey, err, o)
		}

		for pKey, p := range ap.Policies {
			err = json.Unmarshal(p, &rawPolicyMsg) // maybe we got a struct?
			if err != nil {
				err = json.Unmarshal(p, &rawPolicyMsg.Policy) // maybe we got a string?
				if err != nil {
					return fmt.Errorf("cannot parse error at application point %q, policy %q: %w", apKey, pKey, o) // don't wrap either err here
				}
			}

			policyIdSubMatches := policyIdRegexp.FindStringSubmatch(rawPolicyMsg.Policy)

			switch {
			case len(policyIdSubMatches) == 2:
				invalidConnectivityTemplateIds[ObjectId(policyIdSubMatches[1])] = struct{}{}
				//detail.InvalidConnectivityTemplateIds = append(detail.InvalidConnectivityTemplateIds, ObjectId(policyIdSubMatches[1]))
			case rawPolicyMsg.Policy == "Not a valid application point":
				invalidApplicationPointIndexes[apIdx] = struct{}{}
			default:
				return fmt.Errorf("cannot parse error at application point %q, policy %q: %w", apKey, pKey, o)
			}
		}
	}

	for invalidApplicationPointIdx := range invalidApplicationPointIndexes {
		detail.InvalidApplicationPointIndexes = append(detail.InvalidApplicationPointIndexes, invalidApplicationPointIdx)
	}
	for invalidConnectivityTemplateId := range invalidConnectivityTemplateIds {
		detail.InvalidConnectivityTemplateIds = append(detail.InvalidConnectivityTemplateIds, invalidConnectivityTemplateId)
	}

	return ClientErr{
		errType:   ErrCtAssignmentFailed,
		err:       fmt.Errorf("assigning connectivity templates: %s", o),
		detail:    &detail,
		retryable: false,
	}
}

// newTalkToApstraErr returns a TalkToApstraErr. It's intended to be called after the
// http.Request has been executed with Do(), so the request body has already
// been "spent" by Read(). We'll fill it back in. The response body is likely to
// be closed by a 'defer body.Close()' somewhere, so we'll replace that as well,
// up to some reasonable limit (don't try to buffer gigabytes of data from the
// webserver).
func newTalkToApstraErr(req *http.Request, reqBody []byte, resp *http.Response, errMsg string) TalkToApstraErr {
	// don't include secret in error
	req.Header.Del(apstraAuthHeader)

	// redact request body for sensitive URLs
	switch req.URL.Path {
	case apiUrlUserLogin:
		req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("request body for '%s' redacted", req.URL.Path)))
	default:
		rehydratedRequest := bytes.NewBuffer(reqBody)
		req.Body = io.NopCloser(rehydratedRequest)
	}

	// redact response body for sensitive URLs
	switch req.URL.Path {
	case apiUrlUserLogin:
		_ = resp.Body.Close() // close the real network socket
		resp.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("resposne body for '%s' redacted", req.URL.Path)))
	default:
		// prepare a stunt double response body for the one that's likely attached to a network
		// socket, and likely to be closed by a `defer` somewhere
		rehydratedResponse := &bytes.Buffer{}
		_, _ = io.CopyN(rehydratedResponse, resp.Body, errResponseBodyLimit) // size limit
		resp.Body = io.NopCloser(rehydratedResponse)                         // replace the original body
	}

	// use first part of response body if errMsg empty
	if errMsg == "" {
		peekAbleBodyReader := bufio.NewReader(resp.Body)
		resp.Body = io.NopCloser(peekAbleBodyReader)
		peek, _ := peekAbleBodyReader.Peek(errResponseStringLimit)
		errMsg = string(peek)
	}

	return TalkToApstraErr{
		Request:  req,
		Response: resp,
		Msg:      errMsg,
	}
}
