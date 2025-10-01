// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"time"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

var testProfileGeneric1x10 = Profile{
	id:             "Generic_Server_1RU_1x10G",
	createdAt:      pointer.To(time.Date(2024, time.January, 2, 3, 4, 5, 6000, time.UTC)),
	lastModifiedAt: pointer.To(time.Date(2025, time.January, 2, 3, 4, 5, 6000, time.UTC)),
	Label:          "Generic_Server_1RU_1x10G",
	Selector: Selector{
		OsVersion:    ".*",
		Model:        "Generic Model",
		Os:           "Ubuntu GNU/Linux",
		Manufacturer: "Generic Manufacturer",
	},
	DeviceProfileType: enum.DeviceProfileTypeMonolithic,
	DualRoutingEngine: false,
	SoftwareCapabilities: SoftwareCapabilities{
		Onie:               false,
		ConfigApplySupport: "complete_only",
		LxcSupport:         false,
	},
	ReferenceDesignCapabilities: &ReferenceDesignCapabilities{
		Datacenter: enum.RefDesignCapabilityFullSupport,
		Freeform:   enum.RefDesignCapabilityFullSupport,
	},
	ChassisProfileID: "",
	HardwareCapabilities: HardwareCapabilities{
		ASIC:       "",
		ECMPLimit:  64,
		RAM:        16,
		CPU:        "x86",
		FormFactor: "1RU",
		Userland:   64,
	},
	Predefined: true,
	SlotCount:  pointer.To(0),
	Ports: []Port{
		{
			ConnectorType: "sfp",
			Panel:         1,
			Transformations: []Transformation{
				{
					ID:        1,
					IsDefault: true,
					Interfaces: []TransformationInterface{
						{
							ID:      1,
							Name:    "eth0",
							State:   enum.InterfaceStateActive,
							Setting: pointer.To(""),
							Speed:   "10G",
						},
					},
				},
			},
			Column:        1,
			ID:            1,
			Row:           1,
			FailureDomain: 1,
			Display:       pointer.To(1),
			Slot:          0,
		},
	},
	PhysicalDevice: true,
}

const testProfileGeneric1x10JSON = `{
  "id": "Generic_Server_1RU_1x10G",
  "created_at": "2024-01-02T03:04:05.000006Z",
  "last_modified_at": "2025-01-02T03:04:05.000006Z",
  "label": "Generic_Server_1RU_1x10G",
  "selector": {
    "os": "Ubuntu GNU/Linux",
    "os_version": ".*",
    "manufacturer": "Generic Manufacturer",
    "model": "Generic Model"
  },
  "device_profile_type": "monolithic",
  "dual_routing_engine": false,
  "software_capabilities": {
    "onie": false,
    "lxc_support": false,
    "config_apply_support": "complete_only"
  },
  "reference_design_capabilities": {
    "datacenter": "full_support",
    "freeform": "full_support"
  },
  "hardware_capabilities": {
    "asic": "",
    "ecmp_limit": 64,
    "ram": 16,
    "cpu": "x86",
    "form_factor": "1RU",
    "userland": 64
  },
  "predefined": true,
  "slot_count": 0,
  "ports": [
    {
      "connector_type": "sfp",
      "panel_id": 1,
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "eth0",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": ""
            }
          ]
        }
      ],
      "column_id": 1,
      "port_id": 1,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 1,
      "slot_id": 0
    }
  ],
  "physical_device": true
}`

