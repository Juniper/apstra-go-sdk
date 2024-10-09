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
	all.IncludeAllUses()

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
	all.IncludeAllUses()

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

func TestLogicalDevicePortRoles_IncludeAllUses(t *testing.T) {
	var data apstra.LogicalDevicePortRoles
	data.IncludeAllUses()

	expected := apstra.LogicalDevicePortRoles{
		enum.PortRoleAccess,
		enum.PortRoleGeneric,
		// enum.PortRoleL3Server, <---- TEST VALIDATES THAT THIS ONE IS OMITTED
		enum.PortRoleLeaf,
		enum.PortRolePeer,
		enum.PortRoleSpine,
		enum.PortRoleSuperspine,
		enum.PortRoleUnused,
	}

	require.Equal(t, expected, data)
}

func TestLogicalDevicePortRoles_Validate(t *testing.T) {
	type testCase struct {
		roles       []string
		errContains string
	}

	testCases := map[string]testCase{
		"okay": {
			roles: []string{"access", "leaf", "spine"},
		},
		"empty": {},
		"single_err": {
			roles:       []string{"l3_server"},
			errContains: "l3_server",
		},
		"multiple_err": {
			roles:       []string{"access", "l3_server", "spine"},
			errContains: "l3_server",
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			var portRoles apstra.LogicalDevicePortRoles
			err := portRoles.FromStrings(tCase.roles)
			require.NoError(t, err)

			err = portRoles.Validate()
			if tCase.errContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tCase.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
