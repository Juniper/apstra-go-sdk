// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"strings"
	"testing"

	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
)

func TestGetTemplate(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx = testutils.WrapCtxWithTestId(t, ctx)

			templateIds, err := client.Client.ListAllTemplateIds(ctx)
			require.NoError(t, err)

			templates, err := client.Client.GetAllTemplates(ctx)
			require.NoError(t, err)

			require.Equal(t, len(templateIds), len(templates))

			log.Printf("fetching %d templateIds...", len(templates))

			for _, i := range testutils.SampleIndexes(t, len(templates)) {
				switch templates[i].Type() {
				case apstra.TemplateTypeRackBased:
					template, err := client.Client.GetRackBasedTemplate(ctx, templates[i].ID())
					require.NoError(t, err)
					require.Equal(t, templates[i].ID(), template.ID())
				case apstra.TemplateTypePodBased:
					template, err := client.Client.GetPodBasedTemplate(ctx, templates[i].ID())
					require.NoError(t, err)
					require.Equal(t, templates[i].ID(), template.ID())
				case apstra.TemplateTypeL3Collapsed:
					template, err := client.Client.GetL3CollapsedTemplate(ctx, templates[i].ID())
					require.NoError(t, err)
					require.Equal(t, templates[i].ID(), template.ID())
				default:
					tt := templates[i].Type()
					t.Fatalf("unknown template type: %d (%s)", tt, tt)
				}
			}
		})
	}
}

func TestGetTemplateMethods(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	var n int

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx = testutils.WrapCtxWithTestId(t, ctx)

			templates, err := client.Client.GetAllTemplates(ctx)
			require.NoError(t, err)
			log.Printf("got %d templates", len(templates))

			// rack-based templates
			rackBasedTemplates, err := client.Client.GetAllRackBasedTemplates(ctx)
			require.NoError(t, err)
			log.Printf("    got %d rack-based templates\n", len(rackBasedTemplates))

			n = rand.Intn(len(rackBasedTemplates))
			log.Printf("  using randomly-selected index %d from the %d available\n", n, len(rackBasedTemplates))
			rackBasedTemplate, err := client.Client.GetRackBasedTemplate(ctx, rackBasedTemplates[n].Id)
			require.NoError(t, err)
			log.Printf("    got template type '%s', ID '%s'\n", rackBasedTemplate.Type, rackBasedTemplate.Id)

			// pod-based templates
			podBasedTemplates, err := client.Client.GetAllPodBasedTemplates(ctx)
			require.NoError(t, err)
			log.Printf("    got %d pod-based templates\n", len(podBasedTemplates))

			n = rand.Intn(len(podBasedTemplates))
			log.Printf("  using randomly-selected index %d from the %d available\n", n, len(podBasedTemplates))
			podBasedTemplate, err := client.Client.GetPodBasedTemplate(ctx, podBasedTemplates[n].Id)
			require.NoError(t, err)
			log.Printf("    got template type '%s', ID '%s'\n", podBasedTemplate.Type, podBasedTemplate.Id)

			// l3-collapsed templates
			l3CollapsedTemplates, err := client.Client.GetAllL3CollapsedTemplates(ctx)
			require.NoError(t, err)
			log.Printf("  got %d pod-based templates\n", len(l3CollapsedTemplates))

			n = rand.Intn(len(l3CollapsedTemplates))
			log.Printf("  using randomly-selected index %d from the %d available\n", n, len(l3CollapsedTemplates))
			l3CollapsedTemplate, err := client.Client.GetL3CollapsedTemplate(ctx, l3CollapsedTemplates[n].Id)
			require.NoError(t, err)
			log.Printf("    got template type '%s', ID '%s'\n", l3CollapsedTemplate.Type(), l3CollapsedTemplate.Id)

			require.Equal(t, len(templates), len(rackBasedTemplates), len(l3CollapsedTemplates)+len(podBasedTemplates)+len(l3CollapsedTemplates))
		})
	}
}

func TestGetTemplateType(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	type testData struct {
		templateId   apstra.ObjectId
		templateType apstra.TemplateType
	}

	data := []testData{
		{"pod1", apstra.TemplateTypeRackBased},
		{"L2_superspine_multi_plane", apstra.TemplateTypePodBased},
		{"L3_Collapsed_ACS", apstra.TemplateTypeL3Collapsed},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx = testutils.WrapCtxWithTestId(t, ctx)

			for _, d := range data {
				t.Run(d.templateId.String(), func(t *testing.T) {
					t.Parallel()
					ctx = testutils.WrapCtxWithTestId(t, ctx)

					ttype, err := client.Client.GetTemplateType(ctx, d.templateId)
					require.NoError(t, err)
					require.Equal(t, d.templateType, ttype)
				})
			}
		})
	}
}

