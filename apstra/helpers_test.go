package apstra

import (
	"context"
	crand "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

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

func randJwt() string {
	return randString(36, "b64") + "." +
		randString(178, "b64") + "." +
		randString(86, "b64")
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

func randomPrefix(t testing.TB, cidrBlock string, bits int) net.IPNet {
	t.Helper()

	ip, block, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		t.Fatalf("randomPrefix cannot parse cidrBlock - %s", err)
	}
	if block.IP.String() != ip.String() {
		t.Fatal("invocation of randomPrefix doesn't use a base block address")
	}

	mOnes, mBits := block.Mask.Size()
	if mOnes >= bits {
		t.Fatalf("cannot select a random /%d from within %s", bits, cidrBlock)
	}

	// generate a completely random address
	randomIP := make(net.IP, mBits/8)
	_, err = crand.Read(randomIP)
	if err != nil {
		t.Fatalf("rand read failed")
	}

	// mask off the "network" bits
	for i, b := range randomIP {
		mBitsThisByte := min(mOnes, 8)
		mOnes -= mBitsThisByte
		block.IP[i] = block.IP[i] | (b & byte(math.MaxUint8>>mBitsThisByte))
	}

	block.Mask = net.CIDRMask(bits, mBits)

	_, result, err := net.ParseCIDR(block.String())
	if err != nil {
		t.Fatal("failed to parse own CIDR block")
	}

	return *result
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

func testBlueprintA(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
		Label:      randString(5, "hex"),
		TemplateId: "L3_Collapsed_ESI",
	})
	require.NoError(t, err)

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	return bpClient
}

func testBlueprintB(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual",
	})
	require.NoError(t, err)

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	return bpClient
}

func testBlueprintC(ctx context.Context, t testing.TB, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual_EVPN",
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.DeleteBlueprint(ctx, bpId)) })

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	require.NoError(t, err)

	return bpClient
}

func testBlueprintD(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func(context.Context) error) {
	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual_ESI_2x_Links",
	})
	if err != nil {
		t.Fatal(err)
	}

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(err)
	}

	bpDeleteFunc := func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	}

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
	err = query.Do(ctx, &response)
	if err != nil {
		t.Fatal(errors.Join(err, bpDeleteFunc(ctx)))
	}

	assignments := make(SystemIdToInterfaceMapAssignment)
	for _, item := range response.Items {
		assignments[item.Leaf.ID] = "Juniper_vQFX__AOS-7x10-Leaf"
	}

	err = bpClient.SetInterfaceMapAssignments(ctx, assignments)
	if err != nil {
		t.Fatal(errors.Join(err, bpDeleteFunc(ctx)))
	}

	return bpClient, bpDeleteFunc
}

func testBlueprintE(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func(context.Context) error) {
	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
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

	bpDeleteFunc := func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	}

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
	if err != nil {
		t.Fatal(errors.Join(err, bpDeleteFunc(ctx)))
	}
	leafAssignements := make(SystemIdToInterfaceMapAssignment)
	for _, item := range leafResponse.Items {
		leafAssignements[item.Leaf.ID] = "Juniper_vQFX__AOS-7x10-Leaf"
	}
	err = bpClient.SetInterfaceMapAssignments(ctx, leafAssignements)
	if err != nil {
		t.Fatal(errors.Join(err, bpDeleteFunc(ctx)))
	}

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
	if err != nil {
		t.Fatal(errors.Join(err, bpDeleteFunc(ctx)))
	}
	accessAssignements := make(SystemIdToInterfaceMapAssignment)
	for _, item := range accessResponse.Items {
		accessAssignements[item.Leaf.ID] = "Juniper_vQFX__AOS-8x10-1"
	}
	err = bpClient.SetInterfaceMapAssignments(ctx, accessAssignements)
	if err != nil {
		t.Fatal(errors.Join(err, bpDeleteFunc(ctx)))
	}

	return bpClient, bpDeleteFunc
}

