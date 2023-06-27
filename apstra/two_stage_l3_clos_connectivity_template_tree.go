package apstra

import "encoding/json"

type ObjPolicyPipelineAttributes struct {
	FirstSubpolicy  *ObjectId `json:"first_subpolicy"`
	SecondSubpolicy *ObjectId `json:"second_subpolicy"`
}

func (o ObjPolicyPipelineAttributes) marshal() (json.RawMessage, error) {
	return json.Marshal(&o)
}

func (o ObjPolicyPipelineAttributes) typeName() string {
	return ObjPolicyTypeNamePipeline.String()
}

func (o ObjPolicyPipelineAttributes) pipeline() string {
	return ""
}

type ObjPolicyBatchAttributes struct {
	Subpolicies []ObjectId `json:"subpolicies"`
}

func (o ObjPolicyBatchAttributes) marshal() (json.RawMessage, error) {
	return json.Marshal(&o)
}

func (o ObjPolicyBatchAttributes) typeName() string {
	return ObjPolicyTypeNameBatch.String()
}

type TwoStageL3ClosObjPolicyUserData struct {
	IsSausage bool               `json:"isSausage"`
	Positions map[ObjectId][]int `json:"positions"`
}

//func pipeline(in ...ObjectId) ObjPolicyPipelineAttributes {
//	var sp1, sp2 *ObjectId
//	switch len(in) {
//	case 2:
//		sp1 = &in[0]
//		sp2 = &in[1]
//	case 1:
//		sp1 = &in[0]
//	}
//	return ObjPolicyPipelineAttributes{
//		FirstSubpolicy:  sp1,
//		SecondSubpolicy: sp2,
//	}
//}

//func batch(in ...ObjectId) ObjPolicyBatchAttributes {
//	return ObjPolicyBatchAttributes{
//		Subpolicies: in,
//	}
//}
