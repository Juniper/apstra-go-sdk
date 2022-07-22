package goapstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

const (
	apiUrlDesignTemplates       = apiUrlDesignPrefix + "templates"
	apiUrlDesignTemplatesPrefix = apiUrlDesignTemplates + apiUrlPathDelim
	apiUrlDesignTemplateById    = apiUrlDesignTemplatesPrefix + "%s"

	templateResponseTypeField  = "type"
	msgTemplateNotImplemented  = "template type '%s' not yet implemented"
	msgTemplateUnpackWrongType = "cannot unpack template of type '%s' into '%s'"
)

type TemplateType int
type templateType string
type AsnAllocationScheme int
type asnAllocationScheme string
type AddressingScheme int
type addressingScheme string

type OverlayControlProtocol int
type overlayControlProtocol string
type TemplateCapability int
type templateCapability string

const (
	TemplateTypeRackBased = TemplateType(iota)
	TemplateTypePodBased
	TemplateTypeL3Collapsed
	TemplateTypeNone
	TemplateTypeUnknown = "unknown template type '%s'"

	templateTypeRackBased   = templateType("rack_based")
	templateTypePodBased    = templateType("pod_based")
	templateTypeL3Collapsed = templateType("l3_collapsed")
	templateTypeNone        = templateType("")
	templateTypeUnknown     = "unknown template type '%d'"
)

const (
	AsnAllocationSchemeDistinct = AsnAllocationScheme(iota)
	AsnAllocationSchemeSingle
	AsnAllocationSchemeUnknown = "unknown asn allocation scheme '%s'"

	asnAllocationSchemeDistinct = asnAllocationScheme("distinct")
	asnAllocationSchemeSingle   = asnAllocationScheme("single")
	asnAllocationUnknown        = "unknown asn allocation scheme '%d'"
)

const (
	AddressingSchemeIp4 = AddressingScheme(iota)
	AddressingSchemeIp6
	AddressingSchemeIp46
	AddressingSchemeUnknown = "unknown asn allocation scheme '%s'"

	addressingSchemeIp4     = addressingScheme("ipv4")
	addressingSchemeIp6     = addressingScheme("ipv6")
	addressingSchemeIp46    = addressingScheme("ipv4_ipv6")
	addressingSchemeUnknown = "unknown asn allocation scheme '%d'"
)

const (
	OverlayControlProtocolNone = OverlayControlProtocol(iota)
	OverlayControlProtocolEvpn
	OverlayControlProtocolUnknown = "unknown overlay control protocol '%s'"

	overlayControlProtocolNone    = overlayControlProtocol("")
	overlayControlProtocolEvpn    = overlayControlProtocol("evpn")
	overlayControlProtocolUnknown = "unknown overlay control protocol '%d'"
)

const (
	TemplateCapabilityBlueprint = TemplateCapability(iota)
	TemplateCapabilityPod
	TemplateCapabilityNone
	TemplateCapabilityUnknown = "unknown template capability '%s'"

	templateCapabilityBlueprint = templateCapability("blueprint")
	templateCapabilityPod       = templateCapability("pod")
	templateCapabilityNone      = templateCapability("")
	templateCapabilityUnknown   = "unknown template capability '%d'"
)

func (o TemplateType) Int() int {
	return int(o)
}

func (o TemplateType) String() string {
	switch o {
	case TemplateTypeRackBased:
		return string(templateTypeRackBased)
	case TemplateTypePodBased:
		return string(templateTypePodBased)
	case TemplateTypeL3Collapsed:
		return string(templateTypeL3Collapsed)
	case TemplateTypeNone:
		return string(templateTypeNone)
	default:
		return fmt.Sprintf(templateTypeUnknown, o)
	}
}

func (o TemplateType) raw() templateType {
	return templateType(o.String())
}

func (o templateType) string() string {
	return string(o)
}

func (o templateType) parse() (int, error) {
	switch o {
	case templateTypeRackBased:
		return int(TemplateTypeRackBased), nil
	case templateTypePodBased:
		return int(TemplateTypePodBased), nil
	case templateTypeL3Collapsed:
		return int(TemplateTypeL3Collapsed), nil
	case templateTypeNone:
		return int(TemplateTypeNone), nil
	default:
		return 0, fmt.Errorf(TemplateTypeUnknown, o)
	}
}

