package apstra

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
