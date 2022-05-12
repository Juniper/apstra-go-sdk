package apstra

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
)

const (
	apstraApiAsyncParamKey          = "async"
	apstraApiAsyncParamValFull      = "full"
	apstraApiAsyncParamValPartial   = "partial" // default?
	errResponseLimit                = 4096
	peekSizeForApstraTaskIdResponse = math.MaxUint8
)

// talkToApstraIn is the input structure for the Client.talkToApstra() function
type talkToApstraIn struct {
	method         string      // how to talk to Apstra
	url            *url.URL    // where to talk to Aptstra (as little as /path/to/thing ok)
	apiInput       interface{} // if non-nil we'll JSON encode this prior to sending it
	apiResponse    interface{} // if non-nil we'll JSON decode Apstra response here
	doNotLogin     bool        // when set, Client will not attempt login (we set for anti-recursion)
	unsynchronized bool        // default behavior is to send apstraApiAsyncParamValFull, block until task completion
}

// craftUrl combines o.baseUrl (probably "http://host:port") with in.url
// (probably "/api/something/something", might have a query string).
// The assumption is that o.baseUrl contains the scheme, host (host+port) and
// leading path components, while `in` (talkToApstraIn) is responsible for the
// path to the specific API endpoint and any required query parameters.
// When `in.unsychronized` is false (the default), Apstra's 'async=full' query
// string parameter is added to the returned result.
func (o Client) craftUrl(in *talkToApstraIn) *url.URL {
	result := in.url
	if result.Scheme == "" {
		result.Scheme = o.baseUrl.Scheme // copy baseUrl scheme
	}
	if result.Host == "" {
		result.Host = o.baseUrl.Host // copy baseUrl host
	}

	result.Path = o.baseUrl.Path + in.url.Path // path is cumulative, baseUrl can be empty

	if !in.unsynchronized {
		params := result.Query()
		params.Set(apstraApiAsyncParamKey, apstraApiAsyncParamValFull)
		result.RawQuery = params.Encode()
	}

	return result
}

// talkToApstra talks to the Apstra server using in.method. If in.apiInput is
// not nil, it JSON-encodes that data structure and sends it. In case the
// in.apiResponse is not nil, the server response is extracted into it.
func (o Client) talkToApstra(ctx context.Context, in *talkToApstraIn) error {
	var err error
	var requestBody []byte
	if ctx == nil {
		ctx = context.TODO()
	}

	// create URL
	apstraUrl := o.craftUrl(in)

	// are we sending data to the server?
	if in.apiInput != nil {
		requestBody, err = json.Marshal(in.apiInput)
		if err != nil {
			return fmt.Errorf("error marshaling payload in talkToApstra for url '%s' - %w", in.url, err)
		}
	}

	// wrap supplied context with timeout (maybe)
	_, contextHasDeadline := ctx.Deadline()
	if o.cfg.Timeout != 0 && !contextHasDeadline {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, o.cfg.Timeout)
		defer cancel()
	}

	// create request
	req, err := http.NewRequestWithContext(ctx, string(in.method), apstraUrl.String(), bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("error creating http Request for url '%s' - %w", in.url, err)
	}

	// set request httpHeaders
	if in.apiInput != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range o.httpHeaders {
		req.Header.Set(k, v)
	}

	// talk to the server
	resp, err := o.httpClient.Do(req)

	// trim authentication token from request - Do() has been called - get this out of the way quickly
	req.Header.Del(apstraAuthHeader)
	if err != nil { // check error from req.Do()
		return fmt.Errorf("error calling http.client.Do for url '%s' - %w", in.url, err)
	}

	// noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// response not okay?
	if resp.StatusCode/100 != 2 {

		// Auth fail?
		if resp.StatusCode == 401 {
			// Auth fail at login API is fatal for this transaction
			if in.url.Path == apiUrlUserLogin {
				return newTalkToApstraErr(req, requestBody, resp,
					fmt.Sprintf("http %d at '%s' - check username/password",
						resp.StatusCode, in.url))
			}

			// Auth fail with "doNotLogin == true" is fatal for this transaction
			if in.doNotLogin {
				return newTalkToApstraErr(req, requestBody, resp,
					fmt.Sprintf("http %d at '%s' and doNotLogin is %t",
						resp.StatusCode, in.url, in.doNotLogin))
			}

			// Try logging in
			err := o.login(ctx)
			if err != nil {
				return fmt.Errorf("error attempting login after initial AuthFail - %w", err)
			}

			// Try the request again
			in.doNotLogin = true
			return o.talkToApstra(ctx, in)
		} // HTTP 401

		return newTalkToApstraErr(req, requestBody, resp, "")
	}

	// caller not expecting any response?
	if in.apiResponse == nil {
		return nil
	}

	// caller is expecting a response, but we don't know if Apstra will return
	// the desired data structure, or a taskIdResponse.
	var tIdR taskIdResponse
	ok, err := peekParseResponseBodyAsTaskId(resp, &tIdR)
	if err != nil {
		return newTalkToApstraErr(req, requestBody, resp, "")
	}
	if !ok {
		// no task ID, decode response body into the caller-specified structure
		return json.NewDecoder(resp.Body).Decode(in.apiResponse)
	}

	// we got a task ID, instead of the expected response object
	bpId, err := blueprintIdFromUrl(apstraUrl)
	if err != nil {
		return fmt.Errorf("error parsing blueprint ID from URL '%s' - %w", apstraUrl.String(), err)
	}

	// get (wait for) full detailed response on the outstanding task ID
	taskResponse, err := waitForTaskCompletion(bpId, tIdR.TaskId, o.taskMonChan)
	if err != nil {
		return fmt.Errorf("error in task monitor - %w\n API result:\n", err)
	}

	// the getTaskResponse data structure is only partially unmarshaled because
	// it's impossible to know exactly what'll be in there. Extract it now into
	// whatever in.apiResponse (interface{} controlled by the caller) is.
	return json.Unmarshal(taskResponse.DetailedStatus.ApiResponse, in.apiResponse)

}

