// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateGetDeleteConfiglet(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())

	clients := testclient.GetTestClients(t, ctx)

	type testStep struct {
		data apstra.ConfigletData
	}

	type testCase struct {
		versionConstraint version.Constraints
		steps             []testStep
	}

	testCases := map[string]testCase{
		"simple_test": {
			steps: []testStep{
				{
					data: apstra.ConfigletData{
						DisplayName: testutils.RandString(6, "hex"),
						RefArchs:    []enum.RefDesign{enum.RefDesignDatacenter},
						Generators: []apstra.ConfigletGenerator{
							{
								ConfigStyle:  enum.ConfigletStyleJunos,
								Section:      enum.ConfigletSectionSystem,
								TemplateText: "set " + testutils.RandString(6, "hex"),
							},
						},
					},
				},
				{
					data: apstra.ConfigletData{
						DisplayName: testutils.RandString(6, "hex"),
						RefArchs:    []enum.RefDesign{enum.RefDesignDatacenter},
						Generators: []apstra.ConfigletGenerator{
							{
								ConfigStyle:          enum.ConfigletStyleNxos,
								Section:              enum.ConfigletSectionInterface,
								TemplateText:         "set " + testutils.RandString(6, "hex"),
								NegationTemplateText: "no " + testutils.RandString(6, "hex"),
							},
							{
								ConfigStyle:          enum.ConfigletStyleEos,
								Section:              enum.ConfigletSectionInterface,
								TemplateText:         "set " + testutils.RandString(6, "hex"),
								NegationTemplateText: "no " + testutils.RandString(6, "hex"),
							},
							{
								ConfigStyle:  enum.ConfigletStyleSonic,
								Section:      enum.ConfigletSectionFile,
								TemplateText: "set " + testutils.RandString(6, "hex"),
								Filename:     "/etc/" + testutils.RandString(6, "hex"),
							},
						},
					},
				},
				{
					data: apstra.ConfigletData{
						DisplayName: testutils.RandString(6, "hex"),
						RefArchs:    []enum.RefDesign{enum.RefDesignDatacenter},
						Generators: []apstra.ConfigletGenerator{
							{
								ConfigStyle:  enum.ConfigletStyleJunos,
								Section:      enum.ConfigletSectionSystem,
								TemplateText: "set " + testutils.RandString(6, "hex"),
							},
						},
					},
				},
			},
		},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Parallel()

					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(client.APIVersion()) {
						t.Skipf("skipping %q due to version constraints: %q. API version: %q",
							tName, tCase.versionConstraint, client.Client.ApiVersion())
					}

					require.GreaterOrEqualf(t, len(tCase.steps), 1, "test %s has no test steps!", tName)

					id, err := client.Client.CreateConfiglet(ctx, &tCase.steps[0].data)
					require.NoError(t, err)
					require.NotEmpty(t, id)

					configlet, err := client.Client.GetConfiglet(ctx, id)
					require.NoError(t, err)
					require.Equal(t, id, configlet.Id)
					compare.ConfigletData(t, &tCase.steps[0].data, configlet.Data)

					configlet, err = client.Client.GetConfigletByName(ctx, tCase.steps[0].data.DisplayName)
					require.NoError(t, err)
					require.Equal(t, id, configlet.Id)
					compare.ConfigletData(t, &tCase.steps[0].data, configlet.Data)

					configlets, err := client.Client.GetAllConfiglets(ctx)
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
					compare.ConfigletData(t, &tCase.steps[0].data, configlet.Data)

					for _, step := range tCase.steps {
						err = client.Client.UpdateConfiglet(ctx, id, &step.data)
						require.NoError(t, err)

						configlet, err = client.Client.GetConfiglet(ctx, id)
						require.NoError(t, err)
						require.Equal(t, id, configlet.Id)
						compare.ConfigletData(t, &step.data, configlet.Data)

						configlet, err = client.Client.GetConfigletByName(ctx, step.data.DisplayName)
						require.NoError(t, err)
						require.Equal(t, id, configlet.Id)
						compare.ConfigletData(t, &step.data, configlet.Data)

						configlets, err = client.Client.GetAllConfiglets(ctx)
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
						compare.ConfigletData(t, &step.data, configlet.Data)
					}

					err = client.Client.DeleteConfiglet(ctx, id)
					require.NoError(t, err)

					var ace apstra.ClientErr

					_, err = client.Client.GetConfiglet(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), apstra.ErrNotfound)

					_, err = client.Client.GetConfigletByName(ctx, tCase.steps[len(tCase.steps)-1].data.DisplayName)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), apstra.ErrNotfound)

					err = client.Client.DeleteConfiglet(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, ace.Type(), apstra.ErrNotfound)
				})
			}
		})
	}
}
