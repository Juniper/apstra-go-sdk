package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestFFLinkA(t *testing.T) {
	var x FreeformLink
	x.Id = "foo"
	x.Data = &FreeformLinkData{
		Type:  LinkTypeEthernet,
		Label: "test_ff_link",
		Speed: "10G",
	}
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestFFLinkB(t *testing.T) {
	var y ConfigTemplate
	rawjson := []byte(`{"id":"foo","link_type":"1","label":"test_ff_link","speed","10G"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
