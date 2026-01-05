// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package apstra_test

import (
	"net"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetSetSecurityZoneDHCPServers(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)

			sz, err := bpClient.GetSecurityZoneByVRFName(ctx, "default")
			require.NoError(t, err)
			require.NotNil(t, sz)

			getIPs, err := bpClient.GetSecurityZoneDhcpServers(ctx, *sz.ID())
			require.NoError(t, err)
			require.Empty(t, getIPs)

			setIPs := []net.IP{
				[]byte{1, 2, 3, 4},
				[]byte{5, 6, 7, 8},
				[]byte{9, 10, 11, 12},
				[]byte{1, 2, 3, 4},
			}

			err = bpClient.SetSecurityZoneDhcpServers(ctx, *sz.ID(), setIPs)
			require.NoError(t, err)

			getIPs, err = bpClient.GetSecurityZoneDhcpServers(ctx, *sz.ID())
			require.NoError(t, err)
			require.Equal(t, len(setIPs), len(getIPs))
			for i := 0; i < len(getIPs); i++ {
				if !setIPs[i].Equal(getIPs[i]) {
					t.Fatalf("dhcp server at index %d: expected %s, got %s", i, setIPs[i].String(), getIPs[i].String())
				}
			}

			err = bpClient.SetSecurityZoneDhcpServers(ctx, *sz.ID(), nil)
			require.NoError(t, err)

			getIPs, err = bpClient.GetSecurityZoneDhcpServers(ctx, *sz.ID())
			require.NoError(t, err)
			require.Empty(t, getIPs)
		})
	}
}
