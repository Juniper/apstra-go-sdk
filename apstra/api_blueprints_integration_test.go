//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestListAllBlueprintIds(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
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
	clients, err := getTestClients(context.Background(), t)
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
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing createBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		name := randString(10, "hex")
		id, err := client.client.CreateBlueprintFromTemplate(context.TODO(), &CreateBlueprintFromTemplateRequest{
			RefDesign:  RefDesignTwoStageL3Clos,
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
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

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

		type nodes struct {
			Nodes map[string]metadataNode `json:"nodes"`
		}
		var nodesA, nodesB nodes

		// fetch all metadata nodes into nodesA
		log.Printf("testing getNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.GetNodes(ctx, NodeTypeMetadata, &nodesA)
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
			log.Printf("testing patchNode(%s) against %s %s (%s)", bpClient.Id(), client.clientType, clientName, client.client.ApiVersion())
			err := bpClient.PatchNode(ctx, nodeA.Id, req, resp)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Label != newName {
				t.Fatalf("expected new blueprint name %q, got %q", newName, resp.Label)
			}
			log.Printf("response indicates name changed '%s' -> '%s'", nodeA.Label, resp.Label)

			// fetch changed node(s) (still expecting one) into nodesB
			log.Printf("testing getNodes(%s) against %s %s (%s)", bpClient.Id(), client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.GetNodes(ctx, NodeTypeMetadata, &nodesB)
			if err != nil {
				t.Fatal()
			}
			for idB, nodeB := range nodesB.Nodes {
				log.Printf("node id: %s ; label: %s\n", idB, nodeB.Label)
				if nodeB.Label != newName {
					t.Fatalf("expected new blueprint name %q, got %q", newName, nodeB.Label)
				}
			}
		}
	}
}

func TestGetNodes(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintB(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		type node struct {
			Id         ObjectId `json:"id"`
			Label      string   `json:"label"`
			SystemType string   `json:"system_type"`
		}
		equal := func(a, b node) bool {
			return a.Id == b.Id &&
				a.Label == b.Label &&
				a.SystemType == b.SystemType
		}

		var response struct {
			Nodes map[ObjectId]node `json:"nodes"`
		}
		log.Printf("testing GetNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.Client().GetNodes(ctx, bpClient.Id(), NodeTypeSystem, &response)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("got %d nodes. Fetch each one...", len(response.Nodes))
		var nodeB node
		for id, nodeA := range response.Nodes {
			log.Printf("testing GetNode() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.Client().GetNode(ctx, bpClient.Id(), id, &nodeB)
			if err != nil {
				t.Fatal()
			}
			if !equal(nodeA, nodeB) {
				t.Fatalf("nodes don't match:\n%v\n%v", nodeA, nodeB)
			}
		}
	}
}

func TestPatchNodes(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintB(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		type node struct {
			Id         ObjectId `json:"id"`
			Label      string   `json:"label"`
			SystemType string   `json:"system_type,omitempty"`
		}

		var getResponse struct {
			Nodes map[ObjectId]node `json:"nodes"`
		}
		log.Printf("testing GetNodes() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.Client().GetNodes(ctx, bpClient.Id(), NodeTypeSystem, &getResponse)
		if err != nil {
			t.Fatal(err)
		}

		var patch []interface{}
		for k, v := range getResponse.Nodes {
			if v.SystemType == "server" {
				patch = append(patch, node{
					Id:    k,
					Label: randString(5, "hex"),
				})
			}
		}

		err = client.client.PatchNodes(ctx, bpClient.Id(), patch)
		if err != nil {
			t.Fatal(err)
		}

		for _, n := range patch {
			var result node
			err = client.client.GetNode(ctx, bpClient.Id(), n.(node).Id, &result)
			if err != nil {
				t.Fatal(err)
			}

			if n.(node).Label != result.Label {
				t.Fatalf("patch expected label %s, got label %s", n.(node).Label, result.Label)
			}
		}
	}
}
