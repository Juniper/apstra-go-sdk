// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package testclient

import (
	"net/netip"

	"github.com/Juniper/apstra-go-sdk/apstra"
)

type SwitchInfo struct {
	ManagementIP netip.Addr
	Username     string
	Password     string
	Platform     apstra.AgentPlatform
}
