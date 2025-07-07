// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"log"
	"math/rand"
	"strings"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listAllTemplateIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		templateIds, err := client.client.listAllTemplateIds(ctx)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("fetching %d templateIds...", len(templateIds))

		for _, i := range sampleIndexes(t, len(templateIds)) {
			templateId := templateIds[i]
			log.Printf("testing getTemplate(%s) against %s %s (%s)", templateId, client.clientType, clientName, client.client.ApiVersion())
			x, err := client.client.getTemplate(ctx, templateId)
			if err != nil {
				t.Fatal(err)
			}

			var name string
			tType, err := x.templateType()
			if err != nil {
				t.Fatal(err)
			}
			switch tType {
			case enum.TemplateTypeRackBased:
				var rawTemplate rawTemplateRackBased
				err = json.Unmarshal(x, &rawTemplate)
				if err != nil {
					t.Fatal(err)
				}
				name = rawTemplate.DisplayName
				rbt2, err := client.client.GetRackBasedTemplate(ctx, templateId)
				if err != nil {
					t.Fatal(err)
				}
				if templateId != rbt2.Id {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", templateId, rbt2.Id)
				}
				if name != rbt2.Data.DisplayName {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.Data.DisplayName)
				}
			case enum.TemplateTypePodBased:
				var rawTemplate rawTemplatePodBased
				err = json.Unmarshal(x, &rawTemplate)
				if err != nil {
					t.Fatal(err)
				}
				name = rawTemplate.DisplayName
				rbt2, err := client.client.GetPodBasedTemplate(ctx, templateId)
				if err != nil {
					t.Fatal(err)
				}
				if templateId != rbt2.Id {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", templateId, rbt2.Id)
				}
				if name != rbt2.Data.DisplayName {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.Data.DisplayName)
				}
			case enum.TemplateTypeL3Collapsed:
				var rawTemplate rawTemplateL3Collapsed
				err = json.Unmarshal(x, &rawTemplate)
				if err != nil {
					t.Fatal(err)
				}
				name = rawTemplate.DisplayName
				rbt2, err := client.client.GetL3CollapsedTemplate(ctx, templateId)
				if err != nil {
					t.Fatal(err)
				}
				if templateId != rbt2.Id {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", templateId, rbt2.Id)
				}
				if name != rbt2.Data.DisplayName {
					t.Fatalf("template ID mismatch: '%s' vs. '%s", name, rbt2.Data.DisplayName)
				}
			}
			log.Printf("template '%s' '%s'", templateId, name)
		}
	}
}

