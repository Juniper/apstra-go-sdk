package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDesignRackTypes       = apiUrlDesignPrefix + "rack-types"
	apiUrlDesignRackTypesPrefix = apiUrlDesignRackTypes + apiUrlPathDelim
	apiUrlDesignRackTypeById    = apiUrlDesignRackTypesPrefix + "%s"
)

type AccessRedundancyProtocol int
type accessRedundancyProtocol string

const (
	AccessRedundancyProtocolNone = AccessRedundancyProtocol(iota)
	AccessRedundancyProtocolEsi
	AccessRedundancyProtocolUnknown = "unknown redundancy protocol '%s'"

	accessRedundancyProtocolNone    = accessRedundancyProtocol("")
	accessRedundancyProtocolEsi     = accessRedundancyProtocol("esi")
	accessRedundancyProtocolUnknown = "unknown redundancy protocol '%d'"
)

type LeafRedundancyProtocol int
type leafRedundancyProtocol string

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

type FabricConnectivityDesign int
type fabricConnectivityDesign string

const (
	FabricConnectivityDesignL3Clos = FabricConnectivityDesign(iota)
	FabricConnectivityDesignL3Collapsed
	FabricConnectivityDesignUnknown = "unknown connectivity design '%s'"

	fabricConnectivityDesignL3Clos      = fabricConnectivityDesign("l3clos")
	fabricConnectivityDesignL3Collapsed = fabricConnectivityDesign("l3collapsed")
	fabricConnectivityDesignUnknown     = "unknown connectivity design '%d'"
)

type FeatureSwitch int
type featureSwitch string

const (
	FeatureSwitchDisabled = FeatureSwitch(iota)
	FeatureSwitchEnabled
	FeatureSwitchUnknown = "unknown feature switch state '%s'"

	featureSwitchDisabled = featureSwitch("disabled")
	featureSwitchEnabled  = featureSwitch("enabled")
	featureSwitchUnknown  = "unknown feature switch state '%d'"
)

type GenericSystemManagementLevel int
type genericSystemManagementLevel string

const (
	GenericSystemUnmanaged = GenericSystemManagementLevel(iota)
	GenericSystemTelemetryOnly
	GenericSystemFullControl
	GenericSystemUnknown = "unknown generic system management level '%s'"

	genericSystemUnmanaged     = genericSystemManagementLevel("unmanaged")
	genericSystemTelemetryOnly = genericSystemManagementLevel("telemetry_only")
	genericSystemFullControl   = genericSystemManagementLevel("full_control")
	genericSystemUnknown       = "unknown generic system management level '%d'"
)

type RackLinkAttachmentType int
type rackLinkAttachmentType string

const (
	RackLinkAttachmentTypeSingle = RackLinkAttachmentType(iota)
	RackLinkAttachmentTypeDual
	RackLinkAttachmentTypeUnknown = "unknown link attachment scheme '%s'"

	rackLinkAttachmentTypeSingle  = rackLinkAttachmentType("singleAttached")
	rackLinkAttachmentTypeDual    = rackLinkAttachmentType("dualAttached")
	rackLinkAttachmentTypeUnknown = "unknown link attachment scheme '%d'"
)

type RackLinkLagMode int
type rackLinkLagMode string

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

type RackLinkSwitchPeer int
type rackLinkSwitchPeer string

