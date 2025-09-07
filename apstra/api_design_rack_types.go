// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlDesignRackTypes       = apiUrlDesignPrefix + "rack-types"
	apiUrlDesignRackTypesPrefix = apiUrlDesignRackTypes + apiUrlPathDelim
	apiUrlDesignRackTypeById    = apiUrlDesignRackTypesPrefix + "%s"
)

type (
	AccessRedundancyProtocol int
	accessRedundancyProtocol string
)

const (
	AccessRedundancyProtocolNone = AccessRedundancyProtocol(iota)
	AccessRedundancyProtocolEsi
	AccessRedundancyProtocolUnknown = "unknown redundancy protocol '%s'"

	accessRedundancyProtocolNone    = accessRedundancyProtocol("")
	accessRedundancyProtocolEsi     = accessRedundancyProtocol("esi")
	accessRedundancyProtocolUnknown = "unknown redundancy protocol '%d'"
)

type (
	LeafRedundancyProtocol int
	leafRedundancyProtocol string
)

const (
	LeafRedundancyProtocolNone = LeafRedundancyProtocol(iota)
	LeafRedundancyProtocolMlag
	LeafRedundancyProtocolEsi
	LeafRedundancyProtocolUnknown = "unknown redundancy protocol '%s'"

	leafRedundancyProtocolNone    = leafRedundancyProtocol("")
	leafRedundancyProtocolMlag    = leafRedundancyProtocol("mlag")
	leafRedundancyProtocolEsi     = leafRedundancyProtocol("esi")
	leafRedundancyProtocolUnknown = "unknown type %d"
)

type (
	FeatureSwitch int
	featureSwitch string
)

const (
	FeatureSwitchDisabled = FeatureSwitch(iota)
	FeatureSwitchEnabled
	FeatureSwitchUnknown = "unknown feature switch state '%s'"

	featureSwitchDisabled = featureSwitch("disabled")
	featureSwitchEnabled  = featureSwitch("enabled")
	featureSwitchUnknown  = "unknown feature switch state '%d'"
)

type (
	SystemManagementLevel int
	systemManagementLevel string
)

const (
	SystemManagementLevelUnmanaged = SystemManagementLevel(iota)
	SystemManagementLevelTelemetryOnly
	SystemManagementLevelFullControl
	SystemManagementLevelNotInstalled
	SystemManagementLevelNone
	SystemManagementLevelUnknown = "unknown generic system management level '%s'"

	systemManagementLevelUnmanaged     = systemManagementLevel("unmanaged")
	systemManagementLevelTelemetryOnly = systemManagementLevel("telemetry_only")
	systemManagementLevelFullControl   = systemManagementLevel("full_control")
	systemManagementLevelNotInstalled  = systemManagementLevel("not_installed")
	systemManagementLevelNone          = systemManagementLevel("")
	systemManagementLevelUnknown       = "unknown generic system management level '%d'"
)

type (
	RackLinkAttachmentType int
	rackLinkAttachmentType string
)

const (
	RackLinkAttachmentTypeSingle = RackLinkAttachmentType(iota)
	RackLinkAttachmentTypeDual
	RackLinkAttachmentTypeUnknown = "unknown link attachment scheme '%s'"

	rackLinkAttachmentTypeSingle  = rackLinkAttachmentType("singleAttached")
	rackLinkAttachmentTypeDual    = rackLinkAttachmentType("dualAttached")
	rackLinkAttachmentTypeUnknown = "unknown link attachment scheme '%d'"
)

type (
	RackLinkLagMode int
	rackLinkLagMode string
)

const (
	RackLinkLagModeNone = RackLinkLagMode(iota)
	RackLinkLagModeActive
	RackLinkLagModePassive
	RackLinkLagModeStatic
	RackLinkLagModeUnknown = "unknown lag mode '%s'"

	rackLinkLagModeNone    = rackLinkLagMode("")
	rackLinkLagModeActive  = rackLinkLagMode("lacp_active")
	rackLinkLagModePassive = rackLinkLagMode("lacp_passive")
	rackLinkLagModeStatic  = rackLinkLagMode("static_lag")
	rackLinkLagModeUnknown = "unknown lag mode '%d'"
)

type (
	RackLinkSwitchPeer int
	rackLinkSwitchPeer string
)

const (
	RackLinkSwitchPeerNone = RackLinkSwitchPeer(iota)
	RackLinkSwitchPeerFirst
	RackLinkSwitchPeerSecond
	RackLinkSwitchPeerUnknown = "unknown switch peer '%s'"

	rackLinkSwitchPeerNone    = rackLinkSwitchPeer("")
	rackLinkSwitchPeerFirst   = rackLinkSwitchPeer("first")
	rackLinkSwitchPeerSecond  = rackLinkSwitchPeer("second")
	rackLinkSwitchPeerUnknown = "unknown switch peer '%d'"
)

func (o AccessRedundancyProtocol) Int() int {
	return int(o)
}

func (o AccessRedundancyProtocol) String() string {
	switch o {
	case AccessRedundancyProtocolNone:
		return string(accessRedundancyProtocolNone)
	case AccessRedundancyProtocolEsi:
		return string(accessRedundancyProtocolEsi)
	default:
		return fmt.Sprintf(accessRedundancyProtocolUnknown, o)
	}
}

func (o AccessRedundancyProtocol) raw() accessRedundancyProtocol {
	return accessRedundancyProtocol(o.String())
}

