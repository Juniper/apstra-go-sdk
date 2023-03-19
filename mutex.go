package goapstra

import (
	"context"
)

var _ error = MutexErr{}

type MutexErr struct {
	LockInfo *LockInfo
	Mutex    Mutex
	err      error
}

func (o MutexErr) Error() string {
	return o.err.Error()
}

type Mutex interface {
	GetMessage() string
	SetMessage(string) error
	BlueprintID() ObjectId
	Lock(context.Context) error
	TryLock(context.Context) error
	Unlock(context.Context) error
}
