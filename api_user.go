package goapstra

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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
	// and start the tasskMonitor (these go together)
	r := rand.Intn(100)
	os.Stderr.WriteString(fmt.Sprintf("xxxxx locking auth token %d...\n", r))
	o.lock(clientAuthTokenMutex)
	os.Stderr.WriteString(fmt.Sprintf("xxxxx locking auth token locked %d.", r))
	defer o.unlock(clientAuthTokenMutex)
	o.httpHeaders[apstraAuthHeader] = response.Token

	os.Stderr.WriteString(fmt.Sprintf("xxxxx unlocking auth token %d.", r))
	return nil
}

func (o Client) logout(ctx context.Context) error {
	// presence of an auth token is proxy for both
	// - "logged in" state and
	// - operation of a task monitor routine
	if _, tokenFound := o.httpHeaders[apstraAuthHeader]; !tokenFound {
		return nil
	}
	defer func() {
		// presence of auth token and taskMonitor go together
		delete(o.httpHeaders, apstraAuthHeader) // delete the auth token
		defer close(o.tmQuit)                   // shut down the task monitor gothread
	}()

	apstraUrl, err := url.Parse(apiUrlUserLogout)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlUserLogout, err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:     http.MethodPost,
		url:        apstraUrl,
		doNotLogin: true,
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' - %w", apiUrlUserLogout, err)
	}
	return nil
}
