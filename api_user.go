package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlUserLogin  = "/api/user/login"
	apiUrlUserLogout = "/api/user/logout"
)

// userLoginRequest token to the apiUrlUserLogin API endpoint
type userLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// userLoginResponse token returned by the apiUrlUserLogin API endpoint
type userLoginResponse struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

func (o *Client) login(ctx context.Context) error {
	apstraUrl, err := url.Parse(apiUrlUserLogin)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlUserLogin, err)
	}
	response := &userLoginResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost,
		url:    apstraUrl,
		apiInput: &userLoginRequest{
			Username: o.cfg.User,
			Password: o.cfg.Pass,
		},
		doNotLogin:  true,
		apiResponse: response,
	})
	if err != nil {
		return fmt.Errorf("error talking to AOS in Login - %w", err)
	}

	// stash auth token in client's default set of apstra http httpHeaders
	o.httpHeaders[apstraAuthHeader] = response.Token

	return nil
}

func (o Client) logout(ctx context.Context) error {
	defer close(o.tmQuit) // shut down the task monitor gothread

	apstraUrl, err := url.Parse(apiUrlUserLogout)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlUserLogout, err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost,
		url:    apstraUrl,
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' - %w", apiUrlUserLogout, err)
	}
	delete(o.httpHeaders, apstraAuthHeader)
	return nil
}
