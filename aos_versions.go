package apstraTelemetry

const (
	aosApiVersionsPrefix = "/api/versions/"
	aosApiVersionsAosdi  = aosApiVersionsPrefix + "aosdi"
	aosApiVersionsApi    = aosApiVersionsPrefix + "api"
	aosApiVersionsBuild  = aosApiVersionsPrefix + "build"
	aosApiVersionsDevice = aosApiVersionsPrefix + "device"
	aosApiVersionsIba    = aosApiVersionsPrefix + "iba"
	aosApiVersionsNode   = aosApiVersionsPrefix + "node"
	aosApiVersionsServer = aosApiVersionsPrefix + "server"
)

type aosApiVersionsAosdiResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

type aosApiVersionsApiResponse struct {
	Major   string `json:"major"`
	Version string `json:"version"`
	Build   string `json:"build"`
	Minor   string `json:"minor"`
}

type aosApiVersionsBuildResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}

type aosApiVersionsDeviceRequest struct {
	SerialNumber string `json:"serial_number"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
}

type aosApiVersionsDeviceResponse struct {
	Status       string `json:"status"`
	Url          string `json:"url"`
	RetryTimeout int    `json:"retry_timeout"`
	Cksum        string `json:"cksum"`
}

type aosApiVersionsIbaRequest struct {
	Version  string `json:"version""`
	SystemId string `json:"system_id""`
}

type aosApiVersionsIbaResponse struct {
	Status       string `json:"status"`
	Url          string `json:"url"`
	RetryTimeout int    `json:"retry_timeout"`
	Cksum        string `json:"cksum"`
}

type aosApiVersionsNodeRequest struct {
	IpAddress string `json:"ip_address"`
	Version   string `json:"version"`
	SystemId  string `json:"system_id"`
}

type aosApiVersionsNodeResponse struct {
	Status       string `json:"status"`
	RetryTimeout int    `json:"retry_timeout"`
}

type aosApiVersionsServerResponse struct {
	Version       string `json:"version"`
	BuildDateTime string `json:"build_datetime"`
}
