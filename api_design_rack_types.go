package goapstra

import "time"

const (
	apiUrlDesignRackTypes = apiUrlDesignPrefix + "rack-types"
)

type rackTypeId string

type optionsRackTypeResponse struct {
	Items   []rackTypeId `json:"items"`
	Methods []string     `json:"methods"`
}

type rackType struct {
	Description              string     `json:"description"`
	Tags                     []string   `json:"tags"`
	Id                       rackTypeId `json:"id"`
	DisplayName              string     `json:"display_name"`
	FabricConnectivityDesign string     `json:"fabric_connectivity_design"`
	CreatedAt                time.Time  `json:"created_at"`
	LastModifiedAt           time.Time  `json:"last_modified_at"`
	LogicalDevices           []struct {
		Panels []struct {
			PanelLayout struct {
				RowCount    int `json:"row_count"`
				ColumnCount int `json:"column_count"`
			} `json:"panel_layout"`
			PortIndexing struct {
				Order      string `json:"order"`
				StartIndex int    `json:"start_index"`
				Schema     string `json:"schema"`
			} `json:"port_indexing"`
			PortGroups []struct {
				Count int `json:"count"`
				Speed struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				} `json:"speed"`
				Roles []string `json:"roles"`
			} `json:"port_groups"`
		} `json:"panels"`
		DisplayName string `json:"display_name"`
		Id          string `json:"id"`
	} `json:"logical_devices"`
	GenericSystems []struct {
		Count     int    `json:"count"`
		AsnDomain string `json:"asn_domain"`
		Links     []struct {
			Tags               []interface{} `json:"tags"`
			LinkPerSwitchCount int           `json:"link_per_switch_count"`
			Label              string        `json:"label"`
			LinkSpeed          struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			} `json:"link_speed"`
			TargetSwitchLabel string `json:"target_switch_label"`
			AttachmentType    string `json:"attachment_type"`
			LagMode           string `json:"lag_mode"`
		} `json:"links"`
		ManagementLevel  string        `json:"management_level"`
		PortChannelIdMin int           `json:"port_channel_id_min"`
		PortChannelIdMax int           `json:"port_channel_id_max"`
		LogicalDevice    string        `json:"logical_device"`
		Loopback         string        `json:"loopback"`
		Tags             []interface{} `json:"tags"`
		Label            string        `json:"label"`
	} `json:"generic_systems"`
	Servers []interface{} `json:"servers"`
	Leafs   []struct {
		LeafLeafL3LinkSpeed       interface{} `json:"leaf_leaf_l3_link_speed"`
		RedundancyProtocol        string      `json:"redundancy_protocol"`
		LeafLeafLinkPortChannelId int         `json:"leaf_leaf_link_port_channel_id"`
		LeafLeafL3LinkCount       int         `json:"leaf_leaf_l3_link_count"`
		LogicalDevice             string      `json:"logical_device"`
		LeafLeafLinkSpeed         struct {
			Unit  string `json:"unit"`
			Value int    `json:"value"`
		} `json:"leaf_leaf_link_speed"`
		LinkPerSpineCount int           `json:"link_per_spine_count"`
		LeafLeafLinkCount int           `json:"leaf_leaf_link_count"`
		Tags              []interface{} `json:"tags"`
		LinkPerSpineSpeed struct {
			Unit  string `json:"unit"`
			Value int    `json:"value"`
		} `json:"link_per_spine_speed"`
		Label                       string `json:"label"`
		MlagVlanId                  int    `json:"mlag_vlan_id"`
		LeafLeafL3LinkPortChannelId int    `json:"leaf_leaf_l3_link_port_channel_id"`
	} `json:"leafs"`
	AccessSwitches []interface{} `json:"access_switches"`
}

type createRackRequest struct {
	Description    string        `json:"description"`
	LastModifiedAt interface{}   `json:"last_modified_at"`
	Tags           []interface{} `json:"tags"`
	Leafs          []struct {
		LinkPerSpineCount   int         `json:"link_per_spine_count"`
		RedundancyProtocol  interface{} `json:"redundancy_protocol"`
		LeafLeafLinkSpeed   interface{} `json:"leaf_leaf_link_speed"`
		LeafLeafL3LinkCount int         `json:"leaf_leaf_l3_link_count"`
		LeafLeafL3LinkSpeed interface{} `json:"leaf_leaf_l3_link_speed"`
		LinkPerSpineSpeed   struct {
			Unit  string `json:"unit"`
			Value int    `json:"value"`
		} `json:"link_per_spine_speed"`
		Label                       string `json:"label"`
		LeafLeafL3LinkPortChannelId int    `json:"leaf_leaf_l3_link_port_channel_id"`
		LeafLeafLinkPortChannelId   int    `json:"leaf_leaf_link_port_channel_id"`
		LogicalDevice               string `json:"logical_device"`
		LeafLeafLinkCount           int    `json:"leaf_leaf_link_count"`
	} `json:"leafs"`
	LogicalDevices []struct {
		CreatedAt time.Time `json:"created_at"`
		Panels    []struct {
			PanelLayout struct {
				RowCount    int `json:"row_count"`
				ColumnCount int `json:"column_count"`
			} `json:"panel_layout"`
			PortIndexing struct {
				Order      string `json:"order"`
				StartIndex int    `json:"start_index"`
				Schema     string `json:"schema"`
			} `json:"port_indexing"`
			PortGroups []struct {
				Count int `json:"count"`
				Speed struct {
					Unit  string `json:"unit"`
					Value int    `json:"value"`
				} `json:"speed"`
				Roles []string `json:"roles"`
			} `json:"port_groups"`
		} `json:"panels"`
		DisplayName    string    `json:"display_name"`
		Id             string    `json:"id"`
		LastModifiedAt time.Time `json:"last_modified_at"`
		Href           string    `json:"href"`
	} `json:"logical_devices"`
	AccessSwitches           []interface{} `json:"access_switches"`
	FabricConnectivityDesign string        `json:"fabric_connectivity_design"`
	Id                       string        `json:"id"`
	GenericSystems           []interface{} `json:"generic_systems"`
	DisplayName              string        `json:"display_name"`
}
