package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDeviceProfiles       = "/api/device-profiles"
	apiUrlDeviceProfilesPrefix = apiUrlDeviceProfiles + apiUrlPathDelim
	apiUrlDeviceProfileById    = apiUrlDeviceProfilesPrefix + "%s"
)

type optionsDeviceProfilessResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type HardwareCapabilities struct {
	FormFactor     string `json:"form_factor"`
	Cpu            string `json:"cpu"`
	Ram            int    `json:"ram"`
	Asic           string `json:"asic"`
	MaxL2Mtu       int    `json:"max_l2_mtu"`
	MaxL3Mtu       int    `json:"max_l3_mtu"`
	Userland       int    `json:"userland"`
	VtepLimit      int    `json:"vtep_limit"`
	BfdSupported   bool   `json:"bfd_supported"`
	VxlanSupported bool   `json:"vxlan_supported"`
	VtepFloodLimit int    `json:"vtep_flood_limit"`
	EcmpLimit      int    `json:"ecmp_limit"`
	VrfLimit       int    `json:"vrf_limit"`
	CoppStrict     []struct {
		Version string `json:"version"`
		Value   bool   `json:"value"`
	} `json:"copp_strict"`
	BreakoutCapable []struct {
		Version string `json:"version"`
		Value   bool   `json:"value"`
		Module  int    `json:"module"`
	} `json:"breakout_capable"`
	RoutingInstanceSupported []struct {
		Version string `json:"version"`
		Value   bool   `json:"value"`
	} `json:"routing_instance_supported"`
	AsSeqNumSupported []struct {
		Version string `json:"version"`
		Value   bool   `json:"value"`
	} `json:"as_seq_num_supported"`
}

type SoftwareCapabilities struct {
	Onie               bool   `json:"onie"`
	ConfigApplySupport string `json:"config_apply_support"`
	LxcSupport         bool   `json:"lxc_support"`
}

type DeviceSelector struct {
	OsVersion    string `json:"os_version"`
	Model        string `json:"model"`
	Os           string `json:"os"`
	Manufacturer string `json:"manufacturer"`
}

type PortInfo struct {
	DisplayId       int
	PanelId         int
	SlotId          int
	ConnectorType   string
	RowId           int
	Transformations []Transformation
	PortId          int
	FailureDomainId int
	ColumnId        int
}

func (o *PortInfo) raw() *rawPortInfo {
	transformations := make([]rawTransformation, len(o.Transformations))
	for i := range o.Transformations {
		transformations[i] = *o.Transformations[i].raw()
	}
	return &rawPortInfo{
		DisplayId:       o.DisplayId,
		PanelId:         o.PanelId,
		SlotId:          o.SlotId,
		ConnectorType:   o.ConnectorType,
		RowId:           o.RowId,
		Transformations: transformations,
		PortId:          o.PortId,
		FailureDomainId: o.FailureDomainId,
		ColumnId:        o.ColumnId,
	}
}

type rawPortInfo struct {
	DisplayId       int                 `json:"display_id"`
	PanelId         int                 `json:"panel_id"`
	SlotId          int                 `json:"slot_id"`
	ConnectorType   string              `json:"connector_type"`
	RowId           int                 `json:"row_id"`
	Transformations []rawTransformation `json:"transformations"`
	PortId          int                 `json:"port_id"`
	FailureDomainId int                 `json:"failure_domain_id"`
	ColumnId        int                 `json:"column_id"`
}

func (o *rawPortInfo) polish() *PortInfo {
	transformations := make([]Transformation, len(o.Transformations))
	for i := range o.Transformations {
		transformations[i] = *o.Transformations[i].polish()
	}
	return &PortInfo{
		DisplayId:       o.DisplayId,
		PanelId:         o.PanelId,
		SlotId:          o.SlotId,
		ConnectorType:   o.ConnectorType,
		RowId:           o.RowId,
		Transformations: transformations,
		PortId:          o.PortId,
		FailureDomainId: o.FailureDomainId,
		ColumnId:        o.ColumnId,
	}
}

type Transformation struct {
	IsDefault        bool
	Interfaces       []TransformInterface
	TransformationId int
}

