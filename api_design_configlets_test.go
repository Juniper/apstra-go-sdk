//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestCreateUpdateGetDeleteConfiglet(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	Name := randString(10, "hex")
	for _, client := range clients {
		var cg []ConfigletGenerator

		cg = append(cg, ConfigletGenerator{
			ConfigStyle:  "junos",
			Section:      "system",
			TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
		})
		var refarchs []string

		refarchs = append(refarchs, "two_stage_l3clos")

		id1, err := client.client.CreateConfiglet(context.Background(), &ConfigletRequest{
			DisplayName: Name,
			RefArchs:    refarchs,
			Generators:  cg,
		})
		if err != nil {
			t.Fatal(err)
		}
		log.Println(id1)
		log.Println("Getting now")
		c, err := client.client.GetConfiglet(context.Background(), id1)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(c)
		g1 := len(c.Data.Generators)
		c.Data.Generators = append(c.Data.Generators, ConfigletGenerator{
			ConfigStyle:  "junos",
			Section:      "system",
			TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
		})
		log.Println("Update Config")

		err = client.client.UpdateConfiglet(context.Background(), id1, &ConfigletRequest{
			DisplayName: c.Data.DisplayName,
			RefArchs:    c.Data.RefArchs,
			Generators:  c.Data.Generators,
		})

		if err != nil {
			t.Fatal(err)
		}
		log.Println("Get Configlet by Name")
		c, err = client.client.GetConfigletByName(context.Background(), Name)
		if err != nil {
			t.Fatal(err)
		}
		g2 := len(c.Data.Generators)
		log.Println(g1, g2)
		if g1 == g2 {
			t.Fatal("append did not work")
		}
		log.Println(c)
		log.Println("Deleting now")

		err = client.client.DeleteConfiglet(context.Background(), id1)
		if err != nil {
			t.Fatal(err)
		}

		log.Println("Testing an incorrect delete")
		err = client.client.DeleteConfiglet(context.Background(), id1)
		log.Println(err)
		if err == nil {
			t.Fatal("Error :: Deleting non-existent item works")
		}
	}
}
