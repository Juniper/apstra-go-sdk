package goapstra

import (
	"context"
	"fmt"
)

const (
	blueprintTypeParam = "type"
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

// GetSecurityZoneByLabel fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o *TwoStageLThreeClosClient) GetSecurityZoneByLabel(ctx context.Context, label string) (*SecurityZone, error) {
	return o.getSecurityZoneByLabel(ctx, label)
}

// GetAllSecurityZones returns []SecurityZone representing all Security Zones /
// Routing Zones / VRFs on the system.
func (o *TwoStageLThreeClosClient) GetAllSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	return o.getAllSecurityZones(ctx)
}
