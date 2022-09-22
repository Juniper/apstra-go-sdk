//go:build integration
// +build integration

package goapstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sort"
	"testing"
	"time"
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
	rand.Seed(time.Now().UnixNano())
	clients, err := getTestClients()
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

		for _ = range asnRanges {
			// delete one randomly selected range
			rangeCount := len(newPool.Ranges)
			deleteMe := newPool.Ranges[rand.Intn(rangeCount)]
			err = client.client.DeleteAsnPoolRange(context.TODO(), newPoolId, &deleteMe)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetAsnPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
			newPool, err = client.client.GetAsnPool(context.TODO(), newPoolId)
			if err != nil {
				t.Fatal(err)
			}

			if rangeCount-1 != len(newPool.Ranges) {
				t.Fatalf("expected new pool to have %d ranges, got %d", rangeCount-1, len(newPool.Ranges))
			}
		}

		if len(newPool.Ranges) != 0 {
			t.Fatalf("expected new pool to have 0 ranges, got %d", len(newPool.Ranges))
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

func TestCreateDeleteAsnPoolRange(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	name := "test-" + randString(10, "hex")
	var tags []string
	tags = append(tags, "tag-"+randString(10, "hex"))
	tags = append(tags, "tag-"+randString(10, "hex"))

	for clientName, client := range clients {
		log.Printf("testing CreateAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		poolId, err := client.client.CreateAsnPool(context.TODO(), &AsnPoolRequest{
			DisplayName: name,
			Tags:        tags,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, poolId)
		var asnRange IntRangeRequest
		for i := 0; i < 3; i++ {
			a := rand.Intn(1000) + (i * 1000 * 2)
			b := rand.Intn(1000) + a
			asnRange.First = uint32(a)
			asnRange.Last = uint32(b)
			log.Printf("testing CreateAsnPoolRange() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.CreateAsnPoolRange(context.TODO(), poolId, &asnRange)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing GetAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		asnPool, err := client.client.GetAsnPool(context.TODO(), poolId)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range asnPool.Ranges {
			log.Printf("testing DeleteAsnPoolRange() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err := client.client.DeleteAsnPoolRange(context.TODO(), asnPool.Id, &IntRangeRequest{First: r.First, Last: r.Last})
			if err != nil {
				t.Fatal(err)
			}
		}
		log.Printf("testing DeleteAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteAsnPool(context.TODO(), asnPool.Id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetAsnPoolByName(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	poolName := randString(10, "hex")

	for clientName, client := range clients {
		log.Printf("testing getAsnPoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		p, err := client.client.getAsnPoolByName(context.Background(), poolName)
		if err == nil {
			t.Fatal("fetching pool with random name should have earned us a 404")
		}

		log.Printf("testing createAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.createAsnPool(context.Background(), &AsnPoolRequest{DisplayName: poolName})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getAsnPoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		p, err = client.client.getAsnPoolByName(context.Background(), poolName)
		if err != nil {
			t.Fatal(err)
		}

		if id != p.Id {
			t.Fatalf("expected '%s', got '%s", id, p.Id)
		}

		log.Printf("testing createAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		dupId, err := client.client.createAsnPool(context.Background(), &AsnPoolRequest{DisplayName: poolName})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getAsnPoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		p, err = client.client.getAsnPoolByName(context.Background(), poolName)
		if err == nil {
			t.Fatalf("expected error: pool '%s' and '%s' both should be named '%s'", id, dupId, poolName)
		}

		err = client.client.deleteAsnPool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		err = client.client.deleteAsnPool(context.TODO(), dupId)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestListIpPools(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listIp4PoolIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		poolIds, err := client.client.listIp4PoolIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(poolIds) <= 0 {
			t.Fatalf("only got %d pools", len(poolIds))
		}
	}
}

func TestGetAllIpPools(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getIp4Pools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.getIp4Pools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(pools) <= 0 {
			t.Fatalf("only got %d pools", len(pools))
		}
		log.Printf("pool count: %d", len(pools))
	}
}

func TestGetIp4PoolByName(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getIp4Pools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.getIp4Pools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		poolNames := make(map[string]struct{})
		for _, p := range pools {
			poolNames[p.DisplayName] = struct{}{}
		}

		delete(poolNames, "")
		for name := range poolNames {
			log.Printf("testing getIp4PoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			_, err := client.client.getIp4PoolByName(context.TODO(), name)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestCreateGetDeleteIp4Pool(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing createIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.createIp4Pool(context.TODO(), &NewIp4PoolRequest{
			DisplayName: randString(10, "hex"),
			Tags:        []string{"tag one", "tag two"},
		})
		if err != nil {
			t.Fatal(err)
		}

		_, s, err := net.ParseCIDR("10.1.2.3/24")
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing addSubnetToIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.addSubnetToIp4Pool(context.TODO(), id, s)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pool, err := client.client.getIp4Pool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pool.Id, pool.Total)

		log.Printf("testing deleteSubnetFromIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteSubnetFromIp4Pool(context.TODO(), id, s)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing deleteIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteIp4Pool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
