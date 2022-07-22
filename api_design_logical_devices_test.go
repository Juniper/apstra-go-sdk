package goapstra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestParseLogicalDeviceSpeed(t *testing.T) {
	tests := [][]string{
		{"10000000", "10M"},
		{"10M", "10M"},
		{"10Mbps", "10M"},
		{"10Mb/s", "10M"},
		{"100000000", "100M"},
		{"100M", "100M"},
		{"100Mbps", "100M"},
		{"100Mb/s", "100M"},
		{"1000000000", "1G"},
		{"1000M", "1G"},
		{"1000Mbps", "1G"},
		{"1000Mb/s", "1G"},
		{"1000000000", "1G"},
		{"10G", "10G"},
		{"10Gbps", "10G"},
		{"10Gb/s", "10G"},
		{"10000000000", "10G"},
		{"25G", "25G"},
		{"25Gbps", "25G"},
		{"25Gb/s", "25G"},
		{"25000000000", "25G"},
		{"40G", "40G"},
		{"40Gbps", "40G"},
		{"40Gb/s", "40G"},
		{"40000000000", "40G"},
		{"50G", "50G"},
		{"50Gbps", "50G"},
		{"50Gb/s", "50G"},
		{"50000000000", "50G"},
		{"100G", "100G"},
		{"100Gbps", "100G"},
		{"100Gb/s", "100G"},
		{"100000000000", "100G"},
		{"200G", "200G"},
		{"200Gbps", "200G"},
		{"200Gb/s", "200G"},
		{"200000000000", "200G"},
		{"400G", "400G"},
		{"400Gbps", "400G"},
		{"400Gb/s", "400G"},
		{"400000000000", "400G"},
	}
	for _, test := range tests {
		r := LogicalDevicePortSpeed(test[0]).raw()
		s1 := fmt.Sprintf("%d%s", r.Value, r.Unit)
		s2 := fmt.Sprintf("%s", r.parse())
		if s1 != s2 {
			log.Fatalf("conversion problem: %s %s %s %s", test[0], test[1], s1, s2)
		}
	}
}

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
							Speed: "10G",
							Roles: LogicalDevicePortRoleUnused,
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

func TestRawIfy(t *testing.T) {
	testDev := LogicalDevice{
		DisplayName: "name",
		Id:          "id",
		Panels: []LogicalDevicePanel{{

			PanelLayout: LogicalDevicePanelLayout{
				RowCount:    3,
				ColumnCount: 3,
			},
			PortIndexing: LogicalDevicePortIndexing{
				Order:      PortIndexingVerticalFirst,
				StartIndex: 0,
				Schema:     PortIndexingSchemaAbsolute,
			},
			PortGroups: []LogicalDevicePortGroup{{
				Count: 9,
				Speed: "10G",
				Roles: LogicalDevicePortRoleAccess | LogicalDevicePortRoleSpine,
			}},
		}},
		CreatedAt:      time.Now().Add(-time.Hour * 24),
		LastModifiedAt: time.Now(),
	}
	raw := testDev.raw()
	log.Println(raw.Panels[0].PortGroups[0].Roles)
}

func TestGetLogicalDeviceByName(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	logicalDevices, err := client.getAllLogicalDevices(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	for _, test := range logicalDevices {
		logicalDevice, err := client.getLogicalDeviceByName(context.TODO(), test.DisplayName)
		var ace ApstraClientErr
		if err != nil {
			if !(errors.As(err, &ace) && ace.Type() == ErrMultipleMatch) {
				log.Fatal(err)
			}
			continue
		}
		if logicalDevice.Id != test.Id {
			log.Fatal(fmt.Errorf("expected '%s', got '%s'", test.Id, logicalDevice.Id))
		}
	}
}
