package apstra

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	policyTypeNameBatch      = "batch"
	policyTypeNamePipeline   = "pipeline"
	policyTypeBatchSuffix    = " (" + policyTypeNameBatch + ")"
	policyTypePipelineSuffix = " (" + policyTypeNamePipeline + ")"
)

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

type xConnectivityTemplatePrimitive struct {
	id          *ObjectId
	userData    *xConnectivityTemplatePrimitiveUserData
	attributes  xConnectivityTemplateAttributes
	subpolicies []*xConnectivityTemplatePrimitive // batch pointers
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

type xRawPipelineAttributes struct {
	FirstSubpolicy  ObjectId    `json:"first_subpolicy"`
	SecondSubpolicy *ObjectId   `json:"second_subpolicy"`
	Resolver        interface{} `json:"resolver"` // what is this?
}

//func (o xConnectivityTemplatePrimitive) id() ObjectId {
//	id, err := uuid.NewUUID()
//	id.ClockSequence()
//}

func rawBatch(description, label string, tags []string, subpolicies []*xConnectivityTemplatePrimitive) ([]xRawConnectivityTemplatePrimitive, error) {
	id, err := uuid1AsObjectId()
	if err != nil {
		return nil, err
	}

	// build downstream pipelines and collect their IDs
	var pipelines []xRawConnectivityTemplatePrimitive
	subpolicyIds := make([]ObjectId, len(subpolicies))
	for i, subpolicy := range subpolicies {
		pipelineSlice, err := subpolicy.rawPipeline()
		if err != nil {
			return nil, err
		}

		subpolicyIds[i] = pipelineSlice[0].Id
		pipelines = append(pipelines, pipelineSlice...)
	}

	rawAttributes, err := json.Marshal(&struct {
		Subpolicies []ObjectId `json:"subpolicies"`
	}{
		Subpolicies: subpolicyIds,
	})
	if err != nil {
		return nil, fmt.Errorf("failed marshaling subpolicy ids for batch - %w", err)
	}

	batch := xRawConnectivityTemplatePrimitive{
		Description:    description,
		Tags:           tags,
		Label:          label,
		PolicyTypeName: policyTypeNameBatch,
		Attributes:     rawAttributes,
		Id:             id,
	}

	return append([]xRawConnectivityTemplatePrimitive{batch}, pipelines...), nil
}

// raw returns []xRawConnectivityTemplatePrimitive consisting of:
// - a pipeline policy element
// - the actual policy element
// - if there are any children, a batch policy element containing downstream primitives
func (o *xConnectivityTemplatePrimitive) rawPipeline() ([]xRawConnectivityTemplatePrimitive, error) {
	err := o.SetIds()
	if err != nil {
		return nil, err
	}

	if o.attributes == nil {
		return nil, errors.New("rawPipeline() invoked with nil attributes")
	}

	attributes := o.attributes
	rawAttributes, err := attributes.raw()
	if err != nil {
		return nil, err
	}

	// "actual"
	actual := xRawConnectivityTemplatePrimitive{
		Description:    attributes.description(),
		Tags:           []string{}, // always empty slice
		Label:          attributes.label(),
		PolicyTypeName: attributes.policyTypeName(),
		Attributes:     rawAttributes,
		Id:             *o.id,
	}

	var secondSubpolicy *ObjectId
	if len(o.subpolicies) > 0 {
		secondSubpolicy = o.batchId
	}

	pipelineAttributes := xRawPipelineAttributes{
		FirstSubpolicy:  *o.id,
		SecondSubpolicy: secondSubpolicy,
		Resolver:        nil,
	}
	rawPipelineAttribtes, err := json.Marshal(&pipelineAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling pipelineAttributes - %w", err)
	}

	pipeline := xRawConnectivityTemplatePrimitive{
		Description:    attributes.description(),
		Tags:           []string{}, // always empty slice
		Label:          attributes.label() + policyTypePipelineSuffix,
		PolicyTypeName: policyTypeNamePipeline,
		Attributes:     rawPipelineAttribtes,
		Id:             *o.pipelineId,
	}

	result := []xRawConnectivityTemplatePrimitive{pipeline, actual}

	if len(o.subpolicies) > 0 {
		batchPolicies, err := rawBatch(attributes.description(), attributes.label(), []string{}, o.subpolicies)
		if err != nil {
			return nil, err
		}
		result = append(result, batchPolicies...)
	}

	return result, nil
}

func (o *xConnectivityTemplatePrimitive) SetIds() error {
	if o.id == nil {
		uuid, err := uuid1AsObjectId()
		if err != nil {
			return err
		}
		o.id = &uuid
	}

	if o.pipelineId == nil {
		uuid := *o.id + policyTypePipelineSuffix
		o.pipelineId = &uuid
	}

	if o.batchId == nil && len(o.subpolicies) > 0 {
		uuid := *o.id + policyTypeBatchSuffix
		o.batchId = &uuid
	}

	return nil
}
