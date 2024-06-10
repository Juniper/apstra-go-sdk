package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestCtX(t *testing.T) {
	var x ConfigTemplate
	x.Id = "foo"
	x.Data = &ConfigTemplateData{
		Label: "test_config_template",
		Text:  "jinja goes here",
	}
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestCtY(t *testing.T) {
	var y ConfigTemplate
	rawjson := []byte(`{"id":"foo","label":"test_config_template","text":"jinja goes here"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