const (
	RackLinkSwitchPeerNone = RackLinkSwitchPeer(iota)
	RackLinkSwitchPeerFirst
	RackLinkSwitchPeerSecond
	RackLinkSwitchPeerUnknown = "unknown lag mode '%s'"

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

func (o FabricConnectivityDesign) Int() int {
	return int(o)
}

func (o FabricConnectivityDesign) String() string {
	switch o {
	case FabricConnectivityDesignL3Clos:
		return string(fabricConnectivityDesignL3Clos)
	case FabricConnectivityDesignL3Collapsed:
		return string(fabricConnectivityDesignL3Collapsed)
	default:
		return fmt.Sprintf(fabricConnectivityDesignUnknown, o)
	}
}

func (o FabricConnectivityDesign) raw() fabricConnectivityDesign {
	return fabricConnectivityDesign(o.String())
}

func (o fabricConnectivityDesign) string() string {
	return string(o)
}

func (o fabricConnectivityDesign) parse() (int, error) {
	switch o {
	case fabricConnectivityDesignL3Clos:
		return int(FabricConnectivityDesignL3Clos), nil
	case fabricConnectivityDesignL3Collapsed:
		return int(FabricConnectivityDesignL3Collapsed), nil
	default:
		return 0, fmt.Errorf(FabricConnectivityDesignUnknown, o)
	}
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

func (o GenericSystemManagementLevel) Int() int {
	return int(o)
}

func (o GenericSystemManagementLevel) String() string {
	switch o {
	case GenericSystemUnmanaged:
		return string(genericSystemUnmanaged)
	case GenericSystemTelemetryOnly:
		return string(genericSystemTelemetryOnly)
	case GenericSystemFullControl:
		return string(genericSystemFullControl)
	default:
		return fmt.Sprintf(genericSystemUnknown, o)
	}
}

func (o genericSystemManagementLevel) string() string {
	return string(o)
}
func (o genericSystemManagementLevel) parse() (int, error) {
	switch o {
	case genericSystemFullControl:
		return int(GenericSystemFullControl), nil
	case genericSystemUnmanaged:
		return int(GenericSystemUnmanaged), nil
	case genericSystemTelemetryOnly:
		return int(GenericSystemTelemetryOnly), nil
	default:
		return 0, fmt.Errorf(GenericSystemUnknown, o)
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

type RackElementLeafSwitchRequest struct {
	Label                       string
	LeafLeafL3LinkCount         int
	LeafLeafL3LinkPortChannelId int
	LeafLeafL3LinkSpeed         LogicalDevicePortSpeed
	LeafLeafLinkCount           int
	LeafLeafLinkPortChannelId   int
	LeafLeafLinkSpeed           LogicalDevicePortSpeed
	LinkPerSpineCount           int
	LinkPerSpineSpeed           LogicalDevicePortSpeed
	MlagVlanId                  int
	RedundancyProtocol          LeafRedundancyProtocol
	Tags                        []TagLabel
	LogicalDeviceId             ObjectId
}

func (o *RackElementLeafSwitchRequest) raw() *rawRackElementLeafSwitchRequest {
	return &rawRackElementLeafSwitchRequest{
		Label:                       o.Label,
		LeafLeafL3LinkCount:         o.LeafLeafL3LinkCount,
		LeafLeafL3LinkPortChannelId: o.LeafLeafL3LinkPortChannelId,
		LeafLeafL3LinkSpeed:         o.LeafLeafL3LinkSpeed.raw(),
		LeafLeafLinkCount:           o.LeafLeafLinkCount,
		LeafLeafLinkPortChannelId:   o.LeafLeafLinkPortChannelId,
		LeafLeafLinkSpeed:           o.LeafLeafLinkSpeed.raw(),
		LinkPerSpineCount:           o.LinkPerSpineCount,
		LinkPerSpineSpeed:           o.LinkPerSpineSpeed.raw(),
		MlagVlanId:                  o.MlagVlanId,
		RedundancyProtocol:          o.RedundancyProtocol.raw(),
		LogicalDevice:               o.LogicalDeviceId, // needs to be fetched from API, cloned into rack type on create() / update()
		Tags:                        o.Tags,            // needs to be fetched from API, cloned into rack type on create() / update()
	}
}

type rawRackElementLeafSwitchRequest struct {
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
	Tags                        []TagLabel                 `json:"tags,omitempty"`
}

type RackElementLeafSwitch struct {
	Label                       string
	LeafLeafL3LinkCount         int
	LeafLeafL3LinkPortChannelId int
	LeafLeafL3LinkSpeed         LogicalDevicePortSpeed
	LeafLeafLinkCount           int
	LeafLeafLinkPortChannelId   int
	LeafLeafLinkSpeed           LogicalDevicePortSpeed
	LinkPerSpineCount           int
	LinkPerSpineSpeed           LogicalDevicePortSpeed
	MlagVlanId                  int
	RedundancyProtocol          LeafRedundancyProtocol
	Tags                        []DesignTag
	Panels                      []LogicalDevicePanel
	DisplayName                 string
	LogicalDeviceId             ObjectId
}

type rawRackElementLeafSwitch struct {
	Label                       string                    `json:"label"`
	LeafLeafL3LinkCount         int                       `json:"leaf_leaf_l3_link_count"`
	LeafLeafL3LinkPortChannelId int                       `json:"leaf_leaf_l3_link_port_channel_id"`
	LeafLeafL3LinkSpeed         rawLogicalDevicePortSpeed `json:"leaf_leaf_l3_link_speed"`
	LeafLeafLinkCount           int                       `json:"leaf_leaf_link_count"`
	LeafLeafLinkPortChannelId   int                       `json:"leaf_leaf_link_port_channel_id"`
	LeafLeafLinkSpeed           rawLogicalDevicePortSpeed `json:"leaf_leaf_link_speed"`
	LinkPerSpineCount           int                       `json:"link_per_spine_count"`
	LinkPerSpineSpeed           rawLogicalDevicePortSpeed `json:"link_per_spine_speed"`
	LogicalDevice               ObjectId                  `json:"logical_device"`
	MlagVlanId                  int                       `json:"mlag_vlan_id"`
	RedundancyProtocol          leafRedundancyProtocol    `json:"redundancy_protocol,omitempty"`
	Tags                        []TagLabel                `json:"tags"`
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

	var tags []DesignTag
	var tag *DesignTag
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("design tag '%s' not found in rack type '%s' definition", label, rack.Id)
		}
		tags = append(tags, *tag)
	}

	result := &RackElementLeafSwitch{
		Label:                       o.Label,
		LeafLeafL3LinkCount:         o.LeafLeafL3LinkCount,
		LeafLeafL3LinkPortChannelId: o.LeafLeafL3LinkPortChannelId,
		LeafLeafL3LinkSpeed:         o.LeafLeafL3LinkSpeed.parse(),
		LeafLeafLinkCount:           o.LeafLeafLinkCount,
		LeafLeafLinkPortChannelId:   o.LeafLeafLinkPortChannelId,
		LeafLeafLinkSpeed:           o.LeafLeafLinkSpeed.parse(),
		LinkPerSpineCount:           o.LinkPerSpineCount,
		LinkPerSpineSpeed:           o.LinkPerSpineSpeed.parse(),
		MlagVlanId:                  o.MlagVlanId,
		RedundancyProtocol:          LeafRedundancyProtocol(rp),
		Panels:                      pld.Panels,
		DisplayName:                 pld.DisplayName,
		Tags:                        tags,
	}

	return result, nil
}

type RackElementAccessSwitchRequest struct {
	InstanceCount         int
	RedundancyProtocol    AccessRedundancyProtocol
	Links                 []RackLinkRequest
	Label                 string
	LogicalDeviceId       ObjectId
	Tags                  []TagLabel
	AccessAccessLinkCount int
	AccessAccessLinkSpeed LogicalDevicePortSpeed
}

func (o *RackElementAccessSwitchRequest) raw() *rawRackElementAccessSwitchRequest {
	links := make([]rawRackLinkRequest, len(o.Links))
	for i, l := range o.Links {
		links[i] = *l.raw()
	}
	return &rawRackElementAccessSwitchRequest{
		InstanceCount:         o.InstanceCount,
		RedundancyProtocol:    o.RedundancyProtocol.raw(),
		Links:                 links,
		Label:                 o.Label,
		AccessAccessLinkCount: o.AccessAccessLinkCount,
		AccessAccessLinkSpeed: o.AccessAccessLinkSpeed.raw(),
		LogicalDevice:         o.LogicalDeviceId, // needs to be fetched from API, cloned into rack type on create() / update()
		Tags:                  o.Tags,            // needs to be fetched from API, cloned into rack type on create() / update()
	}
}

type rawRackElementAccessSwitchRequest struct {
	InstanceCount         int                        `json:"instance_count"`
	RedundancyProtocol    accessRedundancyProtocol   `json:"redundancy_protocol,omitempty"`
	Links                 []rawRackLinkRequest       `json:"links"`
	Label                 string                     `json:"label"`
	LogicalDevice         ObjectId                   `json:"logical_device"`
	AccessAccessLinkCount int                        `json:"access_access_link_count"`
	AccessAccessLinkSpeed *rawLogicalDevicePortSpeed `json:"access_access_link_speed"`
	Tags                  []TagLabel                 `json:"tags,omitempty"`
}

type RackElementAccessSwitch struct {
	InstanceCount         int
	RedundancyProtocol    AccessRedundancyProtocol
	Links                 []RackLink
	Label                 string
	Panels                []LogicalDevicePanel
	DisplayName           string
	LogicalDeviceId       ObjectId
	Tags                  []DesignTag
	AccessAccessLinkCount int
	AccessAccessLinkSpeed LogicalDevicePortSpeed
}

type rawRackElementAccessSwitch struct {
	InstanceCount         int                       `json:"instance_count"`
	RedundancyProtocol    accessRedundancyProtocol  `json:"redundancy_protocol,omitempty"`
	Links                 []RackLink                `json:"links"`
	Label                 string                    `json:"label"`
	LogicalDevice         ObjectId                  `json:"logical_device"`
	AccessAccessLinkCount int                       `json:"access_access_link_count"`
	AccessAccessLinkSpeed rawLogicalDevicePortSpeed `json:"access_access_link_speed"`
	Tags                  []TagLabel                `json:"tags"`
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

	var tags []DesignTag
	var tag *DesignTag
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("design tag '%s' not found in rack type '%s' definition", label, rack.Id)
		}
		tags = append(tags, *tag)
	}

	return &RackElementAccessSwitch{
		InstanceCount:         o.InstanceCount,
		RedundancyProtocol:    AccessRedundancyProtocol(rp),
		Links:                 o.Links,
		Label:                 o.Label,
		Panels:                pld.Panels,
		DisplayName:           pld.DisplayName,
		AccessAccessLinkCount: o.AccessAccessLinkCount,
		AccessAccessLinkSpeed: o.AccessAccessLinkSpeed.parse(),
		Tags:                  tags,
	}, nil
}

