// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func randBool() bool {
	if rand.Int63()%2 == 0 {
		return true
	}
	return false
}

func randString(n int, style string) string {
	b64Letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")
	hexLetters := []rune("0123456789abcdef")
	var letters []rune
	b := make([]rune, n)
	switch style {
	case "hex":
		letters = hexLetters
	case "b64":
		letters = b64Letters
	}

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randStrings(count, strLen int) []string {
	result := make([]string, count)
	for i := range result {
		result[i] = randString(strLen, "hex")
	}
	return result
}

func randomIpv4() net.IP {
	return []byte{
		byte(rand.Intn(222) + 1),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
		byte(rand.Intn(256)),
	}
}

func randomIpv6() net.IP {
	return []byte{
		0x20, 0x01,
		0x0d, 0xb8,
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
		byte(rand.Intn(256)), byte(rand.Intn(256)),
	}
}

func netIpToNetIpAddr(t *testing.T, ip net.IP) netip.Addr {
	t.Helper()
	result, err := netip.ParseAddr(ip.String())
	require.NoError(t, err)
	return result
}

// randomHardwareAddr returns a net.HardwareAddr. The set and unset arguments
// allow the caller to specify certain bits which must be set or must be unset
// in the result.
// For example, to get a random mac with only the LAA bit set, you'd invoke the
// function with arguments indicating that LAA must be set and all other bits
// in the first byte must be unset:
//
//	set:   []byte{2},
//	unset: []byte{253},
func randomHardwareAddr(set []byte, unset []byte) net.HardwareAddr {
	result := net.HardwareAddr{
		byte(rand.Intn(math.MaxUint8 + 1)),
		byte(rand.Intn(math.MaxUint8 + 1)),
		byte(rand.Intn(math.MaxUint8 + 1)),
		byte(rand.Intn(math.MaxUint8 + 1)),
		byte(rand.Intn(math.MaxUint8 + 1)),
		byte(rand.Intn(math.MaxUint8 + 1)),
	}

	for i := range min(len(set), len(result)) {
		result[i] = result[i] | set[i]
	}

	for i := range min(len(unset), len(result)) {
		result[i] = result[i] & ^unset[i]
	}

	return result
}

func TestRandomHardwareAddr(t *testing.T) {
	type testCase struct {
		set   []byte
		unset []byte
	}

	testCases := map[string]testCase{
		"laa": {
			set: []byte{2},
		},
		"group": {
			set: []byte{1},
		},
		"laa_and_not_group": {
			set:   []byte{2},
			unset: []byte{1},
		},
		"group_and_not_laa": {
			set:   []byte{1},
			unset: []byte{2},
		},
		"laa_and_group": {
			set: []byte{3},
		},
		"last_byte_128": {
			set:   []byte{0, 0, 0, 0, 0, 128},
			unset: []byte{0, 0, 0, 0, 0, 127},
		},
		"last_byte_high": {
			set: []byte{0, 0, 0, 0, 0, 128},
		},
		"last_byte_low": {
			unset: []byte{0, 0, 0, 0, 0, 128},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := randomHardwareAddr(tCase.set, tCase.unset)

			for i, setByte := range tCase.set {
				require.Equal(t, setByte, result[i]&setByte)
			}

			for i, unsetByte := range tCase.unset {
				require.Equal(t, ^unsetByte, result[i]|^unsetByte)
			}
		})
	}
}

// randomIntsN fills the supplied slice with non-negative pseudo-random values in the half-open interval [0,n)
func randomIntsN(s []int, n int) {
	l := len(s)
	m := make(map[int]struct{}, l)
	for len(m) < l {
		m[rand.Intn(n)] = struct{}{}
	}

	i := 0
	for k := range m {
		s[i] = k
		i++
	}
}

func randomSlash31(t testing.TB) net.IPNet {
	t.Helper()

	ip := randomIpv4()
	_, ipNet, err := net.ParseCIDR(ip.String() + "/31")
	require.NoError(t, err)
	return *ipNet
}

