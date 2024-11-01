// Copyright (c) Juniper Networks, Inc., 2024-2024.
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
)

var (
	// Todo: find usages of these constraints, replace them with appropriately named compatibility.Constraints
	EqApstra420  = version.MustConstraints(version.NewConstraint(apstra420))
	EqApstra421  = version.MustConstraints(version.NewConstraint(apstra421))
	EqApstra4211 = version.MustConstraints(version.NewConstraint(apstra4211))
	EqApstra422  = version.MustConstraints(version.NewConstraint(apstra422))

	GeApstra421 = version.MustConstraints(version.NewConstraint(">=" + apstra421))
	GeApstra500 = version.MustConstraints(version.NewConstraint(">=" + apstra500))
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
	}
}

// Make sure the version is 5.0.0 5.1, 5.0.0A etc

func OnlyFiveAndAbove(v *version.Version) bool {
	switch {
	case EqApstra420.Check(v):
	case EqApstra421.Check(v):
	case EqApstra4211.Check(v):
	case EqApstra422.Check(v):
		return false
	default:
		return true
	}
	return true
}