func (o AsnAllocationScheme) Int() int {
	return int(o)
}

func (o AsnAllocationScheme) String() string {
	switch o {
	case AsnAllocationSchemeDistinct:
		return string(asnAllocationSchemeDistinct)
	case AsnAllocationSchemeSingle:
		return string(asnAllocationSchemeSingle)
	default:
		return fmt.Sprintf(asnAllocationUnknown, o)
	}
}

func (o AsnAllocationScheme) raw() asnAllocationScheme {
	return asnAllocationScheme(o.String())
}

func (o asnAllocationScheme) string() string {
	return string(o)
}

func (o asnAllocationScheme) parse() (int, error) {
	switch o {
	case asnAllocationSchemeDistinct:
		return int(AsnAllocationSchemeDistinct), nil
	case asnAllocationSchemeSingle:
		return int(AsnAllocationSchemeSingle), nil
	default:
		return 0, fmt.Errorf(AsnAllocationSchemeUnknown, o)
	}
}

func (o AddressingScheme) Int() int {
	return int(o)
}

func (o AddressingScheme) String() string {
	switch o {
	case AddressingSchemeIp4:
		return string(addressingSchemeIp4)
	case AddressingSchemeIp6:
		return string(addressingSchemeIp6)
	case AddressingSchemeIp46:
		return string(addressingSchemeIp46)
	default:
		return fmt.Sprintf(addressingSchemeUnknown, o)
	}
}

func (o AddressingScheme) raw() addressingScheme {
	return addressingScheme(o.String())
}

func (o addressingScheme) string() string {
	return string(o)
}

func (o addressingScheme) parse() (int, error) {
	switch o {
	case addressingSchemeIp4:
		return int(AddressingSchemeIp4), nil
	case addressingSchemeIp6:
		return int(AddressingSchemeIp6), nil
	case addressingSchemeIp46:
		return int(AddressingSchemeIp46), nil
	default:
		return 0, fmt.Errorf(AddressingSchemeUnknown, o)
	}
}

func (o OverlayControlProtocol) Int() int {
	return int(o)
}
func (o OverlayControlProtocol) String() string {
	switch o {
	case OverlayControlProtocolNone:
		return string(overlayControlProtocolNone)
	case OverlayControlProtocolEvpn:
		return string(overlayControlProtocolEvpn)
	default:
		return fmt.Sprintf(overlayControlProtocolUnknown, o)
	}
}
func (o OverlayControlProtocol) raw() overlayControlProtocol {
	return overlayControlProtocol(o.String())
}
func (o overlayControlProtocol) string() string {
	return string(o)
}
func (o overlayControlProtocol) parse() (int, error) {
	switch o {
	case overlayControlProtocolNone:
		return int(OverlayControlProtocolNone), nil
	case overlayControlProtocolEvpn:
		return int(OverlayControlProtocolEvpn), nil
	default:
		return 0, fmt.Errorf(OverlayControlProtocolUnknown, o)
	}
}

func (o TemplateCapability) Int() int {
	return int(o)
}
func (o TemplateCapability) String() string {
	switch o {
	case TemplateCapabilityBlueprint:
		return string(templateCapabilityBlueprint)
	case TemplateCapabilityPod:
		return string(templateCapabilityPod)
	case TemplateCapabilityNone:
		return string(templateCapabilityNone)
	default:
		return fmt.Sprintf(templateCapabilityUnknown, o)
	}
}
func (o TemplateCapability) raw() templateCapability {
	return templateCapability(o.String())
}
func (o templateCapability) string() string {
	return string(o)
}
func (o templateCapability) parse() (int, error) {
	switch o {
	case templateCapabilityBlueprint:
		return int(TemplateCapabilityBlueprint), nil
	case templateCapabilityPod:
		return int(TemplateCapabilityPod), nil
	case templateCapabilityNone:
		return int(TemplateCapabilityNone), nil
	default:
		return 0, fmt.Errorf(TemplateCapabilityUnknown, o)
	}
}

type rawTemplateResponse map[string]interface{}

func (o rawTemplateResponse) getType() (templateType, bool) {
	if t, ok := o[templateResponseTypeField]; ok {
		return templateType(t.(string)), true
	}
	return "", false
}

