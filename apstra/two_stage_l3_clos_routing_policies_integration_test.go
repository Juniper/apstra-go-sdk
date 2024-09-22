// Copyright (c) Juniper Networks, Inc., 2023-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/stretchr/testify/require"
)

func compareDcRoutingExportPolicies(t *testing.T, a, b DcRoutingExportPolicy) {
	t.Helper()

	if a.StaticRoutes != b.StaticRoutes {
		t.Fatal()
	}
	if a.Loopbacks != b.Loopbacks {
		t.Fatal()
	}
	if a.SpineSuperspineLinks != b.SpineSuperspineLinks {
		t.Fatal()
	}
	if a.L3EdgeServerLinks != b.L3EdgeServerLinks {
		t.Fatal()
	}
	if a.SpineLeafLinks != b.SpineLeafLinks {
		t.Fatal()
	}
	if a.L2EdgeSubnets != b.L2EdgeSubnets {
		t.Fatal()
	}
}

func comparePrefixSlices(t *testing.T, a, b []net.IPNet) {
	t.Helper()

	if len(a) != len(b) {
		t.Fatal()
	}
	for i := range a {
		if a[i].String() != b[i].String() {
			t.Fatal()
		}
	}
}

func comparePrefixFilters(t *testing.T, a, b PrefixFilter) {
	t.Helper()

	if a.Action != b.Action {
		t.Fatal()
	}
	if a.Prefix.String() != b.Prefix.String() {
		t.Fatal()
	}
	if (a.GeMask == nil) != (b.GeMask == nil) {
		t.Fatal() // where one is nil the other must be nil
	}
	if (a.LeMask == nil) != (b.LeMask == nil) {
		t.Fatal() // where one is nil the other must be nil
	}
	if a.GeMask != nil && b.GeMask != nil {
		if *a.GeMask != *b.GeMask {
			t.Fatal()
		}
	}
	if a.LeMask != nil && b.LeMask != nil {
		if *a.LeMask != *b.LeMask {
			t.Fatal()
		}
	}
}

func comparePrefixFilterSlices(t *testing.T, a, b []PrefixFilter) {
	t.Helper()

	if len(a) != len(b) {
		t.Fatal()
	}

	for i := range a {
		comparePrefixFilters(t, a[i], b[i])
	}
}

func compareDcRoutingPolicyData(t *testing.T, a, b *DcRoutingPolicyData) {
	t.Helper()

	require.Equalf(t, a.Label, b.Label, "Label: Expected %q, got %q", a.Label, b.Label)
	require.Equalf(t, a.Description, b.Description, "Description: Expected %q, got %q", a.Description, b.Description)
	require.Equalf(t, a.PolicyType, b.PolicyType, "PolicyType: Expected %q, got %q", a.PolicyType, b.PolicyType)
	require.Equalf(t, a.ImportPolicy, b.ImportPolicy, "ImportPolicy: Expected %q, got %q", a.ImportPolicy, b.ImportPolicy)

	compareDcRoutingExportPolicies(t, a.ExportPolicy, b.ExportPolicy)
	require.Equalf(t, a.ExpectDefaultIpv4Route, b.ExpectDefaultIpv4Route, "ExpectDefaultIpv4Route: Expected %t, got %t", a.ExpectDefaultIpv4Route, b.ExpectDefaultIpv4Route)
	require.Equalf(t, a.ExpectDefaultIpv6Route, b.ExpectDefaultIpv6Route, "ExpectDefaultIpv6Route: Expected %t, got %t", a.ExpectDefaultIpv6Route, b.ExpectDefaultIpv6Route)

	comparePrefixSlices(t, a.AggregatePrefixes, b.AggregatePrefixes)
	comparePrefixFilterSlices(t, a.ExtraImportRoutes, b.ExtraImportRoutes)
	comparePrefixFilterSlices(t, a.ExtraExportRoutes, b.ExtraExportRoutes)
}

