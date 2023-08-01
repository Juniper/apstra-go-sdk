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

type rawConnectivityTemplateState struct {
	Id                ObjectId          `json:"id"`
	Status            ctPrimitiveStatus `json:"status"`
	AppPointsCount    int               `json:"app_points_count"`
	MissingAttributes map[string]string `json:"missing_attributes"`
	Visible           bool              `json:"visible"`
}

func (o rawConnectivityTemplateState) polish() (*ConnectivityTemplateState, error) {
	if !o.Visible {
		return nil, fmt.Errorf("attempt to polish rawConnectivityTemplateState %q which is not visible", o.Id)
	}

	status, err := o.Status.parse()
	if err != nil {
		return nil, err
	}

	return &ConnectivityTemplateState{
		Id:                o.Id,
		Status:            CtPrimitiveStatus(status),
		AppPointsCount:    o.AppPointsCount,
		MissingAttributes: o.MissingAttributes,
	}, nil
}

type ConnectivityTemplateState struct {
	Id                ObjectId
	Status            CtPrimitiveStatus
	AppPointsCount    int
	MissingAttributes map[string]string
}

type ConnectivityTemplate struct {
	Id          *ObjectId
	Label       string
	Description string
	Subpolicies []*ConnectivityTemplatePrimitive // batch pointers
	Tags        []string
	UserData    *ConnectivityTemplatePrimitiveUserData
}

func (o *ConnectivityTemplate) raw() (*rawConnectivityTemplate, error) {
	err := o.SetIds()
	if err != nil {
		return nil, err
	}

	subpolicyIds := make([]ObjectId, len(o.Subpolicies))
	for i, primitivePtr := range o.Subpolicies {
		err = primitivePtr.setIds()
		if err != nil {
			return nil, err
		}

		subpolicyIds[i] = *primitivePtr.PipelineId
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

// SetIds walks the Connectivity Template tree and sets all "batch", "pipeline"
// and "actual" object IDs which aren't set, but need to be. Batch IDs where
// no children exist will not be set.
func (o *ConnectivityTemplate) SetIds() error {
	if o.Id == nil {
		uuid, err := uuid1AsObjectId()
		if err != nil {
			return err
		}
		o.Id = &uuid
	}

	for _, subpolicy := range o.Subpolicies {
		err := subpolicy.setIds()
		if err != nil {
			return err
		}
	}

	return nil
}

// SetUserData builds the top level `user_data` struct. It tries to lay the
// primitive "sausages" out sensibly.
func (o *ConnectivityTemplate) SetUserData() {
	o.UserData = &ConnectivityTemplatePrimitiveUserData{
		IsSausage: true,
		Positions: make(map[ObjectId][]float64),
	}

	for i, subpolicy := range o.Subpolicies {
		additionalPositions := subpolicy.positions(float64(i*xSpacing+xInitialPosition), yInitialPosition)
		mergePositionMaps(&o.UserData.Positions, &additionalPositions)
	}
}

type rawConnectivityTemplate struct {
	Policies []rawConnectivityTemplatePolicy `json:"policies"`
}

func (o *rawConnectivityTemplate) rootBatchIds() []ObjectId {
	var result []ObjectId

	for _, rawPolicy := range o.Policies {
		if rawPolicy.Visible {
			result = append(result, rawPolicy.Id)
		}
	}

	return result
}

func (o *rawConnectivityTemplate) policyMap() map[ObjectId]rawConnectivityTemplatePolicy {
	result := make(map[ObjectId]rawConnectivityTemplatePolicy, len(o.Policies))
	for _, policy := range o.Policies {
		result[policy.Id] = policy
	}
	return result
}

func (o *rawConnectivityTemplate) polish(id ObjectId) (*ConnectivityTemplate, error) {
	if len(o.Policies) == 0 {
		return nil, fmt.Errorf("cannot polish a rawConnectivityTemplate with no policies")
	}

	policyMap := o.policyMap()
	rootBatch, ok := policyMap[id]
	if !ok {
		return nil, fmt.Errorf("root batch policy %q not found", id)
	}
	if rootBatch.UserData == nil {
		return nil, fmt.Errorf("connectivity template root batch has no user data")
	}
	if policyTypeNameBatch != rootBatch.PolicyTypeName {
		return nil, fmt.Errorf("expected policy %q to be type %q, got %q",
			rootBatch.Id, policyTypeNameBatch, rootBatch.PolicyTypeName)
	}

	delete(policyMap, rootBatch.Id)

	var userData ConnectivityTemplatePrimitiveUserData
	err := json.Unmarshal([]byte(*rootBatch.UserData), &userData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling root batch %q user data %q - %w",
			rootBatch.Id, *rootBatch.UserData, err)
	}

	var attributes rawBatchAttributes
	err = json.Unmarshal(rootBatch.Attributes, &attributes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling root batch %q attributes %q - %w",
			rootBatch.Id, *rootBatch.UserData, err)
	}

	subpolicies := make([]*ConnectivityTemplatePrimitive, len(attributes.Subpolicies))
	for i, policyId := range attributes.Subpolicies {
		subpolicies[i], err = parsePrimitiveTreeByPipelineId(policyId, policyMap)
		if err != nil {
			return nil, err
		}
	}

	return &ConnectivityTemplate{
		Id:          &rootBatch.Id,
		Label:       rootBatch.Label,
		Description: rootBatch.Description,
		Subpolicies: subpolicies,
		Tags:        rootBatch.Tags,
		UserData:    &userData,
	}, nil
}

