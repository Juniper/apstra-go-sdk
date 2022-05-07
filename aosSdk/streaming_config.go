package aosSdk

import (
	"fmt"
)

const (
	apiUrlStreamingConfig       = "/api/streaming-config"
	apiUrlStreamingConfigPrefix = apiUrlStreamingConfig + "/"
	iotaStringTypesMax          = 50
)

const (
	StreamingConfigSequencingModeUnknown StreamingConfigSequencingMode = iota
	StreamingConfigSequencingModeSequenced
	StreamingConfigSequencingModeUnsequenced
	StreamingConfigSequencingModeNotFound = "sequencing type %d has no string value"

	StreamingConfigStreamingTypeUnknown StreamingConfigStreamingType = iota
	StreamingConfigStreamingTypeAlerts
	StreamingConfigStreamingTypeEvents
	StreamingConfigStreamingTypePerfmon
	StreamingConfigStreamingTypeNotFound = "streaming type %d has no string value"

	StreamingConfigProtocolUnknown StreamingConfigProtocol = iota
	StreamingConfigProtocolProtoBufOverTcp
	StreamingConfigProtocolNotFound = "protocol %d has no string value"
)

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
		return fmt.Sprintf(StreamingConfigSequencingModeNotFound, o)
	}
}

func StreamingConfigSequencingModeFromString(in string) StreamingConfigSequencingMode {
	for i := StreamingConfigSequencingModeUnknown; i <= StreamingConfigSequencingModeUnknown+iotaStringTypesMax; i++ {
		switch {
		case i.String() == in:
			return i
		case i.String() == fmt.Sprintf(StreamingConfigSequencingModeNotFound, i):
			return StreamingConfigSequencingModeUnknown
		}
	}
	return StreamingConfigSequencingModeUnknown
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
		return fmt.Sprintf(StreamingConfigStreamingTypeNotFound, o)
	}
}

func StreamingConfigStreamingTypeFromString(in string) StreamingConfigStreamingType {
	for i := StreamingConfigStreamingTypeUnknown; i <= StreamingConfigStreamingTypeUnknown+iotaStringTypesMax; i++ {
		switch {
		case i.String() == in:
			return i
		case i.String() == fmt.Sprintf(StreamingConfigStreamingTypeNotFound, i):
			return StreamingConfigStreamingTypeUnknown
		}
	}
	return StreamingConfigStreamingTypeUnknown
}

type StreamingConfigProtocol int

func (o StreamingConfigProtocol) String() string {
	switch o {
	case StreamingConfigProtocolUnknown:
		return "unknown"
	case StreamingConfigProtocolProtoBufOverTcp:
		return "protoBufOverTcp"
	default:
		return fmt.Sprintf(StreamingConfigProtocolNotFound, o)
	}
}

func StreamingConfigProtocolFromString(in string) StreamingConfigProtocol {
	for i := StreamingConfigProtocolUnknown; i <= StreamingConfigProtocolUnknown+iotaStringTypesMax; i++ {
		switch {
		case i.String() == in:
			return i
		case i.String() == fmt.Sprintf(StreamingConfigProtocolNotFound, i):
			return StreamingConfigProtocolUnknown
		}
	}
	return StreamingConfigProtocolUnknown
}

type getStreamingConfigsResponse struct {
	Items []StreamingConfigInfo `json:"items"`
}

// StreamingConfigInfo is returned by Apstra in response to
// 'GET apiUrlStreamingConfig/{id}'
type StreamingConfigInfo struct {
	Status         StreamingConfigStatus         `json:"status"`
	StreamingType  StreamingConfigStreamingType  `json:"streaming_type"`
	SequencingMode StreamingConfigSequencingMode `json:"sequencing_mode"`
	Protocol       StreamingConfigProtocol       `json:"protocol"`
	Hostname       string                        `json:"hostname"`
	Id             ObjectId                      `json:"id"`
	Port           uint16                        `json:"port"`
}

