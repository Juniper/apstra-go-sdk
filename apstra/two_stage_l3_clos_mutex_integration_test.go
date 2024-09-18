//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestLockUnlockBlueprintMutex is a simple test of Lock() followed by Unlock()
func TestLockUnlockBlueprintMutex(t *testing.T) {
	ctx := context.Background()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintA(ctx, t, client.client)

			err = bp.Mutex.SetMessage("locked by apstra test")
			require.NoError(t, err)

			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.Mutex.Lock(ctx)
			require.NoError(t, err)

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.Mutex.Unlock(ctx)
			require.NoError(t, err)
		})
	}
}

// TestLockLockUnlockBlueprintMutex tests locking an already-locked Mutex,
// verifies that we hit context expiration
func TestLockLockUnlockBlueprintMutex(t *testing.T) {
	ctx := context.Background()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpA := testBlueprintA(ctx, t, client.client)

			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpB, err := client.client.NewTwoStageL3ClosClient(ctx, bpA.Id())
			require.NoError(t, err)

			err = bpA.Mutex.SetMessage("locked by test client A")
			require.NoError(t, err)

			err = bpB.Mutex.SetMessage("locked by test client B")
			require.NoError(t, err)

			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Lock(ctx)
			require.NoError(t, err)

			start := time.Now()
			timeout := 2 * time.Second
			timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.Lock(timeoutCtx)
			require.Error(t, err)
			require.Contains(t, err.Error(), "context deadline exceeded")

			elapsed := time.Since(start)
			require.Greaterf(t, elapsed, timeout, "we should have waited %s for lock, but only %s has elapsed", timeout, elapsed)

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Unlock(ctx)
			require.NoError(t, err)
		})
	}
}

// TestLockBlockUnlockLockUnlockBlueprintMutex tests locking an already-locked
// blueprint, then clearing the lock to un-block the 2nd lock attempt. It checks
// to ensure that timing is works out as expected.
func TestLockBlockUnlockLockUnlockBlueprintMutex(t *testing.T) {
	ctx := context.Background()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpA := testBlueprintA(ctx, t, client.client)

			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpB, err := client.client.NewTwoStageL3ClosClient(ctx, bpA.Id())
			require.NoError(t, err)

			err = bpA.Mutex.SetMessage("locked by client A")
			require.NoError(t, err)

			err = bpB.Mutex.SetMessage("locked by client B")
			require.NoError(t, err)

			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Lock(ctx)
			require.NoError(t, err)

			timeout := 10 * time.Second
			timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			go func() {
				delay := 2 * time.Second
				log.Printf("Unlock scheduled for %s", time.Now().Add(delay))
				time.Sleep(delay)
				err = bpA.Mutex.Unlock(ctx)
				require.NoError(t, err, "error unlocking blueprint")
			}()

			start := time.Now()
			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.Lock(timeoutCtx)
			require.NoError(t, err)

			blocked := time.Since(start)
			require.Greater(t, timeout, blocked, "unlock took longer than expected")
			log.Printf("Blocked for %s waiting for lock to clear", blocked)

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.Unlock(ctx)
			require.NoError(t, err)
		})
	}
}

// TestLockTryLockBlueprintMutex verifies expected outcome when TryLock() is
// attempted against a previously Lock()ed blueprint.
func TestLockTryLockBlueprintMutex(t *testing.T) {
	ctx := context.Background()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()
			bpA := testBlueprintA(ctx, t, client.client)

			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpB, err := client.client.NewTwoStageL3ClosClient(ctx, bpA.Id())
			require.NoError(t, err)

			log.Print("testing SetMessage()")
			err = bpA.Mutex.SetMessage("locked by client A")
			require.NoError(t, err)

			log.Print("testing SetMessage()")
			err = bpB.Mutex.SetMessage("locked by client B")
			require.NoError(t, err)

			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Lock(ctx)
			require.NoError(t, err)

			var mutexErr MutexErr
			log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.TryLock(ctx)
			require.Error(t, err, "TryLock should have returned an error")
			require.ErrorAs(t, err, &mutexErr, "TryLock should have returned a MutexErr")
			require.Nil(t, mutexErr.LockInfo, "mutexErr's LockInfo should be nil")
			require.NotNil(t, mutexErr.Mutex, "mutexErr's Mutex should not be nil")
			require.Equal(t, mutexErr.Mutex.GetMessage(), bpA.Mutex.GetMessage(), "reason and original lock messages do not match")

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Unlock(ctx)
			require.NoError(t, err)
		})
	}
}

