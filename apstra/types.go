// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	vlanMin = 1
	vlanMax = 4094

	vniMin = 4096
	vniMax = 16777214
)

type Vlan uint16

func (o Vlan) validate() error {
	if o < vlanMin || o > vlanMax {
		return fmt.Errorf("VLAN %d out of range", o)
	}
	return nil
}

type VNI uint32

type RtPolicy struct {
	ImportRTs []string `json:"import_RTs"`
	ExportRTs []string `json:"export_RTs"`
}

type DurationInSecs time.Duration

func (i DurationInSecs) MarshalJSON() ([]byte, error) {
	return json.Marshal(int((time.Duration)(i).Seconds()))
}

func (o *DurationInSecs) UnmarshalJSON(bytes []byte) error {
	var i int
	err := json.Unmarshal(bytes, &i)
	if err != nil {
		return err
	}
	*o = DurationInSecs(time.Duration(i) * time.Second)
	return nil
}

func (o *DurationInSecs) TimeinSecs() int {
	return int(time.Duration(*o).Seconds())
}

func NewDurationInSecs(timeinsecs int) *DurationInSecs {
	t := DurationInSecs(time.Duration(timeinsecs) * time.Second)
	return &t
}
