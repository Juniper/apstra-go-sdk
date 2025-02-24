// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlFfDp       = apiUrlBlueprintById + apiUrlPathDelim + "device-profiles"
	apiUrlFfDpImport = apiUrlFfDp + apiUrlPathDelim + "import"
	apiUrlFfDpById   = apiUrlFfDp + apiUrlPathDelim + "%s"
)

func (o *FreeformClient) ImportDeviceProfile(ctx context.Context, id ObjectId) (ObjectId, error) {
	var response struct {
		Ids []ObjectId `json:"ids"`
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodPost,
		urlStr: fmt.Sprintf(apiUrlFfDpImport, o.blueprintId),
		apiInput: struct {
			DeviceProfileIds []ObjectId `json:"device_profile_ids"`
		}{
			DeviceProfileIds: []ObjectId{id},
		},
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	if len(response.Ids) != 1 {
		return "", fmt.Errorf("expected one ObjectId got %d", len(response.Ids))
	}

	return response.Ids[0], nil
}

func (o *FreeformClient) DeleteDeviceProfile(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfDpById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
