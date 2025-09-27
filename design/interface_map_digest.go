// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"time"

	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ json.Unmarshaler  = (*InterfaceMapDigest)(nil)
	_ timeutils.Stamper = (*InterfaceMapDigest)(nil)
)

type InterfaceMapDigest struct {
	Label              string
	DeviceProfileID    string
	DeviceProfileLabel string
	LogicalDeviceID    string
	LogicalDeviceLabel string

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (i InterfaceMapDigest) ID() *string {
	if i.id == "" {
		return nil
	}
	return &i.id
}

func (i InterfaceMapDigest) CreatedAt() *time.Time {
	return i.createdAt
}

func (i InterfaceMapDigest) LastModifiedAt() *time.Time {
	return i.lastModifiedAt
}

func (i *InterfaceMapDigest) UnmarshalJSON(bytes []byte) error {
	type idLabel struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	}
	var raw struct {
		Label          string     `json:"label"`
		DeviceProfile  idLabel    `json:"device_profile"`
		LogicalDevice  idLabel    `json:"logical_device"`
		ID             string     `json:"id"`
		CreatedAt      *time.Time `json:"created_at"`
		LastModifiedAt *time.Time `json:"last_modified_at"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling interface map digest: %w", err)
	}

	i.Label = raw.Label
	i.DeviceProfileID = raw.DeviceProfile.ID
	i.DeviceProfileLabel = raw.DeviceProfile.Label
	i.LogicalDeviceID = raw.LogicalDevice.ID
	i.LogicalDeviceLabel = raw.LogicalDevice.Label
	i.id = raw.ID
	i.createdAt = raw.CreatedAt
	i.lastModifiedAt = raw.LastModifiedAt

	return nil
}