type ConnectivityTemplatePrimitiveUserData struct {
	IsSausage bool                   `json:"isSausage"`
	Positions map[ObjectId][]float64 `json:"positions"`
}

type ConnectivityTemplatePrimitive struct {
	Id          *ObjectId
	Attributes  ConnectivityTemplatePrimitiveAttributes
	Subpolicies []*ConnectivityTemplatePrimitive // batch of pointers to pipelines
	BatchId     *ObjectId
	PipelineId  *ObjectId
}

func (o *ConnectivityTemplatePrimitive) positions(x, y float64) map[ObjectId][]float64 {
	positions := make(map[ObjectId][]float64)
	positions[*o.Id] = []float64{x, y, 1}
	for i, subpolicy := range o.Subpolicies {
		additionalPositions := subpolicy.positions(x+float64(i*xSpacing), y+ySpacing)
		mergePositionMaps(&positions, &additionalPositions)
	}
	return positions
}

// rawPipeline returns []rawConnectivityTemplatePolicy consisting of:
//   - a pipeline policy element
//   - the actual policy element
//   - if there are any children, a batch policy element containing downstream primitives
func (o *ConnectivityTemplatePrimitive) rawPipeline() ([]rawConnectivityTemplatePolicy, error) {
	if o.Attributes == nil {
		return nil, errors.New("rawPipeline() invoked with nil attributes")
	}

	err := o.setIds()
	if err != nil {
		return nil, err
	}

	attributes := o.Attributes
	rawAttributes, err := attributes.raw()
	if err != nil {
		return nil, err
	}

	// "actual"
	actual := rawConnectivityTemplatePolicy{
		Description:    attributes.Description(),
		Tags:           []string{}, // always empty slice
		Label:          attributes.Label(),
		PolicyTypeName: attributes.PolicyTypeName().raw(),
		Attributes:     rawAttributes,
		Id:             *o.Id,
	}

	var secondSubpolicy *ObjectId
	if len(o.Subpolicies) > 0 {
		secondSubpolicy = o.BatchId
	}

	pipelineAttributes := rawPipelineAttributes{
		FirstSubpolicy:  *o.Id,
		SecondSubpolicy: secondSubpolicy,
		Resolver:        nil,
	}
	rawPipelineAttribtes, err := json.Marshal(&pipelineAttributes)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling pipelineAttributes - %w", err)
	}

	pipeline := rawConnectivityTemplatePolicy{
		Description:    attributes.Description(),
		Tags:           []string{}, // always empty slice
		Label:          attributes.Label() + policyTypePipelineSuffix,
		PolicyTypeName: policyTypeNamePipeline,
		Attributes:     rawPipelineAttribtes,
		Id:             *o.PipelineId,
	}

	result := []rawConnectivityTemplatePolicy{pipeline, actual}

	if len(o.Subpolicies) > 0 {
		batchPolicies, err := rawBatch(*o.BatchId, attributes.Description(), attributes.Label(), o.Subpolicies)
		if err != nil {
			return nil, err
		}
		result = append(result, batchPolicies...)
	}

	return result, nil
}

