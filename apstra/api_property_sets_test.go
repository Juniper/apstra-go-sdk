//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"testing"
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
		eq, err := areEqualJSON(a.Values, b.Values)
		if err != nil {
			t.Fatal(err)
		}
		return eq
	}

	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	samples := rand.Intn(10) + 3
	m := make(map[string]string, samples)
	ps := &PropertySetData{
		Label: randString(10, "hex"),
	}
	for i := 0; i < samples; i++ {
		m["_"+randString(10, "hex")] = randString(10, "hex")
	}
	b, err := json.Marshal(m)
	ps.Values = b
	for clientName, client := range clients {
		log.Printf("testing CreatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())

		id1, err := client.client.CreatePropertySet(ctx, ps)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing GetPropertySetByLabel() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ps1, err := client.client.GetPropertySetByLabel(ctx, ps.Label)
		if err != nil {
			t.Fatal(err)
		}
		if !equal(ps, ps1.Data) {
			t.Fatal("Created and extracted objects are not equal")
		}

		ps.Label = randString(10, "hex")
		for i := 0; i < samples; i++ {
			m["_"+randString(10, "hex")] = randString(10, "hex")
		}
		b, err = json.Marshal(m)
		s := string(b)
		s = s[:len(s)-1] + `,"json_in_json":{"integer":1, "string":"2"}}`
		fmt.Printf(s)
		b = []byte(s)
		ps.Values = b

		log.Printf("testing UpdatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.UpdatePropertySet(ctx, id1, ps)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ps2, err := client.client.GetPropertySet(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}
		if !equal(ps, ps2.Data) {
			t.Fatal("ps and ps2 are not equal")
		}

		log.Printf("testing DeletePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeletePropertySet(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("Testing GetPropertySet() against %s %s (%s). This should fail", client.clientType, clientName, client.client.ApiVersion())
		ps2, err = client.client.GetPropertySet(ctx, id1)
		if err == nil {
			t.Fatal("This Get Should have failed with a 404")
		}
	}
}
