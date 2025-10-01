// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package slice_test

import (
	"strconv"
	"testing"

	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
	"github.com/Juniper/apstra-go-sdk/internal/slice"
	"github.com/stretchr/testify/require"
)

var _ internal.IDer = (*stringIDer)(nil)

// stringIDer is used in TestObjectWithID
type stringIDer string

func (s stringIDer) ID() *string {
	if s == "nil" {
		return nil
	}
	return (*string)(&s)
}

func newStringIDer(s string) *stringIDer {
	result := stringIDer(s)
	return &result
}

var _ internal.IDer = (*intIDer)(nil)

// intIDer is used in TestObjectWithID
type intIDer int

func (i intIDer) ID() *string {
	if i == -1 {
		return nil
	}
	return pointer.To(strconv.Itoa(int(i)))
}

func newIntIDer(i int) *intIDer {
	result := intIDer(i)
	return &result
}

func TestObjectWithID(t *testing.T) {
	type testCase struct {
		iders     []internal.IDer
		wantID    string
		expectIdx *int
		expectErr bool
	}

	testCases := map[string]testCase{
		"found_first": {
			iders:     []internal.IDer{newStringIDer("a"), newIntIDer(2), newStringIDer("c")},
			wantID:    "a",
			expectIdx: pointer.To(0),
		},
		"found_middle": {
			iders:     []internal.IDer{newStringIDer("a"), newIntIDer(2), newStringIDer("c")},
			wantID:    "2",
			expectIdx: pointer.To(1),
		},
		"found_last": {
			iders:     []internal.IDer{newStringIDer("a"), newIntIDer(2), newStringIDer("c")},
			wantID:    "c",
			expectIdx: pointer.To(2),
		},
		"not_found": {
			iders:     []internal.IDer{newStringIDer("a"), newIntIDer(2), newStringIDer("c")},
			wantID:    "d",
			expectIdx: nil,
		},
		"found_among_nil": {
			iders:     []internal.IDer{newStringIDer("nil"), newIntIDer(-1), newStringIDer("x")},
			wantID:    "x",
			expectIdx: pointer.To(2),
		},
		"multiple_match": {
			iders:     []internal.IDer{newStringIDer("a"), newIntIDer(1), newStringIDer("a")},
			wantID:    "a",
			expectErr: true,
		},
	}

	for tName, tCase := range testCases {
		t.Run("must_"+tName, func(t *testing.T) {
			t.Parallel()

			result, err := slice.FindByID(tCase.iders, tCase.wantID)
			if tCase.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tCase.expectIdx == nil {
				require.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			require.Same(t, &tCase.iders[*tCase.expectIdx], result)
		})
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var result *internal.IDer
			if tCase.expectErr {
				require.Panics(t, func() {
					_ = slice.MustFindByID(tCase.iders, tCase.wantID)
				})
				return
			} else {
				result = slice.MustFindByID(tCase.iders, tCase.wantID)
			}

			if tCase.expectIdx == nil {
				require.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			require.Same(t, &tCase.iders[*tCase.expectIdx], result)
		})
	}
}

func TestRemove(t *testing.T) {
	type testCase struct {
		d        []any
		i        int
		e        []any
		expPanic bool
	}

	testCases := map[string]testCase{
		"ints_from_beginning": {
			d: []any{1, 2, 3, 4, 5},
			i: 0,
			e: []any{2, 3, 4, 5},
		},
		"ints_from_middle": {
			d: []any{1, 2, 3, 4, 5},
			i: 2,
			e: []any{1, 2, 4, 5},
		},
		"ints_from_end": {
			d: []any{1, 2, 3, 4, 5},
			i: 4,
			e: []any{1, 2, 3, 4},
		},
		"ints_invalid_too_big": {
			d:        []any{1, 2, 3, 4, 5},
			i:        5,
			expPanic: true,
		},
		"ints_invalid_negative": {
			d:        []any{1, 2, 3, 4, 5},
			i:        -1,
			expPanic: true,
		},
		"strings_from_beginning": {
			d: []any{"a", "b", "c", "d", "e"},
			i: 0,
			e: []any{"b", "c", "d", "e"},
		},
		"strings_from_middle": {
			d: []any{"a", "b", "c", "d", "e"},
			i: 2,
			e: []any{"a", "b", "d", "e"},
		},
		"strings_from_end": {
			d: []any{"a", "b", "c", "d", "e"},
			i: 4,
			e: []any{"a", "b", "c", "d"},
		},
		"strings_invalid_too_big": {
			d:        []any{"a", "b", "c", "d", "e"},
			i:        5,
			expPanic: true,
		},
		"strings_invalid_negative": {
			d:        []any{"a", "b", "c", "d", "e"},
			i:        -1,
			expPanic: true,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			if tCase.expPanic {
				require.Panics(t, func() {
					slice.Remove(tCase.d, tCase.i)
				})
				return
			}

			r := slice.Remove(tCase.d, tCase.i)
			require.ElementsMatch(t, tCase.e, r)
		})
	}
}
