package aosSdk

import "fmt"

const (
	apiUrlStreamingConfig       = "/api/streaming-config"
	apiUrlStreamingConfigPrefix = apiUrlStreamingConfig + "/"
	iotaStringTypesMax          = 50
)

const (
	StreamingConfigSequencingModeSequenced   = "sequenced"
	StreamingConfigSequencingModeUnsequenced = "unsequenced"

	StreamingConfigStreamingTypeAlerts  = "alerts"
	StreamingConfigStreamingTypeEvents  = "events"
	StreamingConfigStreamingTypePerfmon = "perfmon"

	StreamingConfigProtocolProtoBufOverTcp = "protoBufOverTcp"
)

type getStreamingConfigsResponse struct {
	Items []StreamingConfigInfo `json:"items"`
}

// StreamingConfigInfo is returned by Apstra in response to
// both:
//  - 'GET apiUrlStreamingConfig' (as a member of list 'Items')
//  - 'GET apiUrlStreamingConfigPrefix + {id}'
type StreamingConfigInfo struct {
	Status struct {
		ConnectionLog []struct {
			Date    string `json:"date"`
			Message string `json:"message"`
		} `json:"connectionLog"`
		ConnectionTime       string                `json:"connectionTime"`
		Epoch                string                `json:"epoch"`
		ConnectionResetCount int                   `json:"connectionResetCount"`
		StreamingEndpoint    StreamingConfigParams `json:"streamingEndpoint"`
		DnsLog               []struct {
			Date    string `json:"date"`
			Message string `json:"message"`
		} `json:"dnsLog"`
		Connected         bool   `json:"connected"`
		DisconnectionTime string `json:"disconnectionTime"`
	} `json:"status"`
	StreamingType  string   `json:"streaming_type"`
	SequencingMode string   `json:"sequencing_mode"`
	Protocol       string   `json:"protocol"`
	Hostname       string   `json:"hostname"`
	Id             ObjectId `json:"id"`
	Port           uint16   `json:"port"`
}

// StreamingConfigParams is the minimally required description needed to create,
// compare, and look up an Apstra streaming config / receiver.
type StreamingConfigParams struct {
	StreamingType  string `json:"streaming_type"`
	SequencingMode string `json:"sequencing_mode"`
	Protocol       string `json:"protocol"`
	Hostname       string `json:"hostname"`
	Port           uint16 `json:"port"`
}

func (o Client) getAllStreamingConfigs() ([]StreamingConfigInfo, error) {
	var gscr getStreamingConfigsResponse
	_, err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlStreamingConfig,
		toServerPtr:   nil,
		fromServerPtr: &gscr,
	})
	if err != nil {
		return nil, err
	}
	var result []StreamingConfigInfo
	for i := range gscr.Items {
		result = append(result, gscr.Items[i])
	}

	return result, nil
}

func (o Client) getAllStreamingConfigIds() ([]ObjectId, error) {
	all, err := o.getAllStreamingConfigs()
	if err != nil {
		return nil, err
	}

	var result []ObjectId
	for _, i := range all {
		result = append(result, i.Id)
	}

	return result, nil
}

func (o Client) getStreamingConfig(id ObjectId) (*StreamingConfigInfo, error) {
	var result StreamingConfigInfo
	_, err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlStreamingConfigPrefix + string(id),
		fromServerPtr: &result,
	})
	return &result, err
}

func (o Client) newStreamingConfig(cfg *StreamingConfigParams) (ObjectId, error) {
	var result objectIdResponse
	_, err := o.talkToAos(&talkToAosIn{
		method:        httpMethodPost,
		url:           apiUrlStreamingConfig,
		toServerPtr:   cfg,
		fromServerPtr: &result,
	})
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (o Client) deleteStreamingConfig(id ObjectId) error {
	_, err := o.talkToAos(&talkToAosIn{
		method: httpMethodDelete,
		url:    apiUrlStreamingConfig + "/" + string(id),
	})
	return err
}

// GetStreamingConfigIDByCfg checks current StreamingConfigs (Streaming
// Receivers) against the supplied StreamingConfigInfo. If the stream seems
// to already exist on the AOS server, the returned ObjectId will be
// populated. If not found, it will be empty.
func (o Client) GetStreamingConfigIDByCfg(in *StreamingConfigParams) (ObjectId, error) {
	all, err := o.getAllStreamingConfigs()
	if err != nil {
		return "", fmt.Errorf("error getting streaming configs - %w", err)
	}
	for _, scInfo := range all {
		testParams := streamingConfigParamsFromStreamingConfigInfo(&scInfo)
		if CompareStreamingConfigs(testParams, in) {
			return scInfo.Id, nil
		}
	}
	return "", nil
}

func streamingConfigParamsFromStreamingConfigInfo(in *StreamingConfigInfo) *StreamingConfigParams {
	return &StreamingConfigParams{
		StreamingType:  in.StreamingType,
		SequencingMode: in.SequencingMode,
		Protocol:       in.Protocol,
		Hostname:       in.Hostname,
		Port:           in.Port,
	}
}

// CompareStreamingConfigs returns true if the supplied StreamingConfigInfo
// objects are likely to be recognized as a collision
// (ErrStringStreamingConfigExists) by the AOS API.
func CompareStreamingConfigs(a *StreamingConfigParams, b *StreamingConfigParams) bool {
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
