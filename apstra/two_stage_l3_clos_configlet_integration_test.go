//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"log"
	"testing"
)

func TestImportGetUpdateGetDeleteConfiglet(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	configletData := ConfigletData{
		DisplayName: "TestImportConfiglet",
		RefArchs:    []RefDesign{RefDesignTwoStageL3Clos},
		Generators: []ConfigletGenerator{{
			ConfigStyle:  PlatformOSJunos,
			Section:      ConfigletSectionSystem,
			TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
		}},
	}

	for clientName, client := range clients {
		// Create Configlet
		catalogConfigletId, err := client.client.CreateConfiglet(ctx, &configletData)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err := client.client.DeleteConfiglet(ctx, catalogConfigletId)
			if err != nil {
				t.Fatal(err)
			}
		}()

		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		log.Printf("testing ImportConfigletById() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		bpConfigletId, err := bpClient.ImportConfigletById(ctx, catalogConfigletId, `role in ["spine", "leaf"]`, "")
		log.Printf("%s", bpConfigletId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err = bpClient.GetConfiglet(ctx, bpConfigletId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteConfiglet(ctx, bpConfigletId)
		if err != nil {
			t.Fatal(err)
		}

		_, err = bpClient.GetConfiglet(ctx, bpConfigletId)
		if err == nil {
			t.Fatal("fetch configlet after delete should have produced an error")
		} else {
			var ace ApstraClientErr
			if !errors.As(err, &ace) || ace.Type() != ErrNotfound {
				t.Fatal("fetch configlet after delete should have produced ErrNotFound")
			}
		}

		configletData.DisplayName = "ImportDirect"
		log.Printf("testing ImportConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		c := TwoStageL3ClosConfigletData{
			Data:      &configletData,
			Condition: "role in [\"spine\", \"leaf\"]",
			Label:     "",
		}
		bpConfigletId, err = bpClient.CreateConfiglet(ctx, &c)
		log.Printf("%s", bpConfigletId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())
		icfg1, err := bpClient.GetConfigletByName(ctx, "ImportDirect")
		log.Println(icfg1)
		if err != nil {
			t.Fatal(err)
		}
		icfg1.Data.Label = "new name"
		icfg1.Data.Condition = "role in [\"spine\"]"
		log.Printf("testing UpdateConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.UpdateConfiglet(ctx, icfg1)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		icfg2, err := bpClient.GetConfiglet(ctx, bpConfigletId)
		log.Println(icfg2)
		if err != nil {
			t.Fatal(err)
		}
		if icfg1.Data.Label != icfg2.Data.Label {
			t.Fatal("Name Change Failed")
		}
		if icfg1.Data.Condition != icfg2.Data.Condition {
			t.Fatal("Condition Change Failed")
		}
		log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteConfiglet(ctx, bpConfigletId)
		if err != nil {
			t.Fatal(err)
		}
	}
}
