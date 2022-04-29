package aosSdk

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	schemeHttp        = "http"
	schemeHttps       = "https"
	schemeHttpsUnsafe = "hxxps"

	aosApiUserLogin  = "/api/user/login"
	aosApiUserLogout = "/api/user/logout"
)

// aosUserLoginRequest token to the aosApiUserLogin API endpoint
type aosUserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// aosUserLoginResponse token returned by the aosApiUserLogin API endpoint
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

	url := o.baseUrl + aosApiUserLogin
	err = o.post(url, msg, []int{201}, o.login)
	if err != nil {
		return fmt.Errorf("error calling '%s' - %v", url, err)
	}

	o.loginTime = time.Now()

	return nil
}

func (o AosClient) userLogout() error {
	err := o.post(o.baseUrl+aosApiUserLogout, nil, []int{200}, nil)
	return err
}
