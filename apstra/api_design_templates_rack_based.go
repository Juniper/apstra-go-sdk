// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Juniper/apstra-go-sdk/compatibility"
)

var _ Template = &TemplateRackBased{}

type TemplateRackBased struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	templateType   TemplateType
	Data           *TemplateRackBasedData
}

func (o *TemplateRackBased) Type() TemplateType {
	return o.templateType
}

func (o *TemplateRackBased) ID() ObjectId {
	return o.Id
}

func (o *TemplateRackBased) OverlayControlProtocol() OverlayControlProtocol {
	if o == nil || o.Data == nil {
		return OverlayControlProtocolNone
	}
	return o.Data.VirtualNetworkPolicy.OverlayControlProtocol
}

type rawTemplateRackBased struct {
	Id                   ObjectId                `json:"id"`
	Type                 templateType            `json:"type"`
	DisplayName          string                  `json:"display_name"`
	AntiAffinityPolicy   *rawAntiAffinityPolicy  `json:"anti_affinity_policy,omitempty"`
	CreatedAt            time.Time               `json:"created_at"`
	LastModifiedAt       time.Time               `json:"last_modified_at"`
	VirtualNetworkPolicy rawVirtualNetworkPolicy `json:"virtual_network_policy"`
	AsnAllocationPolicy  rawAsnAllocationPolicy  `json:"asn_allocation_policy"`
	Capability           templateCapability      `json:"capability,omitempty"`
	Spine                rawSpine                `json:"spine"`
	RackTypes            []rawRackType           `json:"rack_types"`
	RackTypeCounts       []RackTypeCount         `json:"rack_type_counts"`
	DhcpServiceIntent    DhcpServiceIntent       `json:"dhcp_service_intent"`
}

func (o rawTemplateRackBased) polish() (*TemplateRackBased, error) {
	tType, err := o.Type.parse()
	if err != nil {
		return nil, err
	}
	v, err := o.VirtualNetworkPolicy.polish()
	if err != nil {
		return nil, err
	}
	a, err := o.AsnAllocationPolicy.polish()
	if err != nil {
		return nil, err
	}
	c, err := o.Capability.parse()
	if err != nil {
		return nil, err
	}
	s, err := o.Spine.polish()
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

	if len(o.RackTypes) != len(o.RackTypeCounts) {
		return nil, fmt.Errorf("template '%s' has %d rack_types and %d rack_type_counts - these should match",
			o.Id, len(o.RackTypes), len(o.RackTypeCounts))
	}

	rackTypeInfos := make(map[ObjectId]TemplateRackBasedRackInfo, len(o.RackTypes))
OUTER:
	for _, rrt := range o.RackTypes { // loop over raw rack types
		prt, err := rrt.polish()
		if err != nil {
			return nil, err
		}
		for _, rtc := range o.RackTypeCounts { // loop over rack type counts looking for matching ID
			if prt.Id == rtc.RackTypeId {
				rackTypeInfos[rtc.RackTypeId] = TemplateRackBasedRackInfo{
					Count:        rtc.Count,
					RackTypeData: prt.Data,
				}
				continue OUTER
			}
		}
		return nil, fmt.Errorf("template contains rack_type '%s' which does not appear among rack_type_counts", rrt.Id)
	}

	return &TemplateRackBased{
		Id:             o.Id,
		templateType:   TemplateType(tType),
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &TemplateRackBasedData{
			DisplayName:          o.DisplayName,
			AntiAffinityPolicy:   aap,
			VirtualNetworkPolicy: *v,
			AsnAllocationPolicy:  *a,
			Capability:           TemplateCapability(c),
			Spine:                *s,
			RackInfo:             rackTypeInfos,
			DhcpServiceIntent:    o.DhcpServiceIntent,
		},
	}, nil
}

type TemplateRackBasedData struct {
	DisplayName          string
	AntiAffinityPolicy   *AntiAffinityPolicy
	VirtualNetworkPolicy VirtualNetworkPolicy
	AsnAllocationPolicy  AsnAllocationPolicy
	Capability           TemplateCapability
	Spine                Spine
	RackInfo             map[ObjectId]TemplateRackBasedRackInfo
	DhcpServiceIntent    DhcpServiceIntent
}

type AsnAllocationPolicy struct {
	SpineAsnScheme AsnAllocationScheme
}

