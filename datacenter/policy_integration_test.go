// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package datacenter_test

import (
	"context"
	"strings"
	"sync"
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

func TestCrudPolicy(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	type testCase struct {
		versionConstraint *compatibility.Constraint
		create            datacenter.Policy
		update            *datacenter.Policy
	}

	testCases := map[string]testCase{
		"start_minimal_src_rz": {
			create: datacenter.Policy{
				Label: testutils.RandString(6, "hex"),
			},
			update: &datacenter.Policy{
				Enabled:             true,
				Label:               testutils.RandString(6, "hex"),
				Description:         testutils.RandString(6, "hex"),
				SrcApplicationPoint: pointer.To("rz:a"),
				Rules: []datacenter.PolicyRule{
					{
						Label:             testutils.RandString(6, "hex"),
						Description:       testutils.RandString(6, "hex"),
						Protocol:          enum.PolicyRuleProtocolTcp,
						Action:            testutils.OneOf(enum.PolicyRuleActionDeny, enum.PolicyRuleActionDenyLog, enum.PolicyRuleActionPermit, enum.PolicyRuleActionPermitLog),
						SrcPort:           dctestobj.RandomPortRanges(3),
						DstPort:           dctestobj.RandomPortRanges(3),
						TcpStateQualifier: testutils.OneOf(nil, pointer.To(enum.TcpStateQualifierEstablished)),
					},
					{
						Label:       testutils.RandString(6, "hex"),
						Description: testutils.RandString(6, "hex"),
						Protocol:    enum.PolicyRuleProtocolUdp,
						Action:      testutils.OneOf(enum.PolicyRuleActionDeny, enum.PolicyRuleActionDenyLog, enum.PolicyRuleActionPermit, enum.PolicyRuleActionPermitLog),
						SrcPort:     dctestobj.RandomPortRanges(3),
						DstPort:     dctestobj.RandomPortRanges(3),
					},
				},
				Tags: testutils.RandomStrings(3, 6, 6, "hex"),
			},
		},
		"start_maximal_dst_rz": {
			create: datacenter.Policy{
				Enabled:             true,
				Label:               testutils.RandString(6, "hex"),
				Description:         testutils.RandString(6, "hex"),
				DstApplicationPoint: pointer.To("rz:b"),
				Rules: []datacenter.PolicyRule{
					{
						Label:             testutils.RandString(6, "hex"),
						Description:       testutils.RandString(6, "hex"),
						Protocol:          enum.PolicyRuleProtocolTcp,
						Action:            testutils.OneOf(enum.PolicyRuleActionDeny, enum.PolicyRuleActionDenyLog, enum.PolicyRuleActionPermit, enum.PolicyRuleActionPermitLog),
						SrcPort:           dctestobj.RandomPortRanges(3),
						DstPort:           dctestobj.RandomPortRanges(3),
						TcpStateQualifier: testutils.OneOf(nil, pointer.To(enum.TcpStateQualifierEstablished)),
					},
					{
						Label:       testutils.RandString(6, "hex"),
						Description: testutils.RandString(6, "hex"),
						Protocol:    enum.PolicyRuleProtocolUdp,
						Action:      testutils.OneOf(enum.PolicyRuleActionDeny, enum.PolicyRuleActionDenyLog, enum.PolicyRuleActionPermit, enum.PolicyRuleActionPermitLog),
						SrcPort:     dctestobj.RandomPortRanges(3),
						DstPort:     dctestobj.RandomPortRanges(3),
					},
				},
				Tags: testutils.RandomStrings(3, 6, 6, "hex"),
			},
			update: &datacenter.Policy{
				Label: testutils.RandString(6, "hex"),
			},
		},
		"start_minimal_intra_rz_vns": {
			create: datacenter.Policy{
				Label: testutils.RandString(6, "hex"),
			},
			update: &datacenter.Policy{
				Enabled:             true,
				Label:               testutils.RandString(6, "hex"),
				Description:         testutils.RandString(6, "hex"),
				SrcApplicationPoint: pointer.To("vn:a:a1"),
				DstApplicationPoint: pointer.To("vn:a:a2"),
				Rules: []datacenter.PolicyRule{
					{
						Label:             testutils.RandString(6, "hex"),
						Description:       testutils.RandString(6, "hex"),
						Protocol:          enum.PolicyRuleProtocolTcp,
						Action:            testutils.OneOf(enum.PolicyRuleActionDeny, enum.PolicyRuleActionDenyLog, enum.PolicyRuleActionPermit, enum.PolicyRuleActionPermitLog),
						SrcPort:           dctestobj.RandomPortRanges(3),
						DstPort:           dctestobj.RandomPortRanges(3),
						TcpStateQualifier: testutils.OneOf(nil, pointer.To(enum.TcpStateQualifierEstablished)),
					},
					{
						Label:       testutils.RandString(6, "hex"),
						Description: testutils.RandString(6, "hex"),
						Protocol:    enum.PolicyRuleProtocolUdp,
						Action:      testutils.OneOf(enum.PolicyRuleActionDeny, enum.PolicyRuleActionDenyLog, enum.PolicyRuleActionPermit, enum.PolicyRuleActionPermitLog),
						SrcPort:     dctestobj.RandomPortRanges(3),
						DstPort:     dctestobj.RandomPortRanges(3),
					},
				},
				Tags: testutils.RandomStrings(3, 6, 6, "hex"),
			},
		},
	}

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			rzMap := make(map[string]string)
			rzMutex := new(sync.Mutex)
			vnMap := make(map[string]string)
			vnMutex := new(sync.Mutex)
			var applicationPointID func(ctx context.Context, bp *apstra.TwoStageL3ClosClient, s string) string
			applicationPointID = func(ctx context.Context, bp *apstra.TwoStageL3ClosClient, s string) string {
				parts := strings.Split(s, ":")
				switch parts[0] {
				case "rz": // format is "rz:name"
					require.Len(t, parts, 2)
					rzMutex.Lock()
					defer rzMutex.Unlock()
					if id, ok := rzMap[s]; ok {
						return id
					}
					id, err := bp.CreateSecurityZone(ctx, datacenter.SecurityZone{
						Label:       parts[1],
						Description: nil,
						Type:        enum.SecurityZoneTypeEVPN,
						VRFName:     parts[1],
					})
					require.NoError(t, err)
					rzMap[s] = id
					return rzMap[s]
				case "vn": // format is "vn:rz-name:name"
					require.Len(t, parts, 3)
					vnMutex.Lock()
					defer vnMutex.Unlock()
					if id, ok := vnMap[s]; ok {
						return id
					}
					rzID := applicationPointID(ctx, bp, "rz:"+parts[1])
					id, err := bp.CreateVirtualNetwork(ctx, datacenter.VirtualNetwork{
						IPv4Enabled:    true,
						Label:          testutils.RandString(6, "hex"),
						SecurityZoneID: rzID,
						Type:           enum.VnTypeVxlan,
					})
					require.NoError(t, err)
					vnMap[s] = id
					return vnMap[s]
				default:
					t.Fatalf("unhandled application point string %q", s)
				}
				return ""
			}

			bp := dctestobj.TestBlueprintA(t, ctx, client.Client)

			for tName, tCase := range testCases {
				t.Run(tName, func(t *testing.T) {
					t.Parallel()
					ctx := testutils.ContextWithTestID(ctx, t)

					if tCase.versionConstraint != nil && !tCase.versionConstraint.Check(client.APIVersion()) {
						t.Skipf("skipping %q due to version constraints: %q. API version: %q",
							tName, tCase.versionConstraint, client.Client.ApiVersion())
					}

					// copy because we modify these values below
					create := tCase.create
					update := pointer.ToCopyOfValue(tCase.update)

					// create application point objects as necessary
					if create.SrcApplicationPoint != nil {
						create.SrcApplicationPoint = pointer.To(applicationPointID(ctx, bp, *create.SrcApplicationPoint))
					}
					if create.DstApplicationPoint != nil {
						create.DstApplicationPoint = pointer.To(applicationPointID(ctx, bp, *create.DstApplicationPoint))
					}
					if update.SrcApplicationPoint != nil {
						update.SrcApplicationPoint = pointer.To(applicationPointID(ctx, bp, *update.SrcApplicationPoint))
					}
					if update.DstApplicationPoint != nil {
						update.DstApplicationPoint = pointer.To(applicationPointID(ctx, bp, *update.DstApplicationPoint))
					}

					var id string
					var err error
					var obj datacenter.Policy

					// create the object
					id, err = bp.CreatePolicy(ctx, create)
					require.NoError(t, err)

					// retrieve the object by ID and validate
					obj, err = bp.GetPolicy(ctx, id)
					require.NoError(t, err)
					require.NotNil(t, obj.ID())
					require.Equal(t, id, *obj.ID())
					comparedatacenter.Policy(t, create, obj)

					// retrieve the object by label and validate
					obj, err = bp.GetPolicyByLabel(ctx, create.Label)
					require.NoError(t, err)
					require.NotNil(t, obj.ID())
					require.Equal(t, id, *obj.ID())
					comparedatacenter.Policy(t, create, obj)

					// retrieve the list of objects (ours must be in there) and validate
					objs, err := bp.GetPolicies(ctx)
					require.NoError(t, err)
					objPtr := slice.MustFindByID(objs, id)
					require.NotNil(t, objPtr)
					require.NotNil(t, obj.ID())
					require.Equal(t, id, *obj.ID())
					comparedatacenter.Policy(t, create, obj)

					// retrieve the list of IDs (ours must be in there) and validate
					ids, err := bp.ListPolicies(ctx)
					require.NoError(t, err)
					require.NotNil(t, ids)
					require.Contains(t, ids, id)

					if update != nil {
						// update the object
						require.NoError(t, update.SetID(id))
						require.NotNil(t, update.ID())
						require.Error(t, update.SetID(id))
						require.Equal(t, id, *update.ID())
						err = bp.UpdatePolicy(ctx, *update)
						require.NoError(t, err)

						// retrieve the object by ID and validate
						obj, err = bp.GetPolicy(ctx, id)
						require.NoError(t, err)
						comparedatacenter.Policy(t, *update, obj)
						require.NotNil(t, obj.ID())
						require.Equal(t, id, *obj.ID())

						// restore the object to the original state
						require.NoError(t, create.SetID(id))
						require.NotNil(t, create.ID())
						require.Error(t, create.SetID(id))
						require.Equal(t, id, *create.ID())
						err = bp.UpdatePolicy(ctx, create)
						require.NoError(t, err)

						// retrieve the object by ID and validate
						obj, err = bp.GetPolicy(ctx, id)
						require.NoError(t, err)
						comparedatacenter.Policy(t, create, obj)
						require.NotNil(t, obj.ID())
						require.Equal(t, id, *obj.ID())
					}

					// delete the object
					err = bp.DeletePolicy(ctx, id)
					require.NoError(t, err)

					// below this point we're expecting to *not* find the object
					var ace apstra.ClientErr

					// get the object by ID
					_, err = bp.GetPolicy(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// get the object by label
					_, err = bp.GetPolicyByLabel(ctx, create.Label)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// retrieve the list of objects (ours must *not* be in there)
					objs, err = bp.GetPolicies(ctx)
					require.NoError(t, err)
					objPtr = slice.MustFindByID(objs, id)
					require.Nil(t, objPtr)

					// retrieve the list of IDs (ours must *not* be in there)
					ids, err = bp.ListPolicies(ctx)
					require.NoError(t, err)
					require.NotNil(t, ids)
					require.NotContains(t, ids, id)

					// update the object
					require.NotNil(t, create.ID())
					require.Equal(t, id, *create.ID())
					err = bp.UpdatePolicy(ctx, create)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())

					// delete the object
					err = bp.DeletePolicy(ctx, id)
					require.Error(t, err)
					require.ErrorAs(t, err, &ace)
					require.Equal(t, apstra.ErrNotfound, ace.Type())
				})
			}
		})
	}
}

