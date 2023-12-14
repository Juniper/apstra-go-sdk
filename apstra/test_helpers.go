package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

const envSampleSize = "GOAPSTRA_TEST_SAMPLE_MAX"

func intsFromZero(length int) []int {
	result := make([]int, length)
	for i := range result {
		result[i] = i
	}
	return result
}

func samples(length int) []int {
	var sampleSizeStr string
	var sampleSizeInt int
	var found bool
	if sampleSizeStr, found = os.LookupEnv(envSampleSize); !found {
		return intsFromZero(length)
	}
	sampleSizeInt, _ = strconv.Atoi(sampleSizeStr)
	if sampleSizeInt == 0 {
		return intsFromZero(length)
	}
	if sampleSizeInt > length {
		return intsFromZero(length)
	}

	sampleMap := make(map[int]struct{})
	for len(sampleMap) < sampleSizeInt {
		sampleMap[rand.Intn(length)] = struct{}{}
	}

	result := make([]int, len(sampleMap))
	i := 0
	for k := range sampleMap {
		result[i] = k
		i++
	}
	return result
}

func compareSlices[A comparable](t *testing.T, a, b []A, info string) {
	if len(a) != len(b) {
		t.Fatalf("%s slice length mismatch: %d vs %d", info, len(a), len(b))
	}

	for i := range a {
		if a[i] != b[i] {
			as, ok := interface{}(a[i]).(fmt.Stringer)
			if !ok {
				t.Fatalf("%s slice element %d mismtach", info, i)
			}

			bs, ok := interface{}(b[i]).(fmt.Stringer)
			if !ok {
				t.Fatalf("%s slice element %d mismtach", info, i)
			}

			t.Fatalf("%s slice element %d mismatch %q vs. %q", info, i, as, bs)
		}
	}
}

func compareSlicesAsSets[A comparable](t *testing.T, a, b []A, info string) {
	if len(a) != len(b) {
		t.Fatalf("%s slice length mismatch: %d vs %d", info, len(a), len(b))
	}

	mapA := make(map[A]bool, len(a))
	for _, v := range a {
		mapA[v] = true
	}

	mapB := make(map[A]bool, len(b))
	for _, v := range b {
		mapB[v] = true
	}

	for k := range mapA {
		if _, ok := mapB[k]; !ok {
			t.Fatalf("%s slice contents mismatch: element %v found only in one slice", info, k)
		}
	}
}

func jsonEqual(t *testing.T, m1, m2 json.RawMessage) bool {
	var map1 interface{}
	var map2 interface{}

	var err error
	err = json.Unmarshal(m1, &map1)
	if err != nil {
		t.Fatalf("error unmarshalling string 1 : %v", err)
	}
	err = json.Unmarshal(m2, &map2)
	if err != nil {
		t.Fatalf("error unmarshalling string 1 : %v", err)
	}
	return reflect.DeepEqual(map1, map2)
}

func getKeysfromMap[A comparable](m map[A]interface{}) []A {
	keys := make([]A, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func nextInterface(in string) string {
	re := regexp.MustCompile(`\d+$`)
	loc := re.FindStringIndex(in)
	portNumStr := in[loc[0]:]
	i, err := strconv.Atoi(portNumStr)
	if err != nil {
		panic("Atoi should not have produced an error because the regex guaranteed digits here.")
	}
	beginStr := in[:loc[0]]
	return beginStr + strconv.Itoa(i+1)
}

func countSystemLinkTypes(ctx context.Context, systemId ObjectId, client *TwoStageL3ClosClient) (map[LinkType]int, int, error) {
	links, err := client.GetCablingMapLinksBySystem(ctx, systemId)
	if err != nil {
		return nil, 0, err
	}

	var lagMembers int

	result := make(map[LinkType]int)
	for _, link := range links {
		result[link.Type]++
		if link.Type == LinkTypeEthernet && link.AggregateLinkId != "" {
			lagMembers++
		}
	}

	return result, lagMembers, nil
}

func getSystemIdsByRole(ctx context.Context, bp *TwoStageL3ClosClient, role string) ([]ObjectId, error) {
	leafQuery := new(PathQuery).
		SetClient(bp.client).
		SetBlueprintId(bp.Id()).
		SetBlueprintType(BlueprintTypeStaging).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{"role", QEStringVal(role)},
			{"name", QEStringVal("n_system")},
		})

	var leafQueryResult struct {
		Items []struct {
			System struct {
				Id ObjectId `json:"id"`
			} `json:"n_system"`
		} `json:"items"`
	}

	err := leafQuery.Do(ctx, &leafQueryResult)
	if err != nil {
		return nil, err
	}

	result := make([]ObjectId, len(leafQueryResult.Items))
	for i, item := range leafQueryResult.Items {
		result[i] = item.System.Id
	}

	return result, nil
}
