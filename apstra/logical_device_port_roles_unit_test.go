// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra_test

import (
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/stretchr/testify/require"
)

func TestLogicalDevicePortRoles_Strings(t *testing.T) {
	type testCase struct {
		r apstra.LogicalDevicePortRoles
		e []string
	}

	var all apstra.LogicalDevicePortRoles
	all.SetAll()

	testCases := map[string]testCase{
		"none": {},
		"generic_only": {
			r: apstra.LogicalDevicePortRoles{enum.PortRoleGeneric},
			e: []string{"generic"},
		},
		"all": {
			r: all,
			e: []string{"access", "generic", "leaf", "peer", "spine", "superspine", "unused"},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := tCase.r.Strings()
			sort.Strings(result)
			sort.Strings(tCase.e)
			require.Equal(t, tCase.e, result)
		})
	}
}

func TestLogicalDevicePortRoles_FromStrings(t *testing.T) {
	type testCase struct {
		s []string
		e apstra.LogicalDevicePortRoles
	}

	var all apstra.LogicalDevicePortRoles
	all.SetAll()

	testCases := map[string]testCase{
		"none": {},
		"generic_only": {
			s: []string{"generic"},
			e: apstra.LogicalDevicePortRoles{enum.PortRoleGeneric},
		},
		"all": {
			s: all.Strings(),
			e: apstra.LogicalDevicePortRoles(enum.PortRoles.Members()),
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var roles apstra.LogicalDevicePortRoles
			err := roles.FromStrings(tCase.s)
			require.NoError(t, err)
			if len(roles) != 0 || len(tCase.s) != 0 {
				require.Equal(t, tCase.e, roles)
			}
		})
	}
}

func TestLogicalDevicePortRoles_Sort(t *testing.T) {
	type testCase struct {
		data     apstra.LogicalDevicePortRoles
		expected apstra.LogicalDevicePortRoles
	}

	testCases := map[string]testCase{
		"empty": {},
		"presorted": {
			data:     apstra.LogicalDevicePortRoles{enum.PortRoleAccess, enum.PortRoleGeneric, enum.PortRoleSpine},
			expected: apstra.LogicalDevicePortRoles{enum.PortRoleAccess, enum.PortRoleGeneric, enum.PortRoleSpine},
		},
		"unsorted": {
			data:     apstra.LogicalDevicePortRoles{enum.PortRoleSpine, enum.PortRoleGeneric, enum.PortRoleAccess},
			expected: apstra.LogicalDevicePortRoles{enum.PortRoleAccess, enum.PortRoleGeneric, enum.PortRoleSpine},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			tCase.data.Sort()
			require.Equal(t, tCase.expected, tCase.data)
		})
	}
}
