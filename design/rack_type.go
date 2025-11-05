// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hash"
	"sort"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ internal.IDer     = (*RackType)(nil)
	_ json.Marshaler    = (*RackType)(nil)
	_ json.Unmarshaler  = (*RackType)(nil)
	_ timeutils.Stamper = (*RackType)(nil)
)

type RackType struct {
	Label                    string
	Description              string
	FabricConnectivityDesign enum.FabricConnectivityDesign
	LeafSwitches             []RackTypeLeafSwitch
	AccessSwitches           []RackTypeAccessSwitch
	GenericSystems           []RackTypeGenericSystem

	id             string
	createdAt      *time.Time
	lastModifiedAt *time.Time
}

func (r RackType) ID() *string {
	if r.id == "" {
		return nil
	}
	return &r.id
}

// Replicate returns a copy of itself with zero values for metadata fields
func (r RackType) Replicate() RackType {
	result := RackType{
		Label:                    r.Label,
		Description:              r.Description,
		FabricConnectivityDesign: r.FabricConnectivityDesign,
		LeafSwitches:             make([]RackTypeLeafSwitch, len(r.LeafSwitches)),
		AccessSwitches:           nil, // don't create an empty slice
		GenericSystems:           nil, // don't create an empty slice
	}

	for i, leafSwitch := range r.LeafSwitches {
		result.LeafSwitches[i] = leafSwitch.Replicate()
	}

	if r.AccessSwitches != nil {
		result.AccessSwitches = make([]RackTypeAccessSwitch, len(r.AccessSwitches))
		for i, accessSwitch := range r.AccessSwitches {
			result.AccessSwitches[i] = accessSwitch.Replicate()
		}
	}

	if r.GenericSystems != nil {
		result.GenericSystems = make([]RackTypeGenericSystem, len(r.GenericSystems))
		for i, genericSystem := range r.GenericSystems {
			result.GenericSystems[i] = genericSystem.Replicate()
		}
	}

	return result
}

func (r RackType) MarshalJSON() ([]byte, error) {
	raw := struct {
		ID                       string                        `json:"id,omitempty"` // ID must be marshaled for template embedding
		Label                    string                        `json:"display_name"`
		Description              string                        `json:"description"`
		FabricConnectivityDesign enum.FabricConnectivityDesign `json:"fabric_connectivity_design"`
		Tags                     []Tag                         `json:"tags"`
		LogicalDevices           []LogicalDevice               `json:"logical_devices"`
		LeafSwitches             []RackTypeLeafSwitch          `json:"leafs"`
		AccessSwitches           []RackTypeAccessSwitch        `json:"access_switches,omitempty"`
		GenericSystems           []RackTypeGenericSystem       `json:"generic_systems,omitempty"`
	}{
		ID:                       r.id,
		Label:                    r.Label,
		Description:              r.Description,
		FabricConnectivityDesign: r.FabricConnectivityDesign,
		Tags:                     nil, // tags collected from various systems below
		LogicalDevices:           nil, // logical devices collected from various systems below
		LeafSwitches:             r.LeafSwitches,
		AccessSwitches:           r.AccessSwitches,
		GenericSystems:           r.GenericSystems,
	}

	tagMap := make(map[string]Tag)
	addTagsToMap := func(tags []Tag) {
		for _, tag := range tags {
			tagMap[tag.Label] = tag.Replicate()
		}
	}

	hasher := md5.New()
	logicalDeviceMap := make(map[string]LogicalDevice)

	// populate the rack-level tags and logical devices maps for each system
	// (all leaf switches, access switches and generic systems) in the rack.
	for _, system := range r.LeafSwitches {
		// collect the tags in the map which will populate the rack-level list
		addTagsToMap(system.Tags)

		// clone the system and set its LD ID to a hash of the LD payload
		system := system.Replicate()
		system.LogicalDevice.setHashID(hasher)

		// add the LD to the rack-wide LD map
		logicalDeviceMap[*system.logicalDeviceID()] = system.LogicalDevice
	}
	for _, system := range r.AccessSwitches {
		// collect system and link tags in the map which will populate the rack-level list
		var tags []Tag
		tags = append(tags, system.Tags...)
		for _, link := range system.Links {
			tags = append(tags, link.Tags...)
		}
		addTagsToMap(tags)

		// clone the system and set its LD ID to a hash of the LD payload
		system := system.Replicate()
		system.LogicalDevice.setHashID(hasher)

		// add the LD to the rack-wide LD map
		logicalDeviceMap[*system.logicalDeviceID()] = system.LogicalDevice
	}
	for _, system := range r.GenericSystems {
		// collect system and link tags in the map which will populate the rack-level list
		var tags []Tag
		tags = append(tags, system.Tags...)
		for _, link := range system.Links {
			tags = append(tags, link.Tags...)
		}
		addTagsToMap(tags)

		// clone the system and set its LD ID to a hash of the LD payload
		system := system.Replicate()
		system.LogicalDevice.setHashID(hasher)

		// add the LD to the rack-wide LD map
		logicalDeviceMap[*system.logicalDeviceID()] = system.LogicalDevice
	}

	// having de-duped tags via map, convert them to sorted slice
	raw.Tags = make([]Tag, 0, len(tagMap))
	for _, tag := range tagMap {
		raw.Tags = append(raw.Tags, tag)
	}
	sort.Slice(raw.Tags, func(i, j int) bool {
		return raw.Tags[i].Label < raw.Tags[j].Label
	})

	// having de-duped logical devices via map, convert them to sorted slice
	raw.LogicalDevices = make([]LogicalDevice, 0, len(logicalDeviceMap))
	for _, logicalDevice := range logicalDeviceMap {
		raw.LogicalDevices = append(raw.LogicalDevices, logicalDevice)
	}
	sort.Slice(raw.LogicalDevices, func(i, j int) bool {
		return raw.LogicalDevices[i].id < raw.LogicalDevices[j].id
	})

	return json.Marshal(&raw)
}

