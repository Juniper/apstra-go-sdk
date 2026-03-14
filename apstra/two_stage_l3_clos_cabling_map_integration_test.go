// Copyright (c) Juniper Networks, Inc., 2023-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"sort"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
)

func TestGetCablingMapLinks(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	linksAreEqual := func(a, b CablingMapLink) bool {
		if a.ID != b.ID {
			return false
		}
		if (a.Label == nil) != (b.Label == nil) {
			return false
		}
		if (a.Label != nil && b.Label != nil) && *a.Label != *b.Label {
			return false
		}
		return true
	}

	for clientName, client := range clients {
		bpClient := testBlueprintB(ctx, t, client.client)

		log.Printf("testing GetCablingMapLinks() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		links, err := bpClient.GetCablingMapLinks(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(links) != 16 {
			t.Fatalf("expectex 16 links, got %d", len(links))
		}

		systemIdToLinks := make(map[string][]CablingMapLink)
		for i, link := range links {
			for _, endpoint := range link.Endpoints {
				systemIdToLinks[endpoint.System.ID] = append(systemIdToLinks[endpoint.System.ID], links[i])
			}
		}

		for systemId, linksA := range systemIdToLinks {
			sort.Slice(linksA, func(i, j int) bool {
				return linksA[i].ID < linksA[j].ID
			})

			linksB, err := bpClient.GetCablingMapLinksBySystem(ctx, systemId)
			if err != nil {
				t.Fatal(err)
			}

			if len(linksA) != len(linksB) {
				t.Fatalf("length of linksA (%d) doesn't match length of linksB(%d)", len(linksA), len(linksB))
			}

			sort.Slice(linksB, func(i, j int) bool {
				return linksB[i].ID < linksB[j].ID
			})

			for i := range linksA {
				if !linksAreEqual(linksA[i], linksB[i]) {
					t.Fatalf("linksA[%d] doesn't match linksB[%d]: %v vs. %v", i, i, linksA[i], linksB[i])
				}
			}
		}
	}
}
