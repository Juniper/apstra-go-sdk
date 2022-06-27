package goapstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
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
		var poolBeginEnds []AsnRange
		for _, p := range pools {
			for _, r := range p.Ranges {
				poolBeginEnds = append(poolBeginEnds, AsnRange{First: r.First, Last: r.Last})
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
		id, err := client.CreateAsnPool(context.TODO(), &AsnPool{
			Ranges: []AsnRange{{
				First: openHoles[r].First,
				Last:  openHoles[r].Last,
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

	name := "test-" + randString(10, "hex")

	for clientName, client := range clients {
		log.Printf("creating empty ASN pool '%s' with %s client", name, clientName)
		newPoolId, err := client.CreateAsnPool(context.TODO(), &AsnPool{DisplayName: name})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, newPoolId)

		pools, err := client.GetAsnPools(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pools)
		var poolBeginEnds []AsnRange
		for _, p := range pools {
			for _, r := range p.Ranges {
				poolBeginEnds = append(poolBeginEnds, AsnRange{First: r.First, Last: r.Last})
			}
		}
		openHoles, err := invertRangesInRange(1, math.MaxUint32, poolBeginEnds)
		if err != nil {
			t.Fatal(err)
		}
		log.Println("open holes in ASN resources: ", openHoles)

		// todo: make sure there's at least one open hole in the plan
		r := rand.Intn(len(openHoles))
		newRange := AsnRange{
			First: openHoles[r].First,
			Last:  openHoles[r].Last,
		}
		newDisplayName := "updated-" + name
		newTags := []string{"updated"}
		err = client.updateAsnPool(context.TODO(), newPoolId, &AsnPool{
			DisplayName: newDisplayName,
			Ranges:      []AsnRange{newRange},
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

func TestGetAsnPoolRangeHash(t *testing.T) {
	testRange := AsnRange{
		First: 66051,    // 0x00 0x01 0x02 0x03,
		Last:  67438087, // 0x04 0x05 0x06 0x07,
	}
	expected := "8a851ff82ee7048ad09ec3847f1ddf44944104d2cbd17ef4e3db22c6785a0d45"
	result := hashAsnPoolRange(&testRange)
	if result != expected {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestCreateDeleteAsnPoolRange(t *testing.T) {
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

	name := "test-" + randString(10, "hex")
	var tags []string
	tags = append(tags, "tag-"+randString(10, "hex"))
	tags = append(tags, "tag-"+randString(10, "hex"))

	for clientName, client := range clients {
		log.Printf("creating empty ASN pool '%s' with %s client", name, clientName)
		poolId, err := client.CreateAsnPool(context.TODO(), &AsnPool{
			DisplayName: name,
			Tags:        tags,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("created ASN pool name %s id %s", name, poolId)
		var asnRange AsnRange
		for i := 0; i < 3; i++ {
			a := rand.Intn(1000) + (i * 1000 * 2)
			b := rand.Intn(1000) + a
			asnRange.First = uint32(a)
			asnRange.Last = uint32(b)
			err = client.CreateAsnPoolRange(context.TODO(), poolId, &asnRange)
			if err != nil {
				t.Fatal(err)
			}
		}

		asnPool, err := client.GetAsnPool(context.TODO(), poolId)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range asnPool.Ranges {
			err := client.DeleteAsnPoolRange(context.TODO(), asnPool.Id, &r)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = client.DeleteAsnPool(context.TODO(), asnPool.Id)
		if err != nil {
			t.Fatal(err)
		}
	}

}
