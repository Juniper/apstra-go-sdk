// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package resourcegroup

import (
	"encoding"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
	oenum "github.com/orsinium-labs/enum"
)

var (
	_ encoding.TextMarshaler   = (*ResourceGroup)(nil)
	_ encoding.TextUnmarshaler = (*ResourceGroup)(nil)
)

type ResourceGroup oenum.Member[string]

func (o ResourceGroup) String() string {
	return o.Value
}

func (o ResourceGroup) Type() *enum.ResourceType {
	switch ResourceGroups.Parse(o.Value) {
	case &AccessAccessIPv4:
		return &enum.ResourceTypeIPv4
	case &AccessAccessIPv6:
		return &enum.ResourceTypeIPv6
	case &AccessASN:
		return &enum.ResourceTypeASN
	case &AccessIPv4:
		return &enum.ResourceTypeIPv4
	case &AccessIPv6:
		return &enum.ResourceTypeIPv6
	case &EVPNL3VNI:
		return &enum.ResourceTypeVNI
	case &ExternalVNLocalVNIs:
		return &enum.ResourceTypeVNI
	case &GenericASN:
		return &enum.ResourceTypeASN
	case &GenericIPv4:
		return &enum.ResourceTypeIPv4
	case &GenericIPv6:
		return &enum.ResourceTypeIPv6
	case &LeafASN:
		return &enum.ResourceTypeASN
	case &LeafIPv4:
		return &enum.ResourceTypeIPv4
	case &LeafIPv6:
		return &enum.ResourceTypeIPv6
	case &LeafL3PeerLinkLinkIPv4:
		return &enum.ResourceTypeIPv4
	case &LeafL3PeerLinkLinkIPv6:
		return &enum.ResourceTypeIPv6
	case &LeafLeafIPv4:
		return &enum.ResourceTypeIPv4
	case &LeafLeafIPv6:
		return &enum.ResourceTypeIPv6
	case &MLAGDomainSVIIPv4:
		return &enum.ResourceTypeIPv4
	case &MLAGDomainSVIIPv6:
		return &enum.ResourceTypeIPv6
	case &SpineASN:
		return &enum.ResourceTypeASN
	case &SpineIPv4:
		return &enum.ResourceTypeIPv4
	case &SpineIPv6:
		return &enum.ResourceTypeIPv6
	case &SpineLeafIPv4:
		return &enum.ResourceTypeIPv4
	case &SpineLeafIPv6:
		return &enum.ResourceTypeIPv6
	case &SuperspineASN:
		return &enum.ResourceTypeASN
	case &SuperspineIPv4:
		return &enum.ResourceTypeIPv4
	case &SuperspineIPv6:
		return &enum.ResourceTypeIPv6
	case &SuperspineSpineIPv4:
		return &enum.ResourceTypeIPv4
	case &SuperspineSpineIPv6:
		return &enum.ResourceTypeIPv6
	case &ToGenericLinkIpv4:
		return &enum.ResourceTypeIPv4
	case &ToGenericLinkIpv6:
		return &enum.ResourceTypeIPv6
	case &VirtualNetworkSviIpv4:
		return &enum.ResourceTypeIPv4
	case &VirtualNetworkSviIpv6:
		return &enum.ResourceTypeIPv6
	case &VTEPIPv4:
		return &enum.ResourceTypeIPv4
	case &VTEPIPv6:
		return &enum.ResourceTypeIPv6
	case &VXLANVNIDs:
		return &enum.ResourceTypeVNI
	}

	return nil
}

func (o ResourceGroup) MarshalText() (text []byte, err error) {
	if ResourceGroups.Parse(o.Value) == nil {
		return nil, fmt.Errorf("attempt to marshal invalid value for ResourceGroup: %q", o.Value)
	}

	return []byte(o.Value), nil
}

func (o *ResourceGroup) UnmarshalText(text []byte) error {
	p := ResourceGroups.Parse(string(text))
	if p == nil {
		return fmt.Errorf("attempt to unmarshal invalid value for ResourceGroup: %q", string(text))
	}
	o.Value = p.Value
	return nil
}