type RackLinkRequest struct {
	Label              string                 // `json:"label"`
	Tags               []TagLabel             // `json:"tags"`
	LinkPerSwitchCount int                    // `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed // `json:"link_speed"`
	TargetSwitchLabel  string                 // `json:"target_switch_label"`
	AttachmentType     RackLinkAttachmentType // `json:"attachment_type"`
	LagMode            RackLinkLagMode        // `json:"lag_mode"`
	SwitchPeer         RackLinkSwitchPeer     // `json:"switch_peer"`
}

func (o RackLinkRequest) raw() *rawRackLinkRequest {
	tags := make([]TagLabel, len(o.Tags))
	for i, tag := range o.Tags {
		tags[i] = tag
	}

	// JSON encoding of lag_mode must be one of the accepted strings or null (nil ptr)
	var lagModePtr *rackLinkLagMode
	lagMode := rackLinkLagMode(o.LagMode.String())
	lagModePtr = &lagMode
	if lagMode == rackLinkLagModeNone {
		lagModePtr = nil
	}

	return &rawRackLinkRequest{
		Label:              o.Label,
		Tags:               tags,
		LinkPerSwitchCount: o.LinkPerSwitchCount,
		LinkSpeed:          o.LinkSpeed.raw(),
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     rackLinkAttachmentType(o.AttachmentType.String()),
		LagMode:            lagModePtr,
		SwitchPeer:         rackLinkSwitchPeer(o.SwitchPeer.String()),
	}
}

