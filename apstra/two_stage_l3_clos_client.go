package apstra

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

const (
	blueprintTypeParam   = "type"
	dcClientMaxRetries   = 10
	dcClientRetryBackoff = 100 * time.Millisecond
)

type BlueprintType int
type blueprintType string

const (
	BlueprintTypeNone = BlueprintType(iota)
	BlueprintTypeConfig
	BlueprintTypeDeployed
	BlueprintTypeOperation
	BlueprintTypeStaging

	blueprintTypeNone      = blueprintType("")
	blueprintTypeConfig    = blueprintType("config")
	blueprintTypeDeployed  = blueprintType("deployed")
	blueprintTypeOperation = blueprintType("operation")
	blueprintTypeStaging   = blueprintType("staging")
	blueprintTypeUnknown   = "unknown-blueprint-type-%d"
)

func (o BlueprintType) raw() blueprintType {
	switch o {
	case BlueprintTypeNone:
		return blueprintTypeNone
	case BlueprintTypeConfig:
		return blueprintTypeConfig
	case BlueprintTypeDeployed:
		return blueprintTypeDeployed
	case BlueprintTypeStaging:
		return blueprintTypeStaging
	case BlueprintTypeOperation:
		return blueprintTypeOperation
	default:
		return blueprintType(fmt.Sprintf(blueprintTypeUnknown, o))
	}
}

func (o BlueprintType) string() string {
	return string(o.raw())
}

type TwoStageL3ClosClient struct {
	client        *Client
	blueprintId   ObjectId
	Mutex         Mutex
	blueprintType BlueprintType
}

// Id returns the client's Blueprint ID
func (o *TwoStageL3ClosClient) Id() ObjectId {
	return o.blueprintId
}

// SetType sets the client's internal BlueprintType value (staging, etc...).
// This value is in HTTP requests as a query string argument, e.g.
//
//	'?type=staging'
func (o *TwoStageL3ClosClient) SetType(bpt BlueprintType) {
	o.blueprintType = bpt
}

// urlWithParam is a helper function which uses the blueprintType element to
// decorate a *URL with the required query parameter.
//
//lint:ignore U1000 keep for future use
func (o *TwoStageL3ClosClient) urlWithParam(in string) (*url.URL, error) {
	apstraUrl, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	if o.blueprintType != BlueprintTypeNone {
		params := apstraUrl.Query()
		params.Set(blueprintTypeParam, o.blueprintType.string())
		apstraUrl.RawQuery = params.Encode()
	}
	return apstraUrl, nil
}

