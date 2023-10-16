package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlIbaProbes       = "/api/blueprints/%s/probes"
	apiUrlIbaProbesPrefix = apiUrlIbaProbes + apiUrlPathDelim
	apiUrlIbaProbesById   = apiUrlIbaProbesPrefix + "%s"
)

type IbaProbe struct {
	Id              ObjectId                 `json:"id"`
	Label           string                   `json:"label"`
	TaskError       string                   `json:"task_error"`
	Stages          []map[string]interface{} `json:"stages"`
	AnomalyCount    int                      `json:"anomaly_count"`
	Tags            []string                 `json:"tags"`
	Disabled        bool                     `json:"disabled"`
	State           string                   `json:"state"`
	Version         int                      `json:"version"`
	TaskState       string                   `json:"task_state"`
	IbaUnit         string                   `json:"iba_unit"`
	PredefinedProbe string                   `json:"predefined_probe"`
	Description     string                   `json:"description"`
}

func (o *Client) getAllIbaProbes(ctx context.Context, bpId ObjectId) ([]IbaProbe, error) {
	response := &struct {
		Items []IbaProbe `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaProbes, bpId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIbaProbeByLabel(ctx context.Context, bpId ObjectId, label string) (*IbaProbe, error) {
	pps, err := o.getAllIbaProbes(ctx, bpId)
	if err != nil {
		return nil, err
	}
	var probe IbaProbe
	i := 0
	for _, p := range pps {
		if p.Label == label {
			probe := p
			i := i + 1
		}
	}
	if i == 0 {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no Predefined Probe with label '%s' found", label),
		}
	}
	if i > 1 {
		return nil, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("too many probes with label %s found, expected 1 got %d", label, i),
		}
	}
	return &probe, nil
}

func (o *Client) getIbaProbe(ctx context.Context, bpId ObjectId, id ObjectId) (*IbaProbe, error) {
	response := &IbaProbe{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlIbaProbesById, bpId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) deleteIbaProbe(ctx context.Context, bpId ObjectId, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlIbaProbesById, bpId, id),
	})
}
