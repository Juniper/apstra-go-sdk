package goapstra

import (
	"testing"
)

func TestInterfaceSettingParam(t *testing.T) {
	expected := `{\"global\":{\"breakout\":false,\"fpc\":0,\"pic\":0,\"port\":0,\"speed\":\"100g\"},\"interface\":{\"speed\":\"\"}}`
	test := InterfaceSettingParam{
		Global: struct {
			Breakout bool   `json:"breakout"`
			Fpc      int    `json:"fpc"`
			Pic      int    `json:"pic"`
			Port     int    `json:"port"`
			Speed    string `json:"speed"`
		}{
			Breakout: false,
			Fpc:      0,
			Pic:      0,
			Port:     0,
			Speed:    "100g",
		},
		Interface: struct {
			Speed string `json:"speed"`
		}{},
	}
	result := test.String()
	if result != expected {
		t.Fatalf("expected '%s', got '%s'", expected, result)
	}
}
