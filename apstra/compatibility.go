package apstra

import "strings"

const (
	apstraSupportedApiVersions = "4.1.0, 4.1.1, 4.1.2, 4.2.0"
	apstraSupportedVersionSep  = ","

	podBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2, 4.2.0, 4.2.1"

	rackBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2, 4.2.0, 4.2.1"

	fabricL3MtuForbiddenVersions = "4.1.0, 4.1.1, 4.1.2"
	fabricL3MtuForbiddenError    = "fabric_l3_mtu permitted only with Apstra 4.2.0 and later"

	integerPoolForbiddenVersions = "4.1.0, 4.1.1"

	policyRuleTcpStateQualifierForbidenVersions = "4.1.0, 4.1.1"
	policyRuleTcpStateQualifierForbidenError    = "tcp_state_qualifier permitted only with Apstra 4.1.2 and later"

	securityZoneJunosEvpnIrbModeRequiredVersions = "4.2.0"
	securityZoneJunosEvpnIrbModeRequiredError    = "junos_evpn_irb_mode is required by Apstra 4.2 and later"

	vnL3MtuForbiddenVersions = "4.1.0, 4.1.1, 4.1.2"
	vnL3MtuForbiddenError    = "virtual network operations support L3 MTU option only with Apstra 4.2 and later"
)

func parseVersionList(s string) StringSliceWithIncludes {
	result := strings.Split(s, apstraSupportedVersionSep)
	for i, s := range result {
		result[i] = strings.TrimSpace(s)
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

func apstraSupportedApi() StringSliceWithIncludes {
	return parseVersionList(apstraSupportedApiVersions)
}

func rackBasedTemplateFabricAddressingPolicyForbidden() StringSliceWithIncludes {
	return parseVersionList(rackBasedTemplateFabricAddressingPolicyForbiddenVersions)
}

func podBasedTemplateFabricAddressingPolicyForbidden() StringSliceWithIncludes {
	return parseVersionList(podBasedTemplateFabricAddressingPolicyForbiddenVersions)
}

func fabricL3MtuForbidden() StringSliceWithIncludes {
	return parseVersionList(fabricL3MtuForbiddenVersions)
}

func integerPoolForbidden() StringSliceWithIncludes {
	return parseVersionList(integerPoolForbiddenVersions)
}

func policyRuleTcpStateQualifierForbidden() StringSliceWithIncludes {
	return parseVersionList(policyRuleTcpStateQualifierForbidenVersions)
}

func securityZoneJunosEvpnIrbModeRequired() StringSliceWithIncludes {
	return parseVersionList(securityZoneJunosEvpnIrbModeRequiredVersions)
}

func vnL3MtuForbidden() StringSliceWithIncludes {
	return parseVersionList(vnL3MtuForbiddenVersions)
}

// ApstraApiSupportedVersions returns the Apstra versions supported by this
// SDK version.
func ApstraApiSupportedVersions() StringSliceWithIncludes {
	return apstraSupportedApi()
}
