package slice_test

import (
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"strconv"
	"testing"

	"github.com/Juniper/apstra-go-sdk/internal/slice"
	"github.com/stretchr/testify/require"
)

func TestUniq_and_IsUniq(t *testing.T) {
	type testCase struct {
		input    []string
		expected []string
		isUniq   bool
	}

	tests := map[string]testCase{
		"empty_slice": {
			input:    []string{},
			expected: []string{},
			isUniq:   true,
		},
		"nil_slice": {
			input:    ([]string)(nil),
			expected: ([]string)(nil),
			isUniq:   true,
		},
		"all_unique": {
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
			isUniq:   true,
		},
		"some_duplicates": {
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
			isUniq:   false,
		},
		"all_duplicates": {
			input:    []string{"x", "x", "x"},
			expected: []string{"x"},
			isUniq:   false,
		},
		"consecutive_duplicates": {
			input:    []string{"a", "a", "b", "b", "c", "c"},
			expected: []string{"a", "b", "c"},
			isUniq:   false,
		},
	}

	for tName, tCase := range tests {
		t.Run(tName, func(t *testing.T) {
			result := slice.Uniq(tCase.input)
			require.Equal(t, tCase.expected, result)

			resultUniq := slice.IsUniq(tCase.input)
			require.Equal(t, tCase.isUniq, resultUniq)
		})
	}
}

// stringIDer is used in TestObjectWithID
type stringIDer string

func (s stringIDer) ID() *string {
	if s == "nil" {
		return nil
	}
	return (*string)(&s)
}

// intIDer is used in TestObjectWithID
type intIDer int

func (i intIDer) ID() *string {
	if i == -1 {
		return nil
	}
	return pointer.To(strconv.Itoa(int(i)))
}

func TestObjectWithID(t *testing.T) {
	type testCase struct {
		iders     []slice.IDer
		wantID    string
		expectIdx *int
	}

	testCases := map[string]testCase{
		"found_first": {
			iders:     []slice.IDer{stringIDer("a"), intIDer(2), stringIDer("c")},
			wantID:    "a",
			expectIdx: pointer.To(0),
		},
		"found_middle": {
			iders:     []slice.IDer{stringIDer("a"), intIDer(2), stringIDer("c")},
			wantID:    "2",
			expectIdx: pointer.To(1),
		},
		"found_last": {
			iders:     []slice.IDer{stringIDer("a"), intIDer(2), stringIDer("c")},
			wantID:    "c",
			expectIdx: pointer.To(2),
		},
		"not_found": {
			iders:     []slice.IDer{stringIDer("a"), intIDer(2), stringIDer("c")},
			wantID:    "d",
			expectIdx: nil,
		},
		"found_among_nil": {
			iders:     []slice.IDer{stringIDer("nil"), intIDer(-1), stringIDer("x")},
			wantID:    "x",
			expectIdx: pointer.To(2),
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			result := slice.ObjectWithID(tCase.iders, tCase.wantID)
			if tCase.expectIdx == nil {
				require.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			require.Same(t, &tCase.iders[*tCase.expectIdx], result)
		})
	}
}
