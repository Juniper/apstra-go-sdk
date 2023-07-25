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
	var cr ConfigletRequest
	var cg []ConfigletGenerator
	var refarchs []RefDesign
	cg = append(cg, ConfigletGenerator{
		ConfigStyle:  PlatformOSJunos,
		Section:      ConfigletSectionSystem,
		TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
	})
	refarchs = append(refarchs, RefDesignTwoStageL3Clos)
	cr = ConfigletRequest{
		DisplayName: "TestImportConfiglet",
		RefArchs:    refarchs,
		Generators:  cg,
	}
	ctx := context.TODO()
	for clientName, client := range clients {
		//Create Configlet
		CatConfId, err := client.client.CreateConfiglet(ctx, &cr)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			client.client.DeleteConfiglet(ctx, CatConfId)
		}()

		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		log.Printf("testing ImportConfigletByID() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips_id, err := bpClient.ImportConfigletByID(ctx, CatConfId, "role in [\"spine\", \"leaf\"]", "")
		log.Printf("%s", ips_id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips, err := bpClient.GetConfiglet(ctx, ips_id)
		log.Println(ips)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteConfiglet(ctx, ips_id)
		if err != nil {
			t.Fatal(err)
		}

		// Delete takes time sometimes
		time.Sleep(3 * time.Second)
		log.Printf("testing ImportConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips_id, err = bpClient.ImportConfiglet(ctx, (ConfigletData)(cr), "role in [\"spine\", \"leaf\"]", "")
		log.Printf("%s", ips_id)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips, err = bpClient.GetConfiglet(ctx, ips_id)
		log.Println(ips)
		if err != nil {
			t.Fatal(err)
		}

		ips.Label = "new name"
		ips.Condition = "role in [\"spine\"]"
		log.Printf("testing UpdateConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.UpdateConfiglet(ctx, ips_id, ips)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ips1, err := bpClient.GetConfiglet(ctx, ips_id)
		log.Println(ips1)
		if err != nil {
			t.Fatal(err)
		}
		if ips1.Label != ips.Label {
			t.Fatal("Name Change Failed")
		}
		if ips1.Condition != ips.Condition {
			t.Fatal("Condition Change Failed")
		}

		log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteConfiglet(ctx, ips_id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
