//go:build integration

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListAndGetSampleDeviceProfiles(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing listDeviceProfileIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ids, err := client.client.listDeviceProfileIds(ctx)
			require.NoError(t, err)
			require.Greater(t, len(ids), 0)

			for _, i := range samples(t, len(ids), 20) {
				id := ids[i]
				t.Run(fmt.Sprintf("GET_Device_Profile_ID_%s", id), func(t *testing.T) {
					t.Parallel()

					log.Printf("testing getDeviceProfile(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
					dp, err := client.client.getDeviceProfile(ctx, id)
					require.NoError(t, err)
					log.Printf("device profile id '%s' label '%s'\n", id, dp.Label)
				})
			}

			log.Printf("testing GetAllDeviceProfiles() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			profiles, err := client.client.GetAllDeviceProfiles(ctx)
			require.NoError(t, err)

			log.Printf("list found %d, GetAll found %d", len(ids), len(profiles))
			require.Equal(t, len(ids), len(profiles))

			for _, i := range samples(t, len(profiles), 5) {
				label := profiles[i].Data.Label
				t.Run(fmt.Sprintf("GET_DeviceProfile_Label_%s", label), func(t *testing.T) {
					//t.Parallel() // this seems to make things worse

					dp, err := client.client.GetDeviceProfileByName(ctx, label)
					require.NoError(t, err)
					require.Equal(t, label, dp.Data.Label)

					log.Printf("device profile id '%s' label '%s'\n", dp.Id, dp.Data.Label)
				})
			}
		})
	}
}

func TestGetDeviceProfile(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	desiredId := ObjectId("Cisco_3172PQ_NXOS")
	desiredLabel := "Cisco 3172PQ"

	for clientName, client := range clients {
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing GetDeviceProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			dp, err := client.client.GetDeviceProfile(ctx, desiredId)
			require.NoError(t, err)
			require.Equal(t, dp.Data.Label, desiredLabel)
		})
	}
}

func TestGetDeviceProfileByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	desiredLabel := "Cisco 3172PQ"
	desiredId := ObjectId("Cisco_3172PQ_NXOS")

	for clientName, client := range clients {
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing GetDeviceProfileByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			dp, err := client.client.GetDeviceProfileByName(ctx, desiredLabel)
			require.NoError(t, err)
			require.Equal(t, desiredId, dp.Id)
		})
	}
}

func TestGetTransformCandidates(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	dpId := ObjectId("Juniper_QFX5120-48T_Junos")
	intfName := "et-0/0/48"
	intfSpeed := LogicalDevicePortSpeed("40G")

	for clientName, client := range clients {
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing GetDeviceProfileByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			dp, err := client.client.GetDeviceProfile(context.Background(), dpId)
			require.NoError(t, err)

			portInfo, err := dp.Data.PortByInterfaceName(intfName)
			require.NoError(t, err)

			candidates := portInfo.TransformationCandidates(intfName, intfSpeed)
			for k, v := range candidates {
				dump, err := json.MarshalIndent(&v, "", "  ")
				require.NoError(t, err)

				log.Printf("port %d (%s) transformations: %s", k, intfName, string(dump))
			}
		})
	}
}
