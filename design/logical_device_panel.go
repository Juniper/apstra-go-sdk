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

var (
	_ json.Marshaler   = (*LogicalDevicePanel)(nil)
	_ json.Unmarshaler = (*LogicalDevicePanel)(nil)
)

type LogicalDevicePanel struct {
	PanelLayout  LogicalDevicePanelLayout
	PortGroups   []LogicalDevicePanelPortGroup
	PortIndexing enum.DesignLogicalDevicePanelPortIndexing
}

func (l LogicalDevicePanel) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawLogicalDevicePanel{
		PanelLayout: l.PanelLayout,
		PortIndexing: rawLogicalDevicePanelPortIndexing{
			Order:      &l.PortIndexing,
			Schema:     "absolute",
			StartIndex: 1,
		},
		PortGroups: l.PortGroups,
	})
}

func (l *LogicalDevicePanel) UnmarshalJSON(bytes []byte) error {
	var raw rawLogicalDevicePanel
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling LogicalDevicePanel: %w", err)
	}

	l.PanelLayout = raw.PanelLayout
	l.PortIndexing = *raw.PortIndexing.Order
	l.PortGroups = raw.PortGroups

	return nil
}

// it is safe and reasonable to have a "raw" type for objects which:
// 1) are marshaled and unmarshaled symmetrically (have no metadata to suppress)
// 2) have JSON layout which doesn't align with their public struct layout
type rawLogicalDevicePanelPortIndexing struct {
	Order      *enum.DesignLogicalDevicePanelPortIndexing `json:"order"`
	Schema     string                                     `json:"schema"`
	StartIndex int                                        `json:"start_index"`
}

// it is safe and reasonable to have a "raw" type for objects which:
// 1) are marshaled and unmarshaled symmetrically (have no metadata to suppress)
// 2) have JSON layout which doesn't align with their public struct layout
type rawLogicalDevicePanel struct {
	PanelLayout  LogicalDevicePanelLayout          `json:"panel_layout"`
	PortIndexing rawLogicalDevicePanelPortIndexing `json:"port_indexing"`
	PortGroups   []LogicalDevicePanelPortGroup     `json:"port_groups"`
}

type LogicalDevicePanelPortGroup struct {
	Count int                    `json:"count"`
	Speed speed.Speed            `json:"speed"`
	Roles LogicalDevicePortRoles `json:"roles"`
}

type LogicalDevicePanelLayout struct {
	RowCount    int `json:"row_count"`
	ColumnCount int `json:"column_count"`
}

type LogicalDevicePortRoles []enum.PortRole

func (o LogicalDevicePortRoles) Strings() []string {
	result := make([]string, len(o))
	for i, pr := range o {
		result[i] = pr.String()
	}

	return result
}

func (o *LogicalDevicePortRoles) FromStrings(in []string) error {
	newPRs := make(LogicalDevicePortRoles, len(in))
	for i, s := range in {
		err := newPRs[i].FromString(s)
		if err != nil {
			return err
		}
	}
	*o = newPRs

	return nil
}

// SetAllRoles ensures that the LogicalDevicePortRoles contains the entire
// set of "available for use" port roles: All roles excluding "l3_server"
// (deprecated)
func (o *LogicalDevicePortRoles) SetAllRoles() {
	// wipe out any existing values
	*o = nil

	for _, member := range enum.PortRoles.Members() {
		switch member {
		case enum.PortRoleL3Server: // don't add this one
		default:
			*o = append(*o, member) // this one's a keeper
		}
	}
}

func (o LogicalDevicePortRoles) Validate() error {
	for _, ldpr := range o {
		if ldpr == enum.PortRoleL3Server {
			return fmt.Errorf("logical device port role %q is no longer supported", ldpr.String())
		}
	}

	return nil
}
