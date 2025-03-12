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

//type (
//	LockStatus int
//	lockStatus string
//)
//
//const (
//	LockStatusUnlocked = LockStatus(iota)
//	LockStatusLocked
//	LockStatusLockedByRestrictedUser
//	LockStatusLockedByAdmin
//	LockStatusLockedByDeletedUser
//	LockStatusUnknown = "unknown lock status %s"
//
//	lockStatusUnlocked               = lockStatus("unlocked")
//	lockStatusLocked                 = lockStatus("locked")
//	lockStatusLockedByRestrictedUser = lockStatus("locked_by_restricted_user")
//	lockStatusLockedByAdmin          = lockStatus("locked_by_admin")
//	lockStatusLockedByDeletedUser    = lockStatus("locked_by_deleted_user")
//	lockStatusUnknown                = "unknown lock status %d"
//)
//
//func (o LockStatus) String() string {
//	switch o {
//	case LockStatusUnlocked:
//		return string(lockStatusUnlocked)
//	case LockStatusLocked:
//		return string(lockStatusLocked)
//	case LockStatusLockedByRestrictedUser:
//		return string(lockStatusLockedByRestrictedUser)
//	case LockStatusLockedByAdmin:
//		return string(lockStatusLockedByAdmin)
//	case LockStatusLockedByDeletedUser:
//		return string(lockStatusLockedByDeletedUser)
//	default:
//		return fmt.Sprintf(lockStatusUnknown, o)
//	}
//}
//
//func (o LockStatus) int() int {
//	return int(o)
//}
//
//func (o lockStatus) parse() (int, error) {
//	switch o {
//	case lockStatusUnlocked:
//		return int(LockStatusUnlocked), nil
//	case lockStatusLocked:
//		return int(LockStatusLocked), nil
//	case lockStatusLockedByRestrictedUser:
//		return int(LockStatusLockedByRestrictedUser), nil
//	case lockStatusLockedByAdmin:
//		return int(LockStatusLockedByAdmin), nil
//	case lockStatusLockedByDeletedUser:
//		return int(LockStatusLockedByDeletedUser), nil
//	default:
//		return 0, fmt.Errorf(LockStatusUnknown, o)
//	}
//}
//
//func (o lockStatus) string() string {
//	return string(o)
//}

//type rawLockInfo struct {
//	UserName         string     `json:"username"`
//	FirstName        string     `json:"first_name"`
//	LastName         string     `json:"last_name"`
//	UserId           ObjectId   `json:"user_id"`
//	PossibleOverride bool       `json:"possible_override"`
//	LockStatus       lockStatus `json:"lock_status"`
//}
//
//func (o *rawLockInfo) polish() (*LockInfo, error) {
//	ls, err := o.LockStatus.parse()
//	if err != nil {
//		return nil, err
//	}
//	return &LockInfo{
//		UserName:         o.UserName,
//		FirstName:        o.FirstName,
//		LastName:         o.LastName,
//		UserId:           o.UserId,
//		PossibleOverride: o.PossibleOverride,
//		LockStatus:       LockStatus(ls),
//	}, nil
//}

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
