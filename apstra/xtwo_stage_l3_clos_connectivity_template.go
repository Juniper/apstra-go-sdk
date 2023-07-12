package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlBlueprintObjPolicyImport = apiUrlBlueprintById + apiUrlPathDelim + "obj-policy-import"

	apiUrlBlueprintObjPolicyExport     = apiUrlBlueprintById + apiUrlPathDelim + "obj-policy-export"
	apiUrlBlueprintObjPolicyExportById = apiUrlBlueprintObjPolicyExport + apiUrlPathDelim + "%s"

	apiUrlBlueprintEndpointPolicies   = apiUrlBlueprintById + apiUrlPathDelim + "endpoint-policies"
	apiUrlBlueprintEndpointPolicyById = apiUrlBlueprintEndpointPolicies + apiUrlPathDelim + "%s"

	deleteRecursive = "delete_recursive"

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

func (o *XConnectivityTemplate) raw() (*rawConnectivityTemplate, error) {
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

	rawPolicies, err := rawBatch(*o.Id, "", "", o.Subpolicies)
	if err != nil {
		return nil, err
	}

	if o.Tags == nil {
		o.Tags = []string{}
	}

	userDataBytes, err := json.Marshal(o.UserData)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling user data - %w", err)
	}
	userDataString := string(userDataBytes)

	// special handling for root batch fields
	rawPolicies[0].Description = o.Description
	rawPolicies[0].Label = o.Label
	rawPolicies[0].Visible = true
	rawPolicies[0].Tags = o.Tags
	rawPolicies[0].UserData = &userDataString

	return &rawConnectivityTemplate{
		Policies: rawPolicies,
	}, nil
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

type rawConnectivityTemplate struct {
	Policies []xRawConnectivityTemplatePolicy `json:"policies"`
}

func (o *rawConnectivityTemplate) rootBatch() (*xRawConnectivityTemplatePolicy, error) {
	rootBatchIdx := -1
	for i, rawPolicy := range o.Policies {
		switch {
		case rawPolicy.Visible && rootBatchIdx < 0:
			rootBatchIdx = i
		case rawPolicy.Visible && rootBatchIdx >= 0:
			return nil, fmt.Errorf(
				"cannot polish rawConnectivityTempalte when policy[%d] and policy[%d] both flagged \"visible\"",
				rootBatchIdx, i)
		}
	}

	if rootBatchIdx < 0 {
		return nil, fmt.Errorf("out of %d raw policies, none are flagged \"visible\"", len(o.Policies))
	}

	return &o.Policies[rootBatchIdx], nil
}

func (o *rawConnectivityTemplate) policyMap() map[ObjectId]xRawConnectivityTemplatePolicy {
	result := make(map[ObjectId]xRawConnectivityTemplatePolicy, len(o.Policies))
	for _, policy := range o.Policies {
		result[policy.Id] = policy
	}
	return result
}

func (o *rawConnectivityTemplate) polish() (*XConnectivityTemplate, error) {
	if len(o.Policies) == 0 {
		return nil, fmt.Errorf("cannot polish a rawConnectivityTemplate with no policies")
	}

	rootBatch, err := o.rootBatch()
	if err != nil {
		return nil, err
	}
	if rootBatch.UserData == nil {
		return nil, fmt.Errorf("connectivity template root batch has no user data")
	}

	policyMap := o.policyMap()
	delete(policyMap, rootBatch.Id)

	var userData xConnectivityTemplatePrimitiveUserData
	err = json.Unmarshal([]byte(*rootBatch.UserData), &userData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling root batch %q user data %q - %w",
			rootBatch.Id, *rootBatch.UserData, err)
	}

	var attributes xRawBatchattributes
	err = json.Unmarshal(rootBatch.Attributes, &attributes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling root batch %q attributes %q - %w",
			rootBatch.Id, *rootBatch.UserData, err)
	}

	subpolicies := make([]*xConnectivityTemplatePrimitive, len(attributes.Subpolicies))
	for i, policyId := range attributes.Subpolicies {
		subpolicies[i], err = parsePrimitiveTreeByPipelineId(policyId, policyMap)
		if err != nil {
			return nil, err
		}
	}

	return &XConnectivityTemplate{
		Id:          &rootBatch.Id,
		Label:       rootBatch.Label,
		Description: rootBatch.Description,
		Subpolicies: subpolicies,
		Tags:        rootBatch.Tags,
		UserData:    &userData,
	}, nil
}

type xConnectivityTemplatePrimitiveUserData struct {
	IsSausage bool               `json:"isSausage"`
	Positions map[ObjectId][]int `json:"positions"`
}

