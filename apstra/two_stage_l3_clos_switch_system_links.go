// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	apiUrlSwitchSystemLinks              = apiUrlBlueprintByIdPrefix + "switch-system-links"
	apiUrlDeleteSwitchSystemLinks        = apiUrlBlueprintByIdPrefix + "delete-switch-system-links"
	apiUrlBlueprintExternalGenericSystem = apiUrlBlueprintByIdPrefix + "external-generic-systems" + apiUrlPathDelim + "%s"
)

type (
	SystemType int
	systemType string
)

const (
	SystemTypeExternal = SystemType(iota)
	SystemTypeInternal
	SystemTypeSwitch
	SystemTypeServer
	SystemTypeUnknown = "unknown system type '%s'"

	systemTypeExternal = systemType("external")
	systemTypeInternal = systemType("internal")
	systemTypeSwitch   = systemType("switch")
	systemTypeServer   = systemType("server")
	systemTypeUnknown  = "unknown system type %d"
)

func (o SystemType) Int() int {
	return int(o)
}

func (o SystemType) String() string {
	switch o {
	case SystemTypeExternal:
		return string(systemTypeExternal)
	case SystemTypeInternal:
		return string(systemTypeInternal)
	case SystemTypeSwitch:
		return string(systemTypeSwitch)
	case SystemTypeServer:
		return string(systemTypeServer)
	default:
		return fmt.Sprintf(systemTypeUnknown, o)
	}
}

func (o SystemType) raw() systemType {
	return systemType(o.String())
}

func (o systemType) string() string {
	return string(o)
}

func (o systemType) parse() (int, error) {
	switch o {
	case systemTypeExternal:
		return int(SystemTypeExternal), nil
	case systemTypeInternal:
		return int(SystemTypeInternal), nil
	case systemTypeSwitch:
		return int(SystemTypeSwitch), nil
	case systemTypeServer:
		return int(SystemTypeServer), nil
	default:
		return 0, fmt.Errorf(SystemTypeUnknown, o)
	}
}

type CreateLinksWithNewSystemRequest struct {
	Links  []CreateLinkRequest
	System CreateLinksWithNewSystemRequestSystem
}

func (o *CreateLinksWithNewSystemRequest) raw(ctx context.Context, client *Client) (*rawCreateLinksWithNewSystemRequest, error) {
	rs, err := o.System.raw(ctx, client)
	if err != nil {
		return nil, err
	}

	links := make([]rawCreateLinkRequest, len(o.Links))
	for i, link := range o.Links {
		links[i] = *link.raw()
	}

	return &rawCreateLinksWithNewSystemRequest{
		NewSystems: []rawCreateLinksWithNewSystemRequestSystem{*rs},
		Links:      links,
	}, nil
}

type rawCreateLinksWithNewSystemRequest struct {
	NewSystems []rawCreateLinksWithNewSystemRequestSystem `json:"new_systems,omitempty"`
	Links      []rawCreateLinkRequest                     `json:"links"`
}

type CreateLinksWithNewSystemRequestSystem struct {
	Hostname         string
	Label            string
	LogicalDeviceId  ObjectId
	LogicalDevice    *LogicalDevice
	PortChannelIdMin int
	PortChannelIdMax int
	Tags             []string
	Type             SystemType
}

func (o *CreateLinksWithNewSystemRequestSystem) raw(ctx context.Context, client *Client) (*rawCreateLinksWithNewSystemRequestSystem, error) {
	if o.LogicalDeviceId != "" && o.LogicalDevice != nil {
		return nil, errors.New("both LogicalDevice (payload) and LogicalDeviceId (catalog ID) specified")
	}

	var err error
	var rawLD *rawLogicalDevice

	if o.LogicalDeviceId != "" {
		rawLD, err = client.getLogicalDevice(ctx, o.LogicalDeviceId)
		if err != nil {
			return nil, fmt.Errorf("error fetching logical device %q - %w", o.LogicalDeviceId, err)
		}
		// wipe out the timestamps so we don't send 'em back to Apstra
		rawLD.CreatedAt = nil
		rawLD.LastModifiedAt = nil
	}

	if o.LogicalDevice != nil {
		rawLD = o.LogicalDevice.raw()
	}

	return &rawCreateLinksWithNewSystemRequestSystem{
		SystemType:       o.Type.String(),
		LogicalDevice:    *rawLD,
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Tags:             o.Tags,
		Label:            o.Label,
		Hostname:         o.Hostname,
	}, nil
}

