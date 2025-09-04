// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"encoding/json"
	"fmt"
	"io"
)

func PrettyPrint(in interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")

	err := enc.Encode(in)
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}

	return nil
}
