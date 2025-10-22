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

func (c Client) CreateTemplateL3Collapsed2(ctx context.Context, v design.TemplateL3Collapsed) (string, error) {
	return c.CreateTemplate2(ctx, &v)
}

func (c Client) GetTemplateL3Collapsed2(ctx context.Context, id string) (design.TemplateL3Collapsed, error) {
	response, err := c.GetTemplate2(ctx, id)
	if err != nil {
		return design.TemplateL3Collapsed{}, err
	}

	if response == nil {
		return design.TemplateL3Collapsed{}, sdk.ErrInternal("template is unexpectedly nil")
	}

	if response.TemplateType() != enum.TemplateTypeL3Collapsed {
		return design.TemplateL3Collapsed{}, sdk.ErrWrongType(fmt.Sprintf("template with id %q has wrong type: expected %q got %q", id, enum.TemplateTypeL3Collapsed, response.TemplateType()))
	}

	if result, ok := response.(*design.TemplateL3Collapsed); ok {
		return *result, nil
	}

	return design.TemplateL3Collapsed{}, sdk.ErrInternal(fmt.Sprintf("response has unexpected underlying type %T", response))
}

func (c Client) UpdateTemplateL3Collapsed2(ctx context.Context, v design.TemplateL3Collapsed) error {
	return c.UpdateTemplate2(ctx, &v)
}

func (c Client) ListTemplatesL3Collapsed2(ctx context.Context) ([]string, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []string
	for i, t := range templates {
		if t.TemplateType() == enum.TemplateTypeL3Collapsed {
			if t.ID() == nil {
				return nil, sdk.ErrAPIResponseInvalid(fmt.Sprintf("template at index %d has nil id", i))
			}
			result = append(result, *t.ID())
		}
	}

	return result, nil
}

func (c Client) GetTemplatesL3Collapsed2(ctx context.Context) ([]design.TemplateL3Collapsed, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []design.TemplateL3Collapsed
	for i, t := range templates {
		if t.TemplateType() != enum.TemplateTypeL3Collapsed {
			continue
		}

		tt, ok := t.(*design.TemplateL3Collapsed)
		if !ok {
			return nil, sdk.ErrInternal(fmt.Sprintf("template at index %d claims to be a %q but has type %T", i, enum.TemplateTypeL3Collapsed, t))
		}
		if tt == nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("template at index %d is unexpectedly nil", i))
		}

		result = append(result, *tt)
	}

	return result, nil
}

func (c Client) GetTemplateL3CollapsedByLabel2(ctx context.Context, label string) (design.TemplateL3Collapsed, error) {
	t, err := c.GetTemplateByLabel2(ctx, label)
	if err != nil {
		return design.TemplateL3Collapsed{}, err
	}

	if t == nil {
		return design.TemplateL3Collapsed{}, sdk.ErrInternal(fmt.Sprintf("template with label %q is unexpectedly nil", label))
	}

	if t.TemplateType() != enum.TemplateTypeL3Collapsed {
		return design.TemplateL3Collapsed{}, sdk.ErrWrongType(fmt.Sprintf("template with label %q has type %q", label, t.TemplateType()))
	}

	result, ok := t.(*design.TemplateL3Collapsed)
	if !ok {
		return design.TemplateL3Collapsed{}, sdk.ErrInternal(fmt.Sprintf("template with label %q has unexpected type %T", label, t))
	}

	return *result, nil
}
