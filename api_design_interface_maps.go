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

type InterfaceMapMapping struct {
	DPPortId      int
	DPTransformId int
	DPInterfaceId int
	LDPanel       int
	LDPort        int
}

func (o *InterfaceMapMapping) raw() *rawInterfaceMapping {
	var dPPortIdPtr, dPTransformIdPtr, dPInterfaceIdPtr, lDPanelPtr, lDPortPtr *int
	var dPPortId, dPTransformId, dPInterfaceId, lDPanel, lDPort int

	// device profile port/interface/transform IDs all seem to be 1-indexed
	// logical device panel/port IDs all seem to be 1-indexed
	// still using "-1" as the invalid/null here for safer compatibility with future versions
	if o.DPPortId < 0 {
		dPPortIdPtr = nil
	} else {
		dPPortId = o.DPPortId
		dPPortIdPtr = &dPPortId
	}

	if o.DPInterfaceId < 0 {
		dPInterfaceIdPtr = nil
	} else {
		dPInterfaceId = o.DPInterfaceId
		dPInterfaceIdPtr = &dPInterfaceId
	}

	if o.DPTransformId < 0 {
		dPTransformIdPtr = nil
	} else {
		dPTransformId = o.DPTransformId
		dPTransformIdPtr = &dPTransformId
	}

	if o.LDPanel < 0 {
		lDPanelPtr = nil
	} else {
		lDPanel = o.LDPanel
		lDPanelPtr = &lDPanel
	}

	if o.LDPort < 0 {
		lDPortPtr = nil
	} else {
		lDPort = o.LDPort
		lDPortPtr = &lDPort
	}

	return &rawInterfaceMapping{dPPortIdPtr, dPTransformIdPtr, dPInterfaceIdPtr, lDPanelPtr, lDPortPtr}
}

type rawInterfaceMapping []*int

func (o rawInterfaceMapping) polish() *InterfaceMapMapping {
	dPPortId, dPTransformId, dPInterfaceId, lDPanel, lDPort := -1, -1, -1, -1, -1
	if o[0] != nil {
		dPPortId = *o[0]
	}
	if o[1] != nil {
		dPTransformId = *o[0]
	}
	if o[2] != nil {
		dPInterfaceId = *o[0]
	}
	if o[3] != nil {
		lDPanel = *o[0]
	}
	if o[4] != nil {
		lDPort = *o[0]
	}

	return &InterfaceMapMapping{
		DPPortId:      dPPortId,
		DPTransformId: dPTransformId,
		DPInterfaceId: dPInterfaceId,
		LDPanel:       lDPanel,
		LDPort:        lDPort,
	}
}

type InterfaceStateActive bool

func (o InterfaceStateActive) raw() rawInterfaceState {
	if o {
		return rawInterfaceStateTrue
	}
	return rawInterfaceStateFalse
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
	Name        string
	Roles       LogicalDevicePortRoleFlags
	Mapping     InterfaceMapMapping
	ActiveState InterfaceStateActive
	Position    int
	Speed       LogicalDevicePortSpeed
	Setting     InterfaceMapInterfaceSetting
}

func (o *InterfaceMapInterface) raw() *rawInterfaceMapInterface {
	return &rawInterfaceMapInterface{
		Name:     o.Name,
		Roles:    o.Roles.raw(),
		Mapping:  *o.Mapping.raw(),
		State:    o.ActiveState.raw(),
		Setting:  o.Setting,
		Position: o.Position,
		Speed:    *o.Speed.raw(),
	}
}

type rawInterfaceMapInterface struct {
	Name     string                       `json:"name"`
	Roles    logicalDevicePortRoles       `json:"roles"`
	Mapping  rawInterfaceMapping          `json:"mapping"`
	State    rawInterfaceState            `json:"state"`
	Setting  InterfaceMapInterfaceSetting `json:"setting"`
	Position int                          `json:"position"`
	Speed    rawLogicalDevicePortSpeed    `json:"speed"`
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
		Speed:       o.Speed.parse(),
		Setting:     o.Setting,
	}, nil
}

type InterfaceMapData struct {
	LogicalDeviceId ObjectId
	DeviceProfileId ObjectId
	Label           string
	Interfaces      []InterfaceMapInterface
}

type InterfaceMap struct {
	Id             ObjectId
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Data           *InterfaceMapData
}

func (o *InterfaceMapData) raw() *rawInterfaceMap {
	rawInterfaces := make([]rawInterfaceMapInterface, len(o.Interfaces))
	for i, intf := range o.Interfaces {
		rawInterfaces[i] = *intf.raw()
	}

	return &rawInterfaceMap{
		LogicalDeviceId: o.LogicalDeviceId,
		DeviceProfileId: o.DeviceProfileId,
		Label:           o.Label,
		Interfaces:      rawInterfaces,
	}
}

type rawInterfaceMap struct {
	LogicalDeviceId ObjectId                   `json:"logical_device_id"`
	DeviceProfileId ObjectId                   `json:"device_profile_id"`
	CreatedAt       time.Time                  `json:"created_at"`
	LastModifiedAt  time.Time                  `json:"last_modified_at"`
	Id              ObjectId                   `json:"id,omitempty"`
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
		Id:             o.Id,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &InterfaceMapData{
			LogicalDeviceId: o.LogicalDeviceId,
			DeviceProfileId: o.DeviceProfileId,
			Label:           o.Label,
			Interfaces:      interfaces,
		},
	}, nil
}

func (o *Client) listAllInterfaceMapIds(ctx context.Context) ([]ObjectId, error) {
	response := &struct {
		Items []ObjectId `json:"items"`
	}{}
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

func (o *Client) getAllInterfaceMaps(ctx context.Context) ([]rawInterfaceMap, error) {
	response := &struct {
		Items []rawInterfaceMap `json:"items"`
	}{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignInterfaceMaps,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getInterfaceMap(ctx context.Context, id ObjectId) (*rawInterfaceMap, error) {
	response := &rawInterfaceMap{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignInterfaceMapById, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getInterfaceMapsByName(ctx context.Context, desired string) ([]rawInterfaceMap, error) {
	interfaceMaps, err := o.getAllInterfaceMaps(ctx)
	if err != nil {
		return nil, err
	}
	var result []rawInterfaceMap
	for _, interfaceMap := range interfaceMaps {
		if interfaceMap.Label == desired {
			result = append(result, interfaceMap)
		}
	}
	return result, nil
}

func (o *Client) getInterfaceMapByName(ctx context.Context, desired string) (*rawInterfaceMap, error) {
	interfaceMaps, err := o.getInterfaceMapsByName(ctx, desired)
	if err != nil {
		return nil, err
	}
	switch len(interfaceMaps) {
	case 0:
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no interface map named '%s' found", desired),
		}
	case 1:
		return &interfaceMaps[0], nil
	default:
		return nil, ApstraClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found multiple interface maps named '%s' found", desired),
		}
	}
}

func (o *Client) createInterfaceMap(ctx context.Context, in *InterfaceMapData) (ObjectId, error) {
	response := &objectIdResponse{}
	return response.Id, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignInterfaceMaps,
		apiInput:    in.raw(),
		apiResponse: response,
	})
}

func (o *Client) updateInterfaceMap(ctx context.Context, id ObjectId, in *InterfaceMapData) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignInterfaceMapById, id),
		apiInput: in.raw(),
	})
}

func (o *Client) deleteInterfaceMap(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignInterfaceMapById, id),
	})
}
