// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"fmt"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var interfaceMapTestX = InterfaceMap{
	Label:           "Juniper_QFX5120-32C__AOS-2x10-1",
	DeviceProfileID: "Juniper_QFX5120-32C_Junos",
	LogicalDeviceID: "AOS-2x10-1",
	Interfaces: []InterfaceMapInterface{
		{
			Name:     "xe-0/0/32",
			Roles:    LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess},
			Position: 1,
			State:    enum.InterfaceMapInterfaceStateActive,
			Speed:    "10G",
			Setting: struct {
				Param string `json:"param"`
			}{
				Param: `{"global": {"speed": ""}, "interface": {"speed": ""}}`,
			},
			Mapping: InterfaceMapInterfaceMapping{
				DeviceProfilePortID:      33,
				DeviceProfileTransformID: 1,
				DeviceProfileInterfaceID: 1,
				LogicalDevicePanel:       pointer.To(1),
				LogicalDevicePort:        pointer.To(1),
			},
		},
		{
			Name:     "xe-0/0/33",
			Roles:    LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess},
			Position: 2,
			State:    enum.InterfaceMapInterfaceStateActive,
			Speed:    "10G",
			Setting: struct {
				Param string `json:"param"`
			}{
				Param: `{"global": {"speed": ""}, "interface": {"speed": ""}}`,
			},
			Mapping: InterfaceMapInterfaceMapping{
				DeviceProfilePortID:      34,
				DeviceProfileTransformID: 1,
				DeviceProfileInterfaceID: 1,
				LogicalDevicePanel:       pointer.To(1),
				LogicalDevicePort:        pointer.To(2),
			},
		},
	},
	id:             "id_Juniper_QFX5120-32C__AOS-2x10-1",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
}

func init() {
	for i := range 32 {
		interfaceMapTestX.Interfaces = append(interfaceMapTestX.Interfaces, InterfaceMapInterface{
			Name:     fmt.Sprintf("et-0/0/%d", i),
			Roles:    LogicalDevicePortRoles{enum.PortRoleUnused},
			Position: i + 3,
			State:    enum.InterfaceMapInterfaceStateActive,
			Speed:    "100G",
			Setting: struct {
				Param string `json:"param"`
			}{
				Param: fmt.Sprintf(`{"global": {"breakout": false, "fpc": 0, "pic": 0, "port": %d, "speed": "100g"}, "interface": {"speed": ""}}`, i),
			},
			Mapping: InterfaceMapInterfaceMapping{
				DeviceProfilePortID:      i + 1,
				DeviceProfileTransformID: 1,
				DeviceProfileInterfaceID: 1,
			},
		})
	}
}

