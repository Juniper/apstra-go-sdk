package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestQEQ(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	blueprints, err := client.listAllBlueprintIds(context.TODO())
	if len(blueprints) == 0 {
		t.Skip("no blueprints, cannot perform query test")
	}

	queryResponse, err := client.newQuery(blueprints[0]).SetType(QEQueryTypeStaging).
		Node([]QEEAttribute{
			{"type", QEStringVal("system")},
			{"name", QEStringVal("system")},
			{"role", QEStringVal("spine")},
			{"label", QEStringValIsIn{"spine1", "spine2"}},
		}).
		Out([]QEEAttribute{
			{"type", QEStringVal("logical_device")},
			{"name", QEStringVal("logical_device")},
		}).
		Do()

	if err != nil {
		t.Fatal(err)
	}

	log.Println(queryResponse)
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

func TestParsingQueryInfo(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	bpIds, err := client.listAllBlueprintIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(bpIds) == 0 {
		t.Skip("cannot test blueprint query with no blueprint")
	}

	// the type of info we expect the query to return (a slice of these)
	type qrItem struct {
		LogicalDevice struct {
			Id    string `json:"id"`
			Label string `json:"label"`
		} `json:"n_logical_device"`
		System struct {
			Id    string `json:"id"`
			Label string `json:"label"`
		} `json:"n_system"`
	}

	qr, err := client.NewQuery(bpIds[0]).
		Node([]QEEAttribute{
			{"type", QEStringVal("system")},
			{"name", QEStringVal("n_system")},
			{"system_type", QEStringVal("switch")},
		}).
		Out([]QEEAttribute{
			{"type", QEStringVal("logical_device")},
		}).
		Node([]QEEAttribute{
			{"type", QEStringVal("logical_device")},
			{"name", QEStringVal("n_logical_device")},
		}).
		Do()
	if err != nil {
		t.Fatal(err)
	}

	qResponses := make([]qrItem, len(qr))
	for i, item := range qr {
		err = json.Unmarshal(item, &qResponses[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	if err != nil {
		t.Fatal(err)
	}

	for _, response := range qResponses {
		log.Println(response)
	}

}
