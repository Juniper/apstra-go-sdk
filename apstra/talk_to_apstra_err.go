package apstra

import (
	"encoding/json"
	"fmt"
)

type apiErrors struct {
	Errors []string
}

func (o *apiErrors) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Errors interface{} `json:"errors"`
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	switch t := raw.Errors.(type) {
	case string:
		o.Errors = []string{t}
	case []string: // this may not be needed. []string comes through as []interface{}
		o.Errors = t
	case []interface{}:
		o.Errors = make([]string, len(t))
		for i, el := range t {
			s, ok := el.(string)
			if !ok {
				return fmt.Errorf("unexpected error type at index %d: %T", i, el)
			}
			o.Errors[i] = s
		}
	case nil:
		o.Errors = nil
		return nil
	default:
		return fmt.Errorf("unexpected error type: %T", raw.Errors)
	}
	
	return nil
}
