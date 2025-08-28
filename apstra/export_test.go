// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "context"

// GetFeatures is in the export_test file because this private function needs to be exposed only for test code
func (o *Client) GetFeatures(ctx context.Context) error {
	return o.getFeatures(ctx)
}
