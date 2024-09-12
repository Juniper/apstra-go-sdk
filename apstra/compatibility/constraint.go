package compatibility

import (
	"strings"

	"github.com/hashicorp/go-version"
)

type Constraint struct {
	constraints             version.Constraints
	considerPreReleaseLabel bool
}

func (o Constraint) Check(v *version.Version) bool {
	if o.considerPreReleaseLabel {
		return o.constraints.Check(v)
	}

	return o.constraints.Check(v.Core())
}

func (o Constraint) String() string {
	result := make([]string, len(o.constraints))
	for i, constraint := range o.constraints {
		result[i] = constraint.String()
	}

	return strings.Join(result, ",")
}

var (
	FabricSettingsApiOk = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra421)),
	}
	PatchNodeSupportsUnsafeArg = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
	TemplateRequestRequiresAntiAffinityPolicy = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<=" + apstra420)),
	}
	ServerVersionSupported = Constraint{
		constraints:             version.MustConstraints(version.NewConstraint(strings.Join(SupportedApiVersions(), ","))),
		considerPreReleaseLabel: true,
	}
)
