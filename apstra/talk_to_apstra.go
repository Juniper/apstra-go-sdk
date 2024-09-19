package apstra

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	CtxKeyTestID   = "Test-ID"   // context.Context key for a test ID string
	CtxKeyTestUUID = "Test-UUID" // context.Context key for a uuid.UUID upon which the test ID string is based

	apstraApiAsyncParamKey          = "async"
	apstraApiAsyncParamValFull      = "full"
	apstraApiAsyncParamValPartial   = "partial" // default?
	apstraApiUnsafeParamKey         = "allow_unsafe"
	errResponseBodyLimit            = 4096
	errResponseStringLimit          = 1024
	peekSizeForApstraTaskIdResponse = math.MaxUint8

	linkHasCtAssignedErrRegexString     = "Link with id (.*) can not be deleted since some of its interfaces have connectivity templates assigned"
	lagHasCtAssignedErrRegexString      = "Deleting all links forming a LAG is not allowed since the LAG has assigned structures: \\[.*'connectivity template'.*]. Link ids: \\[(.*)]"
	linkHasVnEndpointErrRegexString     = "Link with id (.*) can not be deleted since some of its interfaces have VN endpoints"
	linkHasSubinterfacesErrRegexString  = "Link with id (.*) can not be deleted since some of its interfaces have subinterfaces"
	lagHasAssignedStructuresRegexString = "Operation is not permitted because link group (.*) has assigned structures"
)

var (
	regexpApiUrlDeleteSwitchSystemLinks = regexp.MustCompile(strings.ReplaceAll(apiUrlDeleteSwitchSystemLinks, "%s", ".*"))
	regexpApiUrlLeafServerLinkLabels    = regexp.MustCompile(strings.ReplaceAll(apiUrlLeafServerLinkLabels, "%s", ".*"))
	regexpLinkHasCtAssignedErr          = regexp.MustCompile(linkHasCtAssignedErrRegexString)
	regexpLagHasCtAssignedErr           = regexp.MustCompile(lagHasCtAssignedErrRegexString)
	regexpLinkHasVnEndpoint             = regexp.MustCompile(linkHasVnEndpointErrRegexString)
	regexpLinkHasSubinterfaces          = regexp.MustCompile(linkHasSubinterfacesErrRegexString)
	regexpLagHasAssignedStructures      = regexp.MustCompile(lagHasAssignedStructuresRegexString)
)

// talkToApstraIn is the input structure for the Client.talkToApstra() function
type talkToApstraIn struct {
	method         string      // how to talk to Apstra
	url            *url.URL    // where to talk to Aptstra (as little as /path/to/thing ok) this is considered before urlStr
	urlStr         string      // where to talk to Apstra, this one is used if url is nil
	apiInput       interface{} // if non-nil we'll JSON encode this prior to sending it
	apiResponse    interface{} // if non-nil we'll JSON decode Apstra response here
	doNotLogin     bool        // when set, Client will not attempt login (we set for anti-recursion)
	unsynchronized bool        // default behavior is to send apstraApiAsyncParamValFull, block until task completion
	httpBodyWriter io.Writer   // when non-nil, http body will be written here instead of unpacked into apiResponse
	unsafe         bool        // when true, set allow_unsafe=true HTTP query string parameter
}

type apstraErr struct {
	Errors string `json:"errors"`
}

// taskIdResponse data structure is returned by Apstra for *some* operations, when the
// URL Query String includes `async=full`
type taskIdResponse struct {
	BlueprintId ObjectId `json:"id"`
	TaskId      TaskId   `json:"task_id"`
}