type rawRackLinkRequest struct {
	Label              string                     `json:"label"`
	LinkPerSwitchCount int                        `json:"link_per_switch_count"`
	LinkSpeed          *rawLogicalDevicePortSpeed `json:"link_speed"`
	TargetSwitchLabel  string                     `json:"target_switch_label"`
	AttachmentType     rackLinkAttachmentType     `json:"attachment_type"`
	LagMode            *rackLinkLagMode           `json:"lag_mode"` // do not "omitempty" // todo: explore this b/c the API sends 'null'
	SwitchPeer         rackLinkSwitchPeer         `json:"switch_peer,omitempty"`
	Tags               []TagLabel                 `json:"tags"` // needs to be fetched from API, cloned into rack type on create() / update()
}

type RackLink struct {
	Label              string                 // `json:"label"`
	LinkPerSwitchCount int                    // `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed // `json:"link_speed"`
	TargetSwitchLabel  string                 // `json:"target_switch_label"`
	AttachmentType     RackLinkAttachmentType // `json:"attachment_type"`
	LagMode            RackLinkLagMode        // `json:"lag_mode"`
	SwitchPeer         RackLinkSwitchPeer     // `json:"switch_peer"`
	Tags               []DesignTag            // `json:"tags"`
}

type rawRackLink struct {
	Label              string                     `json:"label"`
	LinkPerSwitchCount int                        `json:"link_per_switch_count"`
	LinkSpeed          *rawLogicalDevicePortSpeed `json:"link_speed"`
	TargetSwitchLabel  string                     `json:"target_switch_label"`
	AttachmentType     rackLinkAttachmentType     `json:"attachment_type"`
	LagMode            *rackLinkLagMode           `json:"lag_mode"` // do not "omitempty" // todo: explore this b/c the API sends 'null'
	SwitchPeer         rackLinkSwitchPeer         `json:"switch_peer,omitempty"`
	Tags               []TagLabel                 `json:"tags"`
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
	var tags []DesignTag
	var tag *DesignTag
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("link '%s' in rack '%s' has tag '%s' but tag missing from rack definition", o.Label, rack.Id, label)
		}
		tags = append(tags, *tag)
	}

	return &RackLink{
		Label:              o.Label,
		Tags:               tags,
		LinkPerSwitchCount: o.LinkPerSwitchCount,
		LinkSpeed:          o.LinkSpeed.parse(),
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     RackLinkAttachmentType(attachment),
		LagMode:            RackLinkLagMode(lagMode),
		SwitchPeer:         RackLinkSwitchPeer(switchPeer),
	}, nil
}

