// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestGetVersion(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		ver, err := client.client.getVersion(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		result, err := json.Marshal(ver)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s %s", client.client.baseUrl.String(), string(result))
	}
}
