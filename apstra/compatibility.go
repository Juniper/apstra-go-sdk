package apstra

import "strings"

const (
	apstraSupportedApiVersions = "4.1.0, 4.1.1, 4.1.2"
	apstraSupportedVersionSep  = ","

	PodBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2"

	RackBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2"
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

// ApstraApiSupportedVersions returns the Apstra versions supported by this
// SDK version.
func ApstraApiSupportedVersions() StringSliceWithIncludes {
	return apstraSupportedApi()
}
