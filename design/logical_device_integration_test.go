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
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

var testLogicalDevices = map[string]design.LogicalDevice{
	"spine_32x400": {
		Label: testutils.RandString(6, "hex"),
		Panels: []design.LogicalDevicePanel{
			{
				PanelLayout: design.LogicalDevicePanelLayout{
					RowCount:    2,
					ColumnCount: 16,
				},
				PortGroups: []design.LogicalDevicePanelPortGroup{
					{
						Count: 32,
						Speed: "400G",
						Roles: design.LogicalDevicePortRoles{enum.PortRoleLeaf, enum.PortRoleSuperspine},
					},
				},
				PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
			},
		},
	},
	"leaf_48x25_4x400": {
		Label: testutils.RandString(6, "hex"),
		Panels: []design.LogicalDevicePanel{
			{
				PanelLayout: design.LogicalDevicePanelLayout{
					RowCount:    2,
					ColumnCount: 24,
				},
				PortGroups: []design.LogicalDevicePanelPortGroup{
					{
						Count: 48,
						Speed: "25G",
						Roles: design.LogicalDevicePortRoles{enum.PortRoleGeneric, enum.PortRoleAccess, enum.PortRoleLeaf},
					},
				},
				PortIndexing: enum.DesignLogicalDevicePanelPortIndexingTBLR,
			},
			{
				PanelLayout: design.LogicalDevicePanelLayout{
					RowCount:    2,
					ColumnCount: 3,
				},
				PortGroups: []design.LogicalDevicePanelPortGroup{
					{
						Count: 6,
						Speed: "400G",
						Roles: design.LogicalDevicePortRoles{enum.PortRoleAccess, enum.PortRoleGeneric, enum.PortRoleLeaf, enum.PortRolePeer, enum.PortRoleSpine},
					},
				},
				PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
			},
		},
	},
	"generic_4x25": {
		Label: "leaf_48x25_4x400",
		Panels: []design.LogicalDevicePanel{
			{
				PanelLayout: design.LogicalDevicePanelLayout{
					RowCount:    2,
					ColumnCount: 2,
				},
				PortGroups: []design.LogicalDevicePanelPortGroup{
					{
						Count: 4,
						Speed: "25G",
						Roles: design.LogicalDevicePortRoles{enum.PortRoleAccess, enum.PortRoleLeaf},
					},
				},
				PortIndexing: enum.DesignLogicalDevicePanelPortIndexingLRTB,
			},
		},
	},
}

func TestLogicalDevice_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		create design.LogicalDevice
		update design.LogicalDevice
	}

	testCases := map[string]testCase{
		"spine_to_leaf": {
			create: testLogicalDevices["spine_32x400"],
			update: testLogicalDevices["leaf_48x25_4x400"],
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
					var obj design.LogicalDevice

					// create the object
					id, err = client.Client.CreateLogicalDevice2(ctx, tCase.create)
					require.NoError(t, err)

					// ensure the object is deleted even if tests fail
					testutils.CleanupWithFreshContext(t, 10, func(ctx context.Context) error {
						_ = client.Client.DeleteLogicalDevice2(ctx, id)
						return nil
					})

					// retrieve the object by ID and validate
					obj, err = client.Client.GetLogicalDevice2(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.LogicalDevice2(t, tCase.create, obj)

					// retrieve the object by label and validate
					obj, err = client.Client.GetLogicalDeviceByLabel2(ctx, tCase.create.Label)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.LogicalDevice2(t, tCase.create, obj)

					// retrieve the list of IDs - ours must be in there
					ids, err := client.Client.ListLogicalDevices2(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := client.Client.GetLogicalDevices2(ctx)
					require.NoError(t, err)
					objPtr := slice.ObjectWithID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.LogicalDevice2(t, tCase.create, obj)

					// update the object and validate
					require.NoError(t, tCase.update.SetID(id))
					require.NotNil(t, tCase.update.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateLogicalDevice2(ctx, tCase.update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = client.Client.GetLogicalDevice2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.LogicalDevice2(t, tCase.update, obj)

					// restore the object to the original state
					require.NoError(t, tCase.create.SetID(id))
					require.NotNil(t, tCase.create.ID())
					require.Equal(t, id, *tCase.update.ID())
					err = client.Client.UpdateLogicalDevice2(ctx, tCase.create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = client.Client.GetLogicalDevice2(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					compare.LogicalDevice2(t, tCase.create, obj)

					// delete the object
					err = client.Client.DeleteLogicalDevice2(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = client.Client.GetLogicalDevice2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = client.Client.GetLogicalDeviceByLabel2(ctx, tCase.create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = client.Client.ListLogicalDevices2(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = client.Client.GetLogicalDevices2(ctx)
					require.NoError(t, err)
					objPtr = slice.ObjectWithID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = client.Client.UpdateLogicalDevice2(ctx, tCase.update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = client.Client.DeleteLogicalDevice2(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}
