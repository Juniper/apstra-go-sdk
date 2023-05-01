package apstra

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Copied from https://gist.github.com/turtlemonvh/e4f7404e28387fadb8ad275a99596f67
func areEqualJSON(s1, s2 json.RawMessage) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal(s1, &o1)
	if err != nil {
		return false, fmt.Errorf("Error unmarshalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal(s2, &o2)
	if err != nil {
		return false, fmt.Errorf("Error unmarshalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
