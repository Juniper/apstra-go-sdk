package apstra

import (
	"testing"
)

func TestNextInterface(t *testing.T) {
	type testCase struct {
		t string
		e string
	}

	testCases := []testCase{
		{t: "xe-0/0/0", e: "xe-0/0/1"},
		{t: "xe-0/0/9", e: "xe-0/0/10"},
	}

	for i, tc := range testCases {
		r := nextInterface(tc.t)
		if tc.e != r {
			t.Fatalf("test case %d: expected %s got %s", i, tc.e, r)
		}
	}
}
