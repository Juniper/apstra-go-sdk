// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestCreateGetUpdateGetDeletePropertySet(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	testData := apstra.PropertySetData{
		Label: testutils.RandString(10, "hex"),
	}

	sampleCount := rand.Intn(10) + 3
	vals := make(map[string]string, sampleCount)
	for i := 0; i < sampleCount; i++ {
		vals["_"+testutils.RandString(10, "hex")] = testutils.RandString(10, "hex")
	}

	var err error
	testData.Values, err = json.Marshal(vals)
	require.NoError(t, err)

	for _, client := range clients {
		psData := testData // start with clean copy of psData in each loop
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			id, err := client.Client.CreatePropertySet(ctx, &psData)
			require.NoError(t, err)

			ps, err := client.Client.GetPropertySetByLabel(ctx, psData.Label)
			require.NoError(t, err)
			require.NotNil(t, ps.Data)

			compare.PropertySetData(t, psData, *ps.Data)

			psData.Label = testutils.RandString(10, "hex")
			for i := 0; i < sampleCount; i++ {
				vals["_"+testutils.RandString(10, "hex")] = testutils.RandString(10, "hex")
			}

			psData.Values, err = json.Marshal(vals)
			require.NoError(t, err)

			err = client.Client.UpdatePropertySet(ctx, id, &psData)
			require.NoError(t, err)

			ps, err = client.Client.GetPropertySet(ctx, id)
			require.NoError(t, err)
			require.NotNil(t, ps.Data)
			compare.PropertySetData(t, psData, *ps.Data)

			err = client.Client.DeletePropertySet(ctx, id)
			require.NoError(t, err)

			var ace apstra.ClientErr

			_, err = client.Client.GetPropertySet(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), apstra.ErrNotfound)

			_, err = client.Client.GetPropertySetByLabel(ctx, psData.Label)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), apstra.ErrNotfound)

			err = client.Client.DeletePropertySet(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ace.Type(), apstra.ErrNotfound)
		})
	}
}