type xConnectivityTemplatePrimitive struct {
	id          *ObjectId
	attributes  xConnectivityTemplateAttributes
	subpolicies []*xConnectivityTemplatePrimitive // batch of pointers to pipelines
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

// rawPipeline returns []xRawConnectivityTemplatePolicy consisting of:
//   - a pipeline policy element
//   - the actual policy element
//   - if there are any children, a batch policy element containing downstream primitives
func (o *xConnectivityTemplatePrimitive) rawPipeline() ([]xRawConnectivityTemplatePolicy, error) {
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
	actual := xRawConnectivityTemplatePolicy{
		Description:    attributes.description(),
		Tags:           []string{}, // always empty slice
		Label:          attributes.label(),
		PolicyTypeName: attributes.policyTypeName().raw(),
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

	pipeline := xRawConnectivityTemplatePolicy{
		Description:    attributes.description(),
		Tags:           []string{}, // always empty slice
		Label:          attributes.label() + policyTypePipelineSuffix,
		PolicyTypeName: policyTypeNamePipeline,
		Attributes:     rawPipelineAttribtes,
		Id:             *o.pipelineId,
	}

	result := []xRawConnectivityTemplatePolicy{pipeline, actual}

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

type xRawConnectivityTemplatePolicy struct {
	Id             ObjectId                  `json:"id"`
	Label          string                    `json:"label"`
	Description    string                    `json:"description"`
	Tags           []string                  `json:"tags"`
	UserData       *string                   `json:"user_data,omitempty"`
	Visible        bool                      `json:"visible"`
	PolicyTypeName ctPrimitivePolicyTypeName `json:"policy_type_name"`
	Attributes     json.RawMessage           `json:"attributes"`
}

func (o xRawConnectivityTemplatePolicy) attributes() (xConnectivityTemplateAttributes, error) {
	var result xConnectivityTemplateAttributes

	switch o.PolicyTypeName {
	case ctPrimitivePolicyTypeNameAttachSingleVlan:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachSingleVlan)
	case ctPrimitivePolicyTypeNameAttachMultipleVLAN:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachMultipleVlan)
	case ctPrimitivePolicyTypeNameAttachLogicalLink:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachLogicalLink)
	case ctPrimitivePolicyTypeNameAttachStaticRoute:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachStaticRoute)
	case ctPrimitivePolicyTypeNameAttachCustomStaticRoute:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachCustomStaticRoute)
	case ctPrimitivePolicyTypeNameAttachIpEndpointWithBgpNsxt:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachIpEndpointWithBgpNsxt)
	case ctPrimitivePolicyTypeNameAttachBgpOverSubinterfacesOrSvi:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachBgpOverSubinterfacesOrSvi)
	case ctPrimitivePolicyTypeNameAttachBgpWithPrefixPeeringForSviOrSubinterface:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachBgpWithPrefixPeeringForSviOrSubinterface)
	case ctPrimitivePolicyTypeNameAttachExistingRoutingPolicy:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachExistingRoutingPolicy)
	case ctPrimitivePolicyTypeNameAttachRoutingZoneConstraint:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachRoutingZoneConstraint)
	default:
		return nil, fmt.Errorf("unhandled connectivity template type %q", o.PolicyTypeName)
	}

	return result, result.fromRawJson(o.Attributes)
}