func (o accessRedundancyProtocol) string() string {
	return string(o)
}

func (o accessRedundancyProtocol) parse() (int, error) {
	switch o {
	case accessRedundancyProtocolNone:
		return int(AccessRedundancyProtocolNone), nil
	case accessRedundancyProtocolEsi:
		return int(AccessRedundancyProtocolEsi), nil
	default:
		return 0, fmt.Errorf(AccessRedundancyProtocolUnknown, o)
	}
}

func (o *AccessRedundancyProtocol) FromString(in string) error {
	i, err := accessRedundancyProtocol(in).parse()
	if err != nil {
		return err
	}
	*o = AccessRedundancyProtocol(i)
	return nil
}

func (o LeafRedundancyProtocol) Int() int {
	return int(o)
}

func (o LeafRedundancyProtocol) String() string {
	switch o {
	case LeafRedundancyProtocolNone:
		return string(leafRedundancyProtocolNone)
	case LeafRedundancyProtocolEsi:
		return string(leafRedundancyProtocolEsi)
	case LeafRedundancyProtocolMlag:
		return string(leafRedundancyProtocolMlag)
	default:
		return fmt.Sprintf(leafRedundancyProtocolUnknown, o)
	}
}

func (o LeafRedundancyProtocol) raw() leafRedundancyProtocol {
	return leafRedundancyProtocol(o.String())
}

func (o leafRedundancyProtocol) string() string {
	return string(o)
}

func (o leafRedundancyProtocol) parse() (int, error) {
	switch o {
	case leafRedundancyProtocolNone:
		return int(LeafRedundancyProtocolNone), nil
	case leafRedundancyProtocolEsi:
		return int(LeafRedundancyProtocolEsi), nil
	case leafRedundancyProtocolMlag:
		return int(LeafRedundancyProtocolMlag), nil
	default:
		return 0, fmt.Errorf(LeafRedundancyProtocolUnknown, o)
	}
}

func (o *LeafRedundancyProtocol) FromString(in string) error {
	i, err := leafRedundancyProtocol(in).parse()
	if err != nil {
		return err
	}
	*o = LeafRedundancyProtocol(i)
	return nil
}

func (o FeatureSwitch) Int() int {
	return int(o)
}

func (o FeatureSwitch) String() string {
	switch o {
	case FeatureSwitchDisabled:
		return string(featureSwitchDisabled)
	case FeatureSwitchEnabled:
		return string(featureSwitchEnabled)
	default:
		return fmt.Sprintf(featureSwitchUnknown, o)
	}
}

func (o featureSwitch) string() string {
	return string(o)
}

func (o featureSwitch) parse() (int, error) {
	switch o {
	case featureSwitchDisabled:
		return int(FeatureSwitchDisabled), nil
	case featureSwitchEnabled:
		return int(FeatureSwitchEnabled), nil
	default:
		return 0, fmt.Errorf(FeatureSwitchUnknown, o)
	}
}

func (o SystemManagementLevel) Int() int {
	return int(o)
}

func (o SystemManagementLevel) String() string {
	switch o {
	case SystemManagementLevelUnmanaged:
		return string(systemManagementLevelUnmanaged)
	case SystemManagementLevelTelemetryOnly:
		return string(systemManagementLevelTelemetryOnly)
	case SystemManagementLevelFullControl:
		return string(systemManagementLevelFullControl)
	case SystemManagementLevelNotInstalled:
		return string(systemManagementLevelNotInstalled)
	case SystemManagementLevelNone:
		return string(systemManagementLevelNone)
	default:
		return fmt.Sprintf(systemManagementLevelUnknown, o)
	}
}

func (o systemManagementLevel) string() string {
	return string(o)
}

func (o systemManagementLevel) parse() (int, error) {
	switch o {
	case systemManagementLevelUnmanaged:
		return int(SystemManagementLevelUnmanaged), nil
	case systemManagementLevelTelemetryOnly:
		return int(SystemManagementLevelTelemetryOnly), nil
	case systemManagementLevelFullControl:
		return int(SystemManagementLevelFullControl), nil
	case systemManagementLevelNotInstalled:
		return int(SystemManagementLevelNotInstalled), nil
	case systemManagementLevelNone:
		return int(SystemManagementLevelNone), nil
	default:
		return 0, fmt.Errorf(SystemManagementLevelUnknown, o)
	}
}

func (o RackLinkAttachmentType) Int() int {
	return int(o)
}

func (o RackLinkAttachmentType) String() string {
	switch o {
	case RackLinkAttachmentTypeSingle:
		return string(rackLinkAttachmentTypeSingle)
	case RackLinkAttachmentTypeDual:
		return string(rackLinkAttachmentTypeDual)
	default:
		return fmt.Sprintf(rackLinkAttachmentTypeUnknown, o)
	}
}

func (o rackLinkAttachmentType) string() string {
	return string(o)
}

func (o rackLinkAttachmentType) parse() (int, error) {
	switch o {
	case rackLinkAttachmentTypeSingle:
		return int(RackLinkAttachmentTypeSingle), nil
	case rackLinkAttachmentTypeDual:
		return int(RackLinkAttachmentTypeDual), nil
	default:
		return 0, fmt.Errorf(RackLinkAttachmentTypeUnknown, o)
	}
}

func (o RackLinkLagMode) Int() int {
	return int(o)
}

