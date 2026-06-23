// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

func (pr *PortRange) Canonicalize() {
	pr.canonicalize()
}

func (prs *PortRanges) Canonicalize() {
	prs.canonicalize()
}
