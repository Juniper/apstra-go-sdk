// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestCreateGetDeleteL3CollapsedTemplate(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	dn := testutils.RandString(5, "hex")

	req := &apstra.CreateL3CollapsedTemplateRequest{
		DisplayName:   dn,
		MeshLinkCount: 1,
		MeshLinkSpeed: "10G",
		RackTypeIds:   []apstra.ObjectId{"L3_collapsed_acs"},
		RackTypeCounts: []apstra.RackTypeCount{{
			RackTypeId: "L3_collapsed_acs",
			Count:      1,
		}},
		VirtualNetworkPolicy: apstra.VirtualNetworkPolicy{OverlayControlProtocol: apstra.OverlayControlProtocolEvpn},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			id, err := client.Client.CreateL3CollapsedTemplate(ctx, req)
			require.NoError(t, err)

			template, err := client.Client.GetL3CollapsedTemplate(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, template.Id)

			err = client.Client.DeleteTemplate(ctx, id)
			require.NoError(t, err)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetL3CollapsedTemplateByName(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	name := "Collapsed Fabric ESI"

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			l3ct, err := client.Client.GetL3CollapsedTemplateByName(ctx, name)
			require.NoError(t, err)
			require.NotNil(t, l3ct.Data)
			require.Equal(t, name, l3ct.Data.DisplayName)
		})
	}
}
