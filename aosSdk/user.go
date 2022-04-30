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

// userLogin submits credentials to an API server, collects a login token
// todo - need to handle token timeout
func (o *Client) userLogin() error {
	msg, err := json.Marshal(userLoginRequest{
		Username: o.cfg.User,
		Password: o.cfg.Pass,
	})
	if err != nil {
		return fmt.Errorf("error marshaling userLoginRequest object - %v", err)
	}

	url := o.baseUrl + apiUrlUserLogin
	err = o.post(url, msg, []int{201}, o.login)
	if err != nil {
		return fmt.Errorf("error calling '%s' - %v", url, err)
	}

	o.loginTime = time.Now()

	return nil
}

func (o Client) userLogout() error {
	err := o.post(o.baseUrl+apiUrlUserLogout, nil, []int{200}, nil)
	return err
}
