// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

const (
	apiUrlLockBlueprint         = "/api/blueprints/%s/lock-blueprint"
	apiUrlLockBlueprintRegexStr = "/api/blueprints/[^/]+/lock-blueprint"

	apiUrlUnlockBlueprint         = "/api/blueprints/%s/unlock-blueprint"
	apiUrlUnlockBlueprintRegexStr = "/api/blueprints/[^/]+/unlock-blueprint"

	apiUrlLockBlueprintLockedByUserRegexStr = "^Blueprint ([^ ]+) already locked by user: user_id: ([^ ]+), username: ([^ ]+), first_name: ([^ ]+), last_name: ([^ ]+)$"
)

var (
	apiUrlLockBlueprintRegex   = regexp.MustCompile(apiUrlLockBlueprintRegexStr)
	apiUrlUnlockBlueprintRegex = regexp.MustCompile(apiUrlUnlockBlueprintRegexStr)

	apiUrlLockBlueprintLockedByUserRegex = regexp.MustCompile(apiUrlLockBlueprintLockedByUserRegexStr)
)

type ErrAlreadyLockedDetail struct {
	UserId *ObjectId
}

type ErrCannotUnlockDetail struct {
	NotEnoughPermission *bool
}

func (o *TwoStageL3ClosClient) Lock(ctx context.Context) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPut,
		urlStr:         fmt.Sprintf(apiUrlLockBlueprint, o.Id()),
		unsynchronized: true,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *TwoStageL3ClosClient) Unlock(ctx context.Context) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:         http.MethodPut,
		urlStr:         fmt.Sprintf(apiUrlUnlockBlueprint, o.Id()),
		unsynchronized: true,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
