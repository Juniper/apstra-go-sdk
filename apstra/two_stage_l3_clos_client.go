// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
)

const (
	blueprintTypeParam   = "type"
	dcClientMaxRetries   = 10
	dcClientRetryBackoff = 100 * time.Millisecond
)

type (
	BlueprintType int
	blueprintType string
)

const (
	BlueprintTypeNone = BlueprintType(iota)
	BlueprintTypeConfig
	BlueprintTypeDeployed
	BlueprintTypeOperation
	BlueprintTypeStaging

	blueprintTypeNone      = blueprintType("")
	blueprintTypeConfig    = blueprintType("config")
	blueprintTypeDeployed  = blueprintType("deployed")
	blueprintTypeOperation = blueprintType("operation")
	blueprintTypeStaging   = blueprintType("staging")
	blueprintTypeUnknown   = "unknown-blueprint-type-%d"
)

func (o BlueprintType) raw() blueprintType {
	switch o {
	case BlueprintTypeNone:
		return blueprintTypeNone
	case BlueprintTypeConfig:
		return blueprintTypeConfig
	case BlueprintTypeDeployed:
		return blueprintTypeDeployed
	case BlueprintTypeStaging:
		return blueprintTypeStaging
	case BlueprintTypeOperation:
		return blueprintTypeOperation
	default:
		return blueprintType(fmt.Sprintf(blueprintTypeUnknown, o))
	}
}

func (o BlueprintType) string() string {
	return string(o.raw())
}

type TwoStageL3ClosClient struct {
	client        *Client
	blueprintId   ObjectId
	Mutex         Mutex
	blueprintType BlueprintType
	nodeIdsByType map[NodeType][]ObjectId
}

// Id returns the client's Blueprint ID
func (o *TwoStageL3ClosClient) Id() ObjectId {
	return o.blueprintId
}

// lockId returns a string intended to be used with Client.lock()
func (o *TwoStageL3ClosClient) lockId(ids ...ObjectId) string {
	var buf bytes.Buffer
	buf.WriteString(o.blueprintId.String())
	for _, id := range ids {
		buf.WriteString(mutexKeySeparator + id.String())
	}
	return buf.String()
}

// SetType sets the client's internal BlueprintType value (staging, etc...).
// This value is in HTTP requests as a query string argument, e.g.
//
//	'?type=staging'
func (o *TwoStageL3ClosClient) SetType(bpt BlueprintType) {
	o.blueprintType = bpt
}

// urlWithParam is a helper function which uses the blueprintType element to
// decorate a *URL with the required query parameter.
//
//lint:ignore U1000 keep for future use
func (o *TwoStageL3ClosClient) urlWithParam(in string) (*url.URL, error) {
	apstraUrl, err := url.Parse(in)
	if err != nil {
		return nil, err
	}

	if o.blueprintType != BlueprintTypeNone {
		params := apstraUrl.Query()
		params.Set(blueprintTypeParam, o.blueprintType.string())
		apstraUrl.RawQuery = params.Encode()
	}
	return apstraUrl, nil
}

