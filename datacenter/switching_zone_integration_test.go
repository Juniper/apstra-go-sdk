// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter_test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedatacenter "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/datacenter"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestSwitchingZone_CRUD(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		constraints []compatibility.Constraint
		create      datacenter.SwitchingZone
		update      datacenter.SwitchingZone
	}

	testCases := map[string]testCase{
		"vlan_aware": {
			constraints: []compatibility.Constraint{compatibility.DatacenterSwitchingZoneOK},
			create: datacenter.SwitchingZone{
				Label:             pointer.To(testutils.RandString(6, "hex")),
				MACVRFDescription: pointer.To(testutils.RandString(6, "hex")),
				MACVRFName:        pointer.To(testutils.RandString(6, "hex")),
				MACVRFServiceType: pointer.To(enum.SwitchingZoneMACVRFServiceTypeVLANAware),
				RouteTarget:       pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)), Tags: nil,
			},
			update: datacenter.SwitchingZone{
				Label:             pointer.To(testutils.RandString(6, "hex")),
				MACVRFDescription: pointer.To(testutils.RandString(6, "hex")),
				RouteTarget:       pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)), Tags: nil,
			},
		},
		"vlan_bundle": {
			constraints: []compatibility.Constraint{compatibility.DatacenterSwitchingZoneOK},
			create: datacenter.SwitchingZone{
				Label:             pointer.To(testutils.RandString(6, "hex")),
				MACVRFDescription: pointer.To(testutils.RandString(6, "hex")),
				MACVRFName:        pointer.To(testutils.RandString(6, "hex")),
				MACVRFServiceType: pointer.To(enum.SwitchingZoneMACVRFServiceTypeVLANBundle),
				RouteTarget:       pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)), Tags: nil,
			},
			update: datacenter.SwitchingZone{
				Label:             pointer.To(testutils.RandString(6, "hex")),
				MACVRFDescription: pointer.To(testutils.RandString(6, "hex")),
				RouteTarget:       pointer.To(fmt.Sprintf("%d:%d", rand.Intn(math.MaxUint16)+1, rand.Intn(math.MaxUint16)+1)), Tags: nil,
			},
		},
	}

	for tName, tCase := range testCases {
		ctx := testutils.ContextWithTestID(ctx, t)
		t.Run(tName, func(t *testing.T) {
			for _, client := range clients {
				t.Run(client.Name(), func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					for i, constraint := range tCase.constraints {
						if !constraint.Check(client.APIVersion()) {
							t.Skipf("skipping %s test due to constraint %d: %q", tName, i, constraint.String())
						}
					}

					create, update := tCase.create, tCase.update // because we modify these values below

					bp := dctestobj.TestBlueprintA(t, ctx, client.Client)

					// create the object
					id, err := bp.CreateSwitchingZone(ctx, create)
					require.NoError(t, err)
					require.NotEmpty(t, id)

					// retrieve the object by ID and validate
					obj, err := bp.GetSwitchingZone(ctx, id)
					require.NoError(t, err)
					idPtr := obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedatacenter.SwitchingZone(t, create, obj)

					// retrieve the object by label and validate
					if create.Label != nil {
						obj, err = bp.GetSwitchingZoneByLabel(ctx, *create.Label)
						require.NoError(t, err)
						idPtr = obj.ID()
						require.NotNil(t, idPtr)
						require.Equal(t, id, *idPtr)
						comparedatacenter.SwitchingZone(t, create, obj)
					}
					// retrieve the list of IDs - ours must be in there
					ids, err := bp.ListSwitchingZones(ctx)
					require.NoError(t, err)
					require.Contains(t, ids, id)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := bp.GetSwitchingZones(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					obj = *objPtr
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedatacenter.SwitchingZone(t, create, obj)

					// update the object and validate
					require.NoError(t, update.SetID(id))
					require.NotNil(t, update.ID())
					require.Equal(t, id, *update.ID())
					err = bp.UpdateSwitchingZone(ctx, update)
					require.NoError(t, err)

					// retrieve the updated object by ID and validate
					obj, err = bp.GetSwitchingZone(ctx, id)
					require.NoError(t, err)
					idPtr = obj.ID()
					require.NotNil(t, idPtr)
					require.Equal(t, id, *idPtr)
					comparedatacenter.SwitchingZone(t, update, obj)

					// delete the object
					err = bp.DeleteSwitchingZone(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = bp.GetSwitchingZone(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					if create.Label != nil {
						_, err = bp.GetSwitchingZoneByLabel(ctx, *create.Label)
						require.Error(t, err)
						require.ErrorAs(t, err, &ace)
						require.Equal(t, apstra.ErrNotfound, ace.Type())
					}

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = bp.ListSwitchingZones(ctx)
					require.NoError(t, err)
					require.NotContains(t, ids, id)

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = bp.GetSwitchingZones(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					err = bp.UpdateSwitchingZone(ctx, update)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = bp.DeleteSwitchingZone(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}

func TestSwitchingZone_GetDefaultSwitchingZone(t *testing.T) {
	ctx := testutils.ContextWithTestID(context.Background(), t)
	clients := testclient.GetTestClients(t, ctx)
	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			if !compatibility.DatacenterSwitchingZoneOK.Check(client.APIVersion()) {
				t.Skipf("skipping test due to compatibility constraint: %q", compatibility.DatacenterSwitchingZoneOK.String())
			}

			bp := dctestobj.TestBlueprintA(t, ctx, client.Client)

			sz, err := bp.GetDefaultSwitchingZone(ctx)
			require.NoError(t, err)
			require.NotNil(t, sz.ID())
			require.NotEmpty(t, *sz.ID())
		})
	}
}
