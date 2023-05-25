//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestQEAttrValsString(t *testing.T) {
	type testCase struct {
		v QEAttrVal
		e string
	}

	testCases := []testCase{
		{
			v: QEStringVal("foo"),
			e: "'foo'",
		},
		{
			v: QEStringVal("\"bar\""),
			e: "'\"bar\"'",
		},
		{
			v: QEStringVal("123"),
			e: "'123'",
		},
		{
			v: QEStringVal(""),
			e: "''",
		},
		{
			v: QEStringValIsIn{"foo", "bar"},
			e: "is_in(['foo','bar'])",
		},
		{
			v: QEStringValNotIn{"foo", "bar"},
			e: "not_in(['foo','bar'])",
		},
		{
			v: QEIntVal(7),
			e: "7",
		},
		{
			v: QEIntVal(-7),
			e: "-7",
		},
	}

	for i, tc := range testCases {
		r := tc.v.String()
		if tc.e != r {
			t.Errorf("test case %d expected %q got %q", i, tc.e, r)
		}
	}
}

func TestQEEAttributeString(t *testing.T) {
	test := []struct {
		expected string
		test     QEEAttribute
	}{
		{
			"empty_string=''",
			QEEAttribute{Key: "empty_string", Value: QEStringVal("")},
		},
		{
			"foo='FOO'",
			QEEAttribute{Key: "foo", Value: QEStringVal("FOO")},
		},
		{
			"empty_val_is_in=is_in([])",
			QEEAttribute{Key: "empty_val_is_in", Value: QEStringValIsIn{}},
		},
		{
			"val_is_in_foo_bar=is_in(['foo','bar'])",
			QEEAttribute{Key: "val_is_in_foo_bar", Value: QEStringValIsIn{"foo", "bar"}},
		},
		{
			"empty_val_not_in=not_in([])",
			QEEAttribute{Key: "empty_val_not_in", Value: QEStringValNotIn{}},
		},
		{
			"val_is_in_foo_bar=is_in(['foo','bar'])",
			QEEAttribute{Key: "val_is_in_foo_bar", Value: QEStringValIsIn{"foo", "bar"}},
		},
	}

	for _, testData := range test {
		result := testData.test.String()
		if testData.expected != result {
			t.Fatalf("expected '%s', got '%s'", testData.expected, result)
		}
	}
}

func TestQueryString(t *testing.T) {
	x := PathQuery{}
	y := x.Node([]QEEAttribute{
		{"type", QEStringVal("system")},
		{"name", QEStringVal("n_system")},
		{"system_type", QEStringVal("switch")},
	}).
		Out([]QEEAttribute{{"type", QEStringVal("logical_device")}}).
		Node([]QEEAttribute{
			{"type", QEStringVal("logical_device")},
		}).
		In([]QEEAttribute{{"type", QEStringVal("logical_device")}}).
		Node([]QEEAttribute{
			{"type", QEStringVal("interface_map")},
			{"name", QEStringVal("n_interface_map")},
		}).
		String()
	log.Println("\n", y)
}

func TestParsingQueryInfo(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		bpClient, bpDel := testBlueprintA(ctx, t, client.client)
		defer func() {
			err = bpDel(ctx)
			if err != nil {
				t.Fatal(err)
			}
		}()

		// the type of info we expect the query to return (a slice of these)
		var qResponse struct {
			Count int `json:"count"`
			Items []struct {
				LogicalDevice struct {
					Id    string `json:"id"`
					Label string `json:"label"`
				} `json:"n_logical_device"`
				System struct {
					Id    string `json:"id"`
					Label string `json:"label"`
				} `json:"n_system"`
			} `json:"items"`
		}
		log.Printf("testing NewQuery() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.NewQuery(bpClient.Id()).
			Node([]QEEAttribute{
				{"type", QEStringVal("system")},
				{"name", QEStringVal("n_system")},
				{"role", QEStringValIsIn{"superspine", "spine", "leaf"}},
				{"external", QEBoolVal(false)},
			}).
			Out([]QEEAttribute{
				{"type", QEStringVal("logical_device")},
			}).
			Node([]QEEAttribute{
				{"type", QEStringVal("logical_device")},
				{"name", QEStringVal("n_logical_device")},
			}).
			Do(ctx, &qResponse)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("query produced %d results", qResponse.Count)
		for i, item := range qResponse.Items {
			log.Printf("  %d id: '%s', label: '%s', logical_device: '%s'", i, item.System.Id, item.System.Label, item.LogicalDevice.Label)
		}
	}
}