func (r *RackType) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID                       string                        `json:"id,omitempty"` // ID must be marshaled for template embedding
		Label                    string                        `json:"display_name"`
		Description              string                        `json:"description"`
		FabricConnectivityDesign enum.FabricConnectivityDesign `json:"fabric_connectivity_design"`
		AllTags                  []Tag                         `json:"tags"`
		LogicalDevices           []LogicalDevice               `json:"logical_devices"`
		LeafSwitches             []RackTypeLeafSwitch          `json:"leafs"`
		AccessSwitches           []RackTypeAccessSwitch        `json:"access_switches"`
		GenericSystems           []RackTypeGenericSystem       `json:"generic_systems"`
		CreatedAt                *time.Time                    `json:"created_at"`
		LastModifiedAt           *time.Time                    `json:"last_modified_at"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return fmt.Errorf("unmarshaling rack_type: %w", err)
	}

	logicalDeviceMap := make(map[string]LogicalDevice, len(raw.LogicalDevices))
	for _, ld := range raw.LogicalDevices {
		logicalDeviceMap[ld.id] = ld.Replicate() // Replicate drops metadata
	}

	r.id = raw.ID
	r.Label = raw.Label
	r.Description = raw.Description
	r.FabricConnectivityDesign = raw.FabricConnectivityDesign

	r.LeafSwitches = make([]RackTypeLeafSwitch, len(raw.LeafSwitches))
	for i, system := range raw.LeafSwitches {
		// find logical device with full detail in rack-level map
		logicalDevice, ok := logicalDeviceMap[system.LogicalDevice.id]
		if !ok {
			return fmt.Errorf("leaf switch %d logical device (%q) not found", i, system.LogicalDevice.id)
		}

		r.LeafSwitches[i] = system.Replicate() // Replicate drops metadata
		r.LeafSwitches[i].LogicalDevice = logicalDevice.Replicate()
		r.LeafSwitches[i].Tags = system.Tags // half-baked tags filled in next
		err = populateTagsByLabel(raw.AllTags, r.LeafSwitches[i].Tags)
		if err != nil {
			return fmt.Errorf("populating tags for leaf switch %d: %w", i, err)
		}
	}

	if len(raw.AccessSwitches) > 0 {
		r.AccessSwitches = make([]RackTypeAccessSwitch, len(raw.AccessSwitches))
	}
	for i, system := range raw.AccessSwitches {
		// find logical device with full detail in rack-level map
		logicalDevice, ok := logicalDeviceMap[system.LogicalDevice.id]
		if !ok {
			return fmt.Errorf("access switch %d logical device (%q) not found", i, system.LogicalDevice.id)
		}

		r.AccessSwitches[i] = system.Replicate() // Replicate drops metadata
		r.AccessSwitches[i].LogicalDevice = logicalDevice.Replicate()

		r.AccessSwitches[i].Tags = system.Tags // half-baked tags filled in next
		err = populateTagsByLabel(raw.AllTags, r.AccessSwitches[i].Tags)
		if err != nil {
			return fmt.Errorf("populating tags for access switch %d: %w", i, err)
		}
		for j := range r.AccessSwitches[i].Links {
			err = populateTagsByLabel(raw.AllTags, r.AccessSwitches[i].Links[j].Tags)
			if err != nil {
				return fmt.Errorf("populating tags for access switch %d link %d: %w", i, j, err)
			}
		}

	}

	if len(raw.GenericSystems) > 0 {
		r.GenericSystems = make([]RackTypeGenericSystem, len(raw.GenericSystems))
	}
	for i, system := range raw.GenericSystems {
		// find logical device with full detail in rack-level map
		logicalDevice, ok := logicalDeviceMap[system.LogicalDevice.id]
		if !ok {
			return fmt.Errorf("generic system %d logical device (%q) not found", i, system.LogicalDevice.id)
		}

		r.GenericSystems[i] = system.Replicate() // Replicate drops metadata
		r.GenericSystems[i].LogicalDevice = logicalDevice.Replicate()

		r.GenericSystems[i].Tags = system.Tags // half-baked tags filled in next
		err = populateTagsByLabel(raw.AllTags, r.GenericSystems[i].Tags)
		if err != nil {
			return fmt.Errorf("populating tags for generic system %d: %w", i, err)
		}
		for j := range r.GenericSystems[i].Links {
			err = populateTagsByLabel(raw.AllTags, r.GenericSystems[i].Links[j].Tags)
			if err != nil {
				return fmt.Errorf("populating tags for generic system %d link %d: %w", i, j, err)
			}
		}

	}

	r.createdAt = raw.CreatedAt
	r.lastModifiedAt = raw.LastModifiedAt

	return nil
}

func (r RackType) CreatedAt() *time.Time {
	return r.createdAt
}

func (r RackType) LastModifiedAt() *time.Time {
	return r.lastModifiedAt
}

func (r RackType) digest(h hash.Hash) []byte {
	h.Reset()
	return mustHashForComparison(r, h)
}

func (r *RackType) setHashID(h hash.Hash) {
	if r.id != "" {
		panic(fmt.Sprintf("id already has value %q", r.id))
	}

	r.id = fmt.Sprintf("%x", r.digest(h))
	return
}

func NewRackType(id string) RackType {
	return RackType{id: id}
}
