// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var (
	_ ider                 = (*RackType)(nil)
	_ replicator[RackType] = (*RackType)(nil)
	_ json.Marshaler       = (*RackType)(nil)
	_ json.Unmarshaler     = (*RackType)(nil)
	_ timeutils.Stamper    = (*RackType)(nil)
)

type RackType struct {
	Label                    string
	Description              string
	FabricConnectivityDesign enum.FabricConnectivityDesign
	Status                   *enum.FFEConsistencyStatus
	LeafSwitches             []LeafSwitch
	// AccessSwitches        []AccessSwitch  `json:"access_switches"`
	// GenericSystems        []GenericSystem `json:"generic_systems"`

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

// SetID sets a previously un-set id attribute. If the id attribute is found to
// have an existing value, an error is returned. Presence of an existing value
// is the only reason SetID will return an error. If the id attribute is known
// to be empty, use MustSetID.
func (r *RackType) SetID(id string) error {
	if r.id != "" {
		return IDIsSet(fmt.Errorf("tag id alredy has value %q", r.id))
	}

	r.id = id
	return nil
}

// MustSetID invokes SetID and panics if an error is returned.
func (r *RackType) MustSetID(id string) {
	err := r.SetID(id)
	if err != nil {
		panic(err)
	}
}

func (r RackType) replicate() RackType {
	return RackType{
		Label:                    r.Label,
		Description:              r.Description,
		FabricConnectivityDesign: r.FabricConnectivityDesign,
		Status:                   r.Status,
		LeafSwitches:             nil,
	}
}

func (r RackType) MarshalJSON() ([]byte, error) {
	raw := struct {
		ID                       string                        `json:"id,omitempty"` // ID must be marshaled when embedded in template
		Label                    string                        `json:"display_name"`
		Description              string                        `json:"description"`
		FabricConnectivityDesign enum.FabricConnectivityDesign `json:"fabric_connectivity_design"`
		Tags                     []Tag                         `json:"tags"`
		LogicalDevices           []LogicalDevice               `json:"logical_devices"`
		LeafSwitches             []LeafSwitch                  `json:"leaf_switches"`
		// AccessSwitches        []AccessSwitch                `json:"access_switches"`
		// GenericSystems        []GenericSystem               `json:"generic_systems"`
	}{
		ID:                       r.id,
		Label:                    r.Label,
		Description:              r.Description,
		FabricConnectivityDesign: r.FabricConnectivityDesign,
		Tags:                     nil, // tags collected from various systems below
		LogicalDevices:           nil, // logical devices collected from various systems below
		LeafSwitches:             r.LeafSwitches,
		// AccessSwitches:        r.AccessSwitches,
		// GenericSystems:        r.GenericSystems
	}

	tagMap := make(map[string]Tag)
	addTagsToMap := func(tags []Tag) {
		for _, tag := range tags {
			tagMap[tag.Label] = tag.replicate()
		}
	}

	idHash := md5.New()
	logicalDeviceMap := make(map[string]LogicalDevice)
	addLogicalDeviceToMap := func(ld LogicalDevice) {
		id := fmt.Sprintf("%x", mustDigestSkipID(ld, idHash))
		idHash.Reset()

		ld = ld.replicate()
		ld.MustSetID(id)
		logicalDeviceMap[id] = ld
	}

	// populate the tags and logicalDevices maps
	for _, system := range r.LeafSwitches {
		addTagsToMap(system.Tags)
		addLogicalDeviceToMap(system.LogicalDevice)
	}
	//for _, system := range r.AccessSwitches {
	//	addTagsToMap(system.Tags)
	//	addLogicalDeviceToMap(system.LogicalDevice)
	//}
	//for _, system := range r.GenericSystems {
	//	addTagsToMap(system.Tags)
	//	addLogicalDeviceToMap(system.LogicalDevice)
	//}

	// having de-duped tags and logical devices via map, add the values to the raw slices
	raw.Tags = make([]Tag, 0, len(tagMap))
	for _, tag := range tagMap {
		raw.Tags = append(raw.Tags, tag)
	}
	raw.LogicalDevices = make([]LogicalDevice, 0, len(logicalDeviceMap))
	for _, logicalDevice := range logicalDeviceMap {
		raw.LogicalDevices = append(raw.LogicalDevices, logicalDevice)
	}

	return json.Marshal(&raw)
}

func (r *RackType) UnmarshalJSON(bytes []byte) error {
	// TODO implement me
	panic("implement me")
}

func (r RackType) CreatedAt() *time.Time {
	return r.createdAt
}

func (r RackType) LastModifiedAt() *time.Time {
	return r.lastModifiedAt
}

func NewRackType(id string) RackType {
	return RackType{id: id}
}
