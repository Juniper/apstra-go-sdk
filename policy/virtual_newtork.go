// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package policy

import "github.com/Juniper/apstra-go-sdk/enum"

type VirtualNetwork struct {
	OverlayControlProtocol enum.OverlayControlProtocol `json:"overlay_control_protocol,omitempty"`
}