// GetResourceAllocations returns ResourceGroupAllocations representing
// all allocations of resource pools to blueprint requirements
func (o *TwoStageL3ClosClient) GetResourceAllocations(ctx context.Context) (ResourceGroupAllocations, error) {
	rawRgaSlice, err := o.getAllResourceAllocations(ctx)
	if err != nil {
		return nil, err
	}
	result := make(ResourceGroupAllocations, len(rawRgaSlice))
	for i, raw := range rawRgaSlice {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetResourceAllocation takes a *ResourceGroup and returns a
// *ResourceGroupAllocation with fields populated based on the Apstra API
// response.
func (o *TwoStageL3ClosClient) GetResourceAllocation(ctx context.Context, in *ResourceGroup) (*ResourceGroupAllocation, error) {
	rga, err := o.getResourceAllocation(ctx, in)
	if err != nil {
		return nil, err
	}
	return rga.polish()
}

// SetResourceAllocation sets the supplied resource allocation, overwriting any
// prior allocations with the supplied info.
func (o *TwoStageL3ClosClient) SetResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) error {
	return o.setResourceAllocation(ctx, in.raw())
}

// GetInterfaceMapAssignments returns a SystemIdToInterfaceMapAssignment (a map
// of string (blueprint graph node ID) to interface map ID detailing assignments
// in the specified blueprint:
//
//	x := SystemIdToInterfaceMapAssignment{
//		"BeAyAoCIgqx4r3hiFow": nil,
//		"B3Ym-PBJJEtvXQsnQQM": "VS_SONiC_BUZZNIK_PLUS__slicer-7x10-1",
//		"4gCWV2NRix6MYPm4PHU": "Arista_vEOS__slicer-7x10-1",
//	}
func (o *TwoStageL3ClosClient) GetInterfaceMapAssignments(ctx context.Context) (SystemIdToInterfaceMapAssignment, error) {
	return o.getInterfaceMapAssignments(ctx)
}

// SetInterfaceMapAssignments sets the supplied interface map assignments,
// overwriting any prior assignments with the supplied info. It returns
// the Blueprint config revision number.
func (o *TwoStageL3ClosClient) SetInterfaceMapAssignments(ctx context.Context, assignments SystemIdToInterfaceMapAssignment) error {
	return o.setInterfaceMapAssignments(ctx, assignments)
}

// CreateSecurityZone creates an Apstra Routing Zone / Security Zone / VRF
func (o *TwoStageL3ClosClient) CreateSecurityZone(ctx context.Context, cfg *SecurityZoneData) (ObjectId, error) {
	response, err := o.createSecurityZone(ctx, cfg.raw())
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// DeleteSecurityZone deletes an Apstra Routing Zone / Security Zone / VRF
func (o *TwoStageL3ClosClient) DeleteSecurityZone(ctx context.Context, zoneId ObjectId) error {
	return o.deleteSecurityZone(ctx, zoneId)
}

// GetSecurityZone fetches the Security Zone / Routing Zone / VRF with the given
// zoneId.
func (o *TwoStageL3ClosClient) GetSecurityZone(ctx context.Context, zoneId ObjectId) (*SecurityZone, error) {
	raw, err := o.getSecurityZone(ctx, zoneId)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetSecurityZoneByVrfName fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o *TwoStageL3ClosClient) GetSecurityZoneByVrfName(ctx context.Context, vrfName string) (*SecurityZone, error) {
	raw, err := o.getSecurityZoneByVrfName(ctx, vrfName)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetAllSecurityZones returns []SecurityZone representing all Security Zones /
// Routing Zones / VRFs on the system.
func (o *TwoStageL3ClosClient) GetAllSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	response, err := o.getAllSecurityZones(ctx)
	if err != nil {
		return nil, err
	}

	// This API endpoint returns a map. Convert to list for consistency with other 'GetAll' functions.
	result := make([]SecurityZone, len(response))
	var i int
	for k := range response {
		polished, err := response[k].polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
		i++
	}

	return result, nil
}

// UpdateSecurityZone replaces the configuration of zone zoneId with the supplied CreateSecurityZoneCfg
func (o *TwoStageL3ClosClient) UpdateSecurityZone(ctx context.Context, zoneId ObjectId, cfg *SecurityZoneData) error {
	return o.updateSecurityZone(ctx, zoneId, cfg.raw())
}

// GetAllPolicies returns []Policy representing all policies configured within the DC blueprint
func (o *TwoStageL3ClosClient) GetAllPolicies(ctx context.Context) ([]Policy, error) {
	return o.getAllPolicies(ctx)
}

// GetPolicy returns *Policy representing policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) GetPolicy(ctx context.Context, id ObjectId) (*Policy, error) {
	return o.getPolicy(ctx, id)
}

// CreatePolicy creates a policy within the DC blueprint, returns its ID
func (o *TwoStageL3ClosClient) CreatePolicy(ctx context.Context, policy *Policy) (ObjectId, error) {
	return o.createPolicy(ctx, policy)
}

// DeletePolicy deletes policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) DeletePolicy(ctx context.Context, id ObjectId) error {
	return o.deletePolicy(ctx, id)
}

// UpdatePolicy calls PUT to replace the configuration of policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) UpdatePolicy(ctx context.Context, id ObjectId, policy *Policy) error {
	return o.updatePolicy(ctx, id, policy)
}

// AddPolicyRule adds a policy rule at 'position' (bumping all other rules
// down). Position 0 makes the new policy first on the list, 1 makes it second
// on the list, etc... Use -1 for last on the list. The returned ObjectId
// represents the new rule
func (o *TwoStageL3ClosClient) AddPolicyRule(ctx context.Context, rule *PolicyRule, position int, policyId ObjectId) (ObjectId, error) {
	return o.addPolicyRule(ctx, rule, position, policyId)
}

// DeletePolicyRuleById deletes the given rule. If the rule doesn't exist, an
// ApstraClientErr with ErrNotFound is returned.
func (o *TwoStageL3ClosClient) DeletePolicyRuleById(ctx context.Context, policyId ObjectId, ruleId ObjectId) error {
	return o.deletePolicyRuleById(ctx, policyId, ruleId)
}

// CreateVirtualNetwork creates a new virtual network according to the supplied VirtualNetworkData
func (o *TwoStageL3ClosClient) CreateVirtualNetwork(ctx context.Context, in *VirtualNetworkData) (ObjectId, error) {
	return o.createVirtualNetwork(ctx, in.raw())
}