type RackElementGenericSystemRequest struct {
	Count            int
	AsnDomain        FeatureSwitch
	ManagementLevel  GenericSystemManagementLevel
	PortChannelIdMin int
	PortChannelIdMax int
	Loopback         FeatureSwitch
	Tags             []TagLabel
	Label            string
	Links            []RackLinkRequest
	LogicalDeviceId  ObjectId
}

func (o *RackElementGenericSystemRequest) raw() *rawRackElementGenericSystemRequest {
	var links []rawRackLinkRequest
	for _, link := range o.Links {
		links = append(links, *link.raw())
	}

	return &rawRackElementGenericSystemRequest{
		Count:            o.Count,
		AsnDomain:        featureSwitch(o.AsnDomain.String()),
		ManagementLevel:  genericSystemManagementLevel(o.ManagementLevel.String()),
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Loopback:         featureSwitch(o.Loopback.String()),
		Label:            o.Label,
		Links:            links,
		LogicalDevice:    o.LogicalDeviceId, // needs to be fetched from API, cloned into rack type on create() / update()
		Tags:             o.Tags,            // needs to be fetched from API, cloned into rack type on create() / update()
	}
}

type rawRackElementGenericSystemRequest struct {
	Count            int                          `json:"count"`
	AsnDomain        featureSwitch                `json:"asn_domain"`
	ManagementLevel  genericSystemManagementLevel `json:"management_level"`
	PortChannelIdMin int                          `json:"port_channel_id_min"`
	PortChannelIdMax int                          `json:"port_channel_id_max"`
	Loopback         featureSwitch                `json:"loopback"`
	Label            string                       `json:"label"`
	LogicalDevice    ObjectId                     `json:"logical_device"`
	Links            []rawRackLinkRequest         `json:"links"`
	Tags             []TagLabel                   `json:"tags,omitempty"`
}

