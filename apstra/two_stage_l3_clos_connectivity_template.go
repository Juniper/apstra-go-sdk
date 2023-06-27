package apstra

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type ObjPolicyTypeName int
type objPolicyTypeName string

// 7 ObjPolicyTypeNameBgpOverSubinterfacesOrSvi // BGP Peering (Generic System)
// 8 ObjPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface // + Dynamic BGP Peering
// 5 ObjPolicyTypeNameCustomStaticRoute // Custom Static Route
// 9 ObjPolicyTypeNameRoutingPolicy // Routing Policy
// 6 ObjPolicyTypeNameIpEndpointWithBgpNsxt // BGP Peering (IP Endpoint)
// 3 ObjPolicyTypeNameLogicalLink // IP Link
// 2 ObjPolicyTypeNameMultipleVLAN //Virtual Network (Multiple)
// 10 ObjPolicyTypeNameRoutingZoneConstraint // Routing Zone Constraint
// 1 ObjPolicyTypeNameSingleVlan // Virtual Network (Single)
// 4 ObjPolicyTypeNameStaticRoute // Static Route

const (
	ObjPolicyTypeNameNone = ObjPolicyTypeName(iota)
	ObjPolicyTypeNameBatch
	ObjPolicyTypeNamePipeline
	ObjPolicyTypeNameBgpOverSubinterfacesOrSvi
	ObjPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface
	ObjPolicyTypeNameCustomStaticRoute
	ObjPolicyTypeNameRoutingPolicy
	ObjPolicyTypeNameIpEndpointWithBgpNsxt
	ObjPolicyTypeNameLogicalLink
	ObjPolicyTypeNameMultipleVLAN
	ObjPolicyTypeNameRoutingZoneConstraint
	ObjPolicyTypeNameSingleVlan
	ObjPolicyTypeNameStaticRoute
	ObjPolicyTypeNameUnknown = "unknown policy_type_name %q"

	objPolicyTypeNameNone                                     = objPolicyTypeName("")
	objPolicyTypeNameBatch                                    = objPolicyTypeName("batch")
	objPolicyTypeNamePipeline                                 = objPolicyTypeName("pipeline")
	objPolicyTypeNameBgpOverSubinterfacesOrSvi                = objPolicyTypeName("AttachBgpOverSubinterfacesOrSvi")
	objPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface = objPolicyTypeName("AttachBgpWithPrefixPeeringForSviOrSubinterface")
	objPolicyTypeNameCustomStaticRoute                        = objPolicyTypeName("AttachCustomStaticRoute")
	objPolicyTypeNameRoutingPolicy                            = objPolicyTypeName("AttachExistingRoutingPolicy")
	objPolicyTypeNameIpEndpointWithBgpNsxt                    = objPolicyTypeName("AttachIpEndpointWithBgpNsxt")
	objPolicyTypeNameLogicalLink                              = objPolicyTypeName("AttachLogicalLink")
	objPolicyTypeNameMultipleVLAN                             = objPolicyTypeName("AttachMultipleVLAN")
	objPolicyTypeNameRoutingZoneConstraint                    = objPolicyTypeName("AttachRoutingZoneConstraint")
	objPolicyTypeNameSingleVlan                               = objPolicyTypeName("AttachSingleVlan")
	objPolicyTypeNameStaticRoute                              = objPolicyTypeName("AttachStaticRoute")
	objPolicyTypeNameUnknown                                  = "unknown policy_type_name %d"
)

func (o ObjPolicyTypeName) Int() int {
	return int(o)
}