func randomSlash127(t testing.TB) net.IPNet {
	t.Helper()

	ip := randomIpv6()
	_, ipNet, err := net.ParseCIDR(ip.String() + "/127")
	require.NoError(t, err)
	return *ipNet
}

func TestOurIpForPeer(t *testing.T) {
	test := net.ParseIP("127.0.0.1")
	expected := net.ParseIP("127.0.0.1")
	result, err := ourIpForPeer(test)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case test.String() == "<nil>":
		t.Fatal("test is '<nil>'")
	case expected.String() == "<nil>":
		t.Fatal("expected is '<nil>'")
	case result.String() == "<nil>":
		t.Fatal("result is '<nil>'")
	}
	if expected.String() != result.String() {
		t.Fatalf("expected %s, got %s", expected.String(), result.String())
	}
}

// Deprecated: use testutils.GetRandInts()
func getRandInts(min, max, count int) ([]int, error) {
	if max-min+1 < count {
		return nil, fmt.Errorf("cannot generate %d random numbers between %d and %d inclusive", count, min, max)
	}
	intMap := make(map[int]struct{}, count)
	for len(intMap) < count {
		intMap[rand.Intn(max+1-min)+min] = struct{}{}
	}
	result := make([]int, count)
	i := 0
	for k := range intMap {
		result[i] = k
		i++
	}
	return result, nil
}

// keyLogWriterFromEnv takes an environment variable which might name a logfile for
// exporting TLS session keys. If so, it returns an io.Writer to be used for
// that purpose, and the name of the logfile file.
func keyLogWriterFromEnv(keyLogEnv string) (*os.File, error) {
	fileName, foundKeyLogFile := os.LookupEnv(keyLogEnv)
	if !foundKeyLogFile {
		return nil, nil
	}

	// expand ~ style home directory
	if strings.HasPrefix(fileName, "~/") {
		dirname, _ := os.UserHomeDir()
		fileName = filepath.Join(dirname, fileName[2:])
	}

	err := os.MkdirAll(filepath.Dir(fileName), os.FileMode(0o600))
	if err != nil {
		return nil, err
	}
	return os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
}

// First tick should come immediately, certainly
// before half of the interval has expired.
func TestImmediateTickerFirstTick(t *testing.T) {
	interval := time.Second
	threshold := interval / 2

	start := time.Now()
	ticker := immediateTicker(time.Second)
	defer ticker.Stop()
	firstTick := <-ticker.C

	elapsed := firstTick.Sub(start)
	if elapsed > threshold {
		t.Fatalf("first tick after %q exceeds threshold %q", elapsed, threshold)
	}
	log.Printf("first tick after %q within threshold %q", elapsed, threshold)
	log.Printf("start %s first tick %s", start, firstTick)
}

// Second tick should come between .5 and 1.5 intervals
func TestImmediateTickerSecondTick(t *testing.T) {
	interval := time.Second
	threshold1 := interval / 2
	threshold2 := interval + interval/2

	start := time.Now()
	ticker := immediateTicker(time.Second)
	defer ticker.Stop()
	firstTick := <-ticker.C
	secondTick := <-ticker.C

	elapsed := secondTick.Sub(start)
	if elapsed < threshold1 {
		t.Fatalf("second tick after only %q doesn't meet threshold %q", elapsed, threshold1)
	}
	if elapsed > threshold2 {
		t.Fatalf("second tick after %q exceeds threshold %q", elapsed, threshold2)
	}
	log.Printf("second tick after %q within expected zone %q - %q", elapsed, threshold1, threshold2)
	log.Printf("start %s first tick %s second tick %s", start, firstTick, secondTick)
}

// Deprecated: Use testutils.TestBlueprintA
func testBlueprintA(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L3_Collapsed_ESI",
	})
	require.NoError(t, err)

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	return bpClient
}

// Deprecated: Use testutils.TestBlueprintB
func testBlueprintB(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual",
	})
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	return bpClient
}

// Deprecated: Use testutils.TestBlueprintC
func testBlueprintC(ctx context.Context, t testing.TB, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual_EVPN",
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	return bpClient
}

