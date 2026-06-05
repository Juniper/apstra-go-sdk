// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/errors"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/str"
	"github.com/Juniper/apstra-go-sdk/internal/urls"
)

const (
	defaultSecurityZoneIDLock  = "default_security_zone"
	defaultSecurityZoneVRFName = "default"
)

// CreateSecurityZone creates an Apstra Routing Zone / Security Zone / VRF.
// If cfg.JunosEVPNIRBMode is omitted, but the API's version-dependent behavior
// requires that field, it will be set to JunosEVPNIRBModeAsymmetric in the
// request sent to the API.
func (o TwoStageL3ClosClient) CreateSecurityZone(ctx context.Context, cfg datacenter.SecurityZone) (string, error) {
	szAddressingSupportOK := compatibility.SecurityZoneAddressingSupported.Check(o.client.apiVersion)

	if (cfg.AddressingSupport != nil || cfg.DisableIPv4 != nil) && !szAddressingSupportOK {
		return "", fmt.Errorf("AddressingSupport and DisableIPv4 must be nil with Apstra %s", o.client.apiVersion)
	}

	if cfg.AddressingSupport != nil &&
		*cfg.AddressingSupport != enum.AddressingSchemeIPv6 &&
		cfg.DisableIPv4 != nil &&
		*cfg.DisableIPv4 {
		return "", fmt.Errorf("disabling IPv4 not permitted with addressing scheme %s", *cfg.AddressingSupport)
	}

	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(urls.DatacenterSecurityZones, o.blueprintId),
		apiInput:    cfg,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return string(response.Id), nil
}

