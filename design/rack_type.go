// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
)

//var _ json.Marshaler = (*RackType)(nil)
//var _ json.Unmarshaler = (*RackType)(nil)

type RackType struct {
	ID                       *string
	Label                    string
	Description              string
	FabricConnectivityDesign enum.FabricConnectivityDesign
	Status                   *enum.FFEConsistencyStatus
	LeafSwitches             []LeafSwitch `json:"leafs"`
	//AccessSwitches           []AccessSwitch  `json:"access_switches"`
	//GenericSystems           []GenericSystem `json:"generic_systems"`
	CreatedAt      *time.Time `json:"created_at"`
	LastModifiedAt *time.Time `json:"last_modified_at"`
}