func (o RackLinkLagMode) String() string {
	switch o {
	case RackLinkLagModeNone:
		return string(rackLinkLagModeNone)
	case RackLinkLagModeActive:
		return string(rackLinkLagModeActive)
	case RackLinkLagModePassive:
		return string(rackLinkLagModePassive)
	case RackLinkLagModeStatic:
		return string(rackLinkLagModeStatic)
	default:
		return fmt.Sprintf(rackLinkLagModeUnknown, o)
	}
}

func (o *RackLinkLagMode) FromString(in string) error {
	i, err := rackLinkLagMode(in).parse()
	if err != nil {
		return err
	}
	*o = RackLinkLagMode(i)
	return nil
}

func (o rackLinkLagMode) string() string {
	return string(o)
}

func (o rackLinkLagMode) parse() (int, error) {
	switch o {
	case rackLinkLagModeNone:
		return int(RackLinkLagModeNone), nil
	case rackLinkLagModeActive:
		return int(RackLinkLagModeActive), nil
	case rackLinkLagModePassive:
		return int(RackLinkLagModePassive), nil
	case rackLinkLagModeStatic:
		return int(RackLinkLagModeStatic), nil
	default:
		return 0, fmt.Errorf(RackLinkLagModeUnknown, o)
	}
}

func (o RackLinkSwitchPeer) Int() int {
	return int(o)
}

func (o RackLinkSwitchPeer) String() string {
	switch o {
	case RackLinkSwitchPeerNone:
		return string(rackLinkSwitchPeerNone)
	case RackLinkSwitchPeerFirst:
		return string(rackLinkSwitchPeerFirst)
	case RackLinkSwitchPeerSecond:
		return string(rackLinkSwitchPeerSecond)
	default:
		return fmt.Sprintf(rackLinkSwitchPeerUnknown, o)
	}
}

func (o *RackLinkSwitchPeer) FromString(in string) error {
	i, err := rackLinkSwitchPeer(in).parse()
	if err != nil {
		return err
	}
	*o = RackLinkSwitchPeer(i)
	return nil
}

func (o rackLinkSwitchPeer) string() string {
	return string(o)
}

func (o rackLinkSwitchPeer) parse() (int, error) {
	switch o {
	case rackLinkSwitchPeerNone:
		return int(RackLinkSwitchPeerNone), nil
	case rackLinkSwitchPeerFirst:
		return int(RackLinkSwitchPeerFirst), nil
	case rackLinkSwitchPeerSecond:
		return int(RackLinkSwitchPeerSecond), nil
	default:
		return 0, fmt.Errorf(RackLinkSwitchPeerUnknown, o)
	}
}

type optionsRackTypeResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type LeafMlagInfo struct {
	LeafLeafL3LinkCount         int
	LeafLeafL3LinkPortChannelId int
	LeafLeafL3LinkSpeed         LogicalDevicePortSpeed
	LeafLeafLinkCount           int
	LeafLeafLinkPortChannelId   int
	LeafLeafLinkSpeed           LogicalDevicePortSpeed
	MlagVlanId                  int
}

type RackElementLeafSwitchRequest struct {
	Label              string
	MlagInfo           *LeafMlagInfo
	LinkPerSpineCount  int
	LinkPerSpineSpeed  LogicalDevicePortSpeed
	RedundancyProtocol LeafRedundancyProtocol
	Tags               []ObjectId
	LogicalDeviceId    ObjectId
}

func (o *RackElementLeafSwitchRequest) raw(tagMap map[ObjectId]DesignTagData) (*rawRackElementLeafSwitch, error) {
	tags := make([]string, len(o.Tags))
	for i, tagId := range o.Tags {
		if tagData, found := tagMap[tagId]; found {
			tags[i] = tagData.Label
		} else {
			return nil, fmt.Errorf("tagMap input to RackElementLeafSwitchRequest.raw() missing required tag ID: '%s'", tagId)
		}
	}

	result := &rawRackElementLeafSwitch{
		Label:              o.Label,
		LinkPerSpineCount:  o.LinkPerSpineCount,
		LinkPerSpineSpeed:  o.LinkPerSpineSpeed.raw(),
		RedundancyProtocol: o.RedundancyProtocol.raw(),
		LogicalDevice:      o.LogicalDeviceId,
		Tags:               tags,
	}
	if o.MlagInfo != nil {
		result.LeafLeafL3LinkCount = o.MlagInfo.LeafLeafL3LinkCount
		result.LeafLeafL3LinkPortChannelId = o.MlagInfo.LeafLeafL3LinkPortChannelId
		result.LeafLeafL3LinkSpeed = o.MlagInfo.LeafLeafL3LinkSpeed.raw()
		result.LeafLeafLinkCount = o.MlagInfo.LeafLeafLinkCount
		result.LeafLeafLinkPortChannelId = o.MlagInfo.LeafLeafLinkPortChannelId
		result.LeafLeafLinkSpeed = o.MlagInfo.LeafLeafLinkSpeed.raw()
		result.MlagVlanId = o.MlagInfo.MlagVlanId
	}
	return result, nil
}

type RackElementLeafSwitch struct {
	Label              string
	LinkPerSpineCount  int
	LinkPerSpineSpeed  LogicalDevicePortSpeed
	MlagInfo           *LeafMlagInfo
	RedundancyProtocol LeafRedundancyProtocol
	Tags               []DesignTagData
	LogicalDevice      *LogicalDeviceData
}

