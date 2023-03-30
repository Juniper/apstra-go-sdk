package goapstra

import (
	"testing"
)

func TestDcRoutingPoliciesStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "", intType: DcRoutingPolicyImportPolicyNone, stringType: dcRoutingPolicyImportPolicyNone},
		{stringVal: "default_only", intType: DcRoutingPolicyImportPolicyDefaultOnly, stringType: dcRoutingPolicyImportPolicyDefaultOnly},
		{stringVal: "all", intType: DcRoutingPolicyImportPolicyAll, stringType: dcRoutingPolicyImportPolicyAll},
		{stringVal: "extra_only", intType: DcRoutingPolicyImportPolicyExtraOnly, stringType: dcRoutingPolicyImportPolicyExtraOnly},

		{stringVal: "", intType: DcRoutingPolicyTypeNone, stringType: dcRoutingPolicyTypeNone},
		{stringVal: "default_immutable", intType: DcRoutingPolicyTypeDefault, stringType: dcRoutingPolicyTypeDefault},
		{stringVal: "user_defined", intType: DcRoutingPolicyTypeUser, stringType: dcRoutingPolicyTypeUser},

		{stringVal: "", intType: PrefixFilterActionNone, stringType: prefixFilterActionNone},
		{stringVal: "permit", intType: PrefixFilterActionPermit, stringType: prefixFilterActionPermit},
		{stringVal: "deny", intType: PrefixFilterActionDeny, stringType: prefixFilterActionDeny},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}

func TestAllDcRoutingPolicyImportPolicy(t *testing.T) {
	result := AllDcRoutingPolicyImportPolicies()
	if len(result) != 4 {
		t.Fatalf("expected 4 results, got %d", len(result))
	}
	expected := []DcRoutingPolicyImportPolicy{0, 1, 2, 3}
	for i := range result {
		if result[i] != expected[i] {
			t.Fatalf("index %d expected %q, got %q", i, expected[i], result[i])
		}
	}
}

func TestAllPrefixFilterActions(t *testing.T) {
	result := AllPrefixFilterActions()
	if len(result) != 3 {
		t.Fatalf("expected 4 results, got %d", len(result))
	}
	expected := []PrefixFilterAction{0, 1, 2}
	for i := range result {
		if result[i] != expected[i] {
			t.Fatalf("index %d expected %q, got %q", i, expected[i], result[i])
		}
	}
}
