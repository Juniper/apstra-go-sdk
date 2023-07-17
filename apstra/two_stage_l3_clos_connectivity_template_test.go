package apstra

import (
	"encoding/json"
	"testing"
)

func TestParseCT(t *testing.T) {
	raw1 := `{
  "policies": [
    {
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table. Can be assigned to physical interfaces or single-chassis LAGs (not applicable to ESI LAG or MLAG interfaces).",
      "tags": [],
      "label": "IP Link (batch)",
      "visible": false,
      "policy_type_name": "batch",
      "attributes": {
        "Subpolicies": [
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
        "Subpolicies": [
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
        "Subpolicies": [
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

	raw2 := `{
  "policies": [
    {
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table.",
      "tags": [],
      "label": "IP Link",
      "visible": false,
      "policy_type_name": "AttachLogicalLink",
      "attributes": {
        "security_zone": "6k8Wmo0n1h5b_Mbnmbc",
        "interface_type": "untagged",
        "ipv6_addressing_type": "none",
        "vlan_id": null,
        "ipv4_addressing_type": "numbered"
      },
      "id": "73f80186-7943-4cd1-b009-a4c5fc2229b7"
    },
    {
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table.",
      "tags": [],
      "label": "IP Link (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": "30631fbc-f139-4012-8bd9-b1721e2bc40e",
        "first_subpolicy": "73f80186-7943-4cd1-b009-a4c5fc2229b7",
        "resolver": null
      },
      "id": "7621d11a-ca48-410f-a413-96505bd14812"
    },
    {
      "description": "test description",
      "tags": [
        "bar",
        "foo"
      ],
      "user_data": "{\"isSausage\":true,\"positions\":{\"73f80186-7943-4cd1-b009-a4c5fc2229b7\":[290,84,1],\"c30f9180-e3f8-4fe6-8386-e2926f918ec7\":[290,150,2],\"42b388ea-f27a-4d56-a222-3e23b2f98130\":[199,221,3],\"34a3cf1d-c3be-4ac8-96b9-9d4c23391d48\":[506,154,4],\"114e76ef-63e8-41e5-98b5-4d022abf0039\":[447,219,5]}}",
      "label": "test label",
      "visible": true,
      "policy_type_name": "batch",
      "attributes": {
        "Subpolicies": [
          "7621d11a-ca48-410f-a413-96505bd14812"
        ]
      },
      "id": "b7626bc2-7fba-41bc-99fd-ae35614f6639"
    },
    {
      "description": "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer).",
      "tags": [],
      "label": "BGP Peering (Generic System)",
      "visible": false,
      "policy_type_name": "AttachBgpOverSubinterfacesOrSvi",
      "attributes": {
        "ipv6_safi": false,
        "keepalive_timer": 10,
        "bfd": false,
        "peer_from": "loopback",
        "local_asn": null,
        "session_addressing_ipv4": "addressed",
        "session_addressing_ipv6": "none",
        "ttl": 2,
        "neighbor_asn_type": "static",
        "peer_to": "interface_or_ip_endpoint",
        "password": null,
        "holdtime_timer": 30,
        "ipv4_safi": true
      },
      "id": "c30f9180-e3f8-4fe6-8386-e2926f918ec7"
    },
    {
      "description": "Allocate routing policy to specific BGP sessions.",
      "tags": [],
      "label": "Routing Policy (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": null,
        "first_subpolicy": "42b388ea-f27a-4d56-a222-3e23b2f98130",
        "resolver": null
      },
      "id": "c46106e6-671c-44e1-ae8c-43835509fa5e"
    },
    {
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table.",
      "tags": [],
      "label": "IP Link (batch)",
      "visible": false,
      "policy_type_name": "batch",
      "attributes": {
        "Subpolicies": [
          "5760f3de-8228-46fd-958a-8f3ad1402814",
          "3862d50e-ea13-4d14-8852-fd041a2225ca"
        ]
      },
      "id": "30631fbc-f139-4012-8bd9-b1721e2bc40e"
    },
    {
      "description": "Create a static route to user defined subnet via next hop derived from either IP link or VN endpoint.",
      "tags": [],
      "label": "Static Route (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": null,
        "first_subpolicy": "34a3cf1d-c3be-4ac8-96b9-9d4c23391d48",
        "resolver": null
      },
      "id": "5760f3de-8228-46fd-958a-8f3ad1402814"
    },
    {
      "description": "Allocate routing policy to specific BGP sessions.",
      "tags": [],
      "label": "Routing Policy (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": null,
        "first_subpolicy": "114e76ef-63e8-41e5-98b5-4d022abf0039",
        "resolver": null
      },
      "id": "816f4aad-bee5-47a1-8ddf-30e04ba113ed"
    },
    {
      "description": "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer).",
      "tags": [],
      "label": "BGP Peering (Generic System) (pipeline)",
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "second_subpolicy": "bd2a85c1-07c5-41fb-93a1-dce39738a58c",
        "first_subpolicy": "c30f9180-e3f8-4fe6-8386-e2926f918ec7",
        "resolver": null
      },
      "id": "3862d50e-ea13-4d14-8852-fd041a2225ca"
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
      "id": "114e76ef-63e8-41e5-98b5-4d022abf0039"
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
      "id": "42b388ea-f27a-4d56-a222-3e23b2f98130"
    },
    {
      "description": "Create a static route to user defined subnet via next hop derived from either IP link or VN endpoint.",
      "tags": [],
      "label": "Static Route",
      "visible": false,
      "policy_type_name": "AttachStaticRoute",
      "attributes": {
        "share_ip_endpoint": true,
        "network": "5.0.0.0/8"
      },
      "id": "34a3cf1d-c3be-4ac8-96b9-9d4c23391d48"
    },
    {
      "description": "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer).",
      "tags": [],
      "label": "BGP Peering (Generic System) (batch)",
      "visible": false,
      "policy_type_name": "batch",
      "attributes": {
        "Subpolicies": [
          "c46106e6-671c-44e1-ae8c-43835509fa5e",
          "816f4aad-bee5-47a1-8ddf-30e04ba113ed"
        ]
      },
      "id": "bd2a85c1-07c5-41fb-93a1-dce39738a58c"
    }
  ]
}`

	for i, apiString := range []string{raw1, raw2} {
		var raw rawConnectivityTemplate
		err := json.Unmarshal([]byte(apiString), &raw)
		if err != nil {
			t.Fatalf("error in test case %d", i)
		}

		ids := raw.rootBatchIds()
		if len(ids) != 1 {
			t.Fatalf("expected 1 root batch ID, got %d", len(ids))
		}

		connectivityTemplate, err := raw.polish(ids[0])
		if err != nil {
			t.Fatal(err)
		}

		_ = connectivityTemplate
	}
}
