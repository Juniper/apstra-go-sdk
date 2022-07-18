package goapstra

import (
	"context"
	"log"
	"testing"
)

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
	err = client.NewQuery(bpIds[0]).
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
		Do(&qResponse)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("query produced %d results", qResponse.Count)
	for i, item := range qResponse.Items {
		log.Printf("  %d id: '%s', label: '%s', logical_device: '%s'", i, item.System.Id, item.System.Label, item.LogicalDevice.Label)
	}
}
