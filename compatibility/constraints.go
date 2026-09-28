// Copyright (c) Juniper Networks, Inc., 2024-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package compatibility

import (
	"strings"

	"github.com/hashicorp/go-version"
)

var (
	DatacenterPolicyAddressFamilyNotSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<=" + apstra612)),
	}
	DatacenterPolicyAddressFamilyRequired = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">" + apstra612)),
	}
	DatacenterSwitchingZoneOK = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra620)),
	}
	FabricSettingsIPv6EnabledOK = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra610)),
	}
	FabricSettingsDefaultAnycastGWMacOK = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra610)),
	}
	HasDeviceOsImageDownloadTimeout = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra510)),
	}
	RailCollapsedSupport = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra600)),
	}
	PolicyIDErrorHasNoQuotes = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra620)),
	}
	SecurityZoneAddressingSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra610)),
	}
	ServerVersionSupported = Constraint{
		constraints:             version.MustConstraints(version.NewConstraint(strings.Join(SupportedApiVersions(), ","))),
		considerPreReleaseLabel: true,
		permitAny:               true,
	}
	SwitchSystemLinksSystemIDNotSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint("<" + apstra610)),
	}
	SwitchingZoneSupported = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra620)),
	}
	VirtualNetworkAPITags = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra620)),
	}
	VirtualNetworkEncapsulateInnerVLANOK = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra620)),
	}
	VirtualNetworkZoneIDsRequired = Constraint{
		constraints: version.MustConstraints(version.NewConstraint(">=" + apstra620)),
	}
)
