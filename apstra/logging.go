package apstra

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Logger interface {
	Println(v ...any)
}

// logStr checks if DebugLevel meets the message verbosity specified in
// msgLevel. If so, it logs the supplied message (maybe)
func (o *Client) logStr(msgLevel int, msg string) {
	if o.logger == nil {
		return
	}

	if msgLevel > o.cfg.LogLevel {
		return
	}
	o.logger.Println(msg)
}

// logStrf checks if DebugLevel meets the message verbosity specified in
// msgLevel. If so, it formats the message and logs it.
func (o *Client) logStrf(msgLevel int, msg string, a ...any) {
	if o.logger == nil {
		return
	}

	if msgLevel > o.cfg.LogLevel {
		return
	}
	o.logger.Println(fmt.Sprintf(msg, a...))
}

// logFunc checks if DebugLevel meets the message verbosity specified in
// msgLevel. If so, it runs the supplied function with the supplied params. The
// string returned by the function is logged. If the function produces an
// error, it is logged directly and the intended log message is lost
func (o *Client) logFunc(msgLevel int, f func(int, ...interface{}) (string, error), params ...interface{}) {
	if o.logger == nil {
		return
	}

	if msgLevel > o.cfg.LogLevel {
		return
	}
	msg, err := f(msgLevel, params...)
	if err != nil {
		o.logStr(0, err.Error())
	}
	o.logStr(msgLevel, msg)
}

// dumpHttpRequest string-ifys an http.Request according to the desired
// verbosity of the incoming message (msgLevel) relative to the configured
// LogLevel.
// When msgLevel exceeds LogLevel count nothing is returned.
// When msgLevel matches LogLevel count, a short message is returned.
// When msgLevel is exceeded by LogLevel, progressively more information is
// returned. It is intended to be passed by name, and called by logFunc().
func (o *Client) dumpHttpRequest(msgLevel int, in ...interface{}) (string, error) {
	if len(in) != 1 {
		return "", fmt.Errorf("error dumping http request: expected 1 parameter, got %d", len(in))
	}
	// todo: detect non-printable body, base64 it prior to printing.
	req := in[0].(*http.Request)
	var data []byte
	var err error
	switch {
	case o.cfg.LogLevel-msgLevel < 0: // fallthrough to empty return
	case o.cfg.LogLevel-msgLevel == 0:
		return strings.Join([]string{req.Method, req.URL.String()}, " "), nil
	case o.cfg.LogLevel-msgLevel == 1:
		data, err = httputil.DumpRequestOut(req, false)
	default: // debug deltas > 1 get the request body
		data, err = httputil.DumpRequestOut(req, true)
	}
	return string(data), err
}

// dumpHttpResponse string-ifys an http.Hesponse according to the desired
// verbosity of the incoming message (msgLevel) relative to the configured
// LogLevel
// When msgLevel exceeds LogLevel nothing is returned.
// When msgLevel matches LogLevel, a short message is returned.
// When msgLevel is exceeded by LogLevel, progressively more information is
// returned. It is intended to be passed by name, and called by logFunc().
func (o *Client) dumpHttpResponse(msgLevel int, in ...interface{}) (string, error) {
	if len(in) != 1 {
		return "", fmt.Errorf("error dumping http request: expected 1 parameter, got %d", len(in))
	}
	// todo: detect non-printable body, base64 it prior to printing.
	resp := in[0].(*http.Response)
	var data []byte
	var err error
	switch {
	case o.cfg.LogLevel-msgLevel < 0: // fallthrough to empty return
	case o.cfg.LogLevel-msgLevel == 0:
		return resp.Status, nil
	case o.cfg.LogLevel-msgLevel == 1:
		data, err = httputil.DumpResponse(resp, false)
	default: // debug deltas > 1 get the request body
		data, err = httputil.DumpResponse(resp, true)
	}
	return string(data), err
}
