// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package device_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/device"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

var testProfiles = map[string]string{
	"Juniper_vEX": `{
  "selector": {
    "os": "Junos",
    "os_version": ".*",
    "manufacturer": "Juniper",
    "model": "VIRTUAL-EX9214"
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
      ],
      "column_id": 1,
      "port_id": 1,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 0,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 2,
      "port_id": 2,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 1,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 3,
      "port_id": 3,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 2,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 4,
      "port_id": 4,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 3,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 5,
      "port_id": 5,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 4,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 6,
      "port_id": 6,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 5,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 7,
      "port_id": 7,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 6,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 8,
      "port_id": 8,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 7,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 9,
      "port_id": 9,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 8,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 10,
      "port_id": 10,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 9,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 11,
      "port_id": 11,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 10,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 12,
      "port_id": 12,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 11,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 13,
      "port_id": 13,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 12,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 14,
      "port_id": 14,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 13,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 15,
      "port_id": 15,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 14,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 16,
      "port_id": 16,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 15,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 17,
      "port_id": 17,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 16,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 18,
      "port_id": 18,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 17,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 19,
      "port_id": 19,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 18,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 20,
      "port_id": 20,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 19,
      "slot_id": 0
    }
  ],
  "label": "_test_ Juniper vEX",
  "physical_device": false
}`,
	"Juniper_EX4400-48F": `{
  "selector": {
    "os": "Junos",
    "os_version": ".*",
    "manufacturer": "Juniper",
    "model": "EX4400-48F"
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
      ],
      "column_id": 1,
      "port_id": 1,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 0,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 1,
      "port_id": 2,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 1,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 2,
      "port_id": 3,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 2,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 2,
      "port_id": 4,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 3,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 3,
      "port_id": 5,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 4,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 3,
      "port_id": 6,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 5,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 4,
      "port_id": 7,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 6,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 4,
      "port_id": 8,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 7,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 5,
      "port_id": 9,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 8,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 5,
      "port_id": 10,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 9,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 6,
      "port_id": 11,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 10,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 6,
      "port_id": 12,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 11,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 7,
      "port_id": 13,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 12,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 7,
      "port_id": 14,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 13,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 8,
      "port_id": 15,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 14,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 8,
      "port_id": 16,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 15,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 9,
      "port_id": 17,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 16,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 9,
      "port_id": 18,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 17,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 10,
      "port_id": 19,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 18,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 10,
      "port_id": 20,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 19,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 11,
      "port_id": 21,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 20,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 11,
      "port_id": 22,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 21,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 12,
      "port_id": 23,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 22,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 12,
      "port_id": 24,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 23,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 13,
      "port_id": 25,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 24,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 13,
      "port_id": 26,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 25,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 14,
      "port_id": 27,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 26,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 14,
      "port_id": 28,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 27,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 15,
      "port_id": 29,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 28,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 15,
      "port_id": 30,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 29,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 16,
      "port_id": 31,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 30,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 16,
      "port_id": 32,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 31,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 17,
      "port_id": 33,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 32,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 17,
      "port_id": 34,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 33,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 18,
      "port_id": 35,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 34,
      "slot_id": 0
    },
    {
      "connector_type": "sfp",
      "panel_id": 1,
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
      ],
      "column_id": 18,
      "port_id": 36,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 35,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 1,
      "port_id": 37,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 36,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 1,
      "port_id": 38,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 37,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 2,
      "port_id": 39,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 38,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 2,
      "port_id": 40,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 39,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 3,
      "port_id": 41,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 40,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 3,
      "port_id": 42,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 41,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 4,
      "port_id": 43,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 42,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 4,
      "port_id": 44,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 43,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 5,
      "port_id": 45,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 44,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 5,
      "port_id": 46,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 45,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 6,
      "port_id": 47,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 46,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 2,
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
      ],
      "column_id": 6,
      "port_id": 48,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 47,
      "slot_id": 0
    },
    {
      "connector_type": "qsfp28",
      "panel_id": 3,
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
      ],
      "column_id": 1,
      "port_id": 49,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 48,
      "slot_id": 0
    },
    {
      "connector_type": "qsfp28",
      "panel_id": 3,
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
      ],
      "column_id": 1,
      "port_id": 50,
      "row_id": 2,
      "failure_domain_id": 1,
      "display_id": 49,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 4,
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
      ],
      "column_id": 1,
      "port_id": 51,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 50,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 4,
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
      ],
      "column_id": 2,
      "port_id": 52,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 51,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 4,
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
      ],
      "column_id": 3,
      "port_id": 53,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 52,
      "slot_id": 0
    },
    {
      "connector_type": "sfp+",
      "panel_id": 4,
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
      ],
      "column_id": 4,
      "port_id": 54,
      "row_id": 1,
      "failure_domain_id": 1,
      "display_id": 53,
      "slot_id": 0
    }
  ],
  "label": "_test_ Juniper_EX4400-48F",
  "physical_device": true
}`,
}

