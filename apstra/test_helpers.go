package apstra

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
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

func stringkeysfromMap(m map[string]interface{}) []string {

	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
