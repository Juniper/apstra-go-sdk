// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

func (o *TwoStageL3ClosClient) DropSysRedundancyCache() {
	o.sysRedundancyCache.mu.Lock()
	o.sysRedundancyCache.systemToGroup = make(map[string]*string)
	o.sysRedundancyCache.groupToSystems = make(map[string][2]string)
	o.sysRedundancyCache.mu.Unlock()
}

func (o *TwoStageL3ClosClient) CountSystems() int {
	o.sysRedundancyCache.mu.RLock()
	result := len(o.sysRedundancyCache.systemToGroup)
	o.sysRedundancyCache.mu.RUnlock()

	return result
}

func (o *TwoStageL3ClosClient) CountGroups() int {
	o.sysRedundancyCache.mu.RLock()
	result := len(o.sysRedundancyCache.groupToSystems)
	o.sysRedundancyCache.mu.RUnlock()

	return result
}