type rawRackElementLeafSwitch struct {
	Label                       string                     `json:"label"`
	LeafLeafL3LinkCount         int                        `json:"leaf_leaf_l3_link_count"`
	LeafLeafL3LinkPortChannelId int                        `json:"leaf_leaf_l3_link_port_channel_id"`
	LeafLeafL3LinkSpeed         *rawLogicalDevicePortSpeed `json:"leaf_leaf_l3_link_speed"`
	LeafLeafLinkCount           int                        `json:"leaf_leaf_link_count"`
	LeafLeafLinkPortChannelId   int                        `json:"leaf_leaf_link_port_channel_id"`
	LeafLeafLinkSpeed           *rawLogicalDevicePortSpeed `json:"leaf_leaf_link_speed"`
	LinkPerSpineCount           int                        `json:"link_per_spine_count"`
	LinkPerSpineSpeed           *rawLogicalDevicePortSpeed `json:"link_per_spine_speed"`
	LogicalDevice               ObjectId                   `json:"logical_device"`
	MlagVlanId                  int                        `json:"mlag_vlan_id"`
	RedundancyProtocol          leafRedundancyProtocol     `json:"redundancy_protocol,omitempty"`
	Tags                        []string                   `json:"tags"`
}

func (o *rawRackElementLeafSwitch) polish(rack *rawRackType) (*RackElementLeafSwitch, error) {
	rp, err := o.RedundancyProtocol.parse()
	if err != nil {
		return nil, err
	}

	var found bool

	var rld *rawLogicalDevice
	if rld, found = rack.logicalDeviceById(o.LogicalDevice); !found {
		return nil, fmt.Errorf("logical device '%s' not found in rack type '%s' definition", o.LogicalDevice, rack.Id)
	}
	pld, err := rld.polish()
	if err != nil {
		return nil, err
	}

	var tags []DesignTagData
	var tag *DesignTagData
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("design tag '%s' not found in rack type '%s' definition", label, rack.Id)
		}
		tags = append(tags, *tag)
	}

	var leafLeafL3LinkSpeed LogicalDevicePortSpeed
	if o.LeafLeafL3LinkSpeed != nil {
		leafLeafL3LinkSpeed = o.LeafLeafL3LinkSpeed.parse()
	}
	var leafLeafLinkSpeed LogicalDevicePortSpeed
	if o.LeafLeafLinkSpeed != nil {
		leafLeafLinkSpeed = o.LeafLeafLinkSpeed.parse()
	}
	var linkPerSpineSpeed LogicalDevicePortSpeed
	if o.LinkPerSpineSpeed != nil {
		linkPerSpineSpeed = o.LinkPerSpineSpeed.parse()
	}

	result := &RackElementLeafSwitch{
		Label:              o.Label,
		LinkPerSpineCount:  o.LinkPerSpineCount,
		LinkPerSpineSpeed:  linkPerSpineSpeed,
		RedundancyProtocol: LeafRedundancyProtocol(rp),
		Tags:               tags,
		LogicalDevice: &LogicalDeviceData{
			Panels:      pld.Data.Panels,
			DisplayName: pld.Data.DisplayName,
		},
		MlagInfo: &LeafMlagInfo{
			LeafLeafL3LinkCount:         o.LeafLeafL3LinkCount,
			LeafLeafL3LinkPortChannelId: o.LeafLeafL3LinkPortChannelId,
			LeafLeafL3LinkSpeed:         leafLeafL3LinkSpeed,
			LeafLeafLinkCount:           o.LeafLeafLinkCount,
			LeafLeafLinkPortChannelId:   o.LeafLeafLinkPortChannelId,
			LeafLeafLinkSpeed:           leafLeafLinkSpeed,
			MlagVlanId:                  o.MlagVlanId,
		},
	}

	return result, nil
}

type EsiLagInfo struct {
	AccessAccessLinkCount int
	AccessAccessLinkSpeed LogicalDevicePortSpeed
}

type RackElementAccessSwitchRequest struct {
	InstanceCount      int
	RedundancyProtocol AccessRedundancyProtocol
	Links              []RackLinkRequest
	Label              string
	LogicalDeviceId    ObjectId
	Tags               []ObjectId
	EsiLagInfo         *EsiLagInfo
}

func (o *RackElementAccessSwitchRequest) raw(tagMap map[ObjectId]DesignTagData) (*rawRackElementAccessSwitch, error) {
	tags := make([]string, len(o.Tags))
	for i, tagId := range o.Tags {
		if tagData, found := tagMap[tagId]; found {
			tags[i] = tagData.Label
		} else {
			return nil, fmt.Errorf("tagMap input to RackElementAccessSwitchRequest.raw() missing required tag ID: '%s'", tagId)
		}
	}

	links := make([]rawRackLink, len(o.Links))
	for i, l := range o.Links {
		rawLink, err := l.raw(tagMap)
		if err != nil {
			return nil, err
		}
		links[i] = *rawLink
	}

	var accessAccessLinkCount int
	var accessAccessLinkSpeed *rawLogicalDevicePortSpeed
	if o.EsiLagInfo != nil {
		accessAccessLinkCount = o.EsiLagInfo.AccessAccessLinkCount
		accessAccessLinkSpeed = o.EsiLagInfo.AccessAccessLinkSpeed.raw()
	}
	return &rawRackElementAccessSwitch{
		InstanceCount:         o.InstanceCount,
		RedundancyProtocol:    o.RedundancyProtocol.raw(),
		Links:                 links,
		Label:                 o.Label,
		AccessAccessLinkCount: accessAccessLinkCount,
		AccessAccessLinkSpeed: accessAccessLinkSpeed,
		LogicalDevice:         o.LogicalDeviceId,
		Tags:                  tags,
	}, nil
}

