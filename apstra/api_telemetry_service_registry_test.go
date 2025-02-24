// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestTelemetryServiceRegistry(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		log.Printf("Testing Telemetry Service Registry against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		log.Println("Test Get All Telemetry Service Registry Entries")
		entries, err := client.client.GetAllTelemetryServiceRegistryEntries(ctx)
		require.NoError(t, err)

		for _, e := range entries {
			log.Print(e.ServiceName)
		}

		log.Println("Test Create Telemetry Service Registry Entries")
		name := randString(10, "hex")
		schema := `{
        "required": ["key","value"],
        "type": "object",
        "properties": {
          	"value": {
            "type": "integer",
            "description": "0 in case of blocked, 1 in case of authorized"
          	},
			"key": {
            "required": [
              "authenticated_vlan",
              "authorization_status",
              "fallback_vlan_active",
              "interface",
              "port_status",
              "supplicant_mac"
            ],
            "type": "object",
            "properties": {
              "interface": {
                "type": "string"
              },
              "supplicant_mac": {
                "type": "string",
                "pattern": "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
              },
              "authenticated_vlan": {
                "type": "string"
              },
              "authorization_status": {
                "type": "string"
              },
              "port_status": {
                "enum": [
                  "authorized",
                  "blocked"
                ],
                "type": "string"
              },
              "fallback_vlan_active": {
                "enum": [
                  "True",
                  "False"
                ],
                "type": "string"
              }
            }
          }
		}
	}`

		entry := TelemetryServiceRegistryEntry{
			ServiceName:       name,
			StorageSchemaPath: enum.StorageSchemaPathIbaIntegerData,
			ApplicationSchema: []byte(schema),
			Builtin:           false,
			Description:       "Test Service",
		}

		ServiceName, err := client.client.CreateTelemetryServiceRegistryEntry(ctx, &entry)
		require.NoError(t, err)

		log.Println(ServiceName)
		require.Equal(t, name, ServiceName)

		pentry, err := client.client.GetTelemetryServiceRegistryEntry(ctx, name)
		require.NoError(t, err)

		log.Println(pentry)
		require.JSONEqf(t, string(pentry.ApplicationSchema), string(entry.ApplicationSchema), "expected: %s\nactual: %s", string(pentry.ApplicationSchema), string(entry.ApplicationSchema))

		schema = `{
        "required": ["key","value"],
        "type": "object",
        "properties": {
          	"value": {
            "type": "integer",
            "description": "0 in case of blocked, 1 in case of authorized"
          	},
			"key": {
            "required": [
              "supplicant_mac",
              "authenticated_vlan",
              "authorization_status",
              "port_status",
              "fallback_vlan_active"
            ],
            "type": "object",
            "properties": {
              "supplicant_mac": {
                "type": "string",
                "pattern": "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
              },
              "authenticated_vlan": {
                "type": "string"
              },
              "authorization_status": {
                "type": "string"
              },
              "port_status": {
                "enum": [
                  "authorized",
                  "blocked"
                ],
                "type": "string"
              },
              "fallback_vlan_active": {
                "enum": [
                  "True",
                  "False"
                ],
                "type": "string"
              }
            }
          }
		}
	}`
		entry = TelemetryServiceRegistryEntry{
			ServiceName:       name,
			StorageSchemaPath: enum.StorageSchemaPathIbaIntegerData,
			ApplicationSchema: []byte(schema),
			Builtin:           false,
			Description:       "Test Service",
		}
		log.Println("Test Update Telemetry Service Registry Entry")

		err = client.client.UpdateTelemetryServiceRegistryEntry(ctx, ServiceName, &entry)
		require.NoError(t, err)

		pentry, err = client.client.GetTelemetryServiceRegistryEntry(ctx, name)
		require.NoError(t, err)
		require.JSONEq(t, string(pentry.ApplicationSchema), string(entry.ApplicationSchema))

		log.Println("Test Delete Telemetry Service Registry Entry")
		err = client.client.DeleteTelemetryServiceRegistryEntry(ctx, pentry.ServiceName)
		require.NoError(t, err)

		var ace ClientErr

		_, err = client.client.GetTelemetryServiceRegistryEntry(ctx, name)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())

		err = client.client.DeleteTelemetryServiceRegistryEntry(ctx, name)
		require.Error(t, err)
		require.ErrorAs(t, err, &ace)
		require.Equal(t, ErrNotfound, ace.Type())
	}
}
