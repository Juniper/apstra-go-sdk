package apstra

import "strings"

const (
	apstraSupportedApiVersions = "4.1.0, 4.1.1"
	apstraSupportedVersionSep  = ","

	PodBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1"

	RackBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1"
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
	return parseVersionList(RackBasedTemplateFabricAddressingPolicyForbiddenVersions)
}

func podBasedTemplateFabricAddressingPolicyForbidden() StringSliceWithIncludes {
	return parseVersionList(PodBasedTemplateFabricAddressingPolicyForbiddenVersions)
}