type rawCreateLinksWithNewSystemRequestSystem struct {
	SystemType       string           `json:"system_type"`         // mandatory
	LogicalDevice    rawLogicalDevice `json:"logical_device"`      // mandatory
	PortChannelIdMin int              `json:"port_channel_id_min"` // mandatory; 0 is default
	PortChannelIdMax int              `json:"port_channel_id_max"` // mandatory; 0 is default
	Tags             []string         `json:"tags,omitempty"`
	Label            string           `json:"label,omitempty"`
	Hostname         string           `json:"hostname,omitempty"`
}

type CreateLinkRequest struct {
	Tags           []string
	SystemEndpoint SwitchLinkEndpoint
	SwitchEndpoint SwitchLinkEndpoint
	GroupLabel     string
	LagMode        RackLinkLagMode
}

func (o *CreateLinkRequest) raw() *rawCreateLinkRequest {
	return &rawCreateLinkRequest{
		Tags:           o.Tags,
		SystemEndpoint: o.SystemEndpoint.raw(),
		SwitchEndpoint: o.SwitchEndpoint.raw(),
		GroupLabel:     o.GroupLabel,
		LagMode:        rackLinkLagMode(o.LagMode.String()),
	}
}

type rawCreateLinkRequest struct {
	Tags           []string              `json:"tags,omitempty"`
	SystemEndpoint rawSwitchLinkEndpoint `json:"system"`
	SwitchEndpoint rawSwitchLinkEndpoint `json:"switch"`
	GroupLabel     string                `json:"link_group_label,omitempty"`
	LagMode        rackLinkLagMode       `json:"lag_mode,omitempty"`
}

type SwitchLinkEndpoint struct {
	TransformationId int
	SystemId         ObjectId
	IfName           string
	// LagMode          RackLinkLagMode
}

func (o *SwitchLinkEndpoint) raw() rawSwitchLinkEndpoint {
	var systemIdPtr *ObjectId
	if s := o.SystemId; s != "" {
		systemIdPtr = &s
	}

	return rawSwitchLinkEndpoint{
		TransformationId: o.TransformationId,
		SystemId:         systemIdPtr,
		IfName:           o.IfName,
	}
}

type rawSwitchLinkEndpoint struct {
	TransformationId int       `json:"transformation_id,omitempty"`
	SystemId         *ObjectId `json:"system_id"` // must send `null` when creating a new system, so no `omitempty`
	IfName           string    `json:"if_name,omitempty"`
	LagMode          string    `json:"lag_mode,omitempty"`
}

