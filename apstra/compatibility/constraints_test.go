// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package compatibility_test

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
)

func TestConstraints(t *testing.T) {
	type testCase struct {
		constraint compatibility.Constraint
		version    string
		expected   bool
	}

	testCases := map[string]testCase{
		"FabricSettingsApiOk_4.2.0": {
			constraint: compatibility.FabricSettingsApiOk,
			version:    "4.2.0",
			expected:   false,
		},
		"FabricSettingsApiOk_4.2.1": {
			constraint: compatibility.FabricSettingsApiOk,
			version:    "4.2.1",
			expected:   true,
		},
		"FabricSettingsApiOk_4.2.2": {
			constraint: compatibility.FabricSettingsApiOk,
			version:    "4.2.2",
			expected:   true,
		},
		"FabricSettingsApiOk_5.0.0": {
			constraint: compatibility.FabricSettingsApiOk,
			version:    "5.0.0",
			expected:   true,
		},
		"FabricSettingsApiOk_5.0.0a-6": {
			constraint: compatibility.FabricSettingsApiOk,
			version:    "5.0.0a-6",
			expected:   true,
		},
		"PatchNodeSupportsUnsafeArg_4.99.99": {
			constraint: compatibility.PatchNodeSupportsUnsafeArg,
			version:    "4.99.99",
			expected:   false,
		},
		"PatchNodeSupportsUnsafeArg_5.0.0": {
			constraint: compatibility.PatchNodeSupportsUnsafeArg,
			version:    "5.0.0",
			expected:   true,
		},
		"PatchNodeSupportsUnsafeArg_5.0.1": {
			constraint: compatibility.PatchNodeSupportsUnsafeArg,
			version:    "5.0.1",
			expected:   true,
		},
		"PatchNodeSupportsUnsafeArg_5.0.0a-6": {
			constraint: compatibility.PatchNodeSupportsUnsafeArg,
			version:    "5.0.0a-6",
			expected:   true,
		},
		"CheckTemplateRequestRequiresAntiAffinityPolicy_4.1.2": {
			constraint: compatibility.TemplateRequestRequiresAntiAffinityPolicy,
			version:    "4.1.2",
			expected:   true,
		},
		"CheckTemplateRequestRequiresAntiAffinityPolicy_4.2.0": {
			constraint: compatibility.TemplateRequestRequiresAntiAffinityPolicy,
			version:    "4.2.0",
			expected:   true,
		},
		"CheckTemplateRequestRequiresAntiAffinityPolicy_4.2.1": {
			constraint: compatibility.TemplateRequestRequiresAntiAffinityPolicy,
			version:    "4.2.1",
			expected:   false,
		},
		"CheckTemplateRequestRequiresAntiAffinityPolicy_5.0.0": {
			constraint: compatibility.TemplateRequestRequiresAntiAffinityPolicy,
			version:    "5.0.0",
			expected:   false,
		},
		"CheckTemplateRequestRequiresAntiAffinityPolicy_5.0.0a-6": {
			constraint: compatibility.TemplateRequestRequiresAntiAffinityPolicy,
			version:    "5.0.0a-6",
			expected:   false,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			v, err := version.NewVersion(tCase.version)
			require.NoError(t, err)

			var msg1, msg2 string
			if tCase.expected {
				msg1 = "expected"
				msg2 = "but it does not"
			} else {
				msg1 = "did not expect"
				msg2 = "but it does"
			}

			result := tCase.constraint.Check(v)
			require.Equalf(t, tCase.expected, result, "%s version %s to satisfy constraint %s %s.", msg1, tCase.version, tCase.constraint, msg2)
		})
	}
}
