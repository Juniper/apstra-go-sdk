// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package design

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash"
	"log"
	"reflect"
	"strings"

	"github.com/Juniper/apstra-go-sdk/internal/slice"
)

// hashForComparison returns a hash of the given value, excluding the top-level
// "id" field (case-insensitively). This is useful for generating consistent
// digests of structurally identical values that may differ only in their "id".
//
// The value `a` must be a struct.
func hashForComparison(a any, h hash.Hash) ([]byte, error) {
	// we expect a to be a struct (or struct pointer)
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %T", a)
	}

	// serialize a - this will drop timestamps because we never serialize those
	b, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("marshaling %T: %w", a, err)
	}

	// reconstitute serialized data without digging too deep
	var m map[string]json.RawMessage
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling %T: %w", a, err)
	}

	// drop id element(s?) (case insensitive)
	for k := range m {
		if strings.EqualFold(k, "id") {
			delete(m, k)
			break
		}
	}

	// reserialize and hash
	b, err = json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshaling %T: %w", a, err)
	}
	h.Reset()
	h.Write(b)
	return h.Sum(nil), err
}

// mustHashForComparison is like hashForComparison but panics if an error occurs.
// It should be used only when failure is unexpected or unrecoverable.
func mustHashForComparison(a any, h hash.Hash) []byte {
	result, err := hashForComparison(a, h)
	if err != nil {
		log.Panicf("hashing: %s", err.Error())
	}
	return result
}

func orderedMarshalJSON(keys []string, values map[string]json.RawMessage) ([]byte, error) {
	if len(keys) != len(values) {
		return nil, fmt.Errorf("orderedMarshalJSON: mismatch — %d keys but %d values", len(keys), len(values))
	}

	if !slice.IsUniq(keys) {
		return nil, fmt.Errorf("orderedMarshalJSON: keys contains duplicate elements")
	}

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyJSON, _ := json.Marshal(k) // for quoting/escaping - marshaling a string cannot error
		buf.Write(keyJSON)
		buf.WriteByte(':')
		v, ok := values[k]
		if !ok {
			return nil, fmt.Errorf("key %q has no value", k)
		}
		buf.Write(v)
	}
	buf.WriteByte('}')

	return buf.Bytes(), nil
}
