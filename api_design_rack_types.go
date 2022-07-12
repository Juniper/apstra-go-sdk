package goapstra

// RACK_LINKS_SCHEMA = s.List(s.Object({
//    'label': s.String(validate=s.Length(min=1)),
//    'target_switch_label': s.String(validate=s.Length(min=1)),
//    'link_per_switch_count': s.Integer(validate=s.Range(min=1)),
//    'link_speed': PortSpeedSchema,
//    'attachment_type': s.Enum(['singleAttached', 'dualAttached']),
//    'switch_peer': s.Optional(s.Enum(['first', 'second'])),
//    'lag_mode': s.Optional(s.Enum(['lacp_active', 'lacp_passive', 'static_lag'])),
//    'tags': TAG_LABEL_LIST,
//}, allow_extra_fields=False, default_field_type=s.IndexField))

// RACK_GENERIC_SYSTEM_SCHEMA = s.Object({
//    'label': s.String(validate=s.Length(min=1)),
//    'logical_device': s.String(validate=s.Length(min=1)),
//    'links': RACK_LINKS_SCHEMA,
//    'count': s.Optional(s.Integer(validate=s.Range(min=1)), load_default=1),
//    'management_level': s.ManagementLevel,
//    'port_channel_id_min': s.Optional(
//        s.Integer(validate=s.Range(min=0, max=CONSTANTS.max_port_channel_number)),
//        load_default=0),
//    'port_channel_id_max': s.Optional(
//        s.Integer(validate=s.Range(min=0, max=CONSTANTS.max_port_channel_number)),
//        load_default=0),
//    'loopback': s.Optional(s.OptionalProperty, load_default='disabled'),
//    'asn_domain': s.Optional(s.OptionalProperty, load_default='disabled'),
//    'tags': TAG_LABEL_LIST,
//}, allow_extra_fields=False, default_field_type=s.IndexField)

// RACK_TYPE_SCHEMA = s.Object({
//    'id': s.String(validate=s.Length(min=1)),
//    'display_name': s.Rackname(),
//    'description': s.Optional(s.String(), load_default=''),
//    'fabric_connectivity_design': s.Optional(s.Enum(['l3clos', 'l3collapsed']),
//                                             load_default='l3clos'),
//    'access_switches': s.Optional(
//        s.List(s.Object({
//            'label': s.String(validate=s.Length(min=1)),
//            'logical_device': s.String(),
//            'instance_count': s.Optional(s.Integer(validate=s.Range(min=1)),
//                                         load_default=1),
//            'links': RACK_LINKS_SCHEMA,
//            'access_access_link_count': s.Optional(
//                s.Integer(validate=s.Range(min=0)), load_default=0),
//            'access_access_link_speed': s.Optional(PortSpeedSchema),
//            'redundancy_protocol': s.Optional(s.Enum(['esi'])),
//            'tags': TAG_LABEL_LIST,
//        }, allow_extra_fields=False, default_field_type=s.IndexField)),
//        load_default=[]),
//    'status': s.Optional(s.Enum(['ok', 'inconsistent']), load_default='ok'),
//    'leafs': s.List(s.Object({
//        'label': s.String(validate=s.Length(min=1)),
//        'logical_device': s.String(),
//        'leaf_leaf_link_count': s.Optional(s.Integer(validate=s.Range(min=0)),
//                                           load_default=0),
//        'leaf_leaf_link_speed': s.Optional(PortSpeedSchema),
//        'leaf_leaf_l3_link_count': s.Optional(
//            s.Integer(validate=s.Range(min=0)), load_default=0),
//        'leaf_leaf_l3_link_speed': s.Optional(PortSpeedSchema),
//        'link_per_spine_count': s.Integer(validate=s.Range(min=0)),
//        'link_per_spine_speed': s.Optional(PortSpeedSchema),
//        'redundancy_protocol': s.Optional(s.Enum(['mlag', 'esi'])),
//        'leaf_leaf_link_port_channel_id': s.Optional(
//            s.Integer(validate=s.Range(min=0, max=CONSTANTS.max_port_channel_number)),
//            load_default=0),
//        'leaf_leaf_l3_link_port_channel_id': s.Optional(
//            s.Integer(validate=s.Range(min=0, max=CONSTANTS.max_port_channel_number)),
//            load_default=0),
//        'mlag_vlan_id': s.Optional(
//            s.Integer(validate=s.Range(min=0, max=CONSTANTS.max_vlan_id)),
//            load_default=0),
//        'tags': TAG_LABEL_LIST,
//    }, allow_extra_fields=False, default_field_type=s.IndexField)),
//    'servers': s.Optional(s.List(RACK_SERVER_SCHEMA), load_default=[]),
//    'generic_systems': s.Optional(s.List(RACK_GENERIC_SYSTEM_SCHEMA), load_default=[]),
//    'logical_devices': s.List(LOGICAL_DEVICE_SCHEMA,
//                              validate=s.Unique(key=lambda ld: ld['id'])),
//    'tags': ASSIGNED_TAGS_SCHEMA,
//    'created_at': s.Optional(s.String(),
//                             load_default="1970-01-01T00:00:00.000000Z"),
//    'last_modified_at': s.Optional(s.String(),
//                                   load_default="1970-01-01T00:00:00.000000Z"),
//}, validate=validate_rack_type, allow_extra_fields=False, default_field_type=s.IndexField)

