// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package compatibility

import (
	"github.com/hashicorp/go-version"
)

const (
	apstra420  = "4.2.0"
	apstra421  = "4.2.1"
	apstra4211 = "4.2.1.1"
	apstra422  = "4.2.2"
	apstra500  = "5.0.0"
	apstra501  = "5.0.1"
	apstra510  = "5.1.0"
	apstra600  = "6.0.0"
	apstra610  = "6.1.0"
)

var (
	// Todo: find usages of these constraints, replace them with appropriately named compatibility.Constraints
	EqApstra420 = version.MustConstraints(version.NewConstraint(apstra420))
	GeApstra421 = version.MustConstraints(version.NewConstraint(">=" + apstra421))
	LeApstra500 = version.MustConstraints(version.NewConstraint("<=" + apstra500))
)

// SupportedApiVersions returns []string with each element representing an Apstra version number like "4.2.0"
func SupportedApiVersions() []string {
	return []string{
		apstra420,
		apstra421,
		apstra4211,
		apstra422,
		apstra500,
		apstra501,
		apstra510,
		apstra600,
		apstra610,
	}
}