var testProfileJunipervEX = Profile{
	id:                          "Juniper_vEX",
	createdAt:                   pointer.To(time.Date(2024, time.December, 12, 0, 40, 16, 337983000, time.UTC)),
	lastModifiedAt:              pointer.To(time.Date(2024, time.December, 12, 0, 40, 16, 337983000, time.UTC)),
	Label:                       "Juniper vEX",
	Predefined:                  true,
	Selector:                    Selector{OsVersion: ".*", Model: "VIRTUAL-EX9214", Os: "Junos", Manufacturer: "Juniper"},
	DeviceProfileType:           enum.DeviceProfileTypeMonolithic,
	DualRoutingEngine:           false,
	SoftwareCapabilities:        SoftwareCapabilities{Onie: false, ConfigApplySupport: "complete_only", LxcSupport: false},
	ReferenceDesignCapabilities: &ReferenceDesignCapabilities{Datacenter: enum.RefDesignCapabilityFullSupport, Freeform: enum.RefDesignCapabilityFullSupport},
	ChassisProfileID:            "",
	HardwareCapabilities: HardwareCapabilities{
		FormFactor:      "1RU",
		BFDSupported:    false,
		ECMPLimit:       32,
		RAM:             32,
		BreakoutCapable: nil,
		Userland:        32,
		ASIC:            "Trio",
		RoutingInstance: FeatureVersions{FeatureVersion{Version: ".*", Enabled: true}},
		VxlanSupported:  false,
		CPU:             "x86",
	},
	SlotCount: pointer.To(0),
	Ports: []Port{
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(0), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/1", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/1", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 2, ID: 2, Row: 1, FailureDomain: 1, Display: pointer.To(1), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/2", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/2", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 3, ID: 3, Row: 1, FailureDomain: 1, Display: pointer.To(2), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/3", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/3", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 4, ID: 4, Row: 1, FailureDomain: 1, Display: pointer.To(3), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/4", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/4", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 5, ID: 5, Row: 1, FailureDomain: 1, Display: pointer.To(4), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 6, ID: 6, Row: 1, FailureDomain: 1, Display: pointer.To(5), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/6", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/6", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 7, ID: 7, Row: 1, FailureDomain: 1, Display: pointer.To(6), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/7", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/7", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 8, ID: 8, Row: 1, FailureDomain: 1, Display: pointer.To(7), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/8", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/8", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 9, ID: 9, Row: 1, FailureDomain: 1, Display: pointer.To(8), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/9", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/9", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 10, ID: 10, Row: 1, FailureDomain: 1, Display: pointer.To(9), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/0", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 11, ID: 11, Row: 1, FailureDomain: 1, Display: pointer.To(10), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/1", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/1", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 12, ID: 12, Row: 1, FailureDomain: 1, Display: pointer.To(11), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/2", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/2", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 13, ID: 13, Row: 1, FailureDomain: 1, Display: pointer.To(12), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/3", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/3", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 14, ID: 14, Row: 1, FailureDomain: 1, Display: pointer.To(13), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/4", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/4", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 15, ID: 15, Row: 1, FailureDomain: 1, Display: pointer.To(14), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/5", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 16, ID: 16, Row: 1, FailureDomain: 1, Display: pointer.To(15), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/6", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/6", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 17, ID: 17, Row: 1, FailureDomain: 1, Display: pointer.To(16), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/7", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/7", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 18, ID: 18, Row: 1, FailureDomain: 1, Display: pointer.To(17), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/8", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/8", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 19, ID: 19, Row: 1, FailureDomain: 1, Display: pointer.To(18), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/9", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "10g"}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/1/9", State: enum.InterfaceStateActive, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}}, Column: 20, ID: 20, Row: 1, FailureDomain: 1, Display: pointer.To(19), Slot: 0},
	},
}

