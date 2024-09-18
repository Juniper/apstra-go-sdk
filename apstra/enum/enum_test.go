package enum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnumParsingError(t *testing.T) {
	type testCase struct {
		enum             enum
		stringVal        string
		expEnumErr       bool
		expErrorString   string
		expEnumErrorType ErrorType
	}

	testCases := map[string]testCase{
		"valid_FeatureSwitchEnum": {
			enum:           new(FeatureSwitch),
			stringVal:      "enabled",
			expErrorString: "",
		},
		"invalid_FeatureSwitchEnum": {
			enum:             new(FeatureSwitch),
			stringVal:        "bogus_switch_status",
			expEnumErr:       true,
			expErrorString:   `failed to parse *enum.FeatureSwitch "bogus_switch_status"`,
			expEnumErrorType: ErrorTypeParsingFailed,
		},
		"valid_ApiFeature": {
			enum:           new(ApiFeature),
			stringVal:      "freeform",
			expErrorString: "",
		},
		"invalid_ApiFeature": {
			enum:             new(ApiFeature),
			stringVal:        "bogus_feature",
			expEnumErr:       true,
			expErrorString:   `failed to parse *enum.ApiFeature "bogus_feature"`,
			expEnumErrorType: ErrorTypeParsingFailed,
		},
	}

	for tName, tCase := range testCases {
		t.Run(tName, func(t *testing.T) {
			t.Parallel()

			err := tCase.enum.FromString(tCase.stringVal)
			if tCase.expEnumErr {
				require.Error(t, err)
				var eErr Error
				require.ErrorAs(t, err, &eErr)
			} else {
				require.NoError(t, err)
			}
			if tCase.expErrorString != "" {
				require.Error(t, err)
				require.Equal(t, tCase.expErrorString, err.Error())
			} else {
				require.NoError(t, err)
			}
			if tCase.expEnumErrorType != errorTypeUnknown {
				require.Error(t, err)
				var eErr Error
				require.ErrorAs(t, err, &eErr)
				require.Equal(t, tCase.expEnumErrorType, eErr.errType)
			}

			// no error cases require string match
			if !tCase.expEnumErr && tCase.expErrorString == "" && tCase.expEnumErrorType == errorTypeUnknown {
				require.NoError(t, err)
				require.Equal(t, tCase.stringVal, tCase.enum.String())
			}
		})
	}
}
