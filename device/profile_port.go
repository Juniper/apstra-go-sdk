// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"fmt"

	sdk "github.com/Juniper/apstra-go-sdk"
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

// DefaultTransform returns the Transformation flagged as default.
// If none are default, an error is returned.
func (o Port) DefaultTransform() (Transformation, error) {
	for _, t := range o.Transformations {
		if t.IsDefault {
			return t, nil
		}
	}

	return Transformation{}, sdk.ErrNotFound(fmt.Sprintf("Port %d has no default transformation", o.ID))
}

// transformationCandidates takes an interface name ("xe-0/0/1:1") and a speed,
// and returns a []Transformation populated with candidate transformations
// available according to the PortInfo and keyed by the transformation ID. Only
// "active" transformations are returned.
func (p Port) transformationCandidates(ifName string, ifSpeed speed.Speed) []Transformation {
	var result []Transformation
	for _, transformation := range p.Transformations {
		for _, intf := range transformation.Interfaces {
			if intf.Name == ifName &&
				intf.State == enum.InterfaceStateActive &&
				intf.Speed.Equal(ifSpeed) {
				result = append(result, transformation)
			}
		}
	}

	return result
}

// TransformationCandidates takes an interface name ("xe-0/0/1:1") and a speed,
// and returns a map[int]Transformation populated with candidate transformations
// available according to the PortInfo and keyed by the transformation ID. Only
// "active" transformations are returned.
//
// Deprecated: Use the PortWithMatchingTransforms() device profile method instead
func (p Port) TransformationCandidates(ifName string, ifSpeed speed.Speed) map[int]Transformation {
	transforms := p.transformationCandidates(ifName, ifSpeed)
	if len(transforms) == 0 {
		return nil
	}

	result := make(map[int]Transformation, len(transforms))
	for _, transformation := range transforms {
		result[transformation.ID] = transformation
	}

	return result
}

// Transformation returns the Transformation with the specified ID. If no
// such Transformation exists, an error is returned.
func (p *Port) Transformation(id int) (Transformation, error) {
	for _, t := range p.Transformations {
		if t.ID == id {
			return t, nil
		}
	}

	return Transformation{}, sdk.ErrNotFound(fmt.Sprintf("transformation id %d not found", id))
}