// Deprecated: Use testutils.TestBlueprintD
func testBlueprintD(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual_ESI_2x_Links",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.DeleteBlueprint(ctx, bpId))
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	query := new(PathQuery).
		SetBlueprintId(bpId).
		SetBlueprintType(BlueprintTypeStaging).
		SetClient(client).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{"system_type", QEStringVal("switch")},
			{"role", QEStringVal("leaf")},
			{"name", QEStringVal("n_leaf")},
		})
	var response struct {
		Items []struct {
			Leaf struct {
				ID string `json:"id"`
			} `json:"n_leaf"`
		} `json:"items"`
	}
	require.NoError(t, query.Do(ctx, &response))

	assignments := make(SystemIdToInterfaceMapAssignment)
	for _, item := range response.Items {
		assignments[item.Leaf.ID] = "Juniper_vQFX__AOS-7x10-Leaf"
	}

	require.NoError(t, bpClient.SetInterfaceMapAssignments(ctx, assignments))

	return bpClient
}

// Deprecated: Use testutils.TestBlueprintE
func testBlueprintE(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L2_ESI_Access",
	})
	if err != nil {
		t.Fatal(err)
	}

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	leafQuery := new(PathQuery).
		SetBlueprintId(bpId).
		SetBlueprintType(BlueprintTypeStaging).
		SetClient(client).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{"system_type", QEStringVal("switch")},
			{"role", QEStringVal("leaf")},
			{"name", QEStringVal("n_leaf")},
		})
	var leafResponse struct {
		Items []struct {
			Leaf struct {
				ID string `json:"id"`
			} `json:"n_leaf"`
		} `json:"items"`
	}
	err = leafQuery.Do(ctx, &leafResponse)
	require.NoError(t, err)

	leafAssignements := make(SystemIdToInterfaceMapAssignment)
	for _, item := range leafResponse.Items {
		leafAssignements[item.Leaf.ID] = "Juniper_vQFX__AOS-7x10-Leaf"
	}
	err = bpClient.SetInterfaceMapAssignments(ctx, leafAssignements)
	require.NoError(t, err)

	accessQuery := new(PathQuery).
		SetBlueprintId(bpId).
		SetBlueprintType(BlueprintTypeStaging).
		SetClient(client).
		Node([]QEEAttribute{
			NodeTypeSystem.QEEAttribute(),
			{"system_type", QEStringVal("switch")},
			{"role", QEStringVal("access")},
			{"name", QEStringVal("n_access")},
		})
	var accessResponse struct {
		Items []struct {
			Leaf struct {
				ID string `json:"id"`
			} `json:"n_access"`
		} `json:"items"`
	}
	err = accessQuery.Do(ctx, &accessResponse)
	require.NoError(t, err)

	accessAssignements := make(SystemIdToInterfaceMapAssignment)
	for _, item := range accessResponse.Items {
		accessAssignements[item.Leaf.ID] = "Juniper_vQFX__AOS-8x10-1"
	}
	err = bpClient.SetInterfaceMapAssignments(ctx, accessAssignements)
	require.NoError(t, err)

	return bpClient
}

// testBlueprintH creates a test blueprint using client and returns a *TwoStageL3ClosClient.
// The blueprint will use a dual-stack fabric and have ipv6 enabled.
// Deprecated: Use testutils.TestBlueprintH
func testBlueprintH(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpRequest := CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual_EVPN",
		FabricSettings: &FabricSettings{
			SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
			SpineLeafLinks:       toPtr(AddressingSchemeIp46),
		},
	}

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &bpRequest)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.DeleteBlueprint(ctx, bpId))
	})

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(err)
	}

	// set fabric addressing to enable IPv6
	if compatibility.EqApstra420.Check(client.apiVersion) {
		// todo - this is temporary
		require.NoError(t, bpClient.SetFabricAddressingPolicy(ctx, &TwoStageL3ClosFabricAddressingPolicy{Ipv6Enabled: toPtr(true)}))
	} else {
		require.NoError(t, client.talkToApstra(ctx, &talkToApstraIn{
			method: http.MethodPatch,
			urlStr: fmt.Sprintf("/api/blueprints/%s/fabric-settings", bpId),
			apiInput: struct {
				Ipv6Enabled bool `json:"ipv6_enabled"`
			}{
				Ipv6Enabled: true,
			},
		}))
	}

	return bpClient
}

