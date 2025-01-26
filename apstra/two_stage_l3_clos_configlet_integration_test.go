// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestImportGetUpdateGetDeleteConfiglet(t *testing.T) {
	ctx := context.Background()

	compareGenerators := func(t *testing.T, a, b ConfigletGenerator) {
		require.Equal(t, a.ConfigStyle.String(), b.ConfigStyle.String())
		require.Equal(t, a.Section.String(), b.Section.String())
		require.Equal(t, a.SectionCondition, b.SectionCondition)
		require.Equal(t, a.TemplateText, b.TemplateText)
		require.Equal(t, a.NegationTemplateText, b.NegationTemplateText)
		require.Equal(t, a.Filename, b.Filename)
	}

	compare := func(t *testing.T, a, b *TwoStageL3ClosConfigletData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.Label, b.Label)
		require.Equal(t, a.Condition, b.Condition)
		require.NotNil(t, a.Data)
		require.NotNil(t, b.Data)
		require.Equal(t, a.Data.DisplayName, b.Data.DisplayName)
		require.Equal(t, len(a.Data.Generators), len(b.Data.Generators))
		for i := range a.Data.Generators {
			compareGenerators(t, a.Data.Generators[i], b.Data.Generators[i])
		}
	}

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	type testStep struct {
		data TwoStageL3ClosConfigletData
	}

	type testCase struct {
		versionConstraint version.Constraints
		steps             []testStep
	}

	testCases := map[string]testCase{
		"simple_test": {
			steps: []testStep{
				{
					data: TwoStageL3ClosConfigletData{
						Label:     randString(6, "hex"),
						Condition: fmt.Sprintf(`label in ["%s"]`, randString(6, "hex")),
						Data: &ConfigletData{
							Generators: []ConfigletGenerator{
								{
									ConfigStyle:          enum.ConfigletStyleJunos,
									Section:              enum.ConfigletSectionSystem,
									TemplateText:         "set " + randString(6, "hex"),
									NegationTemplateText: "del " + randString(6, "hex"),
								},
							},
							DisplayName: randString(6, "hex"),
						},
					},
				},
				{
					data: TwoStageL3ClosConfigletData{
						Label:     randString(6, "hex"),
						Condition: fmt.Sprintf(`label in ["%s"]`, randString(6, "hex")),
						Data: &ConfigletData{
							Generators: []ConfigletGenerator{
								{
									ConfigStyle:          enum.ConfigletStyleEos,
									Section:              enum.ConfigletSectionInterface,
									SectionCondition:     `role in ["spine_leaf"]`,
									TemplateText:         "set " + randString(6, "hex"),
									NegationTemplateText: "del " + randString(6, "hex"),
								},
							},
							DisplayName: randString(6, "hex"),
						},
					},
				},
				{
					data: TwoStageL3ClosConfigletData{
						Label:     randString(6, "hex"),
						Condition: fmt.Sprintf("label in [%q]", randString(6, "hex")),
						Data: &ConfigletData{
							Generators: []ConfigletGenerator{
								{
									ConfigStyle:          enum.ConfigletStyleEos,
									Section:              enum.ConfigletSectionOspf,
									TemplateText:         "set " + randString(6, "hex"),
									NegationTemplateText: "del " + randString(6, "hex"),
								},
							},
							DisplayName: randString(6, "hex"),
						},
					},
				},
			},
		},
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			bp := testBlueprintA(ctx, t, client.client)

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Parallel()

					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(client.client.apiVersion) {
						t.Skipf("skipping %q due to version constraints: %q. API version: %q",
							tName, tCase.versionConstraint, client.client.apiVersion)
					}

					require.GreaterOrEqualf(t, len(tCase.steps), 1, "test %s has no test steps!", tName)

					log.Printf("testing CreateConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					id, err := bp.CreateConfiglet(ctx, &tCase.steps[0].data)
					require.NoError(t, err)
					require.NotEmpty(t, id)

					log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					configlet, err := bp.GetConfiglet(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, configlet.Id)
					compare(t, &tCase.steps[0].data, configlet.Data)

					log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					configlet, err = bp.GetConfigletByName(ctx, tCase.steps[0].data.Label)
					require.NoError(t, err)
					require.Equal(t, id, configlet.Id)
					compare(t, &tCase.steps[0].data, configlet.Data)

					log.Printf("testing GetAllConfiglets() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					configlets, err := bp.GetAllConfiglets(ctx)
					require.NoError(t, err)
					configlet = nil
					for _, c := range configlets {
						if c.Id == id {
							configlet = &c
							break
						}
					}
					require.NotNil(t, configlet)
					require.Equal(t, id, configlet.Id)
					compare(t, &tCase.steps[0].data, configlet.Data)

					for _, step := range tCase.steps {
						log.Printf("testing UpdateConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						err = bp.UpdateConfiglet(ctx, id, &step.data)
						require.NoError(t, err)

						log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						configlet, err = bp.GetConfiglet(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, configlet.Id)
						compare(t, &step.data, configlet.Data)

						log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						configlet, err = bp.GetConfigletByName(ctx, step.data.Label)
						require.NoError(t, err)
						require.Equal(t, id, configlet.Id)
						compare(t, &step.data, configlet.Data)

						log.Printf("testing GetAllConfigletIds() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						ids, err := bp.GetAllConfigletIds(ctx)
						require.NoError(t, err)
						require.Contains(t, ids, id)

						log.Printf("testing GetAllConfiglets() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						configlets, err = bp.GetAllConfiglets(ctx)
						require.NoError(t, err)
						configlet = nil
						for _, c := range configlets {
							if c.Id == id {
								configlet = &c
								break
							}
						}
						require.NotNil(t, configlet)
						require.Equal(t, id, configlet.Id)
						compare(t, &step.data, configlet.Data)
					}

					log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					err = bp.DeleteConfiglet(ctx, id)
					require.NoError(t, err)

					var ace ClientErr

					log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					_, err = bp.GetConfiglet(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)

					log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					_, err = bp.GetConfigletByName(ctx, tCase.steps[len(tCase.steps)-1].data.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)

					log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					err = bp.DeleteConfiglet(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)
				})
			}
		})
	}
}
