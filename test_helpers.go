package goapstra

import (
	"math/rand"
	"os"
	"strconv"
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
