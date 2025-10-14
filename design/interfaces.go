// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

type logicalDeviceIDer interface {
	logicalDeviceID() *string
}

type Template interface {
	timeutils.Stamper
	internal.IDer
	TemplateType() enum.TemplateType
}