func TestProfile_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create string
		update string
	}

	testCases := map[string]testCase{
		"vEX_to_EX4400": {
			create: testProfiles["Juniper_vEX"],
			update: testProfiles["Juniper_EX4400-48F"],
		},
	}

	for tName, tCase := range testCases {
		var create, update device.Profile
		require.NoError(t, json.Unmarshal([]byte(tCase.create), &create)) // extract test object from JSON sample from API
		require.NoError(t, json.Unmarshal([]byte(tCase.update), &update)) // extract test object from JSON sample from API
		create.Label = testutils.RandString(6, "hex")                     // randomize test object name to avoid collisions
		update.Label = testutils.RandString(6, "hex")                     // randomize test object name to avoid collisions

		t.Run(tName, func(t *testing.T) {
			for _, client := range clients {
				create, update := create, update
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					// remove features not supported by earlier API versions
					if !compatibility.DeviceProfileHasRefdesignCapabilities.Check(version.Must(version.NewVersion(client.Client.ApiVersion()))) {
						create.ReferenceDesignCapabilities = nil
						update.ReferenceDesignCapabilities = nil
					}

					var id string
					var err error
					var obj device.Profile

					// create the object
					id, err = client.Client.CreateDeviceProfile(ctx, create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, time.Minute, func(ctx context.Context) error {
						_ = client.Client.DeleteDeviceProfile(ctx, id)
						return nil
					})

					// retrieve the object by ID and validate
					obj, err = client.Client.GetDeviceProfile(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.DeviceProfile(t, create, obj)

					// retrieve the object by label and validate
					obj, err = client.Client.GetDeviceProfileByLabel(ctx, create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.DeviceProfile(t, create, obj)

					// retrieve the list of IDs - ours must be in there
					ids, err := client.Client.ListDeviceProfiles(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := client.Client.GetDeviceProfiles(ctx)
					require.NoError(t, err)
					objPtr := slice.ObjectWithID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.DeviceProfile(t, create, obj)

					// update the object and validate
					require.NoError(t, update.SetID(id))
					require.NotNil(t, update.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateDeviceProfile(ctx, update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = client.Client.GetDeviceProfile(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.DeviceProfile(t, update, obj)

					// restore the object to the original state
					require.NoError(t, create.SetID(id))
					require.NotNil(t, create.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateDeviceProfile(ctx, create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = client.Client.GetDeviceProfile(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.DeviceProfile(t, create, obj)

					// delete the object
					err = client.Client.DeleteDeviceProfile(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetDeviceProfile(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetDeviceProfileByLabel(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListDeviceProfiles(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = client.Client.GetDeviceProfiles(ctx)
					require.NoError(t, err)
					objPtr = slice.ObjectWithID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateDeviceProfile(ctx, update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteDeviceProfile(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
