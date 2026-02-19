// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package apstra_test

import (
	"encoding/json"
	"net/netip"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	comparefreeform "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/freeform"
	"github.com/stretchr/testify/require"
)

func TestFreeformAggregateLink_MarshalJSON(t *testing.T) {
	type testcase struct {
		data     apstra.FreeformAggregateLink
		expected string
	}
	testcases := map[string]testcase{
		"one_server_two_switches": {
			data: apstra.FreeformAggregateLink{
				Label:         "label",
				MemberLinkIds: []string{"link_id_1", "link_id_2", "link_id_3"},
				EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
					{
						Label: "server",
						Tags:  []string{"server_tag_1", "server_tag_2", "server_tag_3"},
						Endpoints: []apstra.FreeformAggregateLinkEndpoint{
							{
								SystemID:      "server_id",
								IfName:        "bond0",
								IPv4Addr:      pointer.To(netip.MustParsePrefix("192.0.2.50/24")),
								IPv6Addr:      pointer.To(netip.MustParsePrefix("3fff::50/64")),
								PortChannelID: 1,
								Tags:          []string{"bond0_tag"},
								LAGMode:       enum.LAGModeActiveLACP,
							},
						},
					},
					{
						Label: "switch",
						Tags:  []string{"switch_1", "switch_2"},
						Endpoints: []apstra.FreeformAggregateLinkEndpoint{
							{
								SystemID:      "switch_1_id",
								IfName:        "ae1",
								IPv4Addr:      pointer.To(netip.MustParsePrefix("192.0.2.11/24")),
								IPv6Addr:      pointer.To(netip.MustParsePrefix("3fff::11/64")),
								PortChannelID: 1,
								Tags:          []string{"switch_1_lag_tag"},
								LAGMode:       enum.LAGModeActiveLACP,
							},
							{
								SystemID:      "switch_2_id",
								IfName:        "ae1",
								IPv4Addr:      pointer.To(netip.MustParsePrefix("192.0.2.12/24")),
								IPv6Addr:      pointer.To(netip.MustParsePrefix("3fff::12/64")),
								PortChannelID: 1,
								Tags:          []string{"switch_2_lag_tag"},
								LAGMode:       enum.LAGModeActiveLACP,
							},
						},
					},
				},
				Tags: []string{"lag tag1", "lag tag2", "lag tag3"},
			},
			expected: `{
              "label": "label",
              "member_link_ids": [ "link_id_1", "link_id_2", "link_id_3" ],
              "endpoints": [
                {
                  "system": { "id": "server_id" },
                  "interface": {
                    "if_name": "bond0",
                    "port_channel_id": 1,
                    "lag_mode": "lacp_active",
                    "ipv4_addr": "192.0.2.50/24",
                    "ipv6_addr": "3fff::50/64",
                    "tags": [ "bond0_tag" ]
                  },
                  "endpoint_group": 0
                },
                {
                  "system": { "id": "switch_1_id" },
                  "interface": {
                    "if_name": "ae1",
                    "port_channel_id": 1,
                    "lag_mode": "lacp_active",
                    "ipv4_addr": "192.0.2.11/24",
                    "ipv6_addr": "3fff::11/64",
                    "tags": [ "switch_1_lag_tag" ]
                  },
                  "endpoint_group": 1
                },
                {
                  "system": { "id": "switch_2_id" },
                  "interface": {
                    "if_name": "ae1",
                    "port_channel_id": 1,
                    "lag_mode": "lacp_active",
                    "ipv4_addr": "192.0.2.12/24",
                    "ipv6_addr": "3fff::12/64",
                    "tags": [ "switch_2_lag_tag" ]
                  },
                  "endpoint_group": 1
                }
              ],
              "endpoint_groups": {
                "0": {
                  "label": "server",
                  "tags": [ "server_tag_1", "server_tag_2", "server_tag_3" ]
                },
                "1": {
                  "label": "switch",
                  "tags": [ "switch_1", "switch_2" ]
                }
              },
              "tags": [ "lag tag1", "lag tag2", "lag tag3" ]
            }`,
		},
	}

	for tName, tCase := range testcases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result, err := json.Marshal(tCase.data)
			require.NoError(t, err)
			require.JSONEq(t, tCase.expected, string(result))
		})
	}
}

