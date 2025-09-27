// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

const (
	urlPrefix = "/api/design/"

	InterfaceMapDigestsURL    = urlPrefix + "interface-map-digests"
	InterfaceMapDigestURLByID = InterfaceMapDigestsURL + "/%s"
	InterfaceMapsURL          = urlPrefix + "interface-maps"
	InterfaceMapURLByID       = InterfaceMapsURL + "/%s"
	LogicalDevicesURL         = urlPrefix + "logical-devices"
	LogicalDeviceURLByID      = LogicalDevicesURL + "/%s"
	RackTypesURL              = urlPrefix + "rack-types"
	RackTypeURLByID           = RackTypesURL + "/%s"
	TagsURL                   = urlPrefix + "tags"
	TagURLByID                = TagsURL + "/%s"
)
