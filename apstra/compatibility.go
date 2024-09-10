package apstra

import (
	"github.com/hashicorp/go-version"
)

const (
	apstra420  = "4.2.0"
	apstra421  = "4.2.1"
	apstra4211 = "4.2.1.1"
	apstra422  = "4.2.2"
	apstra500  = "5.0.0"
	apstra510  = "5.1.0"
)

var (
	eqApstra420 = version.MustConstraints(version.NewConstraint(apstra420))
	geApstra421 = version.MustConstraints(version.NewConstraint(">=" + apstra421))
	geApstra500 = version.MustConstraints(version.NewConstraint(">=" + apstra500))

	fabricSettingsApiOk                       = geApstra421
	patchNodeSupportsUnsafeArg                = geApstra500
	templateRequestRequiresAntiAffinityPolicy = eqApstra420
)

// SupportedApiVersions returns []string with each element representing an Apstra version number like "4.2.0"
func SupportedApiVersions() []string {
	return []string{
		apstra420,
		apstra421,
		apstra4211,
		apstra422,
	}
}

func supportedApiVersionsAsConstraints() []version.Constraints {
	s := SupportedApiVersions()
	result := make([]version.Constraints, len(s))
	for i, v := range s {
		result[i] = version.MustConstraints(version.NewConstraint(v))
	}
	return result
}

type StringSliceWithIncludes []string

func (o StringSliceWithIncludes) Includes(s string) bool {
	for _, test := range o {
		if s == test {
			return true
		}
	}
	return false
}