// POST
// https://13.58.9.57:22409/api/design/rack-types
// {
//  "display_name": "yourname-single"
//  "description": "",
//  "last_modified_at": null,
//  "tags": [],
//  "leafs": [
//    {
//      "link_per_spine_count": 1,
//      "redundancy_protocol": null,
//      "leaf_leaf_link_speed": null,
//      "leaf_leaf_l3_link_count": 0,
//      "leaf_leaf_l3_link_speed": null,
//      "link_per_spine_speed": {
//        "unit": "G",
//        "value": 10
//      },
//      "label": "yourname-single",
//      "leaf_leaf_l3_link_port_channel_id": 0,
//      "leaf_leaf_link_port_channel_id": 0,
//      "logical_device": "virtual-7x10-1",
//      "leaf_leaf_link_count": 0
//    }
//  ],
//  "logical_devices": [
//    {
//      "created_at": "2022-07-08T13:48:38.033982Z",
//      "panels": [], // snip
//      "display_name": "virtual-7x10-1",
//      "id": "virtual-7x10-1",
//      "last_modified_at": "2022-07-08T13:48:38.033982Z",
//      "href": "#/design/logical-devices/virtual-7x10-1"
//    },
//    {
//      "created_at": "2022-04-22T06:08:55.568587Z",
//      "panels": [], // snip
//      "display_name": "AOS-1x10-1",
//      "id": "AOS-1x10-1",
//      "last_modified_at": "2022-04-22T06:08:55.568587Z",
//      "href": "#/design/logical-devices/AOS-1x10-1"
//    }
//  ],
//  "access_switches": [],
//  "fabric_connectivity_design": "l3clos", // or "l3collapsed"
//  "id": "yourname-single",
//  "generic_systems": [
//    {
//      "loopback": "disabled",
//      "asn_domain": "disabled",
//      "port_channel_id_max": 0,
//      "label": "single-server",
//      "count": 1,
//      "management_level": "unmanaged",
//      "logical_device": "AOS-1x10-1",
//      "links": [
//        {
//          "link_per_switch_count": 1,
//          "link_speed": {
//            "unit": "G",
//            "value": 10
//          },
//          "target_switch_label": "yourname-single",
//          "lag_mode": null,
//          "switch_peer": null,
//          "attachment_type": "singleAttached",
//          "label": "single-link"
//        }
//      ],
//      "port_channel_id_min": 0
//    }
//  ],
//}

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	apiUrlDesignRackTypes       = apiUrlDesignPrefix + "rack-types"
	apiUrlDesignRackTypesPrefix = apiUrlDesignRackTypes + apiUrlPathDelim
	apiUrlDesignRackTypeById    = apiUrlDesignRackTypesPrefix + "%s"

	leafSwitchLogicalDeviceIdPrefix    = "leaf-"
	accessSwitchLogicalDeviceIdPrefix  = "access-"
	genericSystemLogicalDeviceIdPrefix = "generic-"
)

type AccessRedundancyProtocol int
type accessRedundancyProtocol string

