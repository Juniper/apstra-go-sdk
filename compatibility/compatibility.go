// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package compatibility

import (
	"github.com/hashicorp/go-version"
)

const (
	apstra500 = "5.0.0"
	apstra501 = "5.0.1"
	apstra510 = "5.1.0"
	apstra600 = "6.0.0"
	apstra610 = "6.1.0"
	apstra611 = "6.1.1"
	apstra612 = "6.1.2"
	apstra620 = "6.2.0"
)

var (
	// Todo: find usages of these constraints, replace them with appropriately named compatibility.Constraints
	LeApstra500 = version.MustConstraints(version.NewConstraint("<=" + apstra500))
)

// SupportedApiVersions returns []string with each element representing an Apstra version number like "4.2.0"
func SupportedApiVersions() []string {
	return []string{
		apstra500,
		apstra501,
		apstra510,
		apstra600,
		apstra610,
		apstra611,
		apstra612,
		apstra620,
	}
}
