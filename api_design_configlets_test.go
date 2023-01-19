//go:build integration
// +build integration

package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestCreateUpdateGetDeleteConfiglet(t *testing.T) {
	clients, err := getTestClients(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	Name := randString(10, "hex")
	for _, client := range clients {
		var cg []ConfigletGenerator

		cg = append(cg, ConfigletGenerator{
			ConfigStyle:  ApstraPlatformOSJunos,
			Section:      ApstraConfigletSectionSystem,
			TemplateText: "interfaces {\n   {% if 'leaf1' in hostname %}\n    xe-0/0/3 {\n      disable;\n    }\n   {% endif %}\n   {% if 'leaf2' in hostname %}\n    xe-0/0/2 {\n      disable;\n    }\n   {% endif %}\n}",
		})
		var refarchs []RefDesign

		refarchs = append(refarchs, RefDesignTwoStageL3Clos)

		id1, err := client.client.CreateConfiglet(context.Background(), &ConfigletRequest{
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
			ConfigStyle:  ApstraPlatformOSJunos,
			Section:      ApstraConfigletSectionSystem,
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
		{stringVal: "system", intType: ApstraConfigletSectionSystem, stringType: apstraConfigletSectionSystem},
		{stringVal: "interface", intType: ApstraConfigletSectionInterface, stringType: apstraConfigletSectionInterface},
		{stringVal: "file", intType: ApstraConfigletSectionFile, stringType: apstraConfigletSectionFile},
		{stringVal: "frr", intType: ApstraConfigletSectionFRR, stringType: apstraConfigletSectionFRR},
		{stringVal: "ospf", intType: ApstraConfigletSectionOSPF, stringType: apstraConfigletSectionOSPF},
		{stringVal: "system_top", intType: ApstraConfigletSectionSystemTop, stringType: apstraConfigletSectionSystemTop},
		{stringVal: "set_based_system", intType: ApstraConfigletSectionSetBasedSystem, stringType: apstraConfigletSectionSetBasedSystem},
		{stringVal: "set_based_interface", intType: ApstraConfigletSectionSetBasedInterface, stringType: apstraConfigletSectionSetBasedInterface},
		{stringVal: "delete_based_interface", intType: ApstraConfigletSectionDeleteBasedInterface, stringType: apstraConfigletSectionDeleteBasedInterface},

		{stringVal: "cumulus", intType: ApstraPlatformOSCumulus, stringType: apstraPlatformOSCumulus},
		{stringVal: "nxos", intType: ApstraPlatformOSNxos, stringType: apstraPlatformOSNxos},
		{stringVal: "eos", intType: ApstraPlatformOSEos, stringType: apstraPlatformOSEos},
		{stringVal: "junos", intType: ApstraPlatformOSJunos, stringType: apstraPlatformOSJunos},
		{stringVal: "sonic", intType: ApstraPlatformOSSonic, stringType: apstraPlatformOSSonic},
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
