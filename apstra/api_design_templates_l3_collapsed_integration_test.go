// Copyright (c) Juniper Networks, Inc., 2025-2025.
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

func TestCreateGetDeleteL3CollapsedTemplate(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")

	req := &CreateL3CollapsedTemplateRequest{
		DisplayName:   dn,
		MeshLinkCount: 1,
		MeshLinkSpeed: "10G",
		RackTypeIds:   []ObjectId{"L3_collapsed_acs"},
		RackTypeCounts: []RackTypeCount{{
			RackTypeId: "L3_collapsed_acs",
			Count:      1,
		}},
		VirtualNetworkPolicy: VirtualNetworkPolicy{OverlayControlProtocol: enum.OverlayControlProtocolEvpn},
	}

	for clientName, client := range clients {
		log.Printf("testing CreateL3CollapsedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateL3CollapsedTemplate(ctx, req)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteTemplate(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetL3CollapsedTemplateByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	name := "Collapsed Fabric ESI"

	for _, client := range clients {
		l3ct, err := client.client.GetL3CollapsedTemplateByName(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		if l3ct.templateType != enum.TemplateTypeL3Collapsed {
			t.Fatalf("expected '%s', got '%s'", l3ct.templateType.String(), enum.TemplateTypeL3Collapsed)
		}
		if l3ct.Data.DisplayName != name {
			t.Fatalf("expected '%s', got '%s'", name, l3ct.Data.DisplayName)
		}
	}
}
