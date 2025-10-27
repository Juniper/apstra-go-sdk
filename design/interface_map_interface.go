// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/speed"
)

type InterfaceMapInterface struct {
	Name     string                          `json:"name"`
	Roles    LogicalDevicePortRoles          `json:"roles"`
	Position int                             `json:"position"`
	State    enum.InterfaceMapInterfaceState `json:"state"`
	Speed    speed.Speed                     `json:"speed"`
	Setting  struct {
		Param string `json:"param"`
	} `json:"setting"`
	Mapping InterfaceMapInterfaceMapping `json:"mapping"`
}

var (
	_ json.Marshaler   = (*InterfaceMapInterfaceMapping)(nil)
	_ json.Unmarshaler = (*InterfaceMapInterfaceMapping)(nil)
)

type InterfaceMapInterfaceMapping struct {
	DeviceProfilePortID      int  // slice index 0
	DeviceProfileTransformID int  // slice index 1
	DeviceProfileInterfaceID int  // slice index 2
	LogicalDevicePanel       *int // slice index 3
	LogicalDevicePort        *int // slice index 4
}

func (i InterfaceMapInterfaceMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal([]*int{
		&i.DeviceProfilePortID,
		&i.DeviceProfileTransformID,
		&i.DeviceProfileInterfaceID,
		i.LogicalDevicePanel,
		i.LogicalDevicePort,
	})
}

func (i *InterfaceMapInterfaceMapping) UnmarshalJSON(bytes []byte) error {
	var raw []*int
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling interface map interface mapping: %w", err)
	}
	if len(raw) != 5 {
		return fmt.Errorf("expected 5 elements in interface map interface mapping, got %d", len(raw))
	}
	for idx := range 2 {
		if raw[idx] == nil {
			return fmt.Errorf("the first threee interface mapping elements must be be present, but the value at index %d is nil", idx)
		}
	}

	i.DeviceProfilePortID = *raw[0]
	i.DeviceProfileTransformID = *raw[1]
	i.DeviceProfileInterfaceID = *raw[2]
	i.LogicalDevicePanel = raw[3]
	i.LogicalDevicePort = raw[4]

	return nil
}
