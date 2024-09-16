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

	"github.com/hashicorp/go-version"
)

func TestCreateGetUpdateGetDeletePropertySet(t *testing.T) {
	ctx := context.Background()

	equalFunc := func(a, b *PropertySetData) bool {
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
		return jsonEqual(t, a.Values, b.Values)
	}

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	testData := PropertySetData{
		Label: randString(10, "hex"),
	}

	sampleCount := rand.Intn(10) + 3
	vals := make(map[string]string, sampleCount)
	for i := 0; i < sampleCount; i++ {
		vals["_"+randString(10, "hex")] = randString(10, "hex")
	}
	testData.Values, err = json.Marshal(vals)
	if err != nil {
		t.Fatal(err)
	}

	nestedJsonMinVer, err := version.NewVersion("4.1.2")
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		psData := testData // start with clean copy of psData in each loop
		t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
			t.Parallel()

			apiVersionString, err := version.NewVersion(client.client.apiVersion.String())
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing CreatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			id1, err := client.client.CreatePropertySet(ctx, &psData)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetPropertySetByLabel() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ps, err := client.client.GetPropertySetByLabel(ctx, psData.Label)
			if err != nil {
				t.Fatal(err)
			}

			if !equalFunc(&psData, ps.Data) {
				t.Fatal("Created and fetched objects are not equal")
			}

			psData.Label = randString(10, "hex")
			for i := 0; i < sampleCount; i++ {
				vals["_"+randString(10, "hex")] = randString(10, "hex")
			}

			psData.Values, err = json.Marshal(vals)
			if err != nil {
				t.Fatal(err)
			}

			// nested JSON only supported by Apstra 4.1.2 and later
			if apiVersionString.GreaterThanOrEqual(nestedJsonMinVer) {
				psData.Values = append(psData.Values[:len(psData.Values)-1], []byte(`,"inner_json":{"number":1, "string":"str1"}}`)...)
			}

			log.Printf("testing UpdatePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.UpdatePropertySet(ctx, id1, &psData)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ps, err = client.client.GetPropertySet(ctx, id1)
			if err != nil {
				t.Fatal(err)
			}
			if !equalFunc(&psData, ps.Data) {
				t.Fatal("psData and ps.Data are not equal")
			}

			log.Printf("testingt  DeletePropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.DeletePropertySet(ctx, id1)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("Testing GetPropertySet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			_, err = client.client.GetPropertySet(ctx, id1)
			if err == nil {
				t.Fatal("Fetching a property set after deletion should have produced an error")
			}
		})
	}
}
