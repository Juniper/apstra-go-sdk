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
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateGetDeleteConfiglet(t *testing.T) {
	ctx := context.Background()

	compareGenerators := func(t *testing.T, a, b ConfigletGenerator) {
		require.Equal(t, a.ConfigStyle.String(), b.ConfigStyle.String())
		require.Equal(t, a.Section.String(), b.Section.String())
		require.Equal(t, a.SectionCondition, b.SectionCondition)
		require.Equal(t, a.TemplateText, b.TemplateText)
		require.Equal(t, a.NegationTemplateText, b.NegationTemplateText)
		require.Equal(t, a.Filename, b.Filename)
	}

	compare := func(t *testing.T, a, b *ConfigletData) {
		t.Helper()

		require.NotNil(t, a)
		require.NotNil(t, b)
		require.Equal(t, a.DisplayName, b.DisplayName)
		compareSlicesAsSets(t, a.RefArchs, b.RefArchs, "while comparing configlet refarchs,")
		require.Equal(t, len(a.Generators), len(b.Generators))
		for i := range a.Generators {
			compareGenerators(t, a.Generators[i], b.Generators[i])
		}
	}

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	type testStep struct {
		data ConfigletData
	}

	type testCase struct {
		versionConstraint version.Constraints
		steps             []testStep
	}

	testCases := map[string]testCase{
		"simple_test": {
			steps: []testStep{
				{
					data: ConfigletData{
						DisplayName: randString(6, "hex"),
						RefArchs:    []enum.RefDesign{enum.RefDesignDatacenter},
						Generators: []ConfigletGenerator{
							{
								ConfigStyle:  enum.ConfigletStyleJunos,
								Section:      enum.ConfigletSectionSystem,
								TemplateText: "set " + randString(6, "hex"),
							},
						},
					},
				},
				{
					data: ConfigletData{
						DisplayName: randString(6, "hex"),
						RefArchs:    []enum.RefDesign{enum.RefDesignDatacenter},
						Generators: []ConfigletGenerator{
							{
								ConfigStyle:          enum.ConfigletStyleNxos,
								Section:              enum.ConfigletSectionInterface,
								TemplateText:         "set " + randString(6, "hex"),
								NegationTemplateText: "no " + randString(6, "hex"),
							},
							{
								ConfigStyle:          enum.ConfigletStyleEos,
								Section:              enum.ConfigletSectionInterface,
								TemplateText:         "set " + randString(6, "hex"),
								NegationTemplateText: "no " + randString(6, "hex"),
							},
							{
								ConfigStyle:  enum.ConfigletStyleSonic,
								Section:      enum.ConfigletSectionFile,
								TemplateText: "set " + randString(6, "hex"),
								Filename:     "/etc/" + randString(6, "hex"),
							},
						},
					},
				},
				{
					data: ConfigletData{
						DisplayName: randString(6, "hex"),
						RefArchs:    []enum.RefDesign{enum.RefDesignDatacenter},
						Generators: []ConfigletGenerator{
							{
								ConfigStyle:  enum.ConfigletStyleJunos,
								Section:      enum.ConfigletSectionSystem,
								TemplateText: "set " + randString(6, "hex"),
							},
						},
					},
				},
			},
		},
	}

	for clientName, client := range clients {
		t.Run(client.name(), func(t *testing.T) {
			t.Parallel()

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Parallel()

					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(client.client.apiVersion) {
						t.Skipf("skipping %q due to version constraints: %q. API version: %q",
							tName, tCase.versionConstraint, client.client.apiVersion)
					}

					require.GreaterOrEqualf(t, len(tCase.steps), 1, "test %s has no test steps!", tName)

					log.Printf("testing CreateConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					id, err := client.client.CreateConfiglet(ctx, &tCase.steps[0].data)
					require.NoError(t, err)
					require.NotEmpty(t, id)

					log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					configlet, err := client.client.GetConfiglet(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, configlet.Id)
					compare(t, &tCase.steps[0].data, configlet.Data)

					log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					configlet, err = client.client.GetConfigletByName(ctx, tCase.steps[0].data.DisplayName)
					require.NoError(t, err)
					require.Equal(t, id, configlet.Id)
					compare(t, &tCase.steps[0].data, configlet.Data)

					log.Printf("testing GetAllConfiglets() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					configlets, err := client.client.GetAllConfiglets(ctx)
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
						err = client.client.UpdateConfiglet(ctx, id, &step.data)
						require.NoError(t, err)

						log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						configlet, err = client.client.GetConfiglet(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, configlet.Id)
						compare(t, &step.data, configlet.Data)

						log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						configlet, err = client.client.GetConfigletByName(ctx, step.data.DisplayName)
						require.NoError(t, err)
						require.Equal(t, id, configlet.Id)
						compare(t, &step.data, configlet.Data)

						log.Printf("testing GetAllConfiglets() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
						configlets, err = client.client.GetAllConfiglets(ctx)
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
					err = client.client.DeleteConfiglet(ctx, id)
					require.NoError(t, err)

					var ace ClientErr

					log.Printf("testing GetConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					_, err = client.client.GetConfiglet(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)

					log.Printf("testing GetConfigletByName() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					_, err = client.client.GetConfigletByName(ctx, tCase.steps[len(tCase.steps)-1].data.DisplayName)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)

					log.Printf("testing DeleteConfiglet() against %s %s (%s)", client.clientType, clientName, client.client.apiVersion)
					err = client.client.DeleteConfiglet(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), ErrNotfound)
				})
			}
		})
	}
}
