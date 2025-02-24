// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "testing"

func TestLockStatusStrings(t *testing.T) {
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
		{stringVal: "locked", intType: LockStatusLocked, stringType: lockStatusLocked},
		{stringVal: "unlocked", intType: LockStatusUnlocked, stringType: lockStatusUnlocked},
		{stringVal: "locked_by_restricted_user", intType: LockStatusLockedByRestrictedUser, stringType: lockStatusLockedByRestrictedUser},
		{stringVal: "locked_by_admin", intType: LockStatusLockedByAdmin, stringType: lockStatusLockedByAdmin},
		{stringVal: "locked_by_deleted_user", intType: LockStatusLockedByDeletedUser, stringType: lockStatusLockedByDeletedUser},
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
