// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"encoding/json"
	"github.com/Juniper/apstra-go-sdk/enum"
)

var _ json.Marshaler = (*VirtualNetwork)(nil)

type VirtualNetwork struct {
	OverlayControlProtocol enum.OverlayControlProtocol `json:"overlay_control_protocol"`
}

func (v VirtualNetwork) MarshalJSON() ([]byte, error) {
	// we want to send `"overlay_control_protocol": null` when OverlayControlProtocolNone

	var raw struct {
		OverlayControlProtocol *string `json:"overlay_control_protocol"`
	}

	if v.OverlayControlProtocol != enum.OverlayControlProtocolNone {
		raw.OverlayControlProtocol = &v.OverlayControlProtocol.Value
	}

	return json.Marshal(raw)
}
