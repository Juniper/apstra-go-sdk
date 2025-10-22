// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
)

func (c Client) CreateTemplateRailCollapsed2(ctx context.Context, v design.TemplateRailCollapsed) (string, error) {
	return c.CreateTemplate2(ctx, &v)
}

func (c Client) GetTemplateRailCollapsed2(ctx context.Context, id string) (design.TemplateRailCollapsed, error) {
	response, err := c.GetTemplate2(ctx, id)
	if err != nil {
		return design.TemplateRailCollapsed{}, err
	}

	if response == nil {
		return design.TemplateRailCollapsed{}, sdk.ErrInternal("template is unexpectedly nil")
	}

	if response.TemplateType() != enum.TemplateTypeRailCollapsed {
		return design.TemplateRailCollapsed{}, sdk.ErrWrongType(fmt.Sprintf("template with id %q has wrong type: expected %q got %q", id, enum.TemplateTypeRailCollapsed, response.TemplateType()))
	}

	if result, ok := response.(*design.TemplateRailCollapsed); ok {
		return *result, nil
	}

	return design.TemplateRailCollapsed{}, sdk.ErrInternal(fmt.Sprintf("response has unexpected underlying type %T", response))
}

func (c Client) UpdateTemplateRailCollapsed2(ctx context.Context, v design.TemplateRailCollapsed) error {
	return c.UpdateTemplate2(ctx, &v)
}

func (c Client) ListTemplatesRailCollapsed2(ctx context.Context) ([]string, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []string
	for i, t := range templates {
		if t.TemplateType() == enum.TemplateTypeRailCollapsed {
			if t.ID() == nil {
				return nil, sdk.ErrAPIResponseInvalid(fmt.Sprintf("template at index %d has nil id", i))
			}
			result = append(result, *t.ID())
		}
	}

	return result, nil
}

func (c Client) GetTemplatesRailCollapsed2(ctx context.Context) ([]design.TemplateRailCollapsed, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []design.TemplateRailCollapsed
	for i, t := range templates {
		if t.TemplateType() != enum.TemplateTypeRailCollapsed {
			continue
		}

		tt, ok := t.(*design.TemplateRailCollapsed)
		if !ok {
			return nil, sdk.ErrInternal(fmt.Sprintf("template at index %d claims to be a %q but has type %T", i, enum.TemplateTypeRailCollapsed, t))
		}
		if tt == nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("template at index %d is unexpectedly nil", i))
		}

		result = append(result, *tt)
	}

	return result, nil
}

func (c Client) GetTemplateRailCollapsedByLabel2(ctx context.Context, label string) (design.TemplateRailCollapsed, error) {
	t, err := c.GetTemplateByLabel2(ctx, label)
	if err != nil {
		return design.TemplateRailCollapsed{}, err
	}

	if t == nil {
		return design.TemplateRailCollapsed{}, sdk.ErrInternal(fmt.Sprintf("template with label %q is unexpectedly nil", label))
	}

	if t.TemplateType() != enum.TemplateTypeRailCollapsed {
		return design.TemplateRailCollapsed{}, sdk.ErrWrongType(fmt.Sprintf("template with label %q has type %q", label, t.TemplateType()))
	}

	result, ok := t.(*design.TemplateRailCollapsed)
	if !ok {
		return design.TemplateRailCollapsed{}, sdk.ErrInternal(fmt.Sprintf("template with label %q has unexpected type %T", label, t))
	}

	return *result, nil
}
