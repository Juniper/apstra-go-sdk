// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"fmt"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

func TestRaResourceValidate(t *testing.T) {
	type testCase struct {
		data   FreeformRaResourceData
		expErr bool
	}
	testCases := []testCase{
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeInt,
				Label:           randString(6, "hex"),
				Value:           nil,
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: nil,
				GeneratorId:     nil,
			},
			expErr: false,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeInt,
				Label:           randString(6, "hex"),
				Value:           toPtr("1"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: nil,
				GeneratorId:     nil,
			},
			expErr: false,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeInt,
				Label:           randString(6, "hex"),
				Value:           toPtr("foo"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: nil,
				GeneratorId:     nil,
			},
			expErr: true,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeHostIpv4,
				Label:           randString(6, "hex"),
				Value:           toPtr("192.168.2.1/24"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: toPtr(24),
				GeneratorId:     nil,
			},
			expErr: true,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeIpv4,
				Label:           randString(6, "hex"),
				Value:           toPtr("192.168.2.0/24"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: toPtr(24),
				GeneratorId:     nil,
			},
			expErr: false,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeIpv4,
				Label:           randString(6, "hex"),
				Value:           toPtr("192.168.2.1/24"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: nil,
				GeneratorId:     nil,
			},
			expErr: true,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeIpv4,
				Label:           randString(6, "hex"),
				Value:           toPtr("2001:db8:3333:4444:5555:6666:7777:8888/64"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: nil,
				GeneratorId:     nil,
			},
			expErr: true,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeIpv6,
				Label:           randString(6, "hex"),
				Value:           toPtr("2001:db8:abcd:0012::0/64"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: toPtr(64),
				GeneratorId:     nil,
			},
			expErr: false,
		},
		{
			data: FreeformRaResourceData{
				ResourceType:    enum.FFResourceTypeVlan,
				Label:           randString(6, "hex"),
				Value:           toPtr("blue"),
				AllocatedFrom:   nil,
				GroupId:         "",
				SubnetPrefixLen: nil,
				GeneratorId:     nil,
			},
			expErr: true,
		},
	}

	for i, tc := range testCases {
		i, tc := i, tc
		t.Run(fmt.Sprintf("test testcase %d", i), func(t *testing.T) {
			t.Parallel()
			err := tc.data.validate()
			if (err != nil) != tc.expErr {
				t.Fatalf("test %d expected error %t got error %t ", i, tc.expErr, err != nil)
			}
		})
	}
}