// GetSecurityZone fetches the Security Zone / Routing Zone / VRF with the given id
func (o TwoStageL3ClosClient) GetSecurityZone(ctx context.Context, id string) (datacenter.SecurityZone, error) {
	var response datacenter.SecurityZone
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterSecurityZoneById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return datacenter.SecurityZone{}, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

// GetSecurityZones returns []SecurityZone representing all Security Zones /
// Routing Zones / VRFs on the system.
func (o TwoStageL3ClosClient) GetSecurityZones(ctx context.Context) ([]datacenter.SecurityZone, error) {
	response := &struct {
		Items map[string]datacenter.SecurityZone `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(urls.DatacenterSecurityZones, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}

	// This API endpoint returns a map. Convert to list for consistency with other 'GetAll' functions.
	result := make([]datacenter.SecurityZone, len(response.Items))
	var i int
	for _, v := range response.Items {
		result[i] = v
		i++
	}

	return result, nil
}

// GetSecurityZoneByLabel fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o TwoStageL3ClosClient) GetSecurityZoneByLabel(ctx context.Context, label string) (datacenter.SecurityZone, error) {
	zones, err := o.GetSecurityZones(ctx)
	if err != nil {
		return datacenter.SecurityZone{}, err
	}

	for _, zone := range zones {
		if zone.Label == label {
			return zone, nil
		}
	}

	return datacenter.SecurityZone{}, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("security zone with label %q not found in blueprint %q", label, o.blueprintId),
	}
}

// GetSecurityZoneByVRFName fetches the Security Zone / Routing Zone / VRF with
// the given label.
func (o TwoStageL3ClosClient) GetSecurityZoneByVRFName(ctx context.Context, vrfName string) (datacenter.SecurityZone, error) {
	zones, err := o.GetSecurityZones(ctx)
	if err != nil {
		return datacenter.SecurityZone{}, err
	}

	for _, zone := range zones {
		if zone.VRFName == vrfName {
			return zone, nil
		}
	}

	return datacenter.SecurityZone{}, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("security zone with vrf name %q not found in blueprint %q", vrfName, o.blueprintId),
	}
}

func (c *TwoStageL3ClosClient) DefaultSecurityZoneID(ctx context.Context) (*string, error) {
	c.client.lock(defaultSecurityZoneIDLock)
	defer c.client.unlock(defaultSecurityZoneIDLock)

	// If we know the default Security Zone ID, return it.
	if c.defaultSecurityZoneID != "" {
		return pointer.ToCopy(c.defaultSecurityZoneID), nil
	}

	// Retrieve the default Security Zone.
	sz, err := c.GetSecurityZoneByVRFName(ctx, defaultSecurityZoneVRFName)
	if err != nil {
		return nil, fmt.Errorf("failed while fetching default Routing Zone ID: %w", err)
	}

	id := sz.ID()
	if id == nil {
		return nil, errors.APIResponseInvalid("default Routing Zone ID returned nil ID")
	}

	c.defaultSecurityZoneID = *id // Cache the default Security Zone ID.

	return id, nil
}

func (c *TwoStageL3ClosClient) GetDefaultSecurityZone(ctx context.Context) (datacenter.SecurityZone, error) {
	c.client.lock(defaultSecurityZoneIDLock)
	defer c.client.unlock(defaultSecurityZoneIDLock)

	// If we know the default Security Zone ID, fetch it using the cached ID.
	if c.defaultSecurityZoneID != "" {
		return c.GetSecurityZone(ctx, c.defaultSecurityZoneID)
	}

	// Retrieve the default Security Zone the hard way.
	sz, err := c.GetSecurityZoneByVRFName(ctx, defaultSecurityZoneVRFName)
	if err != nil {
		return datacenter.SecurityZone{}, fmt.Errorf("failed while fetching default Security Zone: %w", err)
	}

	id := sz.ID()
	if id == nil {
		return datacenter.SecurityZone{}, errors.APIResponseInvalid("default Security Zone ID returned nil ID")
	}

	c.defaultSecurityZoneID = *id // Cache the default Security Zone ID.

	return sz, nil
}

func (c *TwoStageL3ClosClient) ListSecurityZones(ctx context.Context) ([]string, error) {
	items, err := c.GetSecurityZones(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(items))
	for i, item := range items {
		idPtr := item.ID()
		if idPtr == nil {
			return nil, ClientErr{
				errType: ErrInvalidId,
				err:     fmt.Errorf("the Security Zone at index %d has nil ID", i),
			}
		}
		result[i] = *idPtr
	}

	return result, nil
}

// UpdateSecurityZone replaces the configuration of zone zoneId with the supplied CreateSecurityZoneCfg
func (o TwoStageL3ClosClient) UpdateSecurityZone(ctx context.Context, v datacenter.SecurityZone) error {
	if v.ID() == nil {
		return fmt.Errorf("id is required in %s", str.FuncName())
	}

	szAddressingSupportOK := compatibility.SecurityZoneAddressingSupported.Check(o.client.apiVersion)

	if (v.AddressingSupport != nil || v.DisableIPv4 != nil) && !szAddressingSupportOK {
		return fmt.Errorf("AddressingSupport and DisableIPv4 must be nil with Apstra %s", o.client.apiVersion)
	}

	if v.AddressingSupport != nil &&
		*v.AddressingSupport != enum.AddressingSchemeIPv6 &&
		v.DisableIPv4 != nil &&
		*v.DisableIPv4 {
		return fmt.Errorf("disabling IPv4 not permitted with addressing scheme %s", *v.AddressingSupport)
	}

	// workaround for error:
	//
	// {
	//  "error_code": 422,
	//  "errors": {
	//    "disable_ipv4": "IPv4 support can only be disabled when addressing_support=\"ipv6\""
	//  }
	// }
	//
	// JP says the API behavior is deliberate: A shim layer sets an omitted `disable_ipv4` attribute
	// to the current value in PUT requests /even when doing so produces an unsupported combination/.
	//
	// Because of this API behavior, when setting addressing_support to something other than IPv6,
	// we will explicitly enable IPv4 support (disable = false)
	if v.AddressingSupport != nil &&
		*v.AddressingSupport != enum.AddressingSchemeIPv6 &&
		v.DisableIPv4 == nil &&
		szAddressingSupportOK {
		v.DisableIPv4 = pointer.To(false)
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(urls.DatacenterSecurityZoneById, o.blueprintId, *v.ID()),
		apiInput: v,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o TwoStageL3ClosClient) DeleteSecurityZone(ctx context.Context, id string) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(urls.DatacenterSecurityZoneById, o.blueprintId, id),
	})
	return convertTtaeToAceWherePossible(err)
}
