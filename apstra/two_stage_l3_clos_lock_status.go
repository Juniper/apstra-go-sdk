package apstra

import (
	"context"
	"fmt"
	"net/http"
)

const (
	apiUrlBlueprintLockStatus = apiUrlBlueprintById + apiUrlPathDelim + "lock-status"
)

type LockStatus int
type lockStatus string

const (
	LockStatusUnlocked = LockStatus(iota)
	LockStatusLockedByRestrictedUser
	LockStatusLockedByAdmin
	LockStatusLockedByDeletedUser
	LockStatusUnknown = "unknown lock status %s"

	lockStatusUnlocked               = lockStatus("unlocked")
	lockStatusLockedByRestrictedUser = lockStatus("locked_by_restricted_user")
	lockStatusLockedByAdmin          = lockStatus("locked_by_admin")
	lockStatusLockedByDeletedUser    = lockStatus("locked_by_deleted_user")
	lockStatusUnknown                = "unknown lock status %d"
)

func (o LockStatus) String() string {
	switch o {
	case LockStatusUnlocked:
		return string(lockStatusUnlocked)
	case LockStatusLockedByRestrictedUser:
		return string(lockStatusLockedByRestrictedUser)
	case LockStatusLockedByAdmin:
		return string(lockStatusLockedByAdmin)
	case LockStatusLockedByDeletedUser:
		return string(lockStatusLockedByDeletedUser)
	default:
		return fmt.Sprintf(lockStatusUnknown, o)
	}
}

func (o LockStatus) int() int {
	return int(o)
}

func (o lockStatus) parse() (int, error) {
	switch o {
	case lockStatusUnlocked:
		return int(LockStatusUnlocked), nil
	case lockStatusLockedByRestrictedUser:
		return int(LockStatusLockedByRestrictedUser), nil
	case lockStatusLockedByAdmin:
		return int(LockStatusLockedByAdmin), nil
	case lockStatusLockedByDeletedUser:
		return int(LockStatusLockedByDeletedUser), nil
	default:
		return 0, fmt.Errorf(LockStatusUnknown, o)
	}
}

func (o lockStatus) string() string {
	return string(o)
}

type rawLockInfo struct {
	UserName         string     `json:"username"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	UserId           ObjectId   `json:"user_id"`
	PossibleOverride bool       `json:"possible_override"`
	LockStatus       lockStatus `json:"lock_status"`
}

func (o *rawLockInfo) polish() (*LockInfo, error) {
	ls, err := o.LockStatus.parse()
	if err != nil {
		return nil, err
	}
	return &LockInfo{
		UserName:         o.UserName,
		FirstName:        o.FirstName,
		LastName:         o.LastName,
		UserId:           o.UserId,
		PossibleOverride: o.PossibleOverride,
		LockStatus:       LockStatus(ls),
	}, nil
}

type LockInfo struct {
	UserName         string
	FirstName        string
	LastName         string
	UserId           ObjectId
	PossibleOverride bool
	LockStatus       LockStatus
}

func (o *LockInfo) String() string {
	return fmt.Sprintf("Lock status %q by %q (%s) override possible: %t",
		o.LockStatus, o.UserName, o.UserId, o.PossibleOverride)
}

func (o *TwoStageL3ClosClient) getLockInfo(ctx context.Context) (*rawLockInfo, error) {
	response := &rawLockInfo{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintLockStatus, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}