type RackElementAccessSwitch struct {
	InstanceCount      int
	RedundancyProtocol AccessRedundancyProtocol
	Links              []RackLink
	Label              string
	Tags               []DesignTagData
	LogicalDevice      *LogicalDeviceData
	EsiLagInfo         *EsiLagInfo
}

type rawRackElementAccessSwitch struct {
	InstanceCount         int                        `json:"instance_count"`
	RedundancyProtocol    accessRedundancyProtocol   `json:"redundancy_protocol,omitempty"`
	Links                 []rawRackLink              `json:"links"`
	Label                 string                     `json:"label"`
	LogicalDevice         ObjectId                   `json:"logical_device"`
	AccessAccessLinkCount int                        `json:"access_access_link_count"`
	AccessAccessLinkSpeed *rawLogicalDevicePortSpeed `json:"access_access_link_speed"`
	Tags                  []string                   `json:"tags"`
}

func (o *rawRackElementAccessSwitch) polish(rack *rawRackType) (*RackElementAccessSwitch, error) {
	rp, err := o.RedundancyProtocol.parse()
	if err != nil {
		return nil, err
	}

	var found bool

	var rld *rawLogicalDevice
	if rld, found = rack.logicalDeviceById(o.LogicalDevice); !found {
		return nil, fmt.Errorf("logical device '%s' not found in rack type '%s' definition", o.LogicalDevice, rack.Id)
	}
	pld, err := rld.polish()
	if err != nil {
		return nil, err
	}

	var tags []DesignTagData
	var tag *DesignTagData
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("design tag '%s' not found in rack type '%s' definition", label, rack.Id)
		}
		tags = append(tags, *tag)
	}

	var accessAccessLinkSpeed LogicalDevicePortSpeed
	if o.AccessAccessLinkSpeed != nil {
		accessAccessLinkSpeed = o.AccessAccessLinkSpeed.parse()
	}

	links := make([]RackLink, len(o.Links))
	for i, link := range o.Links {
		polished, err := link.polish(rack)
		if err != nil {
			return nil, err
		}
		links[i] = *polished
	}

	var esiLagInfo *EsiLagInfo
	if o.AccessAccessLinkCount > 0 {
		esiLagInfo = &EsiLagInfo{
			AccessAccessLinkCount: o.AccessAccessLinkCount,
			AccessAccessLinkSpeed: accessAccessLinkSpeed,
		}
	}

	return &RackElementAccessSwitch{
		InstanceCount:      o.InstanceCount,
		RedundancyProtocol: AccessRedundancyProtocol(rp),
		Links:              links,
		Label:              o.Label,
		EsiLagInfo:         esiLagInfo,
		Tags:               tags,
		LogicalDevice: &LogicalDeviceData{
			Panels:      pld.Data.Panels,
			DisplayName: pld.Data.DisplayName,
		},
	}, nil
}

type RackLinkRequest struct {
	Label              string                 // `json:"label"`
	Tags               []ObjectId             // `json:"tags"`
	LinkPerSwitchCount int                    // `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed // `json:"link_speed"`
	TargetSwitchLabel  string                 // `json:"target_switch_label"`
	AttachmentType     RackLinkAttachmentType // `json:"attachment_type"`
	LagMode            RackLinkLagMode        // `json:"lag_mode"`
	SwitchPeer         RackLinkSwitchPeer     // `json:"switch_peer"`
}

func (o RackLinkRequest) raw(tagMap map[ObjectId]DesignTagData) (*rawRackLink, error) {
	tags := make([]string, len(o.Tags))
	for i, tagId := range o.Tags {
		if tagData, found := tagMap[tagId]; found {
			tags[i] = tagData.Label
		} else {
			return nil, fmt.Errorf("tagMap input to RackLinkRequest.raw() missing required tag ID '%s'", tagId)
		}
	}

	// JSON encoding of lag_mode must be one of the accepted strings or null (nil ptr)
	var lagModePtr *rackLinkLagMode
	lagMode := rackLinkLagMode(o.LagMode.String())
	lagModePtr = &lagMode
	if lagMode == rackLinkLagModeNone {
		lagModePtr = nil
	}

	return &rawRackLink{
		Label:              o.Label,
		Tags:               tags,
		LinkPerSwitchCount: o.LinkPerSwitchCount,
		LinkSpeed:          o.LinkSpeed.raw(),
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     rackLinkAttachmentType(o.AttachmentType.String()),
		LagMode:            lagModePtr,
		SwitchPeer:         rackLinkSwitchPeer(o.SwitchPeer.String()),
	}, nil
}

type RackLink struct {
	Label              string
	LinkPerSwitchCount int
	LinkSpeed          LogicalDevicePortSpeed
	TargetSwitchLabel  string
	AttachmentType     RackLinkAttachmentType
	LagMode            RackLinkLagMode
	SwitchPeer         RackLinkSwitchPeer
	Tags               []DesignTagData
}

