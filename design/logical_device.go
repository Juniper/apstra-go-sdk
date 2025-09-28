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
	_ internal.IDer                      = (*LogicalDevice)(nil)
	_ internal.Replicator[LogicalDevice] = (*LogicalDevice)(nil)
	_ json.Marshaler                     = (*LogicalDevice)(nil)
	_ json.Unmarshaler                   = (*LogicalDevice)(nil)
	_ timeutils.Stamper                  = (*LogicalDevice)(nil)
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

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (l *LogicalDevice) SetID(id string) error {
	if l.id != "" {
		return internal.IDIsSet(fmt.Errorf("id already has value %q", l.id))
	}

	l.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (l *LogicalDevice) MustSetID(id string) {
	err := l.SetID(id)
	if err != nil {
		panic(err)
	}
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

func (l *LogicalDevice) setHashID(h hash.Hash) error {
	return l.SetID(fmt.Sprintf("%x", l.digest(h)))
}

func (l *LogicalDevice) mustSetHashID(h hash.Hash) {
	l.SetID(fmt.Sprintf("%x", l.digest(h)))
}

func NewLogicalDevice(id string) LogicalDevice {
	return LogicalDevice{id: id}
}
