package apstra

import (
	"encoding/json"
	"testing"
)

func TestThing(t *testing.T) {
	vnNodeId := ObjectId("abc")

	xa := ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{
		Tagged:   true,
		VnNodeId: &vnNodeId,
	}

	x := xConnectivityTemplatePrimitive{
		id:          nil,
		attributes:  &xa,
		subpolicies: nil,
		batchId:     nil,
		pipelineId:  nil,
	}

	p, err := x.rawPipeline()
	if err != nil {
		t.Fatal(err)
	}
	_ = p
}

func TestBgpOverL3Connectivity(t *testing.T) {
	attachExistingRoutingPolicy := ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy{
		RpToAttach: "o-ob0kv9g1yniFpiTco",
	}
	rppPipelineId := ObjectId("8c654b0f-3253-45c6-9d8b-88bcc35fb70b")
	rppId := ObjectId("49f36469-7f10-4b10-9102-83654f3fe6a6")
	routingPolicyPrimitive := xConnectivityTemplatePrimitive{
		id:          &rppId,
		attributes:  attachExistingRoutingPolicy,
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
		TTL:                   2,
		BFD:                   false,
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
		attributes:  attachBgpOverSubinterfacesOrSvi,
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
		attributes:  attachLogicalLink,
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

	raw, err := ct.Raw()
	if err != nil {
		t.Fatal(err)
	}

	result, err := json.Marshal(&struct {
		Policies []xRawConnectivityTemplatePrimitive `json:"policies"`
	}{
		Policies: raw,
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = result
}
