// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ internal.IDer     = (*ConfigTemplate)(nil)
	_ json.Unmarshaler  = (*ConfigTemplate)(nil)
	_ timeutils.Stamper = (*ConfigTemplate)(nil)
)

type ConfigTemplate struct {
	Label      string `json:"label"`
	Predefined bool   `json:"predefined"`
	Text       string `json:"text"`

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (c ConfigTemplate) ID() *string {
	if c.id == "" {
		return nil
	}
	return &c.id
}

func (c ConfigTemplate) CreatedAt() *time.Time {
	return c.createdAt
}

func (c ConfigTemplate) LastModifiedAt() *time.Time {
	return c.lastModifiedAt
}

func (c *ConfigTemplate) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string     `json:"id"`
		Label          string     `json:"label"`
		Predefined     bool       `json:"predefined"`
		Text           string     `json:"text"`
		CreatedAt      *time.Time `json:"created_at"`
		LastModifiedAt *time.Time `json:"last_modified_at"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling logical device: %w", err)
	}

	c.id = raw.ID
	c.Label = raw.Label
	c.Predefined = raw.Predefined
	c.Text = raw.Text
	c.createdAt = raw.CreatedAt
	c.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func NewConfigTemplate(id string) ConfigTemplate {
	return ConfigTemplate{id: id}
}
