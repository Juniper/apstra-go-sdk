// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

type FreeformClient struct {
	client      *Client
	blueprintId ObjectId
}

// Id returns the ID of the Freeform Blueprint associated with this client.
func (o *FreeformClient) Id() ObjectId {
	return o.blueprintId
}

// Client returns the Client within this freeform client.
func (o *FreeformClient) Client() *Client {
	return o.client
}
