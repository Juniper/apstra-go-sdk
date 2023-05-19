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

func CreateBluePrint(t *testing.T, client testClient, clientName string) ObjectId {
	// Create Blueprint
	log.Printf("testing createBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
	bp_name := randString(10, "hex")
	bp_id, err := client.client.CreateBlueprintFromTemplate(context.TODO(), &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignDatacenter,
		Label:      bp_name,
		TemplateId: "L2_Virtual_EVPN",
	})
	if err != nil {
		t.Fatal(err)
	}
	return bp_id
}

func TestImportGetUpdateGetDeletePropertySet(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	for clientName, client := range clients {

		// Create Blueprint
		bp_id := CreateBluePrint(t, client, clientName)
		// Create Property Set
		samples := rand.Intn(10) + 4
		ps := &PropertySetData{
			Label: randString(10, "hex"),
		}
		var t TwoStageL3ClosClient
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
		ips_id, err := client.client.ImportPropertySet(ctx, bp_id, ps_id)
		log.Printf("%s", ips_id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetImportedPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips, err := client.client.GetImportedPropertySet(ctx, bp_id, ips_id)
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
		keys := stringkeysfromMap(p)
		log.Printf("%v", keys)
		log.Printf("testing UpdateImportedPropertySet() with a subset of keys against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.UpdateImportedPropertySet(ctx, bp_id, ps_id, keys[2:]...)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetImportedPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips2, err := client.client.GetImportedPropertySet(ctx, bp_id, ps_id)
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
		k2 := stringkeysfromMap(p2)
		log.Println(k2)
		log.Println(len(keys))
		log.Println(len(k2))
		if len(k2)+2 != len(keys) {
			t.Fatalf("Subset of keys not imported. Should have imported %v imported %v", keys[2:], k2)
		}

		log.Printf("testing DeleteImportedPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteImportedPropertySet(ctx, bp_id, ps_id)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			log.Printf("testing DeletePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.DeletePropertySet(ctx, ps_id)
			if err != nil {
				log.Println(err)
			}
			log.Printf("got id '%s', deleting blueprint...\n", bp_id)
			log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.deleteBlueprint(context.TODO(), bp_id)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
