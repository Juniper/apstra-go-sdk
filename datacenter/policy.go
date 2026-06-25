// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding/json"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/errors"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

var (
	_ internal.IDer     = (*Policy)(nil)
	_ internal.IDSetter = (*Policy)(nil)
	_ json.Unmarshaler  = (*Policy)(nil)
)

type Policy struct {
	Enabled             bool         `json:"enabled"`
	Label               string       `json:"label"`
	Description         string       `json:"description"`
	SrcApplicationPoint *string      `json:"src_application_point"`
	DstApplicationPoint *string      `json:"dst_application_point"`
	Rules               []PolicyRule `json:"rules"`
	Tags                []string     `json:"tags"`

	id string
}

func (p Policy) ID() *string {
	if p.id == "" {
		return nil
	}
	return pointer.ToCopyOf(p.id)
}

func (p *Policy) SetID(id string) error {
	if p == nil {
		return fmt.Errorf("cannot set ID of nil %T", p)
	}

	if p.id != "" {
		return errors.IDAlreadySet(fmt.Sprintf("id already has value %q", p.id))
	}

	p.id = id
	return nil
}

func (p *Policy) UnmarshalJSON(bytes []byte) error {
	if p == nil {
		return fmt.Errorf("cannot unmarshal into nil %T", p)
	}

	type policy Policy // type alias prevents recursion

	var target struct {
		policy
		ID  string `json:"id"` // the `id` struct element cannot be unmarshaled so we temporarily stash that value here
		Src struct {
			ID *string `json:"id"`
		} `json:"src_application_point"`
		Dst struct {
			ID *string `json:"id"`
		} `json:"dst_application_point"`
	}

	// unmarshal everything which can be handled by the `json` package.
	if err := json.Unmarshal(bytes, &target); err != nil {
		return err
	}

	*p = Policy(target.policy)
	p.id = target.ID
	p.SrcApplicationPoint = target.Src.ID
	p.DstApplicationPoint = target.Dst.ID

	return nil
}
