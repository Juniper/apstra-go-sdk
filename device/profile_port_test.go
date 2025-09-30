// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPort_TransformationCandidates(t *testing.T) {
	type testCase struct {
		port     Port
		speed    speed.Speed
		name     string
		expected map[int]Transformation
	}

	testCases := map[string]testCase{
		"none": {
			port:     Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "eth0", State: enum.InterfaceStateActive, Setting: pointer.To(""), Speed: "10G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(1)},
			name:     "eth0",
			speed:    "100G",
			expected: nil,
		},
		"eth0": {
			port:  Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "eth0", State: enum.InterfaceStateActive, Setting: pointer.To(""), Speed: "10G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(1)},
			name:  "eth0",
			speed: "10G",
			expected: map[int]Transformation{
				1: {ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "eth0", State: enum.InterfaceStateActive, Setting: pointer.To(""), Speed: "10G"}}},
			},
		},
		"multiple": {
			port:  Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":""},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":"10000"},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":"1000"},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "1G"}}}}, Column: 6, ID: 12, Row: 2, FailureDomain: 1, Slot: 0},
			name:  "Ethernet1/12",
			speed: "10G",
			expected: map[int]Transformation{
				1: {ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":""},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}},
				2: {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":"10000"},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}},
			},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			r := tCase.port.TransformationCandidates(tCase.name, tCase.speed)
			require.Equal(t, tCase.expected, r)
		})
	}
}
