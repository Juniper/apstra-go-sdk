// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package str

import (
	"encoding/json"
	"fmt"
)

// QuoteJSONString returns a JSON-encoded string literal as []byte.
// Example: QuoteJSONString("hi") -> []byte("\"hi\"")
func QuoteJSONString(s string) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		panic(fmt.Sprintf("marshaling a string shouldn't fail, but %q failed with %v", s, err))
	}
	return b
}