// testBlueprintI returns a collapsed fabric which has been committed and has no build errors
// Deprecated: Use testutils.TestBlueprintI
func testBlueprintI(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L3_Collapsed_ESI",
	})
	require.NoError(t, err)

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	// assign leaf interface maps
	leafIds, err := getSystemIdsByRole(ctx, bpClient, "leaf")
	require.NoError(t, err)
	mappings := make(SystemIdToInterfaceMapAssignment, len(leafIds))
	for _, leafId := range leafIds {
		mappings[leafId.String()] = "Juniper_vQFX__AOS-7x10-Leaf"
	}
	err = bpClient.SetInterfaceMapAssignments(ctx, mappings)
	require.NoError(t, err)

	// set leaf loopback pool
	err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
		ResourceGroup: ResourceGroup{
			Type: ResourceTypeIp4Pool,
			Name: ResourceGroupNameLeafIp4,
		},
		PoolIds: []ObjectId{"Private-10_0_0_0-8"},
	})
	require.NoError(t, err)

	// set leaf-leaf pool
	err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
		ResourceGroup: ResourceGroup{
			Type: ResourceTypeIp4Pool,
			Name: ResourceGroupNameLeafLeafIp4,
		},
		PoolIds: []ObjectId{"Private-10_0_0_0-8"},
	})
	require.NoError(t, err)

	// set leaf ASN pool
	err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
		ResourceGroup: ResourceGroup{
			Type: ResourceTypeAsnPool,
			Name: ResourceGroupNameLeafAsn,
		},
		PoolIds: []ObjectId{"Private-64512-65534"},
	})
	require.NoError(t, err)

	// set VN VNI pool
	err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
		ResourceGroup: ResourceGroup{
			Type: ResourceTypeVniPool,
			Name: ResourceGroupNameEvpnL3Vni,
		},
		PoolIds: []ObjectId{"Default-10000-20000"},
	})
	require.NoError(t, err)

	// set VN VNI pool
	err = bpClient.SetResourceAllocation(ctx, &ResourceGroupAllocation{
		ResourceGroup: ResourceGroup{
			Type: ResourceTypeVniPool,
			Name: ResourceGroupNameVxlanVnIds,
		},
		PoolIds: []ObjectId{"Default-10000-20000"},
	})
	require.NoError(t, err)

	// commit
	bpStatus, err := client.GetBlueprintStatus(ctx, bpClient.blueprintId)
	require.NoError(t, err)
	_, err = client.DeployBlueprint(ctx, &BlueprintDeployRequest{
		Id:          bpClient.blueprintId,
		Description: "initial commit in test: " + t.Name(),
		Version:     bpStatus.Version,
	})
	require.NoError(t, err)

	return bpClient
}

func TestItemInSlice(t *testing.T) {
	type testCase struct {
		item     any
		slice    []any
		expected bool
	}

	testCases := []testCase{
		{item: 1, slice: []any{1, 2, 3}, expected: true},
		{item: 1, slice: []any{1, 2, 3, 1}, expected: true},
		{item: 1, slice: []any{3, 2, 1}, expected: true},
		{item: 0, slice: []any{1, 2, 3}, expected: false},
		{item: 0, slice: []any{}, expected: false},
		{item: 1, slice: []any{}, expected: false},
		{item: "foo", slice: []any{"foo", "bar"}, expected: true},
		{item: "foo", slice: []any{"bar", "foo"}, expected: true},
		{item: "foo", slice: []any{"foo", "bar", "foo"}, expected: true},
		{item: "foo", slice: []any{"bar", "baz"}, expected: false},
		{item: "foo", slice: []any{""}, expected: false},
		{item: "foo", slice: []any{"", ""}, expected: false},
		{item: "foo", slice: []any{}, expected: false},
		{item: "", slice: []any{"bar", "foo"}, expected: false},
		{item: "", slice: []any{"bar", "", "foo"}, expected: true},
		{item: "", slice: []any{}, expected: false},
	}

	var result bool
	for i, tc := range testCases {
		switch tc.item.(type) {
		case int:
			item := tc.item.(int)
			slice := make([]int, len(tc.slice))
			for j := range tc.slice {
				slice[j] = tc.slice[j].(int)
			}
			result = itemInSlice(item, slice)
		case string:
			item := tc.item.(string)
			slice := make([]string, len(tc.slice))
			for j := range tc.slice {
				slice[j] = tc.slice[j].(string)
			}
			result = itemInSlice(item, slice)
		}
		if result != tc.expected {
			t.Fatalf("test case %d produced %t, expected %t", i, result, tc.expected)
		}
	}
}

