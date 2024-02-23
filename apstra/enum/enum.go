package enum

import (
	olEnum "github.com/orsinium-labs/enum"
)

type EnumType int

type Value interface {
	Equal(instance Value) bool
	String() string
	Type() EnumType
	member() olEnum.Member[string]
}

func newInstance(t EnumType, s string) Value {
	return value{
		enumType: t,
		value:    &olEnum.Member[string]{Value: s},
	}
}

var _ Value = new(value)

type value struct {
	enumType EnumType
	value    *olEnum.Member[string]
}

func (o value) Equal(e Value) bool {
	if o.enumType != e.Type() {
		return false
	}

	return o.value.Value == e.String()
}

func (o value) String() string {
	return o.value.Value
}

func (o value) Type() EnumType {
	return o.enumType
}

func (o value) member() olEnum.Member[string] {
	return *o.value
}

// New returns n Value based on t and s, or nil if t, s or the
// t, s combination is invalid.
func New(t EnumType, s string) Value {
	if valueFuncs, ok := enumTypeToFuncs[t]; ok {
		members := make([]olEnum.Member[string], len(valueFuncs))
		for i, valueFunc := range valueFuncs {
			members[i] = valueFunc().(value).member()
		}

		e := olEnum.New(members...).Parse(s)
		if e == nil {
			return nil // s not a valid member of they enum type
		}

		return value{
			enumType: t,
			value:    e,
		}
	}

	return nil // t not a valid enum type
}
