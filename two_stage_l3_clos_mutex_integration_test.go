//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"
)

func TestLockUnlockBlueprintMutex(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for _, client := range clients {
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Println(bpId)

		bp, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		err = bp.mutex.Lock(context.Background(), "locked by goapstra test")
		if err != nil {
			t.Fatal(err)
		}

		err = bp.mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLockLockUnlockBlueprintMutex(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	bpName := randString(5, "hex")
	for _, client := range clients {
		bpId, err := client.client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      bpName,
			TemplateId: "L3_Collapsed_ESI",
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Println(bpId)

		bp, err := client.client.NewTwoStageL3ClosClient(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}

		err = bp.mutex.Lock(context.Background(), "locked by goapstra test")
		if err != nil {
			t.Fatal(err)
		}

		start := time.Now()
		timeout := 3 * time.Second
		timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		saved := bp.mutex.tagId // save tag ID because...
		bp.mutex.tagId = ""     // bypass local re-lock safety check
		err = bp.mutex.Lock(timeoutCtx, "re-locked by goapstra test")
		if err != nil && !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Fatal(err)
		}

		elapsed := time.Since(start)
		if elapsed < timeout {
			t.Fatalf("we should have waited %s for lock, but only %s has elapsed", timeout.String(), elapsed.String())
		}

		bp.mutex.tagId = saved // restore tag ID

		err = bp.mutex.Unlock(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		err = client.client.deleteBlueprint(context.Background(), bpId)
		if err != nil {
			t.Fatal(err)
		}
	}
}
