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

	queryResponse, err := client.NewQuery(blueprints[0]).Node([]QEEAttributes{
		{"type", QEStringVal("system")},
		{"name", QEStringVal("system")},
		{"role", QEStringVal("spine")},
		{"label", QEStringValIsIn{"spine1", "spine2"}},
	}).Out([]QEEAttributes{
		{"type", QEStringVal("logical_device")},
		{"name", QEStringVal("logical_device")},
	}).Do()

	if err != nil {
		t.Fatal(err)
	}

	log.Println(queryResponse)

	//element := QEElement{
	//	qeeType: "node",
	//	attributes: []QEEAttributes{
	//		{"role", QEStringVal("spine")},
	//		{"name", QEStringVal("switch")},
	//		{"label", QEStringValIsIn{"spine1", "spine2"}},
	//	},
	//}
	//log.Println(q.string())

	//q := client.NewQuery("system",
	//	QEEAttributes{"role", "spine"},
	//	QEEAttributes{"name", "switch"},
	//	QEEAttributes{"label", QEStringValIsIn{"spine1", "spine2"}},
	//	QEEAttributes{"exteernal", false},
	//)
	//log.Println(q.String())
}
