// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

func TestGetLockInfo(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			log.Printf("testing createBlueprintFromTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			name := randString(10, "hex")
			id, err := client.client.CreateBlueprintFromTemplate(context.TODO(), &CreateBlueprintFromTemplateRequest{
				RefDesign:  enum.RefDesignDatacenter,
				Label:      name,
				TemplateId: "L2_Virtual_EVPN",
			})
			if err != nil {
				t.Fatal(err)
			}

			bp, err := client.client.NewTwoStageL3ClosClient(context.Background(), id)
			if err != nil {
				t.Fatal(err)
			}

			l, err := bp.GetLockInfo(context.TODO())
			if err != nil {
				t.Fatal(err)
			}
			log.Println(l)
			log.Printf("got id '%s', deleting blueprint...\n", id)
			log.Printf("testing deleteBlueprint() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.deleteBlueprint(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
