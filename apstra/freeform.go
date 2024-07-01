package apstra

type FreeformClient struct {
	client      *Client
	blueprintId ObjectId
}

// Id returns the ID of the Freeform Blueprint associated with this client.
func (o *FreeformClient) Id() ObjectId {
	return o.blueprintId
}
