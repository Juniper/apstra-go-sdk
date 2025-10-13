// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	timeutils "github.com/Juniper/apstra-go-sdk/internal/time_utils"
)

var configletTestFlowData = Configlet{
	id:             "id_configletTestFlowData",
	createdAt:      pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2006-01-02T15:04:00.000000Z")),
	lastModifiedAt: pointer.To(timeutils.TimeParseMust("2006-01-02T15:04:05.000000Z", "2016-01-02T15:04:00.000000Z")),
	Label:          "Flow Data For Optional Flow Analytics",
	RefArchs:       []enum.RefDesign{enum.RefDesignDatacenter},
	Generators: []ConfigletGenerator{
		{
			ConfigStyle: enum.ConfigletStyleJunos,
			Section:     enum.ConfigletSectionSystem,
			TemplateText: `{% if not os_version.endswith("-EVO") %}
routing-options {
    static {
        route {{collector_ip}}/32 next-table mgmt_junos.inet.0;
    }
}
protocols {
    sflow {
        polling-interval 10;
        sample-rate {
            ingress 1024;
            egress 1024;
        }
    {% if management_ip is defined and management_ip %}
        source-ip {{management_ip}};
    {% endif %}
        collector {{collector_ip}} {
            udp-port 6343;
        }
    {% for interface, settings in portSetting.items() %}
        {% if settings['state'] == 'active' %}
        interfaces {{ interface }};
        {% endif %}
    {% endfor %}
    }
}
{% endif %}
`,
		},
	},
}

const configletTestFlowDataJSON = `{
  "id": "id_configletTestFlowData",
  "ref_archs": [
    "two_stage_l3clos"
  ],
  "created_at": "2006-01-02T15:04:00.000000Z",
  "last_modified_at": "2016-01-02T15:04:00.000000Z",
  "display_name": "Flow Data For Optional Flow Analytics",
  "generators": [
    {
      "config_style": "junos",
      "template_text": "{% if not os_version.endswith(\"-EVO\") %}\nrouting-options {\n    static {\n        route {{collector_ip}}/32 next-table mgmt_junos.inet.0;\n    }\n}\nprotocols {\n    sflow {\n        polling-interval 10;\n        sample-rate {\n            ingress 1024;\n            egress 1024;\n        }\n    {% if management_ip is defined and management_ip %}\n        source-ip {{management_ip}};\n    {% endif %}\n        collector {{collector_ip}} {\n            udp-port 6343;\n        }\n    {% for interface, settings in portSetting.items() %}\n        {% if settings['state'] == 'active' %}\n        interfaces {{ interface }};\n        {% endif %}\n    {% endfor %}\n    }\n}\n{% endif %}\n",
      "section": "system",
      "negation_template_text": "",
      "filename": ""
    }
  ]
}
`