func TestGetTemplateMethods(t *testing.T) {
	var n int

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getAllTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		templates, err := client.client.getAllTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("got %d templates", len(templates))

		// rack-based templates
		log.Printf("testing getAllRackBasedTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		rackBasedTemplates, err := client.client.getAllRackBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d rack-based templates\n", len(rackBasedTemplates))

		n = rand.Intn(len(rackBasedTemplates))
		log.Printf("testing getRackBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(rackBasedTemplates))
		rackBasedTemplate, err := client.client.getRackBasedTemplate(context.TODO(), rackBasedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", rackBasedTemplate.Type, rackBasedTemplate.Id)

		// pod-based templates
		log.Printf("testing getAllPodBasedTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		podBasedTemplates, err := client.client.getAllPodBasedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got %d pod-based templates\n", len(podBasedTemplates))

		n = rand.Intn(len(podBasedTemplates))
		log.Printf("testing getPodBasedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(podBasedTemplates))
		podBasedTemplate, err := client.client.getPodBasedTemplate(context.TODO(), podBasedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", podBasedTemplate.Type, podBasedTemplate.Id)

		// l3-collapsed templates
		log.Printf("testing getAllL3CollapsedTemplates() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		l3CollapsedTemplates, err := client.client.getAllL3CollapsedTemplates(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("  got %d pod-based templates\n", len(l3CollapsedTemplates))

		n = rand.Intn(len(l3CollapsedTemplates))
		log.Printf("testing getL3CollapsedTemplate() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		log.Printf("  using randomly-selected index %d from the %d available\n", n, len(l3CollapsedTemplates))
		l3CollapsedTemplate, err := client.client.getL3CollapsedTemplate(context.TODO(), l3CollapsedTemplates[n].Id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("    got template type '%s', ID '%s'\n", l3CollapsedTemplate.Type, l3CollapsedTemplate.Id)
	}
}

func TestGetTemplateType(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	type testData struct {
		templateId   ObjectId
		templateType enum.TemplateType
	}

	data := []testData{
		{"pod1", enum.TemplateTypeRackBased},
		{"L2_superspine_multi_plane", enum.TemplateTypePodBased},
		{"L3_Collapsed_ACS", enum.TemplateTypeL3Collapsed},
	}

	for clientName, client := range clients {
		for _, d := range data {
			log.Printf("testing getTemplateType(%s) against %s %s (%s)", d.templateType, client.clientType, clientName, client.client.ApiVersion())
			ttype, err := client.client.GetTemplateType(ctx, d.templateId)
			if err != nil {
				t.Fatal(err)
			}
			if ttype != d.templateType {
				t.Fatalf("expected '%s', got '%s'", ttype, d.templateType)
			}
		}
	}
}

func TestGetTemplateIdsTypesByName(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	templateName := randString(10, "hex")
	for clientName, client := range clients {
		// fetch all template IDs
		templateIds, err := client.client.listAllTemplateIds(ctx)
		if err != nil {
			t.Fatal(err)
		}

		// choose a random template for cloning
		cloneMeId := templateIds[rand.Intn(len(templateIds))]
		cloneMeType, err := client.client.GetTemplateType(ctx, cloneMeId)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("cloning template '%s' (%s) for this test", cloneMeId, cloneMeType)

		cloneCount := rand.Intn(5) + 2
		cloneIds := make([]ObjectId, cloneCount)
		for i := 0; i < cloneCount; i++ {
			switch cloneMeType {
			case enum.TemplateTypeRackBased:
				cloneMe, err := client.client.getRackBasedTemplate(ctx, cloneMeId)
				if err != nil {
					t.Fatal(err)
				}
				id, err := client.client.createRackBasedTemplate(ctx, &rawCreateRackBasedTemplateRequest{
					Type:                 cloneMe.Type,
					DisplayName:          fmt.Sprintf("%s-%d", templateName, i),
					Spine:                cloneMe.Spine,
					RackTypes:            cloneMe.RackTypes,
					RackTypeCounts:       cloneMe.RackTypeCounts,
					DhcpServiceIntent:    cloneMe.DhcpServiceIntent,
					AntiAffinityPolicy:   cloneMe.AntiAffinityPolicy,
					AsnAllocationPolicy:  cloneMe.AsnAllocationPolicy,
					VirtualNetworkPolicy: cloneMe.VirtualNetworkPolicy,
				})
				if err != nil {
					t.Fatal(err)
				}
				cloneIds[i] = id
			case enum.TemplateTypePodBased:
				cloneMe, err := client.client.getPodBasedTemplate(ctx, cloneMeId)
				if err != nil {
					t.Fatal(err)
				}
				id, err := client.client.createPodBasedTemplate(ctx, &rawCreatePodBasedTemplateRequest{
					Type:                    cloneMe.Type,
					DisplayName:             fmt.Sprintf("%s-%d", templateName, i),
					Superspine:              cloneMe.Superspine,
					RackBasedTemplates:      cloneMe.RackBasedTemplates,
					RackBasedTemplateCounts: cloneMe.RackBasedTemplateCounts,
					AntiAffinityPolicy:      cloneMe.AntiAffinityPolicy,
				})
				if err != nil {
					t.Fatal(err)
				}
				cloneIds[i] = id
			case enum.TemplateTypeL3Collapsed:
				cloneMe, err := client.client.getL3CollapsedTemplate(ctx, cloneMeId)
				if err != nil {
					t.Fatal(err)
				}
				id, err := client.client.createL3CollapsedTemplate(ctx, &rawCreateL3CollapsedTemplateRequest{
					Type:                 cloneMe.Type,
					DisplayName:          fmt.Sprintf("%s-%d", templateName, i),
					MeshLinkCount:        cloneMe.MeshLinkCount,
					MeshLinkSpeed:        *cloneMe.MeshLinkSpeed,
					RackTypes:            cloneMe.RackTypes,
					RackTypeCounts:       cloneMe.RackTypeCounts,
					DhcpServiceIntent:    cloneMe.DhcpServiceIntent,
					AntiAffinityPolicy:   cloneMe.AntiAffinityPolicy,
					VirtualNetworkPolicy: cloneMe.VirtualNetworkPolicy,
				})
				if err != nil {
					t.Fatal(err)
				}
				cloneIds[i] = id
			}
		}
		clones := make([]string, len(cloneIds))
		for i, clone := range cloneIds {
			clones[i] = string(clone)
		}
		log.Printf("clone IDs: '%s'", strings.Join(clones, ", "))

		templateIdsToType := make(map[ObjectId]enum.TemplateType)
		for i := 0; i < cloneCount; i++ {
			log.Printf("testing getTemplateIdsTypesByName(%s) against %s %s (%s)", templateName, client.clientType, clientName, client.client.ApiVersion())
			temp, err := client.client.getTemplateIdsTypesByName(ctx, fmt.Sprintf("%s-%d", templateName, i))
			if err != nil {
				t.Fatal(err)
			}
			for k, v := range temp {
				templateIdsToType[k] = v
			}
		}

		if cloneCount != len(templateIdsToType) {
			t.Fatalf("expected %d, got %d", cloneCount, len(templateIdsToType))
		}
		for _, v := range templateIdsToType {
			if cloneMeType != v {
				t.Fatalf("expected %s, got %s", cloneMeType, v)
			}
		}

		for i, cloneId := range cloneIds {
			name := fmt.Sprintf("%s-%d", templateName, i)
			if i+1 == len(cloneIds) { // last one before they're all deleted
				log.Printf("testing getTemplateIdTypeByName(%s) against %s %s (%s)", name, client.clientType, clientName, client.client.ApiVersion())
				tId, tType, err := client.client.getTemplateIdTypeByName(ctx, name)
				if err != nil {
					t.Fatal(err)
				}
				if cloneId != tId {
					t.Fatalf("expected template id '%s', got '%s'", cloneIds, tId)
				}
				if cloneMeType != tType {
					t.Fatalf("expected template type '%s', got '%s'", cloneMeType, tType.String())
				}

			}
			log.Printf("deleting clone '%s'", cloneId)
			err = client.client.deleteTemplate(ctx, cloneId)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
