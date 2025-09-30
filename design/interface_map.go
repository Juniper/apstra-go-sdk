// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ internal.IDer     = (*InterfaceMap)(nil)
	_ json.Marshaler    = (*InterfaceMap)(nil)
	_ json.Unmarshaler  = (*InterfaceMap)(nil)
	_ timeutils.Stamper = (*InterfaceMap)(nil)
)

type InterfaceMap struct {
	Label           string
	DeviceProfileID string
	LogicalDeviceID string
	Interfaces      []InterfaceMapInterface

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (i InterfaceMap) ID() *string {
	if i.id == "" {
		return nil
	}
	return &i.id
}

// SetID sets a the value returned by ID only if it was previously un-set. An
// error is returned If the value was previously set. Presence of an existing
// value is the only reason SetID will return an error. If the value is known to
// be empty, use MustSetID.
func (i *InterfaceMap) SetID(id string) error {
	if i.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", i.id))
	}

	i.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (i *InterfaceMap) MustSetID(id string) {
	err := i.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (i InterfaceMap) CreatedAt() *time.Time {
	return i.createdAt
}

func (i InterfaceMap) LastModifiedAt() *time.Time {
	return i.lastModifiedAt
}

func (i InterfaceMap) MarshalJSON() ([]byte, error) {
	raw := rawInterfaceMap{
		Label:           i.Label,
		DeviceProfileID: i.DeviceProfileID,
		LogicalDeviceID: i.LogicalDeviceID,
		Interfaces:      i.Interfaces,
		ID:              i.id,
	}
	return json.Marshal(raw)
}

func (i *InterfaceMap) UnmarshalJSON(bytes []byte) error {
	var raw rawInterfaceMap
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshalling interface map: %w", err)
	}

	i.Label = raw.Label
	i.DeviceProfileID = raw.DeviceProfileID
	i.LogicalDeviceID = raw.LogicalDeviceID
	i.Interfaces = raw.Interfaces
	i.id = raw.ID
	i.createdAt = raw.CreatedAt
	i.lastModifiedAt = raw.LastModifiedAt

	return nil
}

type rawInterfaceMap struct {
	Label           string                  `json:"label"`
	DeviceProfileID string                  `json:"device_profile_id"`
	LogicalDeviceID string                  `json:"logical_device_id"`
	Interfaces      []InterfaceMapInterface `json:"interfaces"`

	ID             string     `json:"id,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	LastModifiedAt *time.Time `json:"last_modified_at,omitempty"`
}
