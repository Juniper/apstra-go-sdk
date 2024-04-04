package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

const (
	apiUrlSystemsConfiguration = apiUrlSystemsById + apiUrlPathDelim + "configuration"
)

var _ json.Unmarshaler = new(SystemConfig)

type SystemConfig struct {
	SystemId ObjectId
	//DeployState                   string
	//ConfigurationServiceState     string
	LastBootTime                  time.Time
	Deviated                      bool
	ErrorMessage                  *string
	ContiguousFailures            uint32
	UserGoldenConfigUpdateVersion uint32
	UserFullConfigDeployVersion   uint32
	AosConfigVersion              uint32
	ExpectedConfig                string
	ActualConfig                  string
}

func (o *SystemConfig) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		SystemId ObjectId `json:"system_id"`
		//DeployState                   string   `json:"deploy_state"`
		//ConfigurationServiceState     string   `json:"configuration_service_state"`
		LastBootTime                  float64 `json:"last_boot_time"`
		Deviated                      bool    `json:"deviated"`
		ErrorMessage                  string  `json:"error_message,omitempty"`
		ContiguousFailures            uint32  `json:"contiguous_failures"`
		UserGoldenConfigUpdateVersion uint32  `json:"user_golden_config_update_version"`
		UserFullConfigDeployVersion   uint32  `json:"user_full_config_deploy_version"`
		AosConfigVersion              uint32  `json:"aos_config_version"`
		Expected                      struct {
			Config string `json:"config"`
		} `json:"expected"`
		Actual struct {
			Config string `json:"config"`
		} `json:"actual"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	sec, nsec := math.Modf(raw.LastBootTime)
	lastBootTime := time.Unix(int64(sec), int64(nsec*(1e9)))

	var errorMessage *string
	if raw.ErrorMessage != "" {
		errorMessage = &raw.ErrorMessage
	}

	o.SystemId = raw.SystemId
	o.LastBootTime = lastBootTime
	o.Deviated = raw.Deviated
	o.ErrorMessage = errorMessage
	o.ContiguousFailures = raw.ContiguousFailures
	o.UserGoldenConfigUpdateVersion = raw.UserGoldenConfigUpdateVersion
	o.UserFullConfigDeployVersion = raw.UserFullConfigDeployVersion
	o.AosConfigVersion = raw.AosConfigVersion
	o.ExpectedConfig = raw.Expected.Config
	o.ActualConfig = raw.Actual.Config

	return nil
}

func (o *Client) GetSystemConfig(ctx context.Context, id ObjectId) (SystemConfig, error) {
	var response SystemConfig
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlSystemsConfiguration, id),
		apiResponse: &response,
	})
	if err != nil {
		return response, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}