// GetResourceAllocations returns ResourceGroupAllocations representing
// all allocations of resource pools to blueprint requirements
func (o *TwoStageL3ClosClient) GetResourceAllocations(ctx context.Context) (ResourceGroupAllocations, error) {
	rawRgaSlice, err := o.getAllResourceAllocations(ctx)
	if err != nil {
		return nil, err
	}
	result := make(ResourceGroupAllocations, len(rawRgaSlice))
	for i, raw := range rawRgaSlice {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetResourceAllocation takes a *ResourceGroup and returns a
// *ResourceGroupAllocation with fields populated based on the Apstra API
// response.
func (o *TwoStageL3ClosClient) GetResourceAllocation(ctx context.Context, in *ResourceGroup) (*ResourceGroupAllocation, error) {
	rga, err := o.getResourceAllocation(ctx, in)
	if err != nil {
		return nil, err
	}
	return rga.polish()
}

// SetResourceAllocation sets the supplied resource allocation, overwriting any
// prior allocations with the supplied info.
func (o *TwoStageL3ClosClient) SetResourceAllocation(ctx context.Context, in *ResourceGroupAllocation) error {
	return o.setResourceAllocation(ctx, in.raw())
}

// GetInterfaceMapAssignments returns a SystemIdToInterfaceMapAssignment (a map
// of string (blueprint graph node ID) to interface map ID detailing assignments
// in the specified blueprint:
//
//	x := SystemIdToInterfaceMapAssignment{
//		"BeAyAoCIgqx4r3hiFow": nil,
//		"B3Ym-PBJJEtvXQsnQQM": "VS_SONiC_BUZZNIK_PLUS__slicer-7x10-1",
//		"4gCWV2NRix6MYPm4PHU": "Arista_vEOS__slicer-7x10-1",
//	}
func (o *TwoStageL3ClosClient) GetInterfaceMapAssignments(ctx context.Context) (SystemIdToInterfaceMapAssignment, error) {
	return o.getInterfaceMapAssignments(ctx)
}

// SetInterfaceMapAssignments sets the supplied interface map assignments,
// overwriting any prior assignments with the supplied info. It returns
// the Blueprint config revision number.
func (o *TwoStageL3ClosClient) SetInterfaceMapAssignments(ctx context.Context, assignments SystemIdToInterfaceMapAssignment) error {
	return o.setInterfaceMapAssignments(ctx, assignments)
}

// CreateSecurityZone creates an Apstra Routing Zone / Security Zone / VRF.
// If cfg.JunosEvpnIrbMode is omitted, but the API's version-dependent behavior
// requires that field, it will be set to JunosEvpnIrbModeAsymmetric in the
// request sent to the API.
func (o *TwoStageL3ClosClient) CreateSecurityZone(ctx context.Context, cfg *SecurityZoneData) (ObjectId, error) {
	raw := cfg.raw()
	if raw.JunosEvpnIrbMode == "" {
		raw.JunosEvpnIrbMode = enum.JunosEvpnIrbModeAsymmetric.Value
	}

	response, err := o.createSecurityZone(ctx, raw)
	if err != nil {
		return "", err
	}
	return response.Id, nil
}

// DeleteSecurityZone deletes an Apstra Routing Zone / Security Zone / VRF
func (o *TwoStageL3ClosClient) DeleteSecurityZone(ctx context.Context, zoneId ObjectId) error {
	return o.deleteSecurityZone(ctx, zoneId)
}

// GetSecurityZoneDhcpServers returns []net.IP representing the DHCP relay
// targets for the security zone specified by zoneId.
func (o *TwoStageL3ClosClient) GetSecurityZoneDhcpServers(ctx context.Context, zoneId ObjectId) ([]net.IP, error) {
	var err error
	ips, err := o.getSecurityZoneDhcpServers(ctx, zoneId)
	if err != nil {
		return nil, err
	}

	result := make([]net.IP, len(ips))
	for i, s := range ips {
		result[i] = net.ParseIP(s)
		if result[i] == nil {
			err = errors.Join(err, fmt.Errorf("failed to parse blueprint %s security zone %s dhcp server"+
				" at index %d; expected an IP address, got %q", o.blueprintId, zoneId, i, s))
		}
	}
	return result, err
}

// SetSecurityZoneDhcpServers assigns the []net.IP as DHCP relay targets for
// the specified security zone, overwriting whatever is there. On the Apstra
// side, the servers seem to be maintained as an ordered list with duplicates
// permitted (though the web UI sorts the data prior to display)
func (o *TwoStageL3ClosClient) SetSecurityZoneDhcpServers(ctx context.Context, zoneId ObjectId, IPs []net.IP) error {
	ips := make([]string, len(IPs))
	for i, ip := range IPs {
		ips[i] = ip.String()
	}
	return o.setSecurityZoneDhcpServers(ctx, zoneId, ips)
}

// GetSecurityZone fetches the Security Zone / Routing Zone / VRF with the given
// zoneId.
func (o *TwoStageL3ClosClient) GetSecurityZone(ctx context.Context, zoneId ObjectId) (*SecurityZone, error) {
	raw, err := o.getSecurityZone(ctx, zoneId)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetSecurityZoneByVrfName fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o *TwoStageL3ClosClient) GetSecurityZoneByVrfName(ctx context.Context, vrfName string) (*SecurityZone, error) {
	raw, err := o.getSecurityZoneByVrfName(ctx, vrfName)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// GetAllSecurityZones returns []SecurityZone representing all Security Zones /
// Routing Zones / VRFs on the system.
func (o *TwoStageL3ClosClient) GetAllSecurityZones(ctx context.Context) ([]SecurityZone, error) {
	response, err := o.getAllSecurityZones(ctx)
	if err != nil {
		return nil, err
	}

	// This API endpoint returns a map. Convert to list for consistency with other 'GetAll' functions.
	result := make([]SecurityZone, len(response))
	var i int
	for k := range response {
		polished, err := response[k].polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
		i++
	}

	return result, nil
}

// UpdateSecurityZone replaces the configuration of zone zoneId with the supplied CreateSecurityZoneCfg
func (o *TwoStageL3ClosClient) UpdateSecurityZone(ctx context.Context, zoneId ObjectId, cfg *SecurityZoneData) error {
	if cfg.JunosEvpnIrbMode == nil {
		return errors.New("junos_evpn_irb_mode cannot be nil")
	}

	return o.updateSecurityZone(ctx, zoneId, cfg.raw())
}

// GetAllPolicies returns []Policy representing all policies configured within the DC blueprint
func (o *TwoStageL3ClosClient) GetAllPolicies(ctx context.Context) ([]Policy, error) {
	policies, err := o.getAllPolicies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Policy, len(policies))
	for i, raw := range policies {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}
		result[i] = *polished
	}
	return result, nil
}

// GetPolicy returns *Policy representing policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) GetPolicy(ctx context.Context, id ObjectId) (*Policy, error) {
	raw, err := o.getPolicy(ctx, id)
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

// GetPolicyByLabel returns *Policy representing policy identified by 'label' within the DC blueprint
func (o *TwoStageL3ClosClient) GetPolicyByLabel(ctx context.Context, label string) (*Policy, error) {
	raw, err := o.getPolicyByLabel(ctx, label)
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

// CreatePolicy creates a policy within the DC blueprint, returns its ID
func (o *TwoStageL3ClosClient) CreatePolicy(ctx context.Context, data *PolicyData) (ObjectId, error) {
	return o.createPolicy(ctx, data.request())
}

// DeletePolicy deletes policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) DeletePolicy(ctx context.Context, id ObjectId) error {
	return o.deletePolicy(ctx, id)
}

// UpdatePolicy calls PUT to replace the configuration of policy 'id' within the DC blueprint
func (o *TwoStageL3ClosClient) UpdatePolicy(ctx context.Context, id ObjectId, data *PolicyData) error {
	return o.updatePolicy(ctx, id, data.request())
}

// AddPolicyRule adds a policy rule at 'position' (bumping all other rules
// down). Position 0 makes the new policy first on the list, 1 makes it second
// on the list, etc... Use -1 for last on the list. The returned ObjectId
// represents the new rule
func (o *TwoStageL3ClosClient) AddPolicyRule(ctx context.Context, rule *PolicyRuleData, position int, policyId ObjectId) (ObjectId, error) {
	return o.addPolicyRule(ctx, rule.raw(), position, policyId)
}

// DeletePolicyRuleById deletes the given rule. If the rule doesn't exist, a
// ClientErr with ErrNotFound is returned.
func (o *TwoStageL3ClosClient) DeletePolicyRuleById(ctx context.Context, policyId ObjectId, ruleId ObjectId) error {
	return o.deletePolicyRuleById(ctx, policyId, ruleId)
}

// CreateVirtualNetwork creates a new virtual network according to the supplied VirtualNetworkData
func (o *TwoStageL3ClosClient) CreateVirtualNetwork(ctx context.Context, in *VirtualNetworkData) (ObjectId, error) {
	if in.Tags != nil {
		return "", ClientErr{
			errType: ErrNotSupported,
			err:     errors.New("tags must be nil when creating virtual network"),
		}
	}

	return o.createVirtualNetwork(ctx, in)
}

// ListAllVirtualNetworkIds returns []ObjectId representing virtual networks configured in the blueprint
func (o *TwoStageL3ClosClient) ListAllVirtualNetworkIds(ctx context.Context) ([]ObjectId, error) {
	return o.listAllVirtualNetworkIds(ctx)
}

// GetVirtualNetwork returns *VirtualNetwork representing the given vnId
func (o *TwoStageL3ClosClient) GetVirtualNetwork(ctx context.Context, vnId ObjectId) (*VirtualNetwork, error) {
	return o.getVirtualNetwork(ctx, vnId)
}

// GetVirtualNetworkByName returns *VirtualNetwork representing the given VN name
func (o *TwoStageL3ClosClient) GetVirtualNetworkByName(ctx context.Context, name string) (*VirtualNetwork, error) {
	return o.getVirtualNetworkByName(ctx, name)
}

// GetAllVirtualNetworks return map[ObjectId]VirtualNetwork representing all
// virtual networks configured in Apstra. NOTE: the underlying API call DOES NOT
// RETURN the SVI information, so each map entry will have a nil slice at it's
// Data.SviIps struct element.
func (o *TwoStageL3ClosClient) GetAllVirtualNetworks(ctx context.Context) (map[ObjectId]VirtualNetwork, error) {
	return o.getAllVirtualNetworks(ctx)
}

// UpdateVirtualNetwork updates the virtual network specified by ID using the
// VirtualNetworkData and HTTP method PUT.
func (o *TwoStageL3ClosClient) UpdateVirtualNetwork(ctx context.Context, id ObjectId, in *VirtualNetworkData) error {
	if in.Tags != nil {
		return ClientErr{
			errType: ErrNotSupported,
			err:     errors.New("tags must be nil when updating virtual network"),
		}
	}

	return o.updateVirtualNetwork(ctx, id, in)
}

// DeleteVirtualNetwork deletes the virtual network specified by id from the
// blueprint.
func (o *TwoStageL3ClosClient) DeleteVirtualNetwork(ctx context.Context, id ObjectId) error {
	return o.deleteVirtualNetwork(ctx, id)
}

// GetNodes fetches the node of the specified type, unpacks the API response
// into 'response'
func (o *TwoStageL3ClosClient) GetNodes(ctx context.Context, nodeType NodeType, response interface{}) error {
	return o.client.getNodes(ctx, o.blueprintId, nodeType, response)
}

// PatchNode patches (only submitted fields are changed) the specified node
// using the contents of 'request'. The server's response (whole node info
// without map wrapper?) is returned in 'response'
func (o *TwoStageL3ClosClient) PatchNode(ctx context.Context, node ObjectId, request interface{}, response interface{}) error {
	return o.client.PatchNode(ctx, o.blueprintId, node, request, response)
}

// PatchNodeUnsafe patches (only submitted fields are changed) the specified node
// using the contents of 'request', and the "allow_unsafe=true" request parameter
// required by Apstra 5.0.0. The server's response (whole node info
// without map wrapper?) is returned in 'response'
func (o *TwoStageL3ClosClient) PatchNodeUnsafe(ctx context.Context, node ObjectId, request interface{}, response interface{}) error {
	return o.client.PatchNodeUnsafe(ctx, o.blueprintId, node, request, response)
}

// Client returns the embedded *Client
func (o *TwoStageL3ClosClient) Client() *Client {
	return o.client
}

// GetAllRoutingPolicies returns []DcRoutingPolicy representing all routing
// policies in the blueprint.
func (o *TwoStageL3ClosClient) GetAllRoutingPolicies(ctx context.Context) ([]DcRoutingPolicy, error) {
	raw, err := o.getAllRoutingPolicies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]DcRoutingPolicy, len(raw))
	for i := range raw {
		polished, err := raw[i].polish()
		if err != nil {
			return nil, err
		}

		result[i] = *polished
	}
	return result, nil
}

