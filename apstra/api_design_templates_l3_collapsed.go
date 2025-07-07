// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var _ Template = &TemplateL3Collapsed{}

type TemplateL3Collapsed struct {
	Id             ObjectId
	templateType   TemplateType
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *TemplateL3CollapsedData
}

func (o *TemplateL3Collapsed) Type() TemplateType {
	return o.templateType
}

func (o *TemplateL3Collapsed) ID() ObjectId {
	return o.Id
}

func (o *TemplateL3Collapsed) OverlayControlProtocol() OverlayControlProtocol {
	if o == nil || o.Data == nil {
		return OverlayControlProtocolNone
	}
	return o.Data.VirtualNetworkPolicy.OverlayControlProtocol
}

type rawTemplateL3Collapsed struct {
	Id                   ObjectId                   `json:"id"`
	Type                 templateType               `json:"type"`
	DisplayName          string                     `json:"display_name"`
	AntiAffinityPolicy   *rawAntiAffinityPolicy     `json:"anti_affinity_policy,omitempty"`
	CreatedAt            time.Time                  `json:"created_at"`
	LastModifiedAt       time.Time                  `json:"last_modified_at"`
	RackTypes            []rawRackType              `json:"rack_types"`
	Capability           templateCapability         `json:"capability"`
	MeshLinkSpeed        *rawLogicalDevicePortSpeed `json:"mesh_link_speed"`
	VirtualNetworkPolicy rawVirtualNetworkPolicy    `json:"virtual_network_policy"`
	MeshLinkCount        int                        `json:"mesh_link_count"`
	RackTypeCounts       []RackTypeCount            `json:"rack_type_counts"`
	DhcpServiceIntent    DhcpServiceIntent          `json:"dhcp_service_intent"`
}

func (o rawTemplateL3Collapsed) polish() (*TemplateL3Collapsed, error) {
	var prt []RackType
	for _, rrt := range o.RackTypes {
		polished, err := rrt.polish()
		if err != nil {
			return nil, err
		}
		prt = append(prt, *polished)
	}
	capability, err := o.Capability.parse()
	if err != nil {
		return nil, err
	}
	vnp, err := o.VirtualNetworkPolicy.polish()
	if err != nil {
		return nil, err
	}
	tType, err := o.Type.parse()
	if err != nil {
		return nil, err
	}
	var aap *AntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		aap, err = o.AntiAffinityPolicy.polish()
		if err != nil {
			return nil, err
		}
	}
	return &TemplateL3Collapsed{
		Id:             o.Id,
		templateType:   TemplateType(tType),
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &TemplateL3CollapsedData{
			DisplayName:          o.DisplayName,
			AntiAffinityPolicy:   aap,
			RackTypes:            prt,
			Capability:           TemplateCapability(capability),
			MeshLinkSpeed:        o.MeshLinkSpeed.parse(),
			VirtualNetworkPolicy: *vnp,
			MeshLinkCount:        o.MeshLinkCount,
			RackTypeCounts:       o.RackTypeCounts,
			DhcpServiceIntent:    o.DhcpServiceIntent,
		},
	}, nil
}

type TemplateL3CollapsedData struct {
	DisplayName          string
	AntiAffinityPolicy   *AntiAffinityPolicy
	RackTypes            []RackType
	Capability           TemplateCapability
	MeshLinkSpeed        LogicalDevicePortSpeed
	VirtualNetworkPolicy VirtualNetworkPolicy
	MeshLinkCount        int
	RackTypeCounts       []RackTypeCount
	DhcpServiceIntent    DhcpServiceIntent
}

func (o *Client) getL3CollapsedTemplate(ctx context.Context, id ObjectId) (*rawTemplateL3Collapsed, error) {
	rawTemplate, err := o.getTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	tType, err := rawTemplate.templateType()
	if err != nil {
		return nil, err
	}

	if tType != templateTypeL3Collapsed {
		return nil, ClientErr{
			errType: ErrWrongType,
			err:     fmt.Errorf("template '%s' is of type '%s', not '%s'", id, tType, templateTypeL3Collapsed),
		}
	}

	template := &rawTemplateL3Collapsed{}
	return template, json.Unmarshal(rawTemplate, template)
}

