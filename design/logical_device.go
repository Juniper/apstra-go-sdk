// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"hash"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ internal.IDer     = (*LogicalDevice)(nil)
	_ json.Marshaler    = (*LogicalDevice)(nil)
	_ json.Unmarshaler  = (*LogicalDevice)(nil)
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

// Replicate returns a copy of itself with zero values for metadata fields
func (l LogicalDevice) Replicate() LogicalDevice {
	return LogicalDevice{
		Label:  l.Label,
		Panels: l.Panels,
	}
}

func (l LogicalDevice) CreatedAt() *time.Time {
	return l.createdAt
}

func (l LogicalDevice) LastModifiedAt() *time.Time {
	return l.lastModifiedAt
}

func (l LogicalDevice) MarshalJSON() ([]byte, error) {
	raw := struct {
		ID     string               `json:"id,omitempty"` // ID must be marshaled for rack-type embedding
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

func (l LogicalDevice) digest(h hash.Hash) []byte {
	h.Reset()
	return mustHashForComparison(l, h)
}

func (l *LogicalDevice) setHashID(h hash.Hash) {
	if l.id != "" {
		panic(fmt.Sprintf("id already has value %q", l.id))
	}

	l.id = fmt.Sprintf("%x", l.digest(h))
	return
}

func NewLogicalDevice(id string) LogicalDevice {
	return LogicalDevice{id: id}
}
