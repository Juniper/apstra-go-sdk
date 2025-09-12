// Copyright (c) Juniper Networks, Inc., 2023-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/enum"
)

type rawModularDeviceSlotConfiguration struct {
	LinecardProfileId ObjectId `json:"linecard_profile_id"`
	SlotId            uint64   `json:"slot_id"`
}

type ModularDeviceSlotConfiguration struct {
	LinecardProfileId ObjectId
}

type ModularDeviceProfile struct {
	Label              string
	ChassisProfileId   ObjectId
	SlotConfigurations map[uint64]ModularDeviceSlotConfiguration
}

func (o *ModularDeviceProfile) raw() *rawModularDeviceProfile {
	result := &rawModularDeviceProfile{
		DeviceProfileType:  enum.DeviceProfileTypeModular.Value,
		Label:              o.Label,
		ChassisProfileId:   o.ChassisProfileId,
		SlotConfigurations: make([]rawModularDeviceSlotConfiguration, len(o.SlotConfigurations)),
	}

	var i int
	for slotId, slotConfiguration := range o.SlotConfigurations {
		result.SlotConfigurations[i] = rawModularDeviceSlotConfiguration{
			SlotId:            slotId,
			LinecardProfileId: slotConfiguration.LinecardProfileId,
		}
		i++
	}

	return result
}

type rawModularDeviceProfile struct {
	DeviceProfileType  string                              `json:"device_profile_type"`
	Label              string                              `json:"label"`
	ChassisProfileId   ObjectId                            `json:"chassis_profile_id"`
	SlotConfigurations []rawModularDeviceSlotConfiguration `json:"slot_configuration"`
}

func (o *rawModularDeviceProfile) polish() *ModularDeviceProfile {
	result := &ModularDeviceProfile{
		Label:              o.Label,
		ChassisProfileId:   o.ChassisProfileId,
		SlotConfigurations: make(map[uint64]ModularDeviceSlotConfiguration, len(o.SlotConfigurations)),
	}

	for _, slotConfiguration := range o.SlotConfigurations {
		result.SlotConfigurations[slotConfiguration.SlotId] = ModularDeviceSlotConfiguration{
			LinecardProfileId: slotConfiguration.LinecardProfileId,
		}
	}

	return result
}

func (o *Client) createModularDeviceProfile(ctx context.Context, in *rawModularDeviceProfile) (ObjectId, error) {
	response := new(objectIdResponse)
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDeviceProfiles,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *Client) getModularDeviceProfile(ctx context.Context, id ObjectId) (*rawModularDeviceProfile, error) {
	response := new(rawModularDeviceProfile)
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDeviceProfileById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *Client) updateModularDeviceProfile(ctx context.Context, id ObjectId, cfg *rawModularDeviceProfile) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDeviceProfileById, id),
		apiInput: cfg,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *Client) deleteModularDeviceProfile(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDeviceProfileById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
