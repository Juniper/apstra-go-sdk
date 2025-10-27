// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"encoding/json"
	"fmt"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
)

const heuristic = "heuristic"

var (
	_ json.Marshaler   = (*AntiAffinity)(nil)
	_ json.Unmarshaler = (*AntiAffinity)(nil)
)

type AntiAffinity struct {
	MaxLinksPerPort          int
	MaxLinksPerSlot          int
	MaxPerSystemLinksPerPort int
	MaxPerSystemLinksPerSlot int
	Mode                     enum.AntiAffinityMode
}

func (a AntiAffinity) MarshalJSON() ([]byte, error) {
	raw := rawAntiAffinity{
		Algorithm:                heuristic,
		MaxLinksPerPort:          a.MaxLinksPerPort,
		MaxLinksPerSlot:          a.MaxLinksPerSlot,
		MaxPerSystemLinksPerPort: a.MaxPerSystemLinksPerPort,
		MaxPerSystemLinksPerSlot: a.MaxPerSystemLinksPerSlot,
		Mode:                     a.Mode,
	}

	return json.Marshal(raw)
}

func (a *AntiAffinity) UnmarshalJSON(bytes []byte) error {
	var raw rawAntiAffinity

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	if raw.Algorithm != heuristic {
		return sdk.ErrAPIResponseInvalid(fmt.Sprintf("anti affinity policy has invalid algorithm: %q", raw.Algorithm))
	}

	a.MaxLinksPerPort = raw.MaxLinksPerPort
	a.MaxLinksPerSlot = raw.MaxLinksPerSlot
	a.MaxPerSystemLinksPerPort = raw.MaxPerSystemLinksPerPort
	a.MaxPerSystemLinksPerSlot = raw.MaxPerSystemLinksPerSlot
	a.Mode = raw.Mode

	return nil
}

type rawAntiAffinity struct {
	Algorithm                string                `json:"algorithm"` // must be 'heuristic'
	MaxLinksPerPort          int                   `json:"max_links_per_port"`
	MaxLinksPerSlot          int                   `json:"max_links_per_slot"`
	MaxPerSystemLinksPerPort int                   `json:"max_per_system_links_per_port"`
	MaxPerSystemLinksPerSlot int                   `json:"max_per_system_links_per_slot"`
	Mode                     enum.AntiAffinityMode `json:"mode"`
}
