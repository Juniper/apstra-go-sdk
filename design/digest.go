// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"encoding/json"
	"fmt"
	"hash"
	"log"
)

func digestSkipID(a any, h hash.Hash) ([]byte, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("marshaling %T: %w", a, err)
	}

	var m map[string]json.RawMessage
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling %T: %w", a, err)
	}

	delete(m, "id")
	delete(m, "iD")
	delete(m, "Id")
	delete(m, "ID")

	b, err = json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshaling %T: %w", a, err)
	}
	h.Write(b)
	return h.Sum(nil), err
}

func mustDigestSkipID(a any, h hash.Hash) []byte {
	result, err := digestSkipID(a, h)
	if err != nil {
		log.Panicf("hashing: %s", err.Error())
	}
	return result
}