func convertTtaeToAceWherePossible(err error) error {
	var ttae TalkToApstraErr
	if errors.As(err, &ttae) {
		switch ttae.Response.StatusCode {
		case http.StatusNotFound:
			if ttae.Request.URL.Path == apiUrlBlueprints {
				return ClientErr{errType: ErrNotfound, retryable: true, err: err}
			}
			return ClientErr{errType: ErrNotfound, err: err}
		case http.StatusConflict:
			return ClientErr{errType: ErrConflict, err: errors.New(ttae.Msg)}
		case http.StatusUnprocessableEntity:
			switch {
			case strings.Contains(ttae.Msg, "Direct graph modification operation is unsafe") &&
				strings.Contains(ttae.Msg, "If you want to proceed with this PATCH API call"):
				return ClientErr{errType: ErrUnsafePatchProhibited, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "No value in either user config or profile"):
				return ClientErr{errType: ErrAgentProfilePlatformRequired, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "already exists"):
				return ClientErr{errType: ErrExists, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "No node with id: "):
				return ClientErr{errType: ErrNotfound, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "No virtual_network with id: "):
				return ClientErr{errType: ErrNotfound, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "Virtual Network name not unique"):
				return ClientErr{errType: ErrExists, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "Transformation cannot be changed"):
				return ClientErr{errType: ErrCannotChangeTransform, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "does not exist"):
				return ClientErr{errType: ErrNotfound, err: errors.New(ttae.Msg)}
			case regexpApiUrlDeleteSwitchSystemLinks.MatchString(ttae.Request.URL.Path):
				switch {
				case regexpLinkHasCtAssignedErr.MatchString(ttae.Msg):
					return ClientErr{errType: ErrCtAssignedToLink, err: errors.New(ttae.Msg)}
				case regexpLagHasCtAssignedErr.MatchString(ttae.Msg):
					return ClientErr{errType: ErrCtAssignedToLink, err: errors.New(ttae.Msg)}
				}
			case regexpApiUrlLeafServerLinkLabels.MatchString(ttae.Request.URL.Path):
				return ClientErr{errType: ErrLagHasAssignedStructrues, err: errors.New(ttae.Msg)}
			}
		case http.StatusInternalServerError:
			switch {
			case strings.Contains(ttae.Msg, "Error executing facade API GET /obj-policy-export") &&
				strings.Contains(ttae.Msg, "'NoneType' object has no attribute 'id'"):
				return ClientErr{errType: ErrNotfound, err: errors.New(ttae.Msg)}
			case strings.Contains(ttae.Msg, "The current mount is conflicting with an existing mount"):
				return ClientErr{errType: ErrIbaCurrentMountConflictsWithExistingMount, retryable: true, err: errors.New(ttae.Msg)}
			}
		}
	}
	return err
}

// craftUrl combines o.baseUrl (probably "http://host:port") with in.url
// (preferred), or in.urlStr (if in.url is nil). Both options are probably
// something like "/api/something/something". More complicated callers (which
// need to use a query string, etc.) will probably send in.url (*url.URL), while
// simple ones can simply send in.urlStr (string).
// The assumption is that o.baseUrl contains the scheme, host (host+port) and
// leading path components, while `in` (talkToApstraIn) is responsible for the
// path to the specific API endpoint and any required query parameters.
// When `in.unsychronized` is false (the default), Apstra's 'async=full' query
// string parameter is added to the returned result.
func (o *Client) craftUrl(in *talkToApstraIn) (*url.URL, error) {
	result := in.url
	var err error

	if in.url == nil {
		result, err = url.Parse(in.urlStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing url '%s' - %w", in.urlStr, err)
		}
	}

	if result.Scheme == "" {
		result.Scheme = o.baseUrl.Scheme // copy baseUrl scheme
	}

	if result.Host == "" {
		result.Host = o.baseUrl.Host // copy baseUrl host
	}

	result.Path = o.baseUrl.Path + result.Path // path is cumulative, baseUrl can be empty

	// set query string parameters
	params := result.Query()
	var paramsChanged bool

	if !in.unsynchronized {
		params.Set(apstraApiAsyncParamKey, apstraApiAsyncParamValFull)
		paramsChanged = true
	}

	if in.unsafe {
		params.Set(apstraApiUnsafeParamKey, strconv.FormatBool(in.unsafe))
		paramsChanged = true
	}

	if paramsChanged {
		result.RawQuery = params.Encode()
	}

	return result, nil
}

