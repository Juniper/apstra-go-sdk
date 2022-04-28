package apstraTelemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	aosApiStreamingConfig = "/api/streaming-config"
)

type AosGetStreamingConfigsResponse struct {
	Items []AosStreamingConfig `json:"items"`
}

type AosStreamingConfig struct {
	Status         AosStreamingConfigStatus `json:"status"`
	StreamingType  string                   `json:"streaming_type"`
	SequencingMode string                   `json:"sequencing_mode"`
	Protocol       string                   `json:"protocol"`
	Hostname       string                   `json:"hostname"`
	Id             string                   `json:"id"`
	Port           uint16                   `json:"Port"`
}

type AosStreamingConfigStatus struct {
	Status               AosStreamingConfigConnectionLog     `json:"status"`
	ConnectionTime       string                              `json:"connectionTime"`
	Epoch                string                              `json:"epoch"`
	ConnectionResetCount uint                                `json:"connnectionResetCount"`
	StreamingEndpoint    AosStreamingConfigStreamingEndpoint `json:"streamingEndpoint"`
	DnsLog               AosStreamingConfigDnsLog            `json:"dnsLog"`
	Connected            bool                                `json:"connected"`
	DisconnectionTime    string                              `json:"disconnectionTime"`
}

type AosStreamingConfigConnectionLog struct {
	Date    string `json:"date"'`
	Message string `json:"message"`
}

type AosStreamingConfigStreamingEndpoint struct {
	StreamingType  string `json:"streaming_type"`
	SequencingMode string `json:"sequencing_mode"`
	Protocol       string `json:"protocol"`
	Hostname       string `json:"Hostname"`
	Port           uint16 `json:"Port"`
}

type AosStreamingConfigDnsLog struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

type aosCreateStreamingConfigResponse struct {
	Id string `json:"id"`
}

func (o AosClient) GetAllStreamingConfigs() ([]*AosStreamingConfig, error) {
	req, err := http.NewRequest("GET", o.baseUrl+aosApiStreamingConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Authtoken", o.token)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling http.client.Do - %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing logout http response body - %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code is not '%d' got '%d' at '%s'", 200, resp.StatusCode, aosApiStreamingConfig)
	}

	var streamingConfigs AosGetStreamingConfigsResponse
	err = json.NewDecoder(resp.Body).Decode(&streamingConfigs)
	if err != nil {
		return nil, fmt.Errorf("error decoding aosGetStreamingConfigs JSON - %v", err)
	}

	var response []*AosStreamingConfig
	for _, asc := range streamingConfigs.Items {
		response = append(response, &asc)
	}

	return response, nil
}

func (o AosClient) GetStreamingConfig(id string) (*AosStreamingConfig, error) {
	req, err := http.NewRequest("GET", o.baseUrl+aosApiStreamingConfig+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Authtoken", o.token)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling http.client.Do - %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing logout http response body - %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code is not '%d' got '%d' at '%s'", 200, resp.StatusCode, aosApiStreamingConfig+"/"+id)
	}

	var streamingConfig AosStreamingConfig
	err = json.NewDecoder(resp.Body).Decode(&streamingConfig)
	if err != nil {
		return nil, fmt.Errorf("error decoding aosGetStreamingConfigs JSON - %v", err)
	}

	return &streamingConfig, nil
}

func (o AosClient) CreateStreamingConfig(cfg *AosStreamingConfigStreamingEndpoint) (string, error) {
	msg, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("error marshaling AosStreamingConfigStreamingEndpoint object - %v", err)
	}

	req, err := http.NewRequest("POST", o.baseUrl+aosApiStreamingConfig, bytes.NewBuffer(msg))
	if err != nil {
		return "", fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authtoken", o.token)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("http response code is not '%d' got '%d' at '%s'", 201, resp.StatusCode, aosApiStreamingConfig)
	}

	var createStreamingConfigResp *aosCreateStreamingConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&createStreamingConfigResp)
	if err != nil {
		return "", fmt.Errorf("error decoding aosCreateStreamingConfigResponse JSON - %v", err)
	}

	return createStreamingConfigResp.Id, nil
}
