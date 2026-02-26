package str_test

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/internal/str"
)

func TestQuoteJSONString_Table(t *testing.T) {
	testCases := map[string]struct {
		in   string
		want string // exact JSON string literal (including surrounding quotes)
	}{
		"simple_ASCII": {
			in:   "hello",
			want: `"hello"`,
		},
		"contains_quote_and_backslash": {
			in:   `he said: "hi" and used \ backslash`,
			want: `"he said: \"hi\" and used \\ backslash"`,
		},
		"newline_and_tab_(short_escapes_are_preserved)": {
			in:   "line1\nline2\tend",
			want: `"line1\nline2\tend"`,
		},
		"NUL_byte_(JSON_must_use_\u0000)": {
			in:   "a\x00b",
			want: `"a\u0000b"`,
		},
		"unit_separator_0x1F_(JSON_must_use_\u001f)": {
			in:   "x\x1fy",
			want: `"x\u001fy"`,
		},
		"vertical_tab_0x0B_(JSON_\u000b)": {
			in:   "a\x0bb",
			want: `"a\u000bb"`,
		},
		"form feed 0x0C (JSON short escape \f)": {
			in:   "a\x0cb",
			want: `"a\fb"`,
		},
		"carriage_return_CR_0x0D_(short_escape_allowed)": {
			in:   "a\rb",
			want: `"a\rb"`,
		},
		"less-than,_ampersand,_greater-than_(json.Marshal_escapes_by_default)": {
			in:   "<&>",
			want: `"\u003c\u0026\u003e"`,
		},
		"solidus_slash_left_as_/": {
			in:   "a/b",
			want: `"a/b"`,
		},
		"non-ASCII_rune_stays_UTF-8_(π)": {
			in:   "π",
			want: `"π"`,
		},
		"emoji_stays_UTF-8_(U+1F600)": {
			in:   "😀",
			want: `"😀"`,
		},
		"mixed_tricky_combo": {
			in: "x\x1f\"\n<&> π \\ / \t",
			// Expect: 0x1f -> \u001f, quotes/backslash escaped, newline \n, < & > escaped,
			// slash kept, tab \t, Unicode kept.
			want: `"x\u001f\"\n\u003c\u0026\u003e π \\ / \t"`,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			got := string(str.QuoteJSONString(tCase.in))
			if got != tCase.want {
				t.Fatalf("QuoteJSONString(%q) = %s, want %s", tCase.in, got, tCase.want)
			}
		})
	}
}
