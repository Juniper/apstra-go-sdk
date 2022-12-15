package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDesignConfiglets       = apiUrlDesignPrefix + "configlets"
	apiUrlDesignConfigletsPrefix = apiUrlDesignConfiglets + apiUrlPathDelim
	apiUrlDesignConfigletsById   = apiUrlDesignConfigletsPrefix + "%s"
)

type ConfigletGenerator struct {
	ConfigStyle          string `json:"config_style"`
	Section              string `json:"section"`
	TemplateText         string `json:"template_text"`
	NegationTemplateText string `json:"negation_template_text"`
	Filename             string `json:"filename"`
}

type Configlet struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *ConfigletData
}

type ConfigletData struct {
	RefArchs    []string             `json:"ref_archs"`
	Generators  []ConfigletGenerator `json:"generators"`
	DisplayName string               `json:"display_name"`
}

type rawConfiglet struct {
	RefArchs       []string             `json:"ref_archs"`
	Generators     []ConfigletGenerator `json:"generators"`
	CreatedAt      time.Time            `json:"created_at"`
	Id             ObjectId             `json:"id,omitempty"`
	LastModifiedAt time.Time            `json:"last_modified_at"`
	DisplayName    string               `json:"display_name"`
}

func (o *rawConfiglet) polish() *Configlet {
	return &Configlet{
		Id:             o.Id,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &ConfigletData{
			RefArchs:    o.RefArchs,
			Generators:  o.Generators,
			DisplayName: o.DisplayName,
		},
	}
}

type ConfigletRequest ConfigletData

func (o *Client) listAllConfiglets(ctx context.Context) ([]ObjectId, error) {
	response := &struct {
		Items []ObjectId `json:"items"`
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignConfiglets,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getConfiglet(ctx context.Context, id ObjectId) (*rawConfiglet, error) {
	response := &rawConfiglet{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignConfigletsById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getConfigletByName(ctx context.Context, name string) (*rawConfiglet, error) {
	cgs, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	for _, t := range cgs {
		if t.DisplayName == name {
			return &t, nil
		}
	}

	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf(" Configlet with name '%s' not found", name),
	}
}

func (o *Client) getAllConfiglets(ctx context.Context) ([]rawConfiglet, error) {
	response := &struct {
		Items []rawConfiglet `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignConfiglets,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) createConfiglet(ctx context.Context, in *ConfigletRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignConfiglets,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateConfiglet(ctx context.Context, id ObjectId, in *ConfigletRequest) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignConfigletsById, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteConfiglet(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignConfigletsById, id),
	})
}
