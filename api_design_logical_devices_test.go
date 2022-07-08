package goapstra

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestListAndGetAllLogicalDevices(t *testing.T) {
	DebugLevel = 2
	clients, _, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(len(clients))

	for clientName, client := range clients {
		if clientName == "mock" {
			continue // todo have I given up on mock testing?
		}
		ids, err := client.listLogicalDeviceIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(ids) <= 0 {
			t.Fatalf("only got %d ids from %s client", len(ids), clientName)
		}
		for _, id := range ids {
			ld, err := client.getLogicalDevice(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("logical device id '%s' name '%s'\n", id, ld.DisplayName)
		}
	}
}

func TestCreateGetUpdateDeleteLogicalDevice(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	var deviceConfigs []LogicalDevice
	for i, indexing := range []string{
		PortIndexingVerticalFirst,
		PortIndexingHorizontalFirst,
	} {
		deviceConfigs = append(deviceConfigs, LogicalDevice{
			DisplayName: fmt.Sprintf("AAAA-%s-%d", t.Name(), i),
			Panels: []LogicalDevicePanel{
				{
					PanelLayout: LogicalDevicePanelLayout{
						RowCount:    2,
						ColumnCount: 2,
					},
					PortIndexing: LogicalDevicePortIndexing{
						Order:      indexing,
						StartIndex: 0,
						Schema:     "absolute",
					},
					PortGroups: []LogicalDevicePortGroup{
						{
							Count: 4,
							Speed: LogicalDevicePortSpeed{
								Unit:  "G",
								Value: 1,
							},
							Roles: []string{"unused"},
						},
					},
				},
			},
		})
	}

	var id []ObjectId
	for i := 0; i < len(deviceConfigs); i++ {
		id = append(id, ObjectId(""))
		id[i], err = client.createLogicalDevice(context.TODO(), &deviceConfigs[i])
		if err != nil {
			t.Fatal(err)
		}

		d, err := client.getLogicalDevice(context.TODO(), id[i])
		if err != nil {
			t.Fatal(err)
		}

		log.Println(d.Id)
		deviceConfigs[i].Panels[0].PortIndexing.StartIndex = 1
		err = client.updateLogicalDevice(context.TODO(), d.Id, &deviceConfigs[i])
		if err != nil {
			log.Fatal(err)
		}

		if i > 0 {
			previous, err := client.getLogicalDevice(context.TODO(), id[i-1])
			if err != nil {
				t.Fatal(err)
			}

			err = client.updateLogicalDevice(context.TODO(), id[i], previous)
			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = client.getLogicalDevice(context.TODO(), id[i])
		if err != nil {
			t.Fatal(err)
		}

		err = client.deleteLogicalDevice(context.TODO(), d.Id)
		if err != nil {
			t.Fatal(err)
		}
	}
}
