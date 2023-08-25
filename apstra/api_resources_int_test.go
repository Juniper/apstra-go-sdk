//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"testing"
)

const (
	testPool1 = `{
      "status": "not_in_use",
      "used": "0",
      "display_name": "foo",
      "tags": [],
      "created_at": "2022-06-13T18:44:55.899107Z",
      "last_modified_at": "2022-06-13T18:44:55.899107Z",
      "ranges": [],
      "used_percentage": 0,
      "total": "0",
      "id": "e49e0f45-ecf3-478d-8b1f-901af6d4ed89"
    }`
)

func TestEmptyAsnPool(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	asnRangeCount := rand.Intn(5) + 2 // random number of ASN ranges to add to new pool
	asnBeginEnds, err := getRandInts(1, 100000000, asnRangeCount*2)
	if err != nil {
		t.Fatal(err)
	}
	sort.Ints(asnBeginEnds) // sort so that the ASN ranges will be ([0]...[1], [2]...[3], etc.)
	asnRanges := make([]IntfIntRange, asnRangeCount)
	for i := 0; i < asnRangeCount; i++ {
		asnRanges[i] = IntRangeRequest{
			uint32(asnBeginEnds[2*i]),
			uint32(asnBeginEnds[(2*i)+1]),
		}
	}

	poolName := "test-" + randString(10, "hex")

	for clientName, client := range clients {
		log.Printf("testing CreateAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		newPoolId, err := client.client.CreateAsnPool(context.TODO(), &AsnPoolRequest{
			DisplayName: poolName,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", poolName, newPoolId)

		log.Printf("testing GetAsnPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
		newPool, err := client.client.GetAsnPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}

		if poolName != newPool.DisplayName {
			t.Fatalf("expected pool name '%s', got '%s'", poolName, newPool.DisplayName)
		}
		if len(newPool.Ranges) != 0 {
			t.Fatalf("expected new pool to have 0 ranges, got %d", len(newPool.Ranges))
		}

		for i := range asnRanges {
			newName := fmt.Sprintf("%s-%d", poolName, i)
			err = client.client.updateAsnPool(context.TODO(), newPoolId, &AsnPoolRequest{
				DisplayName: newName,
				Ranges:      asnRanges[:i+1],
			})
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetAsnPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
			newPool, err = client.client.GetAsnPool(context.TODO(), newPoolId)
			if err != nil {
				t.Fatal(err)
			}
			if newName != newPool.DisplayName {
				t.Fatalf("expected pool name '%s', got '%s'", newName, newPool.DisplayName)
			}
			if i+1 != len(newPool.Ranges) {
				t.Fatalf("expected new pool to have %d ranges, got %d", i+1, len(newPool.Ranges))
			}
		}

		log.Printf("testing DeleteAsnPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteAsnPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUnmarshalPool(t *testing.T) {
	var result rawIntPool
	err := json.Unmarshal([]byte(testPool1), &result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAsnPoolByName(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	poolName := randString(10, "hex")

	for clientName, client := range clients {
		log.Printf("testing getAsnPoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err := client.client.getAsnPoolByName(context.Background(), poolName)
		if err == nil {
			t.Fatal("fetching pool with random name should have earned us a 404")
		}

		log.Printf("testing createAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.createAsnPool(context.Background(), &AsnPoolRequest{DisplayName: poolName})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getAsnPoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		p, err := client.client.getAsnPoolByName(context.Background(), poolName)
		if err != nil {
			t.Fatal(err)
		}

		if id != p.Id {
			t.Fatalf("expected '%s', got '%s", id, p.Id)
		}

		err = client.client.deleteAsnPool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListAsnPoolIds(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAsnPoolIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		poolIds, err := client.client.listAsnPoolIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(poolIds) == 0 {
			t.Fatal("no pool IDs on this system?")
		}
	}
}

func TestEmptyVniPool(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	vniRangeCount := rand.Intn(5) + 2 // random number of VNI ranges to add to new pool
	vniBeginEnds, err := getRandInts(vniMin, vniMax, vniRangeCount*2)
	if err != nil {
		t.Fatal(err)
	}
	sort.Ints(vniBeginEnds) // sort so that the VNI ranges will be ([0]...[1], [2]...[3], etc.)
	vniRanges := make([]IntfIntRange, vniRangeCount)
	for i := 0; i < vniRangeCount; i++ {
		vniRanges[i] = IntRangeRequest{
			uint32(vniBeginEnds[2*i]),
			uint32(vniBeginEnds[(2*i)+1]),
		}
	}

	poolName := "test-" + randString(10, "hex")

	for clientName, client := range clients {
		log.Printf("testing CreateVniPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		newPoolId, err := client.client.CreateVniPool(context.TODO(), &VniPoolRequest{
			DisplayName: poolName,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created VNI pool name %s id %s", poolName, newPoolId)

		log.Printf("testing GetVniPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
		newPool, err := client.client.GetVniPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}

		if poolName != newPool.DisplayName {
			t.Fatalf("expected pool name '%s', got '%s'", poolName, newPool.DisplayName)
		}
		if len(newPool.Ranges) != 0 {
			t.Fatalf("expected new pool to have 0 ranges, got %d", len(newPool.Ranges))
		}

		for i := range vniRanges {
			newName := fmt.Sprintf("%s-%d", poolName, i)
			err = client.client.updateVniPool(context.TODO(), newPoolId, &VniPoolRequest{
				DisplayName: newName,
				Ranges:      vniRanges[:i+1],
			})
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetVniPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
			newPool, err = client.client.GetVniPool(context.TODO(), newPoolId)
			if err != nil {
				t.Fatal(err)
			}
			if newName != newPool.DisplayName {
				t.Fatalf("expected pool name '%s', got '%s'", newName, newPool.DisplayName)
			}
			if i+1 != len(newPool.Ranges) {
				t.Fatalf("expected new pool to have %d ranges, got %d", i+1, len(newPool.Ranges))
			}
		}

		log.Printf("testing DeleteVniPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteVniPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetVniPoolByName(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	poolName := randString(10, "hex")

	for clientName, client := range clients {
		log.Printf("testing getVniPoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err := client.client.getVniPoolByName(context.Background(), poolName)
		if err == nil {
			t.Fatal("fetching pool with random name should have earned us a 404")
		}

		log.Printf("testing createVniPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.createVniPool(context.Background(), &VniPoolRequest{DisplayName: poolName})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getVniPoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		p, err := client.client.getVniPoolByName(context.Background(), poolName)
		if err != nil {
			t.Fatal(err)
		}

		if id != p.Id {
			t.Fatalf("expected '%s', got '%s", id, p.Id)
		}

		err = client.client.deleteVniPool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListVniPoolIds(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listVniPoolIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		poolIds, err := client.client.listVniPoolIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(poolIds) == 0 {
			t.Fatal("no pool IDs on this system?")
		}
	}
}
