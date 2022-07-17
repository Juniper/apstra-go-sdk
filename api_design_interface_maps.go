package goapstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	apiUrlDesignInterfaceMaps       = apiUrlDesignPrefix + "interface-maps"
	apiUrlDesignInterfaceMapsPrefix = apiUrlDesignInterfaceMaps + apiUrlPathDelim
	apiUrlDesignInterfaceMapById    = apiUrlDesignInterfaceMapsPrefix + "%s"

	rawInterfaceStateTrue  = rawInterfaceState("active")
	rawInterfaceStateFalse = rawInterfaceState("inactive")
)

// rawInterfaceMapInterface.Setting.Param is a string containing JSON like this.
// it needs double quotes escaped {\"like\": \"this\"}.
// {
//  "global": {
//    "breakout": false,
//    "fpc": 0,
//    "pic": 0,
//    "port": 0,
//    "speed": "100g"
//  },
//  "interface": {
//    "speed": ""
//  }
//}

//     'mapping': s.List(s.Optional(s.Integer()),
//                      validate=[
//                          s.Length(exact=5),
//                          validates_first_three_entries_are_always_non_none],
//                      description='This list of 5 integers represent which '
//                      '(port ID, transformation ID and interface ID) in the '
//                      'device profile and which '
//                      '(panel ID, port ID) in the logical device '
//                      'is this interface coming from')})

type optionsInterfaceMapsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type getAllInterfaceMapsResponse struct {
	Items []rawInterfaceMap `json:"items"`
}

type InterfaceSettingParam struct {
	Global struct {
		Breakout bool   `json:"breakout"`
		Fpc      int    `json:"fpc"`
		Pic      int    `json:"pic"`
		Port     int    `json:"port"`
		Speed    string `json:"speed"`
	} `json:"global"`
	Interface struct {
		Speed string `json:"speed"`
	} `json:"interface"`
}

func (o InterfaceSettingParam) String() string {
	// medium confident we won't provoke UnsupportedTypeError or
	// UnsupportedValueError here, so ignoring the error return.
	payload, _ := json.Marshal(o)
	return strings.Replace(string(payload), `"`, `\"`, -1)
}

type InterfaceMapping struct {
	DPPortId      int
	DPTransformId int
	DPInterfaceId int
	LDPanel       int
	LDPort        int
}

func (o *InterfaceMapping) raw() *rawInterfaceMapping {
	return &rawInterfaceMapping{o.DPPortId, o.DPTransformId, o.DPInterfaceId, o.LDPanel, o.LDPort}
}

type rawInterfaceMapping []int

func (o rawInterfaceMapping) polish() *InterfaceMapping {
	return &InterfaceMapping{DPPortId: o[0], DPTransformId: o[1], DPInterfaceId: o[2], LDPanel: o[3], LDPort: o[4]}
}

type InterfaceStateActive bool

func (o InterfaceStateActive) raw() rawInterfaceState {
	if o {
		return rawInterfaceStateTrue
	} else {
		return rawInterfaceStateFalse
	}
}

type rawInterfaceState string

func (o rawInterfaceState) polish() (InterfaceStateActive, error) {
	switch o {
	case rawInterfaceStateTrue:
		return true, nil
	case rawInterfaceStateFalse:
		return false, nil
	default:
		return false, fmt.Errorf("unknown interface state '%s'", o)
	}
}

type InterfaceMapInterfaceSetting struct {
	Param string `json:"param"`
}

type InterfaceMapInterface struct {
	Name        string                       `json:"name"`
	Roles       LogicalDevicePortRoleFlags   `json:"roles"`
	Mapping     InterfaceMapping             `json:"mapping"`
	ActiveState InterfaceStateActive         `json:"state"`
	Position    int                          `json:"position"`
	Speed       LogicalDevicePortSpeed       `json:"speed"`
	Setting     InterfaceMapInterfaceSetting `json:"setting"`
}

func (o *InterfaceMapInterface) raw() *rawInterfaceMapInterface {
	return &rawInterfaceMapInterface{
		Name:     o.Name,
		Roles:    o.Roles.raw(),
		Mapping:  *o.Mapping.raw(),
		State:    o.ActiveState.raw(),
		Setting:  o.Setting,
		Position: o.Position,
		Speed:    o.Speed,
	}
}

type rawInterfaceMapInterface struct {
	Name     string                       `json:"name"`
	Roles    logicalDevicePortRoles       `json:"roles"`
	Mapping  rawInterfaceMapping          `json:"mapping"`
	State    rawInterfaceState            `json:"state"`
	Setting  InterfaceMapInterfaceSetting `json:"setting"`
	Position int                          `json:"position"`
	Speed    LogicalDevicePortSpeed       `json:"speed"`
}

func (o *rawInterfaceMapInterface) polish() (*InterfaceMapInterface, error) {
	roles, err := o.Roles.parse()
	if err != nil {
		return nil, err
	}
	state, err := o.State.polish()
	if err != nil {
		return nil, err
	}

	return &InterfaceMapInterface{
		Name:        o.Name,
		Roles:       roles,
		Mapping:     *o.Mapping.polish(),
		ActiveState: state,
		Position:    o.Position,
		Speed:       o.Speed,
		Setting:     o.Setting,
	}, nil
}

type InterfaceMap struct {
	LogicalDeviceId ObjectId
	DeviceProfileId ObjectId
	CreatedAt       time.Time
	LastModifiedAt  time.Time
	Id              ObjectId
	Label           string
	Interfaces      []InterfaceMapInterface
}

type rawInterfaceMap struct {
	LogicalDeviceId ObjectId                   `json:"logical_device_id"`
	DeviceProfileId ObjectId                   `json:"device_profile_id"`
	CreatedAt       time.Time                  `json:"created_at"`
	LastModifiedAt  time.Time                  `json:"last_modified_at"`
	Id              ObjectId                   `json:"id"`
	Label           string                     `json:"label"`
	Interfaces      []rawInterfaceMapInterface `json:"interfaces"`
}

func (o *rawInterfaceMap) polish() (*InterfaceMap, error) {
	interfaces := make([]InterfaceMapInterface, len(o.Interfaces))
	for i, intf := range o.Interfaces {
		polished, err := intf.polish()
		if err != nil {
			return nil, err
		}
		interfaces[i] = *polished
	}
	return &InterfaceMap{
		LogicalDeviceId: o.LogicalDeviceId,
		DeviceProfileId: o.DeviceProfileId,
		CreatedAt:       o.CreatedAt,
		LastModifiedAt:  o.LastModifiedAt,
		Id:              o.Id,
		Label:           o.Label,
		Interfaces:      interfaces,
	}, nil
}

func (o *Client) listAllInterfaceMapIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsInterfaceMapsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignInterfaceMaps,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (o *Client) GetInterfaceMap(ctx context.Context, id ObjectId) (*InterfaceMap, error) {
	response := &rawInterfaceMap{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignInterfaceMapById, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, err
	}
	return response.polish()
}
