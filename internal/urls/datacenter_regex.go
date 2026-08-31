package urls

import "regexp"

const (
	datacenterEvpnInterconnectGroupByIDRegexStr = blueprintByIDRegexStr + pathDelim + datacenterEvpnInterconnectGroupsPathComponent + pathDelim + "[^/]+$"
)

var DatacenterEvpnInterconnectGroupByIDRegex = regexp.MustCompile(datacenterEvpnInterconnectGroupByIDRegexStr)
