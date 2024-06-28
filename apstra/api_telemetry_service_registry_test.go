//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestTelemetryServiceRegistry(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	for clientName, client := range clients {
		log.Printf("Testing Telemetry Service Registry against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())
		log.Println("Test Get All Telemetry Service Registry Entries")
		entries, err := client.client.GetAllTelemetryServiceRegistryEntries(ctx)
		if err != nil {
			t.Fatal(err)
		}
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
              "interface",
              "supplicant_mac",
              "authenticated_vlan",
              "authorization_status",
              "port_status",
              "fallback_vlan_active"
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
			StorageSchemaPath: StorageSchemaPathIBA_INTEGER_DATA,
			ApplicationSchema: []byte(schema),
			Builtin:           false,
			Description:       "Test Service",
		}
		ServiceName, err := client.client.CreateTelemetryServiceRegistryEntry(ctx, entry)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(ServiceName)
		if ServiceName != name {
			t.Fatalf("Expected Service Name %s, Got %s", name, ServiceName)
		}
		pentry, err := client.client.GetTelemetryServiceRegistryEntry(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(pentry)
		if !jsonEqual(t, pentry.ApplicationSchema, entry.ApplicationSchema) {
			t.Fatal("Application Schema mismatch")
		}
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
			StorageSchemaPath: StorageSchemaPathIBA_INTEGER_DATA,
			ApplicationSchema: []byte(schema),
			Builtin:           false,
			Description:       "Test Service",
		}
		log.Println("Test Update Telemetry Service Registry Entry")

		err = client.client.UpdateTelemetryServiceRegistryEntry(ctx, ServiceName, &entry)
		if err != nil {
			t.Fatal(err)
		}

		pentry, err = client.client.GetTelemetryServiceRegistryEntry(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		if !jsonEqual(t, pentry.ApplicationSchema, entry.ApplicationSchema) {
			t.Fatal("Application Schema mismatch")
		}

		log.Println("Test Delete Telemetry Service Registry Entry")
		err = client.client.DeleteTelemetryServiceRegistryEntry(ctx, pentry.ServiceName)
		if err != nil {
			t.Fatal(err)
		}
		pentry, err = client.client.GetTelemetryServiceRegistryEntry(ctx, name)
		if err != nil {
			t.Log("Telemetry Service Successfully Removed. ", err)
		} else {
			t.Fatal("Delete has not succeeded. returned ", pentry)
		}
	}
}
