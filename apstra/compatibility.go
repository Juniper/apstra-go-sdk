package apstra

import "strings"

const (
	apstraSupportedApiVersions = "4.1.0, 4.1.1, 4.1.2"
	apstraSupportedVersionSep  = ","

	podBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2"

	rackBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1, 4.1.2"

	integerPoolForbiddenVersions = "4.1.0, 4.1.1"
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

func integerPoolForbidden() StringSliceWithIncludes {
	return parseVersionList(integerPoolForbiddenVersions)
}

// ApstraApiSupportedVersions returns the Apstra versions supported by this
// SDK version.
func ApstraApiSupportedVersions() StringSliceWithIncludes {
	return apstraSupportedApi()
}
