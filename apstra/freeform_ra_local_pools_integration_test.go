package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestRaLpA(t *testing.T) {
	var x FreeformRaLocalPools
	x.Id = "foo"
	x.PoolType = "integer"
	x.Label = "foo"
	x.ResourceType = "integer"
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestRaLpB(t *testing.T) {
	var y ConfigTemplate
	rawjson := []byte(`{"id":"foo","pool_type":"integer","resource_type":"integer","label":"foo"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
