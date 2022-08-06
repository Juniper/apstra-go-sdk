package goapstra

import (
	"context"
	"fmt"
	"net"
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

type TwoStageLThreeClosClient struct {
	client      *Client
	blueprintId ObjectId
}

// GetResourceAllocation takes a *ResourceGroupAllocation as input for the
// ResourceGroupAllocation.Type and ResourceGroupAllocation.Name fields (the
// ResourceGroupAllocation.PoolIds is ignored). It returns a fully populated
// *ResourceGroupAllocation with all fields populated based on the Apstra API
// response.
func (o *TwoStageLThreeClosClient) GetResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) (*ResourceGroupAllocation, error) {
	return o.getResourceAllocation(ctx, in)
}

// SetResourceAllocation sets the supplied resource allocation, overwriting any
// prior allocations with the supplied info.
func (o *TwoStageLThreeClosClient) SetResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) error {
	return o.setResourceAllocation(ctx, in)
}

// GetInterfaceMapAssignments takes a *InterfaceMapAssignments as input for the
// ResourceGroupAllocation.Type and ResourceGroupAllocation.Name fields (the
// ResourceGroupAllocation.PoolIds is ignored). It returns a fully populated
// *ResourceGroupAllocation with all fields populated based on the Apstra API
// response.
func (o *TwoStageLThreeClosClient) GetInterfaceMapAssignments(ctx context.Context) (SystemIdToInterfaceMapAssignment, error) {
	return o.getInterfaceMapAssignments(ctx)
}

// SetInterfaceMapAssignments sets the supplied interface map assignments,
// overwriting any prior assignments with the supplied info. It returns
// the Blueprint config revision number.
func (o *TwoStageLThreeClosClient) SetInterfaceMapAssignments(ctx context.Context, assignments SystemIdToInterfaceMapAssignment) error {
	return o.setInterfaceMapAssignments(ctx, assignments)
}

// CreateSecurityZone creates an Apstra Routing Zone / Security Zone / VRF
func (o *TwoStageLThreeClosClient) CreateSecurityZone(ctx context.Context, cfg *CreateSecurityZoneCfg) (ObjectId, error) {
	response, err := o.createSecurityZone(ctx, cfg)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// DeleteSecurityZone deletes an Apstra Routing Zone / Security Zone / VRF
func (o *TwoStageLThreeClosClient) DeleteSecurityZone(ctx context.Context, zoneId ObjectId) error {
	return o.deleteSecurityZone(ctx, zoneId)
}

// GetSecurityZones returns all Apstra Routing Zones / Security Zones / VRFs
// associated with the specified blueprint
func (o *TwoStageLThreeClosClient) GetSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	return o.getAllSecurityZones(ctx)
}

// GetSecurityZone fetches the Security Zone / Routing Zone / VRF with the given
// zoneId.
func (o *TwoStageLThreeClosClient) GetSecurityZone(ctx context.Context, zoneId ObjectId) (*SecurityZone, error) {
	return o.getSecurityZone(ctx, zoneId)
}

// GetSecurityZoneByName fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o *TwoStageLThreeClosClient) GetSecurityZoneByName(ctx context.Context, label string) (*SecurityZone, error) {
	return o.getSecurityZoneByName(ctx, label)
}

// GetAllSecurityZones returns []SecurityZone representing all Security Zones /
// Routing Zones / VRFs on the system.
func (o *TwoStageLThreeClosClient) GetAllSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	return o.getAllSecurityZones(ctx)
}

// UpdateSecurityZone replaces the configuration of zone zoneId with the supplied CreateSecurityZoneCfg
func (o *TwoStageLThreeClosClient) UpdateSecurityZone(ctx context.Context, zoneId ObjectId, cfg *CreateSecurityZoneCfg) error {
	return o.updateSecurityZone(ctx, zoneId, cfg)
}

// GetAllPolicies returns []Policy representing all policies configured within the DC blueprint
func (o *TwoStageLThreeClosClient) GetAllPolicies(ctx context.Context) ([]Policy, error) {
	return o.getAllPolicies(ctx)
}

// GetPolicy returns *Policy representing policy 'id' within the DC blueprint
func (o *TwoStageLThreeClosClient) GetPolicy(ctx context.Context, id ObjectId) (*Policy, error) {
	return o.getPolicy(ctx, id)
}

// CreatePolicy creates a policy within the DC blueprint, returns its ID
func (o *TwoStageLThreeClosClient) CreatePolicy(ctx context.Context, policy *Policy) (ObjectId, error) {
	return o.createPolicy(ctx, policy)
}

// DeletePolicy deletes policy 'id' within the DC blueprint
func (o *TwoStageLThreeClosClient) DeletePolicy(ctx context.Context, id ObjectId) error {
	return o.deletePolicy(ctx, id)
}

// UpdatePolicy calls PUT to replace the configuration of policy 'id' within the DC blueprint
func (o *TwoStageLThreeClosClient) UpdatePolicy(ctx context.Context, id ObjectId, policy *Policy) error {
	return o.updatePolicy(ctx, id, policy)
}

// AddPolicyRule adds a policy rule at 'position' (bumping all other rules
// down). Position 0 makes the new policy first on the list, 1 makes it second
// on the list, etc... Use -1 for last on the list. The returned ObjectId
// represents the new rule
func (o *TwoStageLThreeClosClient) AddPolicyRule(ctx context.Context, rule *PolicyRule, position int, policyId ObjectId) (ObjectId, error) {
	return o.addPolicyRule(ctx, rule, position, policyId)
}

// DeletePolicyRuleById deletes the given rule. If the rule doesn't exist, an
// ApstraClientErr with ErrNotFound is returned.
func (o *TwoStageLThreeClosClient) DeletePolicyRuleById(ctx context.Context, policyId ObjectId, ruleId ObjectId) error {
	return o.deletePolicyRuleById(ctx, policyId, ruleId)
}

// ListAllVirtualNetworkIds returns []ObjectId representing virtual networks configured in the blueprint
func (o *TwoStageLThreeClosClient) ListAllVirtualNetworkIds(ctx context.Context, bpType BlueprintType) ([]ObjectId, error) {
	return o.listAllVirtualNetworkIds(ctx, bpType)
}

// GetVirtualNetwork returns *VirtualNetwork representing the given vnId within the blueprint type
func (o *TwoStageLThreeClosClient) GetVirtualNetwork(ctx context.Context, vnId ObjectId, bpType BlueprintType) (*VirtualNetwork, error) {
	return o.getVirtualNetwork(ctx, vnId, bpType)
}

// GetVirtualNetworkBySubnet returns *VirtualNetwork representing the given desiredNet within the blueprint type
func (o *TwoStageLThreeClosClient) GetVirtualNetworkBySubnet(ctx context.Context, desiredNet *net.IPNet, vrf ObjectId, bpType BlueprintType) (*VirtualNetwork, error) {
	return o.getVirtualNetworkBySubnet(ctx, desiredNet, vrf, bpType)
}