func TestFreeformAggregateLink_UnmarshalJSON(t *testing.T) {
	type testcase struct {
		data     string
		expected apstra.FreeformAggregateLink
	}

	oneServerTwoSwitches := apstra.FreeformAggregateLink{
		Label:         "label",
		MemberLinkIds: []string{"link_id_1", "link_id_2", "link_id_3"},
		EndpointGroups: [2]apstra.FreeformAggregateLinkEndpointGroup{
			{
				Label: "server",
				Tags:  []string{"server_tag_1", "server_tag_2", "server_tag_3"},
				Endpoints: []apstra.FreeformAggregateLinkEndpoint{
					{
						SystemID:      "server_id",
						IfName:        "bond0",
						IPv4Addr:      pointer.To(netip.MustParsePrefix("192.0.2.50/24")),
						IPv6Addr:      pointer.To(netip.MustParsePrefix("3fff::50/64")),
						PortChannelID: 1,
						Tags:          []string{"bond0_tag"},
						LAGMode:       enum.LAGModeActiveLACP,
					},
				},
			},
			{
				Label: "switch",
				Tags:  []string{"switch_1", "switch_2"},
				Endpoints: []apstra.FreeformAggregateLinkEndpoint{
					{
						SystemID:      "switch_1_id",
						IfName:        "ae1",
						IPv4Addr:      pointer.To(netip.MustParsePrefix("192.0.2.11/24")),
						IPv6Addr:      pointer.To(netip.MustParsePrefix("3fff::11/64")),
						PortChannelID: 1,
						Tags:          []string{"switch_1_lag_tag"},
						LAGMode:       enum.LAGModeActiveLACP,
					},
					{
						SystemID:      "switch_2_id",
						IfName:        "ae1",
						IPv4Addr:      pointer.To(netip.MustParsePrefix("192.0.2.12/24")),
						IPv6Addr:      pointer.To(netip.MustParsePrefix("3fff::12/64")),
						PortChannelID: 1,
						Tags:          []string{"switch_2_lag_tag"},
						LAGMode:       enum.LAGModeActiveLACP,
					},
				},
			},
		},
		Tags: []string{"lag tag1", "lag tag2", "lag tag3"},
	}
	oneServerTwoSwitches.SetID("link_id")
	oneServerTwoSwitches.EndpointGroups[0].SetID("endpoint_group_0_id")
	oneServerTwoSwitches.EndpointGroups[0].Endpoints[0].SetID("endpoint_group_0_endpoint_0_id")
	oneServerTwoSwitches.EndpointGroups[1].SetID("endpoint_group_1_id")
	oneServerTwoSwitches.EndpointGroups[1].Endpoints[0].SetID("endpoint_group_1_endpoint_0_id")
	oneServerTwoSwitches.EndpointGroups[1].Endpoints[1].SetID("endpoint_group_1_endpoint_1_id")

	testcases := map[string]testcase{
		"one_server_two_switches": {
			expected: oneServerTwoSwitches,
			data: `{
              "id": "link_id",
              "label": "label",
              "member_link_ids": [ "link_id_1", "link_id_2", "link_id_3" ],
              "endpoints": [
                {
                  "id": "endpoint_group_0_endpoint_0_id",
                  "system": { "id": "server_id" },
                  "interface": {
                    "if_name": "bond0",
                    "port_channel_id": 1,
                    "lag_mode": "lacp_active",
                    "ipv4_addr": "192.0.2.50/24",
                    "ipv6_addr": "3fff::50/64",
                    "tags": [ "bond0_tag" ]
                  },
                  "endpoint_group": 0
                },
                {
                  "id": "endpoint_group_1_endpoint_0_id",
                  "system": { "id": "switch_1_id" },
                  "interface": {
                    "if_name": "ae1",
                    "port_channel_id": 1,
                    "lag_mode": "lacp_active",
                    "ipv4_addr": "192.0.2.11/24",
                    "ipv6_addr": "3fff::11/64",
                    "tags": [ "switch_1_lag_tag" ]
                  },
                  "endpoint_group": 1
                },
                {
                  "id": "endpoint_group_1_endpoint_1_id",
                  "system": { "id": "switch_2_id" },
                  "interface": {
                    "if_name": "ae1",
                    "port_channel_id": 1,
                    "lag_mode": "lacp_active",
                    "ipv4_addr": "192.0.2.12/24",
                    "ipv6_addr": "3fff::12/64",
                    "tags": [ "switch_2_lag_tag" ]
                  },
                  "endpoint_group": 1
                }
              ],
              "endpoint_groups": {
                "0": {
                  "id":"endpoint_group_0_id",
                  "label": "server",
                  "tags": [ "server_tag_1", "server_tag_2", "server_tag_3" ]
                },
                "1": {
                  "id":"endpoint_group_1_id",
                  "label": "switch",
                  "tags": [ "switch_1", "switch_2" ]
                }
              },
              "tags": [ "lag tag1", "lag tag2", "lag tag3" ]
            }`,
		},
	}

	for tName, tCase := range testcases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var result apstra.FreeformAggregateLink
			err := json.Unmarshal([]byte(tCase.data), &result)
			require.NoError(t, err)
			comparefreeform.AggregateLink(t, tCase.expected, result)
		})
	}
}
