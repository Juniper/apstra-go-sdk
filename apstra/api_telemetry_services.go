package apstra

import (
	"context"
	"net/http"
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

func (o *Client) GetTelemetryServicesDeviceMapping(ctx context.Context) (*GetTelemetryServiceMappingResult, error) {
	var result GetTelemetryServiceMappingResult

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlTelemetryServices,
		apiResponse: &result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &result, nil
}
