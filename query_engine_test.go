package goapstra

import (
	"context"
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
