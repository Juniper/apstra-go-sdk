package goapstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestGetTemplate(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
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
	rand.Seed(time.Now().UTC().UnixNano())

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
