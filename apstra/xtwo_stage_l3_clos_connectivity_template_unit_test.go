package apstra

import (
	"encoding/json"
	"testing"
)

func TestBgpOverL3Connectivity(t *testing.T) {
	expected := json.RawMessage(`{
  "policies": [
    {
      "id": "9f6c2ee4-a842-4fc2-979c-afce6c5f0ace",
      "label": "BGP over L3 connectivity",
      "description": "this is the description",
      "tags": [],
      "visible": true,
      "policy_type_name": "batch",
      "attributes": {
        "subpolicies": [
          "31e32ddd-98e9-4f74-8fd7-61bbf9501cfd"
        ]
      }
    },
    {
      "id": "31e32ddd-98e9-4f74-8fd7-61bbf9501cfd",
      "label": "IP Link (pipeline)",
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table.",
      "tags": [],
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "first_subpolicy": "bac16090-88ff-4f8b-9ee6-79b31078e123",
        "second_subpolicy": "e4f0ae44-871e-4002-806e-c61e647e5657",
        "resolver": null
      }
    },
    {
      "id": "bac16090-88ff-4f8b-9ee6-79b31078e123",
      "label": "IP Link",
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table.",
      "tags": [],
      "visible": false,
      "policy_type_name": "AttachLogicalLink",
      "attributes": {
        "security_zone": "6k8Wmo0n1h5b_Mbnmbc",
        "interface_type": "tagged",
        "vlan_id": 5,
        "ipv4_addressing_type": "numbered",
        "ipv6_addressing_type": "link_local"
      }
    },
    {
      "id": "e4f0ae44-871e-4002-806e-c61e647e5657",
      "label": "IP Link (batch)",
      "description": "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table.",
      "tags": [],
      "visible": false,
      "policy_type_name": "batch",
      "attributes": {
        "subpolicies": [
          "de1474c2-f892-4fa6-bef4-e330ae7f9ac7"
        ]
      }
    },
    {
      "id": "de1474c2-f892-4fa6-bef4-e330ae7f9ac7",
      "label": "BGP Peering (Generic System) (pipeline)",
      "description": "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer).",
      "tags": [],
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "first_subpolicy": "498b2502-e062-414b-b401-4e88a08ae8c5",
        "second_subpolicy": "83bd1635-e543-4752-b526-e290b8771285",
        "resolver": null
      }
    },
    {
      "id": "498b2502-e062-414b-b401-4e88a08ae8c5",
      "label": "BGP Peering (Generic System)",
      "description": "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer).",
      "tags": [],
      "visible": false,
      "policy_type_name": "AttachBgpOverSubinterfacesOrSvi",
      "attributes": {
        "ipv4_safi": true,
        "ipv6_safi": false,
        "ttl": 2,
        "bfd": false,
        "password": "foo",
        "keepalive_timer": 10,
        "holdtime_timer": 30,
        "local_asn": 55,
        "neighbor_asn_type": "static",
        "peer_from": "interface",
        "peer_to": "interface_or_ip_endpoint",
        "session_addressing_ipv4": "addressed",
        "session_addressing_ipv6": "link_local"
      }
    },
    {
      "id": "83bd1635-e543-4752-b526-e290b8771285",
      "label": "BGP Peering (Generic System) (batch)",
      "description": "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer).",
      "tags": [],
      "visible": false,
      "policy_type_name": "batch",
      "attributes": {
        "subpolicies": [
          "8c654b0f-3253-45c6-9d8b-88bcc35fb70b"
        ]
      }
    },
    {
      "id": "8c654b0f-3253-45c6-9d8b-88bcc35fb70b",
      "label": "Routing Policy (pipeline)",
      "description": "Allocate routing policy to specific BGP sessions.",
      "tags": [],
      "visible": false,
      "policy_type_name": "pipeline",
      "attributes": {
        "first_subpolicy": "49f36469-7f10-4b10-9102-83654f3fe6a6",
        "second_subpolicy": null,
        "resolver": null
      }
    },
    {
      "id": "49f36469-7f10-4b10-9102-83654f3fe6a6",
      "label": "Routing Policy",
      "description": "Allocate routing policy to specific BGP sessions.",
      "tags": [],
      "visible": false,
      "policy_type_name": "AttachExistingRoutingPolicy",
      "attributes": {
        "rp_to_attach": "o-ob0kv9g1yniFpiTco"
      }
    }
  ]
}`)

	expectedUserData := "{\"isSausage\":true,\"positions\":{\"bac16090-88ff-4f8b-9ee6-79b31078e123\":[290,80,1],\"498b2502-e062-414b-b401-4e88a08ae8c5\":[290,150,1],\"49f36469-7f10-4b10-9102-83654f3fe6a6\":[290,220,1]}}"

	rpId := "o-ob0kv9g1yniFpiTco"
	attachExistingRoutingPolicy := ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy{
		RpToAttach: &rpId,
	}
	rppPipelineId := ObjectId("8c654b0f-3253-45c6-9d8b-88bcc35fb70b")
	rppId := ObjectId("49f36469-7f10-4b10-9102-83654f3fe6a6")
	routingPolicyPrimitive := xConnectivityTemplatePrimitive{
		id:          &rppId,
		attributes:  &attachExistingRoutingPolicy,
		subpolicies: nil,
		batchId:     nil,
		pipelineId:  &rppPipelineId,
	}

	bgpPassword := "foo"
	keepalive := uint16(10)
	holdtime := uint16(30)
	localAsn := uint32(55)
	attachBgpOverSubinterfacesOrSvi := ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi{
		Ipv4Safi:              true,
		Ipv6Safi:              false,
		Ttl:                   2,
		Bfd:                   false,
		Password:              &bgpPassword,
		Keepalive:             &keepalive,
		Holdtime:              &holdtime,
		SessionAddressingIpv4: CtPrimitiveIPv4ProtocolSessionAddressingAddressed,
		SessionAddressingIpv6: CtPrimitiveIPv6ProtocolSessionAddressingLinkLocal,
		LocalAsn:              &localAsn,
		PeerFromLoopback:      false,
		PeerTo:                CtPrimitiveBgpPeerToInterfaceOrIpEndpoint,
		NeighborAsnDynamic:    false,
	}
	bgpPipelineId := ObjectId("de1474c2-f892-4fa6-bef4-e330ae7f9ac7")
	bgpId := ObjectId("498b2502-e062-414b-b401-4e88a08ae8c5")
	bgpBatchId := ObjectId("83bd1635-e543-4752-b526-e290b8771285")
	bgpPrimitive := xConnectivityTemplatePrimitive{
		id:          &bgpId,
		attributes:  &attachBgpOverSubinterfacesOrSvi,
		subpolicies: []*xConnectivityTemplatePrimitive{&routingPolicyPrimitive},
		batchId:     &bgpBatchId,
		pipelineId:  &bgpPipelineId,
	}

	securityZone := ObjectId("6k8Wmo0n1h5b_Mbnmbc")
	vlan := Vlan(5)
	attachLogicalLink := ConnectivityTemplatePrimitiveAttributesAttachLogicalLink{
		SecurityZone:            &securityZone,
		Tagged:                  true,
		Vlan:                    &vlan,
		IPv4AddressingNumbered:  true,
		IPv6AddressingLinkLocal: true,
	}
	IpLinkPipelineId := ObjectId("31e32ddd-98e9-4f74-8fd7-61bbf9501cfd")
	IpLinkId := ObjectId("bac16090-88ff-4f8b-9ee6-79b31078e123")
	IpLinkBatchId := ObjectId("e4f0ae44-871e-4002-806e-c61e647e5657")
	IpLinkPrimitive := xConnectivityTemplatePrimitive{
		id:          &IpLinkId,
		attributes:  &attachLogicalLink,
		subpolicies: []*xConnectivityTemplatePrimitive{&bgpPrimitive},
		batchId:     &IpLinkBatchId,
		pipelineId:  &IpLinkPipelineId,
	}

	ctId := ObjectId("9f6c2ee4-a842-4fc2-979c-afce6c5f0ace")
	ct := XConnectivityTemplate{
		Id:          &ctId,
		Subpolicies: []*xConnectivityTemplatePrimitive{&IpLinkPrimitive},
		Tags:        nil,
		Label:       "BGP over L3 connectivity",
		Description: "this is the description",
	}
	ct.SetUserData()

	raw, err := ct.raw()
	if err != nil {
		t.Fatal(err)
	}

	resultUserData := raw.Policies[0].UserData
	raw.Policies[0].UserData = nil

	result, err := json.Marshal(&struct {
		Policies []xRawConnectivityTemplatePolicy `json:"policies"`
	}{
		Policies: raw.Policies,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !jsonEqual(t, expected, result) {
		t.Fatalf("expected:\n %s\n\n got:\n%s", expected, result)
	}

	if !jsonEqual(t, json.RawMessage(expectedUserData), json.RawMessage(*resultUserData)) {
		t.Fatalf("expected:\n %s\n\n got:\n%s", expectedUserData, *resultUserData)
	}
}