const testProfileJunipervEXJSON = `{
  "id": "Juniper_vEX",
  "created_at": "2024-12-12T00:40:16.337983Z",
  "last_modified_at": "2024-12-12T00:40:16.337983Z",
  "label": "Juniper vEX",
  "predefined": true,
  "dual_routing_engine": false,
  "physical_device": false,
  "slot_count": 0,
  "ports": [
    {
      "port_id": 1,
      "display_id": 0,
      "row_id": 1,
      "column_id": 1,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/0",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/0",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 2,
      "display_id": 1,
      "row_id": 1,
      "column_id": 2,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/1",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/1",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 3,
      "display_id": 2,
      "row_id": 1,
      "column_id": 3,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/2",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/2",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 4,
      "display_id": 3,
      "row_id": 1,
      "column_id": 4,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/3",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/3",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 5,
      "display_id": 4,
      "row_id": 1,
      "column_id": 5,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/4",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/4",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 6,
      "display_id": 5,
      "row_id": 1,
      "column_id": 6,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/5",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/5",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 7,
      "display_id": 6,
      "row_id": 1,
      "column_id": 7,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/6",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/6",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 8,
      "display_id": 7,
      "row_id": 1,
      "column_id": 8,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/7",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/7",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 9,
      "display_id": 8,
      "row_id": 1,
      "column_id": 9,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/8",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/8",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 10,
      "display_id": 9,
      "row_id": 1,
      "column_id": 10,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/9",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/9",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 11,
      "display_id": 10,
      "row_id": 1,
      "column_id": 11,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/0",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/0",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 12,
      "display_id": 11,
      "row_id": 1,
      "column_id": 12,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/1",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/1",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 13,
      "display_id": 12,
      "row_id": 1,
      "column_id": 13,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/2",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/2",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 14,
      "display_id": 13,
      "row_id": 1,
      "column_id": 14,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/3",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/3",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 15,
      "display_id": 14,
      "row_id": 1,
      "column_id": 15,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/4",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/4",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 16,
      "display_id": 15,
      "row_id": 1,
      "column_id": 16,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/5",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/5",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 17,
      "display_id": 16,
      "row_id": 1,
      "column_id": 17,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/6",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/6",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 18,
      "display_id": 17,
      "row_id": 1,
      "column_id": 18,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/7",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/7",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 19,
      "display_id": 18,
      "row_id": 1,
      "column_id": 19,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/8",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/8",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 20,
      "display_id": 19,
      "row_id": 1,
      "column_id": 20,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/1/9",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"10g\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/1/9",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    }
  ],
  "hardware_capabilities": {
    "asic": "Trio",
    "ecmp_limit": 32,
    "ram": 32,
    "routing_instance_supported": [
      {
        "value": true,
        "version": ".*"
      }
    ],
    "cpu": "x86",
    "form_factor": "1RU",
    "userland": 32
  },
  "software_capabilities": {
    "onie": false,
    "lxc_support": false,
    "config_apply_support": "complete_only"
  },
  "selector": {
    "os": "Junos",
    "os_version": ".*",
    "manufacturer": "Juniper",
    "model": "VIRTUAL-EX9214"
  },
  "reference_design_capabilities": {
    "datacenter": "full_support",
    "freeform": "full_support"
  },
  "device_profile_type": "monolithic"
}`

