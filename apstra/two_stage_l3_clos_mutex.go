// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	tagNameLenMax    = 64                    // apstra 4.1.0 limit
	lockTagName      = "blueprint %s locked" // BP GUID is 36 char, this should fit
	lockPollInterval = 500 * time.Millisecond
)

type TwoStageL3ClosMutex struct {
	client   *TwoStageL3ClosClient
	tagId    ObjectId
	readOnly bool
	message  string
}

// GetMessage returns the message embedded in the mutex
func (o *TwoStageL3ClosMutex) GetMessage() string {
	return o.message
}

// SetMessage sets the lock message embedded in the mutex
func (o *TwoStageL3ClosMutex) SetMessage(msg string) error {
	if o.readOnly {
		return ClientErr{
			errType: ErrReadOnly,
			err:     errors.New("attempt to set message of a read-only mutex"),
		}
	}
	if o.tagId != "" {
		return errors.New("attempt to set message of a locked mutex")
	}
	o.message = msg
	return nil
}

// BlueprintID returns the Blueprint ID
func (o *TwoStageL3ClosMutex) BlueprintID() ObjectId {
	return o.client.blueprintId
}

// Lock attempts to assert the blueprint mutex, repeatedly trying until the
// context.Context expires or it encounters an error.
func (o *TwoStageL3ClosMutex) Lock(ctx context.Context) error {
	return o.lock(ctx, false)
}

// TryLock attempts to assert the blueprint mutex without blocking.
func (o *TwoStageL3ClosMutex) TryLock(ctx context.Context) error {
	return o.lock(ctx, true)
}

// lock's behavior is controlled by the nonBlocking boolean. When called with
// nonBlocking == false, it will block until it asserts the mutex/tag, or an
// error is encountered. When called with nonBlocking == true, it will return
// a MutexErr populated with either a *LockInfo, indicating an Apstra blueprint
// lock was in place, or a *Mutex indicating somebody else has asserted the
// tag/mutex. In either case, the caller can inspect the MutexErr to learn
// exactly what went wrong.
func (o *TwoStageL3ClosMutex) lock(ctx context.Context, nonBlocking bool) error {
	if o.readOnly {
		return errors.New("attempt to lock read-only mutex")
	}

	if o.tagId != "" {
		return fmt.Errorf("attempt to lock previously locked mutex - previous lock ID %q", o.tagId)
	}

	lockName := fmt.Sprintf(lockTagName, o.client.blueprintId)
	if len(lockName) > tagNameLenMax {
		return fmt.Errorf("lock name %q exceeds limit (max %d characters)", lockName, tagNameLenMax)
	}

	li := new(LockInfo)

	var err error
	tickerA := immediateTicker(lockPollInterval)
	defer tickerA.Stop()

	for li.LockStatus != enum.LockStatusUnlocked {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for lock status %q - %w",
				enum.LockStatusUnlocked, ctx.Err())
		case <-tickerA.C:
		}

		li, err = o.client.GetLockInfo(ctx)
		if err != nil {
			return err
		}

		// Pass when locked by our own ID.
		if li.LockStatus == enum.LockStatusLocked && li.UserId != nil && *li.UserId == o.client.client.id {
			break
		}

		if nonBlocking && li.LockStatus != enum.LockStatusUnlocked {
			return MutexErr{
				LockInfo: li,
				err:      fmt.Errorf("blueprint %q: %s", o.client.blueprintId, li.String()),
			}
		}
	}

	// loop until we acquire the lock or the context deadline (set by caller) expires.
	tickerB := immediateTicker(lockPollInterval)
	defer tickerB.Stop()
	var ace ClientErr
	var tagID ObjectId
	for tagID == "" {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while trying to establish lock - %w", ctx.Err())
		case <-tickerB.C:
		}

		tagID, err = o.client.client.createTag(ctx, &DesignTagRequest{
			Label:       lockName,
			Description: o.message,
		})
		if err != nil {
			if errors.As(err, &ace) && ace.errType == ErrExists {
				// mutex already exists
				if !nonBlocking {
					// caller specified blocking behavior; nothing to do but try again
					continue
				}

				// retrieve the offending tag so we can inform the caller about it
				tag, err := o.client.client.getTagByLabel(ctx, lockName)
				if err != nil {
					if errors.As(err, &ace) {
						// offending tag deleted in the last few milliseconds? Try again.
						continue
					}
					// error retrieving the offending tag. blow up in the caller's face.
					return err
				}
				tagURL := fmt.Sprintf(apiUrlDesignTagById, tagID)
				return MutexErr{
					err: fmt.Errorf("unable to lock blueprint mutex due to: %q", tagURL),
					Mutex: &TwoStageL3ClosMutex{
						client:   o.client,
						tagId:    tag.Id,
						readOnly: true,
						message:  tag.Description,
					},
				}
			}
			// some other tag creation error
			return err
		}
	}

	o.tagId = tagID
	return nil
}

// Unlock releases the mutex
func (o *TwoStageL3ClosMutex) Unlock(ctx context.Context) error {
	if o.readOnly {
		return ClientErr{
			errType: ErrReadOnly,
			err:     errors.New("attempt to unlock read-only mutex"),
		}
	}

	err := o.client.client.deleteTag(ctx, o.tagId)
	if err != nil {
		var ace ClientErr
		if !errors.As(err, &ace) || ace.Type() != ErrNotfound {
			return err
		}
	}

	o.tagId = ""
	return nil
}
