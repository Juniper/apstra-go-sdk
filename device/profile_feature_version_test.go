// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package device

import (
	"reflect"
	"testing"

	sdk "github.com/Juniper/apstra-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestFeatureVersions_Validate(t *testing.T) {
	type testCase struct {
		v      FeatureVersions
		expErr any
	}

	testCases := map[string]testCase{
		"nil":      {},
		"empty":    {v: FeatureVersions{}},
		"single":   {v: FeatureVersions{{Version: "1.2.3"}}},
		"multiple": {v: FeatureVersions{{Version: "1.2.3"}, {Version: "4.5.6"}}},
		"collision": {
			v:      FeatureVersions{{Version: "1.2.3"}, {Version: "4.5.6"}, {Version: "1.2.3"}},
			expErr: new(sdk.ErrMultipleMatch),
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			err := tCase.v.Validate()

			if tCase.expErr != nil {
				target := reflect.New(reflect.TypeOf(tCase.expErr).Elem()).Interface()
				require.ErrorAs(t, err, target)
				return
			}

			require.NoError(t, err)
		})
	}
}
