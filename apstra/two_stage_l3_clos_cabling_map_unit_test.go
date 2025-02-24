// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "testing"

func TestSwitchLinkStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "ethernet", intType: InterfaceTypeEthernet, stringType: interfaceTypeEthernet},
		{stringVal: "ip", intType: InterfaceTypeIp, stringType: interfaceTypeIp},
		{stringVal: "loopback", intType: InterfaceTypeLoopback, stringType: interfaceTypeLoopback},
		{stringVal: "port_channel", intType: InterfaceTypePortChannel, stringType: interfaceTypePortChannel},
		{stringVal: "svi", intType: InterfaceTypeSvi, stringType: interfaceTypeSvi},
		{stringVal: "logical_vtep", intType: InterfaceTypeLogicalVtep, stringType: interfaceTypeLogicalVtep},
		{stringVal: "anycast_vtep", intType: InterfaceTypeAnycastVtep, stringType: interfaceTypeAnycastVtep},
		{stringVal: "unicast_vtep", intType: InterfaceTypeUnicastVtep, stringType: interfaceTypeUnicastVtep},
		{stringVal: "global_anycast_vtep", intType: InterfaceTypeGlobalAnycastVtep, stringType: interfaceTypeGlobalAnycastVtep},
		{stringVal: "subinterface", intType: InterfaceTypeSubinterface, stringType: interfaceTypeSubinterface},

		{stringVal: "", intType: InterfaceOperationStateNone, stringType: interfaceOperationStateNone},
		{stringVal: "up", intType: InterfaceOperationStateUp, stringType: interfaceOperationStateUp},
		{stringVal: "deduced_down", intType: InterfaceOperationStateDown, stringType: interfaceOperationStateDown},
		{stringVal: "admin_down", intType: InterfaceOperationStateAdminDown, stringType: interfaceOperationStateAdminDown},

		{stringVal: "access_l3_peer_link", intType: LinkRoleAccessL3PeerLink, stringType: linkRoleAccessL3PeerLink},
		{stringVal: "access_server", intType: LinkRoleAccessServer, stringType: linkRoleAccessServer},
		{stringVal: "leaf_access", intType: LinkRoleLeafAccess, stringType: linkRoleLeafAccess},
		{stringVal: "leaf_l2_server", intType: LinkRoleLeafL2Server, stringType: linkRoleLeafL2Server},
		{stringVal: "leaf_l3_peer_link", intType: LinkRoleLeafL3PeerLink, stringType: linkRoleLeafL3PeerLink},
		{stringVal: "leaf_l3_server", intType: LinkRoleLeafL3Server, stringType: linkRoleLeafL3Server},
		{stringVal: "leaf_leaf", intType: LinkRoleLeafLeaf, stringType: linkRoleLeafLeaf},
		{stringVal: "leaf_pair_access", intType: LinkRoleLeafPairAccess, stringType: linkRoleLeafPairAccess},
		{stringVal: "leaf_pair_access_pair", intType: LinkRoleLeafPairAccessPair, stringType: linkRoleLeafPairAccessPair},
		{stringVal: "leaf_pair_l2_server", intType: LinkRoleLeafPairL2Server, stringType: linkRoleLeafPairL2Server},
		{stringVal: "leaf_peer_link", intType: LinkRoleLeafPeerLink, stringType: linkRoleLeafPeerLink},
		{stringVal: "spine_leaf", intType: LinkRoleSpineLeaf, stringType: linkRoleSpineLeaf},
		{stringVal: "spine_superspine", intType: LinkRoleSpineSuperspine, stringType: linkRoleSpineSuperspine},
		{stringVal: "to_external_router", intType: LinkRoleToExternalRouter, stringType: linkRoleToExternalRouter},
		{stringVal: "to_generic", intType: LinkRoleToGeneric, stringType: linkRoleToGeneric},

		{stringVal: "", intType: SystemNodeRoleNone, stringType: systemNodeRoleNone},
		{stringVal: "access", intType: SystemNodeRoleAccess, stringType: systemNodeRoleAccess},
		{stringVal: "generic", intType: SystemNodeRoleGeneric, stringType: systemNodeRoleGeneric},
		{stringVal: "l3_server", intType: SystemNodeRoleL3Server, stringType: systemNodeRoleL3Server},
		{stringVal: "leaf", intType: SystemNodeRoleLeaf, stringType: systemNodeRoleLeaf},
		{stringVal: "remote_gateway", intType: SystemNodeRoleRemoteGateway, stringType: systemNodeRoleRemoteGateway},
		{stringVal: "spine", intType: SystemNodeRoleSpine, stringType: systemNodeRoleSpine},
		{stringVal: "superspine", intType: SystemNodeRoleSuperspine, stringType: systemNodeRoleSuperspine},

		{stringVal: "aggregate_link", intType: LinkTypeAggregateLink, stringType: linkTypeAggregateLink},
		{stringVal: "ethernet", intType: LinkTypeEthernet, stringType: linkTypeEthernet},
		{stringVal: "logical_link", intType: LinkTypeLogicalLink, stringType: linkTypeLogicalLink},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}