const (
	AccessRedundancyProtocolNone = AccessRedundancyProtocol(iota)
	AccessRedundancyProtocolEsi

	accessRedundancyProtocolNone    = accessRedundancyProtocol("")
	accessRedundancyProtocolEsi     = accessRedundancyProtocol("esi")
	accessRedundancyProtocolUnknown = "unknown type %d"
)

type LeafRedundancyProtocol int
type leafRedundancyProtocol string

const (
	LeafRedundancyProtocolNone = LeafRedundancyProtocol(iota)
	LeafRedundancyProtocolMlag
	LeafRedundancyProtocolEsi

	leafRedundancyProtocolNone    = leafRedundancyProtocol("")
	leafRedundancyProtocolMlag    = leafRedundancyProtocol("mlag")
	leafRedundancyProtocolEsi     = leafRedundancyProtocol("esi")
	leafRedundancyProtocolUnknown = "unknown type %d"
)

type FabricConnectivityDesign int
type fabricConnectivityDesign string

const (
	FabricCOnnectivityDesignL3Clos = FabricConnectivityDesign(iota)
	FabricCOnnectivityDesignL3Collapsed

	fabricConnectivityDesignL3Clos      = fabricConnectivityDesign("l3clos")
	fabricConnectivityDesignL3Collapsed = fabricConnectivityDesign("l3collapsed")
	fabricConnectivityDesignUnknown     = "unknown type %d"
)

type FeatureSwitch int
type featureSwitch string

const (
	FeatureSwitchDisabled = FeatureSwitch(iota)
	FeatureSwitchEnabled

	featureSwitchDisabled = featureSwitch("disabled")
	featureSwitchEnabled  = featureSwitch("enabled")
	featureSwitchUnknown  = "unknown feature switch state '%d'"
)

type GenericSystemManagementLevel int
type genericSystemManagementLevel string

const (
	// unmanaged, telemetry_only or full_control
	GenericSystemUnmanaged = GenericSystemManagementLevel(iota)
	GenericSystemTelemetryOnly
	GenericSystemFullControl

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

	rackLinkAttachmentTypeSingle  = rackLinkAttachmentType("singleAttached")
	rackLinkAttachmentTypeDual    = rackLinkAttachmentType("dualAttached")
	rackLinkAttachmentTypeUnknown = "unknown generic system attachment type '%d'"
)

type RackLinkLagMode int
type rackLinkLagMode string

const (
	RackLinkLagModeLacpNone = RackLinkLagMode(iota)
	RackLinkLagModeLacpActive
	RackLinkLagModeLacpPassive
	RackLinkLagModeLacpStatic

	rackLinkLagModeLacpNone    = rackLinkLagMode("")
	rackLinkLagModeLacpActive  = rackLinkLagMode("lacp_active")
	rackLinkLagModeLacpPassive = rackLinkLagMode("lacp_passive")
	rackLinkLagModeLacpStatic  = rackLinkLagMode("static_lag")
	rackLinkLagModeUnknown     = "unknown lag mode '%d'"
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
		return 0, fmt.Errorf("unknown access redundancy protocol '%s'", o)
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
		return 0, fmt.Errorf("unknown leaf redundancy protocol '%s'", o)
	}
}

func (o FabricConnectivityDesign) Int() int {
	return int(o)
}

func (o FabricConnectivityDesign) String() string {
	switch o {
	case FabricCOnnectivityDesignL3Clos:
		return string(fabricConnectivityDesignL3Clos)
	case FabricCOnnectivityDesignL3Collapsed:
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
		return int(FabricCOnnectivityDesignL3Clos), nil
	case fabricConnectivityDesignL3Collapsed:
		return int(FabricCOnnectivityDesignL3Collapsed), nil
	default:
		return 0, fmt.Errorf("unknown fabric connectivity design '%s'", o)
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
		return 0, fmt.Errorf("unknown feature state '%s'", o)
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
		return 0, fmt.Errorf("unknown generic system management state '%s'", o)
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
		return 0, fmt.Errorf("unknown generic system link attachment type '%s'", o)
	}
}

func (o RackLinkLagMode) Int() int {
	return int(o)
}

