package apstra

import "encoding/json"

//	"AttachMultipleVLAN"
// "Virtual Network (Multiple)"
// "Add a list of VLANs to interfaces, as tagged or untagged."

//	"AttachLogicalLink"
// "IP Link"
// "Build an IP link between a fabric node and a generic system. This primitive uses AOS resource pool \"Link IPs - To Generic\" by default to dynamically allocate an IP endpoint (/31) on each side of the link. To allocate different IP endpoints, navigate under Routing Zone>Subinterfaces Table. Can be assigned to physical interfaces or single-chassis LAGs (not applicable to ESI LAG or MLAG interfaces)."

//	"AttachStaticRoute"
// "Static Route"
// "Create a static route to user defined subnet via next hop derived from either IP link or VN endpoint."

//	"AttachCustomStaticRoute"
// "Custom Static Route"
// "Create a static route with user defined next hop and destination network."

//	"AttachIpEndpointWithBgpNsxt"
// "BGP Peering (IP Endpoint)"
// "Create a BGP peering session with a user-specified BGP neighbor addressed peer."

// 	"AttachBgpOverSubinterfacesOrSvi"
// "BGP Peering (Generic System)"
// "Create a BGP peering session with Generic Systems inherited from AOS Generic System properties such as loopback and ASN (addressed, or link-local peer). Static route is automatically created when selecting loopback peering."

//	"AttachBgpWithPrefixPeeringForSviOrSubinterface"
// "Dynamic BGP Peering"
// "Configure dynamic BGP peering with IP prefix specified."

//	"AttachExistingRoutingPolicy"
// "Routing Policy"
// "Allocate routing policy to specific BGP sessions."

// "AttachRoutingZoneConstraint"
// "Routing Zone Constraint"
// "Assign a Routing Zone Constraint"

type xConnectivityTemplateAttributes interface {
	raw() (json.RawMessage, error)
	policyTypeName() string
	label() string
	description() string
}

var _ xConnectivityTemplateAttributes = ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{}

type ConnectivityTemplatePrimitiveAttributesAttachSingleVlan struct {
	Tagged   bool
	VnNodeId *ObjectId
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) raw() (json.RawMessage, error) {
	var tagType string
	if o.Tagged {
		tagType = "vlan_tagged"
	} else {
		tagType = "untagged"
	}

	raw := struct {
		TagType  string    `json:"tag_type"`
		VnNodeId *ObjectId `json:"vn_node_id"`
	}{
		TagType:  tagType,
		VnNodeId: o.VnNodeId,
	}

	return json.Marshal(&raw)
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) policyTypeName() string {
	return "AttachSingleVlan"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) label() string {
	return "Virtual Network (Single)"
}

func (o ConnectivityTemplatePrimitiveAttributesAttachSingleVlan) description() string {
	return "Add a single VLAN to interfaces, as tagged or untagged."

}
