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

func (c Client) CreateTemplateRackBased2(ctx context.Context, v design.TemplateRackBased) (string, error) {
	return c.CreateTemplate2(ctx, &v)
}

func (c Client) GetTemplateRackBased2(ctx context.Context, id string) (design.TemplateRackBased, error) {
	response, err := c.GetTemplate2(ctx, id)
	if err != nil {
		return design.TemplateRackBased{}, err
	}

	if response == nil {
		return design.TemplateRackBased{}, sdk.ErrInternal("template is unexpectedly nil")
	}

	if response.TemplateType() != enum.TemplateTypeRackBased {
		return design.TemplateRackBased{}, sdk.ErrWrongType(fmt.Sprintf("template with id %q has wrong type: expected %q got %q", id, enum.TemplateTypeRackBased, response.TemplateType()))
	}

	if result, ok := response.(*design.TemplateRackBased); ok {
		return *result, nil
	}

	return design.TemplateRackBased{}, sdk.ErrInternal(fmt.Sprintf("response has unexpected underlying type %T", response))
}

func (c Client) UpdateTemplateRackBased2(ctx context.Context, v design.TemplateRackBased) error {
	return c.UpdateTemplate2(ctx, &v)
}

func (c Client) ListTemplatesRackBased2(ctx context.Context) ([]string, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []string
	for i, t := range templates {
		if t.TemplateType() == enum.TemplateTypeRackBased {
			if t.ID() == nil {
				return nil, sdk.ErrAPIResponseInvalid(fmt.Sprintf("template at index %d has nil id", i))
			}
			result = append(result, *t.ID())
		}
	}

	return result, nil
}

func (c Client) GetTemplatesRackBased2(ctx context.Context) ([]design.TemplateRackBased, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []design.TemplateRackBased
	for i, t := range templates {
		if t.TemplateType() != enum.TemplateTypeRackBased {
			continue
		}

		tt, ok := t.(*design.TemplateRackBased)
		if !ok {
			return nil, sdk.ErrInternal(fmt.Sprintf("template at index %d claims to be a %q but has type %T", i, enum.TemplateTypeRackBased, t))
		}
		if tt == nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("template at index %d is unexpectedly nil", i))
		}

		result = append(result, *tt)
	}

	return result, nil
}

func (c Client) GetTemplateRackBasedByLabel2(ctx context.Context, label string) (design.TemplateRackBased, error) {
	t, err := c.GetTemplateByLabel2(ctx, label)
	if err != nil {
		return design.TemplateRackBased{}, err
	}

	if t == nil {
		return design.TemplateRackBased{}, sdk.ErrInternal(fmt.Sprintf("template with label %q is unexpectedly nil", label))
	}

	if t.TemplateType() != enum.TemplateTypeRackBased {
		return design.TemplateRackBased{}, sdk.ErrWrongType(fmt.Sprintf("template with label %q has type %q", label, t.TemplateType()))
	}

	result, ok := t.(*design.TemplateRackBased)
	if !ok {
		return design.TemplateRackBased{}, sdk.ErrInternal(fmt.Sprintf("template with label %q has unexpected type %T", label, t))
	}

	return *result, nil
}
