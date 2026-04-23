// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/errors"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedatacenter "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/datacenter"
	"github.com/stretchr/testify/require"
)

func TestSwitchingZone_ID_SetID(t *testing.T) {
	var o datacenter.SwitchingZone
	require.Nil(t, o.ID())
	id := testutils.RandString(6, "hex")
	require.NoError(t, o.SetID(id))
	err := o.SetID(id)
	require.Error(t, err)
	var eErr errors.IDAlreadySet
	require.ErrorAs(t, err, &eErr)
	require.Equal(t, id, *o.ID())
}

func TestSwitchingZone_MarshalJSON(t *testing.T) {
	type testCase struct {
		d datacenter.SwitchingZone
		e string
	}

	testCases := map[string]testCase{
		"vlan_aware": {
			d: datacenter.SwitchingZone{
				Label:             pointer.To("label1"),
				MACVRFDescription: pointer.To("description1"),
				MACVRFName:        pointer.To("name1"),
				MACVRFServiceType: pointer.To(enum.SwitchingZoneMACVRFServiceTypeVLANAware),
				RouteTarget:       pointer.To("1:1"),
				Tags:              []string{"tag1", "tag2"},
			},
			e: `{
                  "label": "label1",
                  "mac_vrf_description": "description1",
                  "mac_vrf_name": "name1",
                  "mac_vrf_service_type": "vlan_aware",
                  "route_target": "1:1",
                  "tags": ["tag1", "tag2"],
                  "impl_type": "mac_vrf"
                }`,
		},
		"vlan_bundle": {
			d: datacenter.SwitchingZone{
				Label:             pointer.To("label2"),
				MACVRFDescription: pointer.To("description2"),
				MACVRFName:        pointer.To("name2"),
				MACVRFServiceType: pointer.To(enum.SwitchingZoneMACVRFServiceTypeVLANBundle),
				RouteTarget:       pointer.To("2:2"),
				Tags:              []string{"tag3", "tag4"},
			},
			e: `{
                  "label": "label2",
                  "mac_vrf_description": "description2",
                  "mac_vrf_name": "name2",
                  "mac_vrf_service_type": "vlan_bundle",
                  "route_target": "2:2",
                  "tags": ["tag3", "tag4"],
                  "impl_type": "mac_vrf"
                }`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			r, err := json.Marshal(tCase.d)
			require.NoError(t, err)
			log.Println(string(r))
			require.JSONEq(t, tCase.e, string(r))
		})
	}
}

func TestSwitchingZone_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		d   string
		e   datacenter.SwitchingZone
		eid string
	}

	testCases := map[string]testCase{
		"vlan_aware": {
			d: `{
                  "label": "label1",
                  "mac_vrf_description": "description1",
                  "mac_vrf_name": "name1",
                  "mac_vrf_service_type": "vlan_aware",
                  "route_target": "1:1",
                  "tags": ["tag1", "tag2"],
                  "impl_type": "mac_vrf",
                  "id": "abc"
                 }`,
			e: datacenter.SwitchingZone{
				Label:             pointer.To("label1"),
				MACVRFDescription: pointer.To("description1"),
				MACVRFName:        pointer.To("name1"),
				MACVRFServiceType: pointer.To(enum.SwitchingZoneMACVRFServiceTypeVLANAware),
				RouteTarget:       pointer.To("1:1"),
				Tags:              []string{"tag1", "tag2"},
			},
			eid: "abc",
		},
		"vlan_bundle": {
			d: `{
                  "label": "label2",
                  "mac_vrf_description": "description2",
                  "mac_vrf_name": "name2",
                  "mac_vrf_service_type": "vlan_bundle",
                  "route_target": "2:2",
                  "tags": ["tag3", "tag4"],
                  "impl_type": "mac_vrf",
                  "id": "def"
                }`,
			e: datacenter.SwitchingZone{
				Label:             pointer.To("label2"),
				MACVRFDescription: pointer.To("description2"),
				MACVRFName:        pointer.To("name2"),
				MACVRFServiceType: pointer.To(enum.SwitchingZoneMACVRFServiceTypeVLANBundle),
				RouteTarget:       pointer.To("2:2"),
				Tags:              []string{"tag3", "tag4"},
			},
			eid: "def",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			var r datacenter.SwitchingZone
			require.NoError(t, json.Unmarshal([]byte(tCase.d), &r))
			comparedatacenter.SwitchingZone(t, tCase.e, r)
			require.NotNil(t, r.ID())
			require.Equal(t, tCase.eid, *r.ID())
		})
	}
}