func (o *Client) getAllL3CollapsedTemplates(ctx context.Context) ([]rawTemplateL3Collapsed, error) {
	templates, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var result []rawTemplateL3Collapsed
	for _, t := range templates {
		tType, err := t.templateType()
		if err != nil {
			return nil, err
		}
		if tType != templateTypeL3Collapsed {
			continue
		}
		var raw rawTemplateL3Collapsed
		err = json.Unmarshal(t, &raw)
		if err != nil {
			return nil, err
		}
		result = append(result, raw)
	}

	return result, nil
}

type CreateL3CollapsedTemplateRequest struct {
	DisplayName          string                 `json:"display_name"`
	MeshLinkCount        int                    `json:"mesh_link_count"`
	MeshLinkSpeed        LogicalDevicePortSpeed `json:"mesh_link_speed"`
	RackTypeIds          []ObjectId             `json:"rack_types"`
	RackTypeCounts       []RackTypeCount        `json:"rack_type_counts"`
	DhcpServiceIntent    DhcpServiceIntent      `json:"dhcp_service_intent"`
	AntiAffinityPolicy   *AntiAffinityPolicy    `json:"anti_affinity_policy,omitempty"`
	VirtualNetworkPolicy VirtualNetworkPolicy   `json:"virtual_network_policy"`
}

func (o *CreateL3CollapsedTemplateRequest) raw(ctx context.Context, client *Client) (*rawCreateL3CollapsedTemplateRequest, error) {
	rackTypes := make([]rawRackType, len(o.RackTypeIds))
	for i, id := range o.RackTypeIds {
		rt, err := client.getRackType(ctx, id)
		if err != nil {
			return nil, err
		}
		rackTypes[i] = *rt
	}

	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	return &rawCreateL3CollapsedTemplateRequest{
		Type:                 templateTypeL3Collapsed,
		DisplayName:          o.DisplayName,
		MeshLinkCount:        o.MeshLinkCount,
		MeshLinkSpeed:        *o.MeshLinkSpeed.raw(),
		RackTypes:            rackTypes,
		RackTypeCounts:       o.RackTypeCounts,
		DhcpServiceIntent:    o.DhcpServiceIntent,
		AntiAffinityPolicy:   antiAffinityPolicy,
		VirtualNetworkPolicy: *o.VirtualNetworkPolicy.raw(),
	}, nil
}

type rawCreateL3CollapsedTemplateRequest struct {
	Type                 templateType              `json:"type"`
	DisplayName          string                    `json:"display_name"`
	MeshLinkCount        int                       `json:"mesh_link_count"`
	MeshLinkSpeed        rawLogicalDevicePortSpeed `json:"mesh_link_speed"`
	RackTypes            []rawRackType             `json:"rack_types"`
	RackTypeCounts       []RackTypeCount           `json:"rack_type_counts"`
	DhcpServiceIntent    DhcpServiceIntent         `json:"dhcp_service_intent"`
	AntiAffinityPolicy   *rawAntiAffinityPolicy    `json:"anti_affinity_policy,omitempty"`
	VirtualNetworkPolicy rawVirtualNetworkPolicy   `json:"virtual_network_policy"`
}

func (o *Client) createL3CollapsedTemplate(ctx context.Context, in *rawCreateL3CollapsedTemplateRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignTemplates,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", err
	}

	return response.Id, nil
}

func (o *Client) updateL3CollapsedTemplate(ctx context.Context, id ObjectId, in *CreateL3CollapsedTemplateRequest) error {
	apiInput, err := in.raw(ctx, o)
	if err != nil {
		return err
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignTemplateById, id),
		apiInput: apiInput,
	})
	if err != nil {
		return err
	}

	return nil
}
