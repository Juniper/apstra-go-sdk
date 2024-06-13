package apstra

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateDeleteFFLink(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)
	for clientName, client := range clients {
		ffc, sysIds := testFFBlueprintA(ctx, t, client.client)
		_ = clientName
		linkId, err := ffc.CreateLink(ctx, &FreeformLinkRequest{
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
		})
		require.NoError(t, err)
		t.Log(linkId)
		err = ffc.DeleteLink(ctx, linkId)
		require.NoError(t, err)
	}

}
func TestFFLinkB(t *testing.T) {
	var y ConfigTemplate
	rawjson := []byte(`{"id":"foo","link_type":"1","label":"test_ff_link","speed","10G"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
