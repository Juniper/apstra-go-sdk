package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiUrlDesignLogicalDevices       = apiUrlDesignPrefix + "logical-devices"
	apiUrlDesignLogicalDevicesPrefix = apiUrlDesignLogicalDevices + apiUrlPathDelim
	apiUrlDesignLogicalDeviceById    = apiUrlDesignLogicalDevicesPrefix + "%s"
)

type LogicalDeviceId string

type optionsLogicalDevicesResponse struct {
	Items   []LogicalDeviceId `json:"items"`
	Methods []string          `json:"methods"`
}

type LogicalDevicePanel struct {
	PanelLayout struct {
		RowCount    int `json:"row_count"`
		ColumnCount int `json:"column_count"`
	} `json:"panel_layout"`
	PortIndexing struct {
		Order      string `json:"order"`
		StartIndex int    `json:"start_index"`
		Schema     string `json:"schema"`
	} `json:"port_indexing"`
	PortGroups []struct {
		Count int `json:"count"`
		Speed struct {
			Unit  string `json:"unit"`
			Value int    `json:"value"`
		} `json:"speed"`
		Roles []string `json:"roles"`
	} `json:"port_groups"`
}

type LogicalDevice struct {
	DisplayName    string               `json:"display_name"`
	Id             LogicalDeviceId      `json:"id"`
	CreatedAt      time.Time            `json:"created_at"`
	LastModifiedAt time.Time            `json:"last_modified_at"`
	Panels         []LogicalDevicePanel `json:"panels"`
}

func (o *Client) listLogicalDeviceIds(ctx context.Context) ([]LogicalDeviceId, error) {
	method := http.MethodOptions
	urlStr := apiUrlDesignLogicalDevices
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	response := &optionsLogicalDevicesResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s' - %w", method, urlStr, convertTtaeToAceWherePossible(err))
	}
	return response.Items, nil
}

func (o *Client) getLogicalDevice(ctx context.Context, id LogicalDeviceId) (*LogicalDevice, error) {
	method := http.MethodGet
	urlStr := fmt.Sprintf(apiUrlDesignLogicalDeviceById, id)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	response := &LogicalDevice{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s' - %w", method, urlStr, convertTtaeToAceWherePossible(err))
	}
	return response, nil
}
