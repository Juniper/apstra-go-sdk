package aosSdk

import (
	"encoding/json"
	"fmt"
)

const (
	apiUrlStreamingConfig = "/api/streaming-config"

	StreamingConfigSequencingModeUnknown StreamingConfigSequencingMode = iota
	StreamingConfigSequencingModeSequenced
	StreamingConfigSequencingModeUnsequenced

	StreamingConfigStreamingTypeUnknown AosApiStreamingConfigStreamingType = iota
	StreamingConfigStreamingTypeAlerts
	StreamingConfigStreamingTypeEvents
	StreamingConfigStreamingTypePerfmon

	StreamingConfigProtocolUnknown AosApiStreamingConfigProtocol = iota
	StreamingConfigProtocolProtoBufOverTcp
)

type StreamingConfigId string

type StreamingConfigSequencingMode int

func (o StreamingConfigSequencingMode) String() string {
	switch o {
	case StreamingConfigSequencingModeUnknown:
		return "unknown"
	case StreamingConfigSequencingModeSequenced:
		return "sequenced"
	case StreamingConfigSequencingModeUnsequenced:
		return "unsequenced"
	default:
		return fmt.Sprintf("sequencing type %d has no string value", o)
	}
}

type AosApiStreamingConfigStreamingType int

func (o AosApiStreamingConfigStreamingType) String() string {
	switch o {
	case StreamingConfigStreamingTypeUnknown:
		return "unknown"
	case StreamingConfigStreamingTypeAlerts:
		return "alerts"
	case StreamingConfigStreamingTypeEvents:
		return "events"
	case StreamingConfigStreamingTypePerfmon:
		return "perfmon"
	default:
		return fmt.Sprintf("message type %d has no string value", o)
	}
}

type AosApiStreamingConfigProtocol int

func (o AosApiStreamingConfigProtocol) String() string {
	switch o {
	case StreamingConfigProtocolUnknown:
		return "unknown"
	case StreamingConfigProtocolProtoBufOverTcp:
		return "protoBufOverTcp"
	default:
		return fmt.Sprintf("message type %d has no string value", o)
	}
}

type AosGetStreamingConfigsResponse struct {
	Items []AosGetStreamingConfigResponse `json:"items"`
}

type AosGetStreamingConfigResponse struct {
	Status         AosStreamingConfigStatus `json:"status"`
	StreamingType  string                   `json:"streaming_type"`
	SequencingMode string                   `json:"sequencing_mode"`
	Protocol       string                   `json:"protocol"`
	Hostname       string                   `json:"hostname"`
	Id             string                   `json:"id"`
	Port           uint16                   `json:"port"`
}

type AosStreamingConfigStatus struct {
	ConnectionLog        AosStreamingConfigConnectionLog     `json:"connectionLog"`
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
	Hostname       string `json:"hostname"`
	Port           uint16 `json:"port"`
}

type AosStreamingConfigDnsLog struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

type createStreamingConfigResponse struct {
	Id string `json:"id"`
}

func (o Client) getAllStreamingConfigs() ([]*AosGetStreamingConfigResponse, error) {
	var agscr AosGetStreamingConfigsResponse
	url := o.baseUrl + apiUrlStreamingConfig
	err := o.get(url, []int{200}, &agscr)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	var result []*AosGetStreamingConfigResponse
	for _, i := range agscr.Items {
		result = append(result, &i)
	}
	return result, nil
}

func (o Client) getStreamingConfig(id string) (*AosGetStreamingConfigResponse, error) {
	var result AosGetStreamingConfigResponse
	url := o.baseUrl + apiUrlStreamingConfig + "/" + id
	err := o.get(url, []int{200}, result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	return &result, nil
}

func (o Client) postStreamingConfig(cfg *AosStreamingConfigStreamingEndpoint) (*createStreamingConfigResponse, error) {
	msg, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling AosStreamingConfigStreamingEndpoint object - %v", err)
	}

	var result createStreamingConfigResponse
	url := o.baseUrl + apiUrlStreamingConfig
	err = o.post(url, msg, []int{201}, &result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)

	}

	return &result, nil
}

type NewStreamingConfigCfg struct {
	StreamingType  AosApiStreamingConfigStreamingType
	SequencingMode StreamingConfigSequencingMode
	Protocol       AosApiStreamingConfigProtocol
	Hostname       string
	Port           uint16
}

func (o Client) NewStreamingConfig(in *NewStreamingConfigCfg) (*StreamingConfigId, error) {
	cfg := AosStreamingConfigStreamingEndpoint{
		StreamingType:  in.StreamingType.String(),
		SequencingMode: in.SequencingMode.String(),
		Protocol:       in.Protocol.String(),
		Hostname:       in.Hostname,
		Port:           in.Port,
	}
	response, err := o.postStreamingConfig(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error in NewStreamingConfig - %v", err)
	}

	id := StreamingConfigId(response.Id)
	return &id, nil
}
