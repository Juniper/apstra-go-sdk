// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"encoding/json"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/internal"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/speed"
)

var _ internal.IDer = (*Profile)(nil)
var _ json.Unmarshaler = (*Profile)(nil)

type Profile struct {
	Selector                    Selector                     `json:"selector"`
	DeviceProfileType           enum.DeviceProfileType       `json:"device_profile_type"`
	DualRoutingEngine           bool                         `json:"dual_routing_engine"`
	SoftwareCapabilities        SoftwareCapabilities         `json:"software_capabilities"`
	ReferenceDesignCapabilities *ReferenceDesignCapabilities `json:"reference_design_capabilities,omitempty"` // introduced between 4.2.0 and 6.0.0
	ChassisProfileID            string                       `json:"chassis_profile_id,omitempty"`
	HardwareCapabilities        HardwareCapabilities         `json:"hardware_capabilities"`
	Predefined                  bool                         `json:"predefined"`
	SlotCount                   *int                         `json:"slot_count,omitempty"`
	ChassisCount                *int                         `json:"chassis_count,omitempty"`
	Ports                       []Port                       `json:"ports"`
	Label                       string                       `json:"label"`
	ChassisInfo                 *ProfileChassisInfo          `json:"chassis_info,omitempty"`
	LinecardsInfo               []ProfileLinecardInfo        `json:"linecards_info,omitempty"`
	SlotConfiguration           []ProfileSlotConfiguration   `json:"slot_configuration,omitempty"`
	PhysicalDevice              bool                         `json:"physical_device"`

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (p Profile) ID() *string {
	if p.id == "" {
		return nil
	}
	return &p.id
}

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (p *Profile) SetID(id string) error {
	if p.id != "" {
		return internal.IDIsSet(fmt.Errorf("id already has value %q", p.id))
	}

	p.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (p *Profile) MustSetID(id string) {
	err := p.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (p *Profile) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Selector                    Selector                     `json:"selector"`
		DeviceProfileType           enum.DeviceProfileType       `json:"device_profile_type"`
		DualRoutingEngine           bool                         `json:"dual_routing_engine"`
		SoftwareCapabilities        SoftwareCapabilities         `json:"software_capabilities"`
		ReferenceDesignCapabilities *ReferenceDesignCapabilities `json:"reference_design_capabilities"` // introduced between 4.2.0 and 6.0.0
		ChassisProfileID            string                       `json:"chassis_profile_id,omitempty"`
		HardwareCapabilities        HardwareCapabilities         `json:"hardware_capabilities"`
		Predefined                  bool                         `json:"predefined"`
		SlotCount                   *int                         `json:"slot_count,omitempty"`
		ChassisCount                *int                         `json:"chassis_count,omitempty"`
		Ports                       []Port                       `json:"ports"`
		Label                       string                       `json:"label"`
		ChassisInfo                 *ProfileChassisInfo          `json:"chassis_info,omitempty"`
		LinecardsInfo               []ProfileLinecardInfo        `json:"linecards_info,omitempty"`
		SlotConfiguration           []ProfileSlotConfiguration   `json:"slot_configuration,omitempty"`
		PhysicalDevice              bool                         `json:"physical_device"`

		Id             string     `json:"id"`
		CreatedAt      *time.Time `json:"created_at,omitempty"`
		LastModifiedAt *time.Time `json:"last_modified_at,omitempty"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling profile: %w", err)
	}

	p.Selector = raw.Selector
	p.DeviceProfileType = raw.DeviceProfileType
	p.DualRoutingEngine = raw.DualRoutingEngine
	p.SoftwareCapabilities = raw.SoftwareCapabilities
	p.ReferenceDesignCapabilities = raw.ReferenceDesignCapabilities
	p.ChassisProfileID = raw.ChassisProfileID
	p.HardwareCapabilities = raw.HardwareCapabilities
	p.Predefined = raw.Predefined
	p.SlotCount = raw.SlotCount
	p.ChassisCount = raw.ChassisCount
	p.Ports = raw.Ports
	p.Label = raw.Label
	p.ChassisInfo = raw.ChassisInfo
	p.LinecardsInfo = raw.LinecardsInfo
	p.SlotConfiguration = raw.SlotConfiguration
	p.PhysicalDevice = raw.PhysicalDevice
	p.id = raw.Id
	p.createdAt = raw.CreatedAt
	p.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func (p Profile) CreatedAt() *time.Time {
	return p.createdAt
}

func (p Profile) LastModifiedAt() *time.Time {
	return p.lastModifiedAt
}

type HardwareCapabilities struct {
	MaxL3Mtu          *int             `json:"max_l3_mtu,omitempty"`
	MaxL2Mtu          *int             `json:"max_l2_mtu,omitempty"`
	FormFactor        string           `json:"form_factor"`
	VTEPLimit         *int             `json:"vtep_limit,omitempty"`
	BFDSupported      bool             `json:"bfd_supported,omitempty"`
	COPPStrict        []FeatureVersion `json:"copp_strict,omitempty"`
	ECMPLimit         int              `json:"ecmp_limit"`
	AsSeqNumSupported []FeatureVersion `json:"as_seq_num_supported,omitempty"`
	Ram               int              `json:"ram"`
	VTEPFloodLimit    *int             `json:"vtep_flood_limit,omitempty"`
	BreakoutCapable   []struct {
		Module  int    `json:"module"`
		Value   bool   `json:"value"`
		Version string `json:"version"` // introduced and became mandatory some time between 4.2.0 and 6.0.0
	} `json:"breakout_capable,omitempty"`
	Userland                 int              `json:"userland"`
	Asic                     string           `json:"asic"`
	VrfLimit                 *int             `json:"vrf_limit,omitempty"`
	RoutingInstanceSupported []FeatureVersion `json:"routing_instance_supported,omitempty"`
	VxlanSupported           bool             `json:"vxlan_supported,omitempty"`
	CPU                      string           `json:"cpu"`
}

// FeatureVersion details whether a feature is enabled on the given NOS version
type FeatureVersion struct {
	Version string `json:"version"`
	Enabled bool   `json:"value"`
}

type SoftwareCapabilities struct {
	Onie               bool   `json:"onie"`
	ConfigApplySupport string `json:"config_apply_support"`
	LxcSupport         bool   `json:"lxc_support"`
}

type ReferenceDesignCapabilities struct {
	Datacenter enum.RefDesignCapability `json:"datacenter"`
	Freeform   enum.RefDesignCapability `json:"freeform"`
}

type Port struct {
	ConnectorType   string           `json:"connector_type"`
	Panel           int              `json:"panel_id"`
	Transformations []Transformation `json:"transformations"`
	Column          int              `json:"column_id"`
	Port            int              `json:"port_id"`
	Row             int              `json:"row_id"`
	FailureDomain   int              `json:"failure_domain_id"`
	Display         *int             `json:"display_id"`
	Slot            int              `json:"slot_id"`
}

type Transformation struct {
	ID         int                       `json:"transformation_id"`
	IsDefault  bool                      `json:"is_default"`
	Interfaces []TransformationInterface `json:"interfaces"`
}

type TransformationInterface struct {
	ID      int                 `json:"interface_id"`
	Name    string              `json:"name"`
	State   enum.InterfaceState `json:"state"`
	Setting *string             `json:"setting"`
	Speed   speed.Speed         `json:"speed"`
}

type Selector struct {
	OsVersion    string `json:"os_version"`
	Model        string `json:"model"`
	Os           string `json:"os"`
	Manufacturer string `json:"manufacturer"`
}

type ProfileChassisInfo struct {
	ID                          string                       `json:"chassis_profile_id"`
	Selector                    Selector                     `json:"selector"`
	HardwareCapabilities        HardwareCapabilities         `json:"hardware_capabilities"`
	SoftwareCapabilities        SoftwareCapabilities         `json:"software_capabilities"`
	DualRoutingEngine           bool                         `json:"dual_routing_engine"`
	LinecardSlotIDs             []int                        `json:"linecard_slot_ids"`
	PhysicalDevice              bool                         `json:"physical_device"`
	ReferenceDesignCapabilities *ReferenceDesignCapabilities `json:"reference_design_capabilities"` // introduced between 4.2.0 and 6.0.0
}

type ProfileLinecardInfo struct {
	ID                   string               `json:"linecard_profile_id"`
	Selector             Selector             `json:"selector"`
	HardwareCapabilities HardwareCapabilities `json:"hardware_capabilities"`
}

type ProfileSlotConfiguration struct {
	LinecardProfileID string `json:"linecard_profile_id"`
	SlotID            int    `json:"slot_id"`
}
