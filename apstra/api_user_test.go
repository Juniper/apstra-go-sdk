// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestUserLogin(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		t.Run(clientName, func(t *testing.T) {
			t.Parallel()

			if client.clientType == "api-ops" {
				t.Skipf("skipping test - api-ops type clients do not log in or out")
			}

			log.Printf("testing Login() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.Login(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing Logout() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.Logout(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing redundant Logout() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.Logout(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing Login() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.Login(context.Background())
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