func TestGetTemplateIdsTypesByName(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	templateName := testutils.RandString(10, "hex")
	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx = testutils.WrapCtxWithTestId(t, ctx)

			// fetch all template IDs
			templateIds, err := client.Client.ListAllTemplateIds(ctx)
			require.NoError(t, err)

			// choose a random template for cloning
			cloneMeId := templateIds[rand.Intn(len(templateIds))]
			cloneMeType, err := client.Client.GetTemplateType(ctx, cloneMeId)
			require.NoError(t, err)

			log.Printf("cloning template '%s' (%s) for this test", cloneMeId, cloneMeType)

			cloneCount := rand.Intn(5) + 2
			cloneIds := make([]apstra.ObjectId, cloneCount)
			for i := 0; i < cloneCount; i++ {
				switch cloneMeType {
				case apstra.TemplateTypeRackBased:
					cloneMe, err := client.Client.GetRackBasedTemplate(ctx, cloneMeId)
					if err != nil {
						t.Fatal(err)
					}
					id, err := client.Client.CreateRackBasedTemplate(ctx, &apstra.CreateRackBasedTemplateRequest{
						DisplayName: fmt.Sprintf("%s-%d", templateName, i),
						Spine: &apstra.TemplateElementSpineRequest{
							Count:                  cloneMe.Data.Spine.Count,
							LinkPerSuperspineSpeed: cloneMe.Data.Spine.LinkPerSuperspineSpeed,
							LogicalDevice:          cloneMe.Data.Spine.LogicalDevice,
							LinkPerSuperspineCount: 0,
							Tags:                   nil,
						},
						RackTypes:            cloneMe.Data.RackTypes,
						RackTypeCounts:       cloneMe.Data.RackTypeCounts,
						DhcpServiceIntent:    cloneMe.Data.DhcpServiceIntent,
						AntiAffinityPolicy:   cloneMe.Data.AntiAffinityPolicy,
						AsnAllocationPolicy:  cloneMe.Data.AsnAllocationPolicy,
						VirtualNetworkPolicy: cloneMe.Data.VirtualNetworkPolicy,
					})
					if err != nil {
						t.Fatal(err)
					}
					cloneIds[i] = id
				case apstra.TemplateTypePodBased:
					cloneMe, err := client.Client.GetPodBasedTemplate(ctx, cloneMeId)
					if err != nil {
						t.Fatal(err)
					}
					id, err := client.Client.CreatePodBasedTemplate(ctx, &apstra.CreatePodBasedTemplateRequest{
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
				case apstra.TemplateTypeL3Collapsed:
					cloneMe, err := client.Client.GetL3CollapsedTemplate(ctx, cloneMeId)
					if err != nil {
						t.Fatal(err)
					}
					id, err := client.Client.CreateL3CollapsedTemplate(ctx, &apstra.CreateL3CollapsedTemplateRequest{
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

			templateIdsToType := make(map[ObjectId]TemplateType)
			for i := 0; i < cloneCount; i++ {
				log.Printf("testing getTemplateIdsTypesByName(%s) against %s %s (%s)", templateName, client.ClientType, clientName, client.Client.ApiVersion())
				temp, err := client.Client.getTemplateIdsTypesByName(ctx, fmt.Sprintf("%s-%d", templateName, i))
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
				parsed, err := cloneMeType.parse()
				if err != nil {
					t.Fatal(err)
				}
				if parsed != v.Int() {
					t.Fatalf("expected %d, got %d", parsed, v.Int())
				}
			}

			for i, cloneId := range cloneIds {
				name := fmt.Sprintf("%s-%d", templateName, i)
				if i+1 == len(cloneIds) { // last one before they're all deleted
					log.Printf("testing getTemplateIdTypeByName(%s) against %s %s (%s)", name, client.ClientType, clientName, client.Client.ApiVersion())
					tId, tType, err := client.Client.getTemplateIdTypeByName(ctx, name)
					if err != nil {
						t.Fatal(err)
					}
					if cloneId != tId {
						t.Fatalf("expected template id '%s', got '%s'", cloneIds, tId)
					}
					if cloneMeType != templateType(tType.String()) {
						t.Fatalf("expected template type '%s', got '%s'", cloneMeType, tType.String())
					}

				}
				log.Printf("deleting clone '%s'", cloneId)
				err = client.Client.DeleteTemplate(ctx, cloneId)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