// GetRoutingPolicy returns *DcRoutingPolicy representing the specified policy.
func (o *TwoStageL3ClosClient) GetRoutingPolicy(ctx context.Context, id ObjectId) (*DcRoutingPolicy, error) {
	raw, err := o.getRoutingPolicy(ctx, id)
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

func (o *TwoStageL3ClosClient) GetRoutingPolicyByName(ctx context.Context, desired string) (*DcRoutingPolicy, error) {
	raw, err := o.getRoutingPolicyByName(ctx, desired)
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

// GetDefaultRoutingPolicy returns *DcRoutingPolicy representing the
// "default_immutable" routing policy attached to the blueprint.
func (o *TwoStageL3ClosClient) GetDefaultRoutingPolicy(ctx context.Context) (*DcRoutingPolicy, error) {
	raw, err := o.getDefaultRoutingPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return raw.polish()
}

// CreateRoutingPolicy creates a blueprint routing policy according to the
// supplied *DcRoutingPolicyData.
func (o *TwoStageL3ClosClient) CreateRoutingPolicy(ctx context.Context, in *DcRoutingPolicyData) (ObjectId, error) {
	return o.createRoutingPolicy(ctx, in.raw())
}

// UpdateRoutingPolicy modifies the blueprint routing policy specified by 'id'
// according to the supplied *DcRoutingPolicyData.
func (o *TwoStageL3ClosClient) UpdateRoutingPolicy(ctx context.Context, id ObjectId, in *DcRoutingPolicyData) error {
	return o.updateRoutingPolicy(ctx, id, in.raw())
}

// DeleteRoutingPolicy deletes the routing policy specified by id.
func (o *TwoStageL3ClosClient) DeleteRoutingPolicy(ctx context.Context, id ObjectId) error {
	return o.deleteRoutingPolicy(ctx, id)
}

// GetAllPropertySets returns []TwoStageL3ClosPropertySet representing
// all property sets imported into a blueprint
func (o *TwoStageL3ClosClient) GetAllPropertySets(ctx context.Context) ([]TwoStageL3ClosPropertySet, error) {
	return o.getAllPropertySets(ctx)
}

// GetPropertySet returns *TwoStageL3ClosPropertySet representing the
// imported property set with the given ID in the specified blueprint
func (o *TwoStageL3ClosClient) GetPropertySet(ctx context.Context, id ObjectId) (*TwoStageL3ClosPropertySet, error) {
	return o.getPropertySet(ctx, id)
}

// GetPropertySetByName returns *TwoStageL3ClosPropertySet representing
// the only property set with the given label, or an error if multiple
// property sets share the label.
func (o *TwoStageL3ClosClient) GetPropertySetByName(ctx context.Context, in string) (*TwoStageL3ClosPropertySet, error) {
	return o.getPropertySetByName(ctx, in)
}

// ImportPropertySet imports a property set into a blueprint. On success,
// it returns the id of the imported property set. Optionally, a set of keys
// can be part of the request
func (o *TwoStageL3ClosClient) ImportPropertySet(ctx context.Context, psid ObjectId, keys ...string) (ObjectId, error) {
	return o.importPropertySet(ctx, psid, keys...)
}

// UpdatePropertySet updates a property set imported into a blueprint.
// Optionally, a set of keys can be part of the request
func (o *TwoStageL3ClosClient) UpdatePropertySet(ctx context.Context, psid ObjectId, keys ...string) error {
	return o.updatePropertySet(ctx, psid, keys...)
}

// DeletePropertySet deletes a property set given the id
func (o *TwoStageL3ClosClient) DeletePropertySet(ctx context.Context, id ObjectId) error {
	return o.deletePropertySet(ctx, id)
}

// GetAllConfiglets returns []TwoStageL3ClosConfiglet representing all
// configlets imported into a blueprint
func (o *TwoStageL3ClosClient) GetAllConfiglets(ctx context.Context) ([]TwoStageL3ClosConfiglet, error) {
	return o.getAllConfiglets(ctx)
}

// GetAllConfigletIds returns Ids of all the configlets imported into a
// blueprint
func (o *TwoStageL3ClosClient) GetAllConfigletIds(ctx context.Context) ([]ObjectId, error) {
	return o.getAllConfigletIds(ctx)
}

// GetConfiglet returns *TwoStageL3ClosConfiglet representing the imported
// configlet with the given ID in the specified blueprint
func (o *TwoStageL3ClosClient) GetConfiglet(ctx context.Context, id ObjectId) (*TwoStageL3ClosConfiglet, error) {
	return o.getConfiglet(ctx, id)
}

// GetConfigletByName returns *TwoStageL3ClosConfiglet representing the only
// configlet with the given label, or an error if no configlet by that name exists
func (o *TwoStageL3ClosClient) GetConfigletByName(ctx context.Context, in string) (*TwoStageL3ClosConfiglet, error) {
	return o.getConfigletByName(ctx, in)
}

// ImportConfigletById imports a configlet from the catalog into a blueprint.
// cid is the Id catalog configlet of the
// condtion is a string input that indicates which devices it applies to.
// label can be used to rename the configlet in the blue print
// On success, it returns the id of the imported configlet.
func (o *TwoStageL3ClosClient) ImportConfigletById(ctx context.Context, cid ObjectId, condition string, label string) (ObjectId, error) {
	cfg, err := o.client.getConfiglet(ctx, cid)
	if err != nil {
		return "", err
	}

	return o.createConfiglet(ctx, &TwoStageL3ClosConfigletData{
		Condition: condition,
		Label:     label,
		Data: &ConfigletData{
			RefArchs:    cfg.Data.RefArchs,
			Generators:  cfg.Data.Generators,
			DisplayName: cfg.Data.DisplayName,
		},
	})
}

// CreateConfiglet creates a configlet described by a TwoStageL3ClosConfigletData structure
// in a blueprint.
func (o *TwoStageL3ClosClient) CreateConfiglet(ctx context.Context, c *TwoStageL3ClosConfigletData) (ObjectId, error) {
	if c.Data == nil {
		return "", errors.New("cannot create configlet with nil data")
	}
	if len(c.Data.RefArchs) != 0 {
		return "", errors.New("RefArchs not permitted when creating configlet")
	}

	return o.createConfiglet(ctx, c)
}

// UpdateConfiglet updates a configlet imported into a blueprint.
func (o *TwoStageL3ClosClient) UpdateConfiglet(ctx context.Context, id ObjectId, c *TwoStageL3ClosConfigletData) error {
	if c.Data == nil {
		return errors.New("cannot update configlet with nil data")
	}
	if len(c.Data.RefArchs) != 0 {
		return errors.New("RefArchs not permitted when updating configlet")
	}

	return o.updateConfiglet(ctx, id, c)
}

// DeleteConfiglet deletes a configlet from the blueprint given the id
func (o *TwoStageL3ClosClient) DeleteConfiglet(ctx context.Context, id ObjectId) error {
	return o.deleteConfiglet(ctx, id)
}

// GetAllSystemNodeInfos return map[ObjectId]SystemNodeInfo describing all
// "system" nodes in the blueprint.
func (o *TwoStageL3ClosClient) GetAllSystemNodeInfos(ctx context.Context) (map[ObjectId]SystemNodeInfo, error) {
	rawNodeInfos, err := o.getAllSystemNodeInfos(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[ObjectId]SystemNodeInfo, len(rawNodeInfos))
	for _, raw := range rawNodeInfos {
		polished, err := raw.polish()
		if err != nil {
			return nil, err
		}

		result[polished.Id] = *polished
	}

	return result, nil
}

// GetSystemNodeInfo returns a SystemNodeInfo describing a "system" node.
func (o *TwoStageL3ClosClient) GetSystemNodeInfo(ctx context.Context, nodeId ObjectId) (*SystemNodeInfo, error) {
	rawNodeInfos, err := o.getAllSystemNodeInfos(ctx)
	if err != nil {
		return nil, err
	}

	for _, raw := range rawNodeInfos {
		if raw.Id == nodeId {
			return raw.polish()
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("system node %q not found in blueprint %q", nodeId, o.blueprintId),
	}
}

// InstantiateIbaPredefinedProbe instantiates a predefined probe using the name and properties specified in data
// and returns the id of the created probe on success, or a blank and error on failure.
func (o *TwoStageL3ClosClient) InstantiateIbaPredefinedProbe(ctx context.Context, data *IbaPredefinedProbeRequest) (ObjectId, error) {
	return o.client.instantiatePredefinedIbaProbe(ctx, o.blueprintId, data)
}

// GetAllIbaPredefinedProbes lists all the Predefined IBA probes available to a blueprint
func (o *TwoStageL3ClosClient) GetAllIbaPredefinedProbes(ctx context.Context) ([]IbaPredefinedProbe, error) {
	return o.client.getAllIbaPredefinedProbes(ctx, o.blueprintId)
}

// GetIbaPredefinedProbeByName locates a predefined probe by name
func (o *TwoStageL3ClosClient) GetIbaPredefinedProbeByName(ctx context.Context, name string) (*IbaPredefinedProbe, error) {
	return o.client.getIbaPredefinedProbeByName(ctx, o.blueprintId, name)
}

// GetIbaProbe returns the IBA Probe that matches the ID
func (o *TwoStageL3ClosClient) GetIbaProbe(ctx context.Context, id ObjectId) (*IbaProbe, error) {
	probe, err := o.client.getIbaProbe(ctx, o.blueprintId, id)
	if err != nil {
		return nil, err
	}

	return probe, err
}

// GetIbaProbeState returns the State of the IBA Probe that matches the ID
func (o *TwoStageL3ClosClient) GetIbaProbeState(ctx context.Context, id ObjectId) (*IbaProbeState, error) {
	probe, err := o.client.getIbaProbeState(ctx, o.blueprintId, id)
	if err != nil {
		return nil, err
	}

	return probe, err
}

// DeleteIbaProbe deletes an IBA Probe
func (o *TwoStageL3ClosClient) DeleteIbaProbe(ctx context.Context, id ObjectId) error {
	return o.client.deleteIbaProbe(ctx, o.blueprintId, id)
}

// CreateIbaProbeFromJson creates an IBA Probe
func (o *TwoStageL3ClosClient) CreateIbaProbeFromJson(ctx context.Context, probeJson json.RawMessage) (ObjectId, error) {
	return o.client.createIbaProbeFromJson(ctx, o.blueprintId, probeJson)
}

// ListAllIbaPredefinedDashboardIds returns a list of Predefined IBA Dashboards in the blueprint
func (o *TwoStageL3ClosClient) ListAllIbaPredefinedDashboardIds(ctx context.Context) ([]ObjectId, error) {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return nil, fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}
	return o.client.listAllIbaPredefinedDashboardIds(ctx, o.blueprintId)
}

// InstantiateIbaPredefinedDashboard instantiates a Predefined IBA Dashboard
func (o *TwoStageL3ClosClient) InstantiateIbaPredefinedDashboard(ctx context.Context, dashboardId ObjectId, label string) (ObjectId, error) {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return "", fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}
	return o.client.instantiateIbaPredefinedDashboard(ctx, o.blueprintId, dashboardId, label)
}

// GetAllIbaDashboards returns a list of IBA Dashboards in the blueprint
func (o *TwoStageL3ClosClient) GetAllIbaDashboards(ctx context.Context) ([]IbaDashboard, error) {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return nil, fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}
	return o.client.getAllIbaDashboards(ctx, o.blueprintId)
}

