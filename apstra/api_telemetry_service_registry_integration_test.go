// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestTelemetryServiceRegistry(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			log.Println("Test Get All Telemetry Service Registry Entries")
			entries, err := client.Client.GetAllTelemetryServiceRegistryEntries(ctx)
			require.NoError(t, err)

			for _, e := range entries {
				log.Print(e.ServiceName)
			}

			log.Println("Test Create Telemetry Service Registry Entries")
			name := testutils.RandString(10, "hex")
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

			entry := apstra.TelemetryServiceRegistryEntry{
				ServiceName:       name,
				StorageSchemaPath: enum.StorageSchemaPathIbaIntegerData,
				ApplicationSchema: []byte(schema),
				Builtin:           false,
				Description:       "Test Service",
			}

			ServiceName, err := client.Client.CreateTelemetryServiceRegistryEntry(ctx, &entry)
			require.NoError(t, err)

			log.Println(ServiceName)
			require.Equal(t, name, ServiceName)

			pentry, err := client.Client.GetTelemetryServiceRegistryEntry(ctx, name)
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
			entry = apstra.TelemetryServiceRegistryEntry{
				ServiceName:       name,
				StorageSchemaPath: enum.StorageSchemaPathIbaIntegerData,
				ApplicationSchema: []byte(schema),
				Builtin:           false,
				Description:       "Test Service",
			}
			log.Println("Test Update Telemetry Service Registry Entry")

			err = client.Client.UpdateTelemetryServiceRegistryEntry(ctx, ServiceName, &entry)
			require.NoError(t, err)

			pentry, err = client.Client.GetTelemetryServiceRegistryEntry(ctx, name)
			require.NoError(t, err)
			require.JSONEq(t, string(pentry.ApplicationSchema), string(entry.ApplicationSchema))

			log.Println("Test Delete Telemetry Service Registry Entry")
			err = client.Client.DeleteTelemetryServiceRegistryEntry(ctx, pentry.ServiceName)
			require.NoError(t, err)

			var ace apstra.ClientErr

			_, err = client.Client.GetTelemetryServiceRegistryEntry(ctx, name)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())

			err = client.Client.DeleteTelemetryServiceRegistryEntry(ctx, name)
			require.Error(t, err)
			require.ErrorAs(t, err, &ace)
			require.Equal(t, apstra.ErrNotfound, ace.Type())
		})
	}
}
