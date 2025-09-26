// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package testutils

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	CtxKeyTestUUID = "Test-UUID"
	CtxKeyTestID   = "Test-ID"
)

// ContextWithTestID produces contexts with the following values:
// - Test-UUID: a uuid.UUID representing this test and all sub-tests.
// - Test-ID: a string of the form uuid/test/subtest/subsubtest...
// the Test-UUID is generated only if not found.
// HTTP transactions related to these tests can be picked out from wireshark
// using filters like:
// - http.request.line contains "843a754c-cc35-4383-807f-833ad991e554"
// - http.request.line contains "843a754c-cc35-4383-807f-833ad991e554/test/subtest"
func ContextWithTestID(parent context.Context, t testing.TB) context.Context {
	var UUID *uuid.UUID

	switch v := parent.Value(CtxKeyTestUUID).(type) {
	case uuid.UUID:
		UUID = &v
	default:
		UUID = pointer.To(newUUID(t))
		parent = context.WithValue(parent, CtxKeyTestUUID, *UUID)
		log.Println(CtxKeyTestUUID, ": ", UUID.String())
	}

	return context.WithValue(parent, CtxKeyTestID, UUID.String()+"/"+t.Name())
}

func newUUID(t testing.TB) uuid.UUID {
	t.Helper()

	result, err := uuid.NewRandom()
	require.NoError(t, err)
	return result
}

func CleanupWithFreshContext(t testing.TB, timeout time.Duration, f func(ctx context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := f(ctx)
		if !assert.NoError(t, err) {
			t.Logf("Cleanup test %q: %v", t.Name(), err)
		}
	})
}