// GetIbaDashboard returns the IBA Dashboard that matches the ID
func (o *TwoStageL3ClosClient) GetIbaDashboard(ctx context.Context, id ObjectId) (*IbaDashboard, error) {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return nil, fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}
	return o.client.getIbaDashboard(ctx, o.blueprintId, id)
}

// GetIbaDashboardByLabel returns the IBA Dashboard that matches the label.
// It will return an error if more than one IBA dashboard matches the label.
func (o *TwoStageL3ClosClient) GetIbaDashboardByLabel(ctx context.Context, label string) (*IbaDashboard, error) {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return nil, fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}

	return o.client.getIbaDashboardByLabel(ctx, o.blueprintId, label)
}

// CreateIbaDashboard creates an IBA Dashboard and returns the id of the created dashboard on success,
// or a blank and error on failure
func (o *TwoStageL3ClosClient) CreateIbaDashboard(ctx context.Context, data *IbaDashboardData) (ObjectId, error) {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return "", fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}

	return o.client.createIbaDashboard(ctx, o.blueprintId, data)
}

// UpdateIbaDashboard updates an IBA Dashboard and returns an error on failure
func (o *TwoStageL3ClosClient) UpdateIbaDashboard(ctx context.Context, id ObjectId, data *IbaDashboardData) error {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}

	return o.client.updateIbaDashboard(ctx, o.blueprintId, id, data)
}

