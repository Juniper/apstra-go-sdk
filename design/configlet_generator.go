// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import "github.com/Juniper/apstra-go-sdk/enum"

type ConfigletGenerator struct {
	ConfigStyle          enum.ConfigletStyle   `json:"config_style"`
	Section              enum.ConfigletSection `json:"section"`
	SectionCondition     *string               `json:"section_condition,omitempty"`
	TemplateText         string                `json:"template_text"`
	NegationTemplateText string                `json:"negation_template_text"`
	Filename             string                `json:"filename"`
}