// Deprecated: Use testutils.TestRackA
func testRackA(ctx context.Context, t *testing.T, client *Client) (ObjectId, func(context.Context) error) {
	deleteFunc := func(context.Context) error { return nil }
	request := RackTypeRequest{
		DisplayName:              randString(5, "hex"),
		FabricConnectivityDesign: enum.FabricConnectivityDesignL3Clos,
		LeafSwitches: []RackElementLeafSwitchRequest{
			{
				Label:             randString(5, "hex"),
				LinkPerSpineCount: 1,
				LinkPerSpineSpeed: "40G",
				LogicalDeviceId:   "AOS-48x10_6x40-1",
			},
		},
	}

	id, err := client.CreateRackType(ctx, &request)
	if err != nil {
		t.Fatal(err)
	}
	deleteFunc = func(context.Context) error {
		return client.DeleteRackType(ctx, id)
	}

	return id, deleteFunc
}

// Deprecated: Use testutils.TestTemplateA
func testTemplateA(ctx context.Context, t *testing.T, client *Client) (ObjectId, func(context.Context) error) {
	deleteFunc := func(context.Context) error { return nil }
	rackId, rackDeleteFunc := testRackA(ctx, t, client)
	deleteFunc = func(context.Context) error {
		return rackDeleteFunc(ctx)
	}

	request := CreateRackBasedTemplateRequest{
		DisplayName: randString(5, "hex"),
		Spine: &TemplateElementSpineRequest{
			Count:         1,
			LogicalDevice: "AOS-16x40-1",
		},
		RackInfos: map[ObjectId]TemplateRackBasedRackInfo{
			rackId: {Count: 1},
		},
		AntiAffinityPolicy: &AntiAffinityPolicy{
			Algorithm:                AlgorithmHeuristic,
			MaxLinksPerPort:          1,
			MaxLinksPerSlot:          1,
			MaxPerSystemLinksPerPort: 1,
			MaxPerSystemLinksPerSlot: 1,
			Mode:                     AntiAffinityModeDisabled,
		},
		AsnAllocationPolicy:  &AsnAllocationPolicy{SpineAsnScheme: AsnAllocationSchemeDistinct},
		VirtualNetworkPolicy: &VirtualNetworkPolicy{OverlayControlProtocol: OverlayControlProtocolEvpn},
	}

	id, err := client.CreateRackBasedTemplate(ctx, &request)
	if err != nil {
		t.Fatal(errors.Join(err, rackDeleteFunc(ctx)))
	}
	deleteFunc = func(context.Context) error {
		return errors.Join(client.DeleteTemplate(ctx, id), rackDeleteFunc(ctx))
	}

	return id, deleteFunc
}

// Deprecated: Use testutils.TestBlueprintF
func testBlueprintF(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func(context.Context) error) {
	deleteFunc := func(context.Context) error { return nil }
	templateId, templateDeleteFunc := testTemplateA(ctx, t, client)
	deleteFunc = func(context.Context) error {
		return templateDeleteFunc(ctx)
	}

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: templateId,
	})
	if err != nil {
		t.Fatal(errors.Join(err, templateDeleteFunc(ctx)))
	}
	deleteFunc = func(ctx context.Context) error {
		return errors.Join(templateDeleteFunc(ctx), client.DeleteBlueprint(ctx, bpId))
	}

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(errors.Join(err, deleteFunc(ctx)))
	}

	return bpClient, deleteFunc
}

