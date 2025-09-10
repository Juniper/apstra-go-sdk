// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestGetFeatures(t *testing.T) {
	ctx := testutils.ContextWithTestID(t, context.Background())
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		allowedVersions []version.Constraints
		feature         enum.ApiFeature
		expExists       bool
		expEnabled      bool
		expNotEnabled   bool
	}

	testCases := map[string]testCase{
		"task_api": {
			feature:       enum.ApiFeatureTaskApi,
			expExists:     true,
			expEnabled:    true,
			expNotEnabled: false,
		},
		"ai_fabric_exists": {
			allowedVersions: []version.Constraints{
				version.MustConstraints(version.NewConstraint("5.0.0a-6")),
				version.MustConstraints(version.NewConstraint("5.0.0a-7")),
			},
			feature:       enum.ApiFeatureAiFabric,
			expExists:     true,
			expEnabled:    true,
			expNotEnabled: false,
		},
		"ai_fabric_not_exists": {
			allowedVersions: []version.Constraints{compatibility.LeApstra500},
			feature:         enum.ApiFeatureAiFabric,
			expExists:       false,
			expEnabled:      false,
			expNotEnabled:   true,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(t, ctx)

			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(t, ctx)

					// true with no constraints present; defaults false when constraints exist
					versionIsPermitted := len(tCase.allowedVersions) == 0

					for _, allowedVersion := range tCase.allowedVersions {
						if allowedVersion.Check(client.APIVersion()) {
							versionIsPermitted = true
							break // because we've found a sign that it's okay to run
						}
					}

					if !versionIsPermitted {
						t.Skipf("skipping Apstra %s", client.APIVersion())
					}

					// test cached values
					require.Equalf(t, tCase.expEnabled, client.Client.FeatureEnabled(tCase.feature), "feature enabled")
					require.Equalf(t, tCase.expExists, client.Client.FeatureExists(tCase.feature), "feature exists")

					// refresh feature cache
					require.NoError(t, client.Client.GetFeatures(ctx))

					// test refreshed values
					require.Equalf(t, tCase.expEnabled, client.Client.FeatureEnabled(tCase.feature), "feature enabled")
					require.Equalf(t, tCase.expExists, client.Client.FeatureExists(tCase.feature), "feature exists")
				})
			}
		})
	}
}
