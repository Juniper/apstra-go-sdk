package compatibility

import (
	"strings"

	"github.com/hashicorp/go-version"
)

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
		permitAny:               true,
	}
)