// talkToApstra talks to the Apstra server using in.method. If in.apiInput is
// not nil, it JSON-encodes that data structure and sends it. In case the
// in.apiResponse is not nil, the server response is extracted into it.
func (o *Client) talkToApstra(ctx context.Context, in *talkToApstraIn) error {
	var err error
	var requestBody []byte

	// create URL
	apstraUrl, err := o.craftUrl(in)
	if err != nil {
		return err
	}

	// are we sending data to the server?
	if in.apiInput != nil {
		requestBody, err = json.Marshal(in.apiInput)
		if err != nil {
			return fmt.Errorf("error marshaling payload in talkToApstra for url '%s' - %w", apstraUrl.String(), err)
		}
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
	req, err := http.NewRequestWithContext(ctx, in.method, apstraUrl.String(), bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("error creating http Request for url '%s' - %w", apstraUrl.String(), err)
	}

	// set the Content-Type request header if we're sending any payload
	if in.apiInput != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// set the Test-ID header if the context has one.
	switch testId := ctx.Value(CtxKeyTestID).(type) {
	case string:
		req.Header.Set(CtxKeyTestID, testId)
	}

	o.lock(mutexKeyHttpHeaders)
	for k, v := range o.httpHeaders {
		req.Header.Set(k, v)
	}
	o.unlock(mutexKeyHttpHeaders)

	o.logFunc(2, o.dumpHttpRequest, req)

	// talk to the server
	resp, err := o.httpClient.Do(req)

	// trim authentication token from request - Do() has been called - get this out of the way quickly
	req.Header.Del(apstraAuthHeader)
	if err != nil { // check error from req.Do()
		return fmt.Errorf("error calling http.client.Do for url '%s' - %w", apstraUrl.String(), err)
	}

	o.logFunc(2, o.dumpHttpResponse, resp)

	// response not okay?
	if resp.StatusCode/100 != 2 {
		// all paths in here lead to 'return'
		// noinspection GoUnhandledErrorResult
		defer resp.Body.Close()

		// Auth fail?
		if resp.StatusCode == 401 {
			// Auth fail at login API is fatal for this transaction
			if strings.HasSuffix(apstraUrl.Path, apiUrlUserLogin) {
				return newTalkToApstraErr(req, requestBody, resp,
					fmt.Sprintf("http %d at '%s' - check username/password",
						resp.StatusCode, apstraUrl))
			}

			// Auth fail with "doNotLogin == true" is fatal for this transaction
			if in.doNotLogin {
				return newTalkToApstraErr(req, requestBody, resp,
					fmt.Sprintf("http %d at '%s' and doNotLogin is %t",
						resp.StatusCode, apstraUrl, in.doNotLogin))
			}

			o.logStr(1, fmt.Sprintf("got http %d '%s' at '%s' attempting login", resp.StatusCode, resp.Status, apstraUrl.String()))
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

	// noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// If the caller gave us an httpBodyWriter, copy the response body into it and return
	if in.httpBodyWriter != nil {
		_, err := io.CopyBuffer(in.httpBodyWriter, resp.Body, nil)
		if err != nil {
			return fmt.Errorf("error while reading http response body - %w", err)
		}
		return nil
	}

	// figure out whether Apstra responded with a task ID
	var tIdR taskIdResponse
	taskResponseFound, err := peekParseResponseBodyAsTaskId(resp, &tIdR)
	if err != nil {
		return newTalkToApstraErr(req, requestBody, resp, "error peeking response body")
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
		return json.NewDecoder(resp.Body).Decode(in.apiResponse)
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
		return newTalkToApstraErr(req, requestBody, resp, "blueprint id not found in url nor in response body")
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

		detailedStatus, _ := json.Marshal(&taskResponse.DetailedStatus)

		return TalkToApstraErr{
			Request:  request,
			Response: response,
			Msg:      string(detailedStatus),
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

// newTalkToApstraErr returns a TalkToApstraErr. It's intended to be called after the
// http.Request has been executed with Do(), so the request body has already
// been "spent" by Read(). We'll fill it back in. The response body is likely to
// be closed by a 'defer body.Close()' somewhere, so we'll replace that as well,
// up to some reasonable limit (don't try to buffer gigabytes of data from the
// webserver).
func newTalkToApstraErr(req *http.Request, reqBody []byte, resp *http.Response, errMsg string) TalkToApstraErr {
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
	}
	// good unmarshal, but what about the contents?
	return result.TaskId != "", nil // no error; bool depends on string match
}
