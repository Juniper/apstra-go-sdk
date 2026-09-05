// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package urls

import "regexp"

const (
	datacenterEvpnInterconnectGroupByIDRegexStr = blueprintByIDRegexStr + pathDelim + datacenterEvpnInterconnectGroupsPathComponent + pathDelim + "[^/]+$"
	datacenterEvpnInterconnectGroupsRegexStr    = blueprintByIDRegexStr + pathDelim + datacenterEvpnInterconnectGroupsPathComponent + "$"
)

var (
	DatacenterEvpnInterconnectGroupByIDRegex = regexp.MustCompile(datacenterEvpnInterconnectGroupByIDRegexStr)
	DatacenterEvpnInterconnectGroupsRegex    = regexp.MustCompile(datacenterEvpnInterconnectGroupsRegexStr)
)
