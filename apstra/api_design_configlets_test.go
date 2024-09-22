//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
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
			ConfigStyle:  PlatformOSJunos,
			Section:      ConfigletSectionSystem,
			TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
		})
		var refarchs []RefDesign

		refarchs = append(refarchs, RefDesignTwoStageL3Clos)

		id1, err := client.client.CreateConfiglet(context.Background(), &ConfigletData{
			DisplayName: Name,
			RefArchs:    refarchs,
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
			ConfigStyle:  PlatformOSJunos,
			Section:      ConfigletSectionSystem,
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

func TestConfigletStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "system", intType: ConfigletSectionSystem, stringType: configletSectionSystem},
		{stringVal: "interface", intType: ConfigletSectionInterface, stringType: configletSectionInterface},
		{stringVal: "file", intType: ConfigletSectionFile, stringType: configletSectionFile},
		{stringVal: "frr", intType: ConfigletSectionFRR, stringType: configletSectionFRR},
		{stringVal: "ospf", intType: ConfigletSectionOSPF, stringType: configletSectionOSPF},
		{stringVal: "system_top", intType: ConfigletSectionSystemTop, stringType: configletSectionSystemTop},
		{stringVal: "set_based_system", intType: ConfigletSectionSetBasedSystem, stringType: configletSectionSetBasedSystem},
		{stringVal: "set_based_interface", intType: ConfigletSectionSetBasedInterface, stringType: configletSectionSetBasedInterface},
		{stringVal: "delete_based_interface", intType: ConfigletSectionDeleteBasedInterface, stringType: configletSectionDeleteBasedInterface},

		{stringVal: "cumulus", intType: PlatformOSCumulus, stringType: platformOSCumulus},
		{stringVal: "nxos", intType: PlatformOSNxos, stringType: platformOSNxos},
		{stringVal: "eos", intType: PlatformOSEos, stringType: platformOSEos},
		{stringVal: "junos", intType: PlatformOSJunos, stringType: platformOSJunos},
		{stringVal: "sonic", intType: PlatformOSSonic, stringType: platformOSSonic},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}

func TestAllPlatformOSes(t *testing.T) {
	all := AllPlatformOSes()
	expected := 5
	if len(all) != expected {
		t.Fatalf("expected %d platform OSes, got %d", expected, len(all))
	}
}

func TestAllConfigletSections(t *testing.T) {
	all := AllConfigletSections()
	expected := 9
	if len(all) != expected {
		t.Fatalf("expected %d configlet sections, got %d", expected, len(all))
	}
}
