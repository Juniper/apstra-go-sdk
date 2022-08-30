package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net"
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

func TestGetCreateDeleteAsnPools(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetAsnPools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.GetAsnPools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pools)
		var poolBeginEnds []AsnRange
		for _, p := range pools {
			for _, r := range p.Ranges {
				poolBeginEnds = append(poolBeginEnds, AsnRange{First: r.first(), Last: r.last()})
			}
		}
		openHoles, err := invertRangesInRange(1, math.MaxUint32, poolBeginEnds)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("open holes in ASN resources: ", openHoles)

		// todo: make sure there's at least one open hole in the plan
		name := "test-" + randString(10, "hex")
		r := rand.Intn(len(openHoles))
		log.Printf("testing CreateAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		arr := AsnRangeRequest{
			First: openHoles[r].first(),
			Last:  openHoles[r].last(),
		}
		id, err := client.client.CreateAsnPool(context.TODO(), &AsnPoolRequest{
			Ranges:      []IntfAsnRange{arr},
			DisplayName: name,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, id)

		log.Printf("testing GetAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = client.client.GetAsnPool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteAsnPool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUpdateEmptyAsnPool(t *testing.T) {
	t.Skip("this test compares ranges across multiple pools -- needs to be revisited")
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	name := "test-" + randString(10, "hex")

	for clientName, client := range clients {
		log.Printf("testing CreateAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		newPoolId, err := client.client.CreateAsnPool(context.TODO(), &AsnPoolRequest{DisplayName: name})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, newPoolId)

		log.Printf("testing GetAsnPools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.GetAsnPools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pools)
		var poolBeginEnds []AsnRange
		for _, p := range pools {
			for _, r := range p.Ranges {
				poolBeginEnds = append(poolBeginEnds, AsnRange{First: r.First, Last: r.Last})
			}
		}
		openHoles, err := invertRangesInRange(1, math.MaxUint32, poolBeginEnds)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("open holes in ASN resources: ", openHoles)

		// todo: make sure there's at least one open hole in the plan
		r := rand.Intn(len(openHoles))
		newRange := AsnRangeRequest{
			First: openHoles[r].First,
			Last:  openHoles[r].Last,
		}
		newDisplayName := "updated-" + name
		newTags := []string{"updated"}
		log.Printf("testing updateAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.updateAsnPool(context.TODO(), newPoolId, &AsnPoolRequest{
			DisplayName: newDisplayName,
			Ranges:      []IntfAsnRange{newRange},
			Tags:        newTags,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = client.client.GetAsnPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteAsnPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteAsnPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUnmarshalPool(t *testing.T) {
	var result rawAsnPool
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
		var asnRange AsnRangeRequest
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
			err := client.client.DeleteAsnPoolRange(context.TODO(), asnPool.Id, &AsnRangeRequest{First: r.First, Last: r.Last})
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