// DeleteIbaDashboard deletes an IBA Dashboard and returns an error on failure
func (o *TwoStageL3ClosClient) DeleteIbaDashboard(ctx context.Context, id ObjectId) error {
	if !compatibility.IbaDashboardSupported.Check(o.client.apiVersion) {
		return fmt.Errorf("this version of the SDK will not support IBA Dashboards with Asptra %s", o.client.apiVersion)
	}

	return o.client.deleteIbaDashboard(ctx, o.blueprintId, id)
}

func (o *TwoStageL3ClosClient) refreshNodeIdsByType(ctx context.Context, nt NodeType) error {
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetClient(o.client).
		Node([]QEEAttribute{
			nt.QEEAttribute(),
			{Key: "name", Value: QEStringVal("node")},
		})

	var queryResponse struct {
		Items []struct {
			Node struct {
				Id ObjectId `json:"id"`
			} `json:"node"`
		} `json:"items"`
	}

	err := query.Do(ctx, &queryResponse)
	if err != nil {
		return fmt.Errorf("failed to query for %s nodes - %w", nt.String(), convertTtaeToAceWherePossible(err))
	}

	o.nodeIdsByType[nt] = make([]ObjectId, len(queryResponse.Items))
	for i, item := range queryResponse.Items {
		o.nodeIdsByType[nt][i] = item.Node.Id
	}

	return nil
}

