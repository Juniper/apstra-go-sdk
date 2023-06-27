package apstra

import (
	"encoding/json"
	"log"
	"testing"
)

func TestConnectivityTemplateStrings(t *testing.T) {
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
		{stringVal: "", intType: ObjPolicyTypeNameNone, stringType: objPolicyTypeNameNone},
		{stringVal: "batch", intType: ObjPolicyTypeNameBatch, stringType: objPolicyTypeNameBatch},
		{stringVal: "pipeline", intType: ObjPolicyTypeNamePipeline, stringType: objPolicyTypeNamePipeline},
		{stringVal: "AttachBgpOverSubinterfacesOrSvi", intType: ObjPolicyTypeNameBgpOverSubinterfacesOrSvi, stringType: objPolicyTypeNameBgpOverSubinterfacesOrSvi},
		{stringVal: "AttachBgpWithPrefixPeeringForSviOrSubinterface", intType: ObjPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface, stringType: objPolicyTypeNameBgpWithPrefixPeeringForSviOrSubinterface},
		{stringVal: "AttachCustomStaticRoute", intType: ObjPolicyTypeNameCustomStaticRoute, stringType: objPolicyTypeNameCustomStaticRoute},
		{stringVal: "AttachExistingRoutingPolicy", intType: ObjPolicyTypeNameRoutingPolicy, stringType: objPolicyTypeNameRoutingPolicy},
		{stringVal: "AttachIpEndpointWithBgpNsxt", intType: ObjPolicyTypeNameIpEndpointWithBgpNsxt, stringType: objPolicyTypeNameIpEndpointWithBgpNsxt},
		{stringVal: "AttachLogicalLink", intType: ObjPolicyTypeNameLogicalLink, stringType: objPolicyTypeNameLogicalLink},
		{stringVal: "AttachMultipleVLAN", intType: ObjPolicyTypeNameMultipleVLAN, stringType: objPolicyTypeNameMultipleVLAN},
		{stringVal: "AttachRoutingZoneConstraint", intType: ObjPolicyTypeNameRoutingZoneConstraint, stringType: objPolicyTypeNameRoutingZoneConstraint},
		{stringVal: "AttachSingleVlan", intType: ObjPolicyTypeNameSingleVlan, stringType: objPolicyTypeNameSingleVlan},
		{stringVal: "AttachStaticRoute", intType: ObjPolicyTypeNameStaticRoute, stringType: objPolicyTypeNameStaticRoute},
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

func TestObjPolicyAttributesMarshal(t *testing.T) {
	type testCase struct {
		attributes TwoStageL3ClosObjPolicyAttributes
		expected   string
	}

	testCases := []testCase{
		{
			attributes: ObjPolicySingleVlanAttributes{
				VnNodeId: "abc",
				Tagged:   false,
			},
			expected: `{"vn_node_id":"abc","tag_type":"untagged"}`,
		},
		{
			attributes: ObjPolicySingleVlanAttributes{
				VnNodeId: "def",
				Tagged:   true,
			},
			expected: `{"vn_node_id":"def","tag_type":"vlan_tagged"}`,
		},
	}

	for i, tc := range testCases {
		raw, err := tc.attributes.marshal()
		if err != nil {
			t.Fatalf("test case %d failed with error - %s", i, err)
		}
		if tc.expected != string(raw) {
			t.Fatalf("test case %d expected %q got %q", i, tc.expected, string(raw))
		}
	}
}

func TestTwoStageL3ClosObjPolicy(t *testing.T) {
	a := TwoStageL3ClosObjPolicy{
		Description:    "description",
		Tags:           []string{"foo", "bar"},
		Label:          "label",
		PolicyTypeName: ObjPolicyTypeNameSingleVlan, // todo eliminate this
		Attributes: ObjPolicySingleVlanAttributes{
			VnNodeId: "aaaa",
			Tagged:   false,
		},
	}

	r, err := a.Raw()
	if err != nil {
		t.Fatal(err)
	}

	result := struct {
		Policies []RawTwoStageL3ClosObjPolicy `json:"policies"`
	}{
		Policies: r,
	}

	x, err := json.Marshal(&result)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(string(x))
}
