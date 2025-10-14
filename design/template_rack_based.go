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
)

var (
	_ Template                               = (*TemplateRackBased)(nil)
	_ internal.IDSetter                      = (*TemplateRackBased)(nil)
	_ internal.Replicator[TemplateRackBased] = (*TemplateRackBased)(nil)
	_ json.Marshaler                         = (*TemplateRackBased)(nil)
	_ json.Unmarshaler                       = (*TemplateRackBased)(nil)
)

type TemplateRackBased struct {
	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (t TemplateRackBased) TemplateType() enum.TemplateType {
	return enum.TemplateTypeRackBased
}

func (t TemplateRackBased) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (t TemplateRackBased) SetID(id string) error {
	if t.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (t TemplateRackBased) MustSetID(id string) {
	err := t.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (t TemplateRackBased) Replicate() TemplateRackBased {
	//TODO implement me
	panic("implement me")
}

func (t TemplateRackBased) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t TemplateRackBased) UnmarshalJSON(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t TemplateRackBased) CreatedAt() *time.Time {
	return t.createdAt
}

func (t TemplateRackBased) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}
