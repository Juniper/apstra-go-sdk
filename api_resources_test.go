package goapstra

import (
	"context"
	"log"
	"math"
	"math/rand"
	"testing"
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
	var openHoles []NewAsnRange
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
		openHoles, err = invertRangesInRange(1, math.MaxUint32, poolBeginEnds)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("open holes in ASN resources: ", openHoles)

		// todo: make sure there's at least one open hole in the plan
		name := "test-" + randString(10, "hex")
		r := rand.Intn(len(openHoles))
		id, err := client.CreateAsnPool(context.TODO(), &NewAsnPool{
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

func TestEmptyAsnPool(t *testing.T) {
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
		id, err := client.CreateAsnPool(context.TODO(), &NewAsnPool{DisplayName: name})
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
