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

func (i *DurationInSecs) MarshalJSON() ([]byte, error) {
	secs := 1
	if i != nil {
		secs = int((*time.Duration)(i).Seconds())
	}
	return json.Marshal(secs)
}

func (o *DurationInSecs) UnmarshalJSON(bytes []byte) error {
	var i int
	if len(bytes) == 0 {
		return nil
	}
	err := json.Unmarshal(bytes, &i)
	if err != nil {
		return err
	}
	*o = DurationInSecs(time.Duration(i) * time.Second)
	return nil
}
