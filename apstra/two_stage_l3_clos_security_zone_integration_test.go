// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package apstra_test

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/compatibility"
	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	comparedatacenter "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/two_stage_l3_clos"
	dctestobj "github.com/Juniper/apstra-go-sdk/internal/test_utils/datacenter_test_objects"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

func TestCRUDSecurityZone(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		versionConstraint *compatibility.Constraint
		create            apstra.SecurityZone
		update            *apstra.SecurityZone
	}

	testCases := map[string]testCase{
		"start_minimal_4.x+": {
			create: apstra.SecurityZone{
				Label:   testutils.RandString(6, "hex"),
				VRFName: testutils.RandString(6, "hex"),
				Type:    enum.SecurityZoneTypeEVPN,
			},
			update: &apstra.SecurityZone{
				Label:            testutils.RandString(8, "hex"),
				Type:             enum.SecurityZoneTypeEVPN,
				RoutingPolicyID:  "",
				RTPolicy:         nil,
				VLAN:             pointer.To(apstra.VLAN(4000)),
				VNI:              pointer.To(40000),
				JunosEVPNIRBMode: pointer.To(enum.JunosEVPNIRBModeSymmetric),
			},
		},
		"start_maximal_4.x+": {
			create: apstra.SecurityZone{
				Label:            testutils.RandString(8, "hex"),
				Type:             enum.SecurityZoneTypeEVPN,
				RoutingPolicyID:  "",
				RTPolicy:         nil,
				VLAN:             pointer.To(apstra.VLAN(4001)),
				VNI:              pointer.To(40010),
				JunosEVPNIRBMode: pointer.To(enum.JunosEVPNIRBModeSymmetric),
				VRFName:          testutils.RandString(6, "hex"),
			},
			update: &apstra.SecurityZone{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.SecurityZoneTypeEVPN,
			},
		},
		"start_minimal_5.x+": {
			versionConstraint: &compatibility.SecurityZoneDescriptionSupported,
			create: apstra.SecurityZone{
				Label:   testutils.RandString(6, "hex"),
				VRFName: testutils.RandString(6, "hex"),
				Type:    enum.SecurityZoneTypeEVPN,
			},
			update: &apstra.SecurityZone{
				Description:      pointer.To(testutils.RandString(8, "hex")),
				Label:            testutils.RandString(8, "hex"),
				Type:             enum.SecurityZoneTypeEVPN,
				RoutingPolicyID:  "",
				RTPolicy:         nil,
				VLAN:             pointer.To(apstra.VLAN(4002)),
				VNI:              pointer.To(40020),
				JunosEVPNIRBMode: pointer.To(enum.JunosEVPNIRBModeSymmetric),
			},
		},
		"start_maximal_5.x+": {
			versionConstraint: &compatibility.SecurityZoneDescriptionSupported,
			create: apstra.SecurityZone{
				Description:      pointer.To(testutils.RandString(8, "hex")),
				Label:            testutils.RandString(8, "hex"),
				Type:             enum.SecurityZoneTypeEVPN,
				RoutingPolicyID:  "",
				RTPolicy:         nil,
				VLAN:             pointer.To(apstra.VLAN(4003)),
				VNI:              pointer.To(40030),
				JunosEVPNIRBMode: pointer.To(enum.JunosEVPNIRBModeSymmetric),
				VRFName:          testutils.RandString(6, "hex"),
			},
			update: &apstra.SecurityZone{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.SecurityZoneTypeEVPN,
			},
		},
		"start_minimal_6.1+": {
			versionConstraint: &compatibility.SecurityZoneAddressingSupported,
			create: apstra.SecurityZone{
				Label:   testutils.RandString(6, "hex"),
				VRFName: testutils.RandString(6, "hex"),
				Type:    enum.SecurityZoneTypeEVPN,
			},
			update: &apstra.SecurityZone{
				Description:       pointer.To(testutils.RandString(8, "hex")),
				Label:             testutils.RandString(8, "hex"),
				Type:              enum.SecurityZoneTypeEVPN,
				RoutingPolicyID:   "",
				RTPolicy:          nil,
				VLAN:              pointer.To(apstra.VLAN(4004)),
				VNI:               pointer.To(40040),
				JunosEVPNIRBMode:  pointer.To(enum.JunosEVPNIRBModeSymmetric),
				AddressingSupport: &enum.AddressingSchemeIPv6,
				DisableIPv4:       pointer.To(true),
			},
		},
		"start_maximmal_6.1+": {
			versionConstraint: &compatibility.SecurityZoneAddressingSupported,
			create: apstra.SecurityZone{
				Description:       pointer.To(testutils.RandString(8, "hex")),
				Label:             testutils.RandString(8, "hex"),
				VRFName:           testutils.RandString(6, "hex"),
				Type:              enum.SecurityZoneTypeEVPN,
				RoutingPolicyID:   "",
				RTPolicy:          nil,
				VLAN:              pointer.To(apstra.VLAN(4005)),
				VNI:               pointer.To(40050),
				JunosEVPNIRBMode:  pointer.To(enum.JunosEVPNIRBModeSymmetric),
				AddressingSupport: &enum.AddressingSchemeIPv46,
				DisableIPv4:       pointer.To(false),
			},
			update: &apstra.SecurityZone{
				Label: testutils.RandString(6, "hex"),
				Type:  enum.SecurityZoneTypeEVPN,
			},
		},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(client.APIVersion()) {
						t.Skipf("skipping %q due to version constraints: %q. API version: %q",
							tName, tCase.versionConstraint, client.Client.ApiVersion())
					}

					// because we modify these values below
					create := tCase.create
					create.VRFName = create.Label
					var update *apstra.SecurityZone
					if tCase.update != nil {
						update = pointer.To(*tCase.update)
						require.Empty(t, update.VRFName, "vrf name must not be set in test case 'update'")
						update.VRFName = create.VRFName
					}

					var id string
					var err error
					var obj apstra.SecurityZone

					// create the object
					id, err = bpClient.CreateSecurityZone(ctx, create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = bpClient.GetSecurityZone(ctx, id)
					require.NoError(t, err)
					comparedatacenter.SecurityZone(t, create, obj)
					require.NotNil(t, obj.ID())
					require.Equal(t, id, *obj.ID())

					// retrieve the object by label and validate
					obj, err = bpClient.GetSecurityZoneByLabel(ctx, create.Label)
					require.NoError(t, err)
					comparedatacenter.SecurityZone(t, create, obj)
					require.NotNil(t, obj.ID())
					require.Equal(t, id, *obj.ID())

					// retrieve the object by vrf name and validate
					obj, err = bpClient.GetSecurityZoneByVRFName(ctx, create.VRFName)
					require.NoError(t, err)
					comparedatacenter.SecurityZone(t, create, obj)
					require.NotNil(t, obj.ID())
					require.Equal(t, id, *obj.ID())

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := bpClient.GetSecurityZones(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					require.NotNil(t, obj.ID())
					require.Equal(t, id, *obj.ID())
					comparedatacenter.SecurityZone(t, create, obj)

					if update != nil {
						// update the object
						update.SetID(id)
						require.NotNil(t, update.ID())
						require.Equal(t, id, *update.ID())
						if update.JunosEVPNIRBMode == nil && update.Type == enum.SecurityZoneTypeEVPN {
							update.JunosEVPNIRBMode = pointer.To(*obj.JunosEVPNIRBMode)
						}
						err = bpClient.UpdateSecurityZone(ctx, *update)
						require.NoError(t, err)

						// retrieve the object by ID and validate
						obj, err = bpClient.GetSecurityZone(ctx, id)
						require.NoError(t, err)
						comparedatacenter.SecurityZone(t, *update, obj)
						require.NotNil(t, obj.ID())
						require.Equal(t, id, *obj.ID())

						// restore the object to the original state
						create.SetID(id)
						require.NotNil(t, create.ID())
						require.Equal(t, id, *create.ID())
						if create.JunosEVPNIRBMode == nil && create.Type == enum.SecurityZoneTypeEVPN {
							create.JunosEVPNIRBMode = pointer.To(*obj.JunosEVPNIRBMode)
						}
						err = bpClient.UpdateSecurityZone(ctx, create)
						require.NoError(t, err)

						// retrieve the object by ID and validate
						obj, err = bpClient.GetSecurityZone(ctx, id)
						require.NoError(t, err)
						comparedatacenter.SecurityZone(t, create, obj)
						require.NotNil(t, obj.ID())
						require.Equal(t, id, *obj.ID())
					}

					// delete the object
					err = bpClient.DeleteSecurityZone(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = bpClient.GetSecurityZone(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = bpClient.GetSecurityZoneByLabel(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by VRF name
					_, err = bpClient.GetSecurityZoneByVRFName(ctx, create.VRFName)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = bpClient.GetSecurityZones(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// update the object
					create.SetID(id)
					require.NotNil(t, create.ID())
					require.Equal(t, id, *create.ID())
					err = bpClient.UpdateSecurityZone(ctx, create)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = bpClient.DeleteSecurityZone(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}

func TestGetDefaultRoutingZone(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)

	clients := testclient.GetTestClients(t, ctx)

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bpClient := dctestobj.TestBlueprintA(t, ctx, client.Client)
			sz, err := bpClient.GetSecurityZoneByVRFName(ctx, "default")
			require.NoError(t, err)
			require.Equal(t, "default", sz.VRFName)
		})
	}
}
