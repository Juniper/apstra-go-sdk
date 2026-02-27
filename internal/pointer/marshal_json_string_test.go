// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pointer_test

import (
	"encoding/json"
	"testing"

	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

func TestMarshalJSONStringWithEmptyAsNull(t *testing.T) {
	mustMarshal := func(s string) json.RawMessage {
		b, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}
		return json.RawMessage(b)
	}

	tests := map[string]struct {
		in   *string
		want json.RawMessage
	}{
		"nil_pointer": {
			in:   nil,
			want: nil,
		},
		"empty_string_becomes_null": {
			in:   pointer.To(""),
			want: json.RawMessage("null"),
		},
		"simple_string": {
			in:   pointer.To("hello"),
			want: mustMarshal("hello"),
		},
		"string_with_quotes": {
			in:   pointer.To(`he said "hello"`),
			want: mustMarshal(`he said "hello"`),
		},
		"string_with_backslash": {
			in:   pointer.To(`c:\windows\path`),
			want: mustMarshal(`c:\windows\path`),
		},
		"string_with_newline_and_tab": {
			in:   pointer.To("line1\nline2\tend"),
			want: mustMarshal("line1\nline2\tend"),
		},
		"string_with_unicode": {
			in:   pointer.To("こんにちは世界"),
			want: mustMarshal("こんにちは世界"),
		},
		"string_with_control_chars": {
			in:   pointer.To("a\x00b\x1fc"),
			want: mustMarshal("a\x00b\x1fc"),
		},
		"string_that_looks_like_json": {
			in:   pointer.To(`{"key":"value"}`),
			want: mustMarshal(`{"key":"value"}`),
		},
		"string_literal_null": {
			in:   pointer.To("null"),
			want: mustMarshal("null"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := pointer.StringMarshalJSONWithEmptyAsNull(tt.in)

			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %q", string(got))
				}
				return
			}

			if got == nil {
				t.Fatalf("expected %q, got nil", string(tt.want))
			}

			if string(got) != string(tt.want) {
				t.Fatalf("mismatch:\nwant: %s\ngot:  %s", tt.want, got)
			}
		})
	}
}
