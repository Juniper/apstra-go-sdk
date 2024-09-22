// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUD_CT(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compare := func(t *testing.T, a, b *ConfigTemplateData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Text, b.Text)
	}

	for _, client := range clients {
		ffc := testFFBlueprintA(ctx, t, client.client)
		cfg := ConfigTemplateData{
			Label: randString(6, "hex") + ".jinja",
			Text:  randString(30, "hex"),
		}

		id, err := ffc.CreateConfigTemplate(ctx, &cfg)
		require.NoError(t, err)

		configTemplate, err := ffc.GetConfigTemplate(ctx, id)
		require.NoError(t, err)
		compare(t, &cfg, configTemplate.Data)

		cfg.Label = randString(6, "hex") + ".jinja"
		cfg.Text = randString(20, "hex")
		err = ffc.UpdateConfigTemplate(ctx, id, &cfg)
		require.NoError(t, err)

		configTemplate, err = ffc.GetConfigTemplate(ctx, id)
		require.NoError(t, err)
		compare(t, &cfg, configTemplate.Data)

		templates, err := ffc.GetAllConfigTemplates(ctx)
		require.NoError(t, err)
		ids := make([]ObjectId, len(templates))
		for i, template := range templates {
			ids[i] = template.Id
		}
		require.Contains(t, ids, id)

		err = ffc.DeleteConfigTemplate(ctx, id)
		require.NoError(t, err)

		var ace ClientErr

		_, err = ffc.GetConfigTemplate(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = ffc.DeleteRaResource(ctx, id)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
