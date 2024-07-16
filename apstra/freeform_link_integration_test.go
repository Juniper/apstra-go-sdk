//go:build integration
// +build integration

package apstra

import (
	"context"
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDFFLink(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compareEndPoint := func(t testing.TB, req, resp *FreeformEndpoint) {
		t.Helper()

		require.Equal(t, req.SystemId, resp.SystemId)
		if req.Interface.Id != nil {
			require.NotNil(t, resp.Interface)
			require.Equal(t, *req.Interface.Id, *resp.Interface.Id)
		}
		require.Equal(t, req.Interface.Data.IfName, resp.Interface.Data.IfName)
		require.Equal(t, req.Interface.Data.TransformationId, resp.Interface.Data.TransformationId)
		if req.Interface.Data.Ipv4Address != nil && resp.Interface.Data.Ipv4Address != nil {
			require.Equal(t, req.Interface.Data.Ipv4Address.String(), resp.Interface.Data.Ipv4Address.String())
		} else {
			require.Nil(t, req.Interface.Data.Ipv4Address)
			require.Nil(t, resp.Interface.Data.Ipv4Address)
		}
		if req.Interface.Data.Ipv6Address != nil && resp.Interface.Data.Ipv6Address != nil {
			require.Equal(t, req.Interface.Data.Ipv6Address.String(), resp.Interface.Data.Ipv6Address.String())
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
		require.Equal(t, req.Label, resp.Label)
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
		ffc, sysIds := testFFBlueprintB(ctx, t, client.client, 2)

		testCases := map[string]testCase{
			"start_with_minimal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEndpoint{
							{SystemId: sysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/0",
								TransformationId: 1,
							}}},
							{SystemId: sysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/0",
								TransformationId: 1,
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEndpoint{
							{SystemId: sysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/1",
								TransformationId: 2,
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.0.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8::3"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: sysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/1",
								TransformationId: 2,
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.0.0.4"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8::4"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEndpoint{
							{SystemId: sysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/2",
								TransformationId: 1,
							}}},
							{SystemId: sysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/2",
								TransformationId: 1,
							}}},
						},
					},
				},
			},
			"start_with_maximal_config": {
				steps: []FreeformLinkRequest{
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEndpoint{
							{SystemId: sysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/3",
								TransformationId: 2,
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.1"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::1"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: sysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/3",
								TransformationId: 2,
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Endpoints: [2]FreeformEndpoint{
							{SystemId: sysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/4",
								TransformationId: 1,
							}}},
							{SystemId: sysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/4",
								TransformationId: 1,
							}}},
						},
					},
					{
						Label: randString(6, "hex"),
						Tags:  randStrings(rand.Intn(3)+2, 6),
						Endpoints: [2]FreeformEndpoint{
							{SystemId: sysIds[0], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/5",
								TransformationId: 2,
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.2"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::2"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
							}}},
							{SystemId: sysIds[1], Interface: FreeformInterface{Data: &FreeformInterfaceData{
								IfName:           "ge-0/0/5",
								TransformationId: 2,
								Ipv4Address:      &net.IPNet{IP: net.ParseIP("10.1.0.3"), Mask: net.CIDRMask(24, 32)},
								Ipv6Address:      &net.IPNet{IP: net.ParseIP("2001:db8:1::3"), Mask: net.CIDRMask(64, 128)},
								Tags:             randStrings(rand.Intn(3)+2, 6),
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
				require.Equal(t, id, link.Id)
				compare(t, &tCase.steps[0], link.Data)

				// update the link once for each "step", including the first step (values used at creation)
				for _, step := range tCase.steps {
					// record the interface ID in the update request
					step.Endpoints[0].Interface.Id = link.Data.Endpoints[0].Interface.Id
					step.Endpoints[1].Interface.Id = link.Data.Endpoints[1].Interface.Id

					// update the link
					require.NoError(t, ffc.UpdateLink(ctx, id, &step))

					// read the link
					link, err = ffc.GetLink(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, link.Id)
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
