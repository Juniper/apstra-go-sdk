// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal/slice"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

const (
	LogicalDeviceUrl     = urlPrefix + "logical-devices"
	LogicalDeviceUrlByID = LogicalDeviceUrl + "/%s"
)

var (
	_ json.Marshaler    = (*LogicalDevice)(nil)
	_ json.Unmarshaler  = (*LogicalDevice)(nil)
	_ slice.IDer        = (*LogicalDevice)(nil)
	_ timeutils.Stamper = (*LogicalDevice)(nil)
)

type LogicalDevice struct {
	Label  string
	Panels []LogicalDevicePanel

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (l LogicalDevice) ID() *string {
	if l.id == "" {
		return nil
	}
	return &l.id
}

func (l LogicalDevice) CreatedAt() *time.Time {
	return l.createdAt
}

func (l LogicalDevice) LastModifiedAt() *time.Time {
	return l.lastModifiedAt
}

func (l LogicalDevice) MarshalJSON() ([]byte, error) {
	raw := struct {
		ID     string               `json:"id,omitempty"`
		Label  string               `json:"display_name"`
		Panels []LogicalDevicePanel `json:"panels"`
	}{
		ID:     l.id,
		Label:  l.Label,
		Panels: l.Panels,
	}
	return json.Marshal(raw)
}

func (l *LogicalDevice) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string               `json:"id"`
		Label          string               `json:"display_name"`
		Panels         []LogicalDevicePanel `json:"panels"`
		CreatedAt      *time.Time           `json:"created_at"`
		LastModifiedAt *time.Time           `json:"last_modified_at"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling logical device: %w", err)
	}

	l.id = raw.ID
	l.Label = raw.Label
	l.Panels = raw.Panels
	l.createdAt = raw.CreatedAt
	l.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func NewLogicalDevice(id string) LogicalDevice {
	return LogicalDevice{id: id}
}
