// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package design_test

import (
	"context"
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestInterfaceMapDigest_Retrieval(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			ids, err := client.Client.ListInterfaceMapDigests2(ctx)
			require.NoError(t, err)
			sort.Strings(ids)

			objs, err := client.Client.GetInterfaceMapDigests2(ctx)
			require.NoError(t, err)

			require.Equal(t, len(ids), len(objs))

			for _, id := range ids {
				t.Run("check_"+id, func(t *testing.T) {
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					require.NotNil(t, objPtr.ID())
					require.Equal(t, id, *objPtr.ID())
				})
			}

			for _, id := range ids {
				t.Run("get_"+id, func(t *testing.T) {
					ctx := testutils.ContextWithTestID(ctx, t)

					obj, err := client.Client.GetInterfaceMapDigest2(ctx, id)
					require.NoError(t, err)

					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)

					require.Equal(t, *objPtr, obj)
				})
			}
		})
	}
}
