// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

const envSampleSize = "APSTRA_TEST_SAMPLE_MAX"

func RandString(n int, style string) string {
	b64Letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")
	hexLetters := []rune("0123456789abcdef")
	var letters []rune
	b := make([]rune, n)
	switch style {
	case "hex":
		letters = hexLetters
	case "b64":
		letters = b64Letters
	}

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// SampleIndexes is intended to be used to select some sample items from a slice.
// Pass it the size of the slice, and it returns a []int representing indexes (samples)
// to be taken from the slice. The number of elements returned is controlled by an
// environment variable or by the optional "count" argument. If the sample count
// is not supplied by either environment nor count, then all indexes starting with
// zero are returned. When sample count is specified both ways, count wins.
func SampleIndexes(t testing.TB, length int, count ...int) []int {
	t.Helper()

	if len(count) > 1 {
		panic("count must only have a element")
	}

	sampleSizeStr, envFound := os.LookupEnv(envSampleSize)
	if !envFound && len(count) == 0 {
		return Range(length)
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

	if sampleSize == 0 {
		return []int{}
	}

	if sampleSize > length {
		return Range(length)
	}

	if float64(sampleSize) > (float64(length) * .75) {
		return Range(sampleSize)
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
