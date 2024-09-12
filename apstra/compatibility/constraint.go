package compatibility

import (
	"strings"

	"github.com/hashicorp/go-version"
)

type Constraint struct {
	constraints             version.Constraints
	considerPreReleaseLabel bool
	permitAny               bool
}

func (o Constraint) Check(v *version.Version) bool {
	if !o.considerPreReleaseLabel {
		// drop the pre-release label
		v = v.Core()
	}

	if !o.permitAny {
		// v must satisfy all constraints
		return o.constraints.Check(v)
	}

	// v can satisfy any constraint
	for _, constraint := range o.constraints {
		if constraint.Check(v) {
			return true
		}
	}

	// v does not satisfy any constraint
	return false
}

func (o Constraint) String() string {
	result := make([]string, len(o.constraints))
	for i, constraint := range o.constraints {
		result[i] = constraint.String()
	}

	return strings.Join(result, ",")
}
