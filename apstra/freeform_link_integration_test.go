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
		require.Equal(t, a.Interface.Data.IfName, b.Interface.Data.IfName)
		require.Equal(t, a.Interface.Data.TransformationId, b.Interface.Data.TransformationId)
		if a.Interface.Data.Ipv4Address != nil && b.Interface.Data.Ipv4Address != nil {
			require.Equal(t, a.Interface.Data.Ipv4Address.String(), b.Interface.Data.Ipv4Address.String())
		} else {
			require.Nil(t, a.Interface.Data.Ipv4Address)
			require.Nil(t, b.Interface.Data.Ipv4Address)
		}
		if a.Interface.Data.Ipv6Address != nil && b.Interface.Data.Ipv6Address != nil {
			require.Equal(t, a.Interface.Data.Ipv6Address.String(), b.Interface.Data.Ipv6Address.String())
		} else {
			require.Nil(t, a.Interface.Data.Ipv6Address)
			require.Nil(t, b.Interface.Data.Ipv6Address)
		}
		compareSlicesAsSets(t, a.Interface.Data.Tags, b.Interface.Data.Tags, "tag mismatch")
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
			Tags:  []string{"a", "b"},
			Endpoints: [2]FreeformEndpoint{
				{
					SystemId: sysIds[0],
					Interface: FreeformInterface{
						Id: nil,
						Data: &FreeformInterfaceData{
							IfName:           "ge-0/0/0",
							TransformationId: 1,
							Ipv4Address:      nil,
							Ipv6Address:      nil,
							Tags:             nil,
						},
					},
				},
				{
					SystemId: sysIds[1],
					Interface: FreeformInterface{
						Id: nil,
						Data: &FreeformInterfaceData{
							IfName:           "ge-0/0/0",
							TransformationId: 1,
							Ipv4Address:      nil,
							Ipv6Address:      nil,
							Tags:             nil,
						},
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

		// lets see if we can update the link
		// first put the read link endpoint interface ID's into the req data structure.
		req.Endpoints[0].Interface.Id = readLink.Data.Endpoints[0].Interface.Id
		req.Endpoints[1].Interface.Id = readLink.Data.Endpoints[1].Interface.Id
		// Now change the Label and the tags.
		req.Label = randString(7, "hex")
		req.Tags = []string{"a", "b", "c"}
		// update the link.
		err = ffc.UpdateLink(ctx, id, &req)
		require.NoError(t, err)
		// read the link back.
		readLink, err = ffc.GetLink(ctx, id)
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