func (o *Transformation) raw() *rawTransformation {
	interfaces := make([]rawTransformInterface, len(o.Interfaces))
	for i := range o.Interfaces {
		interfaces[i] = *o.Interfaces[i].raw()
	}
	return &rawTransformation{
		IsDefault:        o.IsDefault,
		Interfaces:       interfaces,
		TransformationId: o.TransformationId,
	}
}

type rawTransformation struct {
	IsDefault        bool                    `json:"is_default"`
	Interfaces       []rawTransformInterface `json:"interfaces"`
	TransformationId int                     `json:"transformation_id"`
}

func (o *rawTransformation) polish() *Transformation {
	interfaces := make([]TransformInterface, len(o.Interfaces))
	for i := range o.Interfaces {
		interfaces[i] = *o.Interfaces[i].polish()
	}
	return &Transformation{
		IsDefault:        o.IsDefault,
		Interfaces:       interfaces,
		TransformationId: o.TransformationId,
	}
}

type TransformInterface struct {
	State       string
	Setting     string
	Speed       LogicalDevicePortSpeed
	Name        string
	InterfaceId int
}

func (o *TransformInterface) raw() *rawTransformInterface {
	return &rawTransformInterface{
		State:       o.State,
		Setting:     o.Setting,
		Speed:       *o.Speed.raw(),
		Name:        o.Name,
		InterfaceId: o.InterfaceId,
	}
}

type rawTransformInterface struct {
	State       string
	Setting     string
	Speed       rawLogicalDevicePortSpeed
	Name        string
	InterfaceId int
}

func (o *rawTransformInterface) polish() *TransformInterface {
	return &TransformInterface{
		State:       o.State,
		Setting:     o.Setting,
		Speed:       o.Speed.parse(),
		Name:        o.Name,
		InterfaceId: o.InterfaceId,
	}
}

type DeviceProfile struct {
	Id                   ObjectId
	Label                string
	DeviceProfileType    string
	CreatedAt            time.Time
	LastModifiedAt       time.Time
	ChassisProfileId     string
	ChassisCount         int
	SlotCount            int
	HardwareCapabilities HardwareCapabilities
	SoftwareCapabilities SoftwareCapabilities
	Ports                []PortInfo
	Selector             DeviceSelector
	ChassisInfo          DeviceProfileChassisInfo
	LinecardsInfo        []DeviceProfileLinecardInfo
	SlotConfiguration    []DeviceProfileSlotConfiguration
}

type DeviceProfileChassisInfo struct {
	ChassisProfileId     string               `json:"chassis_profile_id"`
	HardwareCapabilities HardwareCapabilities `json:"hardware_capabilities"`
	SoftwareCapabilities SoftwareCapabilities `json:"software_capabilities"`
	Selector             DeviceSelector       `json:"selector"`
}

type DeviceProfileLinecardInfo struct {
	HardwareCapabilities HardwareCapabilities `json:"hardware_capabilities"`
	LinecardProfileId    string               `json:"linecard_profile_id"`
	Selector             DeviceSelector       `json:"selector"`
}

type DeviceProfileSlotConfiguration struct {
	LinecardProfileId string `json:"linecard_profile_id"`
	SlotId            int    `json:"slot_id"`
}

func (o *DeviceProfile) raw() *rawDeviceProfile {
	ports := make([]rawPortInfo, len(o.Ports))
	for i := range o.Ports {
		ports[i] = *o.Ports[i].raw()
	}
	return &rawDeviceProfile{
		Id:                   o.Id,
		Label:                o.Label,
		DeviceProfileType:    o.DeviceProfileType,
		ChassisProfileId:     o.ChassisProfileId,
		ChassisCount:         o.ChassisCount,
		SlotCount:            o.SlotCount,
		HardwareCapabilities: o.HardwareCapabilities,
		SoftwareCapabilities: o.SoftwareCapabilities,
		Ports:                ports,
		Selector:             o.Selector,
		ChassisInfo:          o.ChassisInfo,
		LinecardsInfo:        o.LinecardsInfo,
		SlotConfiguration:    nil,
	}
}