// Deprecated: Use testutils.TestBlueprintG
func testBlueprintG(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	templateId := testTemplateB(ctx, t, client)

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  enum.RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: templateId,
		FabricSettings: &FabricSettings{
			SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
			SpineLeafLinks:       toPtr(AddressingSchemeIp46),
		},
	})
	require.NoError(t, err)

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	require.NoError(t, bpClient.SetFabricSettings(ctx, &FabricSettings{Ipv6Enabled: toPtr(true)}))

	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	return bpClient
}

// Deprecated: Use testutils.TestTemplateB
func testTemplateB(ctx context.Context, t *testing.T, client *Client) ObjectId {
	t.Helper()

	rbt, err := client.GetRackBasedTemplate(ctx, "L2_Virtual")
	require.NoError(t, err)

	rbt.Data.DisplayName = randString(5, "hex")
	for k, v := range rbt.Data.RackInfo {
		v.RackTypeData = nil
		rbt.Data.RackInfo[k] = v
	}

	id, err := client.CreateRackBasedTemplate(ctx, &CreateRackBasedTemplateRequest{
		DisplayName: rbt.Data.DisplayName,
		Spine: &TemplateElementSpineRequest{
			Count:                  rbt.Data.Spine.Count,
			LinkPerSuperspineSpeed: rbt.Data.Spine.LinkPerSuperspineSpeed,
			LogicalDevice:          "AOS-7x10-Spine",
			LinkPerSuperspineCount: rbt.Data.Spine.LinkPerSuperspineCount,
		},
		RackInfos:            rbt.Data.RackInfo,
		DhcpServiceIntent:    &rbt.Data.DhcpServiceIntent,
		AntiAffinityPolicy:   rbt.Data.AntiAffinityPolicy,
		AsnAllocationPolicy:  &rbt.Data.AsnAllocationPolicy,
		VirtualNetworkPolicy: &rbt.Data.VirtualNetworkPolicy,
	})
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, client.DeleteTemplate(ctx, id)) })

	return id
}

// Deprecated: Use testutils.TestSecurityZoneA
func testSecurityZone(t testing.TB, ctx context.Context, bp *TwoStageL3ClosClient) ObjectId {
	t.Helper()

	rs := randString(6, "hex")

	id, err := bp.CreateSecurityZone(ctx, &SecurityZoneData{
		Label:            rs,
		SzType:           SecurityZoneTypeEVPN,
		VrfName:          rs,
		RoutingPolicyId:  "",
		RouteTarget:      nil,
		RtPolicy:         nil,
		VlanId:           nil,
		VniId:            nil,
		JunosEvpnIrbMode: nil,
	})
	require.NoError(t, err)

	return id
}

// Deprecated: Use testutils.TestVirtualNetworkA
func testVirtualNetwork(t testing.TB, ctx context.Context, bp *TwoStageL3ClosClient, szId ObjectId) ObjectId {
	t.Helper()

	var vnBindings []VnBinding
	nodeMap, err := bp.GetAllSystemNodeInfos(ctx)
	require.NoError(t, err)

	for _, node := range nodeMap {
		if node.Role == SystemRoleLeaf {
			vnBindings = append(vnBindings, VnBinding{SystemId: node.Id})
		}
	}

	id, err := bp.CreateVirtualNetwork(ctx, &VirtualNetworkData{
		Ipv4Enabled:               true,
		Label:                     randString(6, "hex"),
		SecurityZoneId:            szId,
		VirtualGatewayIpv4Enabled: true,
		VnBindings:                vnBindings,
		VnType:                    enum.VnTypeVxlan,
	})
	require.NoError(t, err)

	return id
}

// testFFBlueprintA creates an empty Freeform blueprint
func testFFBlueprintA(ctx context.Context, t testing.TB, client *Client) *FreeformClient {
	t.Helper()

	id, err := client.CreateFreeformBlueprint(ctx, randString(6, "hex"))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, client.DeleteBlueprint(ctx, id))
	})

	c, err := client.NewFreeformClient(ctx, id)
	require.NoError(t, err)

	return c
}

