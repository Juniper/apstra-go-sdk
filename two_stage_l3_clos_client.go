package goapstra

import "context"

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
