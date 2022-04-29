package aosSdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	aosApiStreamingConfig = "/api/streaming-config"

	AosApiStreamingConfigSequencingTypeUnknown AosApiStreamingConfigSequencingType = iota
	AosApiStreamingConfigSequencingTypeSequenced
	AosApiStreamingConfigSequencingTypeUnsequenced

	AosApiStreamingConfigMessageTypeUnknown AosApiStreamingConfigMessageType = iota
	AosApiStreamingConfigMessageTypeAlerts
	AosApiStreamingConfigMessageTypeEvents
	AosApiStreamingConfigMessageTypePerfmon

	AosApiStreamingConfigProtocolTypeUnknown AosApiStreamingConfigProtocol = iota
	AosApiStreamingConfigProtocolTypeProtoBufTcp
)

type AosApiStreamingConfigId string

type AosApiStreamingConfigSequencingType int

func (o AosApiStreamingConfigSequencingType) String() string {
	switch o {
	case AosApiStreamingConfigSequencingTypeUnknown:
		return "unknown"
	case AosApiStreamingConfigSequencingTypeSequenced:
		return "sequenced"
	case AosApiStreamingConfigSequencingTypeUnsequenced:
		return "unsequenced"
	default:
		return fmt.Sprintf("sequencing type %d has no string value", o)
	}
}

type AosApiStreamingConfigMessageType int

func (o AosApiStreamingConfigMessageType) String() string {
	switch o {
	case AosApiStreamingConfigMessageTypeUnknown:
		return "unknown"
	case AosApiStreamingConfigMessageTypeAlerts:
		return "alerts"
	case AosApiStreamingConfigMessageTypeEvents:
		return "events"
	case AosApiStreamingConfigMessageTypePerfmon:
		return "perfmon"
	default:
		return fmt.Sprintf("message type %d has no string value", o)
	}
}

type AosApiStreamingConfigProtocol int

func (o AosApiStreamingConfigProtocol) String() string {
	switch o {
	case AosApiStreamingConfigProtocolTypeUnknown:
		return "unknown"
	case AosApiStreamingConfigProtocolTypeProtoBufTcp:
		return "protoBufOverTcp"
	default:
		return fmt.Sprintf("message type %d has no string value", o)
	}
}

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

func (o AosClient) getAllStreamingConfigs() ([]*AosStreamingConfig, error) {
	var agscr AosGetStreamingConfigsResponse
	url := o.baseUrl + aosApiStreamingConfig
	err := o.get(url, []int{200}, &agscr)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	var result []*AosStreamingConfig
	for _, i := range agscr.Items {
		result = append(result, &i)
	}
	return result, nil
}

func (o AosClient) getStreamingConfig(id string) (*AosStreamingConfig, error) {
	var result AosStreamingConfig
	url := o.baseUrl + aosApiStreamingConfig + "/" + id
	err := o.get(url, []int{200}, result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)
	}
	return &result, nil
}

func (o AosClient) postStreamingConfig(cfg *AosStreamingConfigStreamingEndpoint) (*aosCreateStreamingConfigResponse, error) {
	msg, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling AosStreamingConfigStreamingEndpoint object - %v", err)
	}

	//todo: use post method here

	var result aosCreateStreamingConfigResponse
	url := o.baseUrl + aosApiStreamingConfig
	err = o.post(url, msg, []int{201}, &result)
	if err != nil {
		return nil, fmt.Errorf("error calling %s - %v", url, err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, fmt.Errorf("error creating http Request - %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authtoken", o.login.Token)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling http.client.Do - %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("http response code is not '%d' got '%d' at '%s'", 201, resp.StatusCode, aosApiStreamingConfig)
	}

	var createStreamingConfigResp *aosCreateStreamingConfigResponse
	err = json.NewDecoder(resp.Body).Decode(&createStreamingConfigResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding aosCreateStreamingConfigResponse JSON - %v", err)
	}

	return createStreamingConfigResp, nil
}

type NewStreamingReceiverIn struct {
	StreamType     AosApiStreamingConfigSequencingType
	MessageType    AosApiStreamingConfigMessageType
	SequencingMode AosApiStreamingConfigSequencingType
	Protocol       AosApiStreamingConfigProtocol
	Hostname       string
	Port           uint16
}

func (o AosClient) NewStreamingReceiver(in *NewStreamingReceiverIn) (*AosApiStreamingConfigId, error) {
	cfg := AosStreamingConfigStreamingEndpoint{
		StreamingType:  in.StreamType.String(),
		SequencingMode: in.SequencingMode.String(),
		Protocol:       in.Protocol.String(),
		Hostname:       in.Hostname,
		Port:           in.Port,
	}
	response, err := o.postStreamingConfig(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating NewStreamingReceiver - %v", err)
	}

	id := AosApiStreamingConfigId(response.Id)
	return &id, nil
}