type RackElementGenericSystem struct {
	Count            int
	AsnDomain        FeatureSwitch
	ManagementLevel  GenericSystemManagementLevel
	PortChannelIdMin int
	PortChannelIdMax int
	Loopback         FeatureSwitch
	Tags             []DesignTag
	Label            string
	Links            []RackLink
	Panels           []LogicalDevicePanel
	DisplayName      string
	LogicalDeviceId  ObjectId
}

type rawRackElementGenericSystem struct {
	Count            int                          `json:"count"`
	AsnDomain        featureSwitch                `json:"asn_domain"`
	ManagementLevel  genericSystemManagementLevel `json:"management_level"`
	PortChannelIdMin int                          `json:"port_channel_id_min"`
	PortChannelIdMax int                          `json:"port_channel_id_max"`
	Loopback         featureSwitch                `json:"loopback"`
	Tags             []TagLabel                   `json:"tags"`
	Label            string                       `json:"label"`
	LogicalDevice    ObjectId                     `json:"logical_device"`
	Links            []rawRackLink                `json:"links"`
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

	var tags []DesignTag
	var tag *DesignTag
	for _, label := range o.Tags {
		if tag, found = rack.tagByLabel(label); !found {
			return nil, fmt.Errorf("design tag '%s' not found in rack type '%s' definition", label, rack.Id)
		}
		tags = append(tags, *tag)
	}

	return &RackElementGenericSystem{
		Count:            o.Count,
		AsnDomain:        FeatureSwitch(asnDomain),
		ManagementLevel:  GenericSystemManagementLevel(mgmtLevel),
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Loopback:         FeatureSwitch(loopback),
		Tags:             tags,
		Label:            o.Label,
		Links:            links,
		Panels:           pld.Panels,
		DisplayName:      pld.DisplayName,
	}, nil
}

type RackTypeRequest struct {
	DisplayName              string
	Description              string
	FabricConnectivityDesign FabricConnectivityDesign
	LeafSwitches             []RackElementLeafSwitchRequest
	GenericSystems           []RackElementGenericSystemRequest
	AccessSwitches           []RackElementAccessSwitchRequest
	logicalDevices           []LogicalDevice
}

func (o *RackTypeRequest) raw(ctx context.Context, client *Client) (*rawRackTypeRequest, error) {
	result := &rawRackTypeRequest{
		DisplayName:              o.DisplayName,
		Description:              o.Description,
		FabricConnectivityDesign: o.FabricConnectivityDesign.raw(),
		Tags:                     nil, // populated by API calls below
		LogicalDevices:           nil, // populated by API calls below
		//LeafSwitches: make([]RackElementLeafSwitchRequest, len(o.LeafSwitches)), // todo: init array with correct size?
	}

	ldMap := make(map[ObjectId]struct{})  // collect IDs of all logical devices relevant to this rack
	tagMap := make(map[TagLabel]struct{}) // collect labels of all tags relevant to this rack

	for _, s := range o.LeafSwitches {
		result.LeafSwitches = append(result.LeafSwitches, *s.raw())
		ldMap[s.LogicalDeviceId] = struct{}{}
		for _, t := range s.Tags {
			tagMap[t] = struct{}{}
		}
	}

	for _, s := range o.AccessSwitches {
		result.AccessSwitches = append(result.AccessSwitches, *s.raw())
		ldMap[s.LogicalDeviceId] = struct{}{}
		for _, t := range s.Tags {
			tagMap[t] = struct{}{}
		}
	}

	for _, s := range o.GenericSystems {
		result.GenericSystems = append(result.GenericSystems, *s.raw())
		ldMap[s.LogicalDeviceId] = struct{}{}
		for _, t := range s.Tags {
			tagMap[t] = struct{}{}
		}
		for _, l := range s.Links {
			for _, t := range l.Tags {
				tagMap[t] = struct{}{}
			}
		}
	}

	for id := range ldMap {
		ld, err := client.GetLogicalDevice(ctx, id)
		if err != nil {
			return nil, err
		}
		result.LogicalDevices = append(result.LogicalDevices, *ld.raw())
	}

	for tl := range tagMap {
		tag, err := client.getTagByLabel(ctx, tl)
		if err != nil {
			return nil, err
		}
		result.Tags = append(result.Tags, *tag)
	}

	return result, nil
}