func TestQueryMatchString(t *testing.T) {
	expected := "match(node(type='system',name='n_system',role=is_in(['superspine','spine','leaf']),external=False).in_(type='tag').node(type='tag',label='tag_a',name='n_tag_a'),node(type='system',name='n_system',role=is_in(['superspine','spine','leaf']),external=False).in_(type='tag').node(type='tag',label='tag_b',name='n_tag_b'))"

	queryTagA := new(PathQuery).
		Node([]QEEAttribute{
			{"type", QEStringVal("system")},
			{"name", QEStringVal("n_system")},
			{"role", QEStringValIsIn{"superspine", "spine", "leaf"}},
			{"external", QEBoolVal(false)},
		}).
		In([]QEEAttribute{
			{"type", QEStringVal("tag")},
		}).
		Node([]QEEAttribute{
			{"type", QEStringVal("tag")},
			{"label", QEStringVal("tag_a")},
			{"name", QEStringVal("n_tag_a")},
		})

	queryTagB := new(PathQuery).
		Node([]QEEAttribute{
			{"type", QEStringVal("system")},
			{"name", QEStringVal("n_system")},
			{"role", QEStringValIsIn{"superspine", "spine", "leaf"}},
			{"external", QEBoolVal(false)},
		}).
		In([]QEEAttribute{
			{"type", QEStringVal("tag")},
		}).
		Node([]QEEAttribute{
			{"type", QEStringVal("tag")},
			{"label", QEStringVal("tag_b")},
			{"name", QEStringVal("n_tag_b")},
		})

	result := new(MatchQuery).Match(queryTagA).Match(queryTagB).String()

	if result != expected {
		t.Fatalf("expected %q, got %q\n", expected, result)
	}

	log.Println(result)
}

func TestMatchQueryDistinct_String(t *testing.T) {
	type testCase struct {
		mqd      MatchQueryDistinct
		expected string
	}

	testCases := []testCase{
		{
			mqd:      nil,
			expected: "[]",
		},
		{
			mqd:      MatchQueryDistinct{},
			expected: "[]",
		},
		{
			mqd:      MatchQueryDistinct{"foo", "bar"},
			expected: "['foo','bar']",
		},
	}

	for i, tc := range testCases {
		result := tc.mqd.String()
		if tc.expected != result {
			t.Fatalf("testcase %d expected %q got %q", i, tc.expected, result)
		}
	}
}

