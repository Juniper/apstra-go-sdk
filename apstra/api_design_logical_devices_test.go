//go:build integration
// +build integration

package apstra

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
)

var testSpeedStrings = [][]string{
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

func TestParseLogicalDeviceSpeed(t *testing.T) {
	for _, test := range testSpeedStrings {
		r := LogicalDevicePortSpeed(test[0]).raw()
		s1 := fmt.Sprintf("%d%s", r.Value, r.Unit)
		s2 := string(r.parse())
		if s1 != s2 {
			t.Fatalf("conversion problem: %s %s %s %s", test[0], test[1], s1, s2)
		}
	}
}

func TestListAndGetAllLogicalDevices(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listLogicalDeviceIds() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		ids, err := client.client.listLogicalDeviceIds(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		if len(ids) <= 0 {
			t.Fatalf("only got %d ids", len(ids))
		}
		for _, i := range samples(len(ids)) {
			id := ids[i]
			log.Printf("testing GetLogicalDevice(%s) against %s %s (%s)", id, client.clientType, clientName, client.client.ApiVersion())
			ld, err := client.client.GetLogicalDevice(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("logical device id '%s' name '%s'\n", id, ld.Data.DisplayName)
		}
	}
}

func TestCreateGetUpdateDeleteLogicalDevice(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	indexingTypes := []string{
		PortIndexingVerticalFirst,
		PortIndexingHorizontalFirst,
	}

	randStr := randString(5, "hex")

	for clientName, client := range clients {
		deviceConfigs := make([]LogicalDeviceData, len(indexingTypes))
		for i, indexing := range indexingTypes {
			deviceConfigs[i] = LogicalDeviceData{
				DisplayName: fmt.Sprintf("AAAA-%s-%s-%d", t.Name(), randStr, i),
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
			}
		}

		id := make([]ObjectId, len(deviceConfigs))
		for i, devCfg := range deviceConfigs {
			log.Printf("testing createLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			id[i], err = client.client.createLogicalDevice(context.TODO(), devCfg.raw())
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing GetLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			d, err := client.client.GetLogicalDevice(context.TODO(), id[i])
			if err != nil {
				t.Fatal(err)
			}

			log.Println(d.Id)
			devCfg.Panels[0].PortIndexing.StartIndex = 1
			log.Printf("testing updateLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.updateLogicalDevice(context.TODO(), d.Id, devCfg.raw())
			if err != nil {
				t.Fatal(err)
			}

			if i > 0 {
				log.Printf("testing GetLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				previous, err := client.client.GetLogicalDevice(context.TODO(), id[i-1])
				if err != nil {
					t.Fatal(err)
				}

				previous.Data.DisplayName = d.Data.DisplayName

				log.Printf("testing updateLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				err = client.client.updateLogicalDevice(context.TODO(), id[i], previous.Data.raw())
				if err != nil {
					t.Fatal(err)
				}
			}

			log.Printf("testing GetLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			_, err = client.client.GetLogicalDevice(context.TODO(), id[i])
			if err != nil {
				t.Fatal(err)
			}
		}
		for _, i := range id {
			log.Printf("testing deleteLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.deleteLogicalDevice(context.TODO(), i)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestGetLogicalDeviceByName(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing deleteLogicalDevice() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		logicalDevices, err := client.client.getAllLogicalDevices(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		for _, i := range samples(len(logicalDevices)) {
			test := logicalDevices[i]
			log.Printf("testing GetLogicalDeviceByName(%s) against %s %s (%s)", test.Data.DisplayName, client.clientType, clientName, client.client.ApiVersion())
			logicalDevice, err := client.client.GetLogicalDeviceByName(context.TODO(), test.Data.DisplayName)
			var ace ClientErr
			if err != nil {
				if !(errors.As(err, &ace) && ace.Type() == ErrMultipleMatch) {
					t.Fatal(err)
				}
				continue
			}
			if logicalDevice.Id != test.Id {
				t.Fatalf("expected '%s', got '%s'", test.Id, logicalDevice.Id)
			}
		}
	}
}

func TestLogicalDevicePortSpeed_IsEqual(t *testing.T) {
	for _, s := range testSpeedStrings {
		if !LogicalDevicePortSpeed(s[0]).IsEqual(LogicalDevicePortSpeed(s[1])) {
			t.Fatalf("speeds not equal %s %s", s[0], s[1])
		}
	}
}

func TestLogicalDevicePortRoleFlagsFromStrings(t *testing.T) {
	data := []string{"spine", "leaf"}
	expected := LogicalDevicePortRoleSpine | LogicalDevicePortRoleLeaf
	var result LogicalDevicePortRoleFlags
	err := result.FromStrings(data)
	if err != nil {
		t.Fatal(err)
	}
	if result != expected {
		t.Fatalf("expected %d, got %d", expected, result)
	}
}
