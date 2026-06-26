// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"encoding"
	"fmt"
	"strconv"
	"strings"
)

const portRangeSep = "-"

var (
	_ encoding.TextMarshaler   = (*PortRange)(nil)
	_ encoding.TextUnmarshaler = (*PortRange)(nil)
)

type PortRange struct {
	First uint16
	Last  uint16
}

func (pr PortRange) MarshalText() ([]byte, error) {
	pr.canonicalize()
	err := pr.validate()
	if err != nil {
		return nil, err
	}

	switch {
	case pr.First == pr.Last:
		return []byte(strconv.Itoa(int(pr.First))), nil
	case pr.First < pr.Last:
		return []byte(strconv.Itoa(int(pr.First)) + portRangeSep + strconv.Itoa(int(pr.Last))), nil
	}

	return nil, fmt.Errorf("unhandled port range not caught by validate function: %d - %d", pr.First, pr.Last)
}

func (pr *PortRange) UnmarshalText(in []byte) error {
	if pr == nil {
		return fmt.Errorf("cannot unmarshal into nil *PortRange")
	}

	var result PortRange
	parts := strings.Split(string(in), portRangeSep)
	switch len(parts) {
	case 1: // a port range may contain only a single string value like "80"
		first, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 16)
		if err != nil {
			return fmt.Errorf("invalid port range %q: %w", string(in), err)
		}
		result.First = uint16(first)
		result.Last = uint16(first)
	case 2: // a port range may contain a range of values like "20-21"
		first, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 16)
		if err != nil {
			return fmt.Errorf("invalid port range %q: %w", string(in), err)
		}
		last, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 16)
		if err != nil {
			return fmt.Errorf("invalid port range %q: %w", string(in), err)
		}
		result.First = uint16(first)
		result.Last = uint16(last)
	default:
		return fmt.Errorf("invalid port range %q", string(in))
	}

	result.canonicalize()
	err := result.validate()
	if err != nil {
		return err
	}

	*pr = result
	return nil
}

func (pr *PortRange) canonicalize() {
	if pr.First > pr.Last {
		pr.First, pr.Last = pr.Last, pr.First
	}
}

func (pr PortRange) validate() error {
	switch {
	case pr.First > pr.Last:
		return fmt.Errorf("invalid port range: first %d > last %d", pr.First, pr.Last)
	case pr.First == 0 || pr.Last == 0:
		return fmt.Errorf("port range %d-%d contains zero", pr.First, pr.Last)
	}
	return nil
}
