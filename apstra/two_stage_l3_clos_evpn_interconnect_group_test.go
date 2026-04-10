// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_test

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedatacenter "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/datacenter"
	"github.com/stretchr/testify/require"
)

func TestEVPNInterconnectGroup_MarshalJSON(t *testing.T) {
	type testCase struct {
		data apstra.EVPNInterconnectGroup
		exp  string
	}

	testCases := map[string]testCase{
		"empty": {
			data: apstra.EVPNInterconnectGroup{},
			exp:  "{}",
		},
		"full": {
			data: apstra.EVPNInterconnectGroup{
				Label:       pointer.To("label"),
				RouteTarget: pointer.To("1:1"),
				ESIMAC:      testutils.Must(net.ParseMAC("00:11:22:33:44:55")),
				InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
					"a": {
						L3Enabled:       true,
						RouteTarget:     pointer.To("11:11"),
						RoutingPolicyId: pointer.To("aa"),
					},
					"b": {
						L3Enabled:       false,
						RouteTarget:     pointer.To("22:22"),
						RoutingPolicyId: pointer.To("bb"),
					},
				},
				InterconnectVirtualNetworks: map[string]apstra.InterconnectVirtualNetwork{
					"c": {
						L2Enabled:      true,
						L3Enabled:      true,
						TranslationVNI: pointer.To(uint32(333)),
					},
					"d": {
						L2Enabled:      false,
						L3Enabled:      false,
						TranslationVNI: pointer.To(uint32(444)),
					},
				},
			},
			exp: `{
  "label": "label",
  "interconnect_route_target": "1:1",
  "interconnect_esi_mac": "00:11:22:33:44:55",
  "interconnect_security_zones": {
    "a": {
      "enabled_for_l3": true,
      "interconnect_route_target": "11:11",
      "routing_policy_id": "aa"
    },
    "b": {
      "enabled_for_l3": false,
      "interconnect_route_target": "22:22",
      "routing_policy_id": "bb"
    }
  },
  "interconnect_virtual_networks": {
    "c": {
      "l2": true,
      "l3": true,
      "translation_vni": 333
    },
    "d": {
      "l2": false,
      "l3": false,
      "translation_vni": 444
    }
  }
}`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			result, err := json.Marshal(tCase.data)
			require.NoError(t, err)
			require.JSONEq(t, tCase.exp, string(result))
		})
	}
}

func TestEVPNInterconnectGroup_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		data string
		exp  apstra.EVPNInterconnectGroup
	}

	testCases := map[string]testCase{
		"empty": {
			exp:  apstra.EVPNInterconnectGroup{},
			data: "{}",
		},
		"full": {
			exp: apstra.EVPNInterconnectGroup{
				Label:       pointer.To("label"),
				RouteTarget: pointer.To("1:1"),
				ESIMAC:      testutils.Must(net.ParseMAC("00:11:22:33:44:55")),
				InterconnectSecurityZones: map[string]apstra.InterconnectSecurityZone{
					"a": {
						L3Enabled:       true,
						RouteTarget:     pointer.To("11:11"),
						RoutingPolicyId: pointer.To("aa"),
					},
					"b": {
						L3Enabled:       false,
						RouteTarget:     pointer.To("22:22"),
						RoutingPolicyId: pointer.To("bb"),
					},
				},
				InterconnectVirtualNetworks: map[string]apstra.InterconnectVirtualNetwork{
					"c": {
						L2Enabled:      true,
						L3Enabled:      true,
						TranslationVNI: pointer.To(uint32(333)),
					},
					"d": {
						L2Enabled:      false,
						L3Enabled:      false,
						TranslationVNI: pointer.To(uint32(444)),
					},
				},
			},
			data: `{
  "label": "label",
  "interconnect_route_target": "1:1",
  "interconnect_esi_mac": "00:11:22:33:44:55",
  "interconnect_security_zones": {
    "a": {
      "enabled_for_l3": true,
      "interconnect_route_target": "11:11",
      "routing_policy_id": "aa"
    },
    "b": {
      "enabled_for_l3": false,
      "interconnect_route_target": "22:22",
      "routing_policy_id": "bb"
    }
  },
  "interconnect_virtual_networks": {
    "c": {
      "l2": true,
      "l3": true,
      "translation_vni": 333
    },
    "d": {
      "l2": false,
      "l3": false,
      "translation_vni": 444
    }
  }
}`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			var result apstra.EVPNInterconnectGroup
			require.NoError(t, json.Unmarshal([]byte(tCase.data), &result))
			comparedatacenter.EVPNInterconnectGroup(t, tCase.exp, result)
		})
	}

}
