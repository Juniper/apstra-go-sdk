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
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedesign "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/design"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/Juniper/apstra-go-sdk/internal/zero"
	"github.com/stretchr/testify/require"
)

var testConfiglets = map[string]design.Configlet{
	"flow_data": {
		Label:    testutils.RandString(6, "hex"),
		RefArchs: []enum.RefDesign{enum.RefDesignDatacenter},
		Generators: []design.ConfigletGenerator{
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
	},
	"flow_snmpv2": {
		Label:    testutils.RandString(6, "hex"),
		RefArchs: []enum.RefDesign{enum.RefDesignDatacenter},
		Generators: []design.ConfigletGenerator{
			{
				ConfigStyle: enum.ConfigletStyleJunos,
				Section:     enum.ConfigletSectionSystem,
				TemplateText: `{% if snmpv2_community is defined and snmpv2_community %}
snmp {
    {# This is an example SNMPv2 configlet which can be paired with the example
       SNMPv2 property-set. This will assist Apstra flow analytics to dereference
       interface names exported within flows to their originating interface names.
       Please consult documentation for a more secure solution than SNMPv2 with
       simple well-known communities. This configlet is only an example. #}
    community {{ snmpv2_community }};
}
{% endif %}
`,
			},
		},
	},
}

func TestConfiglet_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.Configlet
		update design.Configlet
	}

	testCases := map[string]testCase{
		"flow_data_to_flow_snmpv2": {
			create: testConfiglets["flow_data"],
			update: testConfiglets["flow_snmpv2"],
		},
		"flow_snmpv2_to_flow_data": {
			create: testConfiglets["flow_snmpv2"],
			update: testConfiglets["flow_data"],
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
					var obj design.Configlet

					// create the object
					id, err = client.Client.CreateConfiglet2(ctx, create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteConfiglet2(ctx, id)
						return nil
					})

					// retrieve the object by ID and validate
					obj, err = client.Client.GetConfiglet2(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Configlet(t, create, obj)

					// retrieve the object by label and validate
					obj, err = client.Client.GetConfigletByLabel2(ctx, create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Configlet(t, create, obj)

					// retrieve the list of IDs - ours must be in there
					ids, err := client.Client.ListConfiglets2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := client.Client.GetConfiglets2(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Configlet(t, create, obj)

					// update the object and validate
					update.SetID(id)
					require.NotNil(t, update.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateConfiglet2(ctx, update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = client.Client.GetConfiglet2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Configlet(t, update, obj)

					// restore the object to the original state
					create.SetID(id)
					require.NotNil(t, create.ID())
					require.Equal(t, id, *update.ID())
					err = client.Client.UpdateConfiglet2(ctx, create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = client.Client.GetConfiglet2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedesign.Configlet(t, create, obj)

					// delete the object
					err = client.Client.DeleteConfiglet2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetConfiglet2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetConfigletByLabel2(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListConfiglets2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = client.Client.GetConfiglets2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateConfiglet2(ctx, update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteConfiglet2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
