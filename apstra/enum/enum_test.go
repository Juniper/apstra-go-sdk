package enum_test

import (
	"sort"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

func TestEnumValues(t *testing.T) {
	type testCase struct {
		enum1   enum.Value
		string1 string
		enum2   enum.Value
		string2 string
		equal   bool
	}

	testCases := map[string]testCase{
		"identical": {
			enum1:   enum.SizeLarge(),
			string1: "large",
			enum2:   enum.SizeLarge(),
			string2: "large",
			equal:   true,
		},
		"same_type": {
			enum1:   enum.SizeSmall(),
			string1: "small",
			enum2:   enum.SizeLarge(),
			string2: "large",
			equal:   false,
		},
		"same_value": {
			enum1:   enum.FlavorChocolate(),
			string1: "chocolate",
			enum2:   enum.SauceChocolate(),
			string2: "chocolate",
			equal:   false,
		},
	}

	checkString := func(t testing.TB, d string, s string, e enum.Value) {
		t.Helper()
		if s != e.String() {
			t.Errorf("%s: expected %s, got %s", d, s, e.String())
		}
	}

	checkEqual := func(t testing.TB, d string, e bool, a, b enum.Value) {
		t.Helper()
		if e && !a.Equal(b) {
			t.Errorf("%s expected enum %q of type %q to be equal enum %q of type %q, but it is not", d, a.String(), a.Type(), b.String(), b.Type())
		}
		if !e && a.Equal(b) {
			t.Errorf("%s did not expect enum %q of type %q to be equal enum %q of type %q, but it is", d, a.String(), a.Type(), b.String(), b.Type())
		}
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			checkString(t, "enum1 string", tCase.string1, tCase.enum1)
			checkString(t, "enum2 string", tCase.string2, tCase.enum2)
			checkEqual(t, "testCase enum1 and enum2", tCase.equal, tCase.enum1, tCase.enum2)
			e := enum.New(tCase.enum1.Type(), tCase.string1)
			checkEqual(t, "testCase enum1 and enum from string", true, tCase.enum1, e)
			e = enum.New(tCase.enum2.Type(), tCase.string2)
			checkEqual(t, "testCase enum2 and enum from string", true, tCase.enum2, e)
		})
	}
}

func TestNewEnum(t *testing.T) {
	type testCase struct {
		t      enum.EnumType
		s      string
		expNil bool
	}

	testCases := map[string]testCase{
		"valid": {
			t:      enum.Size,
			s:      "small",
			expNil: false,
		},
		"bad_string": {
			t:      enum.Size,
			s:      "bogus",
			expNil: true,
		},
		"bad_type": {
			t:      -1,
			s:      "small",
			expNil: true,
		},
		"bad_both": {
			t:      -1,
			s:      "bogus",
			expNil: true,
		},
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			e := enum.New(tCase.t, tCase.s)
			if tCase.expNil != (e == nil) {
				t.Errorf("expected e == nil: %t, got %t", tCase.expNil, e == nil)
			}
		})
	}
}

func TestEnumTypes(t *testing.T) {
	type testCase struct {
		enumType enum.EnumType
		strings  []string
	}

	testCases := map[string]testCase{
		"valid": {
			enumType: enum.Size,
			strings:  []string{"small", "medium", "large"},
		},
		"invalid": {
			enumType: -1,
			strings:  nil,
		},
	}

	compareStringSlices := func(t testing.TB, a, b []string) {
		t.Helper()

		if (a == nil) != (b == nil) {
			t.Errorf("a == nil: %t; b == nil: %t", a == nil, b == nil)
		}

		if len(a) != len(b) {
			t.Errorf("a has %d items, b has %d items", len(a), len(b))
		}

		sort.Strings(a)
		sort.Strings(b)
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				t.Errorf("a[%d] == %q; b[%d] == %q", i, a[i], i, b[i])
			}
		}
	}

	for tName, tCase := range testCases {
		tName, tCase := tName, tCase
		t.Run(tName, func(t *testing.T) {
			t.Parallel()
			compareStringSlices(t, tCase.strings, tCase.enumType.Strings())
			values := tCase.enumType.Values()
			if (values == nil) != (tCase.strings == nil) {
				t.Fatalf("Values() returned nil: %t; expected nil: %t", values == nil, tCase.strings == nil)
			}

			stringsFromValues := make([]string, len(values))
			for i, v := range values {
				stringsFromValues[i] = v.String()
			}
			compareStringSlices(t, tCase.strings, stringsFromValues)
		})
	}
}
