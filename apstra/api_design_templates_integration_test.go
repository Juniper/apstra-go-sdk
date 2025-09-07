// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestGetTemplate(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

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
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			var n int

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
			log.Printf("    got template type '%s', ID '%s'\n", rackBasedTemplate.Type(), rackBasedTemplate.Id)

			// pod-based templates
			podBasedTemplates, err := client.Client.GetAllPodBasedTemplates(ctx)
			require.NoError(t, err)
			log.Printf("    got %d pod-based templates\n", len(podBasedTemplates))

			n = rand.Intn(len(podBasedTemplates))
			log.Printf("  using randomly-selected index %d from the %d available\n", n, len(podBasedTemplates))
			podBasedTemplate, err := client.Client.GetPodBasedTemplate(ctx, podBasedTemplates[n].Id)
			require.NoError(t, err)
			log.Printf("    got template type '%s', ID '%s'\n", podBasedTemplate.Type(), podBasedTemplate.Id)

			// l3-collapsed templates
			l3CollapsedTemplates, err := client.Client.GetAllL3CollapsedTemplates(ctx)
			require.NoError(t, err)
			log.Printf("  got %d pod-based templates\n", len(l3CollapsedTemplates))

			n = rand.Intn(len(l3CollapsedTemplates))
			log.Printf("  using randomly-selected index %d from the %d available\n", n, len(l3CollapsedTemplates))
			l3CollapsedTemplate, err := client.Client.GetL3CollapsedTemplate(ctx, l3CollapsedTemplates[n].Id)
			require.NoError(t, err)
			log.Printf("    got template type '%s', ID '%s'\n", l3CollapsedTemplate.Type(), l3CollapsedTemplate.Id)

			require.Equal(t, len(templates), len(rackBasedTemplates)+len(podBasedTemplates)+len(l3CollapsedTemplates))
		})
	}
}

func TestGetTemplateType(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
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
			ctx := testutils.ContextWithTestID(t, ctx)

			for _, d := range data {
				t.Run(d.templateId.String(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(t, ctx)

					ttype, err := client.Client.GetTemplateType(ctx, d.templateId)
					require.NoError(t, err)
					require.Equal(t, d.templateType, ttype)
				})
			}
		})
	}
}

func TestCRUDTemplates(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct { // fill only one template type
		rackBasedTemplate   *apstra.CreateRackBasedTemplateRequest
		podBasedTemplate    *apstra.CreatePodBasedTemplateRequest
		l3CollapsedTemplate *apstra.CreateL3CollapsedTemplateRequest
	}

	exactlyOneTemplate := func(tc testCase) bool {
		count := 0
		if tc.rackBasedTemplate != nil {
			count++
		}
		if tc.podBasedTemplate != nil {
			count++
		}
		if tc.l3CollapsedTemplate != nil {
			count++
		}
		return count == 1
	}

	testCases := map[string]testCase{
		"simple_rack_based": {
			rackBasedTemplate: &apstra.CreateRackBasedTemplateRequest{
				DisplayName: testutils.RandString(6, "hex"),
				Spine: &apstra.TemplateElementSpineRequest{
					Count:                  1,
					LinkPerSuperspineSpeed: "10G",
					LogicalDevice:          "AOS-7x10-Spine",
					LinkPerSuperspineCount: 1,
				},
				RackInfos: map[apstra.ObjectId]apstra.TemplateRackBasedRackInfo{
					"L2_Virtual_MLAG_2x_Links": {Count: 1},
				},
				DhcpServiceIntent: nil,
				AntiAffinityPolicy: &apstra.AntiAffinityPolicy{
					Algorithm: apstra.AlgorithmHeuristic,
					Mode:      apstra.AntiAffinityModeDisabled,
				},
				AsnAllocationPolicy:  &apstra.AsnAllocationPolicy{SpineAsnScheme: apstra.AsnAllocationSchemeDistinct},
				VirtualNetworkPolicy: &apstra.VirtualNetworkPolicy{OverlayControlProtocol: apstra.OverlayControlProtocolEvpn},
			},
		},
	}

	for tName, tCase := range testCases {
		require.True(t, exactlyOneTemplate(tCase))
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(t, ctx)

					var id apstra.ObjectId
					var err error
					var expectedType apstra.TemplateType
					var expectedName string
					switch {
					case tCase.rackBasedTemplate != nil:
						id, err = client.Client.CreateRackBasedTemplate(ctx, tCase.rackBasedTemplate)
						expectedType = apstra.TemplateTypeRackBased
						expectedName = tCase.rackBasedTemplate.DisplayName
					case tCase.podBasedTemplate != nil:
						id, err = client.Client.CreatePodBasedTemplate(ctx, tCase.podBasedTemplate)
						expectedType = apstra.TemplateTypePodBased
						expectedName = tCase.podBasedTemplate.DisplayName
					case tCase.l3CollapsedTemplate != nil:
						id, err = client.Client.CreateL3CollapsedTemplate(ctx, tCase.l3CollapsedTemplate)
						expectedType = apstra.TemplateTypeL3Collapsed
						expectedName = tCase.l3CollapsedTemplate.DisplayName
					default:
						t.Fatal("we should never get here")
					}
					require.NoError(t, err)

					all, err := client.Client.GetAllTemplates(ctx)
					require.NoError(t, err)
					var template apstra.Template
					for i := range all {
						if all[i].ID() == id {
							template = all[i]
						}
					}
					require.NotNil(t, template)
					require.Equal(t, expectedType, template.Type())

					tType, err := client.Client.GetTemplateType(ctx, id)
					require.NoError(t, err)
					require.Equal(t, expectedType, tType)

					idToType, err := client.Client.GetTemplateIdsTypesByName(ctx, expectedName)
					require.NoError(t, err)
					require.Contains(t, idToType, id)
					require.Equal(t, expectedType, idToType[id])

					tId, tType, err := client.Client.GetTemplateIdTypeByName(ctx, expectedName)
					require.NoError(t, err)
					require.Equal(t, id, tId)
					require.Equal(t, expectedType, tType)

					tType, err = client.Client.GetTemplateType(ctx, id)
					require.NoError(t, err)
					require.Equal(t, expectedType, tType)

					switch expectedType {
					case apstra.TemplateTypeRackBased:
						template, err = client.Client.GetRackBasedTemplate(ctx, id)
					case apstra.TemplateTypePodBased:
						template, err = client.Client.GetPodBasedTemplate(ctx, id)
					case apstra.TemplateTypeL3Collapsed:
						template, err = client.Client.GetL3CollapsedTemplate(ctx, id)
					}
					require.NoError(t, err)
					require.Equal(t, expectedType, template.Type())
					require.Equal(t, id, template.ID())

					err = client.Client.DeleteTemplate(ctx, id)
					require.NoError(t, err)
				})
			}
		})
	}
}