// TestTryLockTryLockBlueprintMutex verifies expected outcome when TryLock() is
// attempted against a previously TryLock()ed blueprint.
func TestTryLockTryLockBlueprintMutex(t *testing.T) {
	ctx := context.Background()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpA := testBlueprintA(ctx, t, client.client)

			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpB, err := client.client.NewTwoStageL3ClosClient(ctx, bpA.Id())
			require.NoError(t, err)

			log.Print("testing SetMessage()")
			err = bpA.Mutex.SetMessage("locked by client A")
			require.NoError(t, err)

			log.Print("testing SetMessage()")
			err = bpB.Mutex.SetMessage("locked by client B")
			require.NoError(t, err)

			log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.TryLock(ctx)
			require.NoError(t, err)

			var mutexErr MutexErr
			log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.TryLock(ctx)
			require.Error(t, err)
			require.ErrorAs(t, err, &mutexErr, "TryLock should have returned a MutexErr")
			require.NotNil(t, mutexErr.Mutex, "TryLock returned a MutexErr with nil Mutex")
			require.NotNil(t, mutexErr.LockInfo, "TryLock returned a MutexEerr with non-nil LockInfo")
			require.Equal(t, mutexErr.Mutex.GetMessage(), bpA.Mutex.GetMessage(), "blocking mutex and original messages do not match")

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Unlock(ctx)
			require.NotNil(t, err)
		})
	}
}

// TestLockBlockTrylockUnlLockTrylockBlueprintMutex ensures that TryLock() fails
// against a locked blueprint, then succeeds after the lock is released.
func TestLockUnlLockTrylockBlueprintMutex(t *testing.T) {
	ctx := context.Background()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpA := testBlueprintA(ctx, t, client.client)

			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpB, err := client.client.NewTwoStageL3ClosClient(ctx, bpA.Id())
			require.NoError(t, err)

			log.Print("testing SetMessage()")
			err = bpA.Mutex.SetMessage("locked by client A")
			require.NoError(t, err)

			log.Print("testing SetMessage()")
			err = bpB.Mutex.SetMessage("locked by client B")
			require.NoError(t, err)

			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Lock(ctx)
			require.NoError(t, err)

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Unlock(ctx)
			require.NoError(t, err)

			log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.TryLock(ctx)
			require.NoError(t, err)

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.Unlock(ctx)
			require.NoError(t, err)
		})
	}
}

// TestReadOnlyBlueprintMutex gets a read-only mutex (reason) by calling
// TryLock() against a locked blueprint. It then calls various "write" methods
// on that mutex to ensure they fail with the expected error type.
func TestReadOnlyBlueprintMutex(t *testing.T) {
	ctx := context.Background()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			bpA := testBlueprintA(ctx, t, client.client)

			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpB, err := client.client.NewTwoStageL3ClosClient(ctx, bpA.Id())
			require.NoError(t, err)

			msgA := "locked by client A"
			log.Print("testing SetMessage()")
			err = bpA.Mutex.SetMessage(msgA)
			require.NoError(t, err)
			require.Equal(t, msgA, bpA.Mutex.GetMessage())

			msgB := "locked by client B"
			log.Print("testing SetMessage()")
			err = bpB.Mutex.SetMessage(msgB)
			require.NoError(t, err)
			require.Equal(t, msgB, bpB.Mutex.GetMessage())

			log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Lock(ctx)
			require.NoError(t, err)

			var mutexErr MutexErr

			log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpB.Mutex.TryLock(ctx)
			require.Error(t, err, "TryLock should have returned an error")
			require.ErrorAs(t, err, &mutexErr)
			require.NotNil(t, mutexErr.Mutex, "TryLock should have returned a *Mutex")

			roMutex := mutexErr.Mutex

			log.Printf("testing Lock() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = roMutex.Lock(ctx)
			require.Errorf(t, err, "locking a read-only Mutex should return an error")
			require.Contains(t, err.Error(), "read", "locking a read-only Mutex complain about it being read-only")

			log.Printf("testing TryLock() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = roMutex.TryLock(ctx)
			require.Errorf(t, err, "try-locking a read-only Mutex should return an error")
			require.Contains(t, err.Error(), "read", "try-locking a read-only Mutex complain about it being read-only")

			log.Printf("testing Unlock() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = roMutex.Unlock(ctx)
			require.Errorf(t, err, "unlocking a read-only Mutex should return an error")
			require.Contains(t, err.Error(), "read", "unlocking a read-only Mutex complain about it being read-only")

			log.Printf("testing SetMessage() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = roMutex.SetMessage("fail please")
			require.Error(t, err, "setting message of a read-only Mutex should return an error")
			require.Contains(t, err.Error(), "read", "setting message of a read-only Mutex complain about it being read-only")

			log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpA.Mutex.Unlock(ctx)
			require.NoError(t, err)
		})
	}
}