func (o *TwoStageL3ClosClient) NodeIdsByType(ctx context.Context, nt NodeType) ([]ObjectId, error) {
	lockId := o.lockId("node_ids")
	o.client.lock(lockId)
	defer o.client.unlock(lockId)

	if nodeIds, ok := o.nodeIdsByType[nt]; ok {
		return nodeIds, nil // already done!
	}

	err := o.refreshNodeIdsByType(ctx, nt)
	if err != nil {
		return nil, err
	}

	return o.nodeIdsByType[nt], nil
}

func (o *TwoStageL3ClosClient) RefreshNodeIdsByType(ctx context.Context, nt NodeType) ([]ObjectId, error) {
	lockId := o.lockId("node_ids")
	o.client.lock(lockId)
	defer o.client.unlock(lockId)

	err := o.refreshNodeIdsByType(ctx, nt)
	if err != nil {
		return nil, err
	}

	return o.nodeIdsByType[nt], nil
}

// GetFabricSettings gets the fabric settings
func (o *TwoStageL3ClosClient) GetFabricSettings(ctx context.Context) (*FabricSettings, error) {
	var raw *rawFabricSettings
	var err error

	switch {
	case compatibility.FabricSettingsApiOk.Check(o.client.apiVersion):
		raw, err = o.getFabricSettings(ctx)
	case compatibility.EqApstra420.Check(o.client.apiVersion):
		raw, err = o.getFabricSettings420(ctx)
	default:
		return nil, fmt.Errorf("cannot invoke GetFabricSettings, not supported with Apstra version %q", o.client.apiVersion)
	}
	if err != nil {
		return nil, err
	}

	return raw.polish()
}

// SetFabricSettings sets the specified fabric settings
func (o *TwoStageL3ClosClient) SetFabricSettings(ctx context.Context, in *FabricSettings) error {
	if in.SpineLeafLinks != nil || in.SpineSuperspineLinks != nil {
		return errors.New("SpineLeafLinks and SpineSuperspineLinks must be nil in SetFabricSettings()")
	}

	switch {
	case compatibility.FabricSettingsApiOk.Check(o.client.apiVersion):
		return o.setFabricSettings(ctx, in.raw())
	case compatibility.EqApstra420.Check(o.client.apiVersion):
		return o.setFabricSettings420(ctx, in.raw())
	}

	return fmt.Errorf("cannot invoke SetFabricSettings, not supported with Apstra version %q", o.client.apiVersion)
}
