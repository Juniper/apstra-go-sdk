package goapstra

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"sort"
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

func randId() ObjectId {
	return ObjectId(
		randString(8, "hex") + "-" +
			randString(4, "hex") + "-" +
			randString(4, "hex") + "-" +
			randString(4, "hex") + "-" +
			randString(12, "hex"))
}

func randJwt() string {
	return randString(36, "b64") + "." +
		randString(178, "b64") + "." +
		randString(86, "b64")
}

func TestKeyLogWriter(t *testing.T) {
	envVarName := randString(10, "hex")

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Fatal(err)
		}
	}()

	testFileName := filepath.Join(dir, randString(10, "b64"))

	err = os.Setenv(envVarName, testFileName)
	if err != nil {
		t.Fatal(err)
	}

	klw, err := keyLogWriter(envVarName)
	if err != nil {
		t.Fatal(err)
	}

	data := randString(100, "b64")
	_, err = klw.Write([]byte(data))
	err = klw.Close()
	if err != nil {
		t.Fatal(err)
	}

	result, err := ioutil.ReadFile(testFileName)
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != data {
		t.Fatal("data read and written do not match")
	}
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

func TestInvertRangesInRange(t *testing.T) {
	_, err := invertRangesInRange(100, 1, nil)
	if err == nil {
		log.Fatalf("expected max/min error")
	}

	var testBegin, testEnd []uint32
	var testUsed [][]NewAsnRange

	testBegin = append(testBegin, 1)
	testEnd = append(testEnd, 100)
	testUsed = append(testUsed, []NewAsnRange{{B: 10, E: 19}, {B: 30, E: 39}, {B: 90, E: 99}})

	testBegin = append(testBegin, 1)
	testEnd = append(testEnd, 100)
	testUsed = append(testUsed, []NewAsnRange{{B: 1, E: 19}, {B: 30, E: 39}, {B: 90, E: 100}})

	testBegin = append(testBegin, 1)
	testEnd = append(testEnd, 100)
	testUsed = append(testUsed, nil)

	for i := range testBegin {
		result, err := invertRangesInRange(testBegin[i], testEnd[i], testUsed[i])
		if err != nil {
			log.Fatal(err)
		}
		log.Println(result)
	}

	_, err = invertRangesInRange(1, 100, []NewAsnRange{{B: 0, E: 19}, {B: 30, E: 39}, {B: 90, E: 100}})
	if err == nil {
		log.Fatal(fmt.Errorf("expected to error on minimum range, but did not"))
	}

	_, err = invertRangesInRange(1, 100, []NewAsnRange{{B: 1, E: 19}, {B: 0, E: 39}, {B: 90, E: 100}})
	if err == nil {
		log.Fatal(fmt.Errorf("expected to error on minimum range, but did not"))
	}

	_, err = invertRangesInRange(1, 100, []NewAsnRange{{B: 0, E: 30}, {B: 30, E: 39}, {B: 90, E: 100}})
	if err == nil {
		log.Fatal(fmt.Errorf("expected to error on range overlap, but did not"))
	}

	_, err = invertRangesInRange(1, 100, []NewAsnRange{{B: 0, E: 19}, {B: 30, E: 39}, {B: 90, E: 101}})
	if err == nil {
		log.Fatal(fmt.Errorf("expected to error on maximum range, but did not"))
	}

	_, err = invertRangesInRange(1, 100, []NewAsnRange{{B: 0, E: 19}, {B: 30, E: 39}, {B: 90, E: 101}})
	if err == nil {
		log.Fatal(fmt.Errorf("expected to error on maximum range, but did not"))
	}

}

// invertRangesInRange was designed to find free space in ASN pool resources.
// Valid ASNs are 1-4294967295.
// If current ASN pools consume 50-100, 64512-65534 and 4200000000-4294967294,
// we'd expect to get back [{1,49}{101,64511}{65535,4199999999}{4294967295,4294967295}]
func invertRangesInRange(min, max uint32, used []NewAsnRange) ([]NewAsnRange, error) {
	if min > max {
		return nil, fmt.Errorf("min > max: %d > %d", min, max) // bad input
	}
	if len(used) == 0 {
		return []NewAsnRange{{B: min, E: max}}, nil // nothing used, return entire range
	}
	sort.Slice(used, func(i, j int) bool {
		return used[i].B < used[j].B
	})

	var result []NewAsnRange
	if used[0].B > min { // if there's room, create the first result item
		result = append(result, NewAsnRange{min, used[0].B - 1})
	}
	for i := 0; i <= len(used)-1; i++ {
		if used[i].B > used[i].E {
			return nil, fmt.Errorf("inverted range element: %s", used[i].String())
		}
		if used[i].B < min || used[i].E > max {
			return nil, fmt.Errorf("'%s' out of of range: min %d, max %d", used[i].String(), min, max)
		}
		if i != len(used)-1 { // don't look past the end of the slice
			if asnOverlap(used[i], used[i+1]) {
				return nil, fmt.Errorf("overlapping ranges %s, %s", used[i].String(), used[i+1].String())
			}
			if used[i].E < used[i+1].B {
				result = append(result, NewAsnRange{used[i].E + 1, used[i+1].B - 1})
			}
		}
	}
	if used[len(used)-1].E < max { // if there's room, create the final result item
		result = append(result, NewAsnRange{used[len(used)-1].E + 1, max})
	}

	return result, nil
}

func asnOverlap(a, b NewAsnRange) bool {
	if a.B >= b.B && a.B <= b.E { // begin 'a' falls within 'b'
		return true
	}
	if a.E <= b.E && a.E >= b.B { // end 'a' falls within 'b'
		return true
	}
	if b.B >= a.B && b.B <= a.E { // begin 'b' falls within 'a'
		return true
	}
	if b.E <= a.E && b.E >= a.B { // end 'b' falls within 'a'
		return true
	}
	return false // no overlap
}