func (o RackLinkLagMode) String() string {
	switch o {
	case RackLinkLagModeLacpNone:
		return string(rackLinkLagModeLacpNone)
	case RackLinkLagModeLacpActive:
		return string(rackLinkLagModeLacpActive)
	case RackLinkLagModeLacpPassive:
		return string(rackLinkLagModeLacpPassive)
	case RackLinkLagModeLacpStatic:
		return string(rackLinkLagModeLacpStatic)
	default:
		return fmt.Sprintf(rackLinkLagModeUnknown, o)
	}
}

func (o rackLinkLagMode) string() string {
	return string(o)
}
func (o rackLinkLagMode) parse() (int, error) {
	switch o {
	case rackLinkLagModeLacpNone:
		return int(RackLinkLagModeLacpNone), nil
	case rackLinkLagModeLacpActive:
		return int(RackLinkLagModeLacpActive), nil
	case rackLinkLagModeLacpPassive:
		return int(RackLinkLagModeLacpPassive), nil
	case rackLinkLagModeLacpStatic:
		return int(RackLinkLagModeLacpStatic), nil
	default:
		return 0, fmt.Errorf("unknown link LACP mode '%s'", o)
	}
}

type RackTag struct {
	Id             ObjectId  `json:"id"`
	Label          string    `json:"label"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

type optionsRackTypeResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type RackElementLeafSwitch struct {
	Label                       string
	LeafLeafL3LinkCount         int
	LeafLeafL3LinkPortChannelId int
	LeafLeafL3LinkSpeed         *LogicalDevicePortSpeed
	LeafLeafLinkCount           int
	LeafLeafLinkPortChannelId   int
	LeafLeafLinkSpeed           *LogicalDevicePortSpeed
	LinkPerSpineCount           int
	LinkPerSpineSpeed           *LogicalDevicePortSpeed
	MlagVlanId                  int
	RedundancyProtocol          LeafRedundancyProtocol
	Tags                        []RackTag
	Panels                      []LogicalDevicePanel
	DisplayName                 string
}

func (o *RackElementLeafSwitch) raw(logicalDeviceId string) *rawRackElementLeaf {
	tags := o.Tags
	if tags == nil {
		tags = []RackTag{}
	}

	return &rawRackElementLeaf{
		Label:                       o.Label,
		LeafLeafL3LinkCount:         o.LeafLeafL3LinkCount,
		LeafLeafL3LinkPortChannelId: o.LeafLeafL3LinkPortChannelId,
		LeafLeafL3LinkSpeed:         o.LeafLeafL3LinkSpeed,
		LeafLeafLinkCount:           o.LeafLeafLinkCount,
		LeafLeafLinkPortChannelId:   o.LeafLeafLinkPortChannelId,
		LeafLeafLinkSpeed:           o.LeafLeafLinkSpeed,
		LinkPerSpineCount:           o.LinkPerSpineCount,
		LinkPerSpineSpeed:           o.LinkPerSpineSpeed,
		LogicalDevice:               logicalDeviceId,
		MlagVlanId:                  o.MlagVlanId,
		RedundancyProtocol:          o.RedundancyProtocol.raw(),
		Tags:                        tags,
	}
}

type rawRackElementLeaf struct {
	Label                       string                  `json:"label"`
	LeafLeafL3LinkCount         int                     `json:"leaf_leaf_l3_link_count"`
	LeafLeafL3LinkPortChannelId int                     `json:"leaf_leaf_l3_link_port_channel_id"`
	LeafLeafL3LinkSpeed         *LogicalDevicePortSpeed `json:"leaf_leaf_l3_link_speed"`
	LeafLeafLinkCount           int                     `json:"leaf_leaf_link_count"`
	LeafLeafLinkPortChannelId   int                     `json:"leaf_leaf_link_port_channel_id"`
	LeafLeafLinkSpeed           *LogicalDevicePortSpeed `json:"leaf_leaf_link_speed"`
	LinkPerSpineCount           int                     `json:"link_per_spine_count"`
	LinkPerSpineSpeed           *LogicalDevicePortSpeed `json:"link_per_spine_speed"`
	LogicalDevice               string                  `json:"logical_device"`
	MlagVlanId                  int                     `json:"mlag_vlan_id"`
	RedundancyProtocol          leafRedundancyProtocol  `json:"redundancy_protocol,omitempty"`
	Tags                        []RackTag               `json:"tags"`
}

func (o *rawRackElementLeaf) polish(ld LogicalDevice) (*RackElementLeafSwitch, error) {
	rp, err := o.RedundancyProtocol.parse()
	if err != nil {
		return nil, err
	}

	return &RackElementLeafSwitch{
		Label:                       o.Label,
		LeafLeafL3LinkCount:         o.LeafLeafL3LinkCount,
		LeafLeafL3LinkPortChannelId: o.LeafLeafL3LinkPortChannelId,
		LeafLeafL3LinkSpeed:         o.LeafLeafL3LinkSpeed,
		LeafLeafLinkCount:           o.LeafLeafLinkCount,
		LeafLeafLinkPortChannelId:   o.LeafLeafLinkPortChannelId,
		LeafLeafLinkSpeed:           o.LeafLeafLinkSpeed,
		LinkPerSpineCount:           o.LinkPerSpineCount,
		LinkPerSpineSpeed:           o.LinkPerSpineSpeed,
		MlagVlanId:                  o.MlagVlanId,
		RedundancyProtocol:          LeafRedundancyProtocol(rp),
		Tags:                        o.Tags,
		Panels:                      ld.Panels,
		DisplayName:                 ld.DisplayName,
	}, nil
}

type AccessSwitchLink struct {
	Label              string                 `json:"label"`
	LinkPerSwitchCount int                    `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed `json:"link_speed"`
	TargetSwitchLabel  string                 `json:"target_switch_label"`
	LagMode            string                 `json:"lag_mode,omitempty"`
	SwitchPeer         interface{}            `json:"switch_peer"` // todo - what is this?
	AttachmentType     string                 `json:"attachment_type"`
}

