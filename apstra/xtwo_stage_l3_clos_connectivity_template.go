package apstra

import (
	"encoding/json"
	"fmt"
)

const (
	policyTypeNameBatch           = "batch"
	policyTypeNamePipelinPipeline = "pipeline"
)

//	"AttachSingleVlan"
// "Virtual Network (Single)"
// "Add a single VLAN to interfaces, as tagged or untagged."

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

type XConnectivityTemplate struct {
	Id          *ObjectId
	SubPolicies []*xConnectivityTemplatePrimitive // batch pointers
	Tags        []string
	Label       string
}

func (o XConnectivityTemplate) Raw() {
	// todo set user data
}

type xConnectivityTemplatePrimitiveUserData struct {
	IsSausage bool  `json:"isSausage"`
	Positions []int `json:"positions"`
}

type xConnectivityTemplateAttributes interface {
	raw() (json.RawMessage, error)
	policyTypeName() string
	label() string
	description() string
}

type xConnectivityTemplatePrimitive struct {
	id          *ObjectId
	userData    *xConnectivityTemplatePrimitiveUserData
	attributes  xConnectivityTemplateAttributes
	subPolicies []*xConnectivityTemplatePrimitive // batch pointers
	batchId     *ObjectId
	pipelineId  *ObjectId
}

type xRawConnectivityTemplatePrimitive struct {
	Id             ObjectId                                `json:"id"`
	Label          string                                  `json:"label"`
	Description    string                                  `json:"description"`
	Tags           []string                                `json:"tags"`
	UserData       *xConnectivityTemplatePrimitiveUserData `json:"user_data,omitempty"`
	Visible        bool                                    `json:"visible"`
	PolicyTypeName string                                  `json:"policy_type_name"`
	Attributes     json.RawMessage                         `json:"attributes"`
}

//func (o xConnectivityTemplatePrimitive) id() ObjectId {
//	id, err := uuid.NewUUID()
//	id.ClockSequence()
//}

func (o xConnectivityTemplatePrimitive) raw(root bool, pl string, pd string, tags []string, batchMembers []xRawConnectivityTemplatePrimitive) ([]xRawConnectivityTemplatePrimitive, error) {
	if root && len(o.subPolicies) == 0 {
		// root batch should always have sub-policies
		return nil, fmt.Errorf("cannot render Connectivity Template with no primitives")
	}

	batchId := o.batchId
	if batchId == nil {
		uuid, err := uuid1AsObjectId()
		if err != nil {
			return nil, err
		}
		batchId = &uuid
	}

	batchPrimitive := xRawConnectivityTemplatePrimitive{
		Id:             *batchId,
		Label:          pl,
		Description:    pd,
		Tags:           tags,
		UserData:       nil, // todo fix after batch members fully populated
		Visible:        root,
		PolicyTypeName: policyTypeNameBatch,
		Attributes:     nil,
	}
	_ = batchPrimitive

	var result []xRawConnectivityTemplatePrimitive
	_ = result

	batch := xRawConnectivityTemplatePrimitive{
		Id:             "",
		Label:          "",
		Description:    "",
		Tags:           nil,
		UserData:       nil,
		Visible:        false,
		PolicyTypeName: "",
		Attributes:     nil,
	}
	_ = batch

	return nil, nil
}
