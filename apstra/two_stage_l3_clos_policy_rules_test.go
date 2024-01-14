package apstra

import (
	"testing"
)

func TestPortRangesSorted(t *testing.T) {
	type testCase struct {
		portRanges PortRanges
		expected   string
	}
	testCases := map[string]testCase{
		"single_val": {
			portRanges: PortRanges{{1, 1}},
			expected:   "1",
		},
		"multiple_val": {
			portRanges: PortRanges{{1, 1}, {2, 2}, {3, 3}},
			expected:   "1,2,3",
		},
		"single_range": {
			portRanges: PortRanges{{1, 3}},
			expected:   "1-3",
		},
		"multiple_range": {
			portRanges: PortRanges{{1, 3}, {5, 7}, {9, 11}},
			expected:   "1-3,5-7,9-11",
		},
		"mixed": {
			portRanges: PortRanges{{1, 1}, {3, 5}, {7, 7}, {9, 11}},
			expected:   "1,3-5,7,9-11",
		},
		"single_val_unordered": {
			portRanges: PortRanges{{3, 3}, {1, 1}, {2, 2}},
			expected:   "1,2,3",
		},
		"mixed_unordered": {
			portRanges: PortRanges{{9, 11}, {3, 5}, {7, 7}, {1, 1}},
			expected:   "1,3-5,7,9-11",
		},
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			result := tCase.portRanges.string()
			if tCase.expected != result {
				t.Fatalf("expected %q, got %q", tCase.expected, result)
			}
		})
	}
}
