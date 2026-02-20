// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "github.com/Juniper/apstra-go-sdk/enum"

// AddressingPolicy is used during Blueprint creation. It was introduced with Apstra 6.1.
// These are Security Zone (Routing Zone) parameters which are used as initial values for
// the Blueprint's default Security Zone.
type AddressingPolicy struct {
	AddressingSupport *enum.AddressingScheme `json:"addressing_support,omitempty"`
	DisableIPv4       *bool                  `json:"disable_ipv4,omitempty"`
	VTEPAddressing    *enum.AddressingScheme `json:"vtep_addressing,omitempty"`
}
