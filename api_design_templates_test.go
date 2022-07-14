package goapstra

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestGetTemplate(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}
	err = client.Login(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	template, err := client.getTemplate(context.TODO(), "1b608f5e-fa27-437d-860f-bb7437f4ff20")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(template.(*TemplateRackBased).Type)
	log.Println(template.(*TemplateRackBased).Id)
	_ = template
}

func TestGetTemplateMethods(t *testing.T) {
	var n int
	rand.Seed(time.Now().UTC().UnixNano())

	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}
	err = client.Login(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	log.Println("testing getAllTemplates()")
	tmap, err := client.getAllTemplates(context.TODO())
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
	log.Printf("testing getAllRackBasedTemplates()\n")
	rackBasedTemplates, err := client.getAllRackBasedTemplates(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("  got %d rack-based templates\n", len(rackBasedTemplates))

	n = rand.Intn(len(rackBasedTemplates))
	log.Printf("testing getRackBasedTemplate() with randomly-selected index %d from the %d available\n", n, len(rackBasedTemplates))
	rackBasedTemplate, err := client.getRackBasedTemplate(context.TODO(), rackBasedTemplates[n].Id)
	log.Printf("  got template type '%s', ID '%s'\n", rackBasedTemplate.Type, rackBasedTemplate.Id)

	// pod-based templates
	log.Printf("testing getAllPodBasedTemplates()\n")
	podBasedTemplates, err := client.getAllPodBasedTemplates(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("  got %d pod-based templates\n", len(podBasedTemplates))

	n = rand.Intn(len(podBasedTemplates))
	log.Printf("testing getPodBasedTemplate() with randomly-selected index %d from the %d available\n", n, len(podBasedTemplates))
	podBasedTemplate, err := client.getPodBasedTemplate(context.TODO(), podBasedTemplates[n].Id)
	log.Printf("  got template type '%s', ID '%s'\n", podBasedTemplate.Type, podBasedTemplate.Id)

	// l3-collapsed templates
	log.Printf("testing getAllL3CollapsedTemplates()\n")
	l3CollapsedTemplates, err := client.getAllL3CollapsedTemplates(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("  got %d pod-based templates\n", len(l3CollapsedTemplates))

	n = rand.Intn(len(l3CollapsedTemplates))
	log.Printf("testing getL3CollapsedTemplate() with randomly-selected index %d from the %d available\n", n, len(l3CollapsedTemplates))
	l3CollapsedTemplate, err := client.getL3CollapsedTemplate(context.TODO(), l3CollapsedTemplates[n].Id)
	log.Printf("  got template type '%s', ID '%s'\n", l3CollapsedTemplate.Type, l3CollapsedTemplate.Id)
}
