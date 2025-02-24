// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDInternalSystem(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformSystemData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Type, b.Type)
		require.Equal(t, a.Label, b.Label)
		if a.SystemId != nil && b.SystemId != nil {
			require.Equal(t, *a.SystemId, *b.SystemId)
		} else {
			require.Nil(t, a.SystemId)
			require.Nil(t, b.SystemId)
		}
		if a.Hostname != "" {
			require.Equal(t, a.Hostname, b.Hostname)
		} else {
			require.Equal(t, a.Label, b.Hostname)
		}
		compareSlicesAsSets(t, a.Tags, b.Tags, "Tags comparison mismatch")
		require.Equal(t, a.DeviceProfileId, b.DeviceProfileId)
	}

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)

		dpIdA, err := ffc.ImportDeviceProfile(ctx, "Juniper_vEX")
		require.NoError(t, err)

		dpIdB, err := ffc.ImportDeviceProfile(ctx, "Juniper_vQFX")
		require.NoError(t, err)

		cfg := FreeformSystemData{
			Label:           randString(6, "hex"),
			DeviceProfileId: &dpIdA,
			Type:            SystemTypeInternal,
		}

		id, err := ffc.CreateSystem(ctx, &cfg)
		require.NoError(t, err)

		ffSystem, err := ffc.GetSystem(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, ffSystem.Id)
		compare(t, &cfg, ffSystem.Data)

		cfg = FreeformSystemData{
			Type:            SystemTypeInternal,
			Label:           randString(6, "hex"),
			Hostname:        randString(6, "hex"),
			Tags:            []string{"tagA", "tagB"},
			DeviceProfileId: &dpIdB,
		}

		err = ffc.UpdateSystem(ctx, id, &cfg)
		require.NoError(t, err)

		ffSystem, err = ffc.GetSystem(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, ffSystem.Id)
		compare(t, &cfg, ffSystem.Data)

		cfg = FreeformSystemData{
			Type:            SystemTypeInternal,
			Label:           randString(6, "hex"),
			DeviceProfileId: &dpIdA,
		}

		err = ffc.UpdateSystem(ctx, id, &cfg)
		require.NoError(t, err)

		ffSystem, err = ffc.GetSystem(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, ffSystem.Id)
		cfg.Hostname = ffSystem.Data.Hostname // compare cannot anticipate this value.
		compare(t, &cfg, ffSystem.Data)

		ffSystems, err := ffc.GetAllSystems(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(ffSystems))
		for i, template := range ffSystems {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteSystem(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetSystem(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteSystem(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}

func TestCRUDExternalSystem(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *FreeformSystemData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Type, b.Type)
		require.Equal(t, a.Label, b.Label)
		if a.SystemId != nil && b.SystemId != nil {
			require.Equal(t, *a.SystemId, *b.SystemId)
		} else {
			require.Nil(t, a.SystemId)
			require.Nil(t, b.SystemId)
		}
		if a.Hostname != "" {
			require.Equal(t, a.Hostname, b.Hostname)
		} else {
			require.Equal(t, a.Label, b.Hostname)
		}
		compareSlicesAsSets(t, a.Tags, b.Tags, "Tags comparison mismatch")
		require.Equal(t, a.DeviceProfileId, b.DeviceProfileId)
	}

	for _, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			ffc := testFFBlueprintA(ctx, t, client.client)

			cfg := FreeformSystemData{
				Label: randString(6, "hex"),
				Type:  SystemTypeExternal,
			}

			id, err := ffc.CreateSystem(ctx, &cfg)
			require.NoError(t, err)

			ffSystem, err := ffc.GetSystem(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, ffSystem.Id)
			compare(t, &cfg, ffSystem.Data)

			cfg = FreeformSystemData{
				Type:     SystemTypeExternal,
				Label:    randString(6, "hex"),
				Hostname: randString(6, "hex"),
				Tags:     []string{"tagA", "tagB"},
			}

			err = ffc.UpdateSystem(ctx, id, &cfg)
			require.NoError(t, err)

			ffSystem, err = ffc.GetSystem(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, ffSystem.Id)
			compare(t, &cfg, ffSystem.Data)

			cfg = FreeformSystemData{
				Type:  SystemTypeExternal,
				Label: randString(6, "hex"),
			}

			err = ffc.UpdateSystem(ctx, id, &cfg)
			require.NoError(t, err)

			ffSystem, err = ffc.GetSystem(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, ffSystem.Id)
			cfg.Hostname = ffSystem.Data.Hostname // compare cannot anticipate this value.
			compare(t, &cfg, ffSystem.Data)

			ffSystems, err := ffc.GetAllSystems(ctx)
			require.NoError(t, err)
			ids := make([]ObjectId, len(ffSystems))
			for i, template := range ffSystems {
				ids[i] = template.Id
			}
			require.Contains(t, ids, id)

			err = ffc.DeleteSystem(ctx, id)
			require.NoError(t, err)

			var ace ClientErr

			_, err = ffc.GetSystem(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ErrNotfound, ace.Type())

			err = ffc.DeleteSystem(ctx, id)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, ErrNotfound, ace.Type())
		})
	}
}
