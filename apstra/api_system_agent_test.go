//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestGetSetSystemAgentManagerConfiguration(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		client.client.Login(context.Background())
		// initial fetch
		log.Printf("testing getSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		mgrCfg, err := client.client.getSystemAgentManagerConfig(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		// new config with opposite values
		testCfg := &SystemAgentManagerConfig{
			SkipRevertToPristineOnUninstall: !mgrCfg.SkipRevertToPristineOnUninstall,
			SkipPristineValidation:          !mgrCfg.SkipPristineValidation,
		}

		// set new config
		log.Printf("testing setSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.setSystemAgentManagerConfig(context.Background(), testCfg)
		if err != nil {
			t.Fatal(err)
		}

		// fetch new config
		log.Printf("testing getSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		testCfgRetrieved, err := client.client.getSystemAgentManagerConfig(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		// validate field as expected
		if testCfgRetrieved.SkipPristineValidation != testCfg.SkipPristineValidation {
			t.Fatalf("error setting skip pristine validation")
		}

		// validate field as expected
		if testCfgRetrieved.SkipRevertToPristineOnUninstall != testCfg.SkipRevertToPristineOnUninstall {
			t.Fatalf("error setting skip revert to pristine on uninstall")
		}

		// reset to original configuration
		log.Printf("testing setSystemAgentManagerConfig() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.setSystemAgentManagerConfig(context.Background(), mgrCfg)
		if err != nil {
			t.Fatal(err)
		}
	}
}
