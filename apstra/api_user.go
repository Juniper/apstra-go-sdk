// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
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
	Token string   `json:"token"`
	Id    ObjectId `json:"id"`
}

func (o *Client) startTaskMonitor() {
	if o.tmQuit == nil {
		o.tmQuit = make(chan struct{})
		o.Log(2, "starting task monitor")
		newTaskMonitor(o).start()
		o.Log(2, "task monitor started")
	}
}

func (o *Client) stopTaskMonitor() {
	if o.tmQuit != nil {
		close(o.tmQuit)
		o.tmQuit = nil
		o.Log(2, "task monitor close requested")
	}
}

func (o *Client) Config() ClientCfg {
	return o.cfg
}

func (o *Client) GetApiToken() string {
	o.lock(mutexKeyHttpHeaders)
	defer o.unlock(mutexKeyHttpHeaders)
	return o.httpHeaders[apstraAuthHeader]
}

func (o *Client) SetApiToken(in string) {
	o.lock(mutexKeyHttpHeaders)
	o.httpHeaders[apstraAuthHeader] = in
	o.unlock(mutexKeyHttpHeaders)
}

func (o *Client) login(ctx context.Context) error {
	response := &userLoginResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost,
		urlStr: apiUrlUserLogin,
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
	o.lock(mutexKeyHttpHeaders)
	o.httpHeaders[apstraAuthHeader] = response.Token
	o.unlock(mutexKeyHttpHeaders)

	o.id = response.Id
	o.startTaskMonitor()
	return nil
}

func (o *Client) logout(ctx context.Context) error {
	o.Log(1, "client logging out")

	defer func() { // clear the auth token and stop the task monitor
		o.Log(1, "deleting auth token")
		o.lock(mutexKeyHttpHeaders)
		delete(o.httpHeaders, apstraAuthHeader)
		o.unlock(mutexKeyHttpHeaders)
		o.Log(1, "shutting down the task monitor")
		o.stopTaskMonitor()
	}()

	o.lock(mutexKeyHttpHeaders)
	if token := o.httpHeaders[apstraAuthHeader]; token == "" { // doesn't exist OR is empty string?
		o.unlock(mutexKeyHttpHeaders)
		return nil // don't need to call the logout API if we have no token
	}
	o.unlock(mutexKeyHttpHeaders)

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:     http.MethodPost,
		urlStr:     apiUrlUserLogout,
		doNotLogin: true,
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' - %w", apiUrlUserLogout, convertTtaeToAceWherePossible(err))
	}
	return nil
}
