package apstra

import "context"

// GetNodes fetches the node of the specified type, unpacks the API response
// into 'response'
func (o *FreeformClient) GetNodes(ctx context.Context, nodeType NodeType, response interface{}) error {
	return o.client.getNodes(ctx, o.blueprintId, nodeType, response)
}