var (
	AccessAccessIPv4       = ResourceGroup{Value: "access_l3_peer_link_link_ips"}
	AccessAccessIPv6       = ResourceGroup{Value: "ipv6_access_l3_peer_link_link_ips"}
	AccessASN              = ResourceGroup{Value: "access_asns"}
	AccessIPv4             = ResourceGroup{Value: "access_loopback_ips"}
	AccessIPv6             = ResourceGroup{Value: "access_loopback_ips_ipv6"}
	EVPNL3VNI              = ResourceGroup{Value: "evpn_l3_vnis"}
	ExternalVNLocalVNIs    = ResourceGroup{Value: "external_vn_local_vnis"}
	GenericASN             = ResourceGroup{Value: "generic_asns"}
	GenericIPv4            = ResourceGroup{Value: "generic_loopback_ips"}
	GenericIPv6            = ResourceGroup{Value: "generic_loopback_ips_ipv6"}
	LeafASN                = ResourceGroup{Value: "leaf_asns"}
	LeafIPv4               = ResourceGroup{Value: "leaf_loopback_ips"}
	LeafIPv6               = ResourceGroup{Value: "leaf_loopback_ips_ipv6"}
	LeafL3PeerLinkLinkIPv4 = ResourceGroup{Value: "leaf_l3_peer_link_link_ips"}
	LeafL3PeerLinkLinkIPv6 = ResourceGroup{Value: "ipv6_leaf_l3_peer_link_link_ips"}
	LeafLeafIPv4           = ResourceGroup{Value: "leaf_leaf_link_ips"}
	LeafLeafIPv6           = ResourceGroup{Value: "ipv6_leaf_leaf_link_ips"}
	MLAGDomainSVIIPv4      = ResourceGroup{Value: "mlag_domain_svi_subnets"}
	MLAGDomainSVIIPv6      = ResourceGroup{Value: "mlag_domain_svi_subnets_ipv6"}
	SpineASN               = ResourceGroup{Value: "spine_asns"}
	SpineIPv4              = ResourceGroup{Value: "spine_loopback_ips"}
	SpineIPv6              = ResourceGroup{Value: "spine_loopback_ips_ipv6"}
	SpineLeafIPv4          = ResourceGroup{Value: "spine_leaf_link_ips"}
	SpineLeafIPv6          = ResourceGroup{Value: "ipv6_spine_leaf_link_ips"}
	SuperspineASN          = ResourceGroup{Value: "superspine_asns"}
	SuperspineIPv4         = ResourceGroup{Value: "superspine_loopback_ips"}
	SuperspineIPv6         = ResourceGroup{Value: "superspine_loopback_ips_ipv6"}
	SuperspineSpineIPv4    = ResourceGroup{Value: "spine_superspine_link_ips"}
	SuperspineSpineIPv6    = ResourceGroup{Value: "ipv6_spine_superspine_link_ips"}
	ToGenericLinkIpv4      = ResourceGroup{Value: "to_generic_link_ips"}
	ToGenericLinkIpv6      = ResourceGroup{Value: "ipv6_to_generic_link_ips"}
	VirtualNetworkSviIpv4  = ResourceGroup{Value: "virtual_network_svi_subnets"}
	VirtualNetworkSviIpv6  = ResourceGroup{Value: "virtual_network_svi_subnets_ipv6"}
	VTEPIPv4               = ResourceGroup{Value: "vtep_ips"}
	VTEPIPv6               = ResourceGroup{Value: "vtep_ips_ipv6"}
	VXLANVNIDs             = ResourceGroup{Value: "vxlan_vn_ids"}

	ResourceGroups = oenum.New(
		AccessAccessIPv4,
		AccessAccessIPv6,
		AccessASN,
		AccessIPv4,
		AccessIPv6,
		EVPNL3VNI,
		ExternalVNLocalVNIs,
		GenericASN,
		GenericIPv4,
		GenericIPv6,
		LeafASN,
		LeafIPv4,
		LeafIPv6,
		LeafL3PeerLinkLinkIPv4,
		LeafL3PeerLinkLinkIPv6,
		LeafLeafIPv4,
		LeafLeafIPv6,
		MLAGDomainSVIIPv4,
		MLAGDomainSVIIPv6,
		SpineASN,
		SpineIPv4,
		SpineIPv6,
		SpineLeafIPv4,
		SpineLeafIPv6,
		SuperspineASN,
		SuperspineIPv4,
		SuperspineIPv6,
		SuperspineSpineIPv4,
		SuperspineSpineIPv6,
		ToGenericLinkIpv4,
		ToGenericLinkIpv6,
		VirtualNetworkSviIpv4,
		VirtualNetworkSviIpv6,
		VTEPIPv4,
		VTEPIPv6,
		VXLANVNIDs,
	)
)
