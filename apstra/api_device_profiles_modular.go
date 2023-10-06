package apstra

import (
	"context"
	"fmt"
	"github.com/orsinium-labs/enum"
	"net/http"
)

type DeviceProfileType enum.Member[string]

var (
	DeviceProfileTypeModular    = DeviceProfileType{Value: "modular"}
	DeviceProfileTypeMonolithic = DeviceProfileType{Value: "monolithic"}
)

// https://35.92.136.236:47759/api/device-profiles
//
// POST
// {
//  "device_profile_type": "modular",
//  "predefined": false,
//  "label": "name",
//  "chassis_profile_id": "Juniper_PTX10008",
//  "slot_configuration": [
//    {
//      "linecard_profile_id": "Juniper_PTX10K_LC1201_36CD",
//      "slot_id": 0
//    }
//  ]
//}
//
// PUT
// {
//  "device_profile_type": "modular",
//  "predefined": false,
//  "id": "dfcfc410-6997-4825-8cb8-41fe8e78d528",
//  "label": "LABEL GOES HERE2",
//  "chassis_profile_id": "Juniper_PTX10008",
//  "slot_configuration": [
//    {
//      "linecard_profile_id": "Juniper_PTX10K_LC1201_36CD",
//      "slot_id": 0
//    }
//  ]
//}

type rawModularDeviceSlotConfiguration struct {
	LinecardProfileId ObjectId `json:"linecard_profile_id"`
	SlotId            uint64   `json:"slot_id"`
}

type ModularDeviceSlotConfiguration struct {
	LinecardProfileId ObjectId
}

type ModularDeviceProfile struct {
	Label             string
	ChassisProfileId  ObjectId
	SlotConfiguration map[uint64]ModularDeviceSlotConfiguration
}

func (o *ModularDeviceProfile) raw() *rawModularDeviceProfile {
	result := &rawModularDeviceProfile{
		DeviceProfileType: DeviceProfileTypeModular.Value,
		Label:             o.Label,
		ChassisProfileId:  o.ChassisProfileId,
		SlotConfiguration: make([]rawModularDeviceSlotConfiguration, len(o.SlotConfiguration)),
	}

	var i int
	for slotId, slotConfiguraiton := range o.SlotConfiguration {
		result.SlotConfiguration[i] = rawModularDeviceSlotConfiguration{
			SlotId:            slotId,
			LinecardProfileId: slotConfiguraiton.LinecardProfileId,
		}
		i++
	}

	return result
}

type rawModularDeviceProfile struct {
	DeviceProfileType string                              `json:"device_profile_type"`
	Label             string                              `json:"label"`
	ChassisProfileId  ObjectId                            `json:"chassis_profile_id"`
	SlotConfiguration []rawModularDeviceSlotConfiguration `json:"slot_configuration"`
}

func (o *rawModularDeviceProfile) polish() *ModularDeviceProfile {
	result := &ModularDeviceProfile{
		Label:             o.Label,
		ChassisProfileId:  o.ChassisProfileId,
		SlotConfiguration: make(map[uint64]ModularDeviceSlotConfiguration, len(o.SlotConfiguration)),
	}

	for _, slotConfiguration := range o.SlotConfiguration {
		result.SlotConfiguration[slotConfiguration.SlotId] = ModularDeviceSlotConfiguration{
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
