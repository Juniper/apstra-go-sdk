// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func PodWithCount(t testing.TB, req, resp design.PodWithCount, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Pod With Count")

	require.Equal(t, req.Count, resp.Count, msg)
	TemplateRackBased(t, req.Pod, resp.Pod, msg...)
}
