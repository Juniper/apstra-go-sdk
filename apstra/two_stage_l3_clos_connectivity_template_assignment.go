package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintEndpointPolicyBatchApply = apiUrlBlueprintById + apiUrlPathDelim + "obj-policy-batch-apply"
)

// GetAllInterfacesConnectivityTemplates returns a map of ConnectivityTemplate
// IDs keyed by Interface (switch port) ID.
func (o *TwoStageL3ClosClient) GetAllInterfacesConnectivityTemplates(ctx context.Context) (map[ObjectId][]ObjectId, error) {
	queryResponse := new(struct {
		Items []struct {
			Interface struct {
				Id ObjectId `json:"id"`
			} `json:"interface"`
			EndpointPolicy struct {
				Id ObjectId `json:"id"`
			} `json:"ep_endpoint_policy"`
		} `json:"items"`
	})

	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetBlueprintType(BlueprintTypeStaging).
		SetClient(o.client).
		Node([]QEEAttribute{
			NodeTypeInterface.QEEAttribute(),
			{Key: "name", Value: QEStringVal(NodeTypeInterface.String())},
		}).
		Out([]QEEAttribute{RelationshipTypeEpMemberOf.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeEpGroup.QEEAttribute()}).
		In([]QEEAttribute{RelationshipTypeEpAffectedBy.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeEpApplicationInstance.QEEAttribute()}).
		Out([]QEEAttribute{RelationshipTypeEpTopLevel.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeEpEndpointPolicy.QEEAttribute(),
			{Key: "visible", Value: QEBoolVal(true)},
			{Key: "name", Value: QEStringVal(NodeTypeEpEndpointPolicy.String())},
		})

	err := query.Do(ctx, queryResponse)
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	intermediateResult := make(map[ObjectId]map[ObjectId]struct{})
	for _, item := range queryResponse.Items {
		if _, ok := intermediateResult[item.Interface.Id]; !ok {
			intermediateResult[item.Interface.Id] = make(map[ObjectId]struct{})
		}
		intermediateResult[item.Interface.Id][item.EndpointPolicy.Id] = struct{}{}
	}

	result := make(map[ObjectId][]ObjectId)
	for interfaceId, ctIds := range intermediateResult {
		for ctId := range ctIds {
			result[interfaceId] = append(result[interfaceId], ctId)
		}
	}

	return result, nil
}

func (o *TwoStageL3ClosClient) GetInterfaceConnectivityTemplates(ctx context.Context, intfId ObjectId) ([]ObjectId, error) {
	allMap, err := o.GetAllInterfacesConnectivityTemplates(ctx)
	if err != nil {
		return nil, err
	}

	if ctIds, ok := allMap[intfId]; ok {
		return ctIds, nil
	}

	// at this point it's unclear whether interface 'intfId' doesn't exist, or
	// merely has no CTs assigned. Returning a nil slice, along with the error
	// returned by GetNode will clear that up.
	return nil, o.client.GetNode(ctx, o.blueprintId, intfId, &struct{}{})
}

// SetInterfaceConnectivityTemplates assigns the listed ConnectivityTemplate IDs
// to the interface specified by intfId.
func (o *TwoStageL3ClosClient) SetInterfaceConnectivityTemplates(ctx context.Context, intfId ObjectId, ctIds []ObjectId) error {
	type policyInfo struct {
		PolicyId ObjectId `json:"policy"`
		Used     bool     `json:"used"`
	}

	type applicationPoints struct {
		InterfaceId ObjectId     `json:"id"`
		PolicyInfo  []policyInfo `json:"policies"`
	}

	appPoints := applicationPoints{
		InterfaceId: intfId,
		PolicyInfo:  make([]policyInfo, len(ctIds)),
	}
	for i, ctId := range ctIds {
		appPoints.PolicyInfo[i] = policyInfo{
			PolicyId: ctId,
			Used:     true,
		}
	}

	apiInput := struct {
		ApplicationPoints []applicationPoints `json:"application_points"`
	}{
		ApplicationPoints: []applicationPoints{appPoints},
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintEndpointPolicyBatchApply, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

// DelInterfaceConnectivityTemplates removes the listed ConnectivityTemplate IDs
// from the interface specified by intfId.
func (o *TwoStageL3ClosClient) DelInterfaceConnectivityTemplates(ctx context.Context, intfId ObjectId, ctIds []ObjectId) error {
	type policyInfo struct {
		PolicyId ObjectId `json:"policy"`
		Used     bool     `json:"used"`
	}

	type applicationPoints struct {
		InterfaceId ObjectId     `json:"id"`
		PolicyInfo  []policyInfo `json:"policies"`
	}

	appPoints := applicationPoints{
		InterfaceId: intfId,
		PolicyInfo:  make([]policyInfo, len(ctIds)),
	}
	for i, ctId := range ctIds {
		appPoints.PolicyInfo[i] = policyInfo{
			PolicyId: ctId,
			Used:     false,
		}
	}

	apiInput := struct {
		ApplicationPoints []applicationPoints `json:"application_points"`
	}{
		ApplicationPoints: []applicationPoints{appPoints},
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintEndpointPolicyBatchApply, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
