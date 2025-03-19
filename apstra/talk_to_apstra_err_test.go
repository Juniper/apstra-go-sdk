package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApiErrors_UnmarshalJSON(t *testing.T) {
	type testCase struct {
		s string
		e []string
	}

	testCases := map[string]testCase{
		"empty": {s: "{}"},
		"single_error": {
			s: `{"errors":"test error"}`,
			e: []string{"test error"},
		},
		"multiple_error": {
			s: `{"errors":["test error one","test error two"]}`,
			e: []string{"test error one", "test error two"},
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			var r apiErrors
			err := json.Unmarshal([]byte(tCase.s), &r)
			require.NoError(t, err)
			require.Equal(t, tCase.e, r.Errors)
		})
	}
}
