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
	_ design.Template   = (*Common)(nil)
	_ internal.IDer     = (*Common)(nil)
	_ json.Unmarshaler  = (*Common)(nil)
	_ timeutils.Stamper = (*Common)(nil)
)

// Common represents any of template types returned by /api/design/templates
// or /api/design/templates/{template_id}. The expected workflow is:
// - unmarshal the API response onto this object
// - determine the type by invoking TemplateType()
type Common struct {
	id             string
	templateType   enum.TemplateType
	createdAt      *time.Time
	lastModifiedAt *time.Time
	Label          string
	L3Collapsed    *design.TemplateL3Collapsed
	PodBased       *design.TemplatePodBased
	RackBased      *design.TemplateRackBased
	RailCollapsed  *design.TemplateRailCollapsed
}

func (t Common) TemplateType() enum.TemplateType {
	return t.templateType
}

func (t Common) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

func (t *Common) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string             `json:"id"`
		DisplayName    string             `json:"display_name"`
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
	t.Label = raw.DisplayName
	t.templateType = *raw.Type
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	switch t.templateType.String() {
	case enum.TemplateTypeL3Collapsed.String():
		t.L3Collapsed = new(design.TemplateL3Collapsed)
		return json.Unmarshal(bytes, t.L3Collapsed)
	case enum.TemplateTypePodBased.String():
		t.PodBased = new(design.TemplatePodBased)
		return json.Unmarshal(bytes, t.PodBased)
	case enum.TemplateTypeRackBased.String():
		t.RackBased = new(design.TemplateRackBased)
		return json.Unmarshal(bytes, t.RackBased)
	case enum.TemplateTypeRailCollapsed.String():
		t.RailCollapsed = new(design.TemplateRailCollapsed)
		return json.Unmarshal(bytes, t.RailCollapsed)
	}

	return sdk.ErrAPIResponseInvalid(fmt.Sprintf("unhandled template type: %q", t.templateType))
}

func (t Common) CreatedAt() *time.Time {
	return t.createdAt
}

func (t Common) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}