type rawRackLink struct {
	Label              string                     `json:"label"`
	LinkPerSwitchCount int                        `json:"link_per_switch_count"`
	LinkSpeed          *rawLogicalDevicePortSpeed `json:"link_speed"`
	TargetSwitchLabel  string                     `json:"target_switch_label"`
	AttachmentType     rackLinkAttachmentType     `json:"attachment_type"`
	LagMode            *rackLinkLagMode           `json:"lag_mode"` // do not "omitempty" // todo: explore this b/c the API sends 'null'
	SwitchPeer         rackLinkSwitchPeer         `json:"switch_peer,omitempty"`
	Tags               []string                   `json:"tags"`
}

func (o rawRackLink) polish(rack *rawRackType) (*RackLink, error) {
	attachment, err := o.AttachmentType.parse()
	if err != nil {
		return nil, err
	}

	var lagMode int
	if o.LagMode != nil {
		lagMode, err = o.LagMode.parse()
		if err != nil {
			return nil, err
		}
	}

	switchPeer, err := o.SwitchPeer.parse()
	if err != nil {
		return nil, err
	}

	var found bool
	var tags []DesignTagData
	var tag *DesignTagData
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("link '%s' in rack '%s' has tag '%s' but tag missing from rack definition", o.Label, rack.Id, label)
		}
		tags = append(tags, *tag)
	}

	var linkSpeed LogicalDevicePortSpeed
	if o.LinkSpeed != nil {
		linkSpeed = o.LinkSpeed.parse()
	}

	return &RackLink{
		Label:              o.Label,
		Tags:               tags,
		LinkPerSwitchCount: o.LinkPerSwitchCount,
		LinkSpeed:          linkSpeed,
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     RackLinkAttachmentType(attachment),
		LagMode:            RackLinkLagMode(lagMode),
		SwitchPeer:         RackLinkSwitchPeer(switchPeer),
	}, nil
}

type RackElementGenericSystemRequest struct {
	Count            int
	AsnDomain        FeatureSwitch
	ManagementLevel  SystemManagementLevel
	PortChannelIdMin int
	PortChannelIdMax int
	Loopback         FeatureSwitch
	Tags             []ObjectId
	Label            string
	Links            []RackLinkRequest
	LogicalDeviceId  ObjectId
}

func (o *RackElementGenericSystemRequest) raw(tagMap map[ObjectId]DesignTagData) (*rawRackElementGenericSystem, error) {
	tags := make([]string, len(o.Tags))
	for i, tagId := range o.Tags {
		if tagData, found := tagMap[tagId]; found {
			tags[i] = tagData.Label
		} else {
			return nil, fmt.Errorf("tagMap input to RackElementGenericSystemRequest.raw() missing required tag ID: '%s'", tagId)
		}
	}

	links := make([]rawRackLink, len(o.Links))
	for i, l := range o.Links {
		rawLink, err := l.raw(tagMap)
		if err != nil {
			return nil, err
		}
		links[i] = *rawLink
	}

	return &rawRackElementGenericSystem{
		Count:            o.Count,
		AsnDomain:        featureSwitch(o.AsnDomain.String()),
		ManagementLevel:  systemManagementLevel(o.ManagementLevel.String()),
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Loopback:         featureSwitch(o.Loopback.String()),
		Label:            o.Label,
		Links:            links,
		LogicalDevice:    o.LogicalDeviceId,
		Tags:             tags,
	}, nil
}

type RackElementGenericSystem struct {
	Count            int
	AsnDomain        FeatureSwitch
	ManagementLevel  SystemManagementLevel
	PortChannelIdMin int
	PortChannelIdMax int
	Loopback         FeatureSwitch
	Tags             []DesignTagData
	Label            string
	Links            []RackLink
	LogicalDevice    *LogicalDeviceData
}

type rawRackElementGenericSystem struct {
	Count            int                   `json:"count"`
	AsnDomain        featureSwitch         `json:"asn_domain"`
	ManagementLevel  systemManagementLevel `json:"management_level"`
	PortChannelIdMin int                   `json:"port_channel_id_min"`
	PortChannelIdMax int                   `json:"port_channel_id_max"`
	Loopback         featureSwitch         `json:"loopback"`
	Tags             []string              `json:"tags"`
	Label            string                `json:"label"`
	LogicalDevice    ObjectId              `json:"logical_device"`
	Links            []rawRackLink         `json:"links"`
}

func (o *rawRackElementGenericSystem) polish(rack *rawRackType) (*RackElementGenericSystem, error) {
	asnDomain, err := o.AsnDomain.parse()
	if err != nil {
		return nil, err
	}

	mgmtLevel, err := o.ManagementLevel.parse()
	if err != nil {
		return nil, err
	}

	loopback, err := o.Loopback.parse()
	if err != nil {
		return nil, err
	}

	var links []RackLink
	for _, link := range o.Links {
		p, err := link.polish(rack)
		if err != nil {
			return nil, err
		}
		links = append(links, *p)
	}

	var found bool

	var rld *rawLogicalDevice
	if rld, found = rack.logicalDeviceById(o.LogicalDevice); !found {
		return nil, fmt.Errorf("logical device '%s' not found in rack type '%s' definition", o.LogicalDevice, rack.Id)
	}
	pld, err := rld.polish()
	if err != nil {
		return nil, err
	}

	var tags []DesignTagData
	var tag *DesignTagData
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("design tag '%s' not found in rack type '%s' definition", label, rack.Id)
		}
		tags = append(tags, *tag)
	}

	return &RackElementGenericSystem{
		Count:            o.Count,
		AsnDomain:        FeatureSwitch(asnDomain),
		ManagementLevel:  SystemManagementLevel(mgmtLevel),
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Loopback:         FeatureSwitch(loopback),
		Tags:             tags,
		Label:            o.Label,
		Links:            links,
		LogicalDevice: &LogicalDeviceData{
			Panels:      pld.Data.Panels,
			DisplayName: pld.Data.DisplayName,
		},
	}, nil
}

