// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlDesignTemplates       = apiUrlDesignPrefix + "templates"
	apiUrlDesignTemplatesPrefix = apiUrlDesignTemplates + apiUrlPathDelim
	apiUrlDesignTemplateById    = apiUrlDesignTemplatesPrefix + "%s"
)

type AntiAffinityPolicy struct {
	Algorithm                enum.AntiAffinityAlgorithm
	MaxLinksPerPort          int
	MaxLinksPerSlot          int
	MaxPerSystemLinksPerPort int
	MaxPerSystemLinksPerSlot int
	Mode                     *enum.AntiAffinityMode
}

func (o *AntiAffinityPolicy) raw() *rawAntiAffinityPolicy {
	//var mode *enum.AntiAffinityMode
	//if o.Mode.Value == "" {
	//	mode = &enum.AntiAffinityModeDisabled
	//} else {
	//	mode = o.Mode
	//}

	return &rawAntiAffinityPolicy{
		Algorithm:                o.Algorithm,
		MaxLinksPerPort:          o.MaxLinksPerPort,
		MaxLinksPerSlot:          o.MaxLinksPerSlot,
		MaxPerSystemLinksPerPort: o.MaxPerSystemLinksPerPort,
		MaxPerSystemLinksPerSlot: o.MaxPerSystemLinksPerSlot,
		Mode:                     o.Mode,
	}
}

type rawAntiAffinityPolicy struct {
	Algorithm                enum.AntiAffinityAlgorithm `json:"algorithm"` // heuristic
	MaxLinksPerPort          int                        `json:"max_links_per_port"`
	MaxLinksPerSlot          int                        `json:"max_links_per_slot"`
	MaxPerSystemLinksPerPort int                        `json:"max_per_system_links_per_port"`
	MaxPerSystemLinksPerSlot int                        `json:"max_per_system_links_per_slot"`
	Mode                     *enum.AntiAffinityMode     `json:"mode,omitempty"` // disabled, enabled_loose, enabled_strict
}

func (o *rawAntiAffinityPolicy) polish() (*AntiAffinityPolicy, error) {
	return &AntiAffinityPolicy{
		Algorithm:                o.Algorithm,
		MaxLinksPerPort:          o.MaxLinksPerPort,
		MaxLinksPerSlot:          o.MaxLinksPerSlot,
		MaxPerSystemLinksPerPort: o.MaxPerSystemLinksPerPort,
		MaxPerSystemLinksPerSlot: o.MaxPerSystemLinksPerSlot,
		Mode:                     o.Mode,
	}, nil
}

type VirtualNetworkPolicy struct {
	OverlayControlProtocol *enum.OverlayControlProtocol
}

func (o *VirtualNetworkPolicy) raw() *rawVirtualNetworkPolicy {
	return &rawVirtualNetworkPolicy{OverlayControlProtocol: o.OverlayControlProtocol}
}

type rawVirtualNetworkPolicy struct {
	OverlayControlProtocol *enum.OverlayControlProtocol `json:"overlay_control_protocol,omitempty"`
}

func (o *rawVirtualNetworkPolicy) polish() (*VirtualNetworkPolicy, error) {
	return &VirtualNetworkPolicy{OverlayControlProtocol: o.OverlayControlProtocol}, nil
}

type RackTypeCount struct {
	RackTypeId ObjectId `json:"rack_type_id"`
	Count      int      `json:"count"`
}

type Template interface {
	Type() enum.TemplateType
	ID() ObjectId
	OverlayControlProtocol() enum.OverlayControlProtocol
}

type template json.RawMessage

func (o *template) templateType() (enum.TemplateType, error) {
	templateProto := &struct {
		Type enum.TemplateType `json:"type"`
	}{}
	return templateProto.Type, json.Unmarshal(*o, templateProto)
}

func (o *template) displayName() (string, error) {
	templateProto := &struct {
		DisplayName string `json:"display_name"`
	}{}
	return templateProto.DisplayName, json.Unmarshal(*o, templateProto)
}