// talkToApstraErr implements error{} and carries around http.Request and
// http.Response object pointers. Error() method produces a string like
// "<error> - http response <status> at url <url>".
// todo: methods like ErrorCRIT() and ErrorWARN()
type talkToApstraErr struct {
	request  *http.Request
	response *http.Response
	error    string
}

func (o talkToApstraErr) Error() string {
	apstraUrl := "nil"
	if o.request != nil {
		apstraUrl = o.request.URL.String()
	}

	status := "nil"
	if o.response != nil {
		status = o.response.Status
	}

	return fmt.Sprintf("%s - http response '%s' at '%s'", o.error, status, apstraUrl)
}

// newTalkToApstraErr returns a talkToApstraErr. It's intended to be called after the
// http.Request has been executed with Do(), so the request body has already
// been "spent" by Read(). We'll fill it back in. The response body is likely to
// be closed by a 'defer body.Close()' somewhere, so we'll replace that as well,
// up to some reasonable limit (don't try to buffer gigabytes of data from the
// webserver).
func newTalkToApstraErr(req *http.Request, reqBody []byte, resp *http.Response, errMsg string) talkToApstraErr {
	apstraUrl := req.URL.String()
	// don't include secret in error
	req.Header.Del(apstraAuthHeader)

	// redact request body for sensitive URLs
	switch apstraUrl {
	case apiUrlUserLogin:
		req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("request body for '%s' redacted", apstraUrl)))
	default:
		rehydratedRequest := bytes.NewBuffer(reqBody)
		req.Body = io.NopCloser(rehydratedRequest)
	}

	// redact response body for sensitive URLs
	switch apstraUrl {
	case apiUrlUserLogin:
		_ = resp.Body.Close() // close the real network socket
		resp.Body = io.NopCloser(strings.NewReader(fmt.Sprintf("resposne body for '%s' redacted", apstraUrl)))
	default:
		// prepare a stunt double response body for the one that's likely attached to a network
		// socket, and likely to be closed by a `defer` somewhere
		rehydratedResponse := &bytes.Buffer{}
		_, _ = io.CopyN(rehydratedResponse, resp.Body, errResponseLimit) // size limit
		_ = resp.Body.Close()                                            // close the real network socket
		resp.Body = io.NopCloser(rehydratedResponse)                     // replace the original body
	}

	return talkToApstraErr{
		request:  req,
		response: resp,
		error:    errMsg,
	}
}

func peekParseResponseBodyAsTaskId(resp *http.Response, result *taskIdResponse) (bool, error) {
	peekAbleBodyReader := bufio.NewReader(resp.Body)
	resp.Body = io.NopCloser(peekAbleBodyReader)
	peek, err := peekAbleBodyReader.Peek(peekSizeForApstraTaskIdResponse)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("error peeking into http response body - %w", err)
	}
	err = json.Unmarshal(peek, result)
	// wild assumption:
	//   Every error means "peek data doesn't look like a taskIdResponse".
	//   There is no error which indicates a problem of any other type.
	if err != nil {
		return false, nil // no error; 'false' b/c unmarshal TaskId failed
	} else { // good unmarshal, but what about the contents?
		return result.TaskId != "", nil // no error; bool depends on string match
	}
}
