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
	apiUrlIbaProbes                 = "/api/blueprints/%s/probes"
	apiUrlIbaProbesPrefix           = apiUrlIbaProbes + apiUrlPathDelim
	apiUrlIbaProbesById             = apiUrlIbaProbesPrefix + "%s"
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

func (o *Client) getAllIbaPredefinedProbes(ctx context.Context, bpId ObjectId) ([]IbaPredefinedProbe, error) {
	response := &struct {
		Items []IbaPredefinedProbe `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbes, bpId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaPredefinedProbeByName(ctx context.Context, bpId ObjectId, name string) (*IbaPredefinedProbe, error) {
	pps, err := o.getAllIbaPredefinedProbes(ctx, bpId)
	if err != nil {
		return nil, err
	}

	for _, p := range pps {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no Predefined Probe with name '%s' found", name),
	}
}

func (o *Client) instantiatePredefinedIbaProbe(ctx context.Context, bpId ObjectId, in *IbaPredefinedProbeRequest) (ObjectId, error) {
	response := &objectIdResponse{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlIbaPredefinedProbesByName, bpId, in.Name),
		apiInput:    in.Data,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) deleteIbaProbe(ctx context.Context, bpId ObjectId, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlIbaProbesById, bpId, id),
	})
}
