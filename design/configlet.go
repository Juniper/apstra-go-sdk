// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ internal.IDSetter = (*Configlet)(nil)
	_ json.Unmarshaler  = (*Configlet)(nil)
	_ timeutils.Stamper = (*Configlet)(nil)
)

type Configlet struct {
	Label      string               `json:"display_name"`
	Generators []ConfigletGenerator `json:"generators"`
	RefArchs   []enum.RefDesign     `json:"ref_archs"`

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (c Configlet) ID() *string {
	if c.id == "" {
		return nil
	}
	return &c.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (c *Configlet) SetID(id string) error {
	if c.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", c.id))
	}

	c.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (c *Configlet) MustSetID(id string) {
	err := c.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (c Configlet) CreatedAt() *time.Time {
	return c.createdAt
}

func (c Configlet) LastModifiedAt() *time.Time {
	return c.lastModifiedAt
}

func (c *Configlet) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string               `json:"id"`
		DisplayName    string               `json:"display_name"`
		Generators     []ConfigletGenerator `json:"generators"`
		RefArchs       []enum.RefDesign     `json:"ref_archs"`
		CreatedAt      *time.Time           `json:"created_at"`
		LastModifiedAt *time.Time           `json:"last_modified_at"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	c.id = raw.ID
	c.Label = raw.DisplayName
	c.Generators = raw.Generators
	c.RefArchs = raw.RefArchs
	c.createdAt = raw.CreatedAt
	c.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func NewConfiglet(id string) Configlet {
	return Configlet{id: id}
}
