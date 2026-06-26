// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package datacenter

import (
	"bytes"
	"encoding"
	"fmt"
	"sort"
	"strings"
)

const (
	portAny       = "any"
	portRangesSep = ","
)

var (
	_ encoding.TextMarshaler   = PortRanges(nil)
	_ encoding.TextUnmarshaler = (*PortRanges)(nil)
)

type PortRanges []PortRange

func (prs PortRanges) MarshalText() ([]byte, error) {
	err := prs.validate()
	if err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return []byte(portAny), nil
	}

	var buf bytes.Buffer
	b, err := prs[0].MarshalText()
	if err != nil {
		return nil, fmt.Errorf("marshaling PortRange at index 0: %w", err)
	}
	_, _ = buf.Write(b)

	for i, pr := range prs[1:] {
		b, err = pr.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("marshaling PortRange at index %d: %w", 1+i, err)
		}

		buf.WriteString(portRangesSep)
		_, _ = buf.Write(b)
	}

	return buf.Bytes(), nil
}

func (prs *PortRanges) UnmarshalText(in []byte) error {
	if prs == nil {
		return fmt.Errorf("cannot unmarshal into nil *PortRanges")
	}

	if string(strings.TrimSpace(string(in))) == portAny {
		*prs = nil
		return nil
	}

	parts := strings.Split(string(in), portRangesSep)
	result := make(PortRanges, 0, len(parts))

	for i, part := range parts {
		var pr PortRange
		err := pr.UnmarshalText([]byte(part))
		if err != nil {
			return fmt.Errorf("unmarshaling PortRange at index %d: %w", i, err)
		}

		result = append(result, pr)
	}

	err := result.validate()
	if err != nil {
		return fmt.Errorf("unmarshaling PortRanges %q: %w", string(in), err)
	}

	*prs = result
	return nil
}

func (prs PortRanges) canonicalize() {
	for i := range prs {
		prs[i].canonicalize()
	}

	sort.Slice(prs, func(i, j int) bool { return prs[i].First < prs[j].First })
}

func (prs PortRanges) validate() error {
	if len(prs) == 0 {
		return nil
	}

	// validate the first entry
	err := prs[0].validate()
	if err != nil {
		return fmt.Errorf("validating PortRange at index 0: %w", err)
	}

	// validate and compare to the previous each entry following the first one
	for i, pr := range prs[1:] {
		err = pr.validate()
		if err != nil {
			return fmt.Errorf("validating PortRange at index %d: %w", 1+i, err)
		}

		previous := prs[i] // i is indexed off prs[1:], so prs[i] refers to the one before pr

		if previous.Last >= pr.First {
			return fmt.Errorf("PortRanges must be sorted and must not overlap")
		}
	}

	return nil
}
