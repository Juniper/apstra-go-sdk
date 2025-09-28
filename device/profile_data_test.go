// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"time"
)

var testProfileGeneric1x10 = Profile{
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
		Asic:       "",
		ECMPLimit:  64,
		Ram:        16,
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
			Port:          1,
			Row:           1,
			FailureDomain: 1,
			Display:       pointer.To(1),
			Slot:          0,
		},
	},
	Label:          "test_Generic_Server_1RU_1x10G",
	PhysicalDevice: true,
	id:             "_test_Generic_Server_1RU_1x10G",
	createdAt:      pointer.To(time.Date(2024, time.January, 2, 3, 4, 5, 6000, time.UTC)),
	lastModifiedAt: pointer.To(time.Date(2025, time.January, 2, 3, 4, 5, 6000, time.UTC)),
}

const testProfileGeneric1x10JSON = `{
  "last_modified_at": "2025-01-02T03:04:05.000006Z",
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
  "label": "test_Generic_Server_1RU_1x10G",
  "id": "_test_Generic_Server_1RU_1x10G",
  "created_at": "2024-01-02T03:04:05.000006Z",
  "physical_device": true
}`
