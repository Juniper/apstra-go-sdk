//go:build integration
// +build integration

package goapstra

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

		err = bp.Mutex.SetMessage("locked by goapstra test")
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
		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
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

		err = bp.Mutex.SetMessage("locked by goapstra test")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bp.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		start := time.Now()
		timeout := 2 * time.Second
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		saved := bp.Mutex.tagId // save tag ID because...
		bp.Mutex.tagId = ""     // bypass local re-lock safety check

		err = bp.Mutex.SetMessage("re-locked by goapstra test")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bp.Mutex.Lock(timeoutCtx)
		if err != nil && !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Fatal(err)
		}

		elapsed := time.Since(start)
		if elapsed < timeout {
			t.Fatalf("we should have waited %s for lock, but only %s has elapsed", timeout.String(), elapsed.String())
		}

		bp.Mutex.tagId = saved // restore tag ID

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
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
		ok, reason, err := bpB.Mutex.TryLock(context.Background(), false)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatal("TryLock should have failed but reports OK")
		}
		if reason == nil {
			t.Fatal("TryLock should have failed, but reason is nil")
		}

		if reason.ID() != bpA.Mutex.ID() {
			t.Fatal("reason and original lock IDs do not match")
		}

		if reason.GetMessage() != bpA.Mutex.GetMessage() {
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
		ok, reason, err := bpA.Mutex.TryLock(context.Background(), false)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("TryLock should have succeeded but reports not OK")
		}
		if reason != nil {
			t.Fatal("TryLock should have succeeded, but reason is not nil")
		}

		log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ok, reason, err = bpB.Mutex.TryLock(context.Background(), false)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatal("TryLock should have failed but reports OK")
		}
		if reason == nil {
			t.Fatal("TryLock should have failed, but reason is nil")
		}

		if reason.ID() != bpA.Mutex.ID() {
			t.Fatal("reason and original lock IDs do not match")
		}

		if reason.GetMessage() != bpA.Mutex.GetMessage() {
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

// TestLockBlockTrylockUnlLockTrylockBlueprintMutex ensures that TryLock() fails
// against a locked blueprint, then succeeds after the lock is released.
func TestLockBlockTrylockUnlLockTrylockBlueprintMutex(t *testing.T) {
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
		ok, reason, err := bpB.Mutex.TryLock(context.Background(), false)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatal("TryLock returned OK when blueprint should have been locked")
		} else {
			log.Printf("trylock failed (good) with reason: %q (tagId: %s)", reason.GetMessage(), reason.ID())
		}

		if reason == nil {
			t.Fatal("reason is nil when blueprint should have been locked")
		}

		if reason.ID() != bpA.Mutex.ID() {
			t.Fatal("tagIDs of blocking lock and reported reason do not match")
		}

		if reason.GetMessage() != bpA.Mutex.GetMessage() {
			t.Fatal("messages in blocking lock and reported reason do not match")
		}

		log.Printf("testing Unlock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing TryLock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ok, reason, err = bpB.Mutex.TryLock(context.Background(), false)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("TryLock returned not OK when blueprint should have been unlocked")
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
		ok, reason, err := bpB.Mutex.TryLock(context.Background(), false)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatal("TryLock should have failed but reports OK")
		}
		if reason == nil {
			t.Fatal("TryLock should have failed, but reason is nil")
		}

		if reason.ID() != bpA.Mutex.ID() {
			t.Fatal("reason and original lock IDs do not match")
		}

		if reason.GetMessage() != bpA.Mutex.GetMessage() {
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

		var ace ApstraClientErr

		log.Printf("testing Unlock() with read-only mutex against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = reason.Unlock(context.Background())
		if err == nil {
			t.Fatal("Unlock() of read-only mutex should have failed")
		}
		if !errors.As(err, &ace) {
			t.Fatal("Unlock() of read-only mutex should have produced an ApstraClientErr")
		}
		if ace.errType != ErrReadOnly {
			t.Fatal("Unlock() of read-only mutex should have produced an ApstraClientErr of type ErrReadOnly")
		}

		log.Printf("testing Lock() with read-only mutex against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err2 := reason.Lock(context.Background())
		if err2 == nil {
			t.Fatal("Lock() of read-only mutex should have failed")
		}
		if !errors.As(err2, &ace) {
			t.Fatal("Lock() of read-only mutex should have produced an ApstraClientErr")
		}
		if ace.errType != ErrReadOnly {
			t.Fatal("Lock() of read-only mutex should have produced an ApstraClientErr of type ErrReadOnly")
		}

		log.Printf("testing TryLock() with read-only mutex against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, _, err = reason.TryLock(context.Background(), false)
		if err == nil {
			t.Fatal("TryLock() of read-only mutex should have failed")
		}
		if !errors.As(err, &ace) {
			t.Fatal("TryLock() of read-only mutex should have produced an ApstraClientErr")
		}
		if ace.errType != ErrReadOnly {
			t.Fatal("TryLock() of read-only mutex should have produced an ApstraClientErr of type ErrReadOnly")
		}

		log.Printf("testing SetMessage() with read-only mutex against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = reason.SetMessage("")
		if err == nil {
			t.Fatal("SetMessage() of read-only mutex should have failed")
		}
		if !errors.As(err, &ace) {
			t.Fatal("SetMessage() of read-only mutex should have produced an ApstraClientErr")
		}
		if ace.errType != ErrReadOnly {
			t.Fatal("SetMessage() of read-only mutex should have produced an ApstraClientErr of type ErrReadOnly")
		}
	}
}

// TestLockClearBlueprintMutex is a simple test of Lock() followed by ClearUnsafely()
func TestLockClearBlueprintMutex(t *testing.T) {
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

		err = bpA.Mutex.SetMessage("locked by goapstra test client A")
		if err != nil {
			t.Fatal(err)
		}

		err = bpB.Mutex.SetMessage("locked by goapstra test client B")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpA.Mutex.Lock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing ClearUnsafely() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.ClearUnsafely(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Lock() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpB.Mutex.Lock(context.Background())
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
