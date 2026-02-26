// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pointer

import (
	"encoding/json"

	"github.com/Juniper/apstra-go-sdk/internal/str"
)

func StringMarshalJSONWithEmptyAsNull(p *string) json.RawMessage {
	if p == nil {
		return nil
	}

	if *p == "" {
		return json.RawMessage("null")
	}

	return str.QuoteJSONString(*p)
}
