package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDesignLogicalDevices       = apiUrlDesignPrefix + "logical-devices"
	apiUrlDesignLogicalDevicesPrefix = apiUrlDesignLogicalDevices + apiUrlPathDelim
	apiUrlDesignLogicalDeviceById    = apiUrlDesignLogicalDevicesPrefix + "%s"

	PortIndexingVerticalFirst   = "L-R, T-B"
	PortIndexingHorizontalFirst = "T-B, L-R"
	PortIndexingSchemaAbsolute  = "absolute"
)

type optionsLogicalDevicesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type LogicalDevicePanelLayout struct {
	RowCount    int `json:"row_count"`
	ColumnCount int `json:"column_count"`
}

type LogicalDevicePortIndexing struct {
	Order      string `json:"order"`
	StartIndex int    `json:"start_index"`
	Schema     string `json:"schema"` // Valid choices: absolute
}

type LogicalDevicePortGroup struct {
	Count int                    `json:"count"`
	Speed LogicalDevicePortSpeed `json:"speed"`
	Roles []string               `json:"roles"` // Valid choices: spine, leaf, l3_server, server, access, peer, unused, superspine, generic,
}

type LogicalDevicePortSpeed struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

type LogicalDevicePanel struct {
	PanelLayout  LogicalDevicePanelLayout  `json:"panel_layout"`
	PortIndexing LogicalDevicePortIndexing `json:"port_indexing"`
	PortGroups   []LogicalDevicePortGroup  `json:"port_groups"`
}

type LogicalDevice struct {
	DisplayName    string               `json:"display_name"`
	Id             ObjectId             `json:"id,omitempty"`
	Panels         []LogicalDevicePanel `json:"panels"`
	CreatedAt      time.Time            `json:"created_at"`
	LastModifiedAt time.Time            `json:"last_modified_at"`
}

func (o *Client) listLogicalDeviceIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsLogicalDevicesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignLogicalDevices,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getLogicalDevice(ctx context.Context, id ObjectId) (*LogicalDevice, error) {
	response := &LogicalDevice{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
		apiResponse: response,
	})
	return response, convertTtaeToAceWherePossible(err)
}

func (o *Client) createLogicalDevice(ctx context.Context, in *LogicalDevice) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignLogicalDevices,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateLogicalDevice(ctx context.Context, id ObjectId, in *LogicalDevice) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
		apiInput: in,
	})
}

func (o *Client) deleteLogicalDevice(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
	})
}
