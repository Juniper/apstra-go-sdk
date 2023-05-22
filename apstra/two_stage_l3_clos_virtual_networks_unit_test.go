//go:build integration
// +build integration

package apstra

import (
	"testing"
)

func TestTwoStageL3ClosVirtualNetworkStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "", intType: SviIpRequirementNone, stringType: sviIpRequirementNone},
		{stringVal: "optional", intType: SviIpRequirementOptional, stringType: sviIpRequirementOptional},
		{stringVal: "forbidden", intType: SviIpRequirementForbidden, stringType: sviIpRequirementForbidden},
		{stringVal: "mandatory", intType: SviIpRequirementMandatory, stringType: sviIpRequirementMandatory},
		{stringVal: "intention_conflict", intType: SviIpRequirementIntentionConflict, stringType: sviIpRequirementIntentionConflict},

		{stringVal: "", intType: Ipv4ModeNone, stringType: ipv4ModeNone},
		{stringVal: "disabled", intType: Ipv4ModeDisabled, stringType: ipv4ModeDisabled},
		{stringVal: "enabled", intType: Ipv4ModeEnabled, stringType: ipv4ModeEnabled},
		{stringVal: "forced", intType: Ipv4ModeForced, stringType: ipv4ModeForced},

		{stringVal: "", intType: Ipv6ModeNone, stringType: ipv6ModeNone},
		{stringVal: "disabled", intType: Ipv6ModeDisabled, stringType: ipv6ModeDisabled},
		{stringVal: "enabled", intType: Ipv6ModeEnabled, stringType: ipv6ModeEnabled},
		{stringVal: "forced", intType: Ipv6ModeForced, stringType: ipv6ModeForced},
		{stringVal: "link_local", intType: Ipv6ModeLinkLocal, stringType: ipv6ModeLinkLocal},

		{stringVal: "vlan", intType: VnTypeVlan, stringType: vnTypeVlan},
		{stringVal: "vxlan", intType: VnTypeVxlan, stringType: vnTypeVxlan},
		{stringVal: "external", intType: VnTypeExternal, stringType: vnTypeExternal},

		{stringVal: "", intType: SystemRoleNone, stringType: systemRoleNone},
		{stringVal: "access", intType: SystemRoleAccess, stringType: systemRoleAccess},
		{stringVal: "leaf", intType: SystemRoleLeaf, stringType: systemRoleLeaf},
	}

	for i, td := range testData {
		ii := td.intType.int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}
