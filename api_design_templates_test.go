package goapstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestTemplateStrings(t *testing.T) {
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
		{stringVal: "heuristic", intType: AlgorithmHeuristic, stringType: algorithmHeuristic},

		{stringVal: "", intType: TemplateTypeNone, stringType: templateTypeNone},
		{stringVal: "rack_based", intType: TemplateTypeRackBased, stringType: templateTypeRackBased},
		{stringVal: "pod_based", intType: TemplateTypePodBased, stringType: templateTypePodBased},
		{stringVal: "l3_collapsed", intType: TemplateTypeL3Collapsed, stringType: templateTypeL3Collapsed},

		{stringVal: "distinct", intType: AsnAllocationSchemeDistinct, stringType: asnAllocationSchemeDistinct},
		{stringVal: "single", intType: AsnAllocationSchemeSingle, stringType: asnAllocationSchemeSingle},

		{stringVal: "ipv4", intType: AddressingSchemeIp4, stringType: addressingSchemeIp4},
		{stringVal: "ipv6", intType: AddressingSchemeIp6, stringType: addressingSchemeIp6},
		{stringVal: "ipv4_ipv6", intType: AddressingSchemeIp46, stringType: addressingSchemeIp46},

		{stringVal: "", intType: OverlayControlProtocolNone, stringType: overlayControlProtocolNone},
		{stringVal: "evpn", intType: OverlayControlProtocolEvpn, stringType: overlayControlProtocolEvpn},

		{stringVal: "", intType: TemplateCapabilityNone, stringType: templateCapabilityNone},
		{stringVal: "blueprint", intType: TemplateCapabilityBlueprint, stringType: templateCapabilityBlueprint},
		{stringVal: "pod", intType: TemplateCapabilityPod, stringType: templateCapabilityPod},
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

func TestGetTemplate(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing listAllTemplateIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		templateIds, err := client.client.listAllTemplateIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("fetching %d templateIds...", len(templateIds))

		for _, id := range templateIds {
			log.Printf("testing getTemplate(%s) against %s %s (%s)", id, client.clientType, client.clientName, client.client.ApiVersion())
			x, err := client.client.getTemplate(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}

			var id ObjectId
			var name string
			switch x.(Template).getType() {
			case TemplateTypeRackBased:
				id = x.(*TemplateRackBased).Id
				name = x.(*TemplateRackBased).DisplayName
			case TemplateTypePodBased:
				id = x.(*TemplatePodBased).Id
				name = x.(*TemplatePodBased).DisplayName
			case TemplateTypeL3Collapsed:
				id = x.(*TemplateL3Collapsed).Id
				name = x.(*TemplateL3Collapsed).DisplayName
			}
			log.Printf("template '%s' '%s'", id, name)
		}
	}
}

func TestGetTemplateMethods(t *testing.T) {
	var n int
	rand.Seed(time.Now().UnixNano())

	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing getAllTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		tmap, err := client.client.getAllTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		keys := make([]TemplateType, len(tmap))
		templateCount := make([]int, len(tmap))
		i := 0
		for k, v := range tmap {
			keys[i] = k
			templateCount[i] = len(v)
			i++
		}

		for i := 0; i < len(tmap); i++ {
			log.Printf("  %s template map has %d elements: ", keys[i].String(), templateCount[i])
		}

		// rack-based templates
		log.Printf("testing getAllRackBasedTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		rackBasedTemplates, err := client.client.getAllRackBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d rack-based templates\n", len(rackBasedTemplates))

		n = rand.Intn(len(rackBasedTemplates))
		log.Printf("testing getRackBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(rackBasedTemplates))
		rackBasedTemplate, err := client.client.getRackBasedTemplate(context.TODO(), rackBasedTemplates[n].Id)
		log.Printf("    got template type '%s', ID '%s'\n", rackBasedTemplate.Type, rackBasedTemplate.Id)

		// pod-based templates
		log.Printf("testing getAllPodBasedTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		podBasedTemplates, err := client.client.getAllPodBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d pod-based templates\n", len(podBasedTemplates))

		n = rand.Intn(len(podBasedTemplates))
		log.Printf("testing getPodBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(podBasedTemplates))
		podBasedTemplate, err := client.client.getPodBasedTemplate(context.TODO(), podBasedTemplates[n].Id)
		log.Printf("    got template type '%s', ID '%s'\n", podBasedTemplate.Type, podBasedTemplate.Id)

		// l3-collapsed templates
		log.Printf("testing getAllL3CollapsedTemplates() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		l3CollapsedTemplates, err := client.client.getAllL3CollapsedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("  got %d pod-based templates\n", len(l3CollapsedTemplates))

		n = rand.Intn(len(l3CollapsedTemplates))
		log.Printf("testing getL3CollapsedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(l3CollapsedTemplates))
		l3CollapsedTemplate, err := client.client.getL3CollapsedTemplate(context.TODO(), l3CollapsedTemplates[n].Id)
		log.Printf("    got template type '%s', ID '%s'\n", l3CollapsedTemplate.Type, l3CollapsedTemplate.Id)
	}
}

func TestGetTemplateAndType(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {

		log.Printf("testing ListAllTemplateIds() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		templateIds, err := client.client.ListAllTemplateIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		randomTemplateId := templateIds[rand.Intn(len(templateIds))]

		log.Printf("testing getTemplateAndType() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		tmplType, tmpl, err := client.client.getTemplateAndType(context.TODO(), randomTemplateId)
		if err != nil {
			t.Fatal(err)
		}

		var name string
		switch tmplType {
		case TemplateTypeRackBased:
			name = tmpl.(*TemplateRackBased).DisplayName
		case TemplateTypePodBased:
			name = tmpl.(*TemplatePodBased).DisplayName
		case TemplateTypeL3Collapsed:
			name = tmpl.(*TemplateL3Collapsed).DisplayName
		default:
			t.Fatalf("unknown template type '%d'", tmplType)
		}
		log.Printf("random template '%s' named '%s' has type '%s'", randomTemplateId, name, tmplType.String())
	}
}

func TestCreateDeleteRackBasedTemplate(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	dn := randString(5, "hex")
	req := CreateRackBasedTemplateRequest{
		DisplayName: dn,
		Capability:  TemplateCapabilityBlueprint,
		Spine: TemplateElementSpineRequest{
			Count:                   2,
			ExternalLinkSpeed:       "10G",
			LinkPerSuperspineSpeed:  "10G",
			LogicalDevice:           "AOS-7x10-Spine",
			LinkPerSuperspineCount:  1,
			Tags:                    nil,
			ExternalLinksPerNode:    0,
			ExternalFacingNodeCount: 0,
			ExternalLinkCount:       0,
		},
		RackTypeIds: []ObjectId{"evpn-single"},
		RackTypeCounts: []RackTypeCounts{{
			Count:      1,
			RackTypeId: "evpn-single",
		}},
		DhcpServiceIntent:      DhcpServiceIntent{Active: true},
		AntiAffinityPolicy:     AntiAffinityPolicy{Algorithm: AlgorithmHeuristic},
		AsnAllocationPolicy:    AsnAllocationPolicy{},
		FabricAddressingPolicy: FabricAddressingPolicy{},
		VirtualNetworkPolicy:   VirtualNetworkPolicy{},
	}

	for _, client := range clients {
		log.Printf("testing createRackBasedTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		id, err := client.client.createRackBasedTemplate(context.TODO(), &req)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing deleteTemplate() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.deleteTemplate(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
