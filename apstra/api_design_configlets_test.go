// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

func TestCreateUpdateGetDeleteConfiglet(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	Name := randString(10, "hex")
	for _, client := range clients {
		var cg []ConfigletGenerator

		cg = append(cg, ConfigletGenerator{
			ConfigStyle:  enum.ConfigletStyleJunos,
			Section:      enum.ConfigletSectionSystem,
			TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
		})

		id1, err := client.client.CreateConfiglet(context.Background(), &ConfigletData{
			DisplayName: Name,
			RefArchs:    []enum.RefDesign{enum.RefDesignDatacenter},
			Generators:  cg,
		})
		if err != nil {
			t.Fatal(err)
		}

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
			ConfigStyle:  enum.ConfigletStyleJunos,
			Section:      enum.ConfigletSectionSystem,
			TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
		})
		log.Println("Update Config")

		err = client.client.UpdateConfiglet(context.Background(), id1, &ConfigletData{
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
