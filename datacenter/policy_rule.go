// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding/json"
	"fmt"

	"github.com/Juniper/apstra-go-sdk/enum"
	"github.com/Juniper/apstra-go-sdk/errors"
	"github.com/Juniper/apstra-go-sdk/internal"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

var (
	_ internal.IDer     = (*PolicyRule)(nil)
	_ internal.IDSetter = (*PolicyRule)(nil)
	_ json.Unmarshaler  = (*PolicyRule)(nil)
)

type PolicyRule struct {
	Label             string                  `json:"label"`
	Description       string                  `json:"description"`
	Protocol          enum.PolicyRuleProtocol `json:"protocol"`
	Action            enum.PolicyRuleAction   `json:"action"`
	SrcPort           PortRanges              `json:"src_port"`
	DstPort           PortRanges              `json:"dst_port"`
	TcpStateQualifier *enum.TcpStateQualifier `json:"tcp_state_qualifier,omitempty"`

	id string
}

func (pr PolicyRule) ID() *string {
	if pr.id == "" {
		return nil
	}
	return pointer.ToCopyOf(pr.id)
}

func (pr *PolicyRule) SetID(id string) error {
	if pr == nil {
		return fmt.Errorf("cannot set ID of nil %T", pr)
	}

	if pr.id != "" {
		return errors.IDAlreadySet(fmt.Sprintf("id already has value %q", pr.id))
	}

	pr.id = id
	return nil
}

func (pr *PolicyRule) UnmarshalJSON(bytes []byte) error {
	if pr == nil {
		return fmt.Errorf("cannot unmarshal into nil %T", pr)
	}

	type policyRule PolicyRule // type alias prevents recursion

	var target struct {
		policyRule
		ID string `json:"id"` // the `id` struct element cannot be unmarshaled so we temporarily stash that value here
	}

	// unmarshal everything which can be handled by the `json` package.
	if err := json.Unmarshal(bytes, &target); err != nil {
		return err
	}

	*pr = PolicyRule(target.policyRule)
	pr.id = target.ID

	return nil
}
