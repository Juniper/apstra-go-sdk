// Copyright (c) Juniper Networks, Inc., 2025-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	comparepolicy "github.com/Juniper/apstra-go-sdk/internal/test_utils/compare/policy"
	"github.com/Juniper/apstra-go-sdk/internal/test_utils/test_message"
	"github.com/stretchr/testify/require"
)

func TemplatePodBased(t testing.TB, req, resp design.TemplatePodBased, msg ...string) {
	msg = testmessage.Add(msg, "Comparing Pod Based Template")

	require.Equal(t, req.Label, resp.Label, msg)

	if req.AntiAffinityPolicy != nil {
		require.NotNil(t, req.AntiAffinityPolicy)
		comparepolicy.AntiAffinity(t, *req.AntiAffinityPolicy, *resp.AntiAffinityPolicy, msg...)
	}

	Superspine(t, req.Superspine, resp.Superspine, msg...)

	require.Equal(t, len(req.Pods), len(resp.Pods), msg)
	for i := range len(req.Pods) {
		PodWithCount(t, req.Pods[i], resp.Pods[i], testmessage.Add(msg, "Comparing Pod %d", i)...)
	}

	if req.ID() != nil && resp.ID() != nil {
		require.Equal(t, req.ID(), resp.ID(), msg)
	}
	if req.CreatedAt() != nil && resp.CreatedAt() != nil {
		require.Equal(t, req.CreatedAt(), resp.CreatedAt(), msg)
	}
	if req.LastModifiedAt() != nil && resp.LastModifiedAt() != nil {
		require.Equal(t, req.LastModifiedAt(), resp.LastModifiedAt(), msg)
	}
}
