// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDFFLink(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compareEndPoint := func(t testing.TB, req, resp *FreeformEthernetEndpoint) {
		t.Helper()

		require.Equal(t, req.SystemId, resp.SystemId, "system id ")
		if req.Interface.Id != nil {
			require.NotNil(t, resp.Interface)
			require.Equal(t, *req.Interface.Id, *resp.Interface.Id, "interface_id")
		}
		require.Equal(t, req.Interface.Data.IfName, resp.Interface.Data.IfName, "if_name")
		require.Equal(t, req.Interface.Data.TransformationId, resp.Interface.Data.TransformationId, "transformation_id")
		if req.Interface.Data.Ipv4Address != nil && resp.Interface.Data.Ipv4Address != nil {
			require.Equal(t, req.Interface.Data.Ipv4Address.String(), resp.Interface.Data.Ipv4Address.String(), "ipv4_address")
		} else {
			require.Nil(t, req.Interface.Data.Ipv4Address)
			require.Nil(t, resp.Interface.Data.Ipv4Address)
		}
		if req.Interface.Data.Ipv6Address != nil && resp.Interface.Data.Ipv6Address != nil {
			require.Equal(t, req.Interface.Data.Ipv6Address.String(), resp.Interface.Data.Ipv6Address.String(), "ipv6_address")
		} else {
			require.Nil(t, req.Interface.Data.Ipv6Address)
			require.Nil(t, resp.Interface.Data.Ipv6Address)
		}
		compareSlicesAsSets(t, req.Interface.Data.Tags, resp.Interface.Data.Tags, "tag mismatch")
	}

	compare := func(t testing.TB, req *FreeformLinkRequest, resp *FreeformLinkData) {
		t.Helper()

		require.NotNil(t, req)
		require.NotNil(t, resp)
		require.Equal(t, req.Label, resp.Label, "label")
		if req.Endpoints[0].SystemId == resp.Endpoints[0].SystemId {
			compareEndPoint(t, &req.Endpoints[0], &resp.Endpoints[0])
			compareEndPoint(t, &req.Endpoints[1], &resp.Endpoints[1])
		} else {
			compareEndPoint(t, &req.Endpoints[0], &resp.Endpoints[1])
			compareEndPoint(t, &req.Endpoints[1], &resp.Endpoints[0])
		}
		compareSlicesAsSets(t, req.Tags, resp.Tags, "tag mismatch")
	}

	type testCase struct {
		steps []FreeformLinkRequest
	}

	for _, client := range clients {
		ffc, intSysIds, extSysIds := testFFBlueprintB(ctx, t, client.client, 2, 2)

		testCases := map[string]testCase{
			"int_int_start_with_minimal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/0"),
								TransformationId: toPtr(1),
							}}},
							{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/0"),
								TransformationId: toPtr(1),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/1"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.0.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8::3"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/1"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.0.0.4"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8::4"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/2"),
								TransformationId: toPtr(1),
							}}},
							{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/2"),
								TransformationId: toPtr(1),
							}}},
						},
					},
				},
			},
			"int_int_start_with_maximal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/3"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.1"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::1"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/3"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/4"),
								TransformationId: toPtr(1),
							}}},
							{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/4"),
								TransformationId: toPtr(1),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/5"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: intSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/5"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::3"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
				},
			},
			"int_ext_start_with_minimal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/6"),
								TransformationId: toPtr(1),
							}}},
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/7"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.0.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8::3"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.0.0.4"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8::4"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/8"),
								TransformationId: toPtr(1),
							}}},
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
						},
					},
				},
			},
			"int_ext_start_with_maximal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/9"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.1"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::1"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/10"),
								TransformationId: toPtr(1),
							}}},
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: intSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           toPtr("ge-0/0/11"),
								TransformationId: toPtr(2),
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.1.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8:1::3"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
				},
			},
			"ext_ext_start_with_minimal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
							{SystemId: extSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.0.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8::3"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: extSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.0.0.4"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8::4"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
							{SystemId: extSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
						},
					},
				},
			},
			"ext_ext_start_with_maximal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.1.0.1"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8:1::1"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: extSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
							{SystemId: extSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEthernetEndpoint{
							{SystemId: extSysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: extSysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								Ipv4Address: &net.IPNet{IP: net.ParseIP("10.1.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address: &net.IPNet{IP: net.ParseIP("2001:db8:1::3"), Mask: net.CIDRMask(64, 128)},
								Tags:        randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
				},
			},
		}

		for tName, tCase := range testCases {
			tName, tCase := tName, tCase

			t.Run(tName, func(t *testing.T) {
				t.Parallel()

				// create the link
				id, err := ffc.CreateLink(ctx, &tCase.steps[0])
				require.NoError(t, err)

				// read the link
				link, err := ffc.GetLink(ctx, id)
				require.NoError(t, err)
				require.Equal(t, id, link.Id, "link Id After Create")
				compare(t, &tCase.steps[0], link.Data)

				// update the link once for each "step", including the first step (values used at creation)
				for i, step := range tCase.steps {
					// record the interface ID in the update request
					step.Endpoints[0].Interface.Id = link.Data.Endpoints[0].Interface.Id
					step.Endpoints[1].Interface.Id = link.Data.Endpoints[1].Interface.Id

					// update the link
					require.NoError(t, ffc.UpdateLink(ctx, id, &step))

					// read the link
					link, err = ffc.GetLink(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, link.Id, fmt.Sprintf("linkId after update iteration %d", i))
					compare(t, &step, link.Data)
				}

				// delete the link
				err = ffc.DeleteLink(ctx, id)
				require.NoError(t, err)

				var ace ClientErr

				// fetching a previously deleted link should fail
				_, err = ffc.GetLink(ctx, id)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ErrNotfound, ace.Type())

				// deleting a previously deleted link should fail
				err = ffc.DeleteLink(ctx, id)
				require.Error(t, err)
				require.ErrorAs(t, err, &ace)
				require.Equal(t, ErrNotfound, ace.Type())
			})
		}
	}
}