func (o *TwoStageL3ClosClient) CreateLinksWithNewSystem(ctx context.Context, req *CreateLinksWithNewSystemRequest) ([]ObjectId, error) {
	apiInput, err := req.raw(ctx, o.Client())
	if err != nil {
		return nil, fmt.Errorf("error processing CreateLinksWithNewSystemRequest, - %w", err)
	}

	apiResponse := struct {
		IDs []ObjectId `json:"ids"`
	}{}

	err = o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlSwitchSystemLinks, o.blueprintId),
		apiInput:    apiInput,
		apiResponse: &apiResponse,
	})

	return apiResponse.IDs, convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) DeleteLinksFromSystem(ctx context.Context, ids []ObjectId) error {
	apiInput := struct {
		LinkIds []ObjectId `json:"link_ids"`
	}{ids}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   fmt.Sprintf(apiUrlDeleteSwitchSystemLinks, o.blueprintId),
		apiInput: &apiInput,
	})
	if err == nil {
		return nil // success!
	}

	// if we got here, then we have an error
	err = convertTtaeToAceWherePossible(err)
	var ace ClientErr
	if !errors.As(err, &ace) {
		return err // cannot handle
	}

	if ace.Type() != ErrCtAssignedToLink {
		return err // cannot handle
	}

	var ds detailedStatus
	if json.Unmarshal([]byte(ace.Error()), &ds) != nil {
		return err // unmarshal fail - surface the original error
	}

	// unpack the error
	var e struct {
		LinkIds interface{} `json:"link_ids"` // this might be a string or a []string (sigh)
	}
	if json.Unmarshal(ds.Errors, &e) != nil {
		return err // unmarshal fail - surface the original error
	}

	// unpack e.LinkIds into a []string, whether it's returned as a string or a []string
	var linkIds []string
	switch t := e.LinkIds.(type) {
	case string:
		linkIds = []string{t}
	case []interface{}:
		linkIds = make([]string, len(t))
		for i := range t {
			if s, ok := t[i].(string); ok {
				linkIds[i] = s
			} else {
				return err // cannot handle non-string - surface the original error
			}
		}
	}

	// we know about two categories of error - use regexes to filter 'em out - examples:
	// 	 "Link with id l2_virtual_004_leaf1<->l2_virtual_004_sys003(link-000000002)[1] can not be deleted since some of its interfaces have connectivity templates assigned",
	//   "Deleting all links forming a LAG is not allowed since the LAG has assigned structures: ['connectivity template', 'VN endpoint']. Link ids: ['l2_virtual_003_leaf1<->l2_virtual_003_sys003(b)[1]', 'l2_virtual_003_leaf1<->l2_virtual_003_sys003(b)[2]']",
	var linkErrs []string
	var lagErrs []string
	for _, le := range linkIds {
		switch {
		case regexpLinkHasCtAssignedErr.MatchString(le):
			linkErrs = append(linkErrs, le)
		case regexpLagHasCtAssignedErr.MatchString(le):
			lagErrs = append(lagErrs, le)
		case regexpLinkHasVnEndpoint.MatchString(le):
			// do nothing - this condition should trigger the regexpLinkHasCtAssignedErr also
		case regexpLinkHasSubinterfaces.MatchString(le):
			// do nothing - this condition should trigger the regexpLinkHasCtAssignedErr also
		default: // cannot handle error - surface it to the user
			return fmt.Errorf("cannot handle link error %q - %w", le, err)
		}
	}

	// Collect the IDs of links with errors
	linkErrCount := len(linkErrs)
	lagErrCount := len(lagErrs)
	linkIdsWithCts := make([]ObjectId, linkErrCount+lagErrCount)

	// extract ids of naked links with errors
	for i, s := range linkErrs {
		m := regexpLinkHasCtAssignedErr.FindStringSubmatch(s)
		if len(m) != 2 {
			return fmt.Errorf("cannot handle link error %q - %w", s, err)
		}

		linkIdsWithCts[i] = ObjectId(m[1])
	}

	// determine ids of aggregate links with errors
	for i, s := range lagErrs {
		m := regexpLagHasCtAssignedErr.FindStringSubmatch(s)
		if len(m) != 2 {
			return fmt.Errorf("cannot handle lag link error %q - %w", s, err)
		}

		// each lag error enumerates all member links. Extract them
		var lagMembers []ObjectId
		for _, quotedId := range strings.Split(m[1], ",") {
			lagMembers = append(lagMembers, ObjectId(strings.Trim(quotedId, "' ")))
		}

		// find the LAG ID common to these member IDs
		lagId, err := o.lagIdFromMemberIds(ctx, lagMembers)
		if err != nil {
			return errors.Join(ace, err)
		}

		linkIdsWithCts[linkErrCount+i] = lagId
	}

	ace.detail = ErrCtAssignedToLinkDetail{LinkIds: linkIdsWithCts}
	return ace
}

