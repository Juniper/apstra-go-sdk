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

func TestGetPatchGetPatchNode(t *testing.T) {
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

	type metadataNode struct {
		Tags         interface{} `json:"tags,omitempty"`
		PropertySet  interface{} `json:"property_set,omitempty"`
		Label        string      `json:"label,omitempty"`
		UserIp       interface{} `json:"user_ip,omitempty"`
		TemplateJson interface{} `json:"template_json,omitempty"`
		Design       string      `json:"design,omitempty"`
		User         interface{} `json:"user,omitempty"`
		Type         string      `json:"type,omitempty"`
		Id           ObjectId    `json:"id,omitempty"`
	}

	nodesA := struct {
		Nodes map[string]metadataNode `json:"nodes"`
	}{}
	nodesB := struct {
		Nodes map[string]metadataNode `json:"nodes"`
	}{}
	err = client.getNodes(context.TODO(), bpIds[0], NodeTypeMetadata, &nodesA)
	if err != nil {
		t.Fatal()
	}

	if len(nodesA.Nodes) != 1 {
		t.Fatalf("not expecting %d '%s' nodes", len(nodesA.Nodes), NodeTypeMetadata)
	}

	newName := randString(10, "hex")
	// loop should run just once (len check above)
	for idA, nodeA := range nodesA.Nodes {
		log.Printf("node id: %s ; label: %s\n", idA, nodeA.Label)

		req := metadataNode{Label: newName}
		resp := &metadataNode{}
		err := client.patchNode(context.TODO(), bpIds[0], nodeA.Id, req, resp)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("response indicates name changed '%s' -> '%s'", nodeA.Label, resp.Label)

		err = client.getNodes(context.TODO(), bpIds[0], NodeTypeMetadata, &nodesB)
		if err != nil {
			t.Fatal()
		}
		for idB, nodeB := range nodesB.Nodes {
			log.Printf("node id: %s ; label: %s\n", idB, nodeB.Label)

		}

		req = metadataNode{Label: nodeA.Label}
		err = client.patchNode(context.TODO(), bpIds[0], nodeA.Id, req, resp)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("response indicates name changed '%s' -> '%s'", newName, resp.Label)
	}

}
