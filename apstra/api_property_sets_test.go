//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestCreateGetUpdateGetDeletePropertySet(t *testing.T) {
	equal := func(a, b *PropertySetData) bool {
		if a.Label != b.Label {
			return false
		}
		if len(a.Blueprints) != len(b.Blueprints) {
			return false
		}
		for i := range a.Blueprints {
			if a.Blueprints[i] != b.Blueprints[i] {
				return false
			}
		}
		if len(a.Values) != len(b.Values) {
			return false
		}
		for k := range a.Values {
			if _, ok := b.Values[k]; !ok {
				return false
			}
			if a.Values[k] != b.Values[k] {
				return false
			}
		}
		return true
	}

	ctx := context.Background()
	clients, err := getTestClients(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	samples := rand.Intn(10) + 3
	ps := &PropertySetRequest{
		Label:  randString(10, "hex"),
		Values: make(map[string]string, samples),
	}
	for i := 0; i < samples; i++ {
		ps.Values["_"+randString(10, "hex")] = randString(10, "hex")
	}

	for clientName, client := range clients {
		log.Printf("testing CreatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id1, err := client.client.CreatePropertySet(ctx, ps)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("Testing duplicate CreatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = client.client.CreatePropertySet(ctx, ps)
		if err == nil {
			t.Fatal("Expected to fail duplicate")
		}
		log.Printf("Error Received %q", err)

		log.Printf("testing GetPropertySetByLabel() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ps1, err := client.client.GetPropertySetByLabel(ctx, ps.Label)
		if err != nil {
			t.Fatalf("Error %q", err)
		}
		if !equal((*PropertySetData)(ps), ps1.Data) {
			t.Fatal("Created and extracted objects are not equal")
		}

		ps.Label = randString(10, "hex")
		for i := 0; i < samples; i++ {
			ps.Values["_"+randString(10, "hex")] = randString(10, "hex")
		}

		log.Printf("testing updatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.UpdatePropertySet(ctx, id1, ps)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ps2, err := client.client.GetPropertySet(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}
		if !equal((*PropertySetData)(ps), ps2.Data) {
			t.Fatal("ps and ps2 are not equal")
		}

		log.Printf("testing deletePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeletePropertySet(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}

	}
}
