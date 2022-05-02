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

	StreamingConfigStreamingTypeUnknown StreamingConfigStreamingType = iota
	StreamingConfigStreamingTypeAlerts
	StreamingConfigStreamingTypeEvents
	StreamingConfigStreamingTypePerfmon

	StreamingConfigProtocolUnknown StreamingConfigProtocol = iota
	StreamingConfigProtocolProtoBufOverTcp

	ErrStringStreamingConfigExists   = "Entity already exists"
	ErrStringStreamingConfigNotFound = "streaming config not found"
)

type ErrStreamingConfigExists struct{}

func (o *ErrStreamingConfigExists) Error() string {
	return ErrStringStreamingConfigExists
}

type ErrStreamingConfigNotFound struct{}

func (o ErrStreamingConfigNotFound) Error() string {
	return ErrStringStreamingConfigNotFound
}

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

type StreamingConfigStreamingType int

func (o StreamingConfigStreamingType) String() string {
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
		return fmt.Sprintf("streaming type %d has no string value", o)
	}
}

type StreamingConfigProtocol int

func (o StreamingConfigProtocol) String() string {
	switch o {
	case StreamingConfigProtocolUnknown:
		return "unknown"
	case StreamingConfigProtocolProtoBufOverTcp:
		return "protoBufOverTcp"
	default:
		return fmt.Sprintf("message type %d has no string value", o)
	}
}

type getStreamingConfigsResponse struct {
	Items []StreamingConfigCfg `json:"items"`
}

type StreamingConfigCfg struct {
	Status         StreamingConfigStatus         `json:"status"`
	StreamingType  StreamingConfigStreamingType  `json:"streaming_type"`
	SequencingMode StreamingConfigSequencingMode `json:"sequencing_mode"`
	Protocol       StreamingConfigProtocol       `json:"protocol"`
	Hostname       string                        `json:"hostname"`
	Id             StreamingConfigId             `json:"id"`
	Port           uint16                        `json:"port"`
}

type StreamingConfigStatus struct {
	ConnectionLog        []StreamingConfigConnectionLog   `json:"connectionLog"`
	ConnectionTime       string                           `json:"connectionTime"`
	Epoch                string                           `json:"epoch"`
	ConnectionResetCount uint                             `json:"connnectionResetCount"`
	StreamingEndpoint    StreamingConfigStreamingEndpoint `json:"streamingEndpoint"`
	DnsLog               []StreamingConfigDnsLog          `json:"dnsLog"`
	Connected            bool                             `json:"connected"`
	DisconnectionTime    string                           `json:"disconnectionTime"`
}

type StreamingConfigConnectionLog struct {
	Date    string `json:"date"'`
	Message string `json:"message"`
}

type StreamingConfigStreamingEndpoint struct {
	StreamingType  string `json:"streaming_type"`
	SequencingMode string `json:"sequencing_mode"`
	Protocol       string `json:"protocol"`
	Hostname       string `json:"hostname"`
	Port           uint16 `json:"port"`
}

type StreamingConfigDnsLog struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

type createStreamingConfigResponse struct {
	Id     string `json:"id"`
	Errors string `json:"errors"`
}

func (o Client) getAllStreamingConfigs() ([]StreamingConfigCfg, error) {
	var result getStreamingConfigsResponse
	url := o.baseUrl + apiUrlStreamingConfig
	err := o.get(url, []int{200}, &result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	return result.Items, nil
}

func (o Client) getStreamingConfig(id string) (*StreamingConfigCfg, error) {
	var result StreamingConfigCfg
	url := o.baseUrl + apiUrlStreamingConfig + "/" + id
	err := o.get(url, []int{200}, result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	return &result, nil
}

func (o Client) postStreamingConfig(cfg *StreamingConfigStreamingEndpoint) (*createStreamingConfigResponse, error) {
	msg, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling StreamingConfigStreamingEndpoint object - %v", err)
	}

	var result createStreamingConfigResponse
	url := o.baseUrl + apiUrlStreamingConfig
	err = o.post(url, msg, []int{201, 409}, &result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	if result.Errors == ErrStringStreamingConfigExists {
		return nil, &ErrStreamingConfigExists{}
	}
	if result.Errors != "" {
		return nil, fmt.Errorf("server error calling %s - %v", url, result.Errors)
	}

	return &result, nil
}

// NewStreamingConfig creates a StreamingConfig (Streaming Receiver) on the AOS server.
func (o Client) NewStreamingConfig(in *StreamingConfigCfg) (StreamingConfigId, error) {
	cfg := StreamingConfigStreamingEndpoint{
		StreamingType:  in.StreamingType.String(),
		SequencingMode: in.SequencingMode.String(),
		Protocol:       in.Protocol.String(),
		Hostname:       in.Hostname,
		Port:           in.Port,
	}
	response, err := o.postStreamingConfig(&cfg)
	if err != nil {
		return "", fmt.Errorf("error in NewStreamingConfig - %v", err)
	}

	id := StreamingConfigId(response.Id)
	return id, nil
}

// DeleteStreamingConfig removes the specified StreamingConfig (Streaming
// Receiver) on the Aos server.
func (o Client) DeleteStreamingConfig(id StreamingConfigId) error {
	return o.talkToAos(&talkToAosIn{
		method: httpMethodDelete,
		url:    apiUrlStreamingConfig + "/" + string(id),
	})
}

// GetStreamingConfigIDByCfg checks current StreamingConfigs (Streaming
// Receivers) against the supplied StreamingConfigCfg. If the stream seems
// to already exist on the AOS server, the returned StreamingConfigId will be
// populated. If not found, it will be empty.
func (o Client) GetStreamingConfigIDByCfg(in *StreamingConfigCfg) (StreamingConfigId, error) {
	all, err := o.GetStreamingConfigs()
	if err != nil {
		return "", fmt.Errorf("error getting streaming configs - %v", err)
	}
	for _, sc := range all {
		if CompareStreamingConfigs(&sc, in) {
			return sc.Id, nil
		}
	}
	return "", nil
}

// CompareStreamingConfigs returns true if the supplied StreamingConfigCfg
// objects are likely to be recognized as a collision
// (ErrStringStreamingConfigExists) by the AOS API.
func CompareStreamingConfigs(a *StreamingConfigCfg, b *StreamingConfigCfg) bool {
	if a.Hostname != b.Hostname {
		return false
	}
	if a.Port != b.Port {
		return false
	}
	if a.StreamingType != b.StreamingType {
		return false
	}
	return true
}
