// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"encoding/json"
	"reflect"
	"testing"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/stretchr/testify/require"
)

func TestProfile_MarshalJSON(t *testing.T) {
	type testCase struct {
		v Profile
		e string
	}

	testCases := map[string]testCase{
		"generic_1x10": {
			v: testProfileGeneric1x10,
			e: testProfileGeneric1x10JSON,
		},
		"Juniper_vEX": {
			v: testProfileJunipervEX,
			e: testProfileJunipervEXJSON,
		},
		"Juniper_EX4400-48F": {
			v: testProfileJuniperEX440048F,
			e: testProfileJuniperEX440048FJSON,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			r, err := json.Marshal(tCase.v)
			require.NoError(t, err)

			// get rid of extraneous fields in the expected string value
			eMap := map[string]json.RawMessage{}
			require.NoError(t, json.Unmarshal([]byte(tCase.e), &eMap))
			delete(eMap, "id")
			delete(eMap, "created_at")
			delete(eMap, "last_modified_at")
			delete(eMap, "predefined")
			e, err := json.Marshal(eMap)
			require.NoError(t, err)

			require.JSONEq(t, string(e), string(r))

			/*
			   inspect raw json with this
			   pbpaste | jq 'walk(if type == "object" then . | to_entries | sort_by(.key) | from_entries else . end)'
			*/
		})
	}
}

func TestProfile_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		e Profile
		v string
	}

	testCases := map[string]testCase{
		"generic_1x10": {
			v: testProfileGeneric1x10JSON,
			e: testProfileGeneric1x10,
		},
		"Juniper_vEX": {
			v: testProfileJunipervEXJSON,
			e: testProfileJunipervEX,
		},
		"Juniper_EX4400-48F": {
			v: testProfileJuniperEX440048FJSON,
			e: testProfileJuniperEX440048F,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			var r Profile
			err := json.Unmarshal([]byte(tCase.v), &r)
			require.NoError(t, err)

			require.Equal(t, tCase.e, r)
		})
	}
}

func TestProfile_PortByID(t *testing.T) {
	type testCase struct {
		profile Profile
		id      int
		exp     Port
		expErr  any
	}

	testCases := map[string]testCase{
		"generic_1x10": {
			profile: testProfileGeneric1x10,
			id:      1,
			exp:     Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "eth0", State: enum.InterfaceStateActive, Setting: pointer.To(""), Speed: "10G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(1)},
			expErr:  nil,
		},
		"Juniper_vEX": {
			profile: testProfileJunipervEX,
			id:      6,
			exp:     Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 6, ID: 6, Row: 1, FailureDomain: 1, Display: pointer.To(5), Slot: 0},
			expErr:  nil,
		},
		"Juniper_EX4400-48F": {
			profile: testProfileJuniperEX440048F,
			id:      50,
			exp:     Port{ConnectorType: "qsfp28", Panel: 3, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "et-0/1/1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "100G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "et-0/1/1:0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 2, Name: "et-0/1/1:1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 3, Name: "et-0/1/1:2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 4, Name: "et-0/1/1:3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/1/1:0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 2, Name: "xe-0/1/1:1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 3, Name: "xe-0/1/1:2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 4, Name: "xe-0/1/1:3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}}}}, Column: 1, ID: 50, Row: 2, FailureDomain: 1, Display: pointer.To(49), Slot: 0},
			expErr:  nil,
		},
		"generic_1x10_not_found": {
			profile: testProfileGeneric1x10,
			id:      2,
			exp:     Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "eth0", State: enum.InterfaceStateActive, Setting: pointer.To(""), Speed: "10G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(1)},
			expErr:  new(sdk.ErrNotFound),
		},
		"bogus_multiple_match": {
			profile: Profile{
				Ports: []Port{
					{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(0), Slot: 0},
					{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/1", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/1", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 2, ID: 2, Row: 1, FailureDomain: 1, Display: pointer.To(1), Slot: 0},
					{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/2", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/2", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 3, ID: 3, Row: 1, FailureDomain: 1, Display: pointer.To(2), Slot: 0},
					{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/3", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/3", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 4, ID: 2, Row: 1, FailureDomain: 1, Display: pointer.To(3), Slot: 0}, // <- duplicate ID: 2 in here
					{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/4", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/4", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 5, ID: 5, Row: 1, FailureDomain: 1, Display: pointer.To(4), Slot: 0},
				},
			},
			id:     2,
			exp:    Port{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 6, ID: 6, Row: 1, FailureDomain: 1, Display: pointer.To(5), Slot: 0},
			expErr: new(sdk.ErrMultipleMatch),
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result, err := tCase.profile.PortByID(tCase.id)
			if tCase.expErr != nil {
				// create an error of the expected type
				target := reflect.New(reflect.TypeOf(tCase.expErr).Elem()).Interface()
				require.ErrorAs(t, err, target)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tCase.exp, result)
		})
	}
}
