// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/errors"
)

func (c Client) CreateTemplatePodBased2(ctx context.Context, v design.TemplatePodBased) (string, error) {
	return c.CreateTemplate2(ctx, &v)
}

func (c Client) GetTemplatePodBased2(ctx context.Context, id string) (design.TemplatePodBased, error) {
	response, err := c.GetTemplate2(ctx, id)
	if err != nil {
		return design.TemplatePodBased{}, err
	}

	if response == nil {
		return design.TemplatePodBased{}, errors.Internal("template is unexpectedly nil")
	}

	if response.TemplateType() != enum.TemplateTypePodBased {
		return design.TemplatePodBased{}, errors.WrongType(fmt.Sprintf("template with id %q has wrong type: expected %q got %q", id, enum.TemplateTypePodBased, response.TemplateType()))
	}

	if result, ok := response.(*design.TemplatePodBased); ok {
		return *result, nil
	}

	return design.TemplatePodBased{}, errors.Internal(fmt.Sprintf("response has unexpected underlying type %T", response))
}

func (c Client) UpdateTemplatePodBased2(ctx context.Context, v design.TemplatePodBased) error {
	return c.UpdateTemplate2(ctx, &v)
}

func (c Client) ListTemplatesPodBased2(ctx context.Context) ([]string, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []string
	for i, t := range templates {
		if t.TemplateType() == enum.TemplateTypePodBased {
			if t.ID() == nil {
				return nil, errors.APIResponseInvalid(fmt.Sprintf("template at index %d has nil id", i))
			}
			result = append(result, *t.ID())
		}
	}

	return result, nil
}

func (c Client) GetTemplatesPodBased2(ctx context.Context) ([]design.TemplatePodBased, error) {
	templates, err := c.GetTemplates2(ctx)
	if err != nil {
		return nil, err
	}

	var result []design.TemplatePodBased
	for i, t := range templates {
		if t.TemplateType() != enum.TemplateTypePodBased {
			continue
		}

		tt, ok := t.(*design.TemplatePodBased)
		if !ok {
			return nil, errors.Internal(fmt.Sprintf("template at index %d claims to be a %q but has type %T", i, enum.TemplateTypePodBased, t))
		}
		if tt == nil {
			return nil, errors.Internal(fmt.Sprintf("template at index %d is unexpectedly nil", i))
		}

		result = append(result, *tt)
	}

	return result, nil
}

func (c Client) GetTemplatePodBasedByLabel2(ctx context.Context, label string) (design.TemplatePodBased, error) {
	t, err := c.GetTemplateByLabel2(ctx, label)
	if err != nil {
		return design.TemplatePodBased{}, err
	}

	if t == nil {
		return design.TemplatePodBased{}, errors.Internal(fmt.Sprintf("template with label %q is unexpectedly nil", label))
	}

	if t.TemplateType() != enum.TemplateTypePodBased {
		return design.TemplatePodBased{}, errors.WrongType(fmt.Sprintf("template with label %q has type %q", label, t.TemplateType()))
	}

	result, ok := t.(*design.TemplatePodBased)
	if !ok {
		return design.TemplatePodBased{}, errors.Internal(fmt.Sprintf("template with label %q has unexpected type %T", label, t))
	}

	return *result, nil
}
