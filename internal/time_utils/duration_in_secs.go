// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package timeutils

import (
	"encoding/json"
	"time"
)

type DurationInSecs time.Duration

func (i DurationInSecs) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(time.Duration(i).Seconds()))
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
