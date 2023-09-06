package apstra

import (
	"context"
	"log"
	"testing"
)

func TestSystemNodeInfo(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing ListAllBlueprintIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpIds, err := client.client.ListAllBlueprintIds(ctx)
		if err != nil {
			t.Fatal(err)
		}

		var bpClient *TwoStageL3ClosClient
		if len(bpIds) > 0 {
			log.Printf("testing NewTwoStageL3ClosClient() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			bpClient, err = client.client.NewTwoStageL3ClosClient(ctx, bpIds[0])
			if err != nil {
				t.Fatal(err)
			}
		} else {
			var deleteFunc func(ctx2 context.Context) error
			bpClient, deleteFunc = testBlueprintA(ctx, t, client.client)
			defer deleteFunc(ctx)
		}

		log.Printf("testing GetAllSystemNodeInfos() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		nodeInfos, err := bpClient.GetAllSystemNodeInfos(ctx)
		if err != nil {
			t.Fatal(err)
		}

		for nodeId := range nodeInfos {
			log.Printf("testing GetSystemNodeInfo() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			nodeInfo, err := bpClient.GetSystemNodeInfo(ctx, nodeId)
			if err != nil {
				t.Fatal(err)
			}
			log.Println(nodeInfo.Id)
		}
	}
}
