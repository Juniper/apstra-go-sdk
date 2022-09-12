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
	client *TwoStageL3ClosClient
	tagId  ObjectId
}

func (o *TwoStageL3ClosMutex) Lock(ctx context.Context, message string) error {
	if o.tagId != "" {
		return fmt.Errorf("attempt to lock previously locked mutex - previous lock ID '%s'", o.tagId)
	}

	lockName := fmt.Sprintf(lockTagName, o.client.blueprintId)
	if len(lockName) > tagNameLenMax {
		return fmt.Errorf("lock name '%s' exceeds limit (max %d characters)", lockName, tagNameLenMax)
	}

	// set initial LockStatus to bogus value b/c desired state is "0"
	li := &LockInfo{LockStatus: -1}

	// refuse to lock this mutex while the actual blueprint is locked
	var err error
	for li.LockStatus != LockStatusUnlocked {
		li, err = o.client.GetLockInfo(ctx)
		if err != nil {
			return err
		}

		// mind the timeout
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for lock status '%s' - %w",
				LockStatusUnlocked.String(), ctx.Err())
		default:
		}
	}

	// loop until we acquire the lock or the context deadline (set by caller) expires.
	for {
		tagId, err := o.client.client.createTag(ctx, &DesignTag{
			Label:       TagLabel(lockName),
			Description: message,
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

func (o *TwoStageL3ClosMutex) Unlock(ctx context.Context) error {
	err := o.client.client.deleteTag(ctx, o.tagId)
	if err != nil {
		return err
	}
	o.tagId = ""
	return nil
}

func (o *TwoStageL3ClosMutex) Id() ObjectId {
	return o.tagId
}
