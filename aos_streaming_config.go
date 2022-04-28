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
	Port           uint16                   `json:"port"`
}

type AosStreamingConfigStatus struct {
	Status               AosStreamingConfigConnectionLog     `json:"status"`
	ConnectionTime       string                              `json:"connectionTime"`
	Epoch                string                              `json:"epoch"`
	ConnectionResetCount uint                                `json:"connnectionResetCount"`
	StreamingEndpoint    AosStreamingConfigStreamingEndpoint `json:"streamingEndpoint"`
	DnsLog               []AosStreamingConfigDnsLog          `json:"dnsLog"`
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
	var agscr AosGetStreamingConfigsResponse
	url := o.baseUrl + aosApiStreamingConfig
	err := o.newGet(url, []int{200}, &agscr)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	var result []*AosStreamingConfig
	for _, i := range agscr.Items {
		result = append(result, &i)
	}
	return result, nil
}

func (o AosClient) GetStreamingConfig(id string) (*AosStreamingConfig, error) {
	var result AosStreamingConfig
	url := o.baseUrl + aosApiStreamingConfig + "/" + id
	err := o.newGet(url, []int{200}, result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	return &result, nil
}

func (o AosClient) CreateStreamingConfig(cfg *AosStreamingConfigStreamingEndpoint) (string, error) {
	msg, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("error marshaling AosStreamingConfigStreamingEndpoint object - %v", err)
	}

	var result aosCreateStreamingConfigResponse
	url := o.baseUrl + aosApiStreamingConfig
	err = o.newPost(url, msg, []int{201}, &result)
	if err != nil {
		return "", fmt.Errorf("error calling %s - %v", url, err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
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
