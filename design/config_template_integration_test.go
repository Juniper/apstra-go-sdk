// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package design_test

import (
	"context"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/design"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/stretchr/testify/require"
)

var testConfigTemplates = map[string]design.ConfigTemplate{
	"junos_system": {
		Label: testutils.RandString(6, "hex") + ".jinja",
		Text:  `{# A system hostname is required on the system for purposes such as hostname and\n   lldp cabling telemetry. #}\n{% if hostname %}\nsystem {\n    host-name {{hostname}};\n}\n{% endif %}\n`,
	},
	"junos_chassis": {
		Label: testutils.RandString(6, "hex") + ".jinja",
		Text:  `{# Count the number of interfaces which are port channels (ae interfaces) #}\n{% set ae_count = interfaces.values()|selectattr('if_type', 'eq', 'port_channel')|list|count %}\n{# junos_chassis is a dictionary mapping from interface_name: port_settings_dict,\n   in which each port_settings_dict entry is shadowed from the associated interface\n   according to its transformation_id.\n\n   For instance, on a Juniper QFX10002, on transformation ID #3 and port ID #2:\n   The example interface 'settings' entry from a device profile transformation for interface xe-0/0/1:0 is:\n\n   {\"global\": {\"breakout\": true, \"fpc\": 0, \"master_port\": 0, \"pic\": 0, \"port\": 1, \"speed\": \"10g\"}, \"interface\": {\"speed\": \"\"}}\n\n   This information is copied into the device model dump as the below.\n   The 'global' component is used by the chassis_config dictionary within\n   junos_chassis.jinja, and most of the 'interface' component is documented and used\n   within junos_interfaces.jinja.  Different platforms have different rendering\n   requirements for channelization and physical port configuration, some requiring\n   chassis config and some requiring statements underneath the 'interfaces' section.\n    \"chassis_config\": {\n        \"0\": {\n            \"0\": {\n                \"1\": {\n                    \"unused\": true\n                },\n                \"0\": {\n                    \"sub-ports-count\": 4,\n                    \"breakout\": true,\n                    \"speed\": \"25g\"\n                },\n                \"3\": {\n                    \"unused\": true\n                },\n                \"2\": {\n                    \"unused\": true\n                }\n            }\n        }\n    },\n\n\n    which will render config as:\n    chassis {\n        fpc 0 {\n            pic 0 {\n                port 0 {\n                    speed 25g;\n                    number-of-sub-ports 4;\n                }\n                port 1 {\n                    unused;\n                }\n                port 2 {\n                    unused;\n                }\n                port 3 {\n                    unused;\n                }\n            }\n        }\n    }\n\n    The freeform blueprint will index this information and add it to the\n    'chassis_config' device model key for usage by junos_chassis.jinja.\n\n    This schema is described as:\n\n    \"chassis_config\": {\n        <fpc_id> {\n        The FPC (Flexible PIC Concentrator) ID that relates to this port\n            <pic_id> {\n            The PIC (Physical Interface Card) ID that relates to this port\n                {port_id}: {\n                    The key 'port_id' represents either the master_port in a\n                    port-group or a regular port which requires chassis speed\n                    configuration. The semantics of the key are already calculated\n                    in the chassis_config in the device model.\n\n                    1) If both 'port_id' and 'master_port' are None\n                       This port does not require any explicit chassis speed\n                       commands rendered and will not be included in the\n                       chassis_config dictionary or rendered underneath the\n                       chassis section jinja.\n                    2) If both 'port_id' and 'master_port' have values\n                       This is a port in a port-group that requires the master's\n                       port speed to be set.  In this case, the port_id will be\n                       taken from 'master_port' for the chassis dictionary for\n                       configuration.\n                    3) If 'port_id' has a value but 'master_port_id' is None\n                       This is a regular port that is not in a port-group, and it\n                       requires explicit chassis speed commands to be set.\n                    4) If 'master_port_id' has a value and 'port_id' is None\n                       This is not a valid combination, and will not be included\n                       in the chassis_config dictionary or rendered by jinja.\n\n                    'speed-keyword': <optional str>,\n                        Some devices require different commands based on the type of\n                        port or the type of hardware platform.  Some devices need\n                        'speed 10g', some need 'channel-speed 10g', and others need\n                        'channelization-speed 10g'.  If the speed-keyword is not\n                        explicitly specified within the device-profile then the\n                        dictionary dump will contain 'channel-speed' as it is the\n                        most common for chassis configurations.\n\n                    'speed': <str>,\n                        The value to use when rendering the 'speed-keyword' command.\n                        Different interface types, device platforms, transceiver\n                        types have different speed capabilities. This value\n                        comes from the device profile.\n                        Examples are 10g | 25g | 40g | 100g | 400g\n\n                    'breakout': <bool>,\n                        Indicates whether a port is broken out or not, which is\n                        used together with 'speed' for port & channel-speed\n                        configuration rendering.\n\n                    'fpc': <int as str>,\n                        The FPC (Flexible PIC Concentrator) ID that relates to this port\n                    'pic': <int as str>,\n                        The PIC (Physical Interface Card) ID that relates to this port\n\n                    'master_port': <Optional int>,\n                        Physical master port ID for a member of a port-group\n\n                    'port': <Optional int>,\n                        Device's physical port ID\n\n                    'sub-ports-count': <optional int>,\n                        If this is a parent of a port-group, this indicates the number\n                        of subports that this master port will be broken out into.\n\n                    'unused-port-list': <optional list of ints>,\n                        On certain platforms, if a port is used for breakout\n                        transformation, the platform may require to set other adjacent\n                        ports as 'unused' underneath the chassis configuration.\n                        This list of port IDs is used to render 'unused' for those\n                        ports.\n                }\n            }\n        }\n    }\n\n#}\n{% if ae_count or chassis_config %}\nchassis {\n    {# Junos requires explicitly configuring device-count for the number of\n       aggregate-ethernet devices which are rendered by junos_interfaces.jinja #}\n    {% if ae_count %}\n    aggregated-devices {\n        ethernet {\n            device-count {{ae_count}};\n        }\n    }\n    {% endif %}\n    {% if chassis_config %}\n        {% for fpc, pic_info in function.sorted_dict(chassis_config) %}\n    fpc {{ fpc }} {\n            {% for pic, port_info in function.sorted_dict(pic_info) %}\n        pic {{ pic }} {\n                {% for port, speed_info in function.sorted_dict(port_info) %}\n            port {{ port }} {\n                    {% if speed_info.get('unused') %}\n                unused;\n                    {% elif speed_info.get('breakout') %}\n                        {% if speed_info.get('sub-ports-count') %}\n                speed {{ speed_info['speed'] }};\n                number-of-sub-ports {{ speed_info['sub-ports-count'] }};\n                        {% else %}\n                {{ speed_info.get('speed-keyword') or 'channel-speed' }} {{ speed_info['speed'] }};\n                        {% endif %}\n                    {% else %}\n                speed {{ speed_info['speed'] }};\n                    {% endif %}\n            }\n                {% endfor %}\n        }\n            {% endfor %}\n    }\n        {% endfor %}\n    {% endif %}\n}\n{% endif %}\n`,
	},
}

