// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package testclient

import (
	"fmt"
	"github.com/orsinium-labs/enum"
)

var _ fmt.Stringer = (*ClientType)(nil)

type ClientType enum.Member[string]

func (c ClientType) String() string {
	return c.Value
}

var (
	ClientTypeAPIOps    = ClientType{"api-ops"}
	ClientTypeAWS       = ClientType{"aws"}
	ClientTypeCloudLabs = ClientType{"cloudlabs"}
	ClientTypeSlicer    = ClientType{"slicer"}
	ClientTypes         = enum.New(
		ClientTypeAPIOps,
		ClientTypeAWS,
		ClientTypeCloudLabs,
		ClientTypeSlicer,
	)
)
