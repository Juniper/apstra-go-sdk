package apstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

const (
	apiUrlBlueprintExperienceWebSystemInfo = apiUrlBlueprintById + apiUrlPathDelim + "experience/web/system-info"
)

type rawSystemNodeInfoLoopback struct {
	Ipv6Addr *string   `json:"ipv6_addr"`
	Ipv4Addr *string   `json:"ipv4_addr"`
	Id       *ObjectId `json:"id"`
}

type rawSystemNodeInfo struct {
	DomainId *string `json:"domain_id"`
	//UplinkedSystemIds interface{} `json:"uplinked_system_ids"`
	//DeployMode        interface{} `json:"deploy_mode"`
	Id             ObjectId                   `json:"id"`
	Label          string                     `json:"label"`
	PodId          ObjectId                   `json:"pod_id"`
	InterfaceMapId *ObjectId                  `json:"interface_map_id"`
	Loopback       *rawSystemNodeInfoLoopback `json:"loopback"`
	Hostname       string                     `json:"hostname"`
	//UnicastVtep        interface{}            `json:"unicast_vtep"`
	Role systemRole `json:"role"`
	//SystemId           interface{}   `json:"system_id"`
	//Hidden             bool          `json:"hidden"`
	//AnycastVtep        interface{}   `json:"anycast_vtep"`
	RackId *ObjectId `json:"rack_id"`
	//LogicalVtep        interface{}   `json:"logical_vtep"`
	//SuperspinePlaneId  interface{}   `json:"superspine_plane_id"`
	LogicalDeviceId ObjectId `json:"logical_device_id"`
	Tags            []string `json:"tags"`
	//HypervisorId       interface{} `json:"hypervisor_id"`
	//RedundancyProtocol interface{} `json:"redundancy_protocol"`
	External         bool `json:"external"`
	PortChannelIdMin int  `json:"port_channel_id_min"`
	//PositionData      interface{} `json:"position_data"`
	PortChannelIdMax  int                   `json:"port_channel_id_max"`
	RedundancyGroupId *ObjectId             `json:"redundancy_group_id"`
	GroupLabel        *string               `json:"group_label"`
	ManagementLevel   systemManagementLevel `json:"management_level"`
	DeviceProfileId   *ObjectId             `json:"device_profile_id"`
}

func (o rawSystemNodeInfo) polish() (*SystemNodeInfo, error) {
	var asn *uint32
	if o.DomainId != nil {
		i64, err := strconv.ParseInt(*o.DomainId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("while parsing ASN string %q - %w", *o.DomainId, err)
		}
		i32 := uint32(i64)
		asn = &i32
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
