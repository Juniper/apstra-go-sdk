// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiUrlDesignTemplates       = apiUrlDesignPrefix + "templates"
	apiUrlDesignTemplatesPrefix = apiUrlDesignTemplates + apiUrlPathDelim
	apiUrlDesignTemplateById    = apiUrlDesignTemplatesPrefix + "%s"
)

type (
	AntiAffninityAlgorithm int
	antiAffinityAlgorithm  string
	AntiAffinityMode       int
	antiAffinityMode       string
	TemplateType           int
	templateType           string
	AsnAllocationScheme    int
	asnAllocationScheme    string
	AddressingScheme       int
	addressingScheme       string
)

type (
	OverlayControlProtocol int
	overlayControlProtocol string
	TemplateCapability     int
	templateCapability     string
)

const (
	AntiAffinityModeDisabled = AntiAffinityMode(iota)
	AntiAffinityModeEnabledLoose
	AntiAffinityModeEnabledStrict
	AntiAffinityModeLoose
	AntiAffinityModeStrict
	AntiAffinityModeUnknown = "unknown anti affinity mode %s"

	antiAffinityModeDisabled      = antiAffinityMode("disabled")
	antiAffinityModeEnabledLoose  = antiAffinityMode("enabled_loose")
	antiAffinityModeEnabledStrict = antiAffinityMode("enabled_strict")
	antiAffinityModeLoose         = antiAffinityMode("loose")
	antiAffinityModeStrict        = antiAffinityMode("strict")
	antiAffinityModeUnknown       = "unknown anti affinity mode %d"
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
	case AntiAffinityModeEnabledLoose:
		return string(antiAffinityModeEnabledLoose)
	case AntiAffinityModeEnabledStrict:
		return string(antiAffinityModeEnabledStrict)
	case AntiAffinityModeLoose:
		return string(antiAffinityModeLoose)
	case AntiAffinityModeStrict:
		return string(antiAffinityModeStrict)
	default:
		return fmt.Sprintf(antiAffinityModeUnknown, o)
	}
}

func (o *AntiAffinityMode) FromString(s string) error {
	i, err := antiAffinityMode(s).parse()
	if err != nil {
		return err
	}
	*o = AntiAffinityMode(i)
	return nil
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
	case antiAffinityModeEnabledLoose:
		return int(AntiAffinityModeEnabledLoose), nil
	case antiAffinityModeEnabledStrict:
		return int(AntiAffinityModeEnabledStrict), nil
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

type RackTypeCount struct {
	RackTypeId ObjectId `json:"rack_type_id"`
	Count      int      `json:"count"`
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

type DhcpServiceIntent struct {
	Active bool `json:"active"`
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
