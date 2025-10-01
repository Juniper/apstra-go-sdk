// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/speed"
)

var (
	_ internal.IDer     = (*Profile)(nil)
	_ internal.IDSetter = (*Profile)(nil)
	_ json.Marshaler    = (*Profile)(nil)
	_ json.Unmarshaler  = (*Profile)(nil)
)

type Profile struct {
	Selector                    Selector
	DeviceProfileType           enum.DeviceProfileType
	DualRoutingEngine           bool
	SoftwareCapabilities        SoftwareCapabilities
	ReferenceDesignCapabilities *ReferenceDesignCapabilities // introduced between 5.0.0
	ChassisProfileID            string
	HardwareCapabilities        HardwareCapabilities
	Predefined                  bool
	SlotCount                   *int
	ChassisCount                *int
	Ports                       []Port
	Label                       string
	ChassisInfo                 *ProfileChassisInfo
	LinecardsInfo               []ProfileLinecardInfo
	SlotConfiguration           []ProfileSlotConfiguration
	PhysicalDevice              bool

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
		return sdk.ErrIDIsSet(fmt.Sprintf("id already has value %q", p.id))
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

func (p Profile) MarshalJSON() ([]byte, error) {
	raw := struct {
		Selector                    Selector                     `json:"selector"`
		DeviceProfileType           enum.DeviceProfileType       `json:"device_profile_type"`
		DualRoutingEngine           bool                         `json:"dual_routing_engine"`
		SoftwareCapabilities        SoftwareCapabilities         `json:"software_capabilities"`
		ReferenceDesignCapabilities *ReferenceDesignCapabilities `json:"reference_design_capabilities"` // introduced between 4.2.0 and 6.0.0
		ChassisProfileID            string                       `json:"chassis_profile_id,omitempty"`
		HardwareCapabilities        HardwareCapabilities         `json:"hardware_capabilities"`
		SlotCount                   *int                         `json:"slot_count,omitempty"`
		ChassisCount                *int                         `json:"chassis_count,omitempty"`
		Ports                       []Port                       `json:"ports"`
		Label                       string                       `json:"label"`
		ChassisInfo                 *ProfileChassisInfo          `json:"chassis_info,omitempty"`
		LinecardsInfo               []ProfileLinecardInfo        `json:"linecards_info,omitempty"`
		SlotConfiguration           []ProfileSlotConfiguration   `json:"slot_configuration,omitempty"`
		PhysicalDevice              bool                         `json:"physical_device"`
	}{
		Selector:                    p.Selector,
		DeviceProfileType:           p.DeviceProfileType,
		DualRoutingEngine:           p.DualRoutingEngine,
		SoftwareCapabilities:        p.SoftwareCapabilities,
		ReferenceDesignCapabilities: p.ReferenceDesignCapabilities,
		ChassisProfileID:            p.ChassisProfileID,
		HardwareCapabilities:        p.HardwareCapabilities,
		SlotCount:                   p.SlotCount,
		ChassisCount:                p.ChassisCount,
		Ports:                       p.Ports,
		Label:                       p.Label,
		ChassisInfo:                 p.ChassisInfo,
		LinecardsInfo:               p.LinecardsInfo,
		SlotConfiguration:           p.SlotConfiguration,
		PhysicalDevice:              p.PhysicalDevice,
	}

	return json.Marshal(&raw)
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

		ID             string     `json:"id"`
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
	p.id = raw.ID
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

// PortByID returns the Port with the given ID. If no port uses that ID, or if
// mulitple ports use that ID (unlikely), an error is returned.
func (p *Profile) PortByID(id int) (Port, error) {
	var result *Port

	for _, port := range p.Ports {
		if port.ID == id {
			if result != nil {
				return Port{}, sdk.ErrMultipleMatch(fmt.Sprintf("found multiple ports with ID %d", id))
			}
			result = &port
		}
	}

	if result == nil {
		return Port{}, sdk.ErrNotFound(fmt.Sprintf("found no ports with ID %d", id))
	}

	return *result, nil
}

// PortsByInterfaceName returns []Port containing all Ports which contain a
// Transformation which contains a TransformationInterface with the given name.
func (p Profile) PortsByInterfaceName(name string) []Port {
	var result []Port

portloop:
	for _, port := range p.Ports {
		for _, transformation := range port.Transformations {
			for _, intf := range transformation.Interfaces {
				if intf.Name == name {
					result = append(result, port)
					continue portloop
				}
			}
		}
	}

	return result
}

// PortByInterfaceName returns the Port which has at least one Transformation
// which contains a TransformationInterface with the given name. If zero ports
// or multiple ports use the desired name, an error is returned.
func (p Profile) PortByInterfaceName(name string) (Port, error) {
	ports := p.PortsByInterfaceName(name)

	switch len(ports) {
	case 0:
		return Port{}, sdk.ErrNotFound(fmt.Sprintf("found no ports with intinterface name %q", name))
	case 1:
		return ports[0], nil
	default:
		return Port{}, sdk.ErrMultipleMatch(fmt.Sprintf("found %d ports with intinterface name %q", len(ports), name))
	}
}

type HardwareCapabilities struct {
	MaxL3Mtu        *int            `json:"max_l3_mtu,omitempty"`
	MaxL2Mtu        *int            `json:"max_l2_mtu,omitempty"`
	FormFactor      string          `json:"form_factor"`
	VTEPLimit       *int            `json:"vtep_limit,omitempty"`
	BFDSupported    bool            `json:"bfd_supported,omitempty"`
	COPPStrict      FeatureVersions `json:"copp_strict,omitempty"`
	ECMPLimit       int             `json:"ecmp_limit"`
	ASNSequencing   FeatureVersions `json:"as_seq_num_supported,omitempty"`
	RAM             int             `json:"ram"`
	VTEPFloodLimit  *int            `json:"vtep_flood_limit,omitempty"`
	BreakoutCapable []struct {
		Module  int    `json:"module"`
		Value   bool   `json:"value"`
		Version string `json:"version"` // introduced and became mandatory some time between 4.2.0 and 6.0.0
	} `json:"breakout_capable,omitempty"`
	Userland        int             `json:"userland"`
	ASIC            string          `json:"asic"`
	VRFLimit        *int            `json:"vrf_limit,omitempty"`
	RoutingInstance FeatureVersions `json:"routing_instance_supported,omitempty"`
	VxlanSupported  bool            `json:"vxlan_supported,omitempty"`
	CPU             string          `json:"cpu"`
}

// PortWithMatchingTransforms searches its ports and transforms based on the
// supplied ifName (e.g. "xe-0/0/0:1") and speed. It returns the lone matching
// Port with that Port's Transformation slice filtered to contain only matching
// entries. For example, if we specified "ge-0/0/37" and "1G" with the
// Juniper_EX4400-48F device profile from Apstra 5.1.0, only one Port (id 38,
// panel 2, column 1, row 2) could match. That Port has 3 transforms, but only
// two match:
// - #2 with name "ge-0/0/37", speed "1G" and a setting which permits autonegotation
// - #3 with name "ge-0/0/37", speed "1G" and a setting which forbids autonegotation
// Transform #1 would be omitted from the returned Port both for its ifName
// ("xe-0/0/37") and its speed ("10G")
func (p Profile) PortWithMatchingTransforms(ifName string, ifSpeed speed.Speed) (Port, error) {
	port, err := p.PortByInterfaceName(ifName)
	if err != nil {
		return Port{}, fmt.Errorf("finding port known by %q in device profile %q: %w", ifName, p.id, err)
	}

	transformations := port.transformationCandidates(ifName, ifSpeed)
	if len(transformations) == 0 {
		return Port{}, sdk.ErrNotFound(fmt.Sprintf("port %d in device profile %q has no transformations named %s which operate at %s", port.ID, p.id, ifName, ifSpeed))
	}

	port.Transformations = transformations // replace the transform slice with the set of matching values

	return port, nil
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