func TestMatchQueryElement_String(t *testing.T) {
	mqe := MatchQueryElement{
		mqeType: "distinct",
		value:   MatchQueryDistinct{"foo", "bar"},
		next:    nil,
	}

	result := mqe.String()
	expected := "distinct(['foo','bar'])"

	if expected != result {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestMatchQueryDistinct(t *testing.T) {
	pq1 := new(PathQuery).
		Node([]QEEAttribute{
			{Key: "type", Value: QEStringVal("system")},
			{Key: "name", Value: QEStringVal("n_switch")},
			{Key: "system_type", Value: QEStringVal("switch")},
		}).
		Out([]QEEAttribute{
			{Key: "type", Value: QEStringVal("hosted_interfaces")},
		}).
		Node([]QEEAttribute{
			{Key: "type", Value: QEStringVal("interface")},
			{Key: "if_type", Value: QEStringVal("loopback")},
			{Key: "name", Value: QEStringVal("n_interface")},
		})

	type testCase struct {
		q QEQuery
		e string
	}

	testCases := []testCase{
		{
			q: new(MatchQuery).
				Match(pq1).
				Distinct(MatchQueryDistinct{"n_switch", "n_interface"}),
			e: "match(" + pq1.String() + ").distinct(['n_switch','n_interface'])",
		},
		{
			q: new(MatchQuery).
				Match(pq1).
				Distinct(MatchQueryDistinct{"n_switch", "n_interface"}).
				Distinct(MatchQueryDistinct{"foo", "bar"}),
			e: "match(" + pq1.String() + ").distinct(['n_switch','n_interface']).distinct(['foo','bar'])",
		},
	}

	for i, tc := range testCases {
		result := tc.q.String()
		if tc.e != result {
			t.Fatalf("test case %d expected %q, got %q", i, tc.e, result)
		}
	}
}

func TestPathQueryWhere(t *testing.T) {
	type testCase struct {
		q QEQuery
		e string
	}

	testCases := []testCase{
		{
			q: new(PathQuery).
				Node([]QEEAttribute{
					{Key: "type", Value: QEStringVal("system")},
					{Key: "name", Value: QEStringVal("n_switch")},
					{Key: "system_type", Value: QEStringVal("switch")},
				}).
				Out([]QEEAttribute{
					{Key: "type", Value: QEStringVal("hosted_interfaces")},
				}).
				Node([]QEEAttribute{
					{Key: "type", Value: QEStringVal("interface")},
					{Key: "if_type", Value: QEStringVal("loopback")},
					{Key: "name", Value: QEStringVal("n_interface")},
				}).
				Where("lambda n_switch: n_switch.role in ('leaf', 'spine')"),
			e: "node(type='system',name='n_switch',system_type='switch')." +
				"out(type='hosted_interfaces')." +
				"node(type='interface',if_type='loopback',name='n_interface')." +
				"where(lambda n_switch: n_switch.role in ('leaf', 'spine'))",
		},
	}

	for i, tc := range testCases {
		r := tc.q.String()
		if tc.e != r {
			t.Fatalf("test case %d expected %q, got: %q", i, tc.e, r)
		}
	}
}

func TestMatchQueryWhere(t *testing.T) {
	pq1 := new(PathQuery).
		Node([]QEEAttribute{
			{Key: "type", Value: QEStringVal("system")},
			{Key: "name", Value: QEStringVal("n_switch")},
			{Key: "system_type", Value: QEStringVal("switch")},
		}).
		Out([]QEEAttribute{
			{Key: "type", Value: QEStringVal("hosted_interfaces")},
		}).
		Node([]QEEAttribute{
			{Key: "type", Value: QEStringVal("interface")},
			{Key: "if_type", Value: QEStringVal("loopback")},
			{Key: "name", Value: QEStringVal("n_interface")},
		})

	type testCase struct {
		q QEQuery
		e string
	}

	testCases := []testCase{
		{
			q: new(MatchQuery).
				Match(pq1).
				Distinct(MatchQueryDistinct{"n_switch", "n_interface"}).
				Where("foo"),
			e: "match(" + pq1.String() + ").distinct(['n_switch','n_interface']).where(foo)",
		},
		{
			q: new(MatchQuery).
				Match(pq1).
				Distinct(MatchQueryDistinct{"n_switch", "n_interface"}).
				Distinct(MatchQueryDistinct{"foo", "bar"}).
				Where("foo"),
			e: "match(" + pq1.String() + ").distinct(['n_switch','n_interface']).distinct(['foo','bar']).where(foo)",
		},
	}

	for i, tc := range testCases {
		r := tc.q.String()
		if tc.e != r {
			t.Fatalf("test case %d expected %q, got: %q", i, tc.e, r)
		}
	}
}
