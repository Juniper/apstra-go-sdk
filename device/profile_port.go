// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/speed"
)

type Port struct {
	ConnectorType   string           `json:"connector_type"`
	Panel           int              `json:"panel_id"`
	Transformations []Transformation `json:"transformations"`
	Column          int              `json:"column_id"`
	ID              int              `json:"port_id"`
	Row             int              `json:"row_id"`
	FailureDomain   int              `json:"failure_domain_id"`
	Display         *int             `json:"display_id"`
	Slot            int              `json:"slot_id"`
}

// TransformationCandidates takes an interface name ("xe-0/0/1:1") and a speed,
// and returns a map[int]Transformation populated with candidate transformations
// available according to the PortInfo and keyed by the transformation ID. Only
// "active" transformations are returned.
func (p Port) TransformationCandidates(ifName string, ifSpeed speed.Speed) map[int]Transformation {
	result := make(map[int]Transformation)
	for _, transformation := range p.Transformations {
		for _, intf := range transformation.Interfaces {
			if intf.Name == ifName &&
				intf.State == enum.InterfaceStateActive &&
				intf.Speed.Equal(ifSpeed) {
				result[transformation.ID] = transformation
			}
		}
	}
	if len(result) == 0 {
		result = nil
	}
	return result
}
