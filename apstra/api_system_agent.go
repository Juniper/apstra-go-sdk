package apstra

import (
	"context"
	"net/http"
)

const apiUrlSystemAgentManagerConfig = "/api/system-agent/manager-config"

type SystemAgentManagerConfig struct {
	SkipRevertToPristineOnUninstall bool `json:"skip_revert_to_pristine_on_uninstall"`
	SkipPristineValidation          bool `json:"skip_pristine_validation"`
}

func (o *Client) getSystemAgentManagerConfig(ctx context.Context) (*SystemAgentManagerConfig, error) {
	result := &SystemAgentManagerConfig{}
	return result, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlSystemAgentManagerConfig,
		apiResponse: result,
	})
}

func (o *Client) setSystemAgentManagerConfig(ctx context.Context, cfg *SystemAgentManagerConfig) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   apiUrlSystemAgentManagerConfig,
		apiInput: cfg,
	})
}
