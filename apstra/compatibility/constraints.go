package compatibility

import (
	"strings"

	"github.com/hashicorp/go-version"
)

var (
	BpHasFabricAddressingPolicyNode = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<=" + apstra420)),
	}
	BpHasVirtualNetworkPolicyNode = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<=" + apstra420)),
	}
	FabricSettingsApiOk = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra421)),
	}
	IbaDashboardSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra500)),
	}
	IbaProbeSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra500)),
	}
	IbaWidgetSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra500)),
	}
	PatchNodeSupportsUnsafeArg = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
	TemplateRequestRequiresAntiAffinityPolicy = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<=" + apstra420)),
	}
	RoutingPolicyExportHasL3EdgeLinks = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra500)),
	}
	ServerVersionSupported = Constraint{
		constraints:             version.MustConstraints(version.NewConstraint(strings.Join(SupportedApiVersions(), ","))),
		considerPreReleaseLabel: true,
		permitAny:               true,
	}
	SystemManagerHasSkipInterfaceShutdownOnUpgrade = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
)