var testProfileJuniperEX440048F = Profile{
	id:                          "Juniper_EX4400-48F",
	createdAt:                   pointer.To(time.Date(2024, time.December, 12, 0, 40, 19, 913413000, time.UTC)),
	lastModifiedAt:              pointer.To(time.Date(2024, time.December, 12, 0, 40, 19, 913413000, time.UTC)),
	Label:                       "Juniper_EX4400-48F",
	Predefined:                  true,
	Selector:                    Selector{OsVersion: ".*", Model: "EX4400-48F", Os: "Junos", Manufacturer: "Juniper"},
	DeviceProfileType:           enum.DeviceProfileType{Value: "monolithic"},
	SoftwareCapabilities:        SoftwareCapabilities{Onie: false, ConfigApplySupport: "complete_only", LxcSupport: false},
	ReferenceDesignCapabilities: &ReferenceDesignCapabilities{Datacenter: enum.RefDesignCapabilityFullSupport, Freeform: enum.RefDesignCapabilityFullSupport},
	ChassisProfileID:            "",
	HardwareCapabilities: HardwareCapabilities{
		FormFactor:      "1RU",
		BFDSupported:    false,
		ECMPLimit:       64,
		RAM:             4,
		Userland:        64,
		ASIC:            "T3",
		RoutingInstance: FeatureVersions{FeatureVersion{Version: ".*", Enabled: true}},
		VxlanSupported:  false,
		CPU:             "x86",
	},
	SlotCount: pointer.To(0),
	Ports: []Port{
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 0, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 1, ID: 1, Row: 1, FailureDomain: 1, Display: pointer.To(0), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 1, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 1, ID: 2, Row: 2, FailureDomain: 1, Display: pointer.To(1), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 2, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 2, ID: 3, Row: 1, FailureDomain: 1, Display: pointer.To(2), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 3, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 2, ID: 4, Row: 2, FailureDomain: 1, Display: pointer.To(3), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/4", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/4", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 4, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 3, ID: 5, Row: 1, FailureDomain: 1, Display: pointer.To(4), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/5", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 5, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 3, ID: 6, Row: 2, FailureDomain: 1, Display: pointer.To(5), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/6", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/6", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 6, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 4, ID: 7, Row: 1, FailureDomain: 1, Display: pointer.To(6), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/7", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/7", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 7, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 4, ID: 8, Row: 2, FailureDomain: 1, Display: pointer.To(7), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/8", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/8", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 8, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 5, ID: 9, Row: 1, FailureDomain: 1, Display: pointer.To(8), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/9", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/9", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 9, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 5, ID: 10, Row: 2, FailureDomain: 1, Display: pointer.To(9), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/10", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/10", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 10, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 6, ID: 11, Row: 1, FailureDomain: 1, Display: pointer.To(10), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/11", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/11", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 11, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 6, ID: 12, Row: 2, FailureDomain: 1, Display: pointer.To(11), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/12", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 12, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 7, ID: 13, Row: 1, FailureDomain: 1, Display: pointer.To(12), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/13", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/13", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 13, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 7, ID: 14, Row: 2, FailureDomain: 1, Display: pointer.To(13), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/14", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/14", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 14, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 8, ID: 15, Row: 1, FailureDomain: 1, Display: pointer.To(14), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/15", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/15", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 15, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 8, ID: 16, Row: 2, FailureDomain: 1, Display: pointer.To(15), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/16", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/16", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 16, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 9, ID: 17, Row: 1, FailureDomain: 1, Display: pointer.To(16), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/17", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/17", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 17, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 9, ID: 18, Row: 2, FailureDomain: 1, Display: pointer.To(17), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/18", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/18", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 18, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 10, ID: 19, Row: 1, FailureDomain: 1, Display: pointer.To(18), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/19", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/19", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 19, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 10, ID: 20, Row: 2, FailureDomain: 1, Display: pointer.To(19), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/20", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/20", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 20, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 11, ID: 21, Row: 1, FailureDomain: 1, Display: pointer.To(20), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/21", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/21", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 21, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 11, ID: 22, Row: 2, FailureDomain: 1, Display: pointer.To(21), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/22", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/22", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 22, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 12, ID: 23, Row: 1, FailureDomain: 1, Display: pointer.To(22), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/23", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/23", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 23, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 12, ID: 24, Row: 2, FailureDomain: 1, Display: pointer.To(23), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/24", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/24", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 24, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 13, ID: 25, Row: 1, FailureDomain: 1, Display: pointer.To(24), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/25", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/25", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 25, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 13, ID: 26, Row: 2, FailureDomain: 1, Display: pointer.To(25), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/26", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/26", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 26, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 14, ID: 27, Row: 1, FailureDomain: 1, Display: pointer.To(26), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/27", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/27", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 27, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 14, ID: 28, Row: 2, FailureDomain: 1, Display: pointer.To(27), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/28", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/28", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 28, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 15, ID: 29, Row: 1, FailureDomain: 1, Display: pointer.To(28), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/29", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/29", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 29, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 15, ID: 30, Row: 2, FailureDomain: 1, Display: pointer.To(29), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/30", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/30", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 30, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 16, ID: 31, Row: 1, FailureDomain: 1, Display: pointer.To(30), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/31", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/31", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 31, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 16, ID: 32, Row: 2, FailureDomain: 1, Display: pointer.To(31), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/32", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/32", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 32, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 17, ID: 33, Row: 1, FailureDomain: 1, Display: pointer.To(32), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/33", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/33", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 33, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 17, ID: 34, Row: 2, FailureDomain: 1, Display: pointer.To(33), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/34", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/34", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 34, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 18, ID: 35, Row: 1, FailureDomain: 1, Display: pointer.To(34), Slot: 0},
		{ConnectorType: "sfp", Panel: 1, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/35", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "1G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/35", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": 35, "speed": "100m"}, "interface": {"speed": ""}}`), Speed: "100M"}}}}, Column: 18, ID: 36, Row: 2, FailureDomain: 1, Display: pointer.To(35), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/36", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/36", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/36", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 1, ID: 37, Row: 1, FailureDomain: 1, Display: pointer.To(36), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/37", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/37", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/37", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 1, ID: 38, Row: 2, FailureDomain: 1, Display: pointer.To(37), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/38", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/38", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/38", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 2, ID: 39, Row: 1, FailureDomain: 1, Display: pointer.To(38), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/39", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/39", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/39", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 2, ID: 40, Row: 2, FailureDomain: 1, Display: pointer.To(39), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/40", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/40", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/40", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 3, ID: 41, Row: 1, FailureDomain: 1, Display: pointer.To(40), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/41", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/41", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/41", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 3, ID: 42, Row: 2, FailureDomain: 1, Display: pointer.To(41), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/42", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/42", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/42", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 4, ID: 43, Row: 1, FailureDomain: 1, Display: pointer.To(42), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/43", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/43", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/43", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 4, ID: 44, Row: 2, FailureDomain: 1, Display: pointer.To(43), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/44", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/44", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/44", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 5, ID: 45, Row: 1, FailureDomain: 1, Display: pointer.To(44), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/45", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/45", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/45", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 5, ID: 46, Row: 2, FailureDomain: 1, Display: pointer.To(45), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/46", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/46", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/46", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 6, ID: 47, Row: 1, FailureDomain: 1, Display: pointer.To(46), Slot: 0},
		{ConnectorType: "sfp+", Panel: 2, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/0/47", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/47", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/0/47", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 6, ID: 48, Row: 2, FailureDomain: 1, Display: pointer.To(47), Slot: 0},
		{ConnectorType: "qsfp28", Panel: 3, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "et-0/1/0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "100G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "et-0/1/0:0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 2, Name: "et-0/1/0:1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 3, Name: "et-0/1/0:2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 4, Name: "et-0/1/0:3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/1/0:0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 2, Name: "xe-0/1/0:1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 3, Name: "xe-0/1/0:2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 4, Name: "xe-0/1/0:3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 0, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}}}}, Column: 1, ID: 49, Row: 1, FailureDomain: 1, Display: pointer.To(48), Slot: 0},
		{ConnectorType: "qsfp28", Panel: 3, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "et-0/1/1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "100G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "et-0/1/1:0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 2, Name: "et-0/1/1:1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 3, Name: "et-0/1/1:2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}, {ID: 4, Name: "et-0/1/1:3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "25g"}, "interface": {"speed": ""}}`), Speed: "25G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/1/1:0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 2, Name: "xe-0/1/1:1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 3, Name: "xe-0/1/1:2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}, {ID: 4, Name: "xe-0/1/1:3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"breakout": true, "fpc": 0, "pic": 1, "port": 1, "speed": "10g"}, "interface": {"speed": ""}}`), Speed: "10G"}}}}, Column: 1, ID: 50, Row: 2, FailureDomain: 1, Display: pointer.To(49), Slot: 0},
		{ConnectorType: "sfp+", Panel: 4, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/2/0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/0", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 1, ID: 51, Row: 1, FailureDomain: 1, Display: pointer.To(50), Slot: 0},
		{ConnectorType: "sfp+", Panel: 4, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/2/1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/1", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 2, ID: 52, Row: 1, FailureDomain: 1, Display: pointer.To(51), Slot: 0},
		{ConnectorType: "sfp+", Panel: 4, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/2/2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/2", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 3, ID: 53, Row: 1, FailureDomain: 1, Display: pointer.To(52), Slot: 0},
		{ConnectorType: "sfp+", Panel: 4, Transformations: []Transformation{{ID: 1, IsDefault: true, Interfaces: []TransformationInterface{{ID: 1, Name: "xe-0/2/3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": ""}}`), Speed: "10G"}}}, {ID: 2, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"speed": "1g"}}`), Speed: "1G"}}}, {ID: 3, IsDefault: false, Interfaces: []TransformationInterface{{ID: 1, Name: "ge-0/2/3", State: enum.InterfaceState{Value: "active"}, Setting: pointer.To(`{"global": {"speed": ""}, "interface": {"auto_negotiation": false, "link_mode": "full-duplex", "speed": "1g"}}`), Speed: "1G"}}}}, Column: 4, ID: 54, Row: 1, FailureDomain: 1, Display: pointer.To(53), Slot: 0},
	},
	PhysicalDevice: true,
}

const testProfileJuniperEX440048FJSON = `{
  "id": "Juniper_EX4400-48F",
  "created_at": "2024-12-12T00:40:19.913413Z",
  "last_modified_at": "2024-12-12T00:40:19.913413Z",
  "label": "Juniper_EX4400-48F",
  "predefined": true,
  "dual_routing_engine": false,
  "physical_device": true,
  "slot_count": 0,
  "ports": [
    {
      "port_id": 1,
      "display_id": 0,
      "row_id": 1,
      "column_id": 1,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/0",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/0",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 0, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 2,
      "display_id": 1,
      "row_id": 2,
      "column_id": 1,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/1",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/1",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 1, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 3,
      "display_id": 2,
      "row_id": 1,
      "column_id": 2,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/2",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/2",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 2, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 4,
      "display_id": 3,
      "row_id": 2,
      "column_id": 2,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/3",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/3",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 3, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 5,
      "display_id": 4,
      "row_id": 1,
      "column_id": 3,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/4",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/4",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 4, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 6,
      "display_id": 5,
      "row_id": 2,
      "column_id": 3,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/5",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/5",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 5, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 7,
      "display_id": 6,
      "row_id": 1,
      "column_id": 4,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/6",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/6",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 6, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 8,
      "display_id": 7,
      "row_id": 2,
      "column_id": 4,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/7",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/7",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 7, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 9,
      "display_id": 8,
      "row_id": 1,
      "column_id": 5,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/8",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/8",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 8, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 10,
      "display_id": 9,
      "row_id": 2,
      "column_id": 5,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/9",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/9",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 9, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 11,
      "display_id": 10,
      "row_id": 1,
      "column_id": 6,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/10",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/10",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 10, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 12,
      "display_id": 11,
      "row_id": 2,
      "column_id": 6,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/11",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/11",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 11, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 13,
      "display_id": 12,
      "row_id": 1,
      "column_id": 7,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/12",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/12",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 12, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 14,
      "display_id": 13,
      "row_id": 2,
      "column_id": 7,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/13",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/13",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 13, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 15,
      "display_id": 14,
      "row_id": 1,
      "column_id": 8,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/14",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/14",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 14, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 16,
      "display_id": 15,
      "row_id": 2,
      "column_id": 8,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/15",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/15",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 15, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 17,
      "display_id": 16,
      "row_id": 1,
      "column_id": 9,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/16",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/16",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 16, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 18,
      "display_id": 17,
      "row_id": 2,
      "column_id": 9,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/17",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/17",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 17, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 19,
      "display_id": 18,
      "row_id": 1,
      "column_id": 10,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/18",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/18",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 18, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 20,
      "display_id": 19,
      "row_id": 2,
      "column_id": 10,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/19",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/19",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 19, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 21,
      "display_id": 20,
      "row_id": 1,
      "column_id": 11,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/20",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/20",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 20, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 22,
      "display_id": 21,
      "row_id": 2,
      "column_id": 11,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/21",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/21",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 21, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 23,
      "display_id": 22,
      "row_id": 1,
      "column_id": 12,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/22",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/22",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 22, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 24,
      "display_id": 23,
      "row_id": 2,
      "column_id": 12,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/23",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/23",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 23, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 25,
      "display_id": 24,
      "row_id": 1,
      "column_id": 13,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/24",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/24",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 24, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 26,
      "display_id": 25,
      "row_id": 2,
      "column_id": 13,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/25",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/25",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 25, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 27,
      "display_id": 26,
      "row_id": 1,
      "column_id": 14,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/26",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/26",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 26, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 28,
      "display_id": 27,
      "row_id": 2,
      "column_id": 14,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/27",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/27",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 27, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 29,
      "display_id": 28,
      "row_id": 1,
      "column_id": 15,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/28",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/28",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 28, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 30,
      "display_id": 29,
      "row_id": 2,
      "column_id": 15,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/29",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/29",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 29, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 31,
      "display_id": 30,
      "row_id": 1,
      "column_id": 16,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/30",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/30",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 30, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 32,
      "display_id": 31,
      "row_id": 2,
      "column_id": 16,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/31",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/31",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 31, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 33,
      "display_id": 32,
      "row_id": 1,
      "column_id": 17,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/32",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/32",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 32, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 34,
      "display_id": 33,
      "row_id": 2,
      "column_id": 17,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/33",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/33",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 33, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 35,
      "display_id": 34,
      "row_id": 1,
      "column_id": 18,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/34",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/34",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 34, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 36,
      "display_id": 35,
      "row_id": 2,
      "column_id": 18,
      "panel_id": 1,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "ge-0/0/35",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/35",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "M"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 35, \"speed\": \"100m\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 37,
      "display_id": 36,
      "row_id": 1,
      "column_id": 1,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/36",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/36",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/36",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 38,
      "display_id": 37,
      "row_id": 2,
      "column_id": 1,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/37",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/37",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/37",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 39,
      "display_id": 38,
      "row_id": 1,
      "column_id": 2,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/38",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/38",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/38",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 40,
      "display_id": 39,
      "row_id": 2,
      "column_id": 2,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/39",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/39",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/39",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 41,
      "display_id": 40,
      "row_id": 1,
      "column_id": 3,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/40",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/40",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/40",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 42,
      "display_id": 41,
      "row_id": 2,
      "column_id": 3,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/41",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/41",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/41",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 43,
      "display_id": 42,
      "row_id": 1,
      "column_id": 4,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/42",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/42",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/42",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 44,
      "display_id": 43,
      "row_id": 2,
      "column_id": 4,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/43",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/43",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/43",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 45,
      "display_id": 44,
      "row_id": 1,
      "column_id": 5,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/44",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/44",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/44",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 46,
      "display_id": 45,
      "row_id": 2,
      "column_id": 5,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/45",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/45",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/45",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 47,
      "display_id": 46,
      "row_id": 1,
      "column_id": 6,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/46",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/46",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/46",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 48,
      "display_id": 47,
      "row_id": 2,
      "column_id": 6,
      "panel_id": 2,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/0/47",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/47",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/0/47",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 49,
      "display_id": 48,
      "row_id": 1,
      "column_id": 1,
      "panel_id": 3,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "qsfp28",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "et-0/1/0",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "et-0/1/0:0",
              "interface_id": 1,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "et-0/1/0:1",
              "interface_id": 2,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "et-0/1/0:2",
              "interface_id": 3,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "et-0/1/0:3",
              "interface_id": 4,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "xe-0/1/0:0",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "xe-0/1/0:1",
              "interface_id": 2,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "xe-0/1/0:2",
              "interface_id": 3,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "xe-0/1/0:3",
              "interface_id": 4,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 0, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 50,
      "display_id": 49,
      "row_id": 2,
      "column_id": 1,
      "panel_id": 3,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "qsfp28",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "et-0/1/1",
              "interface_id": 1,
              "speed": {
                "value": 100,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "et-0/1/1:0",
              "interface_id": 1,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "et-0/1/1:1",
              "interface_id": 2,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "et-0/1/1:2",
              "interface_id": 3,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "et-0/1/1:3",
              "interface_id": 4,
              "speed": {
                "value": 25,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"25g\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "xe-0/1/1:0",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "xe-0/1/1:1",
              "interface_id": 2,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "xe-0/1/1:2",
              "interface_id": 3,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            },
            {
              "name": "xe-0/1/1:3",
              "interface_id": 4,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"breakout\": true, \"fpc\": 0, \"pic\": 1, \"port\": 1, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 51,
      "display_id": 50,
      "row_id": 1,
      "column_id": 1,
      "panel_id": 4,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/2/0",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/0",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/0",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 52,
      "display_id": 51,
      "row_id": 1,
      "column_id": 2,
      "panel_id": 4,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/2/1",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/1",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/1",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 53,
      "display_id": 52,
      "row_id": 1,
      "column_id": 3,
      "panel_id": 4,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/2/2",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/2",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/2",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    },
    {
      "port_id": 54,
      "display_id": 53,
      "row_id": 1,
      "column_id": 4,
      "panel_id": 4,
      "slot_id": 0,
      "failure_domain_id": 1,
      "connector_type": "sfp+",
      "transformations": [
        {
          "transformation_id": 1,
          "is_default": true,
          "interfaces": [
            {
              "name": "xe-0/2/3",
              "interface_id": 1,
              "speed": {
                "value": 10,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
            }
          ]
        },
        {
          "transformation_id": 2,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/3",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"1g\"}}"
            }
          ]
        },
        {
          "transformation_id": 3,
          "is_default": false,
          "interfaces": [
            {
              "name": "ge-0/2/3",
              "interface_id": 1,
              "speed": {
                "value": 1,
                "unit": "G"
              },
              "state": "active",
              "setting": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"auto_negotiation\": false, \"link_mode\": \"full-duplex\", \"speed\": \"1g\"}}"
            }
          ]
        }
      ]
    }
  ],
  "hardware_capabilities": {
    "asic": "T3",
    "ecmp_limit": 64,
    "ram": 4,
    "routing_instance_supported": [
      {
        "value": true,
        "version": ".*"
      }
    ],
    "cpu": "x86",
    "form_factor": "1RU",
    "userland": 64
  },
  "software_capabilities": {
    "onie": false,
    "lxc_support": false,
    "config_apply_support": "complete_only"
  },
  "selector": {
    "os": "Junos",
    "os_version": ".*",
    "manufacturer": "Juniper",
    "model": "EX4400-48F"
  },
  "reference_design_capabilities": {
    "datacenter": "full_support",
    "freeform": "full_support"
  },
  "device_profile_type": "monolithic"
}`
