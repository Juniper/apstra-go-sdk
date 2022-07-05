package goapstra

import "context"

// GetResourceAllocation takes a *ResourceGroupAllocation as input for the
// ResourceGroupAllocation.Type and ResourceGroupAllocation.Name fields (the
// ResourceGroupAllocation.PoolIds is ignored). It returns a fully populated
// *ResourceGroupAllocation with all fields populated based on the Apstra API
// response.
func (o *Blueprint) GetResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) (*ResourceGroupAllocation, error) {
	return o.getResourceAllocation(ctx, in)

}

// SetResourceAllocation sets the supplied resource allocation, overwriting any
// allocations with the supplied info.
func (o *Blueprint) SetResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) error {
	return o.setResourceAllocation(ctx, in)
}