type RackTypeRequest struct {
	DisplayName              string
	Description              string
	FabricConnectivityDesign enum.FabricConnectivityDesign
	LeafSwitches             []RackElementLeafSwitchRequest
	AccessSwitches           []RackElementAccessSwitchRequest
	GenericSystems           []RackElementGenericSystemRequest
}

func (o *RackTypeRequest) raw(ctx context.Context, client *Client) (*rawRackTypeRequest, error) {
	result := &rawRackTypeRequest{
		DisplayName:              o.DisplayName,
		Description:              o.Description,
		FabricConnectivityDesign: o.FabricConnectivityDesign,
		LogicalDevices:           nil, // populated based on ldMap below
		Tags:                     nil, // populated based on tagMap below
		LeafSwitches:             make([]rawRackElementLeafSwitch, len(o.LeafSwitches)),
		AccessSwitches:           make([]rawRackElementAccessSwitch, len(o.AccessSwitches)),
		GenericSystems:           make([]rawRackElementGenericSystem, len(o.GenericSystems)),
	}

	// collect IDs of all logical devices relevant to this rack as a "set" of Object IDs
	ldMap := make(map[ObjectId]struct{})

	// collect all DesignTagData objects relevant to this rack, keyed by tag ID
	tagMap := make(map[ObjectId]DesignTagData)

	// getLabelAndCacheTagById populates tagMap (map[tagId]DesignTagData) by
	// calling the Apstra tag API each time it's called with a previously-unseen
	// tag ID. It returns the tag's label (used as a key within the rack-type
	// JSON) and squirrels away the tag payload (DesignTagData) for subsequent
	// use in the rawRackTypeRequest.
	getLabelAndCacheTagById := func(id ObjectId) (string, error) {
		var tagData DesignTagData
		var found bool
		if tagData, found = tagMap[id]; !found {
			tag, err := client.GetTag(ctx, id)
			if err != nil {
				return "", err
			}
			tagMap[id] = *tag.Data
			return tag.Data.Label, nil
		}
		return tagData.Label, nil
	}

	// each leaf switch: logical device and tags
	for i, req := range o.LeafSwitches {
		tagLabels := make([]string, len(req.Tags))
		for j, tagId := range req.Tags {
			label, err := getLabelAndCacheTagById(tagId) // fetch the tag label / cache the payload
			if err != nil {
				return nil, err
			}
			tagLabels[j] = label
		}

		raw, err := req.raw(tagMap) // raw-ify the leaf switch request
		if err != nil {
			return nil, err
		}

		ldMap[req.LogicalDeviceId] = struct{}{} // Add the logical device to our set
		result.LeafSwitches[i] = *raw           // Add the raw leaf switch request to the result
	}

	// each access switch: logical device, tags and link tags
	for i, req := range o.AccessSwitches {
		tagLabels := make([]string, len(req.Tags))
		for j, tagId := range req.Tags {
			label, err := getLabelAndCacheTagById(tagId) // fetch the tag label / cache the payload
			if err != nil {
				return nil, err
			}
			tagLabels[j] = label
		}

		// populate map with tags used by each access switch link
		for _, linkReq := range req.Links {
			for _, tagId := range linkReq.Tags {
				_, err := getLabelAndCacheTagById(tagId) // fetch the tag label / cache the payload
				if err != nil {
					return nil, err
				}
			}
		}

		raw, err := req.raw(tagMap) // raw-ify the access switch request
		if err != nil {
			return nil, err
		}

		ldMap[req.LogicalDeviceId] = struct{}{} // Add the logical device to our set
		result.AccessSwitches[i] = *raw         // Add the raw access switch request to the result
	}

	// each generic system: logical device, tags and link tags
	for i, req := range o.GenericSystems {
		tagLabels := make([]string, len(req.Tags))
		for j, tagId := range req.Tags {
			label, err := getLabelAndCacheTagById(tagId) // fetch the tag label / cache the payload
			if err != nil {
				return nil, err
			}
			tagLabels[j] = label
		}

		// populate map with tags used by each generic system link
		for _, linkReq := range req.Links {
			for _, tagId := range linkReq.Tags {
				_, err := getLabelAndCacheTagById(tagId) // fetch the tag label / cache the payload
				if err != nil {
					return nil, err
				}
			}
		}

		raw, err := req.raw(tagMap) // raw-ify the generic system request
		if err != nil {
			return nil, err
		}

		ldMap[req.LogicalDeviceId] = struct{}{} // Add the logical device to our set
		result.GenericSystems[i] = *raw         // Add the raw generic system request to the result
	}

	// prepare the []rawLogicalDevice we'll submit when creating the rack type
	// using ldMap, which is the set of logical device IDs representing every
	// device in the rack (leaf, access, generic)
	result.LogicalDevices = make([]rawLogicalDevice, len(ldMap))
	i := 0
	for id := range ldMap {
		ld, err := client.getLogicalDevice(ctx, id)
		if err != nil {
			return nil, err
		}

		result.LogicalDevices[i] = rawLogicalDevice{
			Id:          ld.Id,
			DisplayName: ld.DisplayName,
			Panels:      ld.Panels,
		}
		i++
	}

	// prepare the []DesignTagData we'll submit when creating the rack type
	// using tagMap, which is the set of DesignTagData representing every tag
	// applied to every device (leaf, access, generic) and link (from access or
	// generic) found in the rack
	result.Tags = make([]DesignTagData, len(tagMap))
	i = 0
	for _, tagData := range tagMap {
		result.Tags[i] = tagData
		i++
	}

	return result, nil
}

