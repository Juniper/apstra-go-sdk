//go:build integration
// +build integration

package apstra //

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestIbaProbes(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for clientName, client := range clients {
		log.Printf("testing Predefined Probes against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())

		bpClient, _ := testBlueprintA(ctx, t, client.client)
		// defer bpDelete(ctx)
		// pdps, err := bpClient.GetAllIbaPredefinedProbes(ctx)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// expectedToFail := map[string]bool{
		// 	"external_ecmp_imbalance":            true,
		// 	"evpn_vxlan_type5":                   true,
		// 	"eastwest_traffic":                   true,
		// 	"vxlan_floodlist":                    true,
		// 	"fabric_hotcold_ifcounter":           true,
		// 	"specific_interface_flapping":        true,
		// 	"evpn_vxlan_type3":                   true,
		// 	"specific_hotcold_ifcounter":         true,
		// 	"spine_superspine_hotcold_ifcounter": true,
		// }

		// for _, p := range pdps {
		// 	t.Logf("Get Predefined Probe By Name %s", p.Name)
		// 	_, err := bpClient.GetIbaPredefinedProbeByName(ctx, p.Name)
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}
		// 	t.Log(p.Description)
		// 	t.Log(p.Schema)
		//
		// 	t.Logf("Instantiating Probe %s", p.Name)
		//
		// 	probeId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &IbaPredefinedProbeRequest{
		// 		Name: p.Name,
		// 		Data: json.RawMessage([]byte(`{"label":"` + p.Name + `"}`)),
		// 	})
		// 	if err != nil {
		// 		if !expectedToFail[p.Name] {
		// 			t.Fatal(err)
		// 		} else {
		// 			t.Logf("%s was expected to fail", p.Name)
		// 			continue
		// 		}
		// 	}
		//
		// 	t.Logf("Got back Probe Id %s \n Now GET it.", probeId)
		//
		// 	p, err := bpClient.GetIbaProbe(ctx, probeId)
		//
		// 	t.Logf("Label %s", p.Label)
		// 	t.Logf("Description %s", p.Description)
		// 	t.Log(p)
		// 	t.Logf("Delete probe")
		// 	for _, i := range p.Stages {
		// 		t.Logf("Stage name %s", i["name"])
		// 	}
		// 	err = bpClient.DeleteIbaProbe(ctx, probeId)
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}
		// 	t.Logf("Delete Probe again, this should fail")
		// 	err = bpClient.DeleteIbaProbe(ctx, probeId)
		// 	if err == nil {
		// 		t.Fatal("Probe Deletion should have failed")
		// 	} else {
		// 		t.Log(err)
		// 	}
		// }
		probeStr := `{
  "label": "Bandwidth Utilization",
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
		_, err = bpClient.CreateIbaProbe(ctx, json.RawMessage(probeStr))
		log.Println(err)
	}
}
