// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	apiUrlBlueprintExperienceWebSystemInfo = apiUrlBlueprintByIdPrefix + "experience/web/system-info"
	apiUrlBlueprintSetNodeDomain           = apiUrlBlueprintByIdPrefix + "systems" + apiUrlPathDelim + "%s" + apiUrlPathDelim + "domain"
	apiUrlBlueprintSetNodeLoopback         = apiUrlBlueprintByIdPrefix + "systems" + apiUrlPathDelim + "%s" + apiUrlPathDelim + "loopback" + apiUrlPathDelim + "%d"
	apiUrlBlueprintSetPortChannelIdMinMax  = apiUrlBlueprintByIdPrefix + "port-channel-id"
)

type rawSystemNodeInfoLoopback struct {
	Ipv6Addr *string   `json:"ipv6_addr"`
	Ipv4Addr *string   `json:"ipv4_addr"`
	Id       *ObjectId `json:"id"`
}

type rawSystemNodeInfo struct {
	DomainId *string `json:"domain_id"`
	// UplinkedSystemIds interface{} `json:"uplinked_system_ids"`
	DeployMode     *string                    `json:"deploy_mode"`
	Id             ObjectId                   `json:"id"`
	Label          string                     `json:"label"`
	PodId          ObjectId                   `json:"pod_id"`
	InterfaceMapId *ObjectId                  `json:"interface_map_id"`
	Loopback       *rawSystemNodeInfoLoopback `json:"loopback"`
	Hostname       string                     `json:"hostname"`
	// UnicastVtep        interface{}            `json:"unicast_vtep"`
	Role systemRole `json:"role"`
	// SystemId           interface{}   `json:"system_id"`
	// Hidden             bool          `json:"hidden"`
	// AnycastVtep        interface{}   `json:"anycast_vtep"`
	RackId *ObjectId `json:"rack_id"`
	// LogicalVtep        interface{}   `json:"logical_vtep"`
	// SuperspinePlaneId  interface{}   `json:"superspine_plane_id"`
	LogicalDeviceId ObjectId `json:"logical_device_id"`
	Tags            []string `json:"tags"`
	// HypervisorId       interface{} `json:"hypervisor_id"`
	// RedundancyProtocol interface{} `json:"redundancy_protocol"`
	External         bool `json:"external"`
	PortChannelIdMin int  `json:"port_channel_id_min"`
	// PositionData      interface{} `json:"position_data"`
	PortChannelIdMax  int                   `json:"port_channel_id_max"`
	RedundancyGroupId *ObjectId             `json:"redundancy_group_id"`
	GroupLabel        *string               `json:"group_label"`
	ManagementLevel   systemManagementLevel `json:"management_level"`
	DeviceProfileId   *ObjectId             `json:"device_profile_id"`
}

func (o rawSystemNodeInfo) polish() (*SystemNodeInfo, error) {
	var asn *uint32
	if o.DomainId != nil {
		i64, err := strconv.ParseUint(*o.DomainId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("while parsing ASN string %q - %w", *o.DomainId, err)
		}
		i32 := uint32(i64)
		asn = &i32
	}

	var deployMode *enum.DeployMode
	if o.DeployMode != nil {
		deployMode = enum.DeployModes.Parse(*o.DeployMode)
		if deployMode == nil {
			return nil, fmt.Errorf("failed to parse deploy mode %q", deployMode)
		}
	}

	var loopbackId *ObjectId
	var loopbackIPv4, loopbackIPv6 *net.IPNet
	if o.Loopback != nil {
		loopbackId = o.Loopback.Id
		if o.Loopback.Ipv4Addr != nil {
			ip, net, err := net.ParseCIDR(*o.Loopback.Ipv4Addr)
			if err != nil {
				return nil, fmt.Errorf("while parsing Loopback IPv4 address %q - %w", *o.Loopback.Ipv4Addr, err)
			}
			net.IP = ip
			loopbackIPv4 = net
		}
		if o.Loopback.Ipv6Addr != nil {
			ip, net, err := net.ParseCIDR(*o.Loopback.Ipv6Addr)
			if err != nil {
				return nil, fmt.Errorf("while parsing Loopback IPv6 address %q - %w", *o.Loopback.Ipv6Addr, err)
			}
			net.IP = ip
			loopbackIPv6 = net
		}
	}

	managementLevel, err := o.ManagementLevel.parse()
	if err != nil {
		return nil, err
	}

	role, err := o.Role.parse()
	if err != nil {
		return nil, err
	}

	return &SystemNodeInfo{
		Asn:               asn,
		DeployMode:        deployMode,
		DeviceProfileId:   o.DeviceProfileId,
		External:          o.External,
		GroupLabel:        o.GroupLabel,
		Hostname:          o.Hostname,
		Id:                o.Id,
		InterfaceMapId:    o.InterfaceMapId,
		Label:             o.Label,
		LogicalDevice:     o.LogicalDeviceId,
		LoopbackId:        loopbackId,
		LoopbackIpv4:      loopbackIPv4,
		LoopbackIpv6:      loopbackIPv6,
		ManagementLevel:   SystemManagementLevel(managementLevel),
		PodId:             o.PodId,
		PortChannelIdMax:  o.PortChannelIdMax,
		PortChannelIdMin:  o.PortChannelIdMin,
		RackId:            o.RackId,
		RedundancyGroupId: o.RedundancyGroupId,
		Role:              SystemRole(role),
		Tags:              o.Tags,
	}, nil
}

