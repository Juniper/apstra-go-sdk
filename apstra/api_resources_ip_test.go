//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"net"
	"testing"
)

func TestListIp4Pools(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing ListIp4PoolIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		poolIds, err := client.client.ListIp4PoolIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(poolIds) <= 0 {
			t.Fatalf("only got %d pools", len(poolIds))
		}
	}
}

func TestGetAllIp4Pools(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetIp4Pools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.GetIp4Pools(context.TODO())
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
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetIp4Pools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.GetIp4Pools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		poolNames := make(map[string]struct{})
		for _, p := range pools {
			poolNames[p.DisplayName] = struct{}{}
		}

		delete(poolNames, "")
		for name := range poolNames {
			log.Printf("testing GetIp4PoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			pool, err := client.client.GetIp4PoolByName(context.TODO(), name)
			if err != nil {
				t.Fatal(err)
			}

			if pool.Used.Cmp(&pool.Total) == 0 {
				log.Fatal("every IP in the pool is in use? seems unlikely.")
			}

			for _, subnet := range pool.Subnets {
				if subnet.Used.Cmp(&subnet.Total) == 0 {
					log.Fatal("every IP in the subnet is in use? seems unlikely.")
				}
			}
		}
	}
}

func TestCreateGetDeleteIp4Pool(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing CreateIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateIp4Pool(context.TODO(), &NewIpPoolRequest{
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
		log.Printf("testing addSubnetToIpPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.addSubnetToIpPool(context.TODO(), id, s)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pool, err := client.client.GetIp4Pool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pool.Id, pool.Total)

		log.Printf("testing deleteSubnetFromIpPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteSubnetFromIpPool(context.TODO(), id, s)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteIp4Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteIp4Pool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListIp6Pools(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing ListIp6PoolIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		poolIds, err := client.client.ListIp6PoolIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(poolIds) <= 0 {
			t.Fatalf("only got %d pools", len(poolIds))
		}
	}
}

func TestGetAllIp6Pools(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetIp6Pools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.GetIp6Pools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(pools) <= 0 {
			t.Fatalf("only got %d pools", len(pools))
		}
		log.Printf("pool count: %d", len(pools))
	}
}

func TestGetIp6PoolByName(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing GetIp6Pools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.GetIp6Pools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		poolNames := make(map[string]struct{})
		for _, p := range pools {
			poolNames[p.DisplayName] = struct{}{}
		}

		delete(poolNames, "")
		for name := range poolNames {
			log.Printf("testing GetIp6PoolByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			pool, err := client.client.GetIp6PoolByName(context.TODO(), name)
			if err != nil {
				t.Fatal(err)
			}

			if pool.Used.Cmp(&pool.Total) == 0 {
				log.Fatal("every IP in the pool is in use? seems unlikely.")
			}

			for _, subnet := range pool.Subnets {
				if subnet.Used.Cmp(&subnet.Total) == 0 {
					log.Fatal("every IP in the subnet is in use? seems unlikely.")
				}
			}
		}
	}
}

func TestCreateGetDeleteIp6Pool(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing CreateIp6Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateIp6Pool(context.TODO(), &NewIpPoolRequest{
			DisplayName: randString(10, "hex"),
			Tags:        []string{"tag one", "tag two"},
		})
		if err != nil {
			t.Fatal(err)
		}

		_, s, err := net.ParseCIDR("2001:db8::/32")
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing addSubnetToIpPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.addSubnetToIpPool(context.TODO(), id, s)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetIp6Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pool, err := client.client.GetIp6Pool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pool.Id, pool.Total)

		log.Printf("testing deleteSubnetFromIpPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteSubnetFromIpPool(context.TODO(), id, s)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteIp6Pool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteIp6Pool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