// ListAllVirtualNetworkIds returns []ObjectId representing virtual networks configured in the blueprint
func (o *TwoStageL3ClosClient) ListAllVirtualNetworkIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAllVirtualNetworkIds(ctx)
}

// GetVirtualNetwork returns *VirtualNetwork representing the given vnId
func (o *TwoStageL3ClosClient) GetVirtualNetwork(ctx context.Context, vnId ObjectId) (*VirtualNetwork, error) {
	raw, err := o.getVirtualNetwork(ctx, vnId)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// UpdateVirtualNetwork updates the virtual network specified by ID using the
// VirtualNetworkData and HTTP method PUT.
func (o *TwoStageL3ClosClient) UpdateVirtualNetwork(ctx context.Context, id ObjectId, cfg *VirtualNetworkData) error {
	return o.updateVirtualNetwork(ctx, id, cfg.raw())
}

// DeleteVirtualNetwork deletes the virtual network specified by id from the
// blueprint.
func (o *TwoStageL3ClosClient) DeleteVirtualNetwork(ctx context.Context, id ObjectId) error {
	return o.deleteVirtualNetwork(ctx, id)
}

// GetLockInfo returns *LockInfo describing the current state of the blueprint lock
func (o *TwoStageL3ClosClient) GetLockInfo(ctx context.Context) (*LockInfo, error) {
	li, err := o.getLockInfo(ctx)
	if err != nil {
		return nil, err
	}
	return li.polish()
}

// GetNodes fetches the node of the specified type, unpacks the API response
// into 'response'
func (o *TwoStageL3ClosClient) GetNodes(ctx context.Context, nodeType NodeType, response interface{}) error {
	return o.client.getNodes(ctx, o.blueprintId, nodeType, response)
}

// PatchNode patches (only submitted fields are changed) the specified node
// using the contents of 'request', the server's response (whole node info
// without map wrapper?) is returned in 'response'
func (o *TwoStageL3ClosClient) PatchNode(ctx context.Context, node ObjectId, request interface{}, response interface{}) error {
	return o.client.PatchNode(ctx, o.blueprintId, node, request, response)
}

// Client returns the embedded *Client
func (o *TwoStageL3ClosClient) Client() *Client {
	return o.client
}

// GetAllRoutingPolicies returns []DcRoutingPolicy representing all routing
// policies in the blueprint.
func (o *TwoStageL3ClosClient) GetAllRoutingPolicies(ctx context.Context) ([]DcRoutingPolicy, error) {
	raw, err := o.getAllRoutingPolicies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]DcRoutingPolicy, len(raw))
	for i := range raw {
		polished, err := raw[i].polish()
		if err != nil {
			return nil, err
		}

		result[i] = *polished
	}
	return result, nil
}

// GetRoutingPolicy returns *DcRoutingPolicy representing the specified policy.
func (o *TwoStageL3ClosClient) GetRoutingPolicy(ctx context.Context, id ObjectId) (*DcRoutingPolicy, error) {
	raw, err := o.getRoutingPolicy(ctx, id)
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

// GetDefaultRoutingPolicy returns *DcRoutingPolicy representing the
// "default_immutable" routing policy attached to the blueprint.
func (o *TwoStageL3ClosClient) GetDefaultRoutingPolicy(ctx context.Context) (*DcRoutingPolicy, error) {
	raw, err := o.getDefaultRoutingPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// CreateRoutingPolicy creates a blueprint routing policy according to the
// supplied *DcRoutingPolicyData.
func (o *TwoStageL3ClosClient) CreateRoutingPolicy(ctx context.Context, in *DcRoutingPolicyData) (ObjectId, error) {
	return o.createRoutingPolicy(ctx, in.raw())
}

// UpdateRoutingPolicy modifies the blueprint routing policy specified by 'id'
// according to the supplied *DcRoutingPolicyData.
func (o *TwoStageL3ClosClient) UpdateRoutingPolicy(ctx context.Context, id ObjectId, in *DcRoutingPolicyData) error {
	return o.updateRoutingPolicy(ctx, id, in.raw())
}

// DeleteRoutingPolicy deletes the routing policy specified by id.
func (o *TwoStageL3ClosClient) DeleteRoutingPolicy(ctx context.Context, id ObjectId) error {
	return o.deleteRoutingPolicy(ctx, id)
}