package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestPropSetsX(t *testing.T) {
	var x FreeformPropertySet
	x.Id = "foo"
	x.Data = &FFPropertySetData{
		Label:  "test_prop_set",
		Values: "{stuff goes here}",
	}
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestPropSetsY(t *testing.T) {
	var y FreeformPropertySet
	rawjson := []byte(`{"property_set_id":"foo","label":"test_prop_set","Values":"{stuff goes here}"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