func (o ObjPolicyTypeName) String() string {
	switch o {
	case ObjPolicyTypeNameNone:
		return string(objPolicyTypeNameNone)
	case ObjPolicyTypeNameBatch:
		return string(objPolicyTypeNameBatch)
	case ObjPolicyTypeNamePipeline:
		return string(objPolicyTypeNamePipeline)
	case ObjPolicyTypeNameBgpOverSubinterfacesOrSvi:
		return string(objPolicyTypeNameBgpOverSubinterfacesOrSvi)
	case ObjPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface:
		return string(objPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface)
	case ObjPolicyTypeNameCustomStaticRoute:
		return string(objPolicyTypeNameCustomStaticRoute)
	case ObjPolicyTypeNameRoutingPolicy:
		return string(objPolicyTypeNameRoutingPolicy)
	case ObjPolicyTypeNameIpEndpointWithBgpNsxt:
		return string(objPolicyTypeNameIpEndpointWithBgpNsxt)
	case ObjPolicyTypeNameLogicalLink:
		return string(objPolicyTypeNameLogicalLink)
	case ObjPolicyTypeNameMultipleVLAN:
		return string(objPolicyTypeNameMultipleVLAN)
	case ObjPolicyTypeNameRoutingZoneConstraint:
		return string(objPolicyTypeNameRoutingZoneConstraint)
	case ObjPolicyTypeNameSingleVlan:
		return string(objPolicyTypeNameSingleVlan)
	case ObjPolicyTypeNameStaticRoute:
		return string(objPolicyTypeNameStaticRoute)
	default:
		return fmt.Sprintf(objPolicyTypeNameUnknown, o)
	}
}

func (o ObjPolicyTypeName) raw() objPolicyTypeName {
	return objPolicyTypeName(o.String())
}

func (o *ObjPolicyTypeName) FromString(in string) error {
	i, err := objPolicyTypeName(in).parse()
	if err != nil {
		return err
	}
	*o = ObjPolicyTypeName(i)
	return nil
}

func (o objPolicyTypeName) string() string {
	return string(o)
}

func (o objPolicyTypeName) parse() (int, error) {
	switch o {
	case objPolicyTypeNameNone:
		return int(ObjPolicyTypeNameNone), nil
	case objPolicyTypeNameBatch:
		return int(ObjPolicyTypeNameBatch), nil
	case objPolicyTypeNamePipeline:
		return int(ObjPolicyTypeNamePipeline), nil
	case objPolicyTypeNameBgpOverSubinterfacesOrSvi:
		return int(ObjPolicyTypeNameBgpOverSubinterfacesOrSvi), nil
	case objPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface:
		return int(ObjPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface), nil
	case objPolicyTypeNameCustomStaticRoute:
		return int(ObjPolicyTypeNameCustomStaticRoute), nil
	case objPolicyTypeNameRoutingPolicy:
		return int(ObjPolicyTypeNameRoutingPolicy), nil
	case objPolicyTypeNameIpEndpointWithBgpNsxt:
		return int(ObjPolicyTypeNameIpEndpointWithBgpNsxt), nil
	case objPolicyTypeNameLogicalLink:
		return int(ObjPolicyTypeNameLogicalLink), nil
	case objPolicyTypeNameMultipleVLAN:
		return int(ObjPolicyTypeNameMultipleVLAN), nil
	case objPolicyTypeNameRoutingZoneConstraint:
		return int(ObjPolicyTypeNameRoutingZoneConstraint), nil
	case objPolicyTypeNameSingleVlan:
		return int(ObjPolicyTypeNameSingleVlan), nil
	case objPolicyTypeNameStaticRoute:
		return int(ObjPolicyTypeNameStaticRoute), nil
	default:
		return 0, fmt.Errorf(ObjPolicyTypeNameUnknown, o)
	}
}

type TwoStageL3ClosObjPolicyAttributes interface {
	marshal() (json.RawMessage, error)
	typeName() string
}

type TwoStageL3ClosObjPolicy struct {
	Description    string
	Tags           []string
	UserData       *TwoStageL3ClosObjPolicyUserData
	Label          string
	PolicyTypeName ObjPolicyTypeName
	Attributes     TwoStageL3ClosObjPolicyAttributes
	Id             *ObjectId
	Children       []ObjectId
	pipeline       *RawTwoStageL3ClosObjPolicy
	batch          *RawTwoStageL3ClosObjPolicy
}

