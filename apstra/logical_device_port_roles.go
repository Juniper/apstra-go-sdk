// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"sort"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
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

func (o *LogicalDevicePortRoles) SetAll() {
	*o = enum.PortRoles.Members()
}

func (o *LogicalDevicePortRoles) Sort() {
	if o == nil || *o == nil {
		return
	}

	clone := *o

	sort.Slice(*o, func(i, j int) bool {
		return clone[i].Value < clone[j].Value
	})
}
