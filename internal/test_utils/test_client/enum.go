// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

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
	ClientTypeAPIOps    = ClientType{Value: "api-ops"}
	ClientTypeAWS       = ClientType{Value: "aws"}
	ClientTypeCloudLabs = ClientType{Value: "cloudlabs"}
	ClientTypeSlicer    = ClientType{Value: "slicer"}
	ClientTypes         = enum.New(
		ClientTypeAPIOps,
		ClientTypeAWS,
		ClientTypeCloudLabs,
		ClientTypeSlicer,
	)
)
