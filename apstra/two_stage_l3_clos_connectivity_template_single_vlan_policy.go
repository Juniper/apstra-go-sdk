package apstra

import "encoding/json"

var _ TwoStageL3ClosObjPolicyAttributes = ObjPolicySingleVlanAttributes{}

type ObjPolicySingleVlanAttributes struct {
	VnNodeId ObjectId
	Tagged   bool
}

func (o ObjPolicySingleVlanAttributes) marshal() (json.RawMessage, error) {
	tagType := "untagged"
	if o.Tagged {
		tagType = "vlan_tagged"
	}

	raw := struct {
		VnNodeId ObjectId `json:"vn_node_id"`
		TagType  string   `json:"tag_type"`
	}{
		VnNodeId: o.VnNodeId,
		TagType:  tagType,
	}

	return json.Marshal(&raw)
}

func (o ObjPolicySingleVlanAttributes) typeName() string {
	return ObjPolicyTypeNameSingleVlan.String()
}