type rawDeviceProfile struct {
	Id                   ObjectId                         `json:"id"`
	Label                string                           `json:"label"`
	DeviceProfileType    string                           `json:"device_profile_type"`
	CreatedAt            time.Time                        `json:"created_at"`
	LastModifiedAt       time.Time                        `json:"last_modified_at"`
	ChassisProfileId     string                           `json:"chassis_profile_id"`
	ChassisCount         int                              `json:"chassis_count"`
	SlotCount            int                              `json:"slot_count"`
	HardwareCapabilities HardwareCapabilities             `json:"hardware_capabilities"`
	SoftwareCapabilities SoftwareCapabilities             `json:"software_capabilities"`
	Ports                []rawPortInfo                    `json:"ports"`
	Selector             DeviceSelector                   `json:"selector"`
	ChassisInfo          DeviceProfileChassisInfo         `json:"chassis_info"`
	LinecardsInfo        []DeviceProfileLinecardInfo      `json:"linecards_info"`
	SlotConfiguration    []DeviceProfileSlotConfiguration `json:"slot_configuration"`
}

func (o *rawDeviceProfile) polish() *DeviceProfile {
	ports := make([]PortInfo, len(o.Ports))
	for i := range o.Ports {
		ports[i] = *o.Ports[i].polish()
	}
	return &DeviceProfile{
		Id:                   o.Id,
		Label:                o.Label,
		DeviceProfileType:    o.DeviceProfileType,
		CreatedAt:            o.CreatedAt,
		LastModifiedAt:       o.LastModifiedAt,
		ChassisProfileId:     o.ChassisProfileId,
		ChassisCount:         o.ChassisCount,
		SlotCount:            o.SlotCount,
		HardwareCapabilities: o.HardwareCapabilities,
		SoftwareCapabilities: o.SoftwareCapabilities,
		Ports:                ports,
		Selector:             o.Selector,
		ChassisInfo:          o.ChassisInfo,
		LinecardsInfo:        o.LinecardsInfo,
		SlotConfiguration:    o.SlotConfiguration,
	}
}

// TransformationCandidates takes an interface name ("xe-0/0/1:1") and a speed,
// and returns a map[int][]Transformation keyed by PortId. Only "active"
// transformations matching the specified interface name and speed are returned.
func (o *DeviceProfile) TransformationCandidates(intfName string, intfSpeed LogicalDevicePortSpeed) map[int][]Transformation {
	var result []Transformation
	for _, port := range o.Ports {
		for _, transformation := range port.Transformations {
			for _, intf := range transformation.Interfaces {
				if intf.Name == intfName &&
					intf.State == "active" &&
					intf.Speed.IsEqual(intfSpeed) {
					result = append(result, transformation)
				}
			}
		}
	}
	return result
}

func (o *Client) listDeviceProfileIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsDeviceProfilessResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDeviceProfiles,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getAllDeviceProfiles(ctx context.Context) ([]rawDeviceProfile, error) {
	response := &struct {
		Items []rawDeviceProfile
	}{}

	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDeviceProfiles,
		apiResponse: response,
	})
	return response.Items, convertTtaeToAceWherePossible(err)
}

func (o *Client) getDeviceProfilesByName(ctx context.Context, desired string) ([]rawDeviceProfile, error) {
	deviceProfiles, err := o.getAllDeviceProfiles(ctx)
	if err != nil {
		return nil, err
	}
	var result []rawDeviceProfile
	for _, deviceProfile := range deviceProfiles {
		if deviceProfile.Label == desired {
			result = append(result, deviceProfile)
		}
	}
	return result, nil
}

func (o *Client) getDeviceProfileByName(ctx context.Context, desired string) (*rawDeviceProfile, error) {
	deviceProfiles, err := o.getDeviceProfilesByName(ctx, desired)
	if err != nil {
		return nil, err
	}
	switch len(deviceProfiles) {
	case 0:
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no device profile named '%s' found", desired),
		}
	case 1:
		return &deviceProfiles[0], nil
	default:
		return nil, ApstraClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("found multiple device profiles named '%s'", desired),
		}
	}
}

func (o *Client) getDeviceProfile(ctx context.Context, id ObjectId) (*rawDeviceProfile, error) {
	response := &rawDeviceProfile{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDeviceProfileById, id),
		apiResponse: response,
	})
	return response, convertTtaeToAceWherePossible(err)
}

func (o *Client) createDeviceProfile(ctx context.Context, profile *rawDeviceProfile) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDeviceProfiles,
		apiInput:    profile,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateDeviceProfile(ctx context.Context, id ObjectId, profile *rawDeviceProfile) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDeviceProfileById, id),
		apiInput: profile,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteDeviceProfile(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDeviceProfileById, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
