package goapstra

import "context"

// GetResourceAllocation takes a *ResourceGroupAllocation as input for the
// ResourceGroupAllocation.Type and ResourceGroupAllocation.Name fields (the
// ResourceGroupAllocation.PoolIds is ignored). It returns a fully populated
// *ResourceGroupAllocation with all fields populated based on the Apstra API
// response.
func (o *Client) GetResourceAllocation(ctx context.Context, blueprintId ObjectId, in *ResourceGroupAllocation) (*ResourceGroupAllocation, error) {
	return o.getResourceAllocation(ctx, blueprintId, in)

}

// SetResourceAllocation sets the supplied resource allocation, overwriting any
// allocations with the supplied info.
func (o *Client) SetResourceAllocation(ctx context.Context, blueprintId ObjectId, in *ResourceGroupAllocation) error {
	return o.setResourceAllocation(ctx, blueprintId, in)
}
