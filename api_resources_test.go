//go:build integration
// +build integration

package goapstra

import (
	"context"
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

func TestEmptyVniPool(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	vniRangeCount := rand.Intn(5) + 2 // random number of ASN ranges to add to new pool
	vniBeginEnds, err := getRandInts(4096, 16777214, vniRangeCount*2)
	if err != nil {
		t.Fatal(err)
	}
	sort.Ints(vniBeginEnds) // sort so that the ASN ranges will be ([0]...[1], [2]...[3], etc.)
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
		log.Printf("created ASN pool name %s id %s", poolName, newPoolId)

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

		for range vniRanges {
			// delete one randomly selected range
			rangeCount := len(newPool.Ranges)
			deleteMe := newPool.Ranges[rand.Intn(rangeCount)]
			err = client.client.DeleteVniPoolRange(context.TODO(), newPoolId, &deleteMe)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetVniPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
			newPool, err = client.client.GetVniPool(context.TODO(), newPoolId)
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

		log.Printf("testing DeleteVniPool(%s) against %s %s (%s)", newPoolId, client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteVniPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreateDeleteVniPoolRange(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	name := "test-" + randString(10, "hex")
	var tags []string
	tags = append(tags, "tag-"+randString(10, "hex"))
	tags = append(tags, "tag-"+randString(10, "hex"))

	for clientName, client := range clients {
		log.Printf("testing CreateVniPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		poolId, err := client.client.CreateVniPool(context.TODO(), &VniPoolRequest{
			DisplayName: name,
			Tags:        tags,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, poolId)
		var vniRange IntRangeRequest
		for i := 0; i < 3; i++ {
			a := rand.Intn(1000) + (i * 1000 * 2) + 4096
			b := rand.Intn(1000) + a
			vniRange.First = uint32(a)
			vniRange.Last = uint32(b)
			log.Printf("testing CreateVniPoolRange() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.CreateVniPoolRange(context.TODO(), poolId, &vniRange)
			if err != nil {
				t.Fatal(err)
			}
		}

		log.Printf("testing GetVniPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		vniPool, err := client.client.GetVniPool(context.TODO(), poolId)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range vniPool.Ranges {
			log.Printf("testing DeleteVniPoolRange() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err := client.client.DeleteVniPoolRange(context.TODO(), vniPool.Id, &IntRangeRequest{First: r.First, Last: r.Last})
			if err != nil {
				t.Fatal(err)
			}
		}
		log.Printf("testing DeleteVniPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteVniPool(context.TODO(), vniPool.Id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListVniPoolIds(t *testing.T) {
	clients, err := getTestClients()
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
			pool, err := client.client.getIp4PoolByName(context.TODO(), name)
			if err != nil {
				t.Fatal(err)
			}

			if pool.Used == pool.Total {
				log.Fatal("every IP in the pool is in use? seems unlikely.")
			}

			for _, subnet := range pool.Subnets {
				if subnet.Used == subnet.Total {
					log.Fatal("every IP in the subnet is in use? seems unlikely.")
				}
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
