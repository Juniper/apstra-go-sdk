// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
)

type LogicalDevicePortRoles []enum.PortRole

func (o *LogicalDevicePortRoles) Strings() []string {
	if *o == nil {
		return nil
	}

	result := make([]string, len(*o))
	for i, pr := range *o {
		result[i] = pr.String()
	}

	return result
}

func (o *LogicalDevicePortRoles) FromStrings(in []string) error {
	newPRs := make(LogicalDevicePortRoles, len(in))
	for i, s := range in {
		err := newPRs[i].FromString(s)
		if err != nil {
			return err
		}
	}
	*o = newPRs

	return nil
}

// IncludeAllUses ensures that the LogicalDevicePortRoles contains the entire
// set of "available for use" port roles: All roles excluding "l3_server"
// (deprecated) and "unused".
func (o *LogicalDevicePortRoles) IncludeAllUses() {
	// wipe out any existing values
	*o = nil

	for _, member := range enum.PortRoles.Members() {
		switch member {
		case enum.PortRoleL3Server: // don't add this one
		default:
			*o = append(*o, member) // this one's a keeper
		}
	}
}

func (o *LogicalDevicePortRoles) Validate() error {
	if o == nil {
		return nil
	}

	for _, ldpr := range *o {
		if ldpr == enum.PortRoleL3Server {
			return fmt.Errorf("logical device port role %q is no longer supported", ldpr.String())
		}
	}

	return nil
}
