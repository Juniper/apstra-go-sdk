package goapstra

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// DebugLevel is the configured log noisiness ceiling. Set it higher for more
// noise. Default value is 0. Set <0 for no debugs.
var DebugLevel int

// debugStr checks if DebugLevel meets the message verbosity specified in
// msgLevel. If so, it logs the supplied message (maybe)
func debugStr(msgLevel int, msg string) {
	if msgLevel > DebugLevel {
		return
	}
	log.Println(msg)
}

// debugFunc checks if DebugLevel meets the message verbosity specified in
// msgLevel. If so, it runs the supplied function with the supplied params. The
// string returned by the function is logged. If the function produces an
// error, it is logged directly and the intended log message is lost
func debugFunc(msgLevel int, f func(int, ...interface{}) (string, error), params ...interface{}) {
	if msgLevel > DebugLevel {
		return
	}
	msg, err := f(msgLevel, params...)
	if err != nil {
		debugStr(0, err.Error())
	}
	debugStr(msgLevel, msg)
}

// dumpHttpRequest string-ifys an http request according to the desired
// verbosity of the incoming message (msgLevel) and the configured DebugLevel.
// When msgLevel exceeds DebugLevel nothing is returned.
// When msgLevel matches DebugLevel, a short message is returned.
// When msgLevel is exceeded by DebugLevel, progressively more information is
// returned. It is intended to be passed by name, and called by debugFunc().
func dumpHttpRequest(msgLevel int, in ...interface{}) (string, error) {
	if len(in) != 1 {
		return "", fmt.Errorf("error dumping http request: expected 1 parameter, got %d", len(in))
	}
	// todo: detect non-printable body, base64 it prior to printing.
	req := in[0].(*http.Request)
	var data []byte
	var err error
	switch {
	case DebugLevel-msgLevel < 0: // fallthrough to empty return
	case DebugLevel-msgLevel == 0:
		return strings.Join([]string{req.Method, req.URL.String()}, " "), nil
	case DebugLevel-msgLevel == 1:
		data, err = httputil.DumpRequestOut(req, false)
	default: // debug deltas > 1 get the request body
		data, err = httputil.DumpRequestOut(req, true)
	}
	return string(data), err
}

// dumpHttpResponse string-ifys an http request according to the desired
// verbosity of the incoming message (msgLevel) and the configured DebugLevel.
// When msgLevel exceeds DebugLevel nothing is returned.
// When msgLevel matches DebugLevel, a short message is returned.
// When msgLevel is exceeded by DebugLevel, progressively more information is
// returned. It is intended to be passed by name, and called by debugFunc().
func dumpHttpResponse(msgLevel int, in ...interface{}) (string, error) {
	if len(in) != 1 {
		return "", fmt.Errorf("error dumping http response: expected 1 parameter, got %d", len(in))
	}
	// todo: detect non-printable body
	resp := in[0].(*http.Response)
	var data []byte
	var err error
	switch {
	case DebugLevel-msgLevel < 0: // fallthrough to empty return
	case DebugLevel-msgLevel == 0:
		return resp.Status, nil
	case DebugLevel-msgLevel == 1:
		data, err = httputil.DumpResponse(resp, false)
	default: // debug deltas > 1 get the response body
		data, err = httputil.DumpResponse(resp, true)
	}
	return string(data), err
}
