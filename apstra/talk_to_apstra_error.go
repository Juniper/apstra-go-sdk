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
	var err error

	// unpack the request we sent
	var rawReq struct {
		ApplicationPoints []struct {
			Id ObjectId `json:"id"`
			//Policies []struct {
			//	Policy ObjectId `json:"policy"`
			//	Used   bool     `json:"used"`
			//} `json:"policies"`
		} `json:"application_points"`
	}
	err = json.NewDecoder(o.Request.Body).Decode(&rawReq)
	if err != nil {
		// not expected to error, so the returned error is wrapped
		return fmt.Errorf("reading/decoding request body: %w parent error: %w", err, o)
	}

	// unpack the response we received, if possible
	var rawResp struct {
		ApplicationPoints map[int]struct {
			Policies map[int]json.RawMessage `json:"policies"` // these might be strings or a structs
		} `json:"application_points"`
	}
	err = json.NewDecoder(o.Response.Body).Decode(&rawResp)
	if err != nil {
		// Do not wrap the error. Wrapping is not a value-add when we know
		// there are error types which we do not handle. For example:
		//
		// {
		//  "error_code": 422,
		//  "errors": {
		//    "application_points": {
		//      "LKbJ0lJjsTmaPucko48": { <<<=================== Application point ID is string, not int
		//        "71d0aeaf-b686-11f0-8849-617073747261": { <<==== CT ID
		//          "71d0af10-b686-11f0-8849-617073747261": { <<== primitive ID (possibly nested deeper in some cases?)
		//            "conflicts": [
		//              {
		//                "policy": "d1bfe22e-2f5e-11ef-ac36-617073747261", <<<== conflicting CT ID
		//                "error": "Unable to apply Virtual Network template. VN my_VN is already configured on interface of system my_server"
		//              }
		//            ]
		//          }
		//        }
		//      }
		//    }
		//  }
		//}
		return o // return the raw error since we're not going to parse it
	}

	var rawPolicyStruct struct { // we attempt to unpack each policy element here
		Policy string `json:"policy"`
	}

	policyIdErrorRegex := regexp.MustCompile("^Endpoint policy with node id (.*) does not exist$")

	// store collected indexes and ids in maps for de-dup reasons
	invalidApIds := make(map[ObjectId]struct{})
	invalidCTIds := make(map[ObjectId]struct{})

	// loop over the error's application point map
	for apIdx, ap := range rawResp.ApplicationPoints {
		for pKey, p := range ap.Policies {
			err = json.Unmarshal(p, &rawPolicyStruct) // maybe we got a struct?
			if err != nil {
				err = json.Unmarshal(p, &rawPolicyStruct.Policy) // maybe we got a string?
				if err != nil {
					return fmt.Errorf("parsing error at application point %d, policy %q: %w", apIdx, pKey, o) // don't wrap either err here
				}
			}

			policyIdSubMatches := policyIdErrorRegex.FindStringSubmatch(rawPolicyStruct.Policy)

			switch {
			case len(policyIdSubMatches) == 2:
				invalidCTIds[ObjectId(policyIdSubMatches[1])] = struct{}{}
			case rawPolicyStruct.Policy == "Not a valid application point":
				if apIdx < 0 || apIdx >= len(rawReq.ApplicationPoints) {
					return fmt.Errorf("invalid application point index %d in API response", apIdx)
				}
				invalidApIds[rawReq.ApplicationPoints[apIdx].Id] = struct{}{}
			default:
				return fmt.Errorf("cannot parse error at application point %d, policy %q: %w", apIdx, pKey, o)
			}
		}
	}

	var detail ErrCtAssignmentFailedDetail // we return this
	for invalidApId := range invalidApIds {
		detail.InvalidApplicationPointIds = append(detail.InvalidApplicationPointIds, invalidApId)
	}
	for invalidCtId := range invalidCTIds {
		detail.InvalidConnectivityTemplateIds = append(detail.InvalidConnectivityTemplateIds, invalidCtId)
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