func (o *AsnAllocationPolicy) raw() *rawAsnAllocationPolicy {
	return &rawAsnAllocationPolicy{SpineAsnScheme: o.SpineAsnScheme.raw()}
}

type rawAsnAllocationPolicy struct {
	SpineAsnScheme asnAllocationScheme `json:"spine_asn_scheme"`
}

func (o *rawAsnAllocationPolicy) polish() (*AsnAllocationPolicy, error) {
	sas, err := o.SpineAsnScheme.parse()
	return &AsnAllocationPolicy{SpineAsnScheme: AsnAllocationScheme(sas)}, err
}

type Spine struct {
	Count                  int
	LinkPerSuperspineSpeed LogicalDevicePortSpeed
	LogicalDevice          LogicalDeviceData
	LinkPerSuperspineCount int
	Tags                   []DesignTagData
}

type TemplateElementSpineRequest struct {
	Count                  int
	LinkPerSuperspineSpeed LogicalDevicePortSpeed
	LogicalDevice          ObjectId
	LinkPerSuperspineCount int
	Tags                   []ObjectId
}

func (o *TemplateElementSpineRequest) raw(ctx context.Context, client *Client) (*rawSpine, error) {
	logicalDevice, err := client.getLogicalDevice(ctx, o.LogicalDevice)
	if err != nil {
		return nil, err
	}

	tags := make([]DesignTagData, len(o.Tags))
	for i, tagId := range o.Tags {
		rawTag, err := client.getTag(ctx, tagId)
		if err != nil {
			return nil, err
		}
		tags[i] = *rawTag.polish().Data
	}

	return &rawSpine{
		Count:                  o.Count,
		LinkPerSuperspineSpeed: o.LinkPerSuperspineSpeed.raw(),
		LogicalDevice:          *logicalDevice,
		LinkPerSuperspineCount: o.LinkPerSuperspineCount,
		Tags:                   tags,
	}, nil
}

type rawSpine struct {
	Count                  int                        `json:"count"`
	LinkPerSuperspineSpeed *rawLogicalDevicePortSpeed `json:"link_per_superspine_speed"`
	LogicalDevice          rawLogicalDevice           `json:"logical_device"`
	LinkPerSuperspineCount int                        `json:"link_per_superspine_count"`
	Tags                   []DesignTagData            `json:"tags"`
}

func (o rawSpine) polish() (*Spine, error) {
	ld, err := o.LogicalDevice.polish()

	var linkPerSuperspineSpeed LogicalDevicePortSpeed
	if o.LinkPerSuperspineSpeed != nil {
		linkPerSuperspineSpeed = o.LinkPerSuperspineSpeed.parse()
	}

	return &Spine{
		Count:                  o.Count,
		LinkPerSuperspineSpeed: linkPerSuperspineSpeed,
		LogicalDevice: LogicalDeviceData{
			DisplayName: ld.Data.DisplayName,
			Panels:      ld.Data.Panels,
		},
		LinkPerSuperspineCount: o.LinkPerSuperspineCount,
		Tags:                   o.Tags,
	}, err
}

type TemplateRackBasedRackInfo struct {
	Count        int
	RackTypeData *RackTypeData
}

func (o *Client) getRackBasedTemplate(ctx context.Context, id ObjectId) (*rawTemplateRackBased, error) {
	rawTemplate, err := o.getTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	tType, err := rawTemplate.templateType()
	if err != nil {
		return nil, err
	}

	if tType != templateTypeRackBased {
		return nil, ClientErr{
			errType: ErrWrongType,
			err:     fmt.Errorf("template '%s' is of type '%s', not '%s'", id, tType, templateTypeRackBased),
		}
	}

	template := &rawTemplateRackBased{}
	return template, json.Unmarshal(rawTemplate, template)
}

func (o *Client) getAllRackBasedTemplates(ctx context.Context) ([]rawTemplateRackBased, error) {
	templates, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var result []rawTemplateRackBased
	for _, t := range templates {
		tType, err := t.templateType()
		if err != nil {
			return nil, err
		}
		if tType != templateTypeRackBased {
			continue
		}
		var raw rawTemplateRackBased
		err = json.Unmarshal(t, &raw)
		if err != nil {
			return nil, err
		}
		result = append(result, raw)
	}

	return result, nil
}

