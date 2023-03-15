package goapstra

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	tagNameLenMax = 64                    // apstra 4.1.0 limit
	lockTagName   = "blueprint %s locked" // BP GUID is 36 char, this should fit
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
		return ApstraClientErr{
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

// ID returns the lock ID
func (o *TwoStageL3ClosMutex) ID() ObjectId {
	return o.tagId
}

// Lock blocks until the mutex is locked or the context expires/is canceled.
func (o *TwoStageL3ClosMutex) Lock(ctx context.Context) error {
	if o.tagId != "" {
		return ApstraClientErr{
			errType: ErrReadOnly,
			err:     fmt.Errorf("attempt to lock previously locked mutex - previous lock ID %q", o.tagId),
		}
	}
	if o.readOnly {
		return ApstraClientErr{
			errType: ErrReadOnly,
			err:     errors.New("attempt to lock read-only mutex"),
		}
	}

	lockName := fmt.Sprintf(lockTagName, o.client.blueprintId)
	if len(lockName) > tagNameLenMax {
		return fmt.Errorf("lock name %q exceeds limit (max %d characters)", lockName, tagNameLenMax)
	}

	// set initial LockStatus to bogus value b/c desired state is "0"
	li := &LockInfo{LockStatus: -1}

	// refuse to lock this mutex while the apstra blueprint lock is asserted
	var err error
	for li.LockStatus != LockStatusUnlocked {
		// mind the timeout
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for lock status %q - %w",
				LockStatusUnlocked.String(), ctx.Err())
		default:
		}

		li, err = o.client.GetLockInfo(ctx)
		if err != nil {
			return err
		}
	}

	// loop until we acquire the lock or the context deadline (set by caller) expires.
	for {
		// mind the timeout
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while trying to establish lock - %w", ctx.Err())
		default:
		}

		tagId, err := o.client.client.createTag(ctx, &DesignTagRequest{
			Label:       lockName,
			Description: o.message,
		})
		if err != nil {
			var ace ApstraClientErr
			if errors.As(err, &ace) && ace.errType == ErrExists {
				time.Sleep(clientPollingIntervalMs * time.Millisecond)
				continue
			}
			return err
		}
		o.tagId = tagId
		return nil
	}
}

// Unlock releases the mutex
func (o *TwoStageL3ClosMutex) Unlock(ctx context.Context) error {
	if o.readOnly {
		return ApstraClientErr{
			errType: ErrReadOnly,
			err:     errors.New("attempt to unlock read-only mutex"),
		}
	}
	err := o.client.client.deleteTag(ctx, o.tagId)
	if err != nil {
		return err
	}
	o.tagId = ""
	return nil
}

// TryLock attempts to lock the mutex without blocking until success. The
// returned boolean indicates the success of the lock attempt. On success,
// the returned values are <true, nil, nil>.
// When the attempt to lock is blocked by a different TwoStageL3ClosMutex,
// a read-only copy of the offender is returned as *TwoStageL3ClosMutex for
// inspection by the caller.
// When the attempt to lock is blocked by Apstra's user-based blueprint lock
// the return values are <false, nil, nil>
// Setting ignoreApstraLock causes the check of the user-based blueprint
// lock to be skipped.
func (o *TwoStageL3ClosMutex) TryLock(ctx context.Context, ignoreApstraLock bool) (bool, *TwoStageL3ClosMutex, error) {
	if o.readOnly {
		return false, nil, ApstraClientErr{
			errType: ErrReadOnly,
			err:     errors.New("attempt to unlock read-only mutex"),
		}
	}

	lockName := fmt.Sprintf(lockTagName, o.client.blueprintId)
	if len(lockName) > tagNameLenMax {
		return false, nil, fmt.Errorf("lock name %q exceeds limit (max %d characters)", lockName, tagNameLenMax)
	}

	if o.tagId != "" {
		err := o.client.client.UpdateTag(ctx, o.tagId, &DesignTagRequest{
			Label:       lockName,
			Description: o.message,
		})
		if err != nil {
			return false, nil, fmt.Errorf("error updating tag/lock %q - %w", o.tagId, err)
		}
		return true, nil, nil
	}

	if !ignoreApstraLock {
		// refuse to lock this mutex if the apstra blueprint lock exists
		li, err := o.client.GetLockInfo(ctx)
		if err != nil {
			return false, nil, err
		}
		if li.LockStatus != LockStatusUnlocked {
			return false, nil, nil
		}
	}

	tagId, err := o.client.client.createTag(ctx, &DesignTagRequest{
		Label:       lockName,
		Description: o.message,
	})
	if err != nil {
		var ace ApstraClientErr
		if errors.As(err, &ace) && ace.errType == ErrExists {
			blockingTag, err := o.client.client.GetTagByLabel(ctx, lockName)
			if err != nil {
				return false, nil, err
			}
			return false, &TwoStageL3ClosMutex{
				tagId:    blockingTag.Id,
				readOnly: true,
				message:  blockingTag.Data.Description,
			}, nil
		}
		return false, nil, err
	}
	o.tagId = tagId
	return true, nil, nil
}

// ClearUnsafely deletes the mutex regardless of whether we originally held it.
// This method should be used in exceptional circumstances only. NOT SAFE!
func (o *TwoStageL3ClosMutex) ClearUnsafely(ctx context.Context) error {
	lockName := fmt.Sprintf(lockTagName, o.client.blueprintId)
	if len(lockName) > tagNameLenMax {
		return fmt.Errorf("lock name %q exceeds limit (max %d characters)", lockName, tagNameLenMax)
	}

	tag, err := o.client.client.GetTagByLabel(ctx, lockName)
	if err != nil {
		var ace ApstraClientErr
		if errors.As(err, &ace) && ace.errType == ErrNotfound {
			return nil
		}
		return fmt.Errorf("error while fetching tag with name %q - %w", lockName, err)
	}

	err = o.client.client.DeleteTag(ctx, tag.Id)
	if err != nil {
		var ace ApstraClientErr
		if errors.As(err, &ace) && ace.errType == ErrNotfound {
			return nil
		}
		return fmt.Errorf("error while fetching tag with id %q - %w", tag.Id, err)
	}

	return nil
}
