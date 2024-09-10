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
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing listDeviceProfileIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ids, err := client.client.listDeviceProfileIds(ctx)
			require.NoError(t, err)

			for _, id := range ids {
				log.Printf("testing getDeviceProfile(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
				dp, err := client.client.getDeviceProfile(ctx, id)
				require.NoError(t, err)
				log.Printf("device profile id '%s' label '%s'\n", id, dp.Label)
			}

			log.Printf("testing getAllDeviceProfiles() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			profiles, err := client.client.getAllDeviceProfiles(ctx)
			require.NoError(t, err)
			log.Printf("list found %d, getAll found %d", len(ids), len(profiles))
		})
	}
}

func TestGetDeviceProfile(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	desiredId := ObjectId("Cisco_3172PQ_NXOS")
	desiredLabel := "Cisco 3172PQ"

	for clientName, client := range clients {
		log.Printf("testing GetDeviceProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		dp, err := client.client.GetDeviceProfile(context.Background(), desiredId)
		if err != nil {
			t.Fatal(err)
		}
		if dp.Data.Label != desiredLabel {
			t.Fatalf("expected '%s', got '%s'", desiredLabel, dp.Data.Label)
		}
	}
}

func TestGetDeviceProfileByName(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	desiredLabel := "Cisco 3172PQ"
	desiredId := ObjectId("Cisco_3172PQ_NXOS")

	for clientName, client := range clients {
		log.Printf("testing GetDeviceProfileByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		dp, err := client.client.GetDeviceProfileByName(context.Background(), desiredLabel)
		if err != nil {
			t.Fatal(err)
		}
		if dp.Id != desiredId {
			t.Fatalf("expected '%s', got '%s'", desiredId, dp.Id)
		}
	}
}

func TestGetTransformCandidates(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	dpId := ObjectId("Juniper_QFX5120-48T_Junos")
	intfName := "et-0/0/48"
	intfSpeed := LogicalDevicePortSpeed("40G")

	for clientName, client := range clients {
		log.Printf("testing GetDeviceProfileByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		dp, err := client.client.GetDeviceProfile(context.Background(), dpId)
		if err != nil {
			t.Fatal(err)
		}

		portInfo, err := dp.Data.PortByInterfaceName(intfName)
		if err != nil {
			t.Fatal(err)
		}

		candidates := portInfo.TransformationCandidates(intfName, intfSpeed)
		for k, v := range candidates {
			dump, err := json.MarshalIndent(&v, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("port %d (%s) transformations: %s", k, intfName, string(dump))
		}
	}
}