type CreateRackBasedTemplateRequest struct {
	DisplayName          string
	Spine                *TemplateElementSpineRequest
	RackInfos            map[ObjectId]TemplateRackBasedRackInfo
	DhcpServiceIntent    *DhcpServiceIntent
	AntiAffinityPolicy   *AntiAffinityPolicy
	AsnAllocationPolicy  *AsnAllocationPolicy
	VirtualNetworkPolicy *VirtualNetworkPolicy
}

func (o *CreateRackBasedTemplateRequest) raw(ctx context.Context, client *Client) (*rawCreateRackBasedTemplateRequest, error) {
	rackTypes := make([]rawRackType, len(o.RackInfos))
	rackTypeCounts := make([]RackTypeCount, len(o.RackInfos))
	var i int
	for k, ri := range o.RackInfos {
		if ri.RackTypeData != nil {
			return nil, fmt.Errorf("the RackTypeData field must be nil when creating a rack-based template")
		}

		// grab the rack type from the API using the caller's map key (ObjectId) and stash it in rackTypes
		rt, err := client.getRackType(ctx, k)
		if err != nil {
			return nil, err
		}
		rackTypes[i] = *rt

		// prep the RackTypeCount object using the caller's map key (ObjectId) as
		// the link between the rawRackType data copy and the RackTypeCount
		rackTypeCounts[i].RackTypeId = k
		rackTypeCounts[i].Count = ri.Count
		i++
	}

	switch {
	case o.Spine == nil:
		return nil, errors.New("spine cannot be <nil> when creating a rack-based template")
	case o.AntiAffinityPolicy == nil && compatibility.TemplateRequestRequiresAntiAffinityPolicy.Check(client.apiVersion):
		return nil, fmt.Errorf("anti-affinity policy cannot be <nil> when creating a rack-based template with Apstra %s", compatibility.TemplateRequestRequiresAntiAffinityPolicy)
	case o.AsnAllocationPolicy == nil:
		return nil, errors.New("asn allocation policy cannot be <nil> when creating a rack-based template")
	case o.VirtualNetworkPolicy == nil:
		return nil, errors.New("virtual network policy cannot be <nil> when creating a rack-based template")
	}

	var err error
	var dhcpServiceIntent DhcpServiceIntent
	if o.DhcpServiceIntent != nil {
		dhcpServiceIntent = *o.DhcpServiceIntent
	}

	spine, err := o.Spine.raw(ctx, client)
	if err != nil {
		return nil, err
	}

	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	asnAllocationPolicy := o.AsnAllocationPolicy.raw()

	virtualNetworkPolicy := o.VirtualNetworkPolicy.raw()

	return &rawCreateRackBasedTemplateRequest{
		Type:                 templateTypeRackBased,
		DisplayName:          o.DisplayName,
		Spine:                *spine,
		RackTypes:            rackTypes,
		RackTypeCounts:       rackTypeCounts,
		DhcpServiceIntent:    dhcpServiceIntent,
		AntiAffinityPolicy:   antiAffinityPolicy,
		AsnAllocationPolicy:  *asnAllocationPolicy,
		VirtualNetworkPolicy: *virtualNetworkPolicy,
	}, nil
}

type rawCreateRackBasedTemplateRequest struct {
	Type                 templateType            `json:"type"`
	DisplayName          string                  `json:"display_name"`
	Spine                rawSpine                `json:"spine"`
	RackTypes            []rawRackType           `json:"rack_types"`
	RackTypeCounts       []RackTypeCount         `json:"rack_type_counts"`
	DhcpServiceIntent    DhcpServiceIntent       `json:"dhcp_service_intent"`
	AntiAffinityPolicy   *rawAntiAffinityPolicy  `json:"anti_affinity_policy,omitempty"`
	AsnAllocationPolicy  rawAsnAllocationPolicy  `json:"asn_allocation_policy"`
	VirtualNetworkPolicy rawVirtualNetworkPolicy `json:"virtual_network_policy"`
}

func (o *Client) createRackBasedTemplate(ctx context.Context, in *rawCreateRackBasedTemplateRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignTemplates,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateRackBasedTemplate(ctx context.Context, id ObjectId, in *CreateRackBasedTemplateRequest) error {
	raw, err := in.raw(ctx, o)
	if err != nil {
		return err
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignTemplateById, id),
		apiInput: raw,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
