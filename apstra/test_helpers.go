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

const envSampleSize = "APSTRA_TEST_SAMPLE_MAX"

func intsFromZero(length int) []int {
	result := make([]int, length)
	for i := range result {
		result[i] = i
	}
	return result
}

// samples is intended to be used to select some sample items from a slice.
// Pass it the size of the slice, and it returns a []int representing indexes (samples)
// to be taken from the slice. The number of elements returned is controlled by an
// environment variable or by the optional "count" argument. If the sample count
// is not supplied by either environment nor count, then all indexes starting with
// zero are returned. When sample count is specified both ways, count wins.
func samples(t testing.TB, length int, count ...int) []int {
	t.Helper()

	if len(count) > 1 {
		panic("count must only have a element")
	}

	sampleSizeStr, envFound := os.LookupEnv(envSampleSize)
	if !envFound && len(count) == 0 {
		return intsFromZero(length)
	}

	var sampleSize int
	if len(count) > 0 {
		sampleSize = count[0]
	} else {
		var err error
		sampleSize, err = strconv.Atoi(sampleSizeStr)
		if err != nil {
			panic(fmt.Sprintf("env var %q (%s) failed to parse as int - %s", envSampleSize, sampleSizeStr, err))
		}
	}

	if sampleSize > length {
		return intsFromZero(length)
	}

	if float64(sampleSize) > (float64(length) * .75) {
		return intsFromZero(sampleSize)
	}

	sampleMap := make(map[int]struct{}, sampleSize)
	for len(sampleMap) < sampleSize {
		sampleMap[rand.Intn(length)] = struct{}{}
	}

	result := make([]int, sampleSize)
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

func compareSlicesAsSets[A comparable](t testing.TB, a, b []A, info string) {
	t.Helper()

	if len(a) != len(b) {
		t.Fatalf("%s slice length mismatch: %d vs %d", info, len(a), len(b))
	}

	mapA := make(map[A]struct{}, len(a))
	for _, v := range a {
		mapA[v] = struct{}{}
	}

	mapB := make(map[A]struct{}, len(b))
	for _, v := range b {
		mapB[v] = struct{}{}
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

func sliceContains[A comparable](s []A, item A) bool {
	for _, element := range s {
		if element == item {
			return true
		}
	}
	return false
}

func sliceContainsAnyOf[A comparable](s, items []A) bool {
	for _, item := range items {
		if sliceContains(s, item) {
			return true
		}
	}
	return false
}
