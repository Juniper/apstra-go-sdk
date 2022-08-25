package goapstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlStreamingConfig       = "/api/streaming-config"
	apiUrlStreamingConfigPrefix = apiUrlStreamingConfig + apiUrlPathDelim
	apiUrlStreamingConfigById   = apiUrlStreamingConfigPrefix + "%s"
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

type getStreamingConfigsOptionsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
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

func (o *Client) getAllStreamingConfigIds(ctx context.Context) ([]ObjectId, error) {
	result := &getStreamingConfigsOptionsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlStreamingConfig,
		apiInput:    nil,
		apiResponse: result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return result.Items, nil
}

func (o *Client) getAllStreamingConfigs(ctx context.Context) ([]StreamingConfigInfo, error) {
	gscr := &getStreamingConfigsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlStreamingConfig,
		apiInput:    nil,
		apiResponse: gscr,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []StreamingConfigInfo
	for i := range gscr.Items {
		result = append(result, gscr.Items[i])
	}

	return result, nil
}

func (o *Client) getStreamingConfig(ctx context.Context, id ObjectId) (*StreamingConfigInfo, error) {
	result := &StreamingConfigInfo{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlStreamingConfigById, id),
		apiResponse: result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return result, nil
}

func (o *Client) newStreamingConfig(ctx context.Context, cfg *StreamingConfigParams) (*objectIdResponse, error) {
	result := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlStreamingConfig,
		apiInput:    cfg,
		apiResponse: result,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return result, nil
}

func (o *Client) deleteStreamingConfig(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlStreamingConfigById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

// todo make public

// GetAllStreamingConfigIds returns a []ObjectId representing Streaming
// Receivers currently known to Apstra
func (o *Client) GetAllStreamingConfigIds(ctx context.Context) ([]ObjectId, error) {
	return o.getAllStreamingConfigIds(ctx)
}

// todo make public

// GetStreamingConfigIDByCfg checks current StreamingConfigs (Streaming
// Receivers) against the supplied StreamingConfigInfo. If the stream seems
// to already exist on the AOS server, the returned ObjectId will be
// populated. If not found, it will be empty.
func (o *Client) GetStreamingConfigIDByCfg(ctx context.Context, in *StreamingConfigParams) (ObjectId, error) {
	all, err := o.getAllStreamingConfigs(ctx)
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
