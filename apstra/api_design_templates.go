package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDesignTemplates       = apiUrlDesignPrefix + "templates"
	apiUrlDesignTemplatesPrefix = apiUrlDesignTemplates + apiUrlPathDelim
	apiUrlDesignTemplateById    = apiUrlDesignTemplatesPrefix + "%s"
)

type AntiAffninityAlgorithm int
type antiAffinityAlgorithm string
type AntiAffinityMode int
type antiAffinityMode string
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
	AntiAffinityModeDisabled = AntiAffinityMode(iota)
	AntiAffinityModeLoose
	AntiAffinityModeStrict
	AntiAffinityModeUnknown = "unknown anti affinity mode %s"

	antiAffinityModeDisabled = antiAffinityMode("disabled")
	antiAffinityModeLoose    = antiAffinityMode("loose")
	antiAffinityModeStrict   = antiAffinityMode("strict")
	antiAffinityModeUnknown  = "unknown anti affinity mode %d"
)

const (
	AlgorithmHeuristic = AntiAffninityAlgorithm(iota)
	AlgorithmUnknown   = "unknown anti affinity algorithm '%s'"

	algorithmHeuristic = antiAffinityAlgorithm("heuristic")
	algorithmUnknown   = "unknown anti affinity algorithm '%d'"
)

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
	AddressingSchemeUnknown = "unknown addressing scheme '%s'"

	addressingSchemeIp4     = addressingScheme("ipv4")
	addressingSchemeIp6     = addressingScheme("ipv6")
	addressingSchemeIp46    = addressingScheme("ipv4_ipv6")
	addressingSchemeUnknown = "unknown addressing scheme '%d'"
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

func (o AntiAffinityMode) Int() int {
	return int(o)
}

func (o AntiAffinityMode) String() string {
	switch o {
	case AntiAffinityModeDisabled:
		return string(antiAffinityModeDisabled)
	case AntiAffinityModeLoose:
		return string(antiAffinityModeLoose)
	case AntiAffinityModeStrict:
		return string(antiAffinityModeStrict)
	default:
		return fmt.Sprintf(antiAffinityModeUnknown, o)
	}
}

func (o AntiAffinityMode) raw() antiAffinityMode {
	return antiAffinityMode(o.String())
}

func (o antiAffinityMode) string() string {
	return string(o)
}

func (o antiAffinityMode) parse() (int, error) {
	switch o {
	case antiAffinityModeDisabled:
		return int(AntiAffinityModeDisabled), nil
	case antiAffinityModeLoose:
		return int(AntiAffinityModeLoose), nil
	case antiAffinityModeStrict:
		return int(AntiAffinityModeStrict), nil
	default:
		return 0, fmt.Errorf(AntiAffinityModeUnknown, o)
	}
}

func (o AntiAffninityAlgorithm) Int() int {
	return int(o)
}

func (o AntiAffninityAlgorithm) String() string {
	switch o {
	case AlgorithmHeuristic:
		return string(algorithmHeuristic)
	default:
		return fmt.Sprintf(algorithmUnknown, o)
	}
}

func (o AntiAffninityAlgorithm) raw() antiAffinityAlgorithm {
	return antiAffinityAlgorithm(o.String())
}

func (o antiAffinityAlgorithm) string() string {
	return string(o)
}

func (o antiAffinityAlgorithm) parse() (int, error) {
	switch o {
	case algorithmHeuristic:
		return int(AlgorithmHeuristic), nil
	default:
		return 0, fmt.Errorf(AlgorithmUnknown, o)
	}
}

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

