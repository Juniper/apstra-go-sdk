package apstra

import (
	"strings"

	"github.com/hashicorp/go-version"
)

const (
	apstra410 = "4.1.0"
	apstra411 = "4.1.1"
	apstra412 = "4.1.2"
	apstra420 = "4.2.0"
	apstra421 = "4.2.1"

	apstraSupportedApiVersions = "4.1.0, 4.1.1, 4.1.2, 4.2.0"
	apstraSupportedVersionSep  = ","

	podBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2, 4.2.0, 4.2.1"

	rackBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2, 4.2.0, 4.2.1"

	fabricL3MtuForbiddenError = "fabric_l3_mtu permitted only with Apstra 4.2.0 and later"

	integerPoolForbiddenVersions = "4.1.0, 4.1.1"

	policyRuleTcpStateQualifierForbidenVersions = "4.1.0, 4.1.1"
	policyRuleTcpStateQualifierForbidenError    = "tcp_state_qualifier permitted only with Apstra 4.1.2 and later"

	securityZoneJunosEvpnIrbModeRequiredVersions = "4.2.0"
	securityZoneJunosEvpnIrbModeRequiredError    = "junos_evpn_irb_mode is required by Apstra 4.2 and later"

	vnL3MtuForbiddenVersions = "4.1.0, 4.1.1, 4.1.2"
	vnL3MtuForbiddenError    = "virtual network operations support L3 MTU option only with Apstra 4.2 and later"
)

var (
	eqApstra410 = version.MustConstraints(version.NewConstraint("=" + apstra410))
	eqApstra411 = version.MustConstraints(version.NewConstraint("=" + apstra411))
	eqApstra412 = version.MustConstraints(version.NewConstraint("=" + apstra412))
	eqApstra420 = version.MustConstraints(version.NewConstraint("=" + apstra420))
	eqApstra421 = version.MustConstraints(version.NewConstraint("=" + apstra421))

	geApstra410 = version.MustConstraints(version.NewConstraint(">=" + apstra410))
	geApstra411 = version.MustConstraints(version.NewConstraint(">=" + apstra411))
	geApstra412 = version.MustConstraints(version.NewConstraint(">=" + apstra412))
	geApstra420 = version.MustConstraints(version.NewConstraint(">=" + apstra420))
	geApstra421 = version.MustConstraints(version.NewConstraint(">=" + apstra421))

	gtApstra410 = version.MustConstraints(version.NewConstraint(">" + apstra410))
	gtApstra411 = version.MustConstraints(version.NewConstraint(">" + apstra411))
	gtApstra412 = version.MustConstraints(version.NewConstraint(">" + apstra412))
	gtApstra420 = version.MustConstraints(version.NewConstraint(">" + apstra420))
	gtApstra421 = version.MustConstraints(version.NewConstraint(">" + apstra421))

	leApstra410 = version.MustConstraints(version.NewConstraint("<=" + apstra410))
	leApstra411 = version.MustConstraints(version.NewConstraint("<=" + apstra411))
	leApstra412 = version.MustConstraints(version.NewConstraint("<=" + apstra412))
	leApstra420 = version.MustConstraints(version.NewConstraint("<=" + apstra420))
	leApstra421 = version.MustConstraints(version.NewConstraint("<=" + apstra421))

	ltApstra410 = version.MustConstraints(version.NewConstraint("<" + apstra410))
	ltApstra411 = version.MustConstraints(version.NewConstraint("<" + apstra411))
	ltApstra412 = version.MustConstraints(version.NewConstraint("<" + apstra412))
	ltApstra420 = version.MustConstraints(version.NewConstraint("<" + apstra420))
	ltApstra421 = version.MustConstraints(version.NewConstraint("<" + apstra421))

	fabricSettingsApiOk  = geApstra421
	fabricL3MtuForbidden = leApstra412
)

// SupportedApiVersions returns []string with each element representing an Apstra version number like "4.2.0"
func SupportedApiVersions() []string {
	return []string{
		apstra410,
		apstra411,
		apstra412,
		apstra420,
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