func (o rawTemplateResponse) parse() (interface{}, error) {
	t, ok := o.getType()
	if !ok {
		return nil, fmt.Errorf("unable to determine template type")
	}

	switch t {
	case templateTypeRackBased:
		raw := &rawTemplateRackBased{}
		return raw, o.unpack(raw)
	case templateTypePodBased:
		raw := &rawTemplatePodBased{}
		return raw, o.unpack(raw)
		//return nil, fmt.Errorf(msgTemplateNotImplemented, t)
	case templateTypeL3Collapsed:
		raw := &rawTemplateL3Collapsed{}
		return raw, o.unpack(raw)
		//return nil, fmt.Errorf(msgTemplateNotImplemented, t)
	default:
		return nil, fmt.Errorf(TemplateTypeUnknown, t)
	}

}

func (o rawTemplateResponse) unpack(in interface{}) error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, in)
	return err
}

type optionsTemplatesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type getAllTemplatesResponse struct {
	Items []rawTemplateResponse `json:"items"`
}

type AntiAffinityPolicy struct {
	Algorithm                string `json:"algorithm"` // heuristic
	MaxLinksPerPort          int    `json:"max_links_per_port"`
	MaxLinksPerSlot          int    `json:"max_links_per_slot"`
	MaxPerSystemLinksPerPort int    `json:"max_per_system_links_per_port"`
	MaxPerSystemLinksPerSlot int    `json:"max_per_system_links_per_slot"`
	Mode                     string `json:"mode"` // disabled, enabled_loose, enabled_strict
}

type VirtualNetworkPolicy struct {
	OverlayControlProtocol OverlayControlProtocol
}

func (o *VirtualNetworkPolicy) raw() *rawVirtualNetworkPolicy {
	return &rawVirtualNetworkPolicy{OverlayControlProtocol: o.OverlayControlProtocol.raw()}
}

type rawVirtualNetworkPolicy struct {
	OverlayControlProtocol overlayControlProtocol `json:"overlay_control_protocol,omitempty"`
}