const interfaceMapTestXJSON = `{
  "label": "Juniper_QFX5120-32C__AOS-2x10-1",
  "device_profile_id": "Juniper_QFX5120-32C_Junos",
  "logical_device_id": "AOS-2x10-1",
  "id": "id_Juniper_QFX5120-32C__AOS-2x10-1",
  "created_at": "2006-01-02T15:04:05.000000Z",
  "last_modified_at": "2016-01-02T15:04:05.000000Z",
  "interfaces": [
    {
      "name": "xe-0/0/32",
      "roles": [
        "leaf",
        "access"
      ],
      "position": 1,
      "state": "active",
      "mapping": [
        33,
        1,
        1,
        1,
        1
      ],
      "speed": {
        "value": 10,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "xe-0/0/33",
      "roles": [
        "leaf",
        "access"
      ],
      "position": 2,
      "state": "active",
      "mapping": [
        34,
        1,
        1,
        1,
        2
      ],
      "speed": {
        "value": 10,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"speed\": \"\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/0",
      "roles": [
        "unused"
      ],
      "position": 3,
      "state": "active",
      "mapping": [
        1,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 0, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/1",
      "roles": [
        "unused"
      ],
      "position": 4,
      "state": "active",
      "mapping": [
        2,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 1, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/2",
      "roles": [
        "unused"
      ],
      "position": 5,
      "state": "active",
      "mapping": [
        3,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 2, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/3",
      "roles": [
        "unused"
      ],
      "position": 6,
      "state": "active",
      "mapping": [
        4,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 3, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/4",
      "roles": [
        "unused"
      ],
      "position": 7,
      "state": "active",
      "mapping": [
        5,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 4, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/5",
      "roles": [
        "unused"
      ],
      "position": 8,
      "state": "active",
      "mapping": [
        6,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 5, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/6",
      "roles": [
        "unused"
      ],
      "position": 9,
      "state": "active",
      "mapping": [
        7,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 6, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/7",
      "roles": [
        "unused"
      ],
      "position": 10,
      "state": "active",
      "mapping": [
        8,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 7, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/8",
      "roles": [
        "unused"
      ],
      "position": 11,
      "state": "active",
      "mapping": [
        9,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 8, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/9",
      "roles": [
        "unused"
      ],
      "position": 12,
      "state": "active",
      "mapping": [
        10,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 9, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/10",
      "roles": [
        "unused"
      ],
      "position": 13,
      "state": "active",
      "mapping": [
        11,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 10, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/11",
      "roles": [
        "unused"
      ],
      "position": 14,
      "state": "active",
      "mapping": [
        12,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 11, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/12",
      "roles": [
        "unused"
      ],
      "position": 15,
      "state": "active",
      "mapping": [
        13,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 12, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/13",
      "roles": [
        "unused"
      ],
      "position": 16,
      "state": "active",
      "mapping": [
        14,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 13, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/14",
      "roles": [
        "unused"
      ],
      "position": 17,
      "state": "active",
      "mapping": [
        15,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 14, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/15",
      "roles": [
        "unused"
      ],
      "position": 18,
      "state": "active",
      "mapping": [
        16,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 15, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/16",
      "roles": [
        "unused"
      ],
      "position": 19,
      "state": "active",
      "mapping": [
        17,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 16, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/17",
      "roles": [
        "unused"
      ],
      "position": 20,
      "state": "active",
      "mapping": [
        18,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 17, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/18",
      "roles": [
        "unused"
      ],
      "position": 21,
      "state": "active",
      "mapping": [
        19,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 18, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/19",
      "roles": [
        "unused"
      ],
      "position": 22,
      "state": "active",
      "mapping": [
        20,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 19, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/20",
      "roles": [
        "unused"
      ],
      "position": 23,
      "state": "active",
      "mapping": [
        21,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 20, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/21",
      "roles": [
        "unused"
      ],
      "position": 24,
      "state": "active",
      "mapping": [
        22,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 21, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/22",
      "roles": [
        "unused"
      ],
      "position": 25,
      "state": "active",
      "mapping": [
        23,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 22, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/23",
      "roles": [
        "unused"
      ],
      "position": 26,
      "state": "active",
      "mapping": [
        24,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 23, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/24",
      "roles": [
        "unused"
      ],
      "position": 27,
      "state": "active",
      "mapping": [
        25,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 24, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/25",
      "roles": [
        "unused"
      ],
      "position": 28,
      "state": "active",
      "mapping": [
        26,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 25, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/26",
      "roles": [
        "unused"
      ],
      "position": 29,
      "state": "active",
      "mapping": [
        27,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 26, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/27",
      "roles": [
        "unused"
      ],
      "position": 30,
      "state": "active",
      "mapping": [
        28,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 27, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/28",
      "roles": [
        "unused"
      ],
      "position": 31,
      "state": "active",
      "mapping": [
        29,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 28, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/29",
      "roles": [
        "unused"
      ],
      "position": 32,
      "state": "active",
      "mapping": [
        30,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 29, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/30",
      "roles": [
        "unused"
      ],
      "position": 33,
      "state": "active",
      "mapping": [
        31,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 30, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    },
    {
      "name": "et-0/0/31",
      "roles": [
        "unused"
      ],
      "position": 34,
      "state": "active",
      "mapping": [
        32,
        1,
        1,
        null,
        null
      ],
      "speed": {
        "value": 100,
        "unit": "G"
      },
      "setting": {
        "param": "{\"global\": {\"breakout\": false, \"fpc\": 0, \"pic\": 0, \"port\": 31, \"speed\": \"100g\"}, \"interface\": {\"speed\": \"\"}}"
      }
    }
  ]
}`
