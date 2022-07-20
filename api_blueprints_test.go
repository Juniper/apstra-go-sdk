package goapstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"testing"
)

func blueprintsTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestListAllBlueprintIds(t *testing.T) {
	client, err := blueprintsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	blueprints, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	result, err := json.Marshal(blueprints)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(result))
}

func TestGetAllBlueprintStatus(t *testing.T) {
	client, err := blueprintsTestClient1()
	if err != nil {
		log.Fatal(err)
	}

	bps, err := client.getAllBlueprintStatus(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	log.Println(len(bps))
}

func TestCreateDeleteBlueprint(t *testing.T) {
	client, err := blueprintsTestClient1()
	if err != nil {
		log.Fatal(err)
	}

	client.Login(context.TODO())
	name := randString(10, "hex")
	id, err := client.createBlueprintFromTemplate(context.TODO(), &CreateBluePrintFromTemplate{
		RefDesign:  RefDesignDatacenter,
		Label:      name,
		TemplateId: "L2_Virtual_EVPN",
	})

	log.Printf("got id '%s'\n", id)

	err = client.deleteBlueprint(context.TODO(), id)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateDeleteRoutingZone(t *testing.T) {
	DebugLevel = 4
	client, err := blueprintsTestClient1()
	if err != nil {
		log.Fatal(err)
	}

	blueprints, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	if len(blueprints) < 1 {
		t.Skipf("cannot proceed without at least one blueprint")
	}

	blueprintId := blueprints[0]
	randString := randString(10, "hex")

	zoneId, err := client.CreateRoutingZone(context.TODO(), blueprintId, &CreateRoutingZoneCfg{
		SzType:  "evpn",
		VrfName: "test-" + randString,
		Label:   "label-test-" + randString,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("created zone '", zoneId, "' deleting...")

	err = client.DeleteRoutingZone(context.TODO(), blueprintId, zoneId)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetNodes(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	bpIds, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(bpIds) == 0 {
		t.Skip("no blueprints? no nodes.")
	}

	nodes := struct {
		Nodes map[string]struct {
			Tags         interface{} `json:"tags"`
			PropertySet  interface{} `json:"property_set"`
			Label        string      `json:"label"`
			UserIp       interface{} `json:"user_ip"`
			TemplateJson interface{} `json:"template_json"`
			Design       string      `json:"design"`
			User         interface{} `json:"user"`
			Type         string      `json:"type"`
			Id           string      `json:"id"`
		} `json:"nodes"`
	}{}
	err = client.getNodes(context.TODO(), bpIds[0], NodeTypeMetadata, &nodes)
	if err != nil {
		t.Fatal()
	}
	for id, node := range nodes.Nodes {
		log.Printf("node id: %s ; label: %s\n", id, node.Label)
	}
}
