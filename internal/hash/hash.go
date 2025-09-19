// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package hash

import (
	"encoding/json"
	"hash"
	"log"
)

func Struct(v any, h hash.Hash) ([]byte, error) {
	enc := json.NewEncoder(h)
	err := enc.Encode(v)
	return h.Sum(nil), err
}

func StructMust(v any, h hash.Hash) []byte {
	result, err := Struct(v, h)
	if err != nil {
		log.Panicf("hashing struct: %s", err.Error())
	}
	return result
}
