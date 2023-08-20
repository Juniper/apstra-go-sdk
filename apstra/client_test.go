//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

func TestClientLog(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		client.client.Logf(1, "log test - client '%s'", clientName)

	}
}

func TestLoginEmptyPassword(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing empty password Login() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		client.client.cfg.Pass = ""
		err := client.client.Login(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging in with empty password, did not get errror"))
		}
	}
}

func TestLoginBadPassword(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing bad password Login() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		client.client.cfg.Pass = randString(10, "hex")
		err = client.client.Login(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging in with bad password, did not get errror"))
		}
	}
}

func TestLogoutAuthFail(t *testing.T) {
	ctx := context.Background()
	clientCfgs, err := getTestClientCfgs(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for name, cfg := range clientCfgs {
		client, err := cfg.cfg.NewClient(ctx)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing Login() against %s %s (%s)", cfg.cfgType, name, client.ApiVersion())
		err = client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("client has this authtoken: '%s'", client.httpHeaders[apstraAuthHeader])
		client.httpHeaders[apstraAuthHeader] = randJwt()
		log.Printf("client authtoken changed to: '%s'", client.httpHeaders[apstraAuthHeader])
		log.Printf("testing failed Logout() against %s %s (%s)", cfg.cfgType, name, client.ApiVersion())
		err = client.Logout(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging out with bad token, did not get errror"))
		}
	}
}

func TestGetBlueprintOverlayControlProtocol(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		bpFunc      func(context.Context, *testing.T, *Client) (*TwoStageL3ClosClient, func(context.Context) error)
		expectedOcp OverlayControlProtocol
	}

	testCases := []testCase{
		{bpFunc: testBlueprintA, expectedOcp: OverlayControlProtocolEvpn},
		{bpFunc: testBlueprintB, expectedOcp: OverlayControlProtocolNone},
	}

	for clientName, client := range clients {
		for i := range testCases {
			bpClient, bpDel := testCases[i].bpFunc(ctx, t, client.client)
			defer func() {
				err := bpDel(ctx)
				if err != nil {
					t.Fatal(err)
				}
			}()

			log.Printf("testing BlueprintOverlayControlProtocol() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ocp, err := bpClient.client.BlueprintOverlayControlProtocol(ctx, bpClient.blueprintId)
			if err != nil {
				t.Fatal(err)
			}

			if ocp != testCases[i].expectedOcp {
				t.Fatalf("expected overlay control protocol %q, got %q", testCases[i].expectedOcp.String(), ocp.String())
			}
			log.Printf("blueprint %q has overlay control protocol %q", bpClient.blueprintId, ocp.String())
		}
	}
}

func TestCRUDIntegerPools(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	validate := func(req *IntPoolRequest, resp *IntPool) {
		if req.DisplayName != resp.DisplayName {
			t.Fatalf("integer pool name mismatch - Requested: %q Actual %q", req.DisplayName, resp.DisplayName)
		}

		if len(req.Ranges) != len(resp.Ranges) {
			t.Fatalf("integer pool range length mismatch - Requested: %d Actual %d", len(req.Ranges), len(resp.Ranges))
		}

		for i := 0; i < len(req.Ranges); i++ {
			if req.Ranges[i].first() != resp.Ranges[i].first() || req.Ranges[i].last() != resp.Ranges[i].last() {
				t.Fatalf("integer pool range %d mismatch: Requested: %d-%d Actual: %d-%d", i,
					req.Ranges[i].first(), req.Ranges[i].last(),
					resp.Ranges[i].first(), resp.Ranges[i].last())
			}
		}

		if len(req.Tags) != len(resp.Tags) {
			t.Fatalf("integer pool tags length mismatch - Requested: %d Actual %d", len(req.Tags), len(resp.Tags))
		}

		sort.Strings(req.Tags)
		sort.Strings(resp.Tags)
		for i := 0; i < len(req.Tags); i++ {
			if req.Tags[i] != resp.Tags[i] {
				t.Fatalf("integer pool tag set mismatch: Requested: [%s], Actual: [%s]",
					strings.Join(req.Tags, ","), strings.Join(resp.Tags, ","))
			}
		}
	}

	randomTags := func(min, max int) []string {
		var result []string
		for i := 0; i < rand.Intn(max-min)+min; i++ {
			result = append(result, randString(5, "hex"))
		}
		return result
	}

	randomRanges := func(minRanges, maxRanges int, minVal, maxVal uint32) []IntfIntRange {
		rangeCount := rand.Intn(maxRanges-minRanges) + minRanges
		valMap := make(map[int]struct{})
		for len(valMap) < rangeCount*2 {
			valMap[rand.Intn(int(maxVal-minVal))+int(minVal)] = struct{}{}
		}
		valSlice := make([]int, len(valMap))
		var i int
		for k := range valMap {
			valSlice[i] = k
			i++
		}
		sort.Ints(valSlice)

		result := make([]IntfIntRange, rangeCount)
		for i = 0; i < rangeCount; i++ {
			result[i] = &IntRangeRequest{
				First: uint32(valSlice[(i * 2)]),
				Last:  uint32(valSlice[(i*2)+1]),
			}
		}
		return result
	}

	for clientName, client := range clients {
		log.Printf("testing GetIntegerPools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err := client.client.GetIntegerPools(ctx)
		if err != nil {
			t.Fatal(err)
		}

		beforePoolCount := len(pools)
		request := IntPoolRequest{
			DisplayName: randString(5, "hex"),
			Ranges:      randomRanges(2, 5, 1, math.MaxUint32),
			Tags:        randomTags(2, 5),
		}

		log.Printf("testing CreateIntegerPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateIntegerPool(ctx, &request)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetIntegerPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pool, err := client.client.GetIntegerPool(ctx, id)
		if err != nil {
			t.Fatal(err)
		}

		validate(&request, pool)

		log.Printf("testing GetIntegerPools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err = client.client.GetIntegerPools(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(pools) != beforePoolCount+1 {
			t.Fatalf("pools before creation: %d; after creation: %d", beforePoolCount, len(pools))
		}

		poolIdx := -1
		for i, pool := range pools {
			if pool.Id == id {
				poolIdx = i
				break
			}
		}
		if poolIdx < 0 {
			t.Fatal("just-created pool id not found among pools")
		}

		validate(&request, &pools[poolIdx])

		poolIds, err := client.client.ListIntegerPoolIds(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(poolIds) != beforePoolCount+1 {
			t.Fatalf("expected %d pool IDs, got %d", beforePoolCount+1, len(poolIds))
		}

		var found bool
		for _, poolId := range poolIds {
			if poolId == id {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("newly created pool ID not found among pool ID list")
		}

		request = IntPoolRequest{
			DisplayName: randString(5, "hex"),
			Ranges:      randomRanges(2, 5, 1, math.MaxUint32),
			Tags:        randomTags(2, 5),
		}

		log.Printf("testing UpdateIntegerPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.UpdateIntegerPool(ctx, id, &request)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetIntegerPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pool, err = client.client.GetIntegerPool(ctx, id)
		if err != nil {
			t.Fatal(err)
		}

		validate(&request, pool)

		log.Printf("testing DeleteIntegerPool() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteIntegerPool(ctx, id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetIntegerPools() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		pools, err = client.client.GetIntegerPools(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(pools) != beforePoolCount {
			t.Fatalf("pools before creation: %d; after creation: %d", beforePoolCount, len(pools))
		}
	}
}
