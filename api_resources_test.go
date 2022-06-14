package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"testing"
)

const (
	testPool1 = `{
      "status": "not_in_use",
      "used": "0",
      "display_name": "foo",
      "tags": [],
      "created_at": "2022-06-13T18:44:55.899107Z",
      "last_modified_at": "2022-06-13T18:44:55.899107Z",
      "ranges": [],
      "used_percentage": 0,
      "total": "0",
      "id": "e49e0f45-ecf3-478d-8b1f-901af6d4ed89"
    }`
)

// todo asnpool mocking
func TestGetCreateDeleteAsnPools(t *testing.T) {
	DebugLevel = 2
	clients, apis, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}

	_, mockExists := apis["mock"]
	if mockExists {
		err = apis["mock"].createMetricdb()
		if err != nil {
			log.Fatal(err)
		}
	}
	for clientName, client := range clients {
		log.Printf("testing GetAsnPools() with %s client", clientName)
		pools, err := client.GetAsnPools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pools)
		var poolBeginEnds []NewAsnRange
		for _, p := range pools {
			for _, r := range p.Ranges {
				poolBeginEnds = append(poolBeginEnds, NewAsnRange{r.First, r.Last})
			}
		}
		openHoles, err := invertRangesInRange(1, math.MaxUint32, poolBeginEnds)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("open holes in ASN resources: ", openHoles)

		// todo: make sure there's at least one open hole in the plan
		name := "test-" + randString(10, "hex")
		r := rand.Intn(len(openHoles))
		id, err := client.CreateAsnPool(context.TODO(), &NewAsnPoolCfg{
			Ranges: []NewAsnRange{{
				B: openHoles[r].B,
				E: openHoles[r].E,
			}},
			DisplayName: name,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, id)

		_, err = client.GetAsnPool(context.TODO(), id)

		err = client.DeleteAsnPool(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUpdateEmptyAsnPool(t *testing.T) {
	DebugLevel = 4
	clients, apis, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}

	_, mockExists := apis["mock"]
	if mockExists {
		err = apis["mock"].createMetricdb()
		if err != nil {
			log.Fatal(err)
		}
	}

	name := "test-" + randString(10, "hex")

	for clientName, client := range clients {

		log.Printf("creating empty ASN pool '%s' with %s client", name, clientName)
		newPoolId, err := client.CreateAsnPool(context.TODO(), &NewAsnPoolCfg{DisplayName: name})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, newPoolId)

		pools, err := client.GetAsnPools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pools)
		var poolBeginEnds []NewAsnRange
		for _, p := range pools {
			for _, r := range p.Ranges {
				poolBeginEnds = append(poolBeginEnds, NewAsnRange{r.First, r.Last})
			}
		}
		openHoles, err := invertRangesInRange(1, math.MaxUint32, poolBeginEnds)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("open holes in ASN resources: ", openHoles)

		// todo: make sure there's at least one open hole in the plan
		r := rand.Intn(len(openHoles))
		newRange := NewAsnRange{
			B: openHoles[r].B,
			E: openHoles[r].E,
		}
		newDisplayName := "updated-" + name
		newTags := []string{"updated"}
		err = client.updateAsnPool(context.TODO(), newPoolId, &NewAsnPoolCfg{
			DisplayName: newDisplayName,
			Ranges:      []NewAsnRange{newRange},
			Tags:        newTags,
		})
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.GetAsnPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}

		err = client.DeleteAsnPool(context.TODO(), newPoolId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUnmarshalPool(t *testing.T) {
	var result rawAsnPool
	err := json.Unmarshal([]byte(testPool1), &result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAsnPoolRangeId(t *testing.T) {
	testRange := AsnRange{
		First: 66051,    // 0x00 0x01 0x02 0x03,
		Last:  67438087, // 0x04 0x05 0x06 0x07,
	}
	expected := "pool-8a851ff8"
	result := getAsnPoolRangeId(&testRange)
	if result != expected {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