func (o *rawVirtualNetworkPolicy) polish() (*VirtualNetworkPolicy, error) {
	ocp, err := o.OverlayControlProtocol.parse()
	return &VirtualNetworkPolicy{OverlayControlProtocol: OverlayControlProtocol(ocp)}, err
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

type FabricAddressingPolicy struct {
	SpineSuperspineLinks AddressingScheme
	SpineLeafLinks       AddressingScheme
}

func (o *FabricAddressingPolicy) raw() *rawFabricAddressingPolicy {
	return &rawFabricAddressingPolicy{
		SpineSuperspineLinks: o.SpineSuperspineLinks.raw(),
		SpineLeafLinks:       o.SpineLeafLinks.raw(),
	}
}

type rawFabricAddressingPolicy struct {
	SpineSuperspineLinks addressingScheme `json:"spine_superspine_links"`
	SpineLeafLinks       addressingScheme `json:"spine_leaf_links"`
}

func (o *rawFabricAddressingPolicy) polish() (*FabricAddressingPolicy, error) {
	ssl, err := o.SpineSuperspineLinks.parse()
	if err != nil {
		return nil, err
	}

	sll, err := o.SpineLeafLinks.parse()
	if err != nil {
		return nil, err
	}

	return &FabricAddressingPolicy{
		SpineSuperspineLinks: AddressingScheme(ssl),
		SpineLeafLinks:       AddressingScheme(sll),
	}, nil
}

type Spine struct {
	Count                   int
	ExternalLinkSpeed       LogicalDevicePortSpeed
	LinkPerSuperspineSpeed  LogicalDevicePortSpeed
	LogicalDevice           LogicalDevice
	LinkPerSuperspineCount  int
	Tags                    []DesignTag
	ExternalLinksPerNode    int
	ExternalFacingNodeCount int
	ExternalLinkCount       int
}

func (o Spine) raw() *rawSpine {
	return &rawSpine{
		Count:                   o.Count,
		ExternalLinkSpeed:       *o.ExternalLinkSpeed.raw(),
		LinkPerSuperspineSpeed:  *o.LinkPerSuperspineSpeed.raw(),
		LogicalDevice:           o.LogicalDevice.raw(),
		LinkPerSuperspineCount:  o.LinkPerSuperspineCount,
		Tags:                    o.Tags,
		ExternalLinksPerNode:    o.ExternalLinksPerNode,
		ExternalFacingNodeCount: o.ExternalFacingNodeCount,
		ExternalLinkCount:       o.ExternalLinkCount,
	}
}

type rawSpine struct {
	Count                   int                       `json:"count"`
	ExternalLinkSpeed       rawLogicalDevicePortSpeed `json:"external_link_speed"`
	LinkPerSuperspineSpeed  rawLogicalDevicePortSpeed `json:"link_per_superspine_speed"`
	LogicalDevice           *rawLogicalDevice         `json:"logical_device"`
	LinkPerSuperspineCount  int                       `json:"link_per_superspine_count"`
	Tags                    []DesignTag               `json:"tags"`
	ExternalLinksPerNode    int                       `json:"external_links_per_node"`
	ExternalFacingNodeCount int                       `json:"external_facing_node_count"`
	ExternalLinkCount       int                       `json:"external_link_count"`
}

func (o rawSpine) polish() (*Spine, error) {
	ld, err := o.LogicalDevice.polish()
	return &Spine{
		Count:                   o.Count,
		ExternalLinkSpeed:       o.ExternalLinkSpeed.parse(),
		LinkPerSuperspineSpeed:  o.LinkPerSuperspineSpeed.parse(),
		LogicalDevice:           *ld,
		LinkPerSuperspineCount:  o.LinkPerSuperspineCount,
		Tags:                    o.Tags,
		ExternalLinksPerNode:    o.ExternalLinksPerNode,
		ExternalFacingNodeCount: o.ExternalFacingNodeCount,
		ExternalLinkCount:       o.ExternalLinkCount,
	}, err
}

type TemplateRackBased struct {
	Id                     ObjectId               `json:"id"`
	Type                   TemplateType           `json:"type"`
	DisplayName            string                 `json:"display_name"`
	Status                 string                 `json:"status,omitempty"` // inconsistent, ok
	AntiAffinityPolicy     AntiAffinityPolicy     `json:"anti_affinity_policy"`
	CreatedAt              time.Time              `json:"created_at"`
	LastModifiedAt         time.Time              `json:"last_modified_at"`
	VirtualNetworkPolicy   VirtualNetworkPolicy   `json:"virtual_network_policy"`
	AsnAllocationPolicy    AsnAllocationPolicy    `json:"asn_allocation_policy"`
	FabricAddressingPolicy FabricAddressingPolicy `json:"fabric_addressing_policy"`
	Capability             TemplateCapability     `json:"capability"`
	Spine                  Spine                  `json:"spine"`
	RackTypes              []RackType             `json:"rack_types"`
	RackTypeCounts         []struct {
		RackTypeId ObjectId `json:"rack_type_id"`
		Count      int      `json:"count"`
	} `json:"rack_type_counts"`
	DhcpServiceIntent struct {
		Active bool `json:"active"`
	} `json:"dhcp_service_intent"`
}

//func (o TemplateRackBased) raw() *rawTemplateRackBased {
//	var rrt []rawRackType
//	for _, rt := range o.RackTypes {
//		rrt = append(rrt, *rt.raw())
//	}
//	return &rawTemplateRackBased{
//		Id:                     o.Id,
//		Type:                   o.Type.raw(),
//		DisplayName:            o.DisplayName,
//		Status:                 o.Status,
//		AntiAffinityPolicy:     o.AntiAffinityPolicy,
//		CreatedAt:              o.CreatedAt,
//		LastModifiedAt:         o.LastModifiedAt,
//		VirtualNetworkPolicy:   *o.VirtualNetworkPolicy.raw(),
//		AsnAllocationPolicy:    *o.AsnAllocationPolicy.raw(),
//		FabricAddressingPolicy: *o.FabricAddressingPolicy.raw(),
//		Spine:                  *o.Spine.raw(),
//		RackTypes:              rrt,
//		RackTypeCounts:         o.RackTypeCounts,
//		DhcpServiceIntent:      o.DhcpServiceIntent,
//		Capability:             o.Capability.raw(),
//	}
//}

type rawTemplateRackBased struct {
	Id                     ObjectId                  `json:"id"`
	Type                   templateType              `json:"type"`
	DisplayName            string                    `json:"display_name"`
	Status                 string                    `json:"status,omitempty"` // inconsistent, ok
	AntiAffinityPolicy     AntiAffinityPolicy        `json:"anti_affinity_policy"`
	CreatedAt              time.Time                 `json:"created_at"`
	LastModifiedAt         time.Time                 `json:"last_modified_at"`
	VirtualNetworkPolicy   rawVirtualNetworkPolicy   `json:"virtual_network_policy"`
	AsnAllocationPolicy    rawAsnAllocationPolicy    `json:"asn_allocation_policy"`
	FabricAddressingPolicy rawFabricAddressingPolicy `json:"fabric_addressing_policy"`
	Capability             templateCapability        `json:"capability,omitempty"`
	Spine                  rawSpine                  `json:"spine"`
	RackTypes              []rawRackType             `json:"rack_types"`
	RackTypeCounts         []struct {
		RackTypeId ObjectId `json:"rack_type_id"`
		Count      int      `json:"count"`
	} `json:"rack_type_counts"`
	DhcpServiceIntent struct {
		Active bool `json:"active"`
	} `json:"dhcp_service_intent"`
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
	f, err := o.FabricAddressingPolicy.polish()
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
	var rackTypes []RackType
	for _, rt := range o.RackTypes {
		prt, err := rt.polish()
		if err != nil {
			return nil, err
		}
		rackTypes = append(rackTypes, *prt)
	}
	return &TemplateRackBased{
		Id:                     o.Id,
		Type:                   TemplateType(tType),
		DisplayName:            o.DisplayName,
		Status:                 o.Status,
		AntiAffinityPolicy:     o.AntiAffinityPolicy,
		CreatedAt:              o.CreatedAt,
		LastModifiedAt:         o.LastModifiedAt,
		VirtualNetworkPolicy:   *v,
		AsnAllocationPolicy:    *a,
		FabricAddressingPolicy: *f,
		Capability:             TemplateCapability(c),
		Spine:                  *s,
		RackTypes:              rackTypes,
		RackTypeCounts:         o.RackTypeCounts,
		DhcpServiceIntent:      o.DhcpServiceIntent,
	}, nil
}

type Superspine struct {
	PlaneCount         int
	ExternalLinkCount  int
	ExternalLinkSpeed  rawLogicalDevicePortSpeed
	Tags               []DesignTag
	SuperspinePerPlane int
	LogicalDevice      LogicalDevice
}

func (o Superspine) raw() *rawSuperspine {
	return &rawSuperspine{
		PlaneCount:         o.PlaneCount,
		ExternalLinkCount:  o.ExternalLinkCount,
		ExternalLinkSpeed:  o.ExternalLinkSpeed,
		Tags:               o.Tags,
		SuperspinePerPlane: o.SuperspinePerPlane,
		LogicalDevice:      *o.LogicalDevice.raw(),
	}
}

type rawSuperspine struct {
	PlaneCount         int                       `json:"plane_count"`
	ExternalLinkCount  int                       `json:"external_link_count"`
	ExternalLinkSpeed  rawLogicalDevicePortSpeed `json:"external_link_speed"`
	Tags               []DesignTag               `json:"tags"`
	SuperspinePerPlane int                       `json:"superspine_per_plane"`
	LogicalDevice      rawLogicalDevice          `json:"logical_device"`
}

func (o rawSuperspine) polish() (*Superspine, error) {
	ld, err := o.LogicalDevice.polish()
	if err != nil {
		return nil, err
	}
	return &Superspine{
		PlaneCount:         o.PlaneCount,
		ExternalLinkCount:  o.ExternalLinkCount,
		ExternalLinkSpeed:  o.ExternalLinkSpeed,
		Tags:               o.Tags,
		SuperspinePerPlane: o.SuperspinePerPlane,
		LogicalDevice:      *ld,
	}, nil
}

type RackBasedTemplateCount struct {
	RackBasedTemplateId ObjectId `json:"rack_based_template_id"`
	Count               int      `json:"count"`
}

type TemplatePodBased struct {
	Id                      ObjectId
	Type                    string
	Status                  string
	DisplayName             string
	AntiAffinityPolicy      AntiAffinityPolicy
	FabricAddressingPolicy  FabricAddressingPolicy
	Superspine              Superspine
	CreatedAt               time.Time
	LastModifiedAt          time.Time
	Capability              TemplateCapability
	RackBasedTemplates      []TemplateRackBased
	RackBasedTemplateCounts []RackBasedTemplateCount
}

//func (o TemplatePodBased) raw() *rawTemplatePodBased { // todo
//	var rrbt []rawTemplateRackBased
//	for _, rbt := range o.RackBasedTemplates {
//		rrbt = append(rrbt, *rbt.raw())
//	}
//	return &rawTemplatePodBased{
//		Id:                      o.Id,
//		Type:                    o.Type,
//		Status:                  o.Status,
//		DisplayName:             o.DisplayName,
//		AntiAffinityPolicy:      o.AntiAffinityPolicy,
//		FabricAddressingPolicy:  *o.FabricAddressingPolicy.raw(),
//		Superspine:              *o.Superspine.raw(),
//		CreatedAt:               o.CreatedAt,
//		LastModifiedAt:          o.LastModifiedAt,
//		Capability:              o.Capability.raw(),
//		RackBasedTemplates:      rrbt,
//		RackBasedTemplateCounts: o.RackBasedTemplateCounts,
//	}
//}

type rawTemplatePodBased struct {
	Id                      ObjectId                  `json:"id"`
	Type                    string                    `json:"type"`
	Status                  string                    `json:"status"`
	DisplayName             string                    `json:"display_name"`
	AntiAffinityPolicy      AntiAffinityPolicy        `json:"anti_affinity_policy"`
	FabricAddressingPolicy  rawFabricAddressingPolicy `json:"fabric_addressing_policy"`
	Superspine              rawSuperspine             `json:"superspine"`
	CreatedAt               time.Time                 `json:"created_at"`
	LastModifiedAt          time.Time                 `json:"last_modified_at"`
	Capability              templateCapability        `json:"capability,omitempty"`
	RackBasedTemplates      []rawTemplateRackBased    `json:"rack_based_templates"`
	RackBasedTemplateCounts []RackBasedTemplateCount  `json:"rack_based_template_counts"`
}

func (o rawTemplatePodBased) polish() (*TemplatePodBased, error) {
	fap, err := o.FabricAddressingPolicy.polish()
	if err != nil {
		return nil, err
	}
	superspine, err := o.Superspine.polish()
	if err != nil {
		return nil, err
	}
	capability, err := o.Capability.parse()
	if err != nil {
		return nil, err
	}
	var _, rbt []TemplateRackBased
	for _, rrbt := range o.RackBasedTemplates {
		if rrbt.Type == templateTypeNone {
			// because sometimes Apstra doesn't fill this in, but we know based on context
			rrbt.Type = templateTypeRackBased
		}
		polished, err := rrbt.polish()
		if err != nil {
			return nil, err
		}
		rbt = append(rbt, *polished)
	}
	return &TemplatePodBased{
		Id:                      o.Id,
		Type:                    o.Type,
		Status:                  o.Status,
		DisplayName:             o.DisplayName,
		AntiAffinityPolicy:      o.AntiAffinityPolicy,
		FabricAddressingPolicy:  *fap,
		Superspine:              *superspine,
		CreatedAt:               o.CreatedAt,
		LastModifiedAt:          o.LastModifiedAt,
		Capability:              TemplateCapability(capability),
		RackBasedTemplates:      rbt,
		RackBasedTemplateCounts: o.RackBasedTemplateCounts,
	}, nil
}

type TemplateL3Collapsed struct {
	Id                   ObjectId                  `json:"id"`
	Type                 string                    `json:"type"`
	Status               string                    `json:"status"`
	DisplayName          string                    `json:"display_name"`
	AntiAffinityPolicy   AntiAffinityPolicy        `json:"anti_affinity_policy"`
	CreatedAt            time.Time                 `json:"created_at"`
	LastModifiedAt       time.Time                 `json:"last_modified_at"`
	RackTypes            []RackType                `json:"rack_types"`
	Capability           TemplateCapability        `json:"capability"`
	MeshLinkSpeed        rawLogicalDevicePortSpeed `json:"mesh_link_speed"`
	VirtualNetworkPolicy VirtualNetworkPolicy      `json:"virtual_network_policy"`
	MeshLinkCount        int                       `json:"mesh_link_count"`
	RackTypeCounts       []struct {
		RackTypeId ObjectId `json:"rack_type_id"`
		Count      int      `json:"count"`
	} `json:"rack_type_counts"`
	DhcpServiceIntent struct {
		Active bool `json:"active"`
	} `json:"dhcp_service_intent"`
}

//func (o TemplateL3Collapsed) raw() *rawTemplateL3Collapsed {
//	var rrt []rawRackType
//	for _, rt := range o.RackTypes {
//		rrt = append(rrt, *rt.raw())
//	}
//	return &rawTemplateL3Collapsed{
//		Id:                   o.Id,
//		Type:                 o.Type,
//		Status:               o.Status,
//		DisplayName:          o.DisplayName,
//		AntiAffinityPolicy:   o.AntiAffinityPolicy,
//		CreatedAt:            o.CreatedAt,
//		LastModifiedAt:       o.LastModifiedAt,
//		RackTypes:            rrt,
//		Capability:           o.Capability.raw(),
//		MeshLinkSpeed:        o.MeshLinkSpeed,
//		VirtualNetworkPolicy: *o.VirtualNetworkPolicy.raw(),
//		MeshLinkCount:        o.MeshLinkCount,
//		RackTypeCounts:       o.RackTypeCounts,
//		DhcpServiceIntent:    o.DhcpServiceIntent,
//	}
//}

type rawTemplateL3Collapsed struct {
	Id                   ObjectId                  `json:"id"`
	Type                 string                    `json:"type"`
	Status               string                    `json:"status"`
	DisplayName          string                    `json:"display_name"`
	AntiAffinityPolicy   AntiAffinityPolicy        `json:"anti_affinity_policy"`
	CreatedAt            time.Time                 `json:"created_at"`
	LastModifiedAt       time.Time                 `json:"last_modified_at"`
	RackTypes            []rawRackType             `json:"rack_types"`
	Capability           templateCapability        `json:"capability"`
	MeshLinkSpeed        rawLogicalDevicePortSpeed `json:"mesh_link_speed"`
	VirtualNetworkPolicy rawVirtualNetworkPolicy   `json:"virtual_network_policy"`
	MeshLinkCount        int                       `json:"mesh_link_count"`
	RackTypeCounts       []struct {
		RackTypeId ObjectId `json:"rack_type_id"`
		Count      int      `json:"count"`
	} `json:"rack_type_counts"`
	DhcpServiceIntent struct {
		Active bool `json:"active"`
	} `json:"dhcp_service_intent"`
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
	return &TemplateL3Collapsed{
		Id:                   o.Id,
		Type:                 o.Type,
		Status:               o.Status,
		DisplayName:          o.DisplayName,
		AntiAffinityPolicy:   o.AntiAffinityPolicy,
		CreatedAt:            o.CreatedAt,
		LastModifiedAt:       o.LastModifiedAt,
		RackTypes:            prt,
		Capability:           TemplateCapability(capability),
		MeshLinkSpeed:        o.MeshLinkSpeed,
		VirtualNetworkPolicy: *vnp,
		MeshLinkCount:        o.MeshLinkCount,
		RackTypeCounts:       o.RackTypeCounts,
		DhcpServiceIntent:    o.DhcpServiceIntent,
	}, nil
}

func polishAnyTemplate(in interface{}) (interface{}, error) {
	switch in.(type) {
	case *rawTemplateRackBased:
		return in.(*rawTemplateRackBased).polish()
	case *rawTemplatePodBased:
		polished, err := in.(*rawTemplatePodBased).polish()
		return polished, err
	case *rawTemplateL3Collapsed:
		return in.(*rawTemplateL3Collapsed).polish()
	default:
		return nil, fmt.Errorf("unknown raw template type '%s'", reflect.TypeOf(in))
	}
}

func (o *Client) listAllTemplateIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsTemplatesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (o *Client) getTemplate(ctx context.Context, id ObjectId) (interface{}, error) {
	rtr := &rawTemplateResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignTemplateById, id),
		apiResponse: rtr,
	})
	if err != nil {
		return nil, err
	}

	raw, err := rtr.parse()
	if err != nil {
		return nil, err
	}

	return polishAnyTemplate(raw)
}

