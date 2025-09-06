// Copyright (c) Juniper Networks, Inc., 2024-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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
	EmptyVnBindingsOk = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
	FabricSettingsApiOk = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra421)),
	}
	HasDeviceOsImageDownloadTimeout = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra510)),
	}
	IbaDashboardSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
	PatchNodeSupportsUnsafeArg = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
	RailCollapsedSupport = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra600)),
	}
	RoutingPolicyExportHasL3EdgeLinks = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra500)),
	}
	SecurityZoneLoopbackApiSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
	ServerVersionSupported = Constraint{
		constraints:             version.MustConstraints(version.NewConstraint(strings.Join(SupportedApiVersions(), ","))),
		considerPreReleaseLabel: true,
		permitAny:               true,
	}
	SystemManagerHasSkipInterfaceShutdownOnUpgrade = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
	TemplateRequestRequiresAntiAffinityPolicy = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<=" + apstra420)),
	}
	VirtualNetworkTags = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra500)),
	}
)
