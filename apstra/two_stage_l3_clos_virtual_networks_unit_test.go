// Copyright (c) Juniper Networks, Inc., 2022-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"testing"
)

func TestTwoStageL3ClosVirtualNetworkStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		int() int
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
		{stringVal: "", intType: SystemRoleNone, stringType: systemRoleNone},
		{stringVal: "access", intType: SystemRoleAccess, stringType: systemRoleAccess},
		{stringVal: "leaf", intType: SystemRoleLeaf, stringType: systemRoleLeaf},
	}

	for i, td := range testData {
		ii := td.intType.int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatal(err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}
