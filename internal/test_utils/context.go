// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package testutils

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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

	switch v := parent.Value(apstra.CtxKeyTestUUID).(type) {
	case uuid.UUID:
		UUID = &v
	default:
		UUID = ToPtr(newUUID(t))
		parent = context.WithValue(parent, apstra.CtxKeyTestUUID, *UUID)
		log.Println("Test UUID: ", UUID.String())
	}

	return context.WithValue(parent, apstra.CtxKeyTestID, UUID.String()+"/"+t.Name())
}

func newUUID(t testing.TB) uuid.UUID {
	t.Helper()

	result, err := uuid.NewRandom()
	require.NoError(t, err)
	return result
}
