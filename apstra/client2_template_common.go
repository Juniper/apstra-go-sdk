// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	commontemplate "github.com/Juniper/apstra-go-sdk/internal/template"
)

func (c Client) CreateTemplate2(ctx context.Context, v design.Template) (string, error) {
	if v.ID() != nil {
		return "", fmt.Errorf("id must be nil in %s", str.FuncName())
	}

	var response struct {
		ID string `json:"id"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      design.TemplatesURL,
		apiInput:    v,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.ID, nil
}

func (c Client) GetTemplate2(ctx context.Context, id string) (design.Template, error) {
	var response commontemplate.Common
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(design.TemplateURLByID, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	switch response.TemplateType().String() {
	case enum.TemplateTypeL3Collapsed.String():
		if response.L3Collapsed == nil {
			return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
		}
		return response.L3Collapsed, nil
	case enum.TemplateTypePodBased.String():
		if response.PodBased == nil {
			return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
		}
		return response.PodBased, nil
	case enum.TemplateTypeRackBased.String():
		if response.RackBased == nil {
			return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
		}
		return response.RackBased, nil
	case enum.TemplateTypeRailCollapsed.String():
		if response.RailCollapsed == nil {
			return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
		}
		return response.RailCollapsed, nil
	}

	return nil, sdk.ErrInternal(fmt.Sprintf("internal error: unhandled template type %q", response.TemplateType()))
}

func (c Client) UpdateTemplate2(ctx context.Context, v design.Template) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(design.TemplateURLByID, *v.ID()),
		apiInput: v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) DeleteTemplate2(ctx context.Context, id string) error {
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(design.TemplateURLByID, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (c Client) ListTemplates2(ctx context.Context) ([]string, error) {
	var response struct {
		Items []string `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      design.TemplatesURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (c Client) GetTemplates2(ctx context.Context) ([]design.Template, error) {
	var response struct {
		Items []commontemplate.Common `json:"items"`
	}

	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      design.TemplatesURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]design.Template, 0, len(response.Items))
	for _, item := range response.Items {
		switch item.TemplateType().String() {
		case enum.TemplateTypeL3Collapsed.String():
			result = append(result, item.L3Collapsed)
		case enum.TemplateTypePodBased.String():
			result = append(result, item.PodBased)
		case enum.TemplateTypeRackBased.String():
			result = append(result, item.RackBased)
		case enum.TemplateTypeRailCollapsed.String():
			// result = append(result, item.RailCollapsed) // todo: restore this when rail collapsed templates are supported
		default:
			var id string
			if item.ID() != nil {
				id = *item.ID()
			}
			return nil, sdk.ErrInternal(fmt.Sprintf("internal error: template %q has unhandled type %q", id, item.TemplateType()))
		}
	}

	return result, nil
}

func (c Client) GetTemplateByLabel2(ctx context.Context, label string) (design.Template, error) {
	var response struct {
		Items []commontemplate.Common `json:"items"`
	}
	err := c.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      design.TemplatesURL,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var result []commontemplate.Common

	for _, item := range response.Items {
		if item.Label == label {
			result = append(result, item)
		}
	}

	switch len(result) {
	case 0:
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("template with label %s not found", label),
		}
	case 1:
		switch result[0].TemplateType().String() {
		case enum.TemplateTypeL3Collapsed.String():
			if result[0].L3Collapsed == nil {
				return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
			}
			return result[0].L3Collapsed, nil
		case enum.TemplateTypePodBased.String():
			if result[0].PodBased == nil {
				return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
			}
			return result[0].PodBased, nil
		case enum.TemplateTypeRackBased.String():
			if result[0].RackBased == nil {
				return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
			}
			return result[0].RackBased, nil
		case enum.TemplateTypeRailCollapsed.String():
			if result[0].RailCollapsed == nil {
				return nil, sdk.ErrInternal("internal error: embedded template is unexpectedly nil")
			}
			return result[0].RailCollapsed, nil
		default:
			var id string
			if result[0].ID() != nil {
				id = *result[0].ID()
			}
			return nil, sdk.ErrInternal(fmt.Sprintf("internal error: template %q has unhandled type %q", id, result[0].TemplateType()))
		}
	default: // len(result) > 1
		return nil, ClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found multiple candidate templates with label %s", label),
		}
	}
}
