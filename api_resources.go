package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiUrlResources               = "/api/resources"
	apiUrlResourcesAsnPools       = apiUrlResources + "/asn-pools"
	apiUrlResourcesAsnPoolsPrefix = apiUrlResourcesAsnPools + apiUrlPathDelim
)

type NewAsnRange struct {
	B int64 `json:"first"`
	E int64 `json:"last"`
}

func (o NewAsnRange) String() string {
	return fmt.Sprintf("%d-%d", o.B, o.E)
}

type NewAsnPool struct {
	Ranges      []NewAsnRange `json:"ranges"`
	DisplayName string        `json:"display_name"`
}

type AsnPool struct {
	Status         string    `json:"status"`
	Used           string    `json:"used"`
	DisplayName    string    `json:"display_name"`
	Tags           []string  `json:"tags"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	Ranges         []struct {
		Status         string  `json:"status"`
		Used           string  `json:"used"`
		Last           int64   `json:"last"`
		UsedPercentage float64 `json:"used_percentage"`
		Total          string  `json:"total"`
		First          int64   `json:"first"`
	} `json:"ranges"`
	UsedPercentage float64  `json:"used_percentage"`
	Total          string   `json:"total"`
	Id             ObjectId `json:"id"`
}

type getAsnPoolsResponse struct {
	Items []AsnPool `json:"items"`
}

func (o *Client) createAsnPool(ctx context.Context, in *NewAsnPool) (*objectIdResponse, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPools)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}
	response := &objectIdResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         apstraUrl,
		apiInput:    in,
		apiResponse: response,
	})
}

// todo: move to client.go
func (o *Client) GetAsnPools(ctx context.Context) ([]AsnPool, error) {
	return o.getAsnPools(ctx)
}

// todo: move to client.go
func (o *Client) CreateAsnPool(ctx context.Context, in *NewAsnPool) (ObjectId, error) {
	response, err := o.createAsnPool(ctx, in)
	if err != nil {
		return "", fmt.Errorf("error creating ASN pool - %w", err)
	}
	return response.Id, nil
}

// todo: move to client.go
func (o *Client) GetAsnPool(ctx context.Context, in ObjectId) (*AsnPool, error) {
	return o.getAsnPool(ctx, in)
}

// todo: move to client.go
func (o *Client) DeleteAsnPool(ctx context.Context, in ObjectId) error {
	return o.deleteAsnPool(ctx, in)
}

func (o *Client) getAsnPools(ctx context.Context) ([]AsnPool, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPools)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPools, err)
	}
	response := &getAsnPoolsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching ASN pools - %w", err)
	}
	return response.Items, nil
}

func (o *Client) getAsnPool(ctx context.Context, in ObjectId) (*AsnPool, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(in))
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPoolsPrefix+string(in), err)
	}
	response := &AsnPool{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching ASN pool '%s' - %w", in, err)
	}
	return response, nil

}

func (o *Client) deleteAsnPool(ctx context.Context, in ObjectId) error {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(in))
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPoolsPrefix+string(in), err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		url:    apstraUrl,
	})
	if err != nil {
		return fmt.Errorf("error fetching ASN pools - %w", err)
	}
	return nil

}
