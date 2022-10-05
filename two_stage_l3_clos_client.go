package goapstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
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
	client      *Client
	blueprintId ObjectId
	mutex       *TwoStageL3ClosMutex
}

// Id returns the client's Blueprint ID
func (o *TwoStageL3ClosClient) Id() ObjectId {
	return o.blueprintId
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
	return o.setResourceAllocation(ctx, in)
}

// GetInterfaceMapAssignments returns a SystemIdToInterfaceMapAssignment (a map
// of string (blueprint graph node ID) to interface map ID detailing assignments
// in the specified blueprint:
// 	x := SystemIdToInterfaceMapAssignment{
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
func (o *TwoStageL3ClosClient) CreateSecurityZone(ctx context.Context, cfg *CreateSecurityZoneCfg) (ObjectId, error) {
	response, err := o.createSecurityZone(ctx, cfg)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// DeleteSecurityZone deletes an Apstra Routing Zone / Security Zone / VRF
func (o *TwoStageL3ClosClient) DeleteSecurityZone(ctx context.Context, zoneId ObjectId) error {
	return o.deleteSecurityZone(ctx, zoneId)
}

// GetSecurityZones returns all Apstra Routing Zones / Security Zones / VRFs
// associated with the specified blueprint
func (o *TwoStageL3ClosClient) GetSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	return o.getAllSecurityZones(ctx)
}

// GetSecurityZone fetches the Security Zone / Routing Zone / VRF with the given
// zoneId.
func (o *TwoStageL3ClosClient) GetSecurityZone(ctx context.Context, zoneId ObjectId) (*SecurityZone, error) {
	return o.getSecurityZone(ctx, zoneId)
}

// GetSecurityZoneByName fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o *TwoStageL3ClosClient) GetSecurityZoneByName(ctx context.Context, label string) (*SecurityZone, error) {
	return o.getSecurityZoneByName(ctx, label)
}

// GetAllSecurityZones returns []SecurityZone representing all Security Zones /
// Routing Zones / VRFs on the system.
func (o *TwoStageL3ClosClient) GetAllSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	return o.getAllSecurityZones(ctx)
}

// UpdateSecurityZone replaces the configuration of zone zoneId with the supplied CreateSecurityZoneCfg
func (o *TwoStageL3ClosClient) UpdateSecurityZone(ctx context.Context, zoneId ObjectId, cfg *CreateSecurityZoneCfg) error {
	return o.updateSecurityZone(ctx, zoneId, cfg)
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

// ListAllVirtualNetworkIds returns []ObjectId representing virtual networks configured in the blueprint
func (o *TwoStageL3ClosClient) ListAllVirtualNetworkIds(ctx context.Context, bpType BlueprintType) ([]ObjectId, error) {
	return o.listAllVirtualNetworkIds(ctx, bpType)
}

// GetVirtualNetwork returns *VirtualNetwork representing the given vnId within the blueprint type
func (o *TwoStageL3ClosClient) GetVirtualNetwork(ctx context.Context, vnId ObjectId, bpType BlueprintType) (*VirtualNetwork, error) {
	return o.getVirtualNetwork(ctx, vnId, bpType)
}

// GetVirtualNetworkBySubnet returns *VirtualNetwork representing the given desiredNet within the blueprint type
func (o *TwoStageL3ClosClient) GetVirtualNetworkBySubnet(ctx context.Context, desiredNet *net.IPNet, vrf ObjectId, bpType BlueprintType) (*VirtualNetwork, error) {
	return o.getVirtualNetworkBySubnet(ctx, desiredNet, vrf, bpType)
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
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintNodes, o.blueprintId))
	if err != nil {
		return err
	}

	if nodeType != NodeTypeNone {
		params := apstraUrl.Query()
		params.Set(nodeQueryNodeTypeUrlParam, nodeType.String())
		apstraUrl.RawQuery = params.Encode()
	}

	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
}

// PatchNode patches (only submitted fields are changed) the specified node
// using the contents of 'request', the server's response (whole node info
// without map wrapper?) is returned in 'response'
func (o *TwoStageL3ClosClient) PatchNode(ctx context.Context, node ObjectId, request interface{}, response interface{}) error {
	return o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPatch,
		urlStr:      fmt.Sprintf(apiUrlBlueprintNodeById, o.blueprintId, node),
		apiInput:    request,
		apiResponse: response,
	})
}
