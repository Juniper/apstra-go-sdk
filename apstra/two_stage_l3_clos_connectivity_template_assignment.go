package apstra

import (
	"context"
	"errors"
	"fmt"
	"github.com/orsinium-labs/enum"
	"net/http"
	"net/url"
)

const (
	apiUrlBlueprintObjPolicyBatchApply        = apiUrlBlueprintById + apiUrlPathDelim + "obj-policy-batch-apply"
	apiUrlBlueprintObjPolicyApplicationPoints = apiUrlBlueprintById + apiUrlPathDelim + "obj-policy-application-points"
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

// SetApplicationPointConnectivityTemplates assigns the listed
// ConnectivityTemplate IDs to the application point specified by apId
func (o *TwoStageL3ClosClient) SetApplicationPointConnectivityTemplates(ctx context.Context, apId ObjectId, ctIds []ObjectId) error {
	type policyInfo struct {
		PolicyId ObjectId `json:"policy"`
		Used     bool     `json:"used"`
	}

	type applicationPoints struct {
		ApplicationPointId ObjectId     `json:"id"`
		PolicyInfo         []policyInfo `json:"policies"`
	}

	appPoints := applicationPoints{
		ApplicationPointId: apId,
		PolicyInfo:         make([]policyInfo, len(ctIds)),
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
		urlStr:   fmt.Sprintf(apiUrlBlueprintObjPolicyBatchApply, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

// DelApplicationPointConnectivityTemplates removes the listed
// ConnectivityTemplate IDs from the application point specified by apId
func (o *TwoStageL3ClosClient) DelApplicationPointConnectivityTemplates(ctx context.Context, apId ObjectId, ctIds []ObjectId) error {
	type policyInfo struct {
		PolicyId ObjectId `json:"policy"`
		Used     bool     `json:"used"`
	}

	type applicationPoint struct {
		ApplicationPointId ObjectId     `json:"id"`
		PolicyInfo         []policyInfo `json:"policies"`
	}

	appPoint := applicationPoint{
		ApplicationPointId: apId,
		PolicyInfo:         make([]policyInfo, len(ctIds)),
	}
	for i, ctId := range ctIds {
		appPoint.PolicyInfo[i] = policyInfo{
			PolicyId: ctId,
			Used:     false,
		}
	}

	apiInput := struct {
		ApplicationPoints []applicationPoint `json:"application_points"`
	}{
		ApplicationPoints: []applicationPoint{appPoint},
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintObjPolicyBatchApply, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

type applicationPointPolicyInfo struct {
	PolicyId ObjectId `json:"policy"`
	Used     bool     `json:"used"` // true: "assign" the policy to the
}

type applicationPointPolicyAssignment struct {
	Id         ObjectId                     `json:"id"`
	PolicyInfo []applicationPointPolicyInfo `json:"policies"`
}

// SetApplicationPointsConnectivityTemplates takes a map of application point ID
// (ObjectId) to map of policy obj (ObjectId) to bool. The application point ID
// (outer map key) is the ObjectID of a switch port interface, a system loopback
// or SVI interface, a system ID, or any other graph node ID which can serve as an
// "application point" for a connectivity template (endpoint policy ojb). The
// connectivity template ID (obj policy ID / inner map key) will be assigned or
// unassigned from the application point depending on the boolean value:
//
//	true: assign the policy to the application point
//	false: un-assign the policy from the application point
func (o *TwoStageL3ClosClient) SetApplicationPointsConnectivityTemplates(ctx context.Context, assignments map[ObjectId]map[ObjectId]bool) error {
	var apiInput struct {
		ApplicationPoints []applicationPointPolicyAssignment `json:"application_points"`
	}

	apiInput.ApplicationPoints = make([]applicationPointPolicyAssignment, len(assignments))
	var i int
	for apId, policyMap := range assignments {
		apiInput.ApplicationPoints[i].Id = apId
		apiInput.ApplicationPoints[i].PolicyInfo = make([]applicationPointPolicyInfo, len(policyMap))
		var j int
		for policyId, used := range policyMap {
			apiInput.ApplicationPoints[i].PolicyInfo[j] = applicationPointPolicyInfo{
				PolicyId: policyId,
				Used:     used,
			}
			j++
		}
		i++
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintObjPolicyBatchApply, o.blueprintId),
		apiInput: &apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

type appPointChildPolicyState enum.Member[string]

var (
	appPointChildPolicyStateUsedDirectly = appPointChildPolicyState{Value: "used-directly"}
	appPointChildPolicyStateUnused       = appPointChildPolicyState{Value: "unused"}
	appPointChildPolicyStateEnum         = enum.New(appPointChildPolicyStateUnused, appPointChildPolicyStateUsedDirectly)
)

// These types are the building blocks for the application point policy tree returned by the API.
// The tree is structured to accommodate the GUI's checkbox-based assignment scheme:
//
//	pod
//	  +-rack
//	  |  +-leaf
//	  |     +-interface
//	  |     +-interface
//	  |  +-leaf
//	  |     +-interface
//	  |     +-interface
//	  +-rack
//	  |  +-leaf
//	  |     +-interface
//	  |     +-interface
type appPointChildPolicy struct {
	Policy ObjectId `json:"policy"`
	State  string   `json:"state"`
}

type applicationPoint struct {
	Id                     ObjectId              `json:"id"`
	Label                  string                `json:"label"`
	Type                   string                `json:"type"`
	Tags                   []string              `json:"tags"`
	ChildApplicationPoints []applicationPoint    `json:"children"`
	ChildrenCount          int                   `json:"children_count"`
	Policies               []appPointChildPolicy `json:"policies"`
}

// fillMap fills the details of obj (endpoint?) policy usages into the supplied map.
// When force is false, fillMap will not add new entries to the map, using existing
// (outer) map entries to gauge caller interest in each application point (outer map
// key). When force is true, fillMap will add entries to the map.
func (o applicationPoint) fillMap(in map[ObjectId]map[ObjectId]bool, force bool) error {
	if in == nil {
		return errors.New("fillMap must not be called with a nil map")
	}

	if force && len(o.Policies) > 0 {
		// this application point has policies, and "force" is set. Add the application point ID
		// to the map so it *looks interesting*
		in[o.Id] = make(map[ObjectId]bool)
	}

	// only collect policy info from application points which have a map entry
	if _, ok := in[o.Id]; ok {
		for _, policyInfo := range o.Policies {
			// parse the raw `state` string from the API
			state := appPointChildPolicyStateEnum.Parse(policyInfo.State)
			if state == nil {
				return fmt.Errorf(
					"unknown application point policy state %q found at policy %q, application point %q",
					policyInfo.State, policyInfo.Policy, o.Id)
			}

			switch *state {
			case appPointChildPolicyStateUsedDirectly:
				in[o.Id][policyInfo.Policy] = true
			case appPointChildPolicyStateUnused:
				in[o.Id][policyInfo.Policy] = false
			}
		}
	}

	// having added (or not) this application point's policy info, recurse to child application points.
	for _, child := range o.ChildApplicationPoints {
		err := child.fillMap(in, force)
		if err != nil {
			return err
		}
	}

	return nil

}

type objPolicyApplicationPointsApiResponse struct {
	ApplicationPoints struct {
		ChildApplicationPoints []applicationPoint `json:"children"`
		ChildrenCount          int                `json:"children_count"`
	} `json:"application_points"`
	NodeTypes []string `json:"node_types"`
}

func (o *TwoStageL3ClosClient) GetAllApplicationPointsConnectivityTemplates(ctx context.Context) (map[ObjectId]map[ObjectId]bool, error) {
	var apiResponse objPolicyApplicationPointsApiResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintObjPolicyApplicationPoints, o.Id()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	// create the response struct.
	response := make(map[ObjectId]map[ObjectId]bool)
	for _, appPoint := range apiResponse.ApplicationPoints.ChildApplicationPoints {
		err = appPoint.fillMap(response, true)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (o *TwoStageL3ClosClient) GetConnectivityTemplatesByApplicationPoints(ctx context.Context, apIds []ObjectId) (map[ObjectId]map[ObjectId]bool, error) {
	var apiResponse objPolicyApplicationPointsApiResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintObjPolicyApplicationPoints, o.Id()),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	// create the response struct. We pre-populate the map with a nil slice to indicate
	// for each caller-supplied application point ID. Doing so indicates to fillmap()
	// which application points are interesting to the caller.
	response := make(map[ObjectId]map[ObjectId]bool, len(apIds))
	for _, apId := range apIds {
		response[apId] = make(map[ObjectId]bool)
	}

	for _, appPoint := range apiResponse.ApplicationPoints.ChildApplicationPoints {
		err = appPoint.fillMap(response, false)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (o *TwoStageL3ClosClient) GetApplicationPointConnectivityTemplates(ctx context.Context, apId ObjectId) (map[ObjectId]bool, error) {
	mapByApplicationPointId, err := o.GetConnectivityTemplatesByApplicationPoints(ctx, []ObjectId{apId})
	if err != nil {
		return nil, err
	}

	if entry, ok := mapByApplicationPointId[apId]; ok {
		return entry, nil
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("connectivity template usage map for interface %q not found", apId),
	}
}

func (o *TwoStageL3ClosClient) GetApplicationPointsConnectivityTemplatesByCt(ctx context.Context, ctId ObjectId) (map[ObjectId]map[ObjectId]bool, error) {
	var apiResponse objPolicyApplicationPointsApiResponse

	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintObjPolicyApplicationPoints, o.Id()))
	if err != nil {
		return nil, err
	}

	params := apstraUrl.Query()
	params.Set("policy", ctId.String())
	apstraUrl.RawQuery = params.Encode()

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make(map[ObjectId]map[ObjectId]bool)
	for _, x := range apiResponse.ApplicationPoints.ChildApplicationPoints {
		err = x.fillMap(result, true)
		if err != nil {
			return nil, err
		}
	}

	// clear unwanted CTs from the result (the API returns extra info despite the filter
	for k1, v1 := range result {
		for k2 := range v1 {
			if k2 != ctId {
				delete(result[k1], k2)
			}
		}
	}

	return result, nil
}
