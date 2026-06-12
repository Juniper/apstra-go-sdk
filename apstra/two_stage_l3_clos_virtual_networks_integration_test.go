// Copyright (c) Juniper Networks, Inc., 2023-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/stretchr/testify/require"
)

func TestVirtualNetworkTags(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	randStr := randString(5, "hex")
	label := "test-" + randStr
	vrfName := "test-" + randStr

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			if !compatibility.VirtualNetworkTags.Check(client.client.apiVersion) {
				t.Skipf("skipping virtual network tag test with version %s", client.client.apiVersion)
			}

			t.Parallel()

			bpClient := testBlueprintC(ctx, t, client.client)

			log.Printf("testing CreateSecurityZone() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			zoneId, err := bpClient.CreateSecurityZone(ctx, datacenter.SecurityZone{
				Type:    enum.SecurityZoneTypeEVPN,
				VRFName: vrfName,
				Label:   label,
			})
			require.NoError(t, err)

			var result struct {
				Items []struct {
					System struct {
						SystemId string `json:"id"`
					} `json:"system"`
				} `json:"items"`
			}

			query := new(PathQuery).
				SetClient(client.client).
				SetBlueprintId(bpClient.Id()).
				Node([]QEEAttribute{
					{"type", QEStringVal("system")},
					{"system_type", QEStringVal("switch")},
					{"role", QEStringVal("leaf")},
					{"name", QEStringVal("system")},
				})

			err = query.Do(ctx, &result)
			require.NoError(t, err)

			vnBindings := make([]datacenter.VNBinding, len(result.Items))
			for i := range result.Items {
				leafId := result.Items[i].System.SystemId
				vnBindings[i] = datacenter.VNBinding{
					SystemID: leafId,
				}
			}

			create := datacenter.VirtualNetwork{
				IPv4Enabled:               true,
				Label:                     label,
				SecurityZoneID:            zoneId,
				VirtualGatewayIPv4Enabled: true,
				Bindings:                  vnBindings[:1],
				Type:                      enum.VnTypeVxlan,
			}

			log.Printf("testing CreateVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vnId, err := bpClient.CreateVirtualNetwork(ctx, create)
			require.NoError(t, err)

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err := bpClient.GetVirtualNetwork(ctx, vnId)
			require.NoError(t, err)
			require.Equal(t, 0, len(vn.Tags))

			tags := make([]string, rand.Intn(10)+1)
			for i := range tags {
				tags[i] = randString(6, "hex")
			}

			log.Printf("setting tags on virtual network against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetNodeTags(ctx, ObjectId(*vn.ID()), tags)
			require.NoError(t, err)

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err = bpClient.GetVirtualNetwork(ctx, vnId)
			require.NoError(t, err)
			require.Equal(t, len(tags), len(vn.Tags))
			compareSlicesAsSets(t, tags, vn.Tags, "tag set mismatch")

			log.Printf("clearing tags on virtual network against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.SetNodeTags(ctx, ObjectId(*vn.ID()), nil)
			require.NoError(t, err)

			log.Printf("testing GetVirtualNetwork() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			vn, err = bpClient.GetVirtualNetwork(ctx, vnId)
			require.NoError(t, err)
			require.Equal(t, 0, len(vn.Tags))

			create.Tags = tags
			create.Label = randString(6, "hex")

			if compatibility.VirtualNetworkAPITags.Check(client.client.apiVersion) {
				log.Printf("ensuring that CreateVirtualNetwork() with tags succeeds against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				id, err := bpClient.CreateVirtualNetwork(ctx, create)
				require.NoError(t, err)

				require.NoError(t, create.SetID(id))
				log.Printf("ensuring that UpdateVirtualNetwork() with tags succeeds against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = bpClient.UpdateVirtualNetwork(ctx, create)
				require.NoError(t, err)
			} else {
				log.Printf("ensuring that CreateVirtualNetwork() with tags fails against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				var ace ClientErr
				_, err := bpClient.CreateVirtualNetwork(ctx, create)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ErrNotSupported, ace.Type())

				log.Printf("ensuring that UpdateVirtualNetwork() with tags fails against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				require.NoError(t, create.SetID(vnId))
				err = bpClient.UpdateVirtualNetwork(ctx, create)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ErrNotSupported, ace.Type())
			}
		})
	}
}
