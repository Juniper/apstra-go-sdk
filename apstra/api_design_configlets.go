// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	apiUrlDesignConfiglets       = apiUrlDesignPrefix + "configlets"
	apiUrlDesignConfigletsPrefix = apiUrlDesignConfiglets + apiUrlPathDelim
	apiUrlDesignConfigletsById   = apiUrlDesignConfigletsPrefix + "%s"
)

type ConfigletGenerator struct {
	ConfigStyle          enum.ConfigletStyle   `json:"config_style"`
	Section              enum.ConfigletSection `json:"section"`
	SectionCondition     string                `json:"section_condition,omitempty"`
	TemplateText         string                `json:"template_text"`
	NegationTemplateText string                `json:"negation_template_text"`
	Filename             string                `json:"filename"`
}

var _ json.Unmarshaler = (*Configlet)(nil)

type Configlet struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *ConfigletData
}

func (o *Configlet) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		RefArchs       []enum.RefDesign     `json:"ref_archs"`
		Generators     []ConfigletGenerator `json:"generators"`
		CreatedAt      time.Time            `json:"created_at"`
		Id             ObjectId             `json:"id"`
		LastModifiedAt time.Time            `json:"last_modified_at"`
		DisplayName    string               `json:"display_name"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.CreatedAt = raw.CreatedAt
	o.LastModifiedAt = raw.LastModifiedAt
	o.Data = &ConfigletData{
		RefArchs:    raw.RefArchs,
		Generators:  raw.Generators,
		DisplayName: raw.DisplayName,
	}

	return nil
}

type ConfigletData struct {
	RefArchs    []enum.RefDesign     `json:"ref_archs"`
	Generators  []ConfigletGenerator `json:"generators"`
	DisplayName string               `json:"display_name"`
}

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

func (o *Client) getConfiglet(ctx context.Context, id ObjectId) (*Configlet, error) {
	response := &Configlet{}
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

func (o *Client) getConfigletByName(ctx context.Context, name string) (*Configlet, error) {
	configlets, err := o.getAllConfiglets(ctx)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	foundIdx := -1
	for i, configlet := range configlets {
		if configlet.Data.DisplayName == name {
			if foundIdx >= 0 {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple Configlets have name %q", name),
				}
			}
			foundIdx = i
		}
	}

	if foundIdx >= 0 {
		return &configlets[foundIdx], nil
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no Configlet with name '%s' found", name),
	}
}

func (o *Client) getAllConfiglets(ctx context.Context) ([]Configlet, error) {
	var response struct {
		Items []Configlet `json:"items"`
	}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignConfiglets,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *Client) createConfiglet(ctx context.Context, in *ConfigletData) (ObjectId, error) {
	var response objectIdResponse

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignConfiglets,
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) updateConfiglet(ctx context.Context, id ObjectId, in *ConfigletData) error {
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
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignConfigletsById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