func TestRoutingPolicies(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		clientName, client := clientName, client
		t.Run(fmt.Sprintf("%s_%s", client.client.ApiVersion(), clientName), func(t *testing.T) {
			t.Parallel()
			bpClient := testBlueprintA(ctx, t, client.client)

			log.Printf("testing GetDefaultRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			defaultPolicy, err := bpClient.GetDefaultRoutingPolicy(ctx)
			require.NoError(t, err)
			require.Equal(t, defaultPolicy.Data.PolicyType, DcRoutingPolicyTypeDefault)

			var aggregatePrefixes []net.IPNet
			for _, s := range []string{"1.0.0.0/8", "2.0.0.0/7"} {
				_, ipNet, err := net.ParseCIDR(s)
				require.NoError(t, err)
				aggregatePrefixes = append(aggregatePrefixes, *ipNet)
			}

			var f PrefixFilter
			var ipNet *net.IPNet

			eleven := 11
			twelve := 12
			thirteen := 13
			fourteen := 14

			var importFilters []PrefixFilter

			_, ipNet, err = net.ParseCIDR("100.0.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionPermit,
				Prefix: *ipNet,
				GeMask: &eleven,
				LeMask: &thirteen,
			}
			importFilters = append(importFilters, f)

			_, ipNet, err = net.ParseCIDR("100.64.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				GeMask: &twelve,
				LeMask: &fourteen,
			}
			importFilters = append(importFilters, f)

			_, ipNet, err = net.ParseCIDR("100.128.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				GeMask: &eleven,
			}
			importFilters = append(importFilters, f)

			_, ipNet, err = net.ParseCIDR("100.192.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				LeMask: &eleven,
			}
			importFilters = append(importFilters, f)

			var exportFilters []PrefixFilter
			_, ipNet, err = net.ParseCIDR("200.0.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionPermit,
				Prefix: *ipNet,
				GeMask: &eleven,
				LeMask: &thirteen,
			}
			exportFilters = append(exportFilters, f)
			_, ipNet, err = net.ParseCIDR("200.64.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				GeMask: &twelve,
				LeMask: &fourteen,
			}
			exportFilters = append(exportFilters, f)
			_, ipNet, err = net.ParseCIDR("200.128.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				GeMask: &twelve,
			}
			exportFilters = append(exportFilters, f)
			_, ipNet, err = net.ParseCIDR("200.192.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				LeMask: &fourteen,
			}
			exportFilters = append(exportFilters, f)

			randStr := randString(5, "hex")
			policyData := &DcRoutingPolicyData{
				Label:        "test-label-" + randStr,
				Description:  "test-description-" + randStr,
				PolicyType:   DcRoutingPolicyTypeUser,
				ImportPolicy: DcRoutingPolicyImportPolicyAll,
				ExportPolicy: DcRoutingExportPolicy{
					StaticRoutes:         false,
					Loopbacks:            false,
					SpineSuperspineLinks: false,
					L3EdgeServerLinks:    false,
					SpineLeafLinks:       false,
					L2EdgeSubnets:        false,
				},
				ExpectDefaultIpv4Route: false,
				ExpectDefaultIpv6Route: false,
				AggregatePrefixes:      aggregatePrefixes,
				ExtraImportRoutes:      importFilters,
				ExtraExportRoutes:      exportFilters,
			}

			log.Printf("testing CreateRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			policyId, err := bpClient.CreateRoutingPolicy(ctx, policyData)
			require.NoError(t, err)
			log.Printf("testing GetRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			policy, err := bpClient.GetRoutingPolicy(ctx, policyId)
			require.NoError(t, err)
			if policy.Id != policyId {
				t.Fatalf("policy IDs don't match %q vs. %q", policy.Id, policyId)
			}
			compareDcRoutingPolicyData(t, policyData, policy.Data)

			log.Printf("testing GetRoutingPolicyByName() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			policy, err = bpClient.GetRoutingPolicyByName(ctx, policy.Data.Label)
			require.NoError(t, err)
			if policy.Id != policyId {
				t.Fatalf("policy IDs don't match %q vs. %q", policy.Id, policyId)
			}
			compareDcRoutingPolicyData(t, policyData, policy.Data)

			log.Printf("testing GetAllRoutingPolicies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			policies, err := bpClient.GetAllRoutingPolicies(ctx)
			require.NoError(t, err)
			if len(policies) != 2 {
				t.Fatalf("expected 2 policies, got %d", len(policies))
			}

			if policies[0].Data.PolicyType != DcRoutingPolicyTypeDefault && policies[1].Data.PolicyType != DcRoutingPolicyTypeDefault {
				t.Fatalf("neither policy has type %q, got %q and %q",
					DcRoutingPolicyTypeDefault, policies[0].Data.PolicyType.String(), policies[1].Data.PolicyType.String())
			}

			if policies[0].Data.PolicyType != DcRoutingPolicyTypeUser && policies[1].Data.PolicyType != DcRoutingPolicyTypeUser {
				t.Fatalf("neither policy has type %q, got %q and %q",
					DcRoutingPolicyTypeUser, policies[0].Data.PolicyType.String(), policies[1].Data.PolicyType.String())
			}

			if policies[0].Id != defaultPolicy.Id && policies[1].Id != defaultPolicy.Id {
				t.Fatalf("neither policy has ID %q, got %q and %q",
					defaultPolicy.Id, policies[0].Id.String(), policies[1].Id.String())
			}

			if policies[0].Id != policy.Id && policies[1].Id != policy.Id {
				t.Fatalf("neither policy has ID %q, got %q and %q",
					policy.Id, policies[0].Id.String(), policies[1].Id.String())
			}

			_, ipNet, err = net.ParseCIDR("110.0.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionPermit,
				Prefix: *ipNet,
				GeMask: &eleven,
				LeMask: &thirteen,
			}

			importFilters = append(importFilters, f)
			_, ipNet, err = net.ParseCIDR("110.32.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				GeMask: &twelve,
				LeMask: &fourteen,
			}

			importFilters = append(importFilters, f)
			_, ipNet, err = net.ParseCIDR("210.0.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionPermit,
				Prefix: *ipNet,
				GeMask: &eleven,
				LeMask: &thirteen,
			}

			exportFilters = append(exportFilters, f)
			_, ipNet, err = net.ParseCIDR("210.32.0.0/10")
			require.NoError(t, err)
			f = PrefixFilter{
				Action: PrefixFilterActionDeny,
				Prefix: *ipNet,
				GeMask: &twelve,
				LeMask: &fourteen,
			}
			exportFilters = append(exportFilters, f)

			randStr = randString(5, "hex")
			policyData.Label = "test-label-" + randStr
			policyData.Description = "test-description-" + randStr
			policyData.ExpectDefaultIpv4Route = true
			policyData.ExpectDefaultIpv6Route = true
			policyData.ExtraImportRoutes = importFilters
			policyData.ExtraExportRoutes = exportFilters
			policyData.ImportPolicy = DcRoutingPolicyImportPolicyDefaultOnly
			policyData.ExportPolicy = DcRoutingExportPolicy{
				StaticRoutes:         true,
				Loopbacks:            true,
				SpineSuperspineLinks: true,
				L3EdgeServerLinks:    compatibility.RoutingPolicyExportHasL3EdgeLinks.Check(bpClient.client.apiVersion),
				SpineLeafLinks:       true,
				L2EdgeSubnets:        true,
			}

			log.Printf("testing UpdateRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			require.NoError(t, bpClient.UpdateRoutingPolicy(ctx, policy.Id, policyData))

			log.Printf("testing GetRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			updatedPolicy, err := bpClient.GetRoutingPolicy(ctx, policy.Id)
			require.NoError(t, err)
			compareDcRoutingPolicyData(t, policyData, updatedPolicy.Data)

			log.Printf("testing DeleteRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = bpClient.DeleteRoutingPolicy(ctx, policy.Id)
			require.NoError(t, err)
			log.Printf("testing GetAllRoutingPolicies() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			policies, err = bpClient.GetAllRoutingPolicies(ctx)
			require.NoError(t, err)
			if len(policies) != 1 {
				t.Fatalf("expected 1 policies, got %d", len(policies))
			}
			if policies[0].Id != defaultPolicy.Id {
				t.Fatalf("surviving policy ID %q does not match previously noted default policy ID %q", policies[0].Id, defaultPolicy.Id)
			}
		})
	}
}

func TestRoutingPolicy404(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient := testBlueprintA(ctx, t, client.client)

		log.Printf("testing GetDefaultRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		_, err := bpClient.GetRoutingPolicy(ctx, "bogus")
		if err == nil {
			t.Fatal("should have gotten an error")
		} else {
			var clientErr ClientErr
			if !errors.As(err, &clientErr) || clientErr.Type() != ErrNotfound {
				t.Fatal("error should have been something 404-ish")
			}
		}

		log.Printf("testing DeleteRoutingPolicy() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = bpClient.DeleteRoutingPolicy(ctx, "bogus")
		if err == nil {
			t.Fatal("should have gotten an error")
		} else {
			var clientErr ClientErr
			if !errors.As(err, &clientErr) || clientErr.Type() != ErrNotfound {
				t.Fatal("error should have been something 404-ish")
			}
		}
	}
}
