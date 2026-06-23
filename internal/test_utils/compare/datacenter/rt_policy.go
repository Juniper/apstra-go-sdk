// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedatacenter

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/datacenter"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/compare"
	testmessage "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
)

func RTPolicy(t testing.TB, req, resp datacenter.RTPolicy, msg ...string) {
	msg = testmessage.Add(msg, "Comparing RT Policy")

	compare.SlicesAsSets(t, req.ImportRTs, resp.ImportRTs, "RTPolicy ImportRTs")
	compare.SlicesAsSets(t, req.ExportRTs, resp.ExportRTs, "RTPolicy ExportRTs")
}