func (o *Client) getAllTemplates(ctx context.Context) (map[TemplateType][]interface{}, error) {
	response := &getAllTemplatesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}

	var rackBasedTemplates []interface{}
	var podBasedTemplates []interface{}
	var l3CollapsedTemplates []interface{}

	for _, rtr := range response.Items {
		var tType interface{}
		var ok bool
		if tType, ok = rtr[templateResponseTypeField]; !ok {
			return nil, fmt.Errorf("no template type in response")
		}
		switch tType {
		case string(templateTypeRackBased):
		case string(templateTypePodBased):
		case string(templateTypeL3Collapsed):
		default:
			return nil, fmt.Errorf("unknown template type '%s'", tType)
		}

		raw, err := rtr.parse()
		if err != nil {
			return nil, err
		}

		polished, err := polishAnyTemplate(raw)
		if err != nil {
			return nil, err
		}
		switch polished.(type) {
		case *TemplateRackBased:
			rackBasedTemplates = append(rackBasedTemplates, *polished.(*TemplateRackBased))
		case *TemplatePodBased:
			podBasedTemplates = append(podBasedTemplates, *polished.(*TemplatePodBased))
		case *TemplateL3Collapsed:
			l3CollapsedTemplates = append(l3CollapsedTemplates, *polished.(*TemplateL3Collapsed))
		default:
			return nil, fmt.Errorf("unknown template type '%s'", reflect.TypeOf(polished))
		}

	}

	return map[TemplateType][]interface{}{
		TemplateTypeRackBased:   rackBasedTemplates,
		TemplateTypePodBased:    podBasedTemplates,
		TemplateTypeL3Collapsed: l3CollapsedTemplates,
	}, nil
}

