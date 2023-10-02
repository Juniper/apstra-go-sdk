package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlIbaPredefinedProbes       = "/api/blueprints/%s/iba/predefined-probes"
	apiUrlIbaPredefinedProbesPrefix = apiUrlIbaPredefinedProbes + apiUrlPathDelim
	apiUrlIbaPredefinedProbesByName = apiUrlIbaPredefinedProbesPrefix + "%s"
)

type IbaPredefinedProbe struct {
	Name         string          `json:"name"`
	Experimental bool            `json:"experimental"`
	Description  string          `json:"description"`
	Schema       json.RawMessage `json:"schema"`
}

type IbaPredefinedProbeRequest struct {
	Name string
	Data json.RawMessage
}

func (o *Client) getAllIbaPredefinedProbes(ctx context.Context, bp_id ObjectId) ([]IbaPredefinedProbe, error) {
	response := &struct {
		Items []IbaPredefinedProbe `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbes, bp_id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaPredefinedProbeByName(ctx context.Context, bpId ObjectId, name string) (*IbaPredefinedProbe, error) {
	response := &IbaPredefinedProbe{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbesByName, bpId, name),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) instantiatePredefinedIbaProbe(ctx context.Context, bpid ObjectId, in *IbaPredefinedProbeRequest) (ObjectId, error) {
	response := &objectIdResponse{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbesByName, bpid, in.Name),
		apiInput:    in.Data,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}