func TestPolicyAddDeleteRule(t *testing.T) {
	ctx := testutils.ContextWithTestID(t.Context(), t)
	clients := testclient.GetTestClients(t, ctx)

	//vnCount := 5

	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.ContextWithTestID(ctx, t)

			bp := dctestobj.TestBlueprintA(t, ctx, client.Client)

			pid, err := bp.CreatePolicy(ctx, datacenter.Policy{
				Label: testutils.RandString(6, "hex"),
				Rules: []datacenter.PolicyRule{
					{
						Label:    testutils.RandString(6, "hex"),
						Protocol: enum.PolicyRuleProtocolIp,
						Action:   enum.PolicyRuleActionDeny,
					},
					{
						Label:    testutils.RandString(6, "hex"),
						Protocol: enum.PolicyRuleProtocolIp,
						Action:   enum.PolicyRuleActionDeny,
					},
					{
						Label:    testutils.RandString(6, "hex"),
						Protocol: enum.PolicyRuleProtocolIp,
						Action:   enum.PolicyRuleActionDeny,
					},
				},
			})
			require.NoError(t, err)

			t.Run("add_rule_at_index_0", func(t *testing.T) {
				label := testutils.RandString(6, "hex")
				prid, err := bp.AddPolicyRule(ctx, datacenter.PolicyRule{
					Label:    label,
					Protocol: enum.PolicyRuleProtocolIp,
					Action:   enum.PolicyRuleActionDeny,
				}, 0, pid)
				require.NoError(t, err)

				obj, err := bp.GetPolicy(ctx, pid)
				require.NoError(t, err)
				require.Equal(t, label, obj.Rules[0].Label)
				require.Equal(t, 4, len(obj.Rules))

				err = bp.DeletePolicyRuleByID(ctx, pid, prid)
				require.NoError(t, err)

				obj, err = bp.GetPolicy(ctx, pid)
				require.NoError(t, err)
				require.Equal(t, 3, len(obj.Rules))
				for _, rule := range obj.Rules {
					require.NotEqual(t, label, rule.Label)
				}
			})

			t.Run("add_rule_at_index_2", func(t *testing.T) {
				label := testutils.RandString(6, "hex")
				prid, err := bp.AddPolicyRule(ctx, datacenter.PolicyRule{
					Label:    label,
					Protocol: enum.PolicyRuleProtocolIp,
					Action:   enum.PolicyRuleActionDeny,
				}, 2, pid)
				require.NoError(t, err)

				obj, err := bp.GetPolicy(ctx, pid)
				require.NoError(t, err)
				require.Equal(t, label, obj.Rules[2].Label)
				require.Equal(t, 4, len(obj.Rules))

				err = bp.DeletePolicyRuleByID(ctx, pid, prid)
				require.NoError(t, err)

				obj, err = bp.GetPolicy(ctx, pid)
				require.NoError(t, err)
				require.Equal(t, 3, len(obj.Rules))
				for _, rule := range obj.Rules {
					require.NotEqual(t, label, rule.Label)
				}
			})

			t.Run("add_rule_at_index_-1", func(t *testing.T) {
				label := testutils.RandString(6, "hex")
				prid, err := bp.AddPolicyRule(ctx, datacenter.PolicyRule{
					Label:    label,
					Protocol: enum.PolicyRuleProtocolIp,
					Action:   enum.PolicyRuleActionDeny,
				}, -1, pid)
				require.NoError(t, err)

				obj, err := bp.GetPolicy(ctx, pid)
				require.NoError(t, err)
				require.Equal(t, label, obj.Rules[3].Label)
				require.Equal(t, 4, len(obj.Rules))

				err = bp.DeletePolicyRuleByID(ctx, pid, prid)
				require.NoError(t, err)

				obj, err = bp.GetPolicy(ctx, pid)
				require.NoError(t, err)
				require.Equal(t, 3, len(obj.Rules))
				for _, rule := range obj.Rules {
					require.NotEqual(t, label, rule.Label)
				}
			})
		})

	}
}
