//go:build integration
// +build integration

package apstra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCRUDFFLink(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	compareEndPoint := func(t testing.TB, a, b *FreeformEndpoint) {
		t.Helper()

		require.Equal(t, a.SystemId, b.SystemId)
		require.Equal(t, a.Interface.IfName, b.Interface.IfName)
		require.Equal(t, a.Interface.TransformationId, b.Interface.TransformationId)
		if a.Interface.Ipv4Address != nil && b.Interface.Ipv4Address != nil {
			require.Equal(t, a.Interface.Ipv4Address.String(), b.Interface.Ipv4Address.String())
		} else {
			require.Nil(t, a.Interface.Ipv4Address)
			require.Nil(t, b.Interface.Ipv4Address)
		}
		if a.Interface.Ipv6Address != nil && b.Interface.Ipv6Address != nil {
			require.Equal(t, a.Interface.Ipv6Address.String(), b.Interface.Ipv6Address.String())
		} else {
			require.Nil(t, a.Interface.Ipv6Address)
			require.Nil(t, b.Interface.Ipv6Address)
		}
		compareSlicesAsSets(t, a.Interface.Tags, b.Interface.Tags, "tag mismatch")
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

	for _, client := range clients {
		ffc, sysIds := testFFBlueprintB(ctx, t, client.client, 2)

		req := FreeformLinkRequest{
			Label: randString(6, "hex"),
			Tags:  []ObjectId{"a", "b"},
			Endpoints: [2]FreeformEndpoint{
				{
					SystemId: sysIds[0],
					Interface: FreeformInterfaceData{
						IfName:           "ge-0/0/0",
						TransformationId: 1,
						Ipv4Address:      nil,
						Ipv6Address:      nil,
						Tags:             nil,
					},
				},
				{
					SystemId: sysIds[1],
					Interface: FreeformInterfaceData{
						IfName:           "ge-0/0/0",
						TransformationId: 1,
						Ipv4Address:      nil,
						Ipv6Address:      nil,
						Tags:             nil,
					},
				},
			},
		}

		// create the link
		id, err := ffc.CreateLink(ctx, &req)
		require.NoError(t, err)

		// now lets read the link
		readLink, err := ffc.GetLink(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, readLink.Id)
		compare(t, &req, readLink.Data)

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
	}
}