func (o *TwoStageL3ClosObjPolicy) Raw() ([]RawTwoStageL3ClosObjPolicy, error) {
	initUUID()

	mainRawAttributes, err := o.Attributes.marshal()
	if err != nil {
		return nil, fmt.Errorf("could not marshal CT policy element - %w", err)
	}

	if o.Id == nil {
		uuid1, err := uuid.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("could not generate UUID - %w", err)
		}
		id := ObjectId(uuid1.String())
		o.Id = &id
	}

	resultMain := RawTwoStageL3ClosObjPolicy{
		Description:    o.Description,
		Label:          o.Label,
		Attributes:     mainRawAttributes,
		PolicyTypeName: o.Attributes.typeName(),
		Id:             *o.Id,
	}

	err = o.buildPipeline()
	if err != nil {
		return nil, err
	}

	err = o.buildBatch()
	if err != nil {
		return nil, err
	}

	return []RawTwoStageL3ClosObjPolicy{resultMain, *o.pipeline, *o.batch}, nil
}

func (o *TwoStageL3ClosObjPolicy) buildPipeline() error {
	if o.Id == nil {
		return errors.New("attempt to generate pipeline policy object with nil ID")
	}

	if o.pipeline != nil {
		return errors.New("attempt to re-generate pipeline policy")
	}

	switch o.PolicyTypeName {
	case ObjPolicyTypeNamePipeline:
		fallthrough
	case ObjPolicyTypeNameBatch:
		fallthrough
	case ObjPolicyTypeNameNone:
		return fmt.Errorf("attempt to generate pipeline policy on a policy of type %q", o.PolicyTypeName)
	}

	uuid1, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("could not generate UUID - %w", err)
	}

	rawAttributes, err := json.Marshal(&ObjPolicyPipelineAttributes{FirstSubpolicy: o.Id})
	if err != nil {
		return err
	}

	o.pipeline = &RawTwoStageL3ClosObjPolicy{
		Description:    o.Description,
		Label:          o.Label + "(pipeline)",
		Attributes:     rawAttributes,
		PolicyTypeName: ObjPolicyTypeNamePipeline.String(),
		Id:             ObjectId(uuid1.String()),
	}

	return nil
}

func (o *TwoStageL3ClosObjPolicy) buildBatch() error {
	if o.batch != nil {
		return errors.New("attempt to re-generate batch policy")
	}

	switch o.PolicyTypeName {
	case ObjPolicyTypeNamePipeline:
		fallthrough
	case ObjPolicyTypeNameBatch:
		fallthrough
	case ObjPolicyTypeNameNone:
		return fmt.Errorf("attempt to generate pipeline policy on a policy of type %q", o.PolicyTypeName)
	}

	uuid1, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("could not generate UUID - %w", err)
	}

	var tags []string // starts as nil
	if o.Tags == nil {
		tags = []string{} // make an empty slice (not nil)
	} else {
		tags = o.Tags
	}

	rawAttributes, err := json.Marshal(&ObjPolicyBatchAttributes{Subpolicies: []ObjectId{}})
	if err != nil {
		return err
	}

	o.batch = &RawTwoStageL3ClosObjPolicy{
		Description:    o.Description,
		Tags:           tags,
		Label:          o.Label,
		PolicyTypeName: ObjPolicyTypeNameBatch.String(),
		Attributes:     rawAttributes,
		Id:             ObjectId(uuid1.String()),
	}

	return nil
}

type RawTwoStageL3ClosObjPolicy struct {
	Description    string                           `json:"description"`
	Tags           []string                         `json:"tags,omitempty"`
	UserData       *TwoStageL3ClosObjPolicyUserData `json:"user_data"`
	Label          string                           `json:"label"`
	Visible        bool                             `json:"visible"`
	PolicyTypeName string                           `json:"policy_type_name"`
	Attributes     json.RawMessage                  `json:"attributes"`
	Id             ObjectId                         `json:"id"`
}
