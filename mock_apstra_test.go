package goapstra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	mockApstraUser = "mockAdmin"
	mockApstraPass = "mockPassword"
)

type mockApstraApi struct {
	username         string
	password         string
	authToken        string
	metricdb         metricdb
	virtualIfraMgrs  virtualInfraMgrsResponse
	resourceAsnPools []rawAsnPool
	anomalies        []string
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
	// todo: inspect HTTP method in addition to URL path?
	switch {
	case req.URL.Path == apiUrlUserLogin:
		return o.handleLogin(req)
	case req.URL.Path == apiUrlUserLogout:
		return o.handleLogout(req)
	case req.URL.Path == apiUrlMetricdbMetric:
		return o.handleMetricdbMetric(req)
	case req.URL.Path == apiUrlMetricdbQuery:
		return o.handleMetricdbQuery(req)
	case req.URL.Path == apiUrlVirtualInfraManagers:
		return o.handleVirtualInfraManagers(req)
	case req.URL.Path == apiUrlResourcesAsnPools:
		return o.handleApiUrlResourcesAsnPools(req)
	case req.URL.Path == apiUrlAnomalies:
		return o.handleApiUrlAnomalies(req)
	case strings.HasPrefix(req.URL.Path, apiUrlResourcesAsnPoolsPrefix):
		return o.handleApiUrlResourcesAsnPoolsPrefix(req)
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
		Status:     "401 Unauthorized",
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
		Id:    string(randId()),
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

func (o mockApstraApi) handleMetricdbMetric(req *http.Request) (*http.Response, error) {
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}

	outBody, err := json.Marshal(&o.metricdb.metrics)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(outBody)),
		StatusCode: http.StatusOK,
		Status:     "200 OK",
	}, nil
}

func (o mockApstraApi) handleMetricdbQuery(req *http.Request) (*http.Response, error) {
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}

	outBody, err := json.Marshal(&MetricDbQueryResponse{
		Status: struct {
			TotalCount    int       `json:"total_count"`
			BeginTime     time.Time `json:"begin_time"`
			EndTime       time.Time `json:"end_time"`
			ResultCode    string    `json:"result_code"`
			LastTimestamp time.Time `json:"last_timestamp"`
		}{},
		// todo
		//Items: []json.RawMessage{[]byte{string"{bogus_data: }"}},
		Items: nil,
	})
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

func (o *mockApstraApi) handleApiUrlResourcesAsnPools(req *http.Request) (*http.Response, error) {
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}

	switch req.Method {
	case http.MethodGet:
		body, err := json.Marshal(getAsnPoolsResponse{Items: o.resourceAsnPools})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(bytes.NewReader(body)),
		}, nil
	case http.MethodPost:
		inBody, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		newRawPool := &rawAsnPool{}
		err = json.Unmarshal(inBody, newRawPool)
		if err != nil {
			return nil, err
		}
		var total uint32
		for _, i := range newRawPool.Ranges {
			size := 1 + i.Last - i.First
			total += size
		}
		newRawPool.Used = "0"
		newRawPool.Total = strconv.Itoa(int(total))
		for _, existingRawPool := range o.resourceAsnPools {
			existingPool, err := rawAsnPoolToAsnPool(&existingRawPool)
			if err != nil {
				return nil, err
			}

			newPool, err := rawAsnPoolToAsnPool(newRawPool)
			if err != nil {
				return nil, err
			}

			// todo: Apstra doesn't enforce this check
			if asnPoolOverlap(*existingPool, *newPool) {
				return nil, fmt.Errorf("overlap with existing asn pool %s", existingRawPool.Id)
			}
		}
		newRawPool.Id = randId()
		o.resourceAsnPools = append(o.resourceAsnPools, *newRawPool)
		body, err := json.Marshal(objectIdResponse{Id: newRawPool.Id})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Status:     "202 ACCEPTED",
			Body:       io.NopCloser(bytes.NewReader(body)),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported method '%s'", req.Method)
	}
}

func asnPoolOverlap(a, b AsnPool) bool {
	var rar, rbr AsnRange
	for _, ra := range a.Ranges {
		rar.First = ra.First
		rar.Last = ra.Last
		for _, rb := range b.Ranges {
			rbr.First = rb.First
			rbr.Last = rb.Last
			if asnOverlap(rar, rbr) {
				return true
			}
		}
	}
	return false
}

func (o *mockApstraApi) handleApiUrlResourcesAsnPoolsPrefix(req *http.Request) (*http.Response, error) {
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}

	requested := ObjectId(strings.TrimPrefix(req.URL.Path, apiUrlResourcesAsnPoolsPrefix))
	switch req.Method {
	case http.MethodGet:
		for _, p := range o.resourceAsnPools {
			if p.Id == requested {
				body, err := json.Marshal(p)
				if err != nil {
					return nil, err
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			}
		}
		return &http.Response{
			StatusCode: http.StatusNotFound,
		}, nil
	case http.MethodDelete:
		id := ObjectId(strings.TrimPrefix(req.URL.Path, apiUrlResourcesAsnPoolsPrefix))
		for i, pool := range o.resourceAsnPools {
			if pool.Id == id {
				o.resourceAsnPools = append(o.resourceAsnPools[:i], o.resourceAsnPools[i+1:]...)
				return &http.Response{
					StatusCode: 202,
					Status:     "202 Accepted",
					Body:       io.NopCloser(strings.NewReader("{}"))}, nil
			}
		}
		return nil, fmt.Errorf("ASN resource pool '%s' not found", id)
	case http.MethodPut:
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Status:     "202 ACCEPTED",
			Body:       io.NopCloser(strings.NewReader("{}"))}, nil
	default:
		return nil, fmt.Errorf("method '%s' not supported", req.Method)
	}
}

func (o mockApstraApi) handleApiUrlAnomalies(req *http.Request) (*http.Response, error) {
	if resp, ok := o.auth(req); !ok {
		return resp, nil
	}
	if req.Method != http.MethodGet {
		return nil, fmt.Errorf("mock '%s' only supports '%s'", apiUrlAnomalies, http.MethodGet)
	}
	result := &getAnomaliesResponse{}
	for _, a := range o.anomalies {
		result.Items = append(result.Items, []byte(a))
	}
	body, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("mock anomaly error: %w", err)
	}
	return &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       io.NopCloser(bytes.NewReader(body)),
	}, nil
}
