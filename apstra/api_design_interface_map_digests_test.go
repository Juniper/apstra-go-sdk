//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestGetInterfaceMapDigest(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getInterfaceMapDigests() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		allImd, err := client.client.getInterfaceMapDigests(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		randId := allImd[rand.Intn(len(allImd))].Id

		log.Printf("testing getInterfaceMapDigest('%s') against %s %s (%s)", randId, client.clientType, clientName, client.client.ApiVersion())
		imd, err := client.client.getInterfaceMapDigest(context.Background(), randId)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s: %s -> %s", imd.Label, imd.LogicalDevice.Label, imd.DeviceProfile.Label)
	}
}

func TestGetInterfaceMapDigestsByLogicalDevice(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getAllLogicalDevices() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		allDp, err := client.client.getAllLogicalDevices(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		randId := allDp[rand.Intn(len(allDp))].Id
		log.Printf("testing getInterfaceMapDigestsByLogicalDevice(%s) against %s %s (%s)", randId, client.clientType, clientName, client.client.ApiVersion())
		imds, err := client.client.getInterfaceMapDigestsByLogicalDevice(context.Background(), randId)
		if err != nil {
			t.Fatal(err)
		}
		for _, imd := range imds {
			log.Printf("%s: %s -> %s", imd.Label, imd.LogicalDevice.Label, imd.DeviceProfile.Label)
		}
	}
}

func TestGetInterfaceMapDigestsByDeviceProfile(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getAllDeviceProfiles() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		allDp, err := client.client.getAllDeviceProfiles(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		randId := allDp[rand.Intn(len(allDp))].Id
		log.Printf("testing getInterfaceMapDigestsByDeviceProfile(%s) against %s %s (%s)", randId, client.clientType, clientName, client.client.ApiVersion())
		imds, err := client.client.getInterfaceMapDigestsByDeviceProfile(context.Background(), randId)
		if err != nil {
			t.Fatal(err)
		}
		for _, imd := range imds {
			log.Printf("%s: %s -> %s", imd.Label, imd.LogicalDevice.Label, imd.DeviceProfile.Label)
		}
	}
}
