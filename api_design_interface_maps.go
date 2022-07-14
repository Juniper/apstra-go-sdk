package goapstra

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	apiUrlDesignInterfaceMaps       = apiUrlDesignPrefix + "interface-maps"
	apiUrlDesignInterfaceMapsPrefix = apiUrlDesignInterfaceMaps + apiUrlPathDelim
	apiUrlDesignInterfaceMapById    = apiUrlDesignInterfaceMapsPrefix + "%s"
)

// rawInterface.Setting.Param is a string containing JSON like this.
// it needs double quotes escaped {\"like\": \"this\"}.
// {
//  "global": {
//    "breakout": false,
//    "fpc": 0,
//    "pic": 0,
//    "port": 0,
//    "speed": "100g"
//  },
//  "interface": {
//    "speed": ""
//  }
//}

type InterfaceSettingParam struct {
	Global struct {
		Breakout bool   `json:"breakout"`
		Fpc      int    `json:"fpc"`
		Pic      int    `json:"pic"`
		Port     int    `json:"port"`
		Speed    string `json:"speed"`
	} `json:"global"`
	Interface struct {
		Speed string `json:"speed"`
	} `json:"interface"`
}

func (o InterfaceSettingParam) String() string {
	// medium confident we won't provoke UnsupportedTypeError or
	// UnsupportedValueError here, so ignoring the error return.
	payload, _ := json.Marshal(o)
	return strings.Replace(string(payload), `"`, `\"`, -1)
}

type rawInterface struct {
	Name    string   `json:"name"`
	Roles   []string `json:"roles"`
	Mapping []int    `json:"mapping"`
	State   string   `json:"state"`
	Setting struct {
		Param string `json:"param"`
	} `json:"setting"`
	Position int                    `json:"position"`
	Speed    LogicalDevicePortSpeed `json:"speed"`
}

type rawInterfaceMap struct {
	LogicalDeviceId ObjectId       `json:"logical_device_id"`
	DeviceProfileId ObjectId       `json:"device_profile_id"`
	CreatedAt       time.Time      `json:"created_at"`
	LastModifiedAt  time.Time      `json:"last_modified_at"`
	Id              ObjectId       `json:"id"`
	Label           string         `json:"label"`
	Interfaces      []rawInterface `json:"interfaces"`
}