func (o *template) polish() (Template, error) {
	t, err := o.templateType()
	if err != nil {
		return nil, err
	}

	// quick unmarshal just to get the type
	switch t {
	case enum.TemplateTypeRackBased:
		var t rawTemplateRackBased
		err = json.Unmarshal(*o, &t)
		if err != nil {
			return nil, err
		}
		return t.polish()
	case enum.TemplateTypePodBased:
		var t rawTemplatePodBased
		err = json.Unmarshal(*o, &t)
		if err != nil {
			return nil, err
		}
		return t.polish()
	case enum.TemplateTypeL3Collapsed:
		var t rawTemplateL3Collapsed
		err = json.Unmarshal(*o, &t)
		if err != nil {
			return nil, err
		}
		return t.polish()
	}
	return nil, fmt.Errorf("unhandled template type: %s", t)
}

type DhcpServiceIntent struct {
	Active bool `json:"active"`
}

func (o *Client) listAllTemplateIds(ctx context.Context) ([]ObjectId, error) {
	response := &struct {
		Items []ObjectId `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

// getTemplate returns one of *TemplateRackBased, *TemplatePodBased or
// *TemplateL3Collapsed, each of which have Type() method which should be
// used to cast them into the correct type.
func (o *Client) getTemplate(ctx context.Context, id ObjectId) (template, error) {
	rawMsg := &json.RawMessage{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignTemplateById, id),
		apiResponse: rawMsg,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return template(*rawMsg), nil
}

func (o *Client) getAllTemplates(ctx context.Context) ([]template, error) {
	response := &struct {
		Items []json.RawMessage `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]template, len(response.Items))
	for i, item := range response.Items {
		result[i] = template(item)
	}

	return result, nil
}

func (o *Client) getTemplateByTypeAndName(ctx context.Context, desiredType enum.TemplateType, desiredName string) (*template, error) {
	templates, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var found *template
	for i := range templates {
		foundType, err := templates[i].templateType()
		if err != nil {
			return nil, err
		}
		if foundType != desiredType {
			continue // wrong type
		}

		foundName, err := templates[i].displayName()
		if foundName != desiredName {
			continue // wrong name
		}
		if err != nil {
			return nil, err
		}

		if found != nil { // multiple matches!
			return nil, ClientErr{
				errType: ErrMultipleMatch,
				err:     fmt.Errorf("found multiple %s templates named '%s'", desiredType, desiredName),
			}
		}

		// record this pointer to detect multiple matches
		found = &templates[i]
	}

	if found == nil { // not found!
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no %s templates named '%s'", desiredType, desiredName),
		}
	}

	return found, nil
}

func (o *Client) deleteTemplate(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignTemplateById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) getTemplateIdsTypesByName(ctx context.Context, desired string) (map[ObjectId]enum.TemplateType, error) {
	response := &struct {
		Items []struct {
			Id          ObjectId          `json:"id"`
			Type        enum.TemplateType `json:"type"`
			DisplayName string            `json:"display_name"`
		} `json:"Items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make(map[ObjectId]enum.TemplateType)
	for _, t := range response.Items {
		if t.DisplayName == desired {
			result[t.Id] = t.Type
		}
	}
	return result, nil
}

func (o *Client) getTemplateIdTypeByName(ctx context.Context, desired string) (ObjectId, enum.TemplateType, error) {
	idToType, err := o.getTemplateIdsTypesByName(ctx, desired)
	if err != nil {
		return "", enum.TemplateType{}, fmt.Errorf("error fetching templates by name - %w", err)
	}

	switch len(idToType) {
	case 0:
		return "", enum.TemplateType{}, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("template named '%s' not found", desired),
		}
	case 1:
		for k, v := range idToType {
			return k, v, nil
		}
	}
	return "", enum.TemplateType{}, ClientErr{
		errType: ErrMultipleMatch,
		err:     fmt.Errorf("found multiple templates named '%s'", desired),
	}
}
