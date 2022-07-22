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
	"time"
)

const (
	apiUrlDesignRackTypes       = apiUrlDesignPrefix + "rack-types"
	apiUrlDesignRackTypesPrefix = apiUrlDesignRackTypes + apiUrlPathDelim
	apiUrlDesignRackTypeById    = apiUrlDesignRackTypesPrefix + "%s"

	leafSwitchLogicalDeviceIdPrefix    = "leaf-"
	accessSwitchLogicalDeviceIdPrefix  = "access-"
	genericSystemLogicalDeviceIdPrefix = "generic-"

	errLdIdRequired = "%s logical device id cannot be empty"
	errNotEmpty     = "%s %s element must be empty for this operation"
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
	Tags                        []TagLabel
	Panels                      []LogicalDevicePanel
	DisplayName                 string
	LogicalDeviceId             ObjectId
}

func (o *RackElementLeafSwitch) raw() *rawRackElementLeaf {
	tags := o.Tags
	if tags == nil {
		tags = []TagLabel{}
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
		LogicalDevice:               o.LogicalDeviceId,
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
	LogicalDevice               ObjectId                `json:"logical_device"`
	MlagVlanId                  int                     `json:"mlag_vlan_id"`
	RedundancyProtocol          leafRedundancyProtocol  `json:"redundancy_protocol,omitempty"`
	Tags                        []TagLabel              `json:"tags"`
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

//type AccessSwitchLink struct {
//	Label              string                 `json:"label"`
//	LinkPerSwitchCount int                    `json:"link_per_switch_count"`
//	LinkSpeed          LogicalDevicePortSpeed `json:"link_speed"`
//	TargetSwitchLabel  string                 `json:"target_switch_label"`
//	LagMode            string                 `json:"lag_mode,omitempty"`
//	SwitchPeer         RackLinkSwitchPeer     `json:"switch_peer"` // todo - what is this?
//	AttachmentType     string                 `json:"attachment_type"`
//}

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

func (o *RackElementAccessSwitch) raw() *rawRackElementAccessSwitch {
	return &rawRackElementAccessSwitch{
		InstanceCount:         o.InstanceCount,
		RedundancyProtocol:    o.RedundancyProtocol.raw(),
		Links:                 o.Links,
		Label:                 o.Label,
		LogicalDevice:         o.LogicalDeviceId,
		Tags:                  o.Tags,
		AccessAccessLinkCount: o.AccessAccessLinkCount,
		AccessAccessLinkSpeed: o.AccessAccessLinkSpeed,
	}
}

type rawRackElementAccessSwitch struct {
	InstanceCount         int                      `json:"instance_count"`
	RedundancyProtocol    accessRedundancyProtocol `json:"redundancy_protocol,omitempty"`
	Links                 []RackLink               `json:"links"`
	Label                 string                   `json:"label"`
	LogicalDevice         ObjectId                 `json:"logical_device"`
	AccessAccessLinkCount int                      `json:"access_access_link_count"`
	AccessAccessLinkSpeed LogicalDevicePortSpeed   `json:"access_access_link_speed"`
	Tags                  []DesignTag              `json:"tags"`
}

func (o *rawRackElementAccessSwitch) polish(ld LogicalDevice) (*RackElementAccessSwitch, error) {
	rp, err := o.RedundancyProtocol.parse()
	if err != nil {
		return nil, err
	}

	return &RackElementAccessSwitch{
		InstanceCount:         o.InstanceCount,
		RedundancyProtocol:    AccessRedundancyProtocol(rp),
		Links:                 o.Links,
		Label:                 o.Label,
		Panels:                ld.Panels,
		DisplayName:           ld.DisplayName,
		AccessAccessLinkCount: o.AccessAccessLinkCount,
		AccessAccessLinkSpeed: o.AccessAccessLinkSpeed,
		Tags:                  o.Tags,
	}, nil
}

type RackLink struct {
	Label              string                 // `json:"label"`
	Tags               []DesignTag            // `json:"tags"`
	LinkPerSwitchCount int                    // `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed // `json:"link_speed"`
	TargetSwitchLabel  string                 // `json:"target_switch_label"`
	AttachmentType     RackLinkAttachmentType // `json:"attachment_type"`
	LagMode            RackLinkLagMode        // `json:"lag_mode"`
	SwitchPeer         RackLinkSwitchPeer     // `json:"switch_peer"`
}

func (o RackLink) raw() *rawRackLink {
	var tags []DesignTag
	for _, tag := range o.Tags {
		tags = append(tags, tag)
	}
	if tags == nil {
		tags = []DesignTag{}
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
		LinkSpeed:          o.LinkSpeed,
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     rackLinkAttachmentType(o.AttachmentType.String()),
		LagMode:            lagModePtr,
		SwitchPeer:         rackLinkSwitchPeer(o.LagMode.String()),
	}
}

type rawRackLink struct {
	Label              string                 `json:"label"`
	Tags               []DesignTag            `json:"tags"`
	LinkPerSwitchCount int                    `json:"link_per_switch_count"`
	LinkSpeed          LogicalDevicePortSpeed `json:"link_speed"`
	TargetSwitchLabel  string                 `json:"target_switch_label"`
	AttachmentType     rackLinkAttachmentType `json:"attachment_type"`
	LagMode            *rackLinkLagMode       `json:"lag_mode"` // do not "omitempty" // todo: explore this b/c the API sends 'null'
	SwitchPeer         rackLinkSwitchPeer     `json:"switch_peer,omitempty"`
}

func (o rawRackLink) polish() (*RackLink, error) {
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

	return &RackLink{
		Label:              o.Label,
		Tags:               o.Tags,
		LinkPerSwitchCount: o.LinkPerSwitchCount,
		LinkSpeed:          o.LinkSpeed,
		TargetSwitchLabel:  o.TargetSwitchLabel,
		AttachmentType:     RackLinkAttachmentType(attachment),
		LagMode:            RackLinkLagMode(lagMode),
		SwitchPeer:         RackLinkSwitchPeer(switchPeer),
	}, nil
}

type RackElementGenericSystem struct {
	Count            int
	AsnDomain        FeatureSwitch
	ManagementLevel  GenericSystemManagementLevel
	PortChannelIdMin int
	PortChannelIdMax int
	Loopback         FeatureSwitch
	Tags             []TagLabel
	Label            string
	Links            []RackLink
	Panels           []LogicalDevicePanel
	DisplayName      string
	LogicalDeviceId  ObjectId
}

func (o RackElementGenericSystem) raw() *rawRackElementGenericSystem {
	tags := o.Tags
	if tags == nil {
		tags = []TagLabel{}
	}

	var links []rawRackLink
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
		LogicalDevice:    o.LogicalDeviceId,
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
	Tags             []TagLabel                   `json:"tags"`
	Label            string                       `json:"label"`
	LogicalDevice    ObjectId                     `json:"logical_device"`
	Links            []rawRackLink                `json:"links"`
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

	var links []RackLink
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
	Tags                     []DesignTag
	CreatedAt                time.Time
	LastModifiedAt           time.Time
	LeafSwitches             []RackElementLeafSwitch
	GenericSystems           []RackElementGenericSystem
	AccessSwitches           []RackElementAccessSwitch
	logicalDevices           []LogicalDevice
}

func (o RackType) raw() *rawRackType {
	result := &rawRackType{
		Id:                       o.Id,
		DisplayName:              o.DisplayName,
		Description:              o.Description,
		FabricConnectivityDesign: o.FabricConnectivityDesign.raw(),
		Tags:                     o.Tags,
		CreatedAt:                o.CreatedAt,
		LastModifiedAt:           o.LastModifiedAt,
	}

	for _, ld := range o.logicalDevices {
		result.LogicalDevices = append(result.LogicalDevices, *ld.raw())
	}

	for _, v := range o.LeafSwitches {
		result.LeafSwitchess = append(result.LeafSwitchess, *v.raw())
		//if _, found := result.logicalDeviceById(v.LogicalDeviceId); !found {
		//	result.LogicalDevices = append(result.LogicalDevices, *LogicalDevice{
		//		Panels:      v.Panels,
		//		DisplayName: v.DisplayName,
		//		Id:          ObjectId(leafSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)),
		//	}.raw())
		//}
	}

	for _, v := range o.AccessSwitches {
		result.AccessSwitches = append(result.AccessSwitches, *v.raw())
		//if _, found := result.logicalDeviceById(v.LogicalDeviceId); !found {
		//	result.LogicalDevices = append(result.LogicalDevices, *LogicalDevice{
		//		Panels:      v.Panels,
		//		DisplayName: v.DisplayName,
		//		Id:          ObjectId(accessSwitchLogicalDeviceIdPrefix + strconv.Itoa(k)),
		//	}.raw())
		//}

	}

	for _, v := range o.GenericSystems {
		result.GenericSystems = append(result.GenericSystems, *v.raw())
		//if _, found := result.logicalDeviceById(v.LogicalDeviceId); !found {
		//	result.LogicalDevices = append(result.LogicalDevices, *LogicalDevice{
		//		Panels:      v.Panels,
		//		DisplayName: v.DisplayName,
		//		Id:          ObjectId(genericSystemLogicalDeviceIdPrefix + strconv.Itoa(k)),
		//	}.raw())
		//}
	}
	return result
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

	var ld *rawLogicalDevice
	var found bool

	for _, r := range o.LeafSwitchess {
		if ld, found = o.logicalDeviceById(r.LogicalDevice); !found {
			return nil, fmt.Errorf("logical device '%s' not found in rack '%s'", r.LogicalDevice, o.Id)
		}
		ldp, err := ld.polish()
		if err != nil {
			return nil, err
		}
		p, err := r.polish(*ldp)
		if err != nil {
			return nil, err
		}
		result.LeafSwitches = append(result.LeafSwitches, *p)
	}

	for _, r := range o.AccessSwitches {
		if ld, found = o.logicalDeviceById(r.LogicalDevice); !found {
			return nil, fmt.Errorf("logical device '%s' not found in rack '%s'", r.LogicalDevice, o.Id)
		}
		ldp, err := ld.polish()
		if err != nil {
			return nil, err
		}
		p, err := r.polish(*ldp)
		if err != nil {
			return nil, err
		}
		result.AccessSwitches = append(result.AccessSwitches, *p)
	}

	for _, r := range o.GenericSystems {
		if ld, found = o.logicalDeviceById(r.LogicalDevice); !found {
			return nil, fmt.Errorf("logical device '%s' not found in rack '%s'", r.LogicalDevice, o.Id)
		}
		ldp, err := ld.polish()
		if err != nil {
			return nil, err
		}
		p, err := r.polish(*ldp)
		if err != nil {
			return nil, err
		}
		result.GenericSystems = append(result.GenericSystems, *p)
	}

	return result, nil
}

func (o rawRackType) logicalDeviceById(id ObjectId) (*rawLogicalDevice, bool) {
	for _, ld := range o.LogicalDevices {
		if ld.Id == id {
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

func (o *Client) createRackType(ctx context.Context, rackType *RackType) (ObjectId, error) {
	err := rackType.populateLogicalDeviceDetailsFromGlobalCatalog(ctx, o)
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	response := &objectIdResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
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

func (o *Client) updateRackType(ctx context.Context, id ObjectId, rackType *RackType) (ObjectId, error) {
	err := rackType.populateLogicalDeviceDetailsFromGlobalCatalog(ctx, o)
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	response := &objectIdResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPut,
		urlStr:      fmt.Sprintf(apiUrlDesignRackTypeById, id),
		apiInput:    rackType.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) deleteRackType(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodGet,
		urlStr: fmt.Sprintf(apiUrlDesignRackTypeById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *RackType) populateLogicalDeviceDetailsFromGlobalCatalog(ctx context.Context, client *Client) error {
	ldMap := make(map[ObjectId]struct{}) // for keeping track of logical devices retrieved from the API

	for _, i := range o.LeafSwitches {
		switch {
		case i.LogicalDeviceId == "":
			return fmt.Errorf(errLdIdRequired, "leaf switch")
		case len(i.Panels) != 0:
			return fmt.Errorf(errNotEmpty, "leaf switch", "[]Panels")
		case i.DisplayName != "":
			return fmt.Errorf(errNotEmpty, "leaf switch", "DisplayName")
		}

		if _, ok := ldMap[i.LogicalDeviceId]; !ok {
			ld, err := client.GetLogicalDevice(ctx, i.LogicalDeviceId)
			if err != nil {
				return err
			}
			ldMap[i.LogicalDeviceId] = struct{}{}
			o.logicalDevices = append(o.logicalDevices, *ld)
		}
	}
	for _, i := range o.AccessSwitches {
		switch {
		case i.LogicalDeviceId == "":
			return fmt.Errorf(errLdIdRequired, "access switch")
		case len(i.Panels) != 0:
			return fmt.Errorf(errNotEmpty, "access switch", "[]Panels")
		case i.DisplayName != "":
			return fmt.Errorf(errNotEmpty, "access switch", "DisplayName")
		}

		if _, ok := ldMap[i.LogicalDeviceId]; !ok {
			ld, err := client.GetLogicalDevice(ctx, i.LogicalDeviceId)
			if err != nil {
				return err
			}
			ldMap[i.LogicalDeviceId] = struct{}{}
			o.logicalDevices = append(o.logicalDevices, *ld)
		}
	}
	for _, i := range o.GenericSystems {
		switch {
		case i.LogicalDeviceId == "":
			return fmt.Errorf(errLdIdRequired, "generic system")
		case len(i.Panels) != 0:
			return fmt.Errorf(errNotEmpty, "generic system", "[]Panels")
		case i.DisplayName != "":
			return fmt.Errorf(errNotEmpty, "generic system", "DisplayName")
		}

		if _, ok := ldMap[i.LogicalDeviceId]; !ok {
			ld, err := client.GetLogicalDevice(ctx, i.LogicalDeviceId)
			if err != nil {
				return err
			}
			ldMap[i.LogicalDeviceId] = struct{}{}
			o.logicalDevices = append(o.logicalDevices, *ld)
		}
	}
	return nil
}
