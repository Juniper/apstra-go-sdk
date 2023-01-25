package goapstra

import (
	"context"
	"net/http"
)

const (
	apiUrlTelemetryQuery = "/api/telemetry-query"
)

type TelemetryQueryResponse struct {
	Services []string `json:"services"`
}

func (o *Client) getTelemetryQuery(ctx context.Context) (*TelemetryQueryResponse, error) {
	response := &TelemetryQueryResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlTelemetryQuery,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}