// testBlueprintH creates a test blueprint using client and returns a *TwoStageL3ClosClient and a cleanup function
// which deletes the test blueprint. If the client is newer that Apstra 4.1.0, the blueprint will use a dual-stack
// fabric and have ipv6 enabled.
func testBlueprintH(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func(context.Context) error) {
	bpRequest := CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual_EVPN",
	}

	if rackBasedTemplateFabricAddressingPolicyForbidden().Includes(client.apiVersion.String()) {
		bpRequest.FabricSettings = &FabricSettings{
			SpineSuperspineLinks: toPtr(AddressingSchemeIp46),
			SpineLeafLinks:       toPtr(AddressingSchemeIp46),
		}
	}

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &bpRequest)
	if err != nil {
		t.Fatal(err)
	}

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(err)
	}

	bpDeleteFunc := func(ctx context.Context) error {
		return client.DeleteBlueprint(ctx, bpId)
	}

	// set fabric addressing to enable IPv6
	if rackBasedTemplateFabricAddressingPolicyForbidden().Includes(client.apiVersion.String()) {
		if client.apiVersion.String() == "4.2.1" {
			// todo - this is temporary
			err = client.talkToApstra(ctx, &talkToApstraIn{
				method: http.MethodPatch,
				urlStr: fmt.Sprintf("/api/blueprints/%s/fabric-settings", bpId),
				apiInput: struct {
					Ipv6Enabled bool `json:"ipv6_enabled"`
				}{
					Ipv6Enabled: true,
				},
			})
			if err != nil {
				t.Fatal(err)
			}
		} else {
			err = bpClient.SetFabricAddressingPolicy(ctx, &TwoStageL3ClosFabricAddressingPolicy{Ipv6Enabled: toPtr(true)})
			if err != nil {
				defer bpDeleteFunc(ctx)
				t.Fatal(err)
			}
		}
	}

	return bpClient, bpDeleteFunc
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

func testRackA(ctx context.Context, t *testing.T, client *Client) (ObjectId, func(context.Context) error) {
	deleteFunc := func(context.Context) error { return nil }
	request := RackTypeRequest{
		DisplayName:              randString(5, "hex"),
		FabricConnectivityDesign: FabricConnectivityDesignL3Clos,
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
		FabricAddressingPolicy: &TemplateFabricAddressingPolicy410Only{
			SpineSuperspineLinks: AddressingSchemeIp4,
			SpineLeafLinks:       AddressingSchemeIp4,
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

func testBlueprintF(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func(context.Context) error) {
	deleteFunc := func(context.Context) error { return nil }
	templateId, templateDeleteFunc := testTemplateA(ctx, t, client)
	deleteFunc = func(context.Context) error {
		return templateDeleteFunc(ctx)
	}

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
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

func testBlueprintG(ctx context.Context, t *testing.T, client *Client) *TwoStageL3ClosClient {
	t.Helper()

	templateId := testTemplateB(ctx, t, client)

	bpId, err := client.CreateBlueprintFromTemplate(ctx, &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignTwoStageL3Clos,
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

// testWidgetsAB instantiates two predefined probes and creates widgets from them,
// returning the widget Object Id and the IbaWidgetData object used for creation
func testWidgetsAB(ctx context.Context, t *testing.T, bpClient *TwoStageL3ClosClient) (ObjectId, IbaWidgetData, ObjectId, IbaWidgetData) {
	probeAId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &IbaPredefinedProbeRequest{
		Name: "bgp_session",
		Data: []byte(`{
			"Label":     "BGP Session Flapping",
			"Duration":  300,
			"Threshold": 40
		}`),
	})
	if err != nil {
		t.Fatal(err)
	}

	probeBId, err := bpClient.InstantiateIbaPredefinedProbe(ctx, &IbaPredefinedProbeRequest{
		Name: "drain_node_traffic_anomaly",
		Data: []byte(`{
			"Label":     "Drain Traffic Anomaly",
			"Threshold": 100000
		}`),
	})
	if err != nil {
		t.Fatal(err)
	}

	widgetA := IbaWidgetData{
		Type:      IbaWidgetTypeStage,
		Label:     "BGP Session Flapping",
		ProbeId:   probeAId,
		StageName: "BGP Session",
	}
	widgetAId, err := bpClient.CreateIbaWidget(ctx, &widgetA)
	if err != nil {
		t.Fatal(err)
	}

	widgetB := IbaWidgetData{
		Type:      IbaWidgetTypeStage,
		Label:     "Drain Traffic Anomaly",
		ProbeId:   probeBId,
		StageName: "excess_range",
	}
	widgetBId, err := bpClient.CreateIbaWidget(ctx, &widgetB)
	if err != nil {
		t.Fatal(err)
	}

	return widgetAId, widgetA, widgetBId, widgetB
}

func testTemplateB(ctx context.Context, t *testing.T, client *Client) ObjectId {
	t.Helper()

	rbt, err := client.GetRackBasedTemplate(ctx, "L2_Virtual")
	require.NoError(t, err)

	rbt.Data.DisplayName = randString(5, "hex")
	rbt.Data.FabricAddressingPolicy = &TemplateFabricAddressingPolicy410Only{
		SpineSuperspineLinks: AddressingSchemeIp46,
		SpineLeafLinks:       AddressingSchemeIp46,
	}
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
		RackInfos:              rbt.Data.RackInfo,
		DhcpServiceIntent:      &rbt.Data.DhcpServiceIntent,
		AntiAffinityPolicy:     rbt.Data.AntiAffinityPolicy,
		AsnAllocationPolicy:    &rbt.Data.AsnAllocationPolicy,
		FabricAddressingPolicy: rbt.Data.FabricAddressingPolicy,
		VirtualNetworkPolicy:   &rbt.Data.VirtualNetworkPolicy,
	})
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, client.DeleteTemplate(ctx, id)) })

	return id
}

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
		VnType:                    VnTypeVxlan,
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

