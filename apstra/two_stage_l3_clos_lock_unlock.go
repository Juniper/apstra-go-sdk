// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

func decorateTwoStageL3ClosLockError(err TalkToApstraErr) error {
	switch {
	case strings.HasSuffix(err.Msg, " already locked"):
		return ClientErr{errType: ErrAlreadyLocked, err: err}
	case strings.Contains(err.Msg, " already locked by user: "):
		var errs apiErrors
		if e := json.Unmarshal([]byte(err.Msg), &errs); e != nil {
			return ClientErr{errType: ErrAlreadyLocked, err: err}
		}
		if len(errs.Errors) != 1 {
			return ClientErr{errType: ErrAlreadyLocked, err: err}
		}

		s := apiUrlLockBlueprintLockedByUserRegex.FindStringSubmatch(errs.Errors[0])
		if len(s) != 6 {
			return ClientErr{errType: ErrAlreadyLocked, err: err}
		}
		return ClientErr{
			errType: ErrAlreadyLocked,
			err:     err,
			detail:  ErrAlreadyLockedDetail{UserId: toPtr(ObjectId(s[2]))},
		}
	}

	return err
}

type ErrCannotUnlockDetail struct {
	NotEnoughPermission *bool
}

func decorateTwoStageL3ClosUnlockError(err TalkToApstraErr) error {
	switch {
	case strings.Contains(err.Msg, " does not have enough permissions to unlock blueprint "):
		return ClientErr{
			errType: ErrCannotUnlock,
			err:     errors.New(err.Msg),
			detail:  ErrCannotUnlockDetail{NotEnoughPermission: toPtr(true)},
		}
	}

	return err
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
