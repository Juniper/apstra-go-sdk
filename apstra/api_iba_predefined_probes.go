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
	Label        string          `json:"label"`
	Experimental bool            `json:"experimental"`
	Description  string          `json:"description"`
	Schema       json.RawMessage `json:"schema"`
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

func (o *Client) instantiatePredefinedIbaProbe(ctx context.Context, bpid ObjectId, in *IbaPredefinedProbe) (ObjectId, error) {
	response := &objectIdResponse{}
	input := struct {
		Label string `json:"label"`
	}{Label: in.Label}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPost, urlStr: fmt.Sprintf(apiUrlIbaPredefinedProbesByName, bpid, in.Name),
		apiInput: input, apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}