func (o *ConnectivityTemplatePrimitive) setIds() error {
	if o.Id == nil {
		uuid, err := uuid1AsObjectId()
		if err != nil {
			return err
		}
		o.Id = &uuid
	}

	if o.PipelineId == nil {
		uuid, err := uuid1AsObjectId()
		if err != nil {
			return err
		}
		o.PipelineId = &uuid
	}

	if o.BatchId == nil && len(o.Subpolicies) > 0 {
		uuid, err := uuid1AsObjectId()
		if err != nil {
			return err
		}
		o.BatchId = &uuid
	}

	for _, subpolicy := range o.Subpolicies {
		err := subpolicy.setIds()
		if err != nil {
			return err
		}
	}

	return nil
}

// A rawConnectivityTemplatePolicy is the base building block of a CT primitive
// (CT building block in the web UI) in the Apstra API. Each CT primitive is
// composed of 2 or 3 (when it has children) rawConnectivityTemplatePolicy
// structs.
//
// The Attributes element can take any of 12 forms: "pipeline", "batch", or
// one of the 10 implementations of ConnectivityTemplatePrimitiveAttributes.
// "pipeline" and "batch" are used to provide the tree structure which forms a
// CT as seen in the web UI.
type rawConnectivityTemplatePolicy struct {
	Id             ObjectId                  `json:"id"`
	Label          string                    `json:"label"`
	Description    string                    `json:"description"`
	Tags           []string                  `json:"tags"`
	UserData       *string                   `json:"user_data,omitempty"`
	Visible        bool                      `json:"visible"`
	PolicyTypeName ctPrimitivePolicyTypeName `json:"policy_type_name"`
	Attributes     json.RawMessage           `json:"attributes"`
}

