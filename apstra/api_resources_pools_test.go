// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "testing"

func TestPoolStatusStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "in_use", intType: PoolStatusInUse, stringType: poolStatusInUse},
		{stringVal: "not_in_use", intType: PoolStatusUnused, stringType: poolStatusUnused},
		{stringVal: "creation_in_progress", intType: PoolStatusCreating, stringType: poolStatusCreating},
		{stringVal: "marked_for_deletion", intType: PoolStatusDeleting, stringType: poolStatusDeleting},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}