// StreamingConfigStatus is a member of StreamingConfigInfo which is returned by
// Apstra in response to 'GET apiUrlStreamingConfig/{id}'
type StreamingConfigStatus struct {
	ConnectionLog        []StreamingConfigConnectionLog `json:"connectionLog"`
	ConnectionTime       string                         `json:"connectionTime"`
	Epoch                string                         `json:"epoch"`
	ConnectionResetCount uint                           `json:"connnectionResetCount"`
	StreamingEndpoint    StreamingConfigParams          `json:"streamingEndpoint"`
	DnsLog               []StreamingConfigDnsLog        `json:"dnsLog"`
	Connected            bool                           `json:"connected"`
	DisconnectionTime    string                         `json:"disconnectionTime"`
}

// StreamingConfigConnectionLog is a member of StreamingConfigStatus-StreamingConfigInfo>
type StreamingConfigConnectionLog struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

type StreamingConfigParams struct {
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
	Id ObjectId `json:"id"`
}

func (o Client) getAllStreamingConfigIds() ([]ObjectId, error) {
	var gscr getStreamingConfigsResponse
	err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlStreamingConfig,
		toServerPtr:   nil,
		fromServerPtr: &gscr,
	})
	if err != nil {
		return nil, err
	}

	var result []ObjectId
	for _, i := range gscr.Items {
		result = append(result, i.Id)
	}

	return result, nil
}

func (o Client) getStreamingConfig(id ObjectId) (*StreamingConfigInfo, error) {
	var result StreamingConfigInfo
	return &result, o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlStreamingConfigPrefix + string(id),
		fromServerPtr: &result,
	})
}

func (o Client) postStreamingConfig(cfg *StreamingConfigParams) (ObjectId, error) {
	var result createStreamingConfigResponse
	return result.Id, o.talkToAos(&talkToAosIn{
		method:        httpMethodPost,
		url:           apiUrlStreamingConfig,
		toServerPtr:   cfg,
		fromServerPtr: &result,
	})
}

// NewStreamingConfig creates a StreamingConfig (Streaming Receiver) on the AOS server.
func (o Client) NewStreamingConfig(in *StreamingConfigInfo) (ObjectId, error) {
	cfg := StreamingConfigParams{
		StreamingType:  in.StreamingType.String(),
		SequencingMode: in.SequencingMode.String(),
		Protocol:       in.Protocol.String(),
		Hostname:       in.Hostname,
		Port:           in.Port,
	}
	return o.postStreamingConfig(&cfg)
}

// DeleteStreamingConfig removes the specified StreamingConfig (Streaming
// Receiver) on the Aos server.
func (o Client) DeleteStreamingConfig(id ObjectId) error {
	return o.talkToAos(&talkToAosIn{
		method: httpMethodDelete,
		url:    apiUrlStreamingConfig + "/" + string(id),
	})
}

// todo restore this function
//// GetStreamingConfigIDByCfg checks current StreamingConfigs (Streaming
//// Receivers) against the supplied StreamingConfigInfo. If the stream seems
//// to already exist on the AOS server, the returned ObjectId will be
//// populated. If not found, it will be empty.
//func (o Client) GetStreamingConfigIDByCfg(in *StreamingConfigInfo) (ObjectId, error) {
//	all, err := o.GetStreamingConfigs()
//	if err != nil {
//		return "", fmt.Errorf("error getting streaming configs - %w", err)
//	}
//	for _, sc := range all {
//		if CompareStreamingConfigs(&sc, in) {
//			return sc.Id, nil
//		}
//	}
//	return "", nil
//}

// CompareStreamingConfigs returns true if the supplied StreamingConfigInfo
// objects are likely to be recognized as a collision
// (ErrStringStreamingConfigExists) by the AOS API.
func CompareStreamingConfigs(a *StreamingConfigInfo, b *StreamingConfigInfo) bool {
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
