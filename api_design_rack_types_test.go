package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestListRackTypeIds(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	rackTypes, err := client.listRackTypeIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	log.Println(rackTypes)
}

func TestRackTypeStrings(t *testing.T) {
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
		{stringVal: "", intType: LeafRedundancyProtocolNone, stringType: leafRedundancyProtocolNone},
		{stringVal: "mlag", intType: LeafRedundancyProtocolMlag, stringType: leafRedundancyProtocolMlag},
		{stringVal: "esi", intType: LeafRedundancyProtocolEsi, stringType: leafRedundancyProtocolEsi},

		{stringVal: "", intType: AccessRedundancyProtocolNone, stringType: accessRedundancyProtocolNone},
		{stringVal: "esi", intType: AccessRedundancyProtocolEsi, stringType: accessRedundancyProtocolEsi},

		{stringVal: "l3clos", intType: FabricCOnnectivityDesignL3Clos, stringType: fabricConnectivityDesignL3Clos},
		{stringVal: "l3collapsed", intType: FabricCOnnectivityDesignL3Collapsed, stringType: fabricConnectivityDesignL3Collapsed},
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

func TestGetRackType(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	rackTypes, err := client.listRackTypeIds(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(rackTypes) == 0 {
		t.Skip("no rack types to fetch")
	}

	rt, err := client.getRackType(context.TODO(), rackTypes[0])
	if err != nil {
		t.Fatal(err)
	}

	log.Println(rt)
}

func TestCreateGetRackType(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	id, err := client.createRackType(context.TODO(), &RackType{
		DisplayName:              "name-" + randString(10, "hex"),
		Description:              "description " + randString(10, "hex"),
		FabricConnectivityDesign: FabricCOnnectivityDesignL3Clos,
		LeafSwitches: []RackElementLeafSwitch{
			{
				Label:                       "label-" + randString(10, "hex"),
				LeafLeafL3LinkCount:         0,
				LeafLeafL3LinkPortChannelId: 0,
				LeafLeafL3LinkSpeed:         nil,
				LeafLeafLinkCount:           0,
				LeafLeafLinkPortChannelId:   0,
				LeafLeafLinkSpeed:           nil,
				LinkPerSpineCount:           2,
				LinkPerSpineSpeed: &LogicalDevicePortSpeed{
					Unit:  "G",
					Value: 10,
				},
				MlagVlanId:         0,
				RedundancyProtocol: LeafRedundancyProtocolNone,
				Tags:               nil,
				Panels: []LogicalDevicePanel{
					{
						PanelLayout: LogicalDevicePanelLayout{
							RowCount:    3,
							ColumnCount: 3,
						},
						PortIndexing: LogicalDevicePortIndexing{
							Order:      PortIndexingHorizontalFirst,
							StartIndex: 1,
							Schema:     PortIndexingSchemaAbsolute,
						},
						PortGroups: []LogicalDevicePortGroup{
							{
								Count: 9,
								Speed: LogicalDevicePortSpeed{
									Unit:  "G",
									Value: 10,
								},
								Roles: []string{"access", "generic", "spine", "peer"},
							},
						},
					},
				},
				DisplayName: "logical device display name " + randString(10, "hex"),
			},
		},
		GenericSystems: []RackElementGenericSystem{
			{
				Count:            10,
				AsnDomain:        FeatureSwitchEnabled,
				ManagementLevel:  GenericSystemUnmanaged,
				PortChannelIdMin: 0,
				PortChannelIdMax: 0,
				Loopback:         FeatureSwitchDisabled,
				Tags:             nil,
				Label:            "some generic system",
				Links: []GenericSystemAccessLink{
					{
						Label:              "",
						Tags:               nil,
						LinkPerSwitchCount: 0,
						LinkSpeed:          LogicalDevicePortSpeed{},
						TargetSwitchLabel:  "",
						AttachmentType:     "",
						LagMode:            "",
					},
				},
				Panels: []LogicalDevicePanel{{
					PanelLayout:  LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 2},
					PortIndexing: LogicalDevicePortIndexing{Order: PortIndexingVerticalFirst, Schema: PortIndexingSchemaAbsolute},
					PortGroups: []LogicalDevicePortGroup{{
						Count: 2,
						Speed: LogicalDevicePortSpeed{Unit: "G", Value: 10},
						Roles: []string{"generic"},
					}},
				}},
				DisplayName: "Generic System Display Name",
			},
		},
		AccessSwitches: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("id: '%s'\n", id)
}