type RackElementAccessSwitch struct {
	InstanceCount      int
	RedundancyProtocol AccessRedundancyProtocol
	Links              []AccessSwitchLink
	Label              string
	Panels             []LogicalDevicePanel
	DisplayName        string
}

func (o *RackElementAccessSwitch) raw(logicalDeviceId string) *rawRackElementAccessSwitch {
	return &rawRackElementAccessSwitch{
		InstanceCount:      o.InstanceCount,
		RedundancyProtocol: o.RedundancyProtocol.raw(),
		Links:              o.Links,
		Label:              o.Label,
		LogicalDevice:      logicalDeviceId,
	}
}

type rawRackElementAccessSwitch struct {
	InstanceCount      int                      `json:"instance_count"`
	RedundancyProtocol accessRedundancyProtocol `json:"redundancy_protocol,omitempty"`
	Links              []AccessSwitchLink       `json:"links"`
	Label              string                   `json:"label"`
	LogicalDevice      string                   `json:"logical_device"`
}

func (o *rawRackElementAccessSwitch) polish(ld LogicalDevice) (*RackElementAccessSwitch, error) {
	rp, err := o.RedundancyProtocol.parse()
	if err != nil {
		return nil, err
	}

	return &RackElementAccessSwitch{
		InstanceCount:      o.InstanceCount,
		RedundancyProtocol: AccessRedundancyProtocol(rp),
		Links:              o.Links,
		Label:              o.Label,
		Panels:             ld.Panels,
		DisplayName:        ld.DisplayName,
	}, nil
}

type GenericSystemAccessLink struct {
	Label              string                 `json:"label"`
	Tags               []RackTag              `json:"tags"`
	LinkPerSwitchCount int                    `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed `json:"link_speed"`
	TargetSwitchLabel  string                 `json:"target_switch_label"`
	AttachmentType     RackLinkAttachmentType `json:"attachment_type"`
	LagMode            RackLinkLagMode        `json:"lag_mode,omitempty"`
}

func (o GenericSystemAccessLink) raw() *rawGenericSystemAccessLink {
	var tags []RackTag
	for _, tag := range o.Tags {
		tags = append(tags, tag)
	}
	if tags == nil {
		tags = []RackTag{}
	}
	return &rawGenericSystemAccessLink{
		Label:              o.Label,
		Tags:               tags,
		LinkPerSwitchCount: o.LinkPerSwitchCount,
		LinkSpeed:          o.LinkSpeed,
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     rackLinkAttachmentType(o.AttachmentType.String()),
		LagMode:            rackLinkLagMode(o.LagMode.String()),
	}
}

type rawGenericSystemAccessLink struct {
	Label              string                 `json:"label"`
	Tags               []RackTag              `json:"tags"`
	LinkPerSwitchCount int                    `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed `json:"link_speed"`
	TargetSwitchLabel  string                 `json:"target_switch_label"`
	AttachmentType     rackLinkAttachmentType `json:"attachment_type"`
	LagMode            rackLinkLagMode        `json:"lag_mode"`
}

