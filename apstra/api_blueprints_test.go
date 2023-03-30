//go:build integration
// +build integration

package goapstra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestListAllBlueprintIds(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
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
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getAllBlueprintStatus() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bps, err := client.client.getAllBlueprintStatus(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Println(len(bps))
	}
}

func TestCreateDeleteBlueprint(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing createBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		name := randString(10, "hex")
		id, err := client.client.CreateBlueprintFromTemplate(context.TODO(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignDatacenter,
			Label:      name,
			TemplateId: "L2_Virtual_EVPN",
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("got id '%s', deleting blueprint...\n", id)
		log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.deleteBlueprint(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetPatchGetPatchNode(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	testSkip := make(map[string]string)
	for clientName, client := range clients {
		log.Printf("testing listAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.listAllBlueprintIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(bpIds) == 0 {
			testSkip[clientName] = fmt.Sprintf("skipping %s because no blueprints exist", clientName)
			continue
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

	BLUEPRINT:
		for _, id := range bpIds {
			bpStatus, err := client.client.getBlueprintStatus(context.Background(), id)
			if err != nil {
				t.Fatal(err)
			}
			if bpStatus.Design != refDesignDatacenter {
				continue BLUEPRINT
			}
			nodesA := struct {
				Nodes map[string]metadataNode `json:"nodes"`
			}{}
			nodesB := struct {
				Nodes map[string]metadataNode `json:"nodes"`
			}{}

			bp, err := client.client.NewTwoStageL3ClosClient(context.Background(), id)
			if err != nil {
				t.Fatal(err)
			}
			// fetch all metadata nodes into nodesA
			log.Printf("testing getNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bp.GetNodes(context.TODO(), NodeTypeMetadata, &nodesA)
			if err != nil {
				t.Fatal()
			}

			// sanity check
			if len(nodesA.Nodes) != 1 {
				t.Fatalf("not expecting %d '%s' nodes", len(nodesA.Nodes), NodeTypeMetadata)
			}

			newName := randString(10, "hex")
			// loop should run just once (len check above)
			for idA, nodeA := range nodesA.Nodes {
				log.Printf("node id: %s ; label: %s\n", idA, nodeA.Label)

				// change name to newName
				req := metadataNode{Label: newName}
				resp := &metadataNode{}
				log.Printf("testing patchNode(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
				err := bp.PatchNode(context.TODO(), nodeA.Id, req, resp)
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("response indicates name changed '%s' -> '%s'", nodeA.Label, resp.Label)

				// fetch changed node(s) (still expecting one) into nodesB
				log.Printf("testing getNodes(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
				err = bp.GetNodes(context.TODO(), NodeTypeMetadata, &nodesB)
				if err != nil {
					t.Fatal()
				}
				for idB, nodeB := range nodesB.Nodes {
					log.Printf("node id: %s ; label: %s\n", idB, nodeB.Label)

				}

				req = metadataNode{Label: nodeA.Label}
				log.Printf("testing patchNode() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = bp.PatchNode(context.TODO(), nodeA.Id, req, resp)
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("response indicates name changed '%s' -> '%s'", newName, resp.Label)
			}
		}
	}
	if len(testSkip) > 0 {
		sb := strings.Builder{}
		for _, msg := range testSkip {
			sb.WriteString(msg + ";")
		}
		t.Skip(sb.String())
	}
}
