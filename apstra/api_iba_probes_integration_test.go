// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestIbaProbes(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	probeStr := `{
  "label": "Test Probe",
  "description": "The probe calculates interfaces bandwidth",
  "processors": [
    {
      "name": "Egress traffic",
      "type": "if_counter",
      "properties": {
        "description": "str(interface.description or '')",
        "counter_type": "tx_bps",
        "graph_query": "node('system', name='system', system_type='switch', deploy_mode='deploy').out('hosted_interfaces').node('interface', if_type=is_in(['ip','ethernet', 'port_channel']), name='interface').out('link').node('link', name='link')",
        "query_group_by": [],
        "query_tag_filter": {
          "filter": {},
          "operation": "and"
        },
        "interface": "interface.if_name",
        "system_id": "system.system_id",
        "group": "str('mlag_peer' if link.role in ['leaf_peer_link', 'leaf_l3_peer_link'] else ('leaf_access' if 'access' in link.role else link.role))",
        "query_expansion": {},
        "enable_streaming": false
      },
      "inputs": {},
      "outputs": {
        "out": "egress_traffic"
      }
    },
    {
      "name": "Ingress traffic",
      "type": "if_counter",
      "properties": {
        "description": "str(interface.description or '')",
        "counter_type": "rx_bps",
        "graph_query": "node('system', name='system', system_type='switch', deploy_mode='deploy').out('hosted_interfaces').node('interface', if_type=is_in(['ip','ethernet', 'port_channel']), name='interface').out('link').node('link', name='link')",
        "query_group_by": [],
        "query_tag_filter": {
          "filter": {},
          "operation": "and"
        },
        "interface": "interface.if_name",
        "system_id": "system.system_id",
        "group": "str('mlag_peer' if link.role in ['leaf_peer_link', 'leaf_l3_peer_link'] else ('leaf_access' if 'access' in link.role else link.role))",
        "query_expansion": {},
        "enable_streaming": false
      },
      "inputs": {},
      "outputs": {
        "out": "ingress_traffic"
      }
    },
    {
      "name": "Egress traffic first summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 120,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "egress_traffic",
          "column": "value"
        }
      },
      "outputs": {
        "out": "egress_traffic_first_summary"
      }
    },
    {
      "name": "Ingress traffic first summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 120,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "ingress_traffic",
          "column": "value"
        }
      },
      "outputs": {
        "out": "ingress_traffic_first_summary"
      }
    },
    {
      "name": "Bucketed egress traffic",
      "type": "sum",
      "properties": {
        "group_by": [
          "group"
        ],
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "egress_traffic_first_summary",
          "column": "value"
        }
      },
      "outputs": {
        "out": "egress_by_group_traffic"
      }
    },
    {
      "name": "Bucketed ingress traffic",
      "type": "sum",
      "properties": {
        "group_by": [
          "group"
        ],
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "ingress_traffic_first_summary",
          "column": "value"
        }
      },
      "outputs": {
        "out": "ingress_by_group_traffic"
      }
    },
    {
      "name": "Egress traffic second summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 3600,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "egress_traffic_first_summary",
          "column": "value"
        }
      },
      "outputs": {
        "out": "egress_traffic_second_summary"
      }
    },
    {
      "name": "Ingress traffic second summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 3600,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "ingress_traffic_first_summary",
          "column": "value"
        }
      },
      "outputs": {
        "out": "ingress_traffic_second_summary"
      }
    },
    {
      "name": "Egress by group traffic first summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 120,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "egress_by_group_traffic",
          "column": "value"
        }
      },
      "outputs": {
        "out": "egress_by_group_traffic_first_summary"
      }
    },
    {
      "name": "Ingress by group traffic first summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 120,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "ingress_by_group_traffic",
          "column": "value"
        }
      },
      "outputs": {
        "out": "ingress_by_group_traffic_first_summary"
      }
    },
    {
      "name": "Egress by group traffic second summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 3600,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "egress_by_group_traffic_first_summary",
          "column": "value"
        }
      },
      "outputs": {
        "out": "egress_by_group_traffic_second_summary"
      }
    },
    {
      "name": "Ingress by group traffic second summary",
      "type": "periodic_average",
      "properties": {
        "graph_query": [],
        "period": 3600,
        "enable_streaming": false
      },
      "inputs": {
        "in": {
          "stage": "ingress_by_group_traffic_first_summary",
          "column": "value"
        }
      },
      "outputs": {
        "out": "ingress_by_group_traffic_second_summary"
      }
    }
  ],
  "stages": [
    {
      "name": "ingress_by_group_traffic",
      "retention_duration": 86400,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "ingress_traffic_second_summary",
      "enable_metric_logging": true,
      "retention_duration": 2592000,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "egress_traffic_first_summary",
      "enable_metric_logging": true,
      "retention_duration": 3600,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "ingress_by_group_traffic_second_summary",
      "enable_metric_logging": true,
      "retention_duration": 2592000,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "egress_by_group_traffic_second_summary",
      "enable_metric_logging": true,
      "retention_duration": 2592000,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "egress_traffic",
      "retention_duration": 86400,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "ingress_traffic_first_summary",
      "enable_metric_logging": true,
      "retention_duration": 3600,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "egress_traffic_second_summary",
      "enable_metric_logging": true,
      "retention_duration": 2592000,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "egress_by_group_traffic_first_summary",
      "enable_metric_logging": true,
      "retention_duration": 3600,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "egress_by_group_traffic",
      "retention_duration": 86400,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "ingress_traffic",
      "retention_duration": 86400,
      "units": {
        "value": "bps"
      }
    },
    {
      "name": "ingress_by_group_traffic_first_summary",
      "enable_metric_logging": true,
      "retention_duration": 3600,
      "units": {
        "value": "bps"
      }
    }
  ]
}`

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, ctx)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)
			predefinedProbes, err := bpClient.GetAllIbaPredefinedProbes(ctx)
			require.NoError(t, err)

			expectedToFail := map[string]bool{
				"external_ecmp_imbalance":            true,
				"evpn_vxlan_type5":                   true,
				"eastwest_traffic":                   true,
				"vxlan_floodlist":                    true,
				"fabric_hotcold_ifcounter":           true,
				"evpn_vxlan_type3":                   true,
				"specific_hotcold_ifcounter":         true,
				"spine_superspine_hotcold_ifcounter": true,
			}

			if version.MustConstraints(version.NewConstraint("<5.1.0")).Check(client.APIVersion()) {
				expectedToFail["specific_interface_flapping"] = true
			}

			for _, predefinedProbe := range predefinedProbes {
				t.Run(predefinedProbe.Name, func(t *testing.T) {
					ctx := testutils.WrapCtxWithTestId(t, ctx)

					t.Logf("Get Predefined Probe By Name %s", predefinedProbe.Name)
					_, err := bpClient.GetIbaPredefinedProbeByName(ctx, predefinedProbe.Name)
					require.NoError(t, err)

					t.Logf("Instantiating Probe %s", predefinedProbe.Name)

					probeId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &apstra.IbaPredefinedProbeRequest{
						Name: predefinedProbe.Name,
						Data: json.RawMessage(`{"label":"` + predefinedProbe.Name + `"}`),
					})
					if expectedToFail[predefinedProbe.Name] {
						t.Log(err)
						t.Logf("%s was expected to fail", predefinedProbe.Name)
						return
					}
					require.NoError(t, err)

					_, err = bpClient.GetIbaProbe(ctx, probeId)
					require.NoError(t, err)

					t.Log("Get IBA probe state")
					_, err = bpClient.GetIbaProbeState(ctx, probeId)
					require.NoError(t, err)

					t.Logf("Delete probe")
					require.NoError(t, bpClient.DeleteIbaProbe(ctx, probeId))

					var ace apstra.ClientErr
					t.Logf("Delete Probe again, this should fail")
					err = bpClient.DeleteIbaProbe(ctx, probeId)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
			t.Log("Create Probe With Json")
			id, err := bpClient.CreateIbaProbeFromJson(ctx, json.RawMessage(probeStr))
			require.NoError(t, err)

			t.Log("Test Get Probe")
			p, err := bpClient.GetIbaProbe(ctx, id)
			require.NoError(t, err)
			require.Equal(t, "Test Probe", p.Label, "expected label %q got %q", "Test Probe", p.Label)

			_, err = bpClient.GetIbaProbeState(ctx, id)
			require.NoError(t, err)

			require.NoError(t, bpClient.DeleteIbaProbe(ctx, id))
		})
	}
}
