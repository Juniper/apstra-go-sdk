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
	prs.canonicalize()
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

	result.canonicalize()
	err := result.validate()
	if err != nil {
		return err
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

	// make a clone so that we don't modify the underlying array when we
	// canonicalize it b/c we're supposed to be merely validating
	clone := make(PortRanges, len(prs))
	copy(clone, prs)

	clone.canonicalize()

	// validate the first entry in the clone
	err := clone[0].validate()
	if err != nil {
		return err
	}

	// validate each remaining entry and compare to the previous one
	for i, pr := range clone[1:] {
		err = pr.validate()
		if err != nil {
			return err
		}

		previous := clone[i] // i is indexed off clone[1:], so clone[i] refers to the one before pr

		if previous.Last >= pr.First {
			return fmt.Errorf("ranges must must not overlap")
		}
	}

	return nil
}
