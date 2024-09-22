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

func TestGetVersionsAll(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getVersionsAosdi() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		aosdi, err := client.client.getVersionsAosdi(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing getVersionsApi() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		api, err := client.client.getVersionsApi(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing getVersionsBuild() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		build, err := client.client.getVersionsBuild(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing getVersionsServer() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		server, err := client.client.getVersionsServer(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		body, err := json.Marshal(&struct {
			Aosdi  *versionsAosdiResponse  `json:"aosdi"`
			Api    *versionsApiResponse    `json:"api"`
			Build  *versionsBuildResponse  `json:"build"`
			Server *versionsServerResponse `json:"server"`
		}{
			Aosdi:  aosdi,
			Api:    api,
			Build:  build,
			Server: server,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Println(string(body))
	}
}
