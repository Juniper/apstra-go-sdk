package aosSdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	schemeHttp        = "http"
	schemeHttps       = "https"
	schemeHttpsUnsafe = "hxxps"

	aosApiUserLogin  = "/api/user/login"
	aosApiUserLogout = "/api/user/logout"
)

// aosUserLoginRequest payload to the aosApiUserLogin API endpoint
type aosUserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// aosUserLoginResponse payload returned by the aosApiUserLogin API endpoint
type aosUserLoginResponse struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

// userLogin submits credentials to an API server, collects a login token
// todo - need to handle token timeout
func (o *AosClient) userLogin() error {
	msg, err := json.Marshal(aosUserLoginRequest{
		Username: o.cfg.User,
		Password: o.cfg.Pass,
	})
	if err != nil {
		return fmt.Errorf("error marshaling aosLogin object - %v", err)
	}

	req, err := http.NewRequest("POST", o.baseUrl+aosApiUserLogin, bytes.NewBuffer(msg))
	if err != nil {
		return fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("http response code is not '%d' got '%d' at '%s'", 201, resp.StatusCode, aosApiUserLogin)
	}

	var loginResp *aosUserLoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return fmt.Errorf("error decoding aosUserLoginResponse JSON - %v", err)
	}

	o.token = loginResp.Token

	return nil
}

func (o AosClient) userLogout() error {
	err := o.post(o.baseUrl+aosApiUserLogout, nil, []int{200}, nil)
	return err
}
