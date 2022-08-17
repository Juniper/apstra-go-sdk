package goapstra

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

// logStr checks if DebugLevel meets the message verbosity specified in
// msgLevel. If so, it logs the supplied message (maybe)
func (o *Client) logStr(msgLevel int, msg string) {
	if msgLevel >= len(o.loggers) {
		return
	}
	o.loggers[msgLevel].Println(msg)
}

// logFunc checks if DebugLevel meets the message verbosity specified in
// msgLevel. If so, it runs the supplied function with the supplied params. The
// string returned by the function is logged. If the function produces an
// error, it is logged directly and the intended log message is lost
func (o *Client) logFunc(msgLevel int, f func(int, ...interface{}) (string, error), params ...interface{}) {
	if msgLevel >= len(o.loggers) {
		return
	}
	msg, err := f(msgLevel, params...)
	if err != nil {
		o.logStr(0, err.Error())
	}
	o.logStr(msgLevel, msg)
}

// dumpHttpRequest string-ifys an http request according to the desired
// verbosity of the incoming message (msgLevel) relative to the count of
// configured loggers:
// When msgLevel exceeds logger count nothing is returned.
// When msgLevel matches logger count, a short message is returned.
// When msgLevel is exceeded by logger count, progressively more information is
// returned. It is intended to be passed by name, and called by logFunc().
func (o *Client) dumpHttpRequest(msgLevel int, in ...interface{}) (string, error) {
	// todo: revisit this function with logger count in mind
	if len(in) != 1 {
		return "", fmt.Errorf("error dumping http request: expected 1 parameter, got %d", len(in))
	}
	// todo: detect non-printable body, base64 it prior to printing.
	loggerCount := len(o.loggers)
	req := in[0].(*http.Request)
	var data []byte
	var err error
	switch {
	case loggerCount-msgLevel < 0: // fallthrough to empty return
	case loggerCount-msgLevel == 0:
		return strings.Join([]string{req.Method, req.URL.String()}, " "), nil
	case loggerCount-msgLevel == 1:
		data, err = httputil.DumpRequestOut(req, false)
	default: // debug deltas > 1 get the request body
		data, err = httputil.DumpRequestOut(req, true)
	}
	return string(data), err
}

// dumpHttpResponse string-ifys an http response according to the desired
// verbosity of the incoming message (msgLevel) relative to the count of
// configured loggers:
// When msgLevel exceeds logger count nothing is returned.
// When msgLevel matches logger count, a short message is returned.
// When msgLevel is exceeded by logger count, progressively more information is
// returned. It is intended to be passed by name, and called by logFunc().
func (o *Client) dumpHttpResponse(msgLevel int, in ...interface{}) (string, error) {
	// todo: revisit this function with logger count in mind
	if len(in) != 1 {
		return "", fmt.Errorf("error dumping http request: expected 1 parameter, got %d", len(in))
	}
	// todo: detect non-printable body, base64 it prior to printing.
	loggerCount := len(o.loggers)
	resp := in[0].(*http.Response)
	var data []byte
	var err error
	switch {
	case loggerCount-msgLevel < 0: // fallthrough to empty return
	case loggerCount-msgLevel == 0:
		return resp.Status, nil
	case loggerCount-msgLevel == 1:
		data, err = httputil.DumpResponse(resp, false)
	default: // debug deltas > 1 get the request body
		data, err = httputil.DumpResponse(resp, true)
	}
	return string(data), err
}