type rawRackTypeRequest struct {
	DisplayName              string                               `json:"display_name"`
	Description              string                               `json:"description"`
	FabricConnectivityDesign fabricConnectivityDesign             `json:"fabric_connectivity_design"`
	Tags                     []DesignTag                          `json:"tags,omitempty"`
	LogicalDevices           []rawLogicalDevice                   `json:"logical_devices,omitempty"`
	GenericSystems           []rawRackElementGenericSystemRequest `json:"generic_systems,omitempty"`
	LeafSwitches             []rawRackElementLeafSwitchRequest    `json:"leafs,omitempty"`
	AccessSwitches           []rawRackElementAccessSwitchRequest  `json:"access_switches,omitempty"`
}

type RackType struct {
	DisplayName              string
	Description              string
	FabricConnectivityDesign FabricConnectivityDesign
	Id                       ObjectId
	CreatedAt                time.Time
	LastModifiedAt           time.Time
	LeafSwitches             []RackElementLeafSwitch
	GenericSystems           []RackElementGenericSystem
	AccessSwitches           []RackElementAccessSwitch
}

type rawRackType struct {
	Id                       ObjectId                      `json:"id,omitempty"`
	DisplayName              string                        `json:"display_name"`
	Description              string                        `json:"description"`
	FabricConnectivityDesign fabricConnectivityDesign      `json:"fabric_connectivity_design"`
	Tags                     []DesignTag                   `json:"tags,omitempty"`
	CreatedAt                time.Time                     `json:"created_at"`
	LastModifiedAt           time.Time                     `json:"last_modified_at"`
	LogicalDevices           []rawLogicalDevice            `json:"logical_devices,omitempty"`
	GenericSystems           []rawRackElementGenericSystem `json:"generic_systems,omitempty"`
	LeafSwitches             []rawRackElementLeafSwitch    `json:"leafs,omitempty"`
	AccessSwitches           []rawRackElementAccessSwitch  `json:"access_switches,omitempty"`
}

func (o *rawRackType) polish() (*RackType, error) {
	fcd, err := o.FabricConnectivityDesign.parse()
	if err != nil {
		return nil, err
	}

	result := &RackType{
		DisplayName:              o.DisplayName,
		Description:              o.Description,
		FabricConnectivityDesign: FabricConnectivityDesign(fcd),
		Id:                       o.Id,
		CreatedAt:                o.CreatedAt,
		LastModifiedAt:           o.LastModifiedAt,
	}

	for _, r := range o.LeafSwitches {
		p, err := r.polish(o)
		if err != nil {
			return nil, err
		}
		result.LeafSwitches = append(result.LeafSwitches, *p)
	}

	for _, r := range o.AccessSwitches {
		p, err := r.polish(o)
		if err != nil {
			return nil, err
		}
		result.AccessSwitches = append(result.AccessSwitches, *p)
	}

	for _, r := range o.GenericSystems {
		p, err := r.polish(o)
		if err != nil {
			return nil, err
		}
		result.GenericSystems = append(result.GenericSystems, *p)
	}

	return result, nil
}

func (o rawRackType) logicalDeviceById(desired ObjectId) (*rawLogicalDevice, bool) {
	for _, ld := range o.LogicalDevices {
		if ld.Id == desired {
			return &ld, true
		}
	}
	return nil, false
}

func (o rawRackType) tagByLabel(desired TagLabel) (*DesignTag, bool) {
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

func (o *Client) getRackType(ctx context.Context, id ObjectId) (*RackType, error) {
	response := &rawRackType{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignRackTypeById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.polish()
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
		if rackType.DisplayName == name {
			return &rackType, nil
		}
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("rack type with name '%s' not found", name),
	}
}

func (o *Client) createRackType(ctx context.Context, request *RackTypeRequest) (ObjectId, error) {
	rawRequest, err := request.raw(ctx, o)
	if err != nil {
		return "", err
	}
	response := &objectIdResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignRackTypes,
		apiInput:    rawRequest,
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
