package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestListAllBlueprintIds(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		blueprints, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		result, err := json.Marshal(blueprints)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(string(result))
	}
}

func TestGetAllBlueprintStatus(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing getAllBlueprintStatus() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		bps, err := client.client.getAllBlueprintStatus(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Println(len(bps))
	}
}

func TestCreateDeleteBlueprint(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing createBlueprintFromTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		name := randString(10, "hex")
		id, err := client.client.createBlueprintFromTemplate(context.TODO(), &CreateBlueprintFromTemplate{
			RefDesign:  RefDesignDatacenter,
			Label:      name,
			TemplateId: "L2_Virtual_EVPN",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("got id '%s', deleting blueprint...\n", id)
		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetPatchGetPatchNode(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		bpIds, err := client.client.listAllBlueprintIds(context.TODO())
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
		log.Printf("testing getNodes() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.getNodes(context.TODO(), bpIds[0], NodeTypeMetadata, &nodesA)
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
			log.Printf("testing patchNode() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			err := client.client.patchNode(context.TODO(), bpIds[0], nodeA.Id, req, resp)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("response indicates name changed '%s' -> '%s'", nodeA.Label, resp.Label)

			log.Printf("testing getNodes() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			err = client.client.getNodes(context.TODO(), bpIds[0], NodeTypeMetadata, &nodesB)
			if err != nil {
				t.Fatal()
			}
			for idB, nodeB := range nodesB.Nodes {
				log.Printf("node id: %s ; label: %s\n", idB, nodeB.Label)

			}

			req = metadataNode{Label: nodeA.Label}
			log.Printf("testing patchNode() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			err = client.client.patchNode(context.TODO(), bpIds[0], nodeA.Id, req, resp)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("response indicates name changed '%s' -> '%s'", newName, resp.Label)
		}
	}
}

func TestGetLockInfo(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing createBlueprintFromTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		name := randString(10, "hex")
		id, err := client.client.createBlueprintFromTemplate(context.TODO(), &CreateBlueprintFromTemplate{
			RefDesign:  RefDesignDatacenter,
			Label:      name,
			TemplateId: "L2_Virtual_EVPN",
		})
		if err != nil {
			t.Fatal(err)
		}

		l, err := client.client.getBlueprintLockInfo(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(l)
		log.Printf("got id '%s', deleting blueprint...\n", id)
		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
