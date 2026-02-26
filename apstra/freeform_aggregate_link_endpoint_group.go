// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"encoding/json"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

var (
	_ internal.IDer    = (*FreeformAggregateLinkEndpointGroup)(nil)
	_ json.Unmarshaler = (*FreeformAggregateLinkEndpointGroup)(nil)
	_ json.Marshaler   = (*FreeformAggregateLinkEndpointGroup)(nil)
)

// FreeformAggregateLinkEndpointGroup represents one end of a Freefor LAG link. Each
// FreeformAggregateLinkEndpointGroup contains a collection of FreeformAggregateLinkEndpoint.
// FreeformAggregateLinkEndpoint is the logical LAG interface (ae1, bond0) on a device. Because
// a LAG may be terminated by a multi-chassis scheme (MLAG, ESI LAG), one endpoint is not enough
// to describe one end of a LAG. Each LAG has exactly two FreeformAggregateLinkEndpointGroup,
// and each FreeformAggregateLinkEndpointGroup has at least one FreeformAggregateLinkEndpoint.
// Because of L3 MLAG schemes (e.g. HSRP routing over Cisco VPC) it is possible for each
// FreeformAggregateLinkEndpoint to have its own IPv4 and IPv6 address.`
type FreeformAggregateLinkEndpointGroup struct {
	Label     *string
	Tags      []string
	Endpoints []FreeformAggregateLinkEndpoint

	endpointGroupNumber int
	id                  string
}

func (o FreeformAggregateLinkEndpointGroup) MarshalJSON() ([]byte, error) {
	raw := struct {
		Label json.RawMessage `json:"label,omitempty"`
		Tags  []string        `json:"tags"`
	}{
		Label: pointer.StringMarshalJSONWithEmptyAsNull(o.Label),
		Tags:  o.Tags,
	}

	return json.Marshal(&raw)
}

func (o FreeformAggregateLinkEndpointGroup) ID() *string {
	if o.id == "" {
		return nil
	}
	return &o.id
}

func (o *FreeformAggregateLinkEndpointGroup) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID    string   `json:"id"`
		Label *string  `json:"label"`
		Tags  []string `json:"tags"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	o.id = raw.ID
	o.Label = raw.Label
	o.Tags = raw.Tags

	return nil
}
