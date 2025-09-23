// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

const (
	urlPrefix            = "/api/design/"
	TagsUrl              = urlPrefix + "tags"
	TagUrlByID           = TagsUrl + "/%s"
	LogicalDevicesUrl    = urlPrefix + "logical-devices"
	LogicalDeviceUrlByID = LogicalDevicesUrl + "/%s"
	RackTypesUrl         = urlPrefix + "rack-types"
	RackTypeUrlByID      = RackTypesUrl + "/%s"
)