// testFFBlueprintB creates a freeform blueprint with predefined internal and external generic systems.
// The returned []ObjectIds represent the requested internal and external generic systems.
func testFFBlueprintB(ctx context.Context, t testing.TB, client *Client, intSystemCount, extSystemCount int) (*FreeformClient, []ObjectId, []ObjectId) {
	t.Helper()

	id, err := client.CreateFreeformBlueprint(ctx, randString(6, "hex"))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, client.DeleteBlueprint(ctx, id))
	})

	c, err := client.NewFreeformClient(ctx, id)
	require.NoError(t, err)

	dpId, err := c.ImportDeviceProfile(ctx, "Juniper_EX4400-48T")
	require.NoError(t, err)

	intSystemIds := make([]ObjectId, intSystemCount)
	for i := range intSystemIds {
		intSystemIds[i], err = c.CreateSystem(ctx, &FreeformSystemData{
			Type:            SystemTypeInternal,
			Label:           randString(6, "hex"),
			DeviceProfileId: &dpId,
		})
		require.NoError(t, err)
	}

	extSystemIds := make([]ObjectId, extSystemCount)
	for i := range extSystemIds {
		extSystemIds[i], err = c.CreateSystem(ctx, &FreeformSystemData{
			Type:  SystemTypeExternal,
			Label: randString(6, "hex"),
		})
		require.NoError(t, err)
	}

	return c, intSystemIds, extSystemIds
}

func testAsnPool(ctx context.Context, t testing.TB, client *Client) ObjectId {
	t.Helper()

	asnBeginEnds, err := getRandInts(1, 100000000, (rand.Intn(5)+2)*2)
	require.NoError(t, err)
	sort.Ints(asnBeginEnds) // sort so that the ASN ranges will be ([0]...[1], [2]...[3], etc.)

	asnRanges := make([]IntfIntRange, len(asnBeginEnds)/2)
	for i := range asnRanges {
		asnRanges[i] = IntRangeRequest{
			uint32(asnBeginEnds[2*i]),
			uint32(asnBeginEnds[(2*i)+1]),
		}
	}

	id, err := client.createAsnPool(ctx, &AsnPoolRequest{
		DisplayName: "test-" + randString(6, "hex"),
		Ranges:      asnRanges,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.DeleteAsnPool(ctx, id)) })

	return id
}

func testIntPool(ctx context.Context, t testing.TB, client *Client) ObjectId {
	t.Helper()

	intBeginEnds, err := getRandInts(1, 100000000, (rand.Intn(5)+2)*2)
	require.NoError(t, err)
	sort.Ints(intBeginEnds) // sort so that the Int ranges will be ([0]...[1], [2]...[3], etc.)

	intRanges := make([]IntfIntRange, len(intBeginEnds)/2)
	for i := range intRanges {
		intRanges[i] = IntRangeRequest{
			uint32(intBeginEnds[2*i]),
			uint32(intBeginEnds[(2*i)+1]),
		}
	}
	id, err := client.CreateIntegerPool(ctx, &IntPoolRequest{
		DisplayName: "test-" + randString(6, "hex"),
		Ranges:      intRanges,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.DeleteIntegerPool(ctx, id)) })

	return id
}

func testVniPool(ctx context.Context, t testing.TB, client *Client) ObjectId {
	t.Helper()

	vniBeginEnds, err := getRandInts(5000, 5999, (rand.Intn(5)+2)*2)
	require.NoError(t, err)
	sort.Ints(vniBeginEnds) // sort so that the Int ranges will be ([0]...[1], [2]...[3], etc.)

	vniRanges := make([]IntfIntRange, len(vniBeginEnds)/2)
	for i := range vniRanges {
		vniRanges[i] = IntRangeRequest{
			uint32(vniBeginEnds[2*i]),
			uint32(vniBeginEnds[(2*i)+1]),
		}
	}

	id, err := client.CreateVniPool(ctx, &VniPoolRequest{
		DisplayName: "test-" + randString(6, "hex"),
		Ranges:      vniRanges,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.DeleteVniPool(ctx, id)) })

	return id
}

