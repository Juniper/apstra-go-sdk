// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ design.Template   = (*CommonTemplate)(nil)
	_ internal.IDer     = (*CommonTemplate)(nil)
	_ json.Unmarshaler  = (*CommonTemplate)(nil)
	_ timeutils.Stamper = (*CommonTemplate)(nil)
)

// CommonTemplate represents any of template types returned by /api/design/templates
// or /api/design/templates/{template_id}. The expected workflow is:
// - unmarshal the API response onto this object
// - determine the type by invoking TemplateType()
type CommonTemplate struct {
	id             string
	templateType   enum.TemplateType
	createdAt      *time.Time
	lastModifiedAt *time.Time
	l3Collapsed    *design.TemplateL3Collapsed
	podBased       *design.TemplatePodBased
	rackBased      *design.TemplateRackBased
	railCollapsed  *design.TemplateRailCollapsed
}

func (t CommonTemplate) TemplateType() enum.TemplateType {
	return t.templateType
}

func (t CommonTemplate) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (t *CommonTemplate) SetID(id string) error {
	if t.id != "" {
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", t.id))
	}

	t.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (t *CommonTemplate) MustSetID(id string) {
	err := t.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (t *CommonTemplate) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string             `json:"id"`
		Type           *enum.TemplateType `json:"type"`
		CreatedAt      *time.Time         `json:"created_at"`
		LastModifiedAt *time.Time         `json:"last_modified_at"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("unmarshalling template: %w", err)
	}

	if raw.ID == "" {
		return sdk.ErrAPIResponseInvalid("unmarshaling template: id is empty")
	}

	if raw.Type == nil {
		return sdk.ErrAPIResponseInvalid("unmarshaling template: templateType is nil")
	}

	t.id = raw.ID
	t.templateType = *raw.Type
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	switch t.templateType.String() {
	case enum.TemplateTypeL3Collapsed.String():
		t.l3Collapsed = new(design.TemplateL3Collapsed)
		return json.Unmarshal(bytes, t.l3Collapsed)
	case enum.TemplateTypePodBased.String():
		t.podBased = new(design.TemplatePodBased)
		return json.Unmarshal(bytes, t.podBased)
	case enum.TemplateTypeRackBased.String():
		t.rackBased = new(design.TemplateRackBased)
		return json.Unmarshal(bytes, t.rackBased)
	case enum.TemplateTypeRailCollapsed.String():
		t.railCollapsed = new(design.TemplateRailCollapsed)
		return json.Unmarshal(bytes, t.railCollapsed)
	}

	return sdk.ErrAPIResponseInvalid(fmt.Sprintf("unhandled template type: %q", t.templateType))
}

func (t CommonTemplate) CreatedAt() *time.Time {
	return t.createdAt
}

func (t CommonTemplate) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}