// testFFBlueprintB creates a freeform blueprint with systemCount generic systems.
// The returned []ObjectId represent the requested generic systems.
func testFFBlueprintB(ctx context.Context, t testing.TB, client *Client, systemCount int) (*FreeformClient, []ObjectId) {
	t.Helper()

	id, err := client.CreateFreeformBlueprint(ctx, randString(6, "hex"))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, client.DeleteBlueprint(ctx, id))
	})

	c, err := client.NewFreeformClient(ctx, id)
	require.NoError(t, err)

	dpId, err := c.ImportDeviceProfile(ctx, "Juniper_vEX")
	require.NoError(t, err)

	systemIds := make([]ObjectId, systemCount)
	for i := range systemIds {
		systemIds[i], err = c.CreateSystem(ctx, &FreeformSystemData{
			Type:            SystemTypeInternal,
			Label:           randString(6, "hex"),
			DeviceProfileId: dpId,
		})
		require.NoError(t, err)
	}

	return c, systemIds
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

	vniBeginEnds, err := getRandInts(1, 400000, (rand.Intn(5)+2)*2)
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
		Type:    ResourcePoolTypeAsn,
		PoolIds: []ObjectId{testAsnPool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}

func testResourceGroupInt(ctx context.Context, t testing.TB, client *FreeformClient) (id ObjectId) {
	id, err := client.CreateAllocGroup(ctx, &FreeformAllocGroupData{
		Name:    randString(6, "hex"),
		Type:    ResourcePoolTypeInt,
		PoolIds: []ObjectId{testIntPool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}

func testResourceGroupIpv4(ctx context.Context, t testing.TB, client *FreeformClient) (id ObjectId) {
	id, err := client.CreateAllocGroup(ctx, &FreeformAllocGroupData{
		Name:    randString(6, "hex"),
		Type:    ResourcePoolTypeIpv4,
		PoolIds: []ObjectId{testIpv4Pool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}

func testResourceGroupIpv6(ctx context.Context, t testing.TB, client *FreeformClient) (id ObjectId) {
	id, err := client.CreateAllocGroup(ctx, &FreeformAllocGroupData{
		Name:    randString(6, "hex"),
		Type:    ResourcePoolTypeIpv6,
		PoolIds: []ObjectId{testIpv6Pool(ctx, t, client.client)},
	})
	require.NoError(t, err)

	return id
}