func testIpv4Pool(ctx context.Context, t testing.TB, client *Client) ObjectId {
	t.Helper()

	subnets := make([]NewIpSubnet, rand.Intn(5)+2)
	for i := range subnets {
		randNet := randomSlash31(t)
		subnets[i] = NewIpSubnet{Network: randNet.String()}
	}

	id, err := client.CreateIp4Pool(ctx, &NewIpPoolRequest{
		DisplayName: "test-" + randString(6, "hex"),
		Subnets:     subnets,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.DeleteIp4Pool(ctx, id)) })

	return id
}

func testIpv6Pool(ctx context.Context, t testing.TB, client *Client) ObjectId {
	t.Helper()

	subnets := make([]NewIpSubnet, rand.Intn(5)+2)
	for i := range subnets {
		randNet := randomSlash127(t)
		subnets[i] = NewIpSubnet{Network: randNet.String()}
	}

	id, err := client.CreateIp6Pool(ctx, &NewIpPoolRequest{
		DisplayName: "test-" + randString(6, "hex"),
		Subnets:     subnets,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.DeleteIp6Pool(ctx, id)) })

	return id
}

func testResourceGroup(ctx context.Context, t testing.TB, client *FreeformClient) (groupId ObjectId) {
	id, err := client.CreateRaGroup(ctx, &FreeformRaGroupData{Label: randString(6, "hex")})
	require.NoError(t, err)

	return id
}

func testResourceGroupAsn(ctx context.Context, t testing.TB, client *FreeformClient) (id ObjectId) {
	id, err := client.CreateAllocGroup(ctx, &FreeformAllocGroupData{
		Name:    randString(6, "hex"),
		Type:    enum.ResourcePoolTypeAsn,
		PoolIds: []ObjectId{testAsnPool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}

func testResourceGroupInt(ctx context.Context, t testing.TB, client *FreeformClient) (id ObjectId) {
	id, err := client.CreateAllocGroup(ctx, &FreeformAllocGroupData{
		Name:    randString(6, "hex"),
		Type:    enum.ResourcePoolTypeInt,
		PoolIds: []ObjectId{testIntPool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}

func testResourceGroupIpv4(ctx context.Context, t testing.TB, client *FreeformClient) (id ObjectId) {
	id, err := client.CreateAllocGroup(ctx, &FreeformAllocGroupData{
		Name:    randString(6, "hex"),
		Type:    enum.ResourcePoolTypeIpv4,
		PoolIds: []ObjectId{testIpv4Pool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}

func testResourceGroupIpv6(ctx context.Context, t testing.TB, client *FreeformClient) (id ObjectId) {
	id, err := client.CreateAllocGroup(ctx, &FreeformAllocGroupData{
		Name:    randString(6, "hex"),
		Type:    enum.ResourcePoolTypeIpv6,
		PoolIds: []ObjectId{testIpv6Pool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}

func newUUID(t testing.TB) uuid.UUID {
	t.Helper()

	result, err := uuid.NewRandom()
	require.NoError(t, err)
	return result
}

// wrapCtxWithTestId produces contexts with the following values:
// - Test-UUID: a uuid.UUID representing this test and all sub-tests.
// - Test-ID: a string of the form uuid/test/subtest/subsubtest...
// the Test-UUID is generated only if not found.
// HTTP transactions related to these tests can be picked out from wireshark
// using filters like:
// - http.request.line contains "843a754c-cc35-4383-807f-833ad991e554"
// - http.request.line contains "843a754c-cc35-4383-807f-833ad991e554/test/subtest"
func wrapCtxWithTestId(t testing.TB, ctx context.Context) context.Context {
	var UUID *uuid.UUID

	switch v := ctx.Value(CtxKeyTestUUID).(type) {
	case uuid.UUID:
		UUID = &v
	default:
		UUID = toPtr(newUUID(t))
		ctx = context.WithValue(ctx, CtxKeyTestUUID, *UUID)
		log.Println("Test UUID: ", UUID.String())
	}

	return context.WithValue(ctx, CtxKeyTestID, UUID.String()+"/"+t.Name())
}

func oneOf[A interface{}](i A, s ...A) A {
	return append([]A{i}, s...)[rand.Intn(len(s)+1)]
}
