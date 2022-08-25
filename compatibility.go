package goapstra

import "strings"

const (
	apstraSupportedApiVersions = "4.1.0, 4.1.1"
	apstraSupportedVersionSep  = ","

	PodBasedTemplateFabricAddressingPolicyRequiredVersions  = "4.1.0"
	podBasedTemplateFabricAddressingPolicyRequiredErr       = "PodBasedTemplateRequest validation error: FabricAddressingPolicy required by Apstra %s"
	PodBasedTemplateFabricAddressingPolicyForbiddenVersions = "4.1.1"

	RackBasedTemplateFabricAddressingPolicyRequiredVersions  = "4.1.0"
	rackBasedTemplateFabricAddressingPolicyRequiredErr       = "RackBasedTemplateRequest validation error: FabricAddressingPolicy required by Apstra %s"
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

func rackBasedTemplateFabricAddressingPolicyRequired() StringSliceWithIncludes {
	return parseVersionList(RackBasedTemplateFabricAddressingPolicyRequiredVersions)
}

func rackBasedTemplateFabricAddressingPolicyForbidden() StringSliceWithIncludes {
	return parseVersionList(RackBasedTemplateFabricAddressingPolicyForbiddenVersions)
}

func podBasedTemplateFabricAddressingPolicyRequired() StringSliceWithIncludes {
	return parseVersionList(PodBasedTemplateFabricAddressingPolicyRequiredVersions)
}

func podBasedTemplateFabricAddressingPolicyForbidden() StringSliceWithIncludes {
	return parseVersionList(PodBasedTemplateFabricAddressingPolicyForbiddenVersions)
}
