package apstra

import (
	"fmt"
	"github.com/orsinium-labs/enum"
)

type FeatureSwitchEnum enum.Member[string]

func (o FeatureSwitchEnum) String() string {
	return o.Value
}

func (o *FeatureSwitchEnum) FromString(s string) error {
	t := FeatureSwitchEnums.Parse(s)
	if t == nil {
		return fmt.Errorf("failed to parse FeatureSwitchEnum %q", s)
	}
	o.Value = t.Value
	return nil
}

var (
	FeatureSwitchEnumEnabled  = FeatureSwitchEnum{Value: "enabled"}
	FeatureSwitchEnumDisabled = FeatureSwitchEnum{Value: "disabled"}
	FeatureSwitchEnums        = enum.New(FeatureSwitchEnumEnabled, FeatureSwitchEnumDisabled)
)
