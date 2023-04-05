//go:build integration
// +build integration

package apstra

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

func TestQueryString(t *testing.T) {
	x := QEQuery{}
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
		string()
	log.Println("\n", y)
}

func TestParsingQueryInfo(t *testing.T) {
	ctx := context.Background()
	clients, err := getTestClients(ctx)
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
			Do(&qResponse)
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("query produced %d results", qResponse.Count)
		for i, item := range qResponse.Items {
			log.Printf("  %d id: '%s', label: '%s', logical_device: '%s'", i, item.System.Id, item.System.Label, item.LogicalDevice.Label)
		}
	}
}
