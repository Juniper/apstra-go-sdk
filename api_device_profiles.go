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

type getAllDeviceProfilesResponse struct {
	Items []DeviceProfile
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
	DisplayId       int    `json:"display_id"`
	PanelId         int    `json:"panel_id"`
	SlotId          int    `json:"slot_id"`
	ConnectorType   string `json:"connector_type"`
	RowId           int    `json:"row_id"`
	Transformations []struct {
		IsDefault  bool `json:"is_default"`
		Interfaces []struct {
			State   string `json:"state"`
			Setting string `json:"setting"`
			Speed   struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			} `json:"speed"`
			Name        string `json:"name"`
			InterfaceId int    `json:"interface_id"`
		} `json:"interfaces"`
		TransformationId int `json:"transformation_id"`
	} `json:"transformations"`
	PortId          int `json:"port_id"`
	FailureDomainId int `json:"failure_domain_id"`
	ColumnId        int `json:"column_id"`
}

type DeviceProfile struct {
	Id                   ObjectId             `json:"id"`
	Label                string               `json:"label"`
	DeviceProfileType    string               `json:"device_profile_type"`
	CreatedAt            time.Time            `json:"created_at"`
	LastModifiedAt       time.Time            `json:"last_modified_at"`
	ChassisProfileId     string               `json:"chassis_profile_id"`
	ChassisCount         int                  `json:"chassis_count"`
	SlotCount            int                  `json:"slot_count"`
	HardwareCapabilities HardwareCapabilities `json:"hardware_capabilities"`
	SoftwareCapabilities SoftwareCapabilities `json:"software_capabilities"`
	Ports                []PortInfo           `json:"ports"`
	Selector             DeviceSelector       `json:"selector"`
	ChassisInfo          struct {
		ChassisProfileId     string               `json:"chassis_profile_id"`
		HardwareCapabilities HardwareCapabilities `json:"hardware_capabilities"`
		SoftwareCapabilities SoftwareCapabilities `json:"software_capabilities"`
		Selector             DeviceSelector       `json:"selector"`
	} `json:"chassis_info"`
	LinecardsInfo []struct {
		HardwareCapabilities HardwareCapabilities `json:"hardware_capabilities"`
		LinecardProfileId    string               `json:"linecard_profile_id"`
		Selector             DeviceSelector       `json:"selector"`
	} `json:"linecards_info"`
	SlotConfiguration []struct {
		LinecardProfileId string `json:"linecard_profile_id"`
		SlotId            int    `json:"slot_id"`
	} `json:"slot_configuration"`
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

func (o *Client) getAllDeviceProfiles(ctx context.Context) ([]DeviceProfile, error) {
	response := &getAllDeviceProfilesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDeviceProfiles,
		apiResponse: response,
	})
	return response.Items, convertTtaeToAceWherePossible(err)
}

func (o *Client) getDeviceProfile(ctx context.Context, id ObjectId) (*DeviceProfile, error) {
	response := &DeviceProfile{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDeviceProfileById, id),
		apiResponse: response,
	})
	return response, convertTtaeToAceWherePossible(err)
}

func (o *Client) createDeviceProfile(ctx context.Context, profile DeviceProfile) (ObjectId, error) {
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

func (o *Client) updateDeviceProfile(ctx context.Context, id ObjectId, profile DeviceProfile) error {
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
