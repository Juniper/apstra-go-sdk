// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
)

var (
	_ encoding.TextMarshaler   = (*DHCPServiceEnabled)(nil)
	_ encoding.TextUnmarshaler = (*DHCPServiceEnabled)(nil)
)

type DHCPServiceEnabled bool

func (o DHCPServiceEnabled) MarshalText() ([]byte, error) {
	if o {
		return []byte(enum.DhcpServiceModeEnabled.String()), nil
	}

	return []byte(enum.DhcpServiceModeDisabled.String()), nil
}

func (o *DHCPServiceEnabled) UnmarshalText(bytes []byte) error {
	var dsm enum.DhcpServiceMode

	err := dsm.FromString(string(bytes))
	if err != nil {
		return fmt.Errorf("while parsing dhcp service mode - %w", err)
	}

	*o = dsm == enum.DhcpServiceModeEnabled

	return nil
}

//func (o *DHCPServiceEnabled) FromString(s string) error {
//	var dsm enum.DhcpServiceMode
//
//	err := dsm.FromString(s)
//	if err != nil {
//		return fmt.Errorf("while parsing dhcp service mode - %w", err)
//	}
//
//	*o = dsm == enum.DhcpServiceModeEnabled
//
//	return nil
//}
//
//func (o DHCPServiceEnabled) MarshalJSON() ([]byte, error) {
//	return json.Marshal(o.String())
//}
//
//func (o DHCPServiceEnabled) String() string {
//	if o {
//		return enum.DhcpServiceModeEnabled.String()
//	}
//
//	return enum.DhcpServiceModeDisabled.String()
//}
