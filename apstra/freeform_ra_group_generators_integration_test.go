package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestRaGroupGenA(t *testing.T) {
	var x FreeformRaGroupGenerator
	x.Data = new(FreeformRaGroupGeneratorData)
	x.Id = "foo"
	x.Data.Scope = "node('link', role='internal', name='target')"
	x.Data.Label = "GroupGenTest"
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestRaGroupGenB(t *testing.T) {
	var y ConfigTemplate
	rawjson := []byte(`{"id":"foo","label":"GroupGenTest","scope":"node('link', role='internal', name='target')"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
