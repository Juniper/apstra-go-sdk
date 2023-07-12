package apstra

import (
	"encoding/json"
	"testing"
)

func TestParseCT(t *testing.T) {
	apiResponse := `{
  "policies": [
    {
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table. Can be assigned to physical interfaces or single-chassis LAGs (not applicable to ESI LAG or MLAG interfaces).",
      "tags": [],
      "label": "IP Link (batch)",
      "visible": false,
      "policy_type_name": "batch",
      "attributes": {
        "subpolicies": [
          "2e264a0a-ca5b-413c-9bab-c8c1c4aa4e6e"
        ]
      },
      "id": "1bbeaf90-ad40-4758-a3e6-15091635dd39"
    },
    {
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table. Can be assigned to physical interfaces or single-chassis LAGs (not applicable to ESI LAG or MLAG interfaces).",
      "tags": [],
      "label": "IP Link",
      "visible": false,
      "policy_type_name": "AttachLogicalLink",
      "attributes": {
        "security_zone": "0cxTQTOoDWfzgFS0UH4",
        "interface_type": "tagged",
        "ipv6_addressing_type": "none",
        "vlan_id": 444,
        "ipv4_addressing_type": "numbered"
      },
      "id": "7411ef3a-05fd-4c2e-a11a-ec4636399969"
    },
    {
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table. Can be assigned to physical interfaces or single-chassis LAGs (not applicable to ESI LAG or MLAG interfaces).",
      "tags": [],
      "label": "IP Link (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": "1bbeaf90-ad40-4758-a3e6-15091635dd39",
        "first_subpolicy": "7411ef3a-05fd-4c2e-a11a-ec4636399969",
        "resolver": null
      },
      "id": "c303e784-82e5-48f5-9ddc-4242641e3770"
    },
    {
      "description": "Create a BGP peering session with a user-specified BGP neighbor addressed peer.",
      "tags": [],
      "label": "BGP Peering (IP Endpoint) (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": "643cbf38-f377-4ad4-ab70-b608d5472bc6",
        "first_subpolicy": "816c5d05-4694-4881-acbc-e2bd2c2884be",
        "resolver": null
      },
      "id": "2e264a0a-ca5b-413c-9bab-c8c1c4aa4e6e"
    },
    {
      "description": "Create a BGP peering session with a user-specified BGP neighbor addressed peer.",
      "tags": [],
      "label": "BGP Peering (IP Endpoint)",
      "visible": false,
      "policy_type_name": "AttachIpEndpointWithBgpNsxt",
      "attributes": {
        "ipv6_safi": false,
        "ipv6_addr": null,
        "keepalive_timer": null,
        "bfd": false,
        "local_asn": null,
        "ipv4_addr": "1.1.1.1",
        "ttl": 2,
        "neighbor_asn_type": "static",
        "password": null,
        "holdtime_timer": null,
        "asn": null,
        "ipv4_safi": true
      },
      "id": "816c5d05-4694-4881-acbc-e2bd2c2884be"
    },
    {
      "description": "",
      "tags": [],
      "user_data": "{\"isSausage\":true,\"positions\":{\"7411ef3a-05fd-4c2e-a11a-ec4636399969\":[290,80,1],\"816c5d05-4694-4881-acbc-e2bd2c2884be\":[290,150,2],\"44bd94f1-9730-4427-8a8c-e3ffdb85ea9a\":[290,220,3]}}",
      "label": "The New CT (3)",
      "visible": true,
      "policy_type_name": "batch",
      "attributes": {
        "subpolicies": [
          "c303e784-82e5-48f5-9ddc-4242641e3770"
        ]
      },
      "id": "fbda8b05-16b3-48e9-9a58-50fa662e06b9"
    },
    {
      "description": "Allocate routing policy to specific BGP sessions.",
      "tags": [],
      "label": "Routing Policy (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": null,
        "first_subpolicy": "44bd94f1-9730-4427-8a8c-e3ffdb85ea9a",
        "resolver": null
      },
      "id": "7ae346b7-9136-4d0a-acd1-37856c23b68f"
    },
    {
      "description": "Create a BGP peering session with a user-specified BGP neighbor addressed peer.",
      "tags": [],
      "label": "BGP Peering (IP Endpoint) (batch)",
      "visible": false,
      "policy_type_name": "batch",
      "attributes": {
        "subpolicies": [
          "7ae346b7-9136-4d0a-acd1-37856c23b68f"
        ]
      },
      "id": "643cbf38-f377-4ad4-ab70-b608d5472bc6"
    },
    {
      "description": "Allocate routing policy to specific BGP sessions.",
      "tags": [],
      "label": "Routing Policy",
      "visible": false,
      "policy_type_name": "AttachExistingRoutingPolicy",
      "attributes": {
        "rp_to_attach": "o-ob0kv9g1yniFpiTco"
      },
      "id": "44bd94f1-9730-4427-8a8c-e3ffdb85ea9a"
    }
  ]
}`

	var raw rawConnectivityTemplate
	err := json.Unmarshal([]byte(apiResponse), &raw)
	if err != nil {
		t.Fatal(err)
	}

	connectivityTemplate, err := raw.polish()
	if err != nil {
		t.Fatal(err)
	}

	_ = connectivityTemplate
}
