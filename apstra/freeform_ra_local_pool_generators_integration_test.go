package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestRaLocalPoolGenA(t *testing.T) {
	var x FreeformRaLocalPoolGenerator
	x.Id = "foo"
	x.Scope = "node('link', role='internal', name='target')"
	x.PoolType = 1
	x.Label = "RaLocalPoolGenTest"
	x.ResourceType = "integer"
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestRaLocalPoolGenB(t *testing.T) {
	var y FreeformRaLocalPoolGenerator
	rawjson := []byte(`{"id":"foo","label":"RaLocalPoolGenTest","PoolType":"1","ResourceType":"integer","scope":"node('link', role='internal', name='target')"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
