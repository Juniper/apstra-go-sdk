// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package comparedesign

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/stretchr/testify/require"
)

func Configlet(t testing.TB, req, resp design.Configlet, msg ...string) {
	msg = addMsg(msg, "Comparing Configlet")

	require.Equal(t, req.Label, resp.Label, msg)
	require.Equal(t, len(req.Generators), len(resp.Generators), msg)
	for i := range len(req.Generators) {
		ConfigletGenerator(t, req.Generators[i], resp.Generators[i], addMsg(msg, "Comparing Configlet Generator %d", i)...)
	}
}

func ConfigletGenerator(t testing.TB, req, resp design.ConfigletGenerator, msg ...string) {
	msg = addMsg(msg, "Comparing Configlet Generator")

	require.Equal(t, req.ConfigStyle, resp.ConfigStyle, msg)
	require.Equal(t, req.Section, resp.Section, msg)
	require.Equal(t, req.SectionCondition, resp.SectionCondition, msg)
	require.Equal(t, req.TemplateText, resp.TemplateText, msg)
	require.Equal(t, req.NegationTemplateText, resp.NegationTemplateText, msg)
	require.Equal(t, req.Filename, resp.Filename, msg)
}
