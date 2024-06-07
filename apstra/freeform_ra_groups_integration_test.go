package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestRaGroupA(t *testing.T) {
	var x FreeformRaGroup
	x.Id = "foo"
	x.Label = "RaGroupTest"
	x.Data.Key = "key"
	x.Data.Value = "value"
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestRaGroupB(t *testing.T) {
	var y FreeformRaGroup
	rawjson := []byte(`{"id":"foo","label":"RaGroupTest","Data":"key,value"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
