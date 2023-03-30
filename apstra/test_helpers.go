package apstra

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
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
	rand.Seed(time.Now().UnixNano())
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