type SystemNodeInfo struct {
	Asn               *uint32
	DeployMode        *enum.DeployMode
	DeviceProfileId   *ObjectId
	External          bool
	GroupLabel        *string
	Hostname          string
	Id                ObjectId
	InterfaceMapId    *ObjectId
	Label             string
	LogicalDevice     ObjectId
	LoopbackId        *ObjectId
	LoopbackIpv4      *net.IPNet
	LoopbackIpv6      *net.IPNet
	ManagementLevel   SystemManagementLevel
	PodId             ObjectId
	PortChannelIdMax  int
	PortChannelIdMin  int
	RackId            *ObjectId
	RedundancyGroupId *ObjectId
	Role              SystemRole
	Tags              []string
}

func (o *TwoStageL3ClosClient) getAllSystemNodeInfos(ctx context.Context) ([]rawSystemNodeInfo, error) {
	var apiResponse struct {
		Data []rawSystemNodeInfo `json:"data"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintExperienceWebSystemInfo, o.blueprintId),
		apiResponse: &apiResponse,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return apiResponse.Data, nil
}

func (o *TwoStageL3ClosClient) SetGenericSystemAsn(ctx context.Context, id ObjectId, asn *uint32) error {
	apiInput := struct {
		DomainId *uint32 `json:"domain_id"`
	}{
		DomainId: asn,
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSetNodeDomain, o.blueprintId, id),
		apiInput: &apiInput,
	})
	return convertTtaeToAceWherePossible(err)
}

// SetGenericSystemLoopbackIPs configures the loopback interface identified by
// gslo.LoopbackId belonging to the system identified by gsId to use the ipv4
// and ipv6 addresses specified in gslo. the LoopbackNodeId and SecurityZoneId
// fields within gslo are ignored.
func (o *TwoStageL3ClosClient) SetGenericSystemLoopbackIPs(ctx context.Context, gsId ObjectId, gslo GenericSystemLoopback) error {
	var apiInput struct {
		Ipv4Addr *string `json:"ipv4_addr"` // this is a pointer so we can clear values with `null`
		Ipv6Addr *string `json:"ipv6_addr"` // this is a pointer so we can clear values with `null`
	}

	if gslo.Ipv4Addr != nil {
		if gslo.Ipv4Addr.IP.To4() == nil {
			return fmt.Errorf("ip4 value does not contain a valid IPv4 address - %X", []byte(gslo.Ipv4Addr.IP))
		}

		maskOnes, maskBits := gslo.Ipv4Addr.Mask.Size()
		if maskBits != 32 || maskOnes != maskBits {
			return errors.New("ip4 value does not contain a valid mask for loopback interfaces - " + gslo.Ipv4Addr.Mask.String())
		}

		apiInput.Ipv4Addr = toPtr(gslo.Ipv4Addr.String())
	}

	if gslo.Ipv6Addr != nil {
		if gslo.Ipv6Addr.IP.To16() == nil || !strings.Contains(gslo.Ipv6Addr.IP.String(), ":") {
			return fmt.Errorf("ip6 value does not contain a valid IPv6 address - %X", []byte(gslo.Ipv6Addr.IP))
		}

		maskOnes, maskBits := gslo.Ipv6Addr.Mask.Size()
		if maskBits != 128 || maskOnes != maskBits {
			return errors.New("ip6 value does not contain a valid mask for loopback interfaces - " + gslo.Ipv6Addr.Mask.String())
		}

		apiInput.Ipv6Addr = toPtr(gslo.Ipv6Addr.String())
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSetNodeLoopback, o.blueprintId, gsId, gslo.LoopbackId),
		apiInput: &apiInput,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) SetGenericSystemLoopbackIpv4(ctx context.Context, id ObjectId, ipNet *net.IPNet, instance int) error {
	apiInput := struct {
		Ipv4Addr *string `json:"ipv4_addr"`
	}{}

	if ipNet != nil {
		s := ipNet.String()
		apiInput.Ipv4Addr = &s
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSetNodeLoopback, o.blueprintId, id, instance),
		apiInput: &apiInput,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) SetGenericSystemLoopbackIpv6(ctx context.Context, id ObjectId, ipNet *net.IPNet, instance int) error {
	apiInput := struct {
		Ipv6Addr *string `json:"ipv6_addr"`
	}{}

	if ipNet != nil {
		s := ipNet.String()
		apiInput.Ipv6Addr = &s
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSetNodeLoopback, o.blueprintId, id, instance),
		apiInput: &apiInput,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) SetGenericSystemPortChannelMinMax(ctx context.Context, id ObjectId, min, max int) error {
	type portChannelStruct struct {
		SystemId         ObjectId `json:"system_id"`
		PortChannelIdMin int      `json:"port_channel_id_min"`
		PortChannelIdMax int      `json:"port_channel_id_max"`
	}

	type apiInput struct {
		Systems []portChannelStruct `json:"systems"`
	}

	input := apiInput{
		Systems: []portChannelStruct{
			{
				PortChannelIdMin: min,
				PortChannelIdMax: max,
				SystemId:         id,
			},
		},
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSetPortChannelIdMinMax, o.blueprintId),
		apiInput: &input,
	})
	return convertTtaeToAceWherePossible(err)
}
