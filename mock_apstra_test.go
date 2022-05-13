package goapstra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func newMockTestClient() (*Client, error) {
	c, err := newLiveTestClient()
	c.httpClient = &mockApstraApi{
		username: "admin",
		password: "admin",
	}
	if err != nil {
		log.Fatal(err)
	}
	return c, err
}

type mockApstraApi struct {
	username  string
	password  string
	authToken string
}

func (o mockApstraApi) Do(req *http.Request) (*http.Response, error) {
	switch req.URL.Path {
	case apiUrlUserLogin:
		return o.handleLogin(req)
	case apiUrlUserLogout:
		return o.handleLogout(req)
	default:
		return nil, fmt.Errorf("mock client doesn't handle API path '%s'", req.URL.Path)
	}
}

func (o mockApstraApi) handleLogin(req *http.Request) (*http.Response, error) {
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
	o.authToken = randString(20)
	outBody, err := json.Marshal(userLoginResponse{
		Token: o.authToken,
		Id:    randString(10),
	})

	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(outBody)),
		StatusCode: http.StatusCreated,
	}, nil
}

func (o mockApstraApi) handleLogout(req *http.Request) (*http.Response, error) {
	for _, val := range req.Header.Values(apstraAuthHeader) {
		if val == o.authToken {
			return &http.Response{
				StatusCode: http.StatusOK,
			}, nil
		}
	}
	return nil, fmt.Errorf("logout attempt without valid token in mockApstraApi.handleLogin()")
}
