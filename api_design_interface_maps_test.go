package goapstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestInterfaceSettingParam(t *testing.T) {
	expected := `{\"global\":{\"breakout\":false,\"fpc\":0,\"pic\":0,\"port\":0,\"speed\":\"100g\"},\"interface\":{\"speed\":\"\"}}`
	test := InterfaceSettingParam{
		Global: struct {
			Breakout bool   `json:"breakout"`
			Fpc      int    `json:"fpc"`
			Pic      int    `json:"pic"`
			Port     int    `json:"port"`
			Speed    string `json:"speed"`
		}{
			Breakout: false,
			Fpc:      0,
			Pic:      0,
			Port:     0,
			Speed:    "100g",
		},
		Interface: struct {
			Speed string `json:"speed"`
		}{},
	}
	result := test.String()
	if result != expected {
		t.Fatalf("expected '%s', got '%s'", expected, result)
	}
}

func TestListGetAllInterfaceMaps(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	iMapIds, err := client.listAllInterfaceMapIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(iMapIds) == 0 {
		t.Fatal("we should have gotten some interface maps here")
	}

	log.Println("all interface maps IDs: ", iMapIds)

	iMap, err := client.GetInterfaceMap(context.TODO(), iMapIds[rand.Intn(len(iMapIds))])
	if err != nil {
		t.Fatal(err)
	}
	log.Println("random interface map: ", iMap)
}
