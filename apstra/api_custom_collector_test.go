//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

	for clientName, client := range clients {
		log.Printf("Testing Custom Collectors against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())

		ts, err := client.client.GetAllTelemetryServiceRegistryEntries(ctx)
		for _, tsr := range ts {
			c, err := client.client.GetCollectorsByServiceName(ctx, tsr.ServiceName)
			if err != nil {
				t.Fatalf(err.Error())
			}
			for _, d := range c {
				log.Printf("%v", d)
			}
		}

		name := randString(10, "hex")
		schema := `{
			"properties": {
			  "key": {
				"properties": {
				  "schemakey1": {
					"type": "string"
				  }
				},
				"required": [
				  "schemakey1"
				],
				"type": "object"
			  },
			  "value": {
				"type": "string"
			  }
			},
			"required": [
			  "key",
			  "value"
			],
			"type": "object"
		  }`

		entry := TelemetryServiceRegistryEntry{
			ServiceName:       name,
			StorageSchemaPath: StorageSchemaPathIBA_STRING_DATA,
			ApplicationSchema: []byte(schema),
			Builtin:           false,
			Description:       "Test Service %s",
		}
		ServiceName, err := client.client.CreateTelemetryServiceRegistryEntry(ctx, &entry)
		log.Printf("Service Name %s Created ", ServiceName)
		require.NoError(t, err)
		cs, err := client.client.GetCollectorsByServiceName(ctx, name)
		require.NoError(t, err)
		if len(cs) != 0 {
			log.Println("There should be no collectors, this is a new service")
		}

		c1 := Collector{
			ServiceName: name,
			Platform: CollectorPlatform{
				OsType:    CollectorOSTypeJunosEvo,
				OsVersion: "22.2r2",
				OsFamily:  []CollectorOSFamily{CollectorOSFamilyACX},
				Model:     "",
			},
			SourceType: CollectorSourceTypeCLI,
			Cli:        "cli show interfaces extensive",
			Query: CollectorQuery{
				Accessors: map[string]string{"telemetrykey1": "/interface-information/docsis-information/docsis-media-properties/downstream-buffers-free"},
				Keys:      map[string]string{"schemakey1": "telemetrykey1"},
				Value:     "telemetrykey1",
				Filter:    "",
			},
			RelaxedSchemaValidation: true,
		}
		c2 := c1
		log.Println("Creating First Collector")
		err = client.client.CreateCollector(ctx, &c1)
		require.NoError(t, err)

		cs, err = client.client.GetCollectorsByServiceName(ctx, name)
		require.NoError(t, err)
		if len(cs) != 1 {
			log.Printf("There should be one collector, got %d", len(cs))
		}
		log.Println("Creating Second Collector")

		c1.Platform.OsFamily = []CollectorOSFamily{CollectorOSFamilyACX_F, CollectorOSFamilyJunos}
		err = client.client.CreateCollector(ctx, &c1)
		require.NoError(t, err)
		cs, err = client.client.GetCollectorsByServiceName(ctx, name)
		require.NoError(t, err)
		if len(cs) != 2 {
			log.Printf("There should be two collectors, got %d", len(cs))
		}

		log.Println("Updating Collector")
		c1.Query.Accessors["telemetrykey1"] = "/interface-information/docsis-information/docsis-media-properties/downstream-buffers-used"
		err = client.client.UpdateCollector(ctx, &c1)
		require.NoError(t, err)
		cs, err = client.client.GetCollectorsByServiceName(ctx, name)
		require.NoError(t, err)
		if len(cs) != 2 {
			log.Printf("There should be two collectors, got %d", len(cs))
		}

		err = client.client.DeleteCollector(ctx, &c1)
		require.NoError(t, err)
		cs, err = client.client.GetCollectorsByServiceName(ctx, name)
		require.NoError(t, err)
		if len(cs) != 1 {
			log.Printf("There should be one collector, got %d", len(cs))
		}

		err = client.client.DeleteCollector(ctx, &c2)
		require.NoError(t, err)
		cs, err = client.client.GetCollectorsByServiceName(ctx, name)
		require.NoError(t, err)
		if len(cs) != 0 {
			log.Printf("There should be no collectors, got %d", len(cs))
		}

		err = client.client.DeleteTelemetryServiceRegistryEntry(ctx, ServiceName)
		require.NoError(t, err)
	}
}