func (o *Client) getRackBasedTemplate(ctx context.Context, id ObjectId) (*TemplateRackBased, error) {
	template, err := o.getTemplate(ctx, id)
	return template.(*TemplateRackBased), err
}

func (o *Client) getAllRackBasedTemplates(ctx context.Context) ([]TemplateRackBased, error) {
	tMap, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]TemplateRackBased, len(tMap[TemplateTypeRackBased]))
	for i, t := range tMap[TemplateTypeRackBased] {
		result[i] = t.(TemplateRackBased)
	}

	return result, nil
}

func (o *Client) getPodBasedTemplate(ctx context.Context, id ObjectId) (*TemplatePodBased, error) {
	template, err := o.getTemplate(ctx, id)
	return template.(*TemplatePodBased), err
}

func (o *Client) getAllPodBasedTemplates(ctx context.Context) ([]TemplatePodBased, error) {
	tMap, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]TemplatePodBased, len(tMap[TemplateTypePodBased]))
	for i, t := range tMap[TemplateTypePodBased] {
		result[i] = t.(TemplatePodBased)
	}

	return result, nil
}

func (o *Client) getL3CollapsedTemplate(ctx context.Context, id ObjectId) (*TemplateL3Collapsed, error) {
	template, err := o.getTemplate(ctx, id)
	return template.(*TemplateL3Collapsed), err
}

func (o *Client) getAllL3CollapsedTemplates(ctx context.Context) ([]TemplateL3Collapsed, error) {
	tMap, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]TemplateL3Collapsed, len(tMap[TemplateTypeL3Collapsed]))
	for i, t := range tMap[TemplateTypeL3Collapsed] {
		result[i] = t.(TemplateL3Collapsed)
	}

	return result, nil
}

func (o *Client) getTemplateAndType(ctx context.Context, id ObjectId) (TemplateType, interface{}, error) {
	template, err := o.getTemplate(ctx, id)
	if err != nil {
		return 0, nil, err
	}
	switch template.(type) {
	case *TemplateRackBased:
		return TemplateTypeRackBased, template, nil
	case *TemplatePodBased:
		return TemplateTypePodBased, template, nil
	case *TemplateL3Collapsed:
		return TemplateTypeL3Collapsed, template, nil
	default:
		return 0, template, errors.New("unknown template type at getTemplateAndType()")
	}
}
