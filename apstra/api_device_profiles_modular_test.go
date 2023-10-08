//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"reflect"
	"testing"
)

func TestModularDeviceProfile(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	mdp1 := &ModularDeviceProfile{
		Label:            randString(5, "hex"),
		ChassisProfileId: "Juniper_PTX10008",
		SlotConfigurations: map[uint64]ModularDeviceSlotConfiguration{
			0: {LinecardProfileId: "Juniper_PTX10K_LC1201_36CD"},
			2: {LinecardProfileId: "Juniper_PTX10K_LC1201_36CD"},
			4: {LinecardProfileId: "Juniper_PTX10K_LC1201_36CD"},
		},
	}

	for clientName, client := range clients {
		log.Printf("testing CreateModularDeviceProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateModularDeviceProfile(ctx, mdp1)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetModularDeviceProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		mdp2, err := client.client.GetModularDeviceProfile(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(*mdp1, *mdp2) {
			t.Fatalf("original and retrieved modular device profiles do not match:\n\n%v\n\n%v", *mdp1, *mdp2)
		}

		mdp1.Label = randString(5, "hex")
		mdp1.ChassisProfileId = "Juniper_PTX10016"
		mdp1.SlotConfigurations = map[uint64]ModularDeviceSlotConfiguration{
			1: {LinecardProfileId: "Juniper_PTX10K_LC1202_36MR"},
			3: {LinecardProfileId: "Juniper_PTX10K_LC1202_36MR"},
			5: {LinecardProfileId: "Juniper_PTX10K_LC1202_36MR"},
		}
		log.Printf("testing UpdateModularDeviceProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.UpdateModularDeviceProfile(ctx, id, mdp1)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetModularDeviceProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		mdp2, err = client.client.GetModularDeviceProfile(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(*mdp1, *mdp2) {
			t.Fatalf("original and retrieved modular device profiles do not match:\n\n%v\n\n%v", *mdp1, *mdp2)
		}

		log.Printf("testing DeleteModularDeviceProfile() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteModularDeviceProfile(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