func (o rawGenericSystemAccessLink) polish() (*GenericSystemAccessLink, error) {
	attachment, err := o.AttachmentType.parse()
	if err != nil {
		return nil, err
	}

	lagMode, err := o.LagMode.parse()
	if err != nil {
		return nil, err
	}

	return &GenericSystemAccessLink{
		Label:              o.Label,
		Tags:               o.Tags,
		LinkPerSwitchCount: o.LinkPerSwitchCount,
		LinkSpeed:          o.LinkSpeed,
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     RackLinkAttachmentType(attachment),
		LagMode:            RackLinkLagMode(lagMode),
	}, nil
}

type RackElementGenericSystem struct {
	Count            int
	AsnDomain        FeatureSwitch
	ManagementLevel  GenericSystemManagementLevel
	PortChannelIdMin int
	PortChannelIdMax int
	Loopback         FeatureSwitch
	Tags             []RackTag
	Label            string
	Links            []GenericSystemAccessLink
	Panels           []LogicalDevicePanel
	DisplayName      string
}

func (o RackElementGenericSystem) raw(logicalDeviceId string) *rawRackElementGenericSystem {
	tags := o.Tags
	if tags == nil {
		tags = []RackTag{}
	}

	var links []rawGenericSystemAccessLink
	for _, link := range o.Links {
		links = append(links, *link.raw())
	}

	return &rawRackElementGenericSystem{
		Count:            o.Count,
		AsnDomain:        featureSwitch(o.AsnDomain.String()),
		ManagementLevel:  genericSystemManagementLevel(o.ManagementLevel.String()),
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Loopback:         featureSwitch(o.Loopback.String()),
		Tags:             tags,
		Label:            o.Label,
		LogicalDevice:    logicalDeviceId,
		Links:            links,
	}
}

type rawRackElementGenericSystem struct {
	Count            int                          `json:"count"`
	AsnDomain        featureSwitch                `json:"asn_domain"`
	ManagementLevel  genericSystemManagementLevel `json:"management_level"`
	PortChannelIdMin int                          `json:"port_channel_id_min"`
	PortChannelIdMax int                          `json:"port_channel_id_max"`
	Loopback         featureSwitch                `json:"loopback"`
	Tags             []RackTag                    `json:"tags"`
	Label            string                       `json:"label"`
	LogicalDevice    string                       `json:"logical_device"`
	Links            []rawGenericSystemAccessLink `json:"links"`
}

func (o *rawRackElementGenericSystem) polish(ld LogicalDevice) (*RackElementGenericSystem, error) {
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

	var links []GenericSystemAccessLink
	for _, link := range o.Links {
		p, err := link.polish()
		if err != nil {
			return nil, err
		}
		links = append(links, *p)
	}

	return &RackElementGenericSystem{
		Count:            o.Count,
		AsnDomain:        FeatureSwitch(asnDomain),
		ManagementLevel:  GenericSystemManagementLevel(mgmtLevel),
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Loopback:         FeatureSwitch(loopback),
		Tags:             o.Tags,
		Label:            o.Label,
		Links:            links,
		Panels:           ld.Panels,
		DisplayName:      ld.DisplayName,
	}, nil
}

type RackType struct {
	DisplayName              string
	Description              string
	FabricConnectivityDesign FabricConnectivityDesign
	Id                       ObjectId
	Tags                     []RackTag
	CreatedAt                time.Time
	LastModifiedAt           time.Time
	LeafSwitches             []RackElementLeafSwitch
	GenericSystems           []RackElementGenericSystem
	AccessSwitches           []RackElementAccessSwitch
}

