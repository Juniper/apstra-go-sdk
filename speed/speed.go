// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package speed

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Juniper/apstra-go-sdk/enum"
	oenum "github.com/orsinium-labs/enum"
)

const (
	m = 1_000_000
	g = 1_000_000_000
	// t = 1_000_000_000_000 // todo: if Tbps is introduced revisit todos in this file and in enum/enums.go
)

type speedValue oenum.Member[int]

var (
	speedValue0   = speedValue{Value: 0}
	speedValue1   = speedValue{Value: 1}
	speedValue2   = speedValue{Value: 2}
	speedValue5   = speedValue{Value: 5}
	speedValue10  = speedValue{Value: 10}
	speedValue25  = speedValue{Value: 25}
	speedValue40  = speedValue{Value: 40}
	speedValue50  = speedValue{Value: 50}
	speedValue100 = speedValue{Value: 100}
	speedValue150 = speedValue{Value: 150}
	speedValue200 = speedValue{Value: 200}
	speedValue400 = speedValue{Value: 400}
	speedValue800 = speedValue{Value: 800}
	speedValues   = oenum.New(speedValue0, speedValue1, speedValue2,
		speedValue5, speedValue10, speedValue25, speedValue40, speedValue50,
		speedValue100, speedValue150, speedValue200, speedValue400, speedValue800,
	)
)

var (
	_ json.Marshaler   = (*Speed)(nil)
	_ json.Unmarshaler = (*Speed)(nil)
)

// Speed is a case-insensitive string value representing interface speed in bps.
// Suffixes "M/m", and "G/g" are supported with an optional trailing "bps" or
// "b/s". For example, "1G", "1g", "1Gbps", "1Gb/s", "1000M" and "1000000000"
// all represent 1Gb/s.
type Speed string

func (s Speed) BitsPerSecond() int64 {
	if s == "" {
		return 0
	}

	// normalize the string a bit
	wc := string(s) // working copy of s
	wc = strings.ToLower(string(s))
	wc = strings.TrimSpace(wc)
	wc = strings.TrimSuffix(wc, "bps")
	wc = strings.TrimSuffix(wc, "b/s")
	wc = strings.TrimSpace(wc)

	multiplier := int64(1)
	// look for a letter indicating an SI unit
	switch {
	case strings.HasSuffix(wc, "m"):
		wc = strings.TrimSuffix(wc, "m")
		wc = strings.TrimSpace(wc)
		multiplier = m
	case strings.HasSuffix(wc, "g"):
		wc = strings.TrimSuffix(wc, "g")
		wc = strings.TrimSpace(wc)
		multiplier = g
		// case strings.HasSuffix(wc, "t"): // todo: if Tbps is introduced revisit todos in this file and in enum/enums.go
		//	wc = strings.TrimSuffix(wc, "t")
		//	wc = strings.TrimSpace(wc)
		//	multiplier = t
	}

	// convert the remainder of the string into a numeric value
	wcInt64, err := strconv.ParseInt(wc, 10, 64)
	if err != nil {
		return 0
	}

	return wcInt64 * multiplier
}

func (s Speed) Equal(other Speed) bool {
	return s.BitsPerSecond() == other.BitsPerSecond()
}

func (s Speed) MarshalJSON() ([]byte, error) {
	bps := s.BitsPerSecond()
	if bps == 0 {
		return []byte("null"), nil
	}

	var unit enum.SpeedUnit
	var value int
	switch {
	// case bps >= t: // at least 1Tbps // todo: if Tbps is introduced revisit todos in this file and in enum/enums.go
	//	if bps%t != 0 {
	//		return nil, fmt.Errorf("speed %q cannot be represented in Tbps", s)
	//	}
	//	unit = "T"
	//	value = int(bps / t)
	case bps >= g: // at least 1Gbps
		if bps%g != 0 {
			return nil, fmt.Errorf("speed %q cannot be represented in Gbps", s)
		}
		unit = enum.SpeedUnitG
		value = int(bps / g)
	default:
		if bps%m != 0 {
			return nil, fmt.Errorf("speed %q cannot be represented in Mbps", s)
		}
		unit = enum.SpeedUnitM
		value = int(bps / m)
	}

	// check the integer value (10, 25, 100, etc.) against the enum
	if speedValues.Parse(value) == nil {
		return nil, fmt.Errorf("speed %q is not supported", s)
	}

	return json.Marshal(&rawSpeed{Unit: unit, Value: &value})
}

func (s *Speed) UnmarshalJSON(bytes []byte) error {
	var raw rawSpeed
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		// This code handles Freeform link speeds. Currently handling freeform speeds as enums
		// because, while it's possible to unmarshal both styles, it's not obvious which format
		// to use when marshaling for the API without using different types. Leaving this code
		// around in case it's helpful in the future. /cmm September 2025
		//
		//var jute *json.UnmarshalTypeError
		//if errors.As(err, &jute) && jute.Value == "string" {
		//	var str string
		//	err = json.Unmarshal(bytes, &str)
		//	if err != nil {
		//		return fmt.Errorf("unmarshaling speed to string: %w", err)
		//	}
		//
		//	// Try parsing the API response as a FreeformLinkSpeed (string) type
		//	ffls := enum.FreeformLinkSpeeds.Parse(str)
		//	if ffls == nil {
		//		return fmt.Errorf("failed to parse speed %q", string(bytes))
		//	}
		//
		//	*s = Speed(strings.ToUpper(ffls.Value))
		//	return nil
		//}
		//
		//return fmt.Errorf("unmarshaling speed to struct: %w", err)
		return fmt.Errorf("unmarshaling speed: %w", err)
	}

	if raw.Value == nil {
		*s = ""
		return nil
	}

	*s = Speed(fmt.Sprintf("%d%s", *raw.Value, raw.Unit))
	return nil
}

type rawSpeed struct {
	Unit  enum.SpeedUnit `json:"unit"`
	Value *int           `json:"value"`
}
