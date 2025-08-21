// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mutexmap_test

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"
	"testing"
	"time"

	mutexmap "github.com/Juniper/apstra-go-sdk/mutex_map"
	"github.com/stretchr/testify/require"
)

func TestMutexMap_Lock_Unlock(t *testing.T) {
	t.Parallel()

	idCount := 10
	ids := make([]string, idCount)

	buf := make([]byte, 8)
	for i := range ids {
		_, err := rand.Read(buf)
		require.NoError(t, err)
		ids[i] = base64.RawURLEncoding.EncodeToString(buf)
	}

	mm := mutexmap.NewMutexMap()

	// lock, then unlock each one
	for _, id := range ids {
		mm.Lock(id)
		mm.Unlock(id)
	}

	// lock each one
	for _, id := range ids {
		mm.Lock(id)
	}

	// unlock each one
	for _, id := range ids {
		mm.Unlock(id)
	}
}

func TestMutexMap_Unlock_Panic(t *testing.T) {
	t.Parallel()

	mm := mutexmap.NewMutexMap()

	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	require.NoError(t, err)
	id := base64.RawURLEncoding.EncodeToString(buf)

	require.Panics(t, func() {
		mm.Unlock(id)
	})
}

func TestMutexMap_Synctest(t *testing.T) {
	t.Parallel()

	precision := 100 * time.Millisecond

	type testcase struct {
		unlockAfter time.Duration
	}

	testCases := map[string]time.Duration{
		"one_second":    1 * time.Second,
		"two_seconds":   2 * time.Second,
		"three_seconds": 3 * time.Second,
		"four_seconds":  4 * time.Second,
	}

	testMap := mutexmap.NewMutexMap()

	elapsedMap := make(map[string]time.Duration) // elapsed times stored here
	elapsedMutex := new(sync.Mutex)              // protects elapsedMap from concurrent access

	// lock everything and schedule future unlocks
	for tName, tCase := range testCases {
		testMap.Lock(tName)

		go func() {
			time.Sleep(tCase)
			testMap.Unlock(tName)
		}()
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			start := time.Now()

			testMap.Lock(tName)
			elapsed := time.Since(start)

			elapsedMutex.Lock()
			elapsedMap[tName] = elapsed
			elapsedMutex.Unlock()

			require.Greater(t, elapsedMap[tName], tCase-precision)
			require.Less(t, elapsedMap[tName], tCase+precision)

			log.Printf("minimum: %s actual: %s: maximum: %s", tCase-precision, elapsedMap[tName], tCase+precision)
		})
	}
}