func (o *RackType) raw() *rawRackType {
	result := &rawRackType{
		Id:                       o.Id,
		DisplayName:              o.DisplayName,
		Description:              o.Description,
		FabricConnectivityDesign: o.FabricConnectivityDesign.raw(),
		Tags:                     o.Tags,
		CreatedAt:                o.CreatedAt,
		LastModifiedAt:           o.LastModifiedAt,
	}

	for k, v := range o.LeafSwitches {
		result.LeafSwitchess = append(result.LeafSwitchess, *v.raw(leafSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)))
		if _, found := result.logicalDeviceById(leafSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)); !found {
			result.LogicalDevices = append(result.LogicalDevices, LogicalDevice{
				Panels:      v.Panels,
				DisplayName: v.DisplayName,
				Id:          ObjectId(leafSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)),
			})
		}
	}

	for k, v := range o.AccessSwitches {
		result.AccessSwitches = append(result.AccessSwitches, *v.raw(accessSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)))
		if _, found := result.logicalDeviceById(accessSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)); !found {
			result.LogicalDevices = append(result.LogicalDevices, LogicalDevice{
				Panels:      v.Panels,
				DisplayName: v.DisplayName,
				Id:          ObjectId(accessSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)),
			})
		}

	}

	for k, v := range o.GenericSystems {
		result.GenericSystems = append(result.GenericSystems, *v.raw(genericSystemLogicalDeviceIdPrefix + strconv.Itoa(k)))
		if _, found := result.logicalDeviceById(genericSystemLogicalDeviceIdPrefix + strconv.Itoa(k)); !found {
			result.LogicalDevices = append(result.LogicalDevices, LogicalDevice{
				Panels:      v.Panels,
				DisplayName: v.DisplayName,
				Id:          ObjectId(genericSystemLogicalDeviceIdPrefix + strconv.Itoa(k)),
			})
		}
	}
	return result
}

type rawRackType struct {
	Id                       ObjectId                      `json:"id,omitempty"`
	DisplayName              string                        `json:"display_name"`
	Description              string                        `json:"description"`
	FabricConnectivityDesign fabricConnectivityDesign      `json:"fabric_connectivity_design"`
	Tags                     []RackTag                     `json:"tags,omitempty"`
	CreatedAt                time.Time                     `json:"created_at"`
	LastModifiedAt           time.Time                     `json:"last_modified_at"`
	LogicalDevices           []LogicalDevice               `json:"logical_devices,omitempty"`
	GenericSystems           []rawRackElementGenericSystem `json:"generic_systems,omitempty"`
	LeafSwitchess            []rawRackElementLeaf          `json:"leafs,omitempty"`
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
		Tags:                     o.Tags,
		CreatedAt:                o.CreatedAt,
		LastModifiedAt:           o.LastModifiedAt,
	}

	var ld *LogicalDevice
	var found bool

	for _, r := range o.LeafSwitchess {
		if ld, found = o.logicalDeviceById(r.LogicalDevice); !found {
			return nil, fmt.Errorf("logical device '%s' not found in rack '%s'", r.LogicalDevice, o.Id)
		}
		p, err := r.polish(*ld)
		if err != nil {
			return nil, err
		}
		result.LeafSwitches = append(result.LeafSwitches, *p)
	}

	for _, r := range o.AccessSwitches {
		if ld, found = o.logicalDeviceById(r.LogicalDevice); !found {
			return nil, fmt.Errorf("logical device '%s' not found in rack '%s'", r.LogicalDevice, o.Id)
		}
		p, err := r.polish(*ld)
		if err != nil {
			return nil, err
		}
		result.AccessSwitches = append(result.AccessSwitches, *p)
	}

	for _, r := range o.GenericSystems {
		if ld, found = o.logicalDeviceById(r.LogicalDevice); !found {
			return nil, fmt.Errorf("logical device '%s' not found in rack '%s'", r.LogicalDevice, o.Id)
		}
		p, err := r.polish(*ld)
		if err != nil {
			return nil, err
		}
		result.GenericSystems = append(result.GenericSystems, *p)
	}

	return result, nil
}

func (o rawRackType) logicalDeviceById(id string) (*LogicalDevice, bool) {
	for _, ld := range o.LogicalDevices {
		if ld.Id == ObjectId(id) {
			return &ld, true
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

func (o *Client) createRackType(ctx context.Context, rackType *RackType) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignRackTypes,
		apiInput:    rackType.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}
