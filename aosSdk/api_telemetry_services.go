package aosSdk

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlTelemetryServices = "/api/telemetry/services"
)

type TelemetryServiceMapping struct {
	ServiceName string   `json:"service_name"`
	Devices     []string `json:"devices"`
}

type GetTelemetryServiceMappingResult struct {
	Configured    []TelemetryServiceMapping `json:"configured"`
	LastRunError  []TelemetryServiceMapping `json:"last_run_error"`
	EnablingError []TelemetryServiceMapping `json:"enabling_error"`
}

func (o Client) GetTelemetryServicesDeviceMapping() (*GetTelemetryServiceMappingResult, error) {
	aosUrl, err := url.Parse(apiUrlTelemetryServices)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlTelemetryServices, err)
	}
	result := &GetTelemetryServiceMappingResult{}
	return result, o.talkToAos(&talkToAosIn{
		method:      http.MethodGet,
		url:         aosUrl,
		apiResponse: result,
	})
}
