// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/speed"
)

type Profile struct {
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
	ChassisInfo                 ProfileChassisInfo           `json:"chassis_info"`
	LinecardsInfo               []ProfileLinecardInfo        `json:"linecards_info"`
	SlotConfiguration           []ProfileSlotConfiguration   `json:"slot_configuration"`
	PhysicalDevice              bool                         `json:"physical_device"`

	id             string     `json:"id"`
	createdAt      *time.Time `json:"created_at,omitempty"`
	lastModifiedAt *time.Time `json:"last_modified_at,omitempty"`
}

func (p Profile) ID() *string {
	if p.id == "" {
		return nil
	}
	return &p.id
}

func (p Profile) CreatedAt() *time.Time {
	return p.createdAt
}

func (p Profile) LastModifiedAt() *time.Time {
	return p.lastModifiedAt
}

type HardwareCapabilities struct {
	MaxL3Mtu          *int             `json:"max_l3_mtu"`
	MaxL2Mtu          *int             `json:"max_l2_mtu"`
	FormFactor        string           `json:"form_factor"`
	VTEPLimit         *int             `json:"vtep_limit"`
	BFDSupported      bool             `json:"bfd_supported"`
	COPPStrict        []FeatureVersion `json:"copp_strict"`
	ECMPLimit         int              `json:"ecmp_limit"`
	AsSeqNumSupported []FeatureVersion `json:"as_seq_num_supported"`
	Ram               int              `json:"ram"`
	VTEPFloodLimit    *int             `json:"vtep_flood_limit"`
	BreakoutCapable   []struct {
		Module  int    `json:"module"`
		Value   bool   `json:"value"`
		Version string `json:"version"` // introduced and became mandatory some time between 4.2.0 and 6.0.0
	} `json:"breakout_capable"`
	Userland                 int              `json:"userland"`
	Asic                     string           `json:"asic"`
	VrfLimit                 *int             `json:"vrf_limit"`
	RoutingInstanceSupported []FeatureVersion `json:"routing_instance_supported"`
	VxlanSupported           bool             `json:"vxlan_supported"`
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