func (o rawConnectivityTemplatePolicy) attributes() (ConnectivityTemplatePrimitiveAttributes, error) {
	var result ConnectivityTemplatePrimitiveAttributes

	switch o.PolicyTypeName {
	case ctPrimitivePolicyTypeNameAttachSingleVlan:
		result = new(ConnectivityTemplatePrimitiveAttributesAttachSingleVlan)
	case ctPrimitivePolicyTypeNameAttachMultipleVlan:
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

	err := result.fromRawJson(o.Attributes)
	if err != nil {
		return nil, err
	}

	return result, err
}

// rawBatchAttributes
// Each "batch" policy, including the root batch keeps a list of child policies.
// These sub-policies are identified by ID. Each one is a "pipeline" policy.
type rawBatchAttributes struct {
	Subpolicies []ObjectId `json:"subpolicies"`
}

// rawPipelineAttributes
// each "pipeline" policy identifies an actual CT primitive policy (VLAN, BGP
// stuff, static route, etc...) by ID in the FirstSubpolicy element. When the
// CT primitive has child primitives, the SecondSubpolicy element identifies a
// downstream "batch" policy. When there are no child primitives, SeconSubpolicy
// is nil.
type rawPipelineAttributes struct {
	FirstSubpolicy  ObjectId    `json:"first_subpolicy"`
	SecondSubpolicy *ObjectId   `json:"second_subpolicy"`
	Resolver        interface{} `json:"resolver"` // what is this?
}

func rawBatch(id ObjectId, description, label string, subpolicies []*ConnectivityTemplatePrimitive) ([]rawConnectivityTemplatePolicy, error) {
	// build downstream pipelines and collect their IDs
	var pipelines []rawConnectivityTemplatePolicy
	subpolicyIds := make([]ObjectId, len(subpolicies))
	for i, subpolicy := range subpolicies {
		pipelineSlice, err := subpolicy.rawPipeline()
		if err != nil {
			return nil, err
		}

		subpolicyIds[i] = pipelineSlice[0].Id
		pipelines = append(pipelines, pipelineSlice...)
	}

	rawAttributes, err := json.Marshal(&rawBatchAttributes{Subpolicies: subpolicyIds})
	if err != nil {
		return nil, fmt.Errorf("failed marshaling subpolicy ids for batch - %w", err)
	}

	batch := rawConnectivityTemplatePolicy{
		Description:    description,
		Tags:           []string{},
		Label:          label + policyTypeBatchSuffix,
		PolicyTypeName: policyTypeNameBatch,
		Attributes:     rawAttributes,
		Id:             id,
	}

	return append([]rawConnectivityTemplatePolicy{batch}, pipelines...), nil
}

func mergePositionMaps(dst, src *map[ObjectId][]float64) {
	t := *dst
	for k, v := range *src {
		t[k] = v
	}
}

// parsePrimitiveTreeByPipelineId takes an entrypoint ObjectId representing a
// "pipeline" policy and a map of rawConnectivityTemplatePolicy including
// the specified pipeline and all of its children.
//
// The returned *ConnectivityTemplatePrimitive is a tree built by recursive
// invocations of parsePrimitiveTreeByPipelineId until the tree is complete.
//
// parsePrimitiveTreeByPipelineId should be invoked once for each sub-policy in
// a connectivity template's root batch.
func parsePrimitiveTreeByPipelineId(pipelineId ObjectId, policyMap map[ObjectId]rawConnectivityTemplatePolicy) (*ConnectivityTemplatePrimitive, error) {
	var actual, batch, pipeline rawConnectivityTemplatePolicy
	var ok bool

	if pipeline, ok = policyMap[pipelineId]; !ok {
		return nil, fmt.Errorf("raw policy map doesn't include pipeline policy %q", pipelineId)
	}
	if pipeline.PolicyTypeName != policyTypeNamePipeline {
		return nil, fmt.Errorf("expected policy %q to be type %q, got %q",
			pipeline.Id, policyTypeNamePipeline, pipeline.PolicyTypeName)
	}

	var pipelineAttributes rawPipelineAttributes
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
	var subpolicies []*ConnectivityTemplatePrimitive
	if pipelineAttributes.SecondSubpolicy != nil {
		// a batch ID appears in the pipeline
		if batch, ok = policyMap[*pipelineAttributes.SecondSubpolicy]; ok {
			// the batch was found in the map
			if batch.PolicyTypeName != policyTypeNameBatch {
				// batch ID has wrong policy type (not batch)
				return nil, fmt.Errorf("expected policy %q to be type %q, got %q",
					batch.Id, policyTypeNameBatch, batch.PolicyTypeName)
			}

			var batchAttributes rawBatchAttributes
			err = json.Unmarshal(batch.Attributes, &batchAttributes)
			if err != nil {
				return nil, fmt.Errorf("failed unmarshaling batch attributes %q for policy %q - %w",
					batch.Attributes, batch.Id, err)
			}

			batchId = &batch.Id

			subpolicies = make([]*ConnectivityTemplatePrimitive, len(batchAttributes.Subpolicies))
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

	return &ConnectivityTemplatePrimitive{
		Id:          &actual.Id,
		Attributes:  attributes,
		Subpolicies: subpolicies,
		BatchId:     batchId,
		PipelineId:  &pipeline.Id,
	}, nil
}

func (o *TwoStageL3ClosClient) ListConnectivityTemplates(ctx context.Context) ([]ObjectId, error) {
	var apiResponse struct {
		Policies []rawConnectivityTemplatePolicy `json:"policies"`
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

func (o *TwoStageL3ClosClient) CreateConnectivityTemplate(ctx context.Context, in *ConnectivityTemplate) error {
	apiInput, err := in.raw()
	if err != nil {
		return err
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintObjPolicyImport, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) UpdateConnectivityTemplate(ctx context.Context, in *ConnectivityTemplate) error {
	return o.CreateConnectivityTemplate(ctx, in)
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

func (o *TwoStageL3ClosClient) GetConnectivityTemplate(ctx context.Context, id ObjectId) (*ConnectivityTemplate, error) {
	var response rawConnectivityTemplate

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintObjPolicyExportById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.polish(id)
}

func (o *TwoStageL3ClosClient) GetAllConnectivityTemplates(ctx context.Context) ([]ConnectivityTemplate, error) {
	var response rawConnectivityTemplate
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintObjPolicyExport, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	ids := response.rootBatchIds()
	result := make([]ConnectivityTemplate, len(ids))
	for i, id := range ids {
		polished, err := response.polish(id)
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) GetConnectivityTemplateState(ctx context.Context, id ObjectId) (*ConnectivityTemplateState, error) {
	var response struct {
		EndpointPolicy rawConnectivityTemplateState `json:"endpoint_policy"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintEndpointPolicyById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.EndpointPolicy.polish()
}

func (o *TwoStageL3ClosClient) GetAllConnectivityTemplateStates(ctx context.Context) ([]ConnectivityTemplateState, error) {
	var response struct {
		EndpointPolicies []rawConnectivityTemplateState `json:"endpoint_policies"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintEndpointPolicies, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]ConnectivityTemplateState, 0, len(response.EndpointPolicies)/3)
	for _, rawPolicy := range response.EndpointPolicies {
		if rawPolicy.Visible {
			polished, err := rawPolicy.polish()
			if err != nil {
				return nil, err
			}
			result = append(result, *polished)
		}
	}

	return result, nil
}