func TestConfigTemplate_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.ConfigTemplate
		update design.ConfigTemplate
	}

	testCases := map[string]testCase{
		"junos_system_to_junos_chassis": {
			create: testConfigTemplates["junos_system"],
			update: testConfigTemplates["junos_chassis"],
		},
		"junos_chassis_to_junos_system": {
			create: testConfigTemplates["junos_chassis"],
			update: testConfigTemplates["junos_system"],
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			require.NotEqual(t, tCase.create, zero.Of(tCase.create)) // make sure we didn't use a bogus map key
			require.NotEqual(t, tCase.update, zero.Of(tCase.update)) // make sure we didn't use a bogus map key

			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					create, update := tCase.create, tCase.update // because we modify these values below

					var id string
					var err error
					var obj design.ConfigTemplate

					// create the object
					id, err = client.Client.CreateConfigTemplate2(ctx, create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteConfigTemplate2(ctx, id)
						return nil
					})

					// retrieve the object by ID and validate
					obj, err = client.Client.GetConfigTemplate2(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.ConfigTemplate(t, create, obj)

					// retrieve the object by label and validate
					obj, err = client.Client.GetConfigTemplateByLabel2(ctx, create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.ConfigTemplate(t, create, obj)

					// retrieve the list of IDs - ours must be in there
					ids, err := client.Client.ListConfigTemplates2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := client.Client.GetConfigTemplates2(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.ConfigTemplate(t, create, obj)

					// update the object and validate
					update.SetID(id)
					require.NotNil(t, update.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateConfigTemplate2(ctx, update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = client.Client.GetConfigTemplate2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.ConfigTemplate(t, update, obj)

					// restore the object to the original state
					create.SetID(id)
					require.NotNil(t, create.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateConfigTemplate2(ctx, create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = client.Client.GetConfigTemplate2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.ConfigTemplate(t, create, obj)

					// delete the object
					err = client.Client.DeleteConfigTemplate2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetConfigTemplate2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetConfigTemplateByLabel2(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListConfigTemplates2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = client.Client.GetConfigTemplates2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateConfigTemplate2(ctx, update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteConfigTemplate2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