type xRawBatchattributes struct {
	Subpolicies []ObjectId `json:"subpolicies"`
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

func rawBatch(id ObjectId, description, label string, subpolicies []*xConnectivityTemplatePrimitive) ([]xRawConnectivityTemplatePolicy, error) {
	// build downstream pipelines and collect their IDs
	var pipelines []xRawConnectivityTemplatePolicy
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

	batch := xRawConnectivityTemplatePolicy{
		Description:    description,
		Tags:           []string{},
		Label:          label + policyTypeBatchSuffix,
		PolicyTypeName: policyTypeNameBatch,
		Attributes:     rawAttributes,
		Id:             id,
	}

	return append([]xRawConnectivityTemplatePolicy{batch}, pipelines...), nil
}

func mergePositionMaps(dst, src *map[ObjectId][]int) {
	t := *dst
	for k, v := range *src {
		t[k] = v
	}
}

// parsePrimitiveTreeByPipelineId takes an entrypoint ObjectId representing a
// "pipeline" policy and a map of xRawConnectivityTemplatePolicy including
// the specified pipeline and all of its children.
//
// The returned *xConnectivityTemplatePrimitive is a tree built by recursive
// invocations of parsePrimitiveTreeByPipelineId until the tree is complete.
//
// parsePrimitiveTreeByPipelineId should be invoked once for each sub-policy in
// a connectivity template's root batch.
func parsePrimitiveTreeByPipelineId(pipelineId ObjectId, policyMap map[ObjectId]xRawConnectivityTemplatePolicy) (*xConnectivityTemplatePrimitive, error) {
	var actual, batch, pipeline xRawConnectivityTemplatePolicy
	var ok bool

	if pipeline, ok = policyMap[pipelineId]; !ok {
		return nil, fmt.Errorf("raw policy map doesn't include pipeline policy %q", pipelineId)
	}
	if pipeline.PolicyTypeName != policyTypeNamePipeline {
		return nil, fmt.Errorf("expected policy %q to be type %q, got %q",
			pipeline.Id, policyTypeNamePipeline, pipeline.PolicyTypeName)
	}

	var pipelineAttributes xRawPipelineAttributes
	err := json.Unmarshal(pipeline.Attributes, &pipelineAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling pipeline attributes %q for policy %q - %w",
			pipeline.Attributes, pipeline.Id, err)
	}

	if actual, ok = policyMap[pipelineAttributes.FirstSubpolicy]; !ok {
		return nil, fmt.Errorf("raw policy map doesn't include actual policy %q", pipelineAttributes.FirstSubpolicy)
	}
	var actualType CtPrimitivePolicyTypeName
	err = actualType.FromString(string(actual.PolicyTypeName))
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy type from CT policy %q - %w", actual.Id, err)
	}

	var batchId *ObjectId
	var subpolicies []*xConnectivityTemplatePrimitive
	if pipelineAttributes.SecondSubpolicy != nil {
		// a batch ID appears in the pipeline
		if batch, ok = policyMap[*pipelineAttributes.SecondSubpolicy]; ok {
			// the batch was found in the map
			if batch.PolicyTypeName != policyTypeNameBatch {
				// batch ID has wrong policy type (not batch)
				return nil, fmt.Errorf("expected policy %q to be type %q, got %q",
					batch.Id, policyTypeNameBatch, batch.PolicyTypeName)
			}

			var batchAttributes xRawBatchattributes
			err = json.Unmarshal(batch.Attributes, &batchAttributes)
			if err != nil {
				return nil, fmt.Errorf("failed unmarshaling batch attributes %q for policy %q - %w",
					batch.Attributes, batch.Id, err)
			}

			batchId = &batch.Id

			subpolicies = make([]*xConnectivityTemplatePrimitive, len(batchAttributes.Subpolicies))
			for i, subpolicyId := range batchAttributes.Subpolicies {
				subpolicies[i], err = parsePrimitiveTreeByPipelineId(subpolicyId, policyMap)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	attributes, err := actual.attributes()
	if err != nil {
		return nil, fmt.Errorf("failed to load attributes \"%s\" for policy %q - %w",
			actual.Attributes, actual.Id, err)
	}

	return &xConnectivityTemplatePrimitive{
		id:          &actual.Id,
		attributes:  attributes,
		subpolicies: subpolicies,
		batchId:     batchId,
		pipelineId:  &pipeline.Id,
	}, nil
}

func (o *TwoStageL3ClosClient) ListConnectivityTemplates(ctx context.Context) ([]ObjectId, error) {
	var apiResponse struct {
		Policies []xRawConnectivityTemplatePolicy `json:"policies"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodGet,
		urlStr:         fmt.Sprintf(apiUrlBlueprintObjPolicyExport, o.blueprintId),
		apiResponse:    &apiResponse,
		doNotLogin:     false,
		unsynchronized: false,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []ObjectId
	for _, policy := range apiResponse.Policies {
		if policy.Visible {
			result = append(result, policy.Id)
		}
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) CreateConnectivityTemplate(ctx context.Context, in *XConnectivityTemplate) error {
	apiInput, err := in.raw()
	if err != nil {
		return err
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPut,
		urlStr:         fmt.Sprintf(apiUrlBlueprintObjPolicyImport, o.blueprintId),
		apiInput:       &apiInput,
		apiResponse:    nil,
		doNotLogin:     false,
		unsynchronized: false,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) DeleteConnectivityTemplate(ctx context.Context, id ObjectId) error {
	urlStr := fmt.Sprintf(apiUrlBlueprintEndpointPolicyById, o.blueprintId, id)
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("failed parsing url %q - %w", urlStr, err)
	}

	params := urlObj.Query()
	params.Set(deleteRecursive, "true")
	urlObj.RawQuery = params.Encode()

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		url:    urlObj,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) getConnectivityTemplate(ctx context.Context, id ObjectId) (map[ObjectId]xRawConnectivityTemplatePolicy, error) {
	urlStr := fmt.Sprintf(apiUrlBlueprintObjPolicyExportById, o.blueprintId, id)

	var response struct {
		Policies []xRawConnectivityTemplatePolicy `json:"policies"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      urlStr,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make(map[ObjectId]xRawConnectivityTemplatePolicy, len(response.Policies))
	for _, policy := range response.Policies {
		result[policy.Id] = policy
	}

	if _, ok := result[id]; !ok {
		return nil, fmt.Errorf("policy %q not found in API response to GET %s", id, urlStr)
	}

	return result, nil
}
