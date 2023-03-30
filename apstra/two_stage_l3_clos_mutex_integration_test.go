//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"log"
	"strings"
	"testing"
	"time"
)

// TestLockUnlockBlueprintMutex is a simple test of Lock() followed by Unlock()
func TestLockUnlockBlueprintMutex(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for clientName, client := range clients {
		log.Printf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bp, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		err = bp.Mutex.SetMessage("locked by apstra test")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bp.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bp.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestLockLockUnlockBlueprintMutex tests locking an already-locked Mutex,
// verifies that we hit context expiration
func TestLockLockUnlockBlueprintMutex(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for clientName, client := range clients {
		log.Printf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpA, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpB, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		err = bpA.Mutex.SetMessage("locked by test client A")
		if err != nil {
			t.Fatal(err)
		}

		err = bpB.Mutex.SetMessage("locked by test client B")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		start := time.Now()
		timeout := 2 * time.Second
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.Lock(timeoutCtx)
		if err != nil && !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Fatal(err)
		}

		elapsed := time.Since(start)
		if elapsed < timeout {
			t.Fatalf("we should have waited %s for lock, but only %s has elapsed", timeout.String(), elapsed.String())
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestLockBlockUnlockLockUnlockBlueprintMutex tests locking an already-locked
// blueprint, then clearing the lock to un-block the 2nd lock attempt. It checks
// to ensure that timing is works out as expected.
func TestLockBlockUnlockLockUnlockBlueprintMutex(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for clientName, client := range clients {
		log.Printf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpA, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpB, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		err = bpA.Mutex.SetMessage("locked by client A")
		if err != nil {
			t.Fatal(err)
		}

		err = bpB.Mutex.SetMessage("locked by client B")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		timeout := 10 * time.Second
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		go func() {
			delay := 2 * time.Second
			log.Printf("Unlock scheduled for %s", time.Now().Add(delay))
			time.Sleep(delay)
			bpA.Mutex.Unlock(context.Background())
		}()

		start := time.Now()
		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.Lock(timeoutCtx)
		if err != nil {
			t.Fatal(err)
		}
		blocked := time.Since(start)
		if blocked > timeout {
			t.Fatal("unlock took longer than expected")
		}
		log.Printf("Blocked for %s waiting for lock to clear", blocked)

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestLockTryLockBlueprintMutex verifies expected outcome when TryLock() is
// attempted against a previously Lock()ed blueprint.
func TestLockTryLockBlueprintMutex(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for clientName, client := range clients {
		log.Printf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpA, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpB, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpA.Mutex.SetMessage("locked by client A")
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpB.Mutex.SetMessage("locked by client B")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.TryLock(context.Background())
		if err == nil {
			t.Fatal("TryLock should have returned an error")
		}
		var mutexErr MutexErr
		if !errors.As(err, &mutexErr) {
			t.Fatal("TryLock should have returned a MutexErr")
		} else {
			log.Printf("got expected MutexErr: %s", mutexErr.Error())
		}
		if mutexErr.LockInfo != nil {
			t.Fatal("mutexErr's LockInfo should be nil")
		}
		if mutexErr.Mutex == nil {
			t.Fatal("mutexErr's Mutex should not be nil")
		}
		if mutexErr.Mutex.GetMessage() != bpA.Mutex.GetMessage() {
			t.Fatal("reason and original lock messages do not match")
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestTryLockTryLockBlueprintMutex verifies expected outcome when TryLock() is
// attempted against a previously TryLock()ed blueprint.
func TestTryLockTryLockBlueprintMutex(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for clientName, client := range clients {
		log.Printf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpA, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpB, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpA.Mutex.SetMessage("locked by client A")
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpB.Mutex.SetMessage("locked by client B")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.TryLock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.TryLock(context.Background())
		if err == nil {
			t.Fatal("TryLock should have returned an error")
		}
		var mutexErr MutexErr
		if !errors.As(err, &mutexErr) {
			t.Fatal("TryLock should have returned a MutexErr")
		}
		if mutexErr.Mutex == nil {
			t.Fatal("TryLock returned a MutexErr with nil Mutex")
		}
		if mutexErr.LockInfo != nil {
			t.Fatal("TryLock returned a MutexEerr with non-nil LockInfo")
		}

		if mutexErr.Mutex.GetMessage() != bpA.Mutex.GetMessage() {
			t.Fatal("blocking mutex and original messages do not match")
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestLockBlockTrylockUnlLockTrylockBlueprintMutex ensures that TryLock() fails
// against a locked blueprint, then succeeds after the lock is released.
func TestLockUnlLockTrylockBlueprintMutex(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for clientName, client := range clients {
		log.Printf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpA, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpB, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpA.Mutex.SetMessage("locked by client A")
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpB.Mutex.SetMessage("locked by client B")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.TryLock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestReadOnlyBlueprintMutex gets a read-only mutex (reason) by calling
// TryLock() against a locked blueprint. It then calls various "write" methods
// on that mutex to ensure they fail with the expected error type.
func TestReadOnlyBlueprintMutex(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for clientName, client := range clients {
		log.Printf("testing CreateBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpA, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpB, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpA.Mutex.SetMessage("locked by client A")
		if err != nil {
			t.Fatal(err)
		}

		log.Print("testing SetMessage()")
		err = bpB.Mutex.SetMessage("locked by client B")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.TryLock(context.Background())
		if err == nil {
			t.Fatal("TryLock should have returned an error")
		}
		var mutexErr MutexErr
		if !errors.As(err, &mutexErr) {
			t.Fatal("TryLock should have returned a MutexErr")
		}
		if mutexErr.Mutex == nil {
			t.Fatal("TryLock should have returned a *Mutex")
		}

		roMutex := mutexErr.Mutex

		log.Printf("testing Lock() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = roMutex.Lock(context.Background())
		if err == nil {
			t.Fatal("locking a read-only Mutex should return an error")
		}
		if !strings.Contains(err.Error(), "read") {
			t.Fatal("locking a read-only Mutex complain about it being read-only")
		}

		log.Printf("testing TryLock() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = roMutex.TryLock(context.Background())
		if err == nil {
			t.Fatal("try-locking a read-only Mutex should return an error")
		}
		if !strings.Contains(err.Error(), "read") {
			t.Fatal("try-locking a read-only Mutex complain about it being read-only")
		}

		log.Printf("testing Unlock() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = roMutex.Unlock(context.Background())
		if err == nil {
			t.Fatal("unlocking a read-only Mutex should return an error")
		}
		if !strings.Contains(err.Error(), "read") {
			t.Fatal("unlocking a read-only Mutex complain about it being read-only")
		}

		log.Printf("testing SetMessage() against locked mutex %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = roMutex.SetMessage("fail please")
		if err == nil {
			t.Fatal("setting message of a read-only Mutex should return an error")
		}
		if !strings.Contains(err.Error(), "read") {
			t.Fatal("setting message of a read-only Mutex complain about it being read-only")
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}