func (o *TwoStageL3ClosClient) lagIdFromMemberIds(ctx context.Context, members []ObjectId) (ObjectId, error) {
	mq := new(MatchQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client)

	for _, member := range members {
		mq.Match(new(PathQuery).
			Node([]QEEAttribute{
				NodeTypeLink.QEEAttribute(),
				{Key: "id", Value: QEStringVal(member.String())},
			}).
			In([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
			Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
			In([]QEEAttribute{RelationshipTypeComposedOf.QEEAttribute()}).
			Node([]QEEAttribute{NodeTypeInterface.QEEAttribute()}).
			Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
			Node([]QEEAttribute{
				NodeTypeLink.QEEAttribute(),
				{Key: "link_type", Value: QEStringVal("aggregate_link")},
				{Key: "name", Value: QEStringVal("n_link")},
			}),
		)
	}

	var result struct {
		Items []struct {
			Link struct {
				Id ObjectId `json:"id"`
			} `json:"n_link"`
		} `json:"items"`
	}

	err := mq.Do(ctx, &result)
	if err != nil {
		return "", err
	}

	// turn result into a map keyed by link ID - we expect all results to use the same ID (one map entry)
	ids := make(map[ObjectId]struct{})
	for _, item := range result.Items {
		ids[item.Link.Id] = struct{}{}
	}

	switch len(ids) {
	case 0:
		return "", fmt.Errorf("member-based LAG member query found no LAG ID - %s", mq.String())
	case 1:
		return result.Items[0].Link.Id, nil // we expect exactly one map entry (all lag members point at one parent)
	default:
		return "", fmt.Errorf("member-based LAG member query found more than one LAG ID - %s", mq.String())
	}
}

func (o *TwoStageL3ClosClient) DeleteGenericSystem(ctx context.Context, id ObjectId) error {
	response := struct {
		Items []struct {
			Link struct {
				ID ObjectId `json:"id"`
			} `json:"n_link"`
			System struct {
				External bool `json:"external"`
			} `json:"n_system"`
		} `json:"items"`
	}{}
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetBlueprintType(BlueprintTypeStaging).
		SetClient(o.client).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{Key: "role", Value: QEStringVal("generic")},
			{Key: "id", Value: QEStringVal(id)},
			{Key: "name", Value: QEStringVal("n_system")},
		}).
		Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeInterface.QEEAttribute(),
			{Key: "if_type", Value: QEStringVal("ethernet")},
		}).
		Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
		Node([]QEEAttribute{
			NodeTypeLink.QEEAttribute(),
			{Key: "name", Value: QEStringVal("n_link")},
		})

	err := query.Do(ctx, &response)
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	if len(response.Items) == 0 {
		return ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("query %q returned no results", query),
		}
	}

	linkIds := make([]ObjectId, len(response.Items))
	for i, item := range response.Items {
		linkIds[i] = item.Link.ID
	}

	err = o.DeleteLinksFromSystem(ctx, linkIds)
	if err != nil {
		return fmt.Errorf("failed to delete external system %q - %w", id, err)
	}

	if response.Items[0].System.External {
		err = o.client.talkToApstra(ctx, talkToApstraIn{
			method: http.MethodDelete,
			urlStr: fmt.Sprintf(apiUrlBlueprintExternalGenericSystem, o.blueprintId, id),
		})
		if err != nil {
			return convertTtaeToAceWherePossible(err)
		}
	}

	return nil
}

func (o *TwoStageL3ClosClient) AddLinksToSystem(ctx context.Context, linkRequests []CreateLinkRequest) ([]ObjectId, error) {
	rawLinkRequests := make([]rawCreateLinkRequest, len(linkRequests))
	for i := range rawLinkRequests {
		rawLinkRequests[i] = *linkRequests[i].raw()
	}

	apiInput := rawCreateLinksWithNewSystemRequest{Links: rawLinkRequests}
	var apiResponse struct {
		Ids []ObjectId `json:"ids"`
	}

	err := o.client.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlSwitchSystemLinks, o.blueprintId),
		apiInput:    &apiInput,
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return apiResponse.Ids, nil
}
