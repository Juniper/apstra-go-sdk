// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package antiaffinity

import (
	"encoding/json"
	"fmt"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
)

const heuristic = "heuristic"

var (
	_ json.Marshaler   = (*Policy)(nil)
	_ json.Unmarshaler = (*Policy)(nil)
)

type Policy struct {
	MaxLinksPerPort          int
	MaxLinksPerSlot          int
	MaxPerSystemLinksPerPort int
	MaxPerSystemLinksPerSlot int
	Mode                     enum.AntiAffinityMode
}

func (a Policy) MarshalJSON() ([]byte, error) {
	raw := rawPolicy{
		Algorithm:                heuristic,
		MaxLinksPerPort:          a.MaxLinksPerPort,
		MaxLinksPerSlot:          a.MaxLinksPerSlot,
		MaxPerSystemLinksPerPort: a.MaxPerSystemLinksPerPort,
		MaxPerSystemLinksPerSlot: a.MaxPerSystemLinksPerSlot,
		Mode:                     a.Mode,
	}

	return json.Marshal(raw)
}

func (a *Policy) UnmarshalJSON(bytes []byte) error {
	var raw rawPolicy

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

type rawPolicy struct {
	Algorithm                string                `json:"algorithm"` // must be 'heuristic'
	MaxLinksPerPort          int                   `json:"max_links_per_port"`
	MaxLinksPerSlot          int                   `json:"max_links_per_slot"`
	MaxPerSystemLinksPerPort int                   `json:"max_per_system_links_per_port"`
	MaxPerSystemLinksPerSlot int                   `json:"max_per_system_links_per_slot"`
	Mode                     enum.AntiAffinityMode `json:"mode"`
}