type rawRackTypeRequest struct {
	DisplayName              string                        `json:"display_name"`
	Description              string                        `json:"description"`
	FabricConnectivityDesign enum.FabricConnectivityDesign `json:"fabric_connectivity_design"`
	Tags                     []DesignTagData               `json:"tags,omitempty"`
	LogicalDevices           []rawLogicalDevice            `json:"logical_devices,omitempty"`
	GenericSystems           []rawRackElementGenericSystem `json:"generic_systems,omitempty"`
	LeafSwitches             []rawRackElementLeafSwitch    `json:"leafs,omitempty"`
	AccessSwitches           []rawRackElementAccessSwitch  `json:"access_switches,omitempty"`
}

type RackType struct {
	Id             ObjectId
	CreatedAt      *time.Time
	LastModifiedAt *time.Time
	Data           *RackTypeData
}

type RackTypeData struct {
	DisplayName              string
	Description              string
	FabricConnectivityDesign enum.FabricConnectivityDesign
	LeafSwitches             []RackElementLeafSwitch
	GenericSystems           []RackElementGenericSystem
	AccessSwitches           []RackElementAccessSwitch
}

type rawRackType struct {
	Id                       ObjectId                      `json:"id,omitempty"`
	DisplayName              string                        `json:"display_name"`
	Description              string                        `json:"description"`
	FabricConnectivityDesign enum.FabricConnectivityDesign `json:"fabric_connectivity_design"`
	Tags                     []DesignTagData               `json:"tags,omitempty"`
	CreatedAt                *time.Time                    `json:"created_at,omitempty"`
	LastModifiedAt           *time.Time                    `json:"last_modified_at,omitempty"`
	LogicalDevices           []rawLogicalDevice            `json:"logical_devices,omitempty"`
	GenericSystems           []rawRackElementGenericSystem `json:"generic_systems,omitempty"`
	LeafSwitches             []rawRackElementLeafSwitch    `json:"leafs,omitempty"`
	AccessSwitches           []rawRackElementAccessSwitch  `json:"access_switches,omitempty"`
}

func (o *rawRackType) polish() (*RackType, error) {
	result := &RackType{
		Id:             o.Id,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &RackTypeData{
			DisplayName:              o.DisplayName,
			Description:              o.Description,
			FabricConnectivityDesign: o.FabricConnectivityDesign,
			LeafSwitches:             make([]RackElementLeafSwitch, len(o.LeafSwitches)),
			AccessSwitches:           make([]RackElementAccessSwitch, len(o.AccessSwitches)),
			GenericSystems:           make([]RackElementGenericSystem, len(o.GenericSystems)),
		},
	}

	for i, raw := range o.LeafSwitches {
		polished, err := raw.polish(o)
		if err != nil {
			return nil, err
		}
		result.Data.LeafSwitches[i] = *polished
	}

	for i, raw := range o.AccessSwitches {
		polished, err := raw.polish(o)
		if err != nil {
			return nil, err
		}
		result.Data.AccessSwitches[i] = *polished
	}

	for i, raw := range o.GenericSystems {
		polished, err := raw.polish(o)
		if err != nil {
			return nil, err
		}
		result.Data.GenericSystems[i] = *polished
	}

	return result, nil
}

func (o *rawRackType) logicalDeviceById(desired ObjectId) (*rawLogicalDevice, bool) {
	for _, ld := range o.LogicalDevices {
		if ld.Id == desired {
			return &ld, true
		}
	}
	return nil, false
}

func (o *rawRackType) tagByLabel(desired string) (*DesignTagData, bool) {
	for _, tag := range o.Tags {
		if tag.Label == desired {
			return &tag, true
		}
	}
	return nil, false
}

func (o *Client) listRackTypeIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsRackTypeResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignRackTypes,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getRackType(ctx context.Context, id ObjectId) (*rawRackType, error) {
	response := &rawRackType{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignRackTypeById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *Client) getAllRackTypes(ctx context.Context) ([]RackType, error) {
	response := &struct {
		Items []rawRackType `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignRackTypes,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]RackType, len(response.Items))
	for i, raw := range response.Items {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

func (o *Client) getRackTypeByName(ctx context.Context, name string) (*RackType, error) {
	rackTypes, err := o.getAllRackTypes(ctx)
	if err != nil {
		return nil, err
	}
	for _, rackType := range rackTypes {
		if rackType.Data.DisplayName == name {
			return &rackType, nil
		}
	}
	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("rack type with name '%s' not found", name),
	}
}

func (o *Client) createRackType(ctx context.Context, request *rawRackTypeRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignRackTypes,
		apiInput:    request,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateRackType(ctx context.Context, id ObjectId, request *RackTypeRequest) error {
	rawRequest, err := request.raw(ctx, o)
	if err != nil {
		return err
	}

	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignRackTypeById, id),
		apiInput: rawRequest,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteRackType(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignRackTypeById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
