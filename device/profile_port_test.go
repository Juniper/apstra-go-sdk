// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"reflect"
	"testing"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
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

func TestPort_transformationCandidates(t *testing.T) {
	type testCase struct {
		port     Port
		speed    speed.Speed
		name     string
		expected []Transformation
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
			expected: []Transformation{
				{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "eth0", State: enum.InterfaceStateActive, Setting: pointer.To(""), Speed: "10G"}}},
			},
		},
		"multiple": {
			port:  Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":""},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":"10000"},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":"1000"},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "1G"}}}}, Column: 6, ID: 12, Row: 2, FailureDomain: 1, Slot: 0},
			name:  "Ethernet1/12",
			speed: "10G",
			expected: []Transformation{
				{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":""},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}},
				{ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "Ethernet1/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"interface":{"speed":"10000"},"global":{"speed":"","module_index":-1,"port_index":-1}}`), Speed: "10G"}}},
			},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			r := tCase.port.transformationCandidates(tCase.name, tCase.speed)
			require.Equal(t, tCase.expected, r)
		})
	}
}

func TestPort_DefaultTransform(t *testing.T) {
	type testCase struct {
		port   Port
		expect Transformation
		expErr any
	}

	testCases := map[string]testCase{
		"Juniper_EX4400-48F_p1_t1": {
			port:   testProfileJuniperEX440048F.Ports[0],
			expect: testProfileJuniperEX440048F.Ports[0].Transformations[0],
		},
		"Juniper_EX4400-48F_p49_t1": {
			port:   testProfileJuniperEX440048F.Ports[49],
			expect: testProfileJuniperEX440048F.Ports[49].Transformations[0],
		},
		"err_not_found": {
			port:   Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(0), Slot: 0},
			expErr: new(sdk.ErrNotFound),
		},
		"not_at_index_zero": {
			port:   Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(0), Slot: 0},
			expect: Transformation{ID: 2, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result, err := tCase.port.DefaultTransform()
			if tCase.expErr != nil {
				target := reflect.New(reflect.TypeOf(tCase.expErr).Elem()).Interface()
				require.ErrorAs(t, err, target)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.expect, result)
		})
	}
}

func TestPort_Transformation(t *testing.T) {
	type testCase struct {
		port   Port
		id     int
		expect Transformation
		expErr any
	}
	testCases := map[string]testCase{
		"Juniper_EX4400-48F_p1_t1": {
			port:   testProfileJuniperEX440048F.Ports[0],
			id:     1,
			expect: testProfileJuniperEX440048F.Ports[0].Transformations[0],
		},
		"Juniper_EX4400-48F_p49_t2": {
			port:   testProfileJuniperEX440048F.Ports[49],
			id:     2,
			expect: testProfileJuniperEX440048F.Ports[49].Transformations[1],
		},
		"not_found": {
			port:   testProfileJuniperEX440048F.Ports[49],
			id:     0,
			expErr: new(sdk.ErrNotFound),
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result, err := tCase.port.Transformation(tCase.id)
			if tCase.expErr != nil {
				target := reflect.New(reflect.TypeOf(tCase.expErr).Elem()).Interface()
				require.ErrorAs(t, err, target)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.expect, result)
		})
	}
}
