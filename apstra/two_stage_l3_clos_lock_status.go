// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlBlueprintLockStatus = apiUrlBlueprintById + apiUrlPathDelim + "lock-status"
)

type LockInfo struct {
	LockStatus       enum.LockStatus `json:"lock_status"`
	LockType         *enum.LockType  `json:"lock_type"`
	UserName         *string         `json:"username"`
	FirstName        *string         `json:"first_name"`
	LastName         *string         `json:"last_name"`
	UserId           *ObjectId       `json:"user_id"`
	PossibleOverride *bool           `json:"possible_override"`
}

func (o *LockInfo) String() string {
	result := fmt.Sprintf("Lock status %q", o.LockStatus)
	if o.UserName != nil && o.UserId != nil {
		result += fmt.Sprintf(" by %q (%s)", *o.UserName, *o.UserId)
	}
	if o.PossibleOverride != nil {
		result += fmt.Sprintf(" override possible: %t", *o.PossibleOverride)
	}
	return result
}

// GetLockInfo returns *LockInfo describing the current state of the blueprint lock
func (o *TwoStageL3ClosClient) GetLockInfo(ctx context.Context) (*LockInfo, error) {
	var response LockInfo
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintLockStatus, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return &response, nil
}
