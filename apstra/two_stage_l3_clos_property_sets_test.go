//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"testing"
)

func TestImportGetUpdateGetDeletePropertySet(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	for clientName, client := range clients {

		// Create Blueprint
		bpClient := testBlueprintC(ctx, t, client.client)

		// Create Property Set
		samples := rand.Intn(10) + 4
		ps := &PropertySetData{
			Label: randString(10, "hex"),
		}
		vals := make(map[string]string, samples)

		for i := 0; i < samples; i++ {
			vals["_"+randString(10, "hex")] = randString(10, "hex")
		}
		ps.Values, err = json.Marshal(vals)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("Create Property Set on %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ps_id, err := client.client.CreatePropertySet(ctx, ps)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing ImportPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips_id, err := bpClient.ImportPropertySet(ctx, ps_id)
		log.Printf("%s", ips_id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetImportedPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips, err := bpClient.GetPropertySet(ctx, ips_id)
		if err != nil {
			t.Fatal(err)
		}

		p := make(map[string]interface{})
		err = json.Unmarshal([]byte(ips.Values), &p)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("Ensure you imported the right set of key/value pairs")
		if !jsonEqual(t, ips.Values, ps.Values) {
			t.Fatalf("Import Mismatch. Expected %v Got %v", vals, p)
		}
		log.Printf("%v", p)
		keys := getKeysfromMap(p)
		log.Printf("%v", keys)
		log.Printf("testing UpdateImportedPropertySet() with a subset of keys against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.UpdatePropertySet(ctx, ps_id, keys[2:]...)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetImportedPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips2, err := bpClient.GetPropertySet(ctx, ps_id)
		if err != nil {
			t.Fatal(err)
		}
		if !ips2.Stale {
			t.Fatal("The imported property set must show as stale")
		}
		p2 := make(map[string]interface{})
		err = json.Unmarshal([]byte(ips2.Values), &p2)
		log.Printf("%v", p2)
		if err != nil {
			t.Fatal(err)
		}
		k2 := getKeysfromMap(p2)
		log.Println(k2)
		log.Println(len(keys))
		log.Println(len(k2))
		if len(k2)+2 != len(keys) {
			t.Fatalf("Subset of keys not imported. Should have imported %v imported %v", keys[2:], k2)
		}

		log.Printf("testing DeleteImportedPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeletePropertySet(ctx, ps_id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
