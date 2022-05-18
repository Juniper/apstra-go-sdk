package goapstra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	mockApstraUser = "mockAdmin"
	mockApstraPass = "mockPassword"
)

type mockApstraApi struct {
	username        string
	password        string
	authToken       string
	metricdb        metricdb
	virtualIfraMgrs virtualInfraMgrsResponse
}

type metricdb struct {
	metrics metricdbMetrics
}

type metricdbMetrics struct {
	Items []MetricdbMetric
}

func newMockApstraApi(password string) (*mockApstraApi, error) {
	var err error
	mock := &mockApstraApi{
		username: mockApstraUser,
	}
	mock.changePassword(password)
	err = mock.createMetricdb()
	if err != nil {
		return nil, err
	}

	return mock, nil
}

func (o *mockApstraApi) changePassword(password string) {
	o.password = password
	o.authToken = randJwt()
}

func (o *mockApstraApi) Do(req *http.Request) (*http.Response, error) {
	switch {
	case req.URL.Path == apiUrlUserLogin:
		return o.handleLogin(req)
	case req.URL.Path == apiUrlUserLogout:
		return o.handleLogout(req)
	case req.URL.Path == apiUrlMetricdbMetric:
		return o.handleMetricdb(req)
	case req.URL.Path == apiUrlVirtualInfraManagers:
		return o.handleVirtualInfraManagers(req)
	default:
		return nil, fmt.Errorf("mock client doesn't handle API path '%s'", req.URL.Path)
	}
}

func (o *mockApstraApi) auth(req *http.Request) (*http.Response, bool) {
	for _, val := range req.Header.Values(apstraAuthHeader) {
		if val == o.authToken {
			return nil, true
		}
	}
	return &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewReader(nil)),
	}, false
}

func (o *mockApstraApi) handleLogin(req *http.Request) (*http.Response, error) {
	inBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body in mockApstraApi.handleLogin() - %w", err)
	}

	in := &userLoginRequest{}
	err = json.Unmarshal(inBody, in)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling userLoginRequest in mockApstraApi.handleLogin() - %w", err)
	}

	if in.Username != o.username || in.Password != o.password {
		return nil, fmt.Errorf("error bad authentication in mockApstraApi.handleLogin() '%s:%s' vs. '%s:%s",
			in.Username, in.Password, o.username, o.password)
	}
	o.authToken = randJwt()
	outBody, err := json.Marshal(userLoginResponse{
		Token: o.authToken,
		Id:    randId(),
	})

	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(outBody)),
		StatusCode: http.StatusCreated,
		Status:     "201 CREATED",
	}, nil
}

func (o mockApstraApi) handleLogout(req *http.Request) (*http.Response, error) {
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}
	o.authToken = ""
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(nil)),
		Status:     "200 OK",
	}, nil

}

func (o mockApstraApi) handleMetricdb(req *http.Request) (*http.Response, error) {
	// so far only GET /api/metricdb/metric supported
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}

	outBody, err := json.Marshal(o.metricdb)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(outBody)),
		StatusCode: http.StatusOK,
		Status:     "200 OK",
	}, nil
}

func (o *mockApstraApi) handleVirtualInfraManagers(req *http.Request) (*http.Response, error) {
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}

	body, err := json.Marshal(o.virtualIfraMgrs)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}, nil
}
