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

	xInitialPosition = 290
	yInitialPosition = 80
	xSpacing         = 200
	ySpacing         = 70
)

type XConnectivityTemplate struct {
	Id          *ObjectId
	Label       string
	Description string
	Subpolicies []*xConnectivityTemplatePrimitive // batch pointers
	Tags        []string
	UserData    *xConnectivityTemplatePrimitiveUserData
}

func (o *XConnectivityTemplate) Raw() ([]xRawConnectivityTemplatePrimitive, error) {
	err := o.SetId()
	if err != nil {
		return nil, err
	}

	subpolicyIds := make([]ObjectId, len(o.Subpolicies))
	for i, primitivePtr := range o.Subpolicies {
		err = primitivePtr.SetIds()
		if err != nil {
			return nil, err
		}

		subpolicyIds[i] = *primitivePtr.pipelineId
	}

	result, err := rawBatch(*o.Id, "", "", o.Subpolicies)
	if err != nil {
		return nil, err
	}

	if o.Tags == nil {
		o.Tags = []string{}
	}

	userData, err := json.Marshal(o.UserData)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling user data - %w", err)
	}

	// special handling for root batch fields
	result[0].Description = o.Description
	result[0].Label = o.Label
	result[0].Visible = true
	result[0].Tags = o.Tags
	result[0].UserData = string(userData)

	return result, nil
}

func (o *XConnectivityTemplate) SetId() error {
	if o.Id == nil {
		uuid, err := uuid1AsObjectId()
		if err != nil {
			return err
		}
		o.Id = &uuid
	}

	return nil
}

func (o *XConnectivityTemplate) SetUserData() {
	o.UserData = &xConnectivityTemplatePrimitiveUserData{
		IsSausage: true,
		Positions: make(map[ObjectId][]int),
	}

	for i, subpolicy := range o.Subpolicies {
		additionalPositions := subpolicy.positions(i*xSpacing+xInitialPosition, yInitialPosition)
		mergePositionMaps(&o.UserData.Positions, &additionalPositions)
	}
}

type xConnectivityTemplatePrimitiveUserData struct {
	IsSausage bool               `json:"isSausage"`
	Positions map[ObjectId][]int `json:"positions"`
}

type xConnectivityTemplatePrimitive struct {
	id          *ObjectId
	attributes  xConnectivityTemplateAttributes
	subpolicies []*xConnectivityTemplatePrimitive // batch pointers
	batchId     *ObjectId
	pipelineId  *ObjectId
}

func (o *xConnectivityTemplatePrimitive) positions(x, y int) map[ObjectId][]int {
	positions := make(map[ObjectId][]int)
	positions[*o.id] = []int{x, y, 1}
	for i, subpolicy := range o.subpolicies {
		additionalPositions := subpolicy.positions(x+i*xSpacing, y+ySpacing)
		mergePositionMaps(&positions, &additionalPositions)
	}
	return positions
}

// rawPipeline returns []xRawConnectivityTemplatePrimitive consisting of:
//   - a pipeline policy element
//   - the actual policy element
//   - if there are any children, a batch policy element containing downstream primitives
func (o *xConnectivityTemplatePrimitive) rawPipeline() ([]xRawConnectivityTemplatePrimitive, error) {
	if o.attributes == nil {
		return nil, errors.New("rawPipeline() invoked with nil attributes")
	}

	err := o.SetIds()
	if err != nil {
		return nil, err
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
		batchPolicies, err := rawBatch(*o.batchId, attributes.description(), attributes.label(), o.subpolicies)
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

type xRawConnectivityTemplatePrimitive struct {
	Id             ObjectId        `json:"id"`
	Label          string          `json:"label"`
	Description    string          `json:"description"`
	Tags           []string        `json:"tags"`
	UserData       string          `json:"user_data,omitempty"`
	Visible        bool            `json:"visible"`
	PolicyTypeName string          `json:"policy_type_name"`
	Attributes     json.RawMessage `json:"attributes"`
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

func rawBatch(id ObjectId, description, label string, subpolicies []*xConnectivityTemplatePrimitive) ([]xRawConnectivityTemplatePrimitive, error) {
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
		Tags:           []string{},
		Label:          label + policyTypeBatchSuffix,
		PolicyTypeName: policyTypeNameBatch,
		Attributes:     rawAttributes,
		Id:             id,
	}

	return append([]xRawConnectivityTemplatePrimitive{batch}, pipelines...), nil
}

func mergePositionMaps(dst, src *map[ObjectId][]int) {
	t := *dst
	for k, v := range *src {
		t[k] = v
	}
}
