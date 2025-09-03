// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testclient

import (
	"github.com/Juniper/apstra-go-sdk/apstra"
	"net/netip"
)

type SwitchInfo struct {
	ManagementIP netip.Addr
	Username     string
	Password     string
	Platform     apstra.AgentPlatform
}
