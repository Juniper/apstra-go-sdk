// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"net/http"
)

const apiUrlSystemAgentManagerConfig = "/api/system-agent/manager-config"

type SystemAgentManagerConfig struct {
	DeviceOsImageDownloadTimeout    *int `json:"device_os_image_download_timeout,omitempty"` // introduced in and required by 5.1.0 (1-2700)
	SkipRevertToPristineOnUninstall bool `json:"skip_revert_to_pristine_on_uninstall"`
	SkipPristineValidation          bool `json:"skip_pristine_validation"`
	SkipInterfaceShutdownOnUpgrade  bool `json:"skip_interface_shutdown_on_upgrade"`
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
