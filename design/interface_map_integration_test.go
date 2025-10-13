// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package design_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/Juniper/apstra-go-sdk/speed"
	"github.com/stretchr/testify/require"
)

var testInterfaceMaps = map[string]design.InterfaceMap{
	"Juniper_vQFX__AOS-7x10-Leaf": {
		Label:           testutils.RandString(6, "hex"),
		DeviceProfileID: "Juniper_vQFX",
		LogicalDeviceID: "AOS-7x10-Leaf",
		Interfaces: []design.InterfaceMapInterface{
			{Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleSpine}},     // port 1
			{Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleSpine}},     // port 2
			{Roles: design.LogicalDevicePortRoles{enum.PortRolePeer}},                         // port 3
			{Roles: design.LogicalDevicePortRoles{enum.PortRolePeer}},                         // port 4
			{Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}}, // port 5
			{Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess}}, // port 6
			{Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric}},                      // port 7
		},
	},
	"Generic_Server_1RU_2x10G__AOS-2x10-1": {
		Label:           testutils.RandString(6, "hex"),
		DeviceProfileID: "Generic_Server_1RU_2x10G",
		LogicalDeviceID: "AOS-2x10-1",
		Interfaces: []design.InterfaceMapInterface{
			{Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess}}, // port 1
			{Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleAccess}}, // port 1
		},
	},
}

func init() {
	rules := map[string]struct {
		LDPortCount     int
		DPPortCount     int
		ifNameFmtString string
		speed           speed.Speed
		settingParam    string
	}{
		"Juniper_vQFX__AOS-7x10-Leaf": {
			LDPortCount:     7,
			DPPortCount:     12,
			ifNameFmtString: "xe-0/0/%d",
			speed:           "10G",
			settingParam:    `{"global": {"speed": ""}, "interface": {"speed": ""}}`,
		},
		"Generic_Server_1RU_2x10G__AOS-2x10-1": {
			LDPortCount:     2,
			DPPortCount:     2,
			ifNameFmtString: "eth%d",
			speed:           "10G",
			settingParam:    "",
		},
	}

	for k, v := range rules {
		// pull the test object in the map
		im, ok := testInterfaceMaps[k]
		if !ok {
			panic("rule " + k + " not found")
		}

		// add unused interfaces to match device profile
		for len(im.Interfaces) < v.DPPortCount {
			im.Interfaces = append(im.Interfaces, design.InterfaceMapInterface{
				Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric},
			})
		}

		// add boilerplate stuff for each interface
		for i := range im.Interfaces {
			im.Interfaces[i].Name = fmt.Sprintf(v.ifNameFmtString, i)
			im.Interfaces[i].Position = i + 1
			im.Interfaces[i].State = enum.InterfaceMapInterfaceStateActive
			im.Interfaces[i].Speed = v.speed
			im.Interfaces[i].Setting = struct {
				Param string `json:"param"`
			}{Param: v.settingParam}
			im.Interfaces[i].Mapping = design.InterfaceMapInterfaceMapping{
				DeviceProfilePortID:      i + 1,
				DeviceProfileTransformID: 1,
				DeviceProfileInterfaceID: 1,
			}
			if i < v.LDPortCount { // interfaces beyond LD port count do not reference LD
				im.Interfaces[i].Mapping.LogicalDevicePanel = pointer.To(1)
				im.Interfaces[i].Mapping.LogicalDevicePort = pointer.To(i + 1)
			}

		}

		testInterfaceMaps[k] = im // add the modified test object back to the map
	}
}

func TestInterfaceMap_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.InterfaceMap
		update design.InterfaceMap
	}

	testCases := map[string]testCase{
		"7x10_leaf_to_2x10_server": {
			create: testInterfaceMaps["Juniper_vQFX__AOS-7x10-Leaf"],
			update: testInterfaceMaps["Generic_Server_1RU_2x10G__AOS-2x10-1"],
		},
		"2x10_server_to_7x10_leaf": {
			create: testInterfaceMaps["Juniper_vQFX__AOS-7x10-Leaf"],
			update: testInterfaceMaps["Generic_Server_1RU_2x10G__AOS-2x10-1"],
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					var id string
					var err error
					var obj design.InterfaceMap

					// create the object
					id, err = client.Client.CreateInterfaceMap2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteInterfaceMap2(ctx, id)
						return nil
					})

					// retrieve the object by ID and validate
					obj, err = client.Client.GetInterfaceMap2(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.InterfaceMap(t, tCase.create, obj)

					// retrieve the object by label and validate
					obj, err = client.Client.GetInterfaceMapByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.InterfaceMap(t, tCase.create, obj)

					// retrieve the list of IDs - ours must be in there
					ids, err := client.Client.ListInterfaceMaps2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := client.Client.GetInterfaceMaps2(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.InterfaceMap(t, tCase.create, obj)

					// update the object and validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateInterfaceMap2(ctx, tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = client.Client.GetInterfaceMap2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.InterfaceMap(t, tCase.update, obj)

					// restore the object to the original state
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateInterfaceMap2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = client.Client.GetInterfaceMap2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.InterfaceMap(t, tCase.create, obj)

					// delete the object
					err = client.Client.DeleteInterfaceMap2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetInterfaceMap2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetInterfaceMapByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListInterfaceMaps2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = client.Client.GetInterfaceMaps2(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateInterfaceMap2(ctx, tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteInterfaceMap2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
