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
	_ Template          = (*CommonTemplate)(nil)
	_ internal.IDer     = (*CommonTemplate)(nil)
	_ json.Unmarshaler  = (*CommonTemplate)(nil)
	_ timeutils.Stamper = (*CommonTemplate)(nil)
)

// CommonTemplate represents any of template types returned by /api/design/templates
// or /api/design/templates/{template_id}. The expected workflow is:
// - unmarshal the API response onto this object
// - determine the type by invoking TemplateType()
type CommonTemplate struct {
	id             string `json`
	kind           enum.TemplateType
	raw            json.RawMessage
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (t CommonTemplate) TemplateType() enum.TemplateType {
	return t.kind
}

func (t CommonTemplate) ID() *string {
	if t.id == "" {
		return nil
	}
	return &t.id
}

func (t *CommonTemplate) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID             string             `json:"id"`
		Kind           *enum.TemplateType `json:"type"`
		CreatedAt      *time.Time         `json:"created_at"`
		LastModifiedAt *time.Time         `json:"last_modified_at"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("unmarshalling template: %w", err)
	}

	if raw.ID == "" {
		return sdk.ErrAPIResponseInvalid("unmarshaling template: id is empty")
	}

	if raw.Kind == nil {
		return sdk.ErrAPIResponseInvalid("unmarshaling template: kind is nil")
	}

	t.id = raw.ID
	t.kind = *raw.Kind
	t.raw = bytes
	t.createdAt = raw.CreatedAt
	t.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func (t CommonTemplate) CreatedAt() *time.Time {
	return t.createdAt
}

func (t CommonTemplate) LastModifiedAt() *time.Time {
	return t.lastModifiedAt
}

func (t CommonTemplate) Template() (Template, error) {
	switch t.kind.String() {
	case enum.TemplateTypeL3Collapsed.String():
		return t.L3Collapsed()
	case enum.TemplateTypeRackBased.String():
		return t.RackBased()
	case enum.TemplateTypePodBased.String():
		return t.PodBased()
	case enum.TemplateTypeRailCollapsed.String():
		return t.RailCollapsed()
	default:
		return nil, fmt.Errorf("unhandled template type: %s", t.kind)
	}
}

func (t CommonTemplate) L3Collapsed() (TemplateL3Collapsed, error) {
	var result TemplateL3Collapsed
	if t.TemplateType() != enum.TemplateTypeL3Collapsed {
		return result, fmt.Errorf("cannot extract l3 collapsed template from common template of type %q", t.TemplateType())
	}

	if err := json.Unmarshal(t.raw, &result); err != nil {
		return result, fmt.Errorf("unmarshalling l3 collapsed template: %w", err)
	}

	return result, nil
}

func (t CommonTemplate) RackBased() (TemplateRackBased, error) {
	var result TemplateRackBased
	if t.TemplateType() != enum.TemplateTypeRackBased {
		return result, fmt.Errorf("cannot extract rack based template from common template of type %q", t.TemplateType())
	}

	if err := json.Unmarshal(t.raw, &result); err != nil {
		return result, fmt.Errorf("unmarshalling rack based template: %w", err)
	}

	return result, nil
}

func (t CommonTemplate) PodBased() (TemplatePodBased, error) {
	var result TemplatePodBased
	if t.TemplateType() != enum.TemplateTypePodBased {
		return result, fmt.Errorf("cannot extract pod based template from common template of type %q", t.TemplateType())
	}

	if err := json.Unmarshal(t.raw, &result); err != nil {
		return result, fmt.Errorf("unmarshalling pod based template: %w", err)
	}

	return result, nil
}

func (t CommonTemplate) RailCollapsed() (TemplateRailCollapsed, error) {
	var result TemplateRailCollapsed
	if t.TemplateType() != enum.TemplateTypeRailCollapsed {
		return result, fmt.Errorf("cannot extract rail collapsed template from common template of type %q", t.TemplateType())
	}

	if err := json.Unmarshal(t.raw, &result); err != nil {
		return result, fmt.Errorf("unmarshalling rail collapsed template: %w", err)
	}

	return result, nil
}

type RackTypeCount struct {
	RackTypeId string `json:"rack_type_id"`
	Count      int    `json:"count"`
}
