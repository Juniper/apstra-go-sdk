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
	ps := &PropertySetData{
		Label:  randString(10, "hex"),
		Values: make(map[string]string, samples),
	}
	for i := 0; i < samples; i++ {
		ps.Values["_"+randString(10, "hex")] = randString(10, "hex")
	}

	for clientName, client := range clients {
		log.Printf("testing createPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id1, err := client.client.createPropertySet(ctx, ps)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getPropertySetsByLabel() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		psSlice, err := client.client.GetPropertySetsByLabel(ctx, ps.Label)
		if err != nil {
			t.Fatal(err)
		}
		if len(psSlice) != 1 {
			t.Fatalf("expected 1 property set, got %d", len(psSlice))
		}

		log.Printf("testing getPropertySetByLabel() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = client.client.getPropertySetByLabel(ctx, ps.Label)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ps1, err := client.client.getPropertySet(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}
		polished1, err := ps1.polish()
		if err != nil {
			t.Fatal(err)
		}
		if !equal(ps, polished1.Data) {
			t.Fatal("ps and ps1 are not equal")
		}

		ps.Label = randString(10, "hex")
		for i := 0; i < samples; i++ {
			ps.Values["_"+randString(10, "hex")] = randString(10, "hex")
		}

		log.Printf("testing updatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.updatePropertySet(ctx, id1, ps)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ps2, err := client.client.getPropertySet(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}
		polished2, err := ps2.polish()
		if err != nil {
			t.Fatal(err)
		}
		if !equal(ps, polished2.Data) {
			t.Fatal("ps and ps2 are not equal")
		}

		log.Printf("testing deletePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deletePropertySet(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = client.client.getPropertySet(ctx, id1)
		if err == nil {
			t.Fatal("expected a 404 here, didn't get one")
		}
	}

}
