package apstra

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func randString(n int, style string) string {
	rand.Seed(time.Now().UnixNano())

	var b64Letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")
	var hexLetters = []rune("0123456789abcdef")
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

func randJwt() string {
	return randString(36, "b64") + "." +
		randString(178, "b64") + "." +
		randString(86, "b64")
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

	err := os.MkdirAll(filepath.Dir(fileName), os.FileMode(0600))
	if err != nil {
		return nil, err
	}
	return os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
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

func testBlueprintA(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func() error) {
	bpId, err := client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L3_Collapsed_ESI",
	})
	if err != nil {
		t.Fatal(err)
	}

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(err)
	}

	bpDeleteFunc := func() error {
		return client.DeleteBlueprint(ctx, bpId)
	}

	return bpClient, bpDeleteFunc
}

func testBlueprintB(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func() error) {
	bpId, err := client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual",
	})
	if err != nil {
		t.Fatal(err)
	}

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(err)
	}

	bpDeleteFunc := func() error {
		return client.DeleteBlueprint(ctx, bpId)
	}

	return bpClient, bpDeleteFunc
}

func testBlueprintC(ctx context.Context, t *testing.T, client *Client) (*TwoStageL3ClosClient, func() error) {
	bpId, err := client.CreateBlueprintFromTemplate(context.Background(), &CreateBlueprintFromTemplateRequest{
		RefDesign:  RefDesignDatacenter,
		Label:      randString(5, "hex"),
		TemplateId: "L2_Virtual_EVPN",
	})
	if err != nil {
		t.Fatal(err)
	}

	bpClient, err := client.NewTwoStageL3ClosClient(ctx, bpId)
	if err != nil {
		t.Fatal(err)
	}

	bpDeleteFunc := func() error {
		return client.DeleteBlueprint(ctx, bpId)
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
