// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build requiretestutils

package compare

import (
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/design"
	"github.com/stretchr/testify/require"
)

func ConfigletData(t testing.TB, a, b *apstra.ConfigletData) {
	t.Helper()

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.Equal(t, a.DisplayName, b.DisplayName)
	SlicesAsSets(t, a.RefArchs, b.RefArchs, "while comparing configlet refarchs,")
	require.Equal(t, len(a.Generators), len(b.Generators))
	for i := range a.Generators {
		ConfigletGenerators(t, a.Generators[i], b.Generators[i])
	}
}

func ConfigletGenerators(t testing.TB, a, b apstra.ConfigletGenerator) {
	t.Helper()

	require.Equal(t, a.ConfigStyle.String(), b.ConfigStyle.String())
	require.Equal(t, a.Section.String(), b.Section.String())
	require.Equal(t, a.SectionCondition, b.SectionCondition)
	require.Equal(t, a.TemplateText, b.TemplateText)
	require.Equal(t, a.NegationTemplateText, b.NegationTemplateText)
	require.Equal(t, a.Filename, b.Filename)
}

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
