// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"slices"
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

	// Range returns []int where each value matches the index of that value e.g. Range(3) -> []int{0, 1, 2}
	Range := func(n int) []int {
		r := make([]int, n)
		for i := range n {
			r[i] = i
		}
		return r
	}

	sampleSizeStr, envFound := os.LookupEnv(envSampleSize)
	if !envFound && len(count) == 0 {
		return Range(length) // sample size not specified
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
		// requested sample size larger than available data
		return Range(length) // return available indexes
	}

	if float64(sampleSize) > (float64(length) * .75) {
		// requested sample size is close to actual size -- use random deletions instead of selections
		result := Range(length) // start with a too-big prototype

		for range length - sampleSize { // delete extra elements
			delIdx := rand.Intn(len(result))
			result = slices.Delete(result, delIdx, delIdx+1)
		}

		return result
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

func GetRandInts(min, max, count int) ([]int, error) {
	if max-min+1 < count {
		return nil, fmt.Errorf("cannot generate %d random numbers between %d and %d inclusive", count, min, max)
	}
	intMap := make(map[int]struct{}, count)
	for len(intMap) < count {
		intMap[rand.Intn(max+1-min)+min] = struct{}{}
	}
	result := make([]int, count)
	i := 0
	for k := range intMap {
		result[i] = k
		i++
	}
	return result, nil
}

func RandomPrefix(t testing.TB, cidrBlock string, bits int) net.IPNet {
	t.Helper()

	ip, block, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		t.Fatalf("RandomPrefix cannot parse cidrBlock - %s", err)
	}
	if block.IP.String() != ip.String() {
		t.Fatal("invocation of RandomPrefix doesn't use a base block address")
	}

	mOnes, mBits := block.Mask.Size()
	if mOnes >= bits {
		t.Fatalf("cannot select a random /%d from within %s", bits, cidrBlock)
	}

	// generate a completely random address
	randomIP := make(net.IP, mBits/8)
	_, err = crand.Read(randomIP)
	if err != nil {
		t.Fatalf("rand read failed")
	}

	// mask off the "network" bits
	for i, b := range randomIP {
		mBitsThisByte := min(mOnes, 8)
		mOnes -= mBitsThisByte
		block.IP[i] = block.IP[i] | (b & byte(math.MaxUint8>>mBitsThisByte))
	}

	block.Mask = net.CIDRMask(bits, mBits)

	_, result, err := net.ParseCIDR(block.String())
	if err != nil {
		t.Fatal("failed to parse own CIDR block")
	}

	return *result
}

func RandJWT() string {
	return RandString(36, "b64") + "." +
		RandString(178, "b64") + "." +
		RandString(86, "b64")
}
