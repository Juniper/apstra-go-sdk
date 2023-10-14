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
	Id                ObjectId    `json:"id"`
	Label             string      `json:"label"`
	TaskError         string      `json:"task_error"`
	AnomalyCount      int         `json:"anomaly_count"`
	Tags              []string    `json:"tags"`
	LastError         interface{} `json:"last_error"`
	UpdatedAt         string      `json:"updated_at"`
	UpdatedBy         string      `json:"updated_by"`
	Disabled          bool        `json:"disabled"`
	ConfigCompletedAt string      `json:"config_completed_at"`
	State             string      `json:"state"`
	Version           int         `json:"version"`
	HostNode          string      `json:"host_node"`
	TaskState         string      `json:"task_state"`
	ConfigStartedAt   string      `json:"config_started_at"`
	IbaUnit           string      `json:"iba_unit"`
	PredefinedProbe   string      `json:"predefined_probe"`
	Description       string      `json:"description"`
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

	for _, p := range pps {
		if p.Label == label {
			return &p, nil
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no Predefined Probe with label '%s' found", label),
	}
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