func (o *TemplateType) FromString(s string) error {
	i, err := templateType(s).parse()
	if err != nil {
		return err
	}
	*o = TemplateType(i)
	return nil
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

func (o *AsnAllocationScheme) FromString(in string) error {
	i, err := asnAllocationScheme(in).parse()
	if err != nil {
		return err
	}
	*o = AsnAllocationScheme(i)
	return nil
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

func (o *AddressingScheme) FromString(in string) error {
	i, err := addressingScheme(in).parse()
	if err != nil {
		return err
	}
	*o = AddressingScheme(i)
	return nil
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

func (o *OverlayControlProtocol) FromString(in string) error {
	i, err := overlayControlProtocol(in).parse()
	if err != nil {
		return err
	}
	*o = OverlayControlProtocol(i)
	return nil
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

type AntiAffinityPolicy struct {
	Algorithm                AntiAffninityAlgorithm
	MaxLinksPerPort          int
	MaxLinksPerSlot          int
	MaxPerSystemLinksPerPort int
	MaxPerSystemLinksPerSlot int
	Mode                     AntiAffinityMode
}

func (o *AntiAffinityPolicy) raw() *rawAntiAffinityPolicy {
	return &rawAntiAffinityPolicy{
		Algorithm:                o.Algorithm.raw(),
		MaxLinksPerPort:          o.MaxLinksPerPort,
		MaxLinksPerSlot:          o.MaxLinksPerSlot,
		MaxPerSystemLinksPerPort: o.MaxPerSystemLinksPerPort,
		MaxPerSystemLinksPerSlot: o.MaxPerSystemLinksPerSlot,
		Mode:                     o.Mode.raw(),
	}
}

type rawAntiAffinityPolicy struct {
	Algorithm                antiAffinityAlgorithm `json:"algorithm"` // heuristic
	MaxLinksPerPort          int                   `json:"max_links_per_port"`
	MaxLinksPerSlot          int                   `json:"max_links_per_slot"`
	MaxPerSystemLinksPerPort int                   `json:"max_per_system_links_per_port"`
	MaxPerSystemLinksPerSlot int                   `json:"max_per_system_links_per_slot"`
	Mode                     antiAffinityMode      `json:"mode"` // disabled, enabled_loose, enabled_strict
}

func (o *rawAntiAffinityPolicy) polish() (*AntiAffinityPolicy, error) {
	algorithm, err := o.Algorithm.parse()
	if err != nil {
		return nil, err
	}
	mode, err := o.Mode.parse()
	if err != nil {
		return nil, err
	}
	return &AntiAffinityPolicy{
		Algorithm:                AntiAffninityAlgorithm(algorithm),
		MaxLinksPerPort:          o.MaxLinksPerPort,
		MaxLinksPerSlot:          o.MaxLinksPerSlot,
		MaxPerSystemLinksPerPort: o.MaxPerSystemLinksPerPort,
		MaxPerSystemLinksPerSlot: o.MaxPerSystemLinksPerSlot,
		Mode:                     AntiAffinityMode(mode),
	}, nil

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

type RackTypeCount struct {
	RackTypeId ObjectId `json:"rack_type_id"`
	Count      int      `json:"count"`
}

type RackBasedTemplateCount struct {
	RackBasedTemplateId ObjectId `json:"rack_based_template_id"`
	Count               int      `json:"count"`
}

type Spine struct {
	Count                   int
	ExternalLinkSpeed       LogicalDevicePortSpeed
	LinkPerSuperspineSpeed  LogicalDevicePortSpeed
	LogicalDevice           LogicalDeviceData
	LinkPerSuperspineCount  int
	Tags                    []DesignTagData
	ExternalLinksPerNode    int
	ExternalFacingNodeCount int
	ExternalLinkCount       int
}

type TemplateElementSpineRequest struct {
	Count                   int
	ExternalLinkSpeed       LogicalDevicePortSpeed
	LinkPerSuperspineSpeed  LogicalDevicePortSpeed
	LogicalDevice           ObjectId
	LinkPerSuperspineCount  int
	Tags                    []ObjectId
	ExternalLinksPerNode    int
	ExternalFacingNodeCount int
	ExternalLinkCount       int
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
		Count:                   o.Count,
		ExternalLinkSpeed:       o.ExternalLinkSpeed.raw(),
		LinkPerSuperspineSpeed:  o.LinkPerSuperspineSpeed.raw(),
		LogicalDevice:           *logicalDevice,
		LinkPerSuperspineCount:  o.LinkPerSuperspineCount,
		Tags:                    tags,
		ExternalLinksPerNode:    o.ExternalLinksPerNode,
		ExternalFacingNodeCount: o.ExternalFacingNodeCount,
		ExternalLinkCount:       o.ExternalLinkCount,
	}, nil
}

type rawSpine struct {
	Count                   int                        `json:"count"`
	ExternalLinkSpeed       *rawLogicalDevicePortSpeed `json:"external_link_speed,omitempty"`
	LinkPerSuperspineSpeed  *rawLogicalDevicePortSpeed `json:"link_per_superspine_speed"`
	LogicalDevice           rawLogicalDevice           `json:"logical_device"`
	LinkPerSuperspineCount  int                        `json:"link_per_superspine_count"`
	Tags                    []DesignTagData            `json:"tags"`
	ExternalLinksPerNode    int                        `json:"external_links_per_node,omitempty"`
	ExternalFacingNodeCount int                        `json:"external_facing_node_count,omitempty"`
	ExternalLinkCount       int                        `json:"external_link_count,omitempty"`
}

func (o rawSpine) polish() (*Spine, error) {
	ld, err := o.LogicalDevice.polish()

	var externalLinkSpeed LogicalDevicePortSpeed
	if o.ExternalLinkSpeed != nil {
		externalLinkSpeed = o.ExternalLinkSpeed.parse()
	}

	var linkPerSuperspineSpeed LogicalDevicePortSpeed
	if o.LinkPerSuperspineSpeed != nil {
		linkPerSuperspineSpeed = o.LinkPerSuperspineSpeed.parse()
	}

	return &Spine{
		Count:                  o.Count,
		ExternalLinkSpeed:      externalLinkSpeed,
		LinkPerSuperspineSpeed: linkPerSuperspineSpeed,
		LogicalDevice: LogicalDeviceData{
			DisplayName: ld.Data.DisplayName,
			Panels:      ld.Data.Panels,
		},
		LinkPerSuperspineCount:  o.LinkPerSuperspineCount,
		Tags:                    o.Tags,
		ExternalLinksPerNode:    o.ExternalLinksPerNode,
		ExternalFacingNodeCount: o.ExternalFacingNodeCount,
		ExternalLinkCount:       o.ExternalLinkCount,
	}, err
}

type Superspine struct {
	PlaneCount         int
	ExternalLinkCount  int
	ExternalLinkSpeed  LogicalDevicePortSpeed
	Tags               []DesignTagData
	SuperspinePerPlane int
	LogicalDevice      LogicalDeviceData
}

type TemplateElementSuperspineRequest struct {
	PlaneCount         int
	ExternalLinkCount  int
	ExternalLinkSpeed  LogicalDevicePortSpeed
	Tags               []ObjectId
	SuperspinePerPlane int
	LogicalDeviceId    ObjectId
}

func (o *TemplateElementSuperspineRequest) raw(ctx context.Context, client *Client) (*rawSuperspine, error) {
	tags := make([]DesignTagData, len(o.Tags))
	for i, tagId := range o.Tags {
		rawTag, err := client.getTag(ctx, tagId)
		if err != nil {
			return nil, err
		}
		tags[i] = *rawTag.polish().Data
	}

	logicalDevice, err := client.getLogicalDevice(ctx, o.LogicalDeviceId)
	if err != nil {
		return nil, err
	}

	return &rawSuperspine{
		PlaneCount:         o.PlaneCount,
		ExternalLinkCount:  o.ExternalLinkCount,
		ExternalLinkSpeed:  o.ExternalLinkSpeed.raw(),
		Tags:               tags,
		SuperspinePerPlane: o.SuperspinePerPlane,
		LogicalDevice:      *logicalDevice,
	}, nil
}

type rawSuperspine struct {
	PlaneCount         int                        `json:"plane_count"`
	ExternalLinkCount  int                        `json:"external_link_count"`
	ExternalLinkSpeed  *rawLogicalDevicePortSpeed `json:"external_link_speed"`
	Tags               []DesignTagData            `json:"tags"`
	SuperspinePerPlane int                        `json:"superspine_per_plane"`
	LogicalDevice      rawLogicalDevice           `json:"logical_device"`
}

func (o rawSuperspine) polish() (*Superspine, error) {
	ld, err := o.LogicalDevice.polish()
	if err != nil {
		return nil, err
	}
	var externalLinkSpeed LogicalDevicePortSpeed
	if o.ExternalLinkSpeed != nil {
		externalLinkSpeed = o.ExternalLinkSpeed.parse()
	}
	return &Superspine{
		PlaneCount:         o.PlaneCount,
		ExternalLinkCount:  o.ExternalLinkCount,
		ExternalLinkSpeed:  externalLinkSpeed,
		Tags:               o.Tags,
		SuperspinePerPlane: o.SuperspinePerPlane,
		LogicalDevice: LogicalDeviceData{
			DisplayName: ld.Data.DisplayName,
			Panels:      ld.Data.Panels,
		},
	}, nil
}

type Template interface {
	Type() TemplateType
	ID() ObjectId
	OverlayControlProtocol() OverlayControlProtocol
}

type template json.RawMessage

func (o *template) templateType() (templateType, error) {
	templateProto := &struct {
		Type templateType `json:"type"`
	}{}
	return templateProto.Type, json.Unmarshal(*o, templateProto)
}

func (o *template) displayName() (string, error) {
	templateProto := &struct {
		DisplayName string `json:"display_name"`
	}{}
	return templateProto.DisplayName, json.Unmarshal(*o, templateProto)
}

func (o *template) polish() (Template, error) {
	t, err := o.templateType()
	if err != nil {
		return nil, err
	}

	// quick unmarshal just to get the type
	switch t {
	case templateTypeRackBased:
		var t rawTemplateRackBased
		err = json.Unmarshal(*o, &t)
		if err != nil {
			return nil, err
		}
		return t.polish()
	case templateTypePodBased:
		var t rawTemplatePodBased
		err = json.Unmarshal(*o, &t)
		if err != nil {
			return nil, err
		}
		return t.polish()
	case templateTypeL3Collapsed:
		var t rawTemplateL3Collapsed
		err = json.Unmarshal(*o, &t)
		if err != nil {
			return nil, err
		}
		return t.polish()
	}
	return nil, fmt.Errorf(TemplateTypeUnknown, t)
}

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

type TemplateRackBasedData struct {
	DisplayName            string
	AntiAffinityPolicy     *AntiAffinityPolicy
	VirtualNetworkPolicy   VirtualNetworkPolicy
	AsnAllocationPolicy    AsnAllocationPolicy
	FabricAddressingPolicy *FabricAddressingPolicy
	Capability             TemplateCapability
	Spine                  Spine
	RackInfo               map[ObjectId]TemplateRackBasedRackInfo
	DhcpServiceIntent      DhcpServiceIntent
}

type TemplateRackBasedRackInfo struct {
	Count        int
	RackTypeData *RackTypeData
}

type DhcpServiceIntent struct {
	Active bool `json:"active"`
}

type rawTemplateRackBased struct {
	Id                     ObjectId                   `json:"id"`
	Type                   templateType               `json:"type"`
	DisplayName            string                     `json:"display_name"`
	AntiAffinityPolicy     *rawAntiAffinityPolicy     `json:"anti_affinity_policy,omitempty"`
	CreatedAt              time.Time                  `json:"created_at"`
	LastModifiedAt         time.Time                  `json:"last_modified_at"`
	VirtualNetworkPolicy   rawVirtualNetworkPolicy    `json:"virtual_network_policy"`
	AsnAllocationPolicy    rawAsnAllocationPolicy     `json:"asn_allocation_policy"`
	FabricAddressingPolicy *rawFabricAddressingPolicy `json:"fabric_addressing_policy,omitempty"`
	Capability             templateCapability         `json:"capability,omitempty"`
	Spine                  rawSpine                   `json:"spine"`
	RackTypes              []rawRackType              `json:"rack_types"`
	RackTypeCounts         []RackTypeCount            `json:"rack_type_counts"`
	DhcpServiceIntent      DhcpServiceIntent          `json:"dhcp_service_intent"`
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
	var f *FabricAddressingPolicy
	if o.FabricAddressingPolicy != nil {
		f, err = o.FabricAddressingPolicy.polish()
		if err != nil {
			return nil, err
		}
	}
	c, err := o.Capability.parse()
	if err != nil {
		return nil, err
	}
	s, err := o.Spine.polish()
	if err != nil {
		return nil, err
	}
	var aa *AntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		aa, err = o.AntiAffinityPolicy.polish()
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
			DisplayName:            o.DisplayName,
			AntiAffinityPolicy:     aa,
			VirtualNetworkPolicy:   *v,
			AsnAllocationPolicy:    *a,
			FabricAddressingPolicy: f,
			Capability:             TemplateCapability(c),
			Spine:                  *s,
			RackInfo:               rackTypeInfos,
			DhcpServiceIntent:      o.DhcpServiceIntent,
		},
	}, nil
}

var _ Template = &TemplatePodBased{}

type TemplatePodBased struct {
	Id             ObjectId
	templateType   TemplateType
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *TemplatePodBasedData
}

func (o *TemplatePodBased) Type() TemplateType {
	return o.templateType
}

func (o *TemplatePodBased) ID() ObjectId {
	return o.Id
}

func (o *TemplatePodBased) OverlayControlProtocol() OverlayControlProtocol {
	if o == nil || o.Data == nil || len(o.Data.RackBasedTemplates) == 0 || o.Data.RackBasedTemplates[0].Data == nil {
		return OverlayControlProtocolNone
	}
	return o.Data.RackBasedTemplates[0].Data.VirtualNetworkPolicy.OverlayControlProtocol
}

type TemplatePodBasedData struct {
	DisplayName             string
	AntiAffinityPolicy      AntiAffinityPolicy
	FabricAddressingPolicy  *FabricAddressingPolicy
	Superspine              Superspine
	Capability              TemplateCapability
	RackBasedTemplates      []TemplateRackBased
	RackBasedTemplateCounts []RackBasedTemplateCount
}

type rawTemplatePodBased struct {
	Id                      ObjectId                   `json:"id"`
	Type                    templateType               `json:"type"`
	DisplayName             string                     `json:"display_name"`
	AntiAffinityPolicy      *rawAntiAffinityPolicy     `json:"anti_affinity_policy,omitempty"`
	FabricAddressingPolicy  *rawFabricAddressingPolicy `json:"fabric_addressing_policy,omitempty"`
	Superspine              rawSuperspine              `json:"superspine"`
	CreatedAt               time.Time                  `json:"created_at"`
	LastModifiedAt          time.Time                  `json:"last_modified_at"`
	Capability              templateCapability         `json:"capability,omitempty"`
	RackBasedTemplates      []rawTemplateRackBased     `json:"rack_based_templates"`
	RackBasedTemplateCounts []RackBasedTemplateCount   `json:"rack_based_template_counts"`
}

func (o rawTemplatePodBased) polish() (*TemplatePodBased, error) {
	var err error
	var fap *FabricAddressingPolicy
	if o.FabricAddressingPolicy != nil {
		fap, err = o.FabricAddressingPolicy.polish()
		if err != nil {
			return nil, err
		}
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
	return &TemplatePodBased{
		Id:             o.Id,
		templateType:   TemplateType(tType),
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &TemplatePodBasedData{
			DisplayName:             o.DisplayName,
			AntiAffinityPolicy:      *aap,
			FabricAddressingPolicy:  fap,
			Superspine:              *superspine,
			Capability:              TemplateCapability(capability),
			RackBasedTemplates:      rbt,
			RackBasedTemplateCounts: o.RackBasedTemplateCounts,
		},
	}, nil
}

var _ Template = &TemplatePodBased{}

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

type TemplateL3CollapsedData struct {
	DisplayName          string
	AntiAffinityPolicy   AntiAffinityPolicy
	RackTypes            []RackType
	Capability           TemplateCapability
	MeshLinkSpeed        LogicalDevicePortSpeed
	VirtualNetworkPolicy VirtualNetworkPolicy
	MeshLinkCount        int
	RackTypeCounts       []RackTypeCount
	DhcpServiceIntent    DhcpServiceIntent
}

type rawTemplateL3Collapsed struct {
	Id                   ObjectId                   `json:"id"`
	Type                 templateType               `json:"type"`
	DisplayName          string                     `json:"display_name"`
	AntiAffinityPolicy   rawAntiAffinityPolicy      `json:"anti_affinity_policy"`
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
	antiAffinityPolicy, err := o.AntiAffinityPolicy.polish()
	if err != nil {
		return nil, err
	}
	return &TemplateL3Collapsed{
		Id:             o.Id,
		templateType:   TemplateType(tType),
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &TemplateL3CollapsedData{
			DisplayName:          o.DisplayName,
			AntiAffinityPolicy:   *antiAffinityPolicy,
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

func (o *Client) listAllTemplateIds(ctx context.Context) ([]ObjectId, error) {
	response := &struct {
		Items []ObjectId `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

// getTemplate returns one of *TemplateRackBased, *TemplatePodBased or
// *TemplateL3Collapsed, each of which have Type() method which should be
// used to cast them into the correct type.
func (o *Client) getTemplate(ctx context.Context, id ObjectId) (template, error) {
	rawMsg := &json.RawMessage{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignTemplateById, id),
		apiResponse: rawMsg,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return template(*rawMsg), nil
}

func (o *Client) getAllTemplates(ctx context.Context) ([]template, error) {
	response := &struct {
		Items []json.RawMessage `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]template, len(response.Items))
	for i, item := range response.Items {
		result[i] = template(item)
	}

	return result, nil
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

func (o *Client) getPodBasedTemplate(ctx context.Context, id ObjectId) (*rawTemplatePodBased, error) {
	rawTemplate, err := o.getTemplate(ctx, id)
	if err != nil {
		return nil, err
	}

	tType, err := rawTemplate.templateType()
	if err != nil {
		return nil, err
	}

	if tType != templateTypePodBased {
		return nil, ClientErr{
			errType: ErrWrongType,
			err:     fmt.Errorf("template '%s' is of type '%s', not '%s'", id, tType, templateTypePodBased),
		}
	}

	result := &rawTemplatePodBased{}
	err = json.Unmarshal(rawTemplate, result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling raw pod-based template - %w", err)
	}

	// force 'type' field of included rack-based templates to "rack_based" b/c Apstra rejects empty string.
	for i, rbt := range result.RackBasedTemplates {
		switch rbt.Type {
		case "":
			result.RackBasedTemplates[i].Type = templateTypeRackBased
		case templateTypeRackBased: //fallthrough
		default:
			return nil, fmt.Errorf("rack-based template '%s' within pod-based template '%s' claims to be type '%s', expected '%s'",
				rbt.DisplayName, result.Id, rbt.Type, templateTypeRackBased)
		}
	}
	return result, nil
}

func (o *Client) getAllPodBasedTemplates(ctx context.Context) ([]rawTemplatePodBased, error) {
	templates, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var result []rawTemplatePodBased
	for _, t := range templates {
		tType, err := t.templateType()
		if err != nil {
			return nil, err
		}
		if tType != templateTypePodBased {
			continue
		}
		var raw rawTemplatePodBased
		err = json.Unmarshal(t, &raw)
		if err != nil {
			return nil, err
		}
		result = append(result, raw)
	}

	return result, nil
}

func (o *Client) getTemplateByTypeAndName(ctx context.Context, desiredType templateType, desiredName string) (*template, error) {
	templates, err := o.getAllTemplates(ctx)
	if err != nil {
		return nil, err
	}

	var found *template
	for i := range templates {
		foundType, err := templates[i].templateType()
		if err != nil {
			return nil, err
		}
		if foundType != desiredType {
			continue // wrong type
		}

		foundName, err := templates[i].displayName()
		if foundName != desiredName {
			continue // wrong name
		}
		if err != nil {
			return nil, err
		}

		if found != nil { // multiple matches!
			return nil, ClientErr{
				errType: ErrMultipleMatch,
				err:     fmt.Errorf("found multiple %s templates named '%s'", desiredType, desiredName),
			}
		}

		// record this pointer to detect multiple matches
		found = &templates[i]
	}

	if found == nil { // not found!
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no %s templates named '%s'", desiredType, desiredName),
		}
	}

	return found, nil
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

type CreateRackBasedTemplateRequest struct {
	DisplayName            string
	Capability             TemplateCapability
	Spine                  *TemplateElementSpineRequest
	RackInfos              map[ObjectId]TemplateRackBasedRackInfo
	DhcpServiceIntent      *DhcpServiceIntent
	AntiAffinityPolicy     *AntiAffinityPolicy
	AsnAllocationPolicy    *AsnAllocationPolicy
	FabricAddressingPolicy *FabricAddressingPolicy
	VirtualNetworkPolicy   *VirtualNetworkPolicy
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

		// prep the rackTypeCount object using the caller's map key (ObjectId) as
		// the link between the racktype data copy and the racktypecount
		rackTypeCounts[i].RackTypeId = k
		rackTypeCounts[i].Count = ri.Count
		i++
	}

	var err error
	var dhcpServiceIntent DhcpServiceIntent
	if o.DhcpServiceIntent != nil {
		dhcpServiceIntent = *o.DhcpServiceIntent
	}

	switch {
	case o.Spine == nil:
		return nil, errors.New("spine cannot be <nil> when creating a rack-based template")
	case o.AntiAffinityPolicy == nil:
		return nil, errors.New("anti-affinity policy cannot be <nil> when creating a rack-based template")
	case o.AsnAllocationPolicy == nil:
		return nil, errors.New("asn allocation policy cannot be <nil> when creating a rack-based template")
	case o.VirtualNetworkPolicy == nil:
		return nil, errors.New("virtual network policy cannot be <nil> when creating a rack-based template")
	}

	spine, err := o.Spine.raw(ctx, client)
	if err != nil {
		return nil, err
	}
	var antiAffinityPolicy *rawAntiAffinityPolicy
	asnAllocationPolicy := o.AsnAllocationPolicy.raw()

	var fabricAddressingPolicy *rawFabricAddressingPolicy
	if o.FabricAddressingPolicy != nil && !rackBasedTemplateFabricAddressingPolicyForbidden().Includes(client.apiVersion) {
		fabricAddressingPolicy = o.FabricAddressingPolicy.raw()
	}

	virtualNetworkPolicy := o.VirtualNetworkPolicy.raw()

	return &rawCreateRackBasedTemplateRequest{
		Type:                   templateTypeRackBased,
		DisplayName:            o.DisplayName,
		Capability:             o.Capability.raw(),
		Spine:                  *spine,
		RackTypes:              rackTypes,
		RackTypeCounts:         rackTypeCounts,
		DhcpServiceIntent:      dhcpServiceIntent,
		AntiAffinityPolicy:     antiAffinityPolicy,
		AsnAllocationPolicy:    *asnAllocationPolicy,
		FabricAddressingPolicy: fabricAddressingPolicy,
		VirtualNetworkPolicy:   *virtualNetworkPolicy,
	}, nil
}

type rawCreateRackBasedTemplateRequest struct {
	Type                   templateType               `json:"type"`
	DisplayName            string                     `json:"display_name"`
	Capability             templateCapability         `json:"capability,omitempty"`
	Spine                  rawSpine                   `json:"spine"`
	RackTypes              []rawRackType              `json:"rack_types"`
	RackTypeCounts         []RackTypeCount            `json:"rack_type_counts"`
	DhcpServiceIntent      DhcpServiceIntent          `json:"dhcp_service_intent"`
	AntiAffinityPolicy     *rawAntiAffinityPolicy     `json:"anti_affinity_policy,omitempty"`
	AsnAllocationPolicy    rawAsnAllocationPolicy     `json:"asn_allocation_policy"`
	FabricAddressingPolicy *rawFabricAddressingPolicy `json:"fabric_addressing_policy,omitempty"`
	VirtualNetworkPolicy   rawVirtualNetworkPolicy    `json:"virtual_network_policy"`
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

type CreatePodBasedTemplateRequest struct {
	DisplayName             string
	Capability              TemplateCapability
	Superspine              *TemplateElementSuperspineRequest
	RackBasedTemplateIds    []ObjectId
	RackBasedTemplateCounts []RackBasedTemplateCount
	AntiAffinityPolicy      *AntiAffinityPolicy
	FabricAddressingPolicy  *FabricAddressingPolicy
}

func (o *CreatePodBasedTemplateRequest) raw(ctx context.Context, client *Client) (*rawCreatePodBasedTemplateRequest, error) {
	var err error

	var superspine *rawSuperspine
	if o.Superspine != nil {
		superspine, err = o.Superspine.raw(ctx, client)
		if err != nil {
			return nil, err
		}
	}

	rawRackBasedTemplates := make([]rawTemplateRackBased, len(o.RackBasedTemplateIds))
	for i, id := range o.RackBasedTemplateIds {
		rbt, err := client.getRackBasedTemplate(ctx, id)
		rbt.Type = templateTypeRackBased
		if err != nil {
			return nil, err
		}
		rawRackBasedTemplates[i] = *rbt
	}

	var antiAffinityPolicy *rawAntiAffinityPolicy
	if o.AntiAffinityPolicy != nil {
		antiAffinityPolicy = o.AntiAffinityPolicy.raw()
	}

	var fabricAddressingPolicy *rawFabricAddressingPolicy
	if o.FabricAddressingPolicy != nil && !podBasedTemplateFabricAddressingPolicyForbidden().Includes(client.apiVersion) {
		fabricAddressingPolicy = o.FabricAddressingPolicy.raw()
	}

	return &rawCreatePodBasedTemplateRequest{
		Type:                    templateTypePodBased,
		DisplayName:             o.DisplayName,
		Capability:              o.Capability.raw(),
		Superspine:              *superspine,
		RackBasedTemplates:      rawRackBasedTemplates,
		RackBasedTemplateCounts: o.RackBasedTemplateCounts,
		AntiAffinityPolicy:      antiAffinityPolicy,
		FabricAddressingPolicy:  fabricAddressingPolicy,
	}, nil
}

type rawCreatePodBasedTemplateRequest struct {
	Type                    templateType               `json:"type"`
	DisplayName             string                     `json:"display_name"`
	Capability              templateCapability         `json:"capability"`
	Superspine              rawSuperspine              `json:"superspine"`
	RackBasedTemplates      []rawTemplateRackBased     `json:"rack_based_templates"`
	RackBasedTemplateCounts []RackBasedTemplateCount   `json:"rack_based_template_counts"`
	AntiAffinityPolicy      *rawAntiAffinityPolicy     `json:"anti_affinity_policy,omitempty"`
	FabricAddressingPolicy  *rawFabricAddressingPolicy `json:"fabric_addressing_policy,omitempty"`
}

func (o *Client) createPodBasedTemplate(ctx context.Context, in *rawCreatePodBasedTemplateRequest) (ObjectId, error) {
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

func (o *Client) updatePodBasedTemplate(ctx context.Context, id ObjectId, in *CreatePodBasedTemplateRequest) error {
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

type CreateL3CollapsedTemplateRequest struct {
	DisplayName          string                 `json:"display_name"`
	Capability           TemplateCapability     `json:"capability"`
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
	return &rawCreateL3CollapsedTemplateRequest{
		Type:                 templateTypeL3Collapsed,
		DisplayName:          o.DisplayName,
		Capability:           o.Capability.raw(),
		MeshLinkCount:        o.MeshLinkCount,
		MeshLinkSpeed:        *o.MeshLinkSpeed.raw(),
		RackTypes:            rackTypes,
		RackTypeCounts:       o.RackTypeCounts,
		DhcpServiceIntent:    o.DhcpServiceIntent,
		AntiAffinityPolicy:   o.AntiAffinityPolicy.raw(),
		VirtualNetworkPolicy: *o.VirtualNetworkPolicy.raw(),
	}, nil
}

type rawCreateL3CollapsedTemplateRequest struct {
	Type                 templateType              `json:"type"`
	DisplayName          string                    `json:"display_name"`
	Capability           templateCapability        `json:"capability"`
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

func (o *Client) deleteTemplate(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignTemplateById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) getTemplateType(ctx context.Context, id ObjectId) (templateType, error) {
	response := &struct {
		Type templateType `tfsdk:"type"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignTemplateById, id),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Type, nil
}

func (o *Client) getTemplateIdsTypesByName(ctx context.Context, desired string) (map[ObjectId]TemplateType, error) {
	response := &struct {
		Items []struct {
			Id          ObjectId     `json:"id"`
			Type        templateType `json:"type"`
			DisplayName string       `json:"display_name"`
		} `json:"Items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignTemplates,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make(map[ObjectId]TemplateType)
	for _, t := range response.Items {
		if t.DisplayName == desired {
			parsed, err := t.Type.parse()
			if err != nil {
				return nil, fmt.Errorf("error parsing type of template '%s' - %w", t.Id, err)
			}
			result[t.Id] = TemplateType(parsed)
		}
	}
	return result, nil
}

func (o *Client) getTemplateIdTypeByName(ctx context.Context, desired string) (ObjectId, TemplateType, error) {
	idToType, err := o.getTemplateIdsTypesByName(ctx, desired)
	if err != nil {
		return "", -1, fmt.Errorf("error fetching templates by name - %w", err)
	}

	switch len(idToType) {
	case 0:
		return "", -1, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("template named '%s' not found", desired),
		}
	case 1:
		for k, v := range idToType {
			return k, v, nil
		}
	}
	return "", -1, ClientErr{
		errType: ErrMultipleMatch,
		err:     fmt.Errorf("found multiple templates named '%s'", desired),
	}
}

// AllTemplateTypes returns the []TemplateType representing
// each supported TemplateType
func AllTemplateTypes() []TemplateType {
	i := 0
	var result []TemplateType
	for {
		var tType TemplateType
		err := tType.FromString(TemplateType(i).String())
		if err != nil {
			return result
		}
		if tType != TemplateTypeNone {
			result = append(result, tType)
		}
		i++
	}
}

// AllOverlayControlProtocols returns the []OverlayControlProtocol representing
// each supported OverlayControlProtocol
func AllOverlayControlProtocols() []OverlayControlProtocol {
	i := 0
	var result []OverlayControlProtocol
	for {
		var ocp OverlayControlProtocol
		err := ocp.FromString(OverlayControlProtocol(i).String())
		if err != nil {
			return result
		}
		result = append(result, ocp)
		i++
	}
}
