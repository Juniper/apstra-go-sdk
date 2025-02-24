// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

func TestGetFeatures(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	require.NoError(t, err)

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
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			for clientName, client := range clients {
				clientName, client := clientName, client
				t.Run(fmt.Sprintf("%s_%s", client.client.apiVersion, clientName), func(t *testing.T) {
					t.Parallel()

					// true with no constraints present; defaults false when constraints exist
					versionIsPermitted := len(tCase.allowedVersions) == 0

					for _, allowedVersion := range tCase.allowedVersions {
						if allowedVersion.Check(client.client.apiVersion) {
							versionIsPermitted = true
							break // because we've found a sign that it's okay to run
						}
					}

					if !versionIsPermitted {
						t.Skipf("skipping Apstra %s", client.client.apiVersion)
					}

					// test cached values
					require.Equalf(t, tCase.expEnabled, client.client.FeatureEnabled(tCase.feature), "feature enabled")
					require.Equalf(t, tCase.expExists, client.client.FeatureExists(tCase.feature), "feature exists")

					// refresh feature cache
					require.NoError(t, client.client.getFeatures(ctx))

					// test refreshed values
					require.Equalf(t, tCase.expEnabled, client.client.FeatureEnabled(tCase.feature), "feature enabled")
					require.Equalf(t, tCase.expExists, client.client.FeatureExists(tCase.feature), "feature exists")
				})
			}
		})
	}
}
