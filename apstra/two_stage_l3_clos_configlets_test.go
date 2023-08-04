//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestImportGetUpdateGetDeleteConfiglet(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	var cr ConfigletData
	var cg []ConfigletGenerator
	var refarchs []RefDesign
	cg = append(cg, ConfigletGenerator{
		ConfigStyle:  PlatformOSJunos,
		Section:      ConfigletSectionSystem,
		TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
	})
	refarchs = append(refarchs, RefDesignTwoStageL3Clos)
	cr = ConfigletData{
		DisplayName: "TestImportConfiglet",
		RefArchs:    refarchs,
		Generators:  cg,
	}
	ctx := context.TODO()
	for clientName, client := range clients {
		// Create Configlet
		CatConfId, err := client.client.CreateConfiglet(ctx, &cr)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			client.client.DeleteConfiglet(ctx, CatConfId)
		}()

		bpClient, _ := testBlueprintA(ctx, t, client.client)
		// defer func() {
		// 	err = bpDel(ctx)
		// 	if err != nil {
		// 		t.Fatal(err)
		// 	}
		// }()

		log.Printf("testing ImportConfigletById() against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())
		icfg_id, err := bpClient.ImportConfigletById(ctx, CatConfId, "role in [\"spine\", \"leaf\"]", "")
		log.Printf("%s", icfg_id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		icfg, err := bpClient.GetConfiglet(ctx, icfg_id)
		log.Println(icfg)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteConfiglet(ctx, icfg_id)
		if err != nil {
			t.Fatal(err)
		}

		// Delete takes time sometimes
		time.Sleep(3 * time.Second)
		log.Printf("testing ImportConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		c := TwoStageL3ClosConfigletData{
			Data:      cr,
			Condition: "role in [\"spine\", \"leaf\"]",
			Label:     "",
		}
		icfg_id, err = bpClient.ImportConfiglet(ctx, c)
		log.Printf("%s", icfg_id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName,
			client.client.ApiVersion())
		icfg1, err := bpClient.GetConfigletByName(ctx, "TestImportConfiglet")
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
		icfg2, err := bpClient.GetConfiglet(ctx, icfg_id)
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
		err = bpClient.DeleteConfiglet(ctx, icfg_id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
