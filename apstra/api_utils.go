package apstra

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func jsonEqual(m1, m2 json.RawMessage) (bool, error) {
	var map1 interface{}
	var map2 interface{}

	var err error
	err = json.Unmarshal(m1, &map1)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling string 1 : %s", err.Error())
	}
	err = json.Unmarshal(m2, &map2)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling string 2 : %s", err.Error())
	}
	return reflect.DeepEqual(map1, map2), nil
}
