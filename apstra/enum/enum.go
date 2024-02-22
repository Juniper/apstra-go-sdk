package enum

import (
	ole "github.com/orsinium-labs/enum"
)

const (
	ValueStateUnknown = valueState(iota)
	ValueStateNull
	ValueStateKnown

	Unknown = Type(iota)
	FeatureSwitchx
)

type valueState uint8

type Type int

type enum struct {
	enumType Type
	state    valueState
	value    ole.Member[string]
}

func (o enum) String() string {
	return o.value.Value
}

func (o enum) IsNull() bool {
	return o.state == ValueStateNull
}

func (o enum) IsUnknown() bool {
	return o.state == ValueStateUnknown
}

type FeatureSwitch enum

var (
	featureSwitchValues = ole.New(
		ole.Member[string]{Value: "disabled"},
		ole.Member[string]{Value: "enabled"},
	)
)

func NewFeatureSwitchFromString(s string) (FeatureSwitch, error) {
	e := featureSwitchValues.Parse(s)
	if e == nil {
		return
	}
}
