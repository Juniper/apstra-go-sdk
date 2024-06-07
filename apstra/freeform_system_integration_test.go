package apstra

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestGSa(t *testing.T) {
	var x FreeformSystem
	x.Id = "foo"
	var devprofileid DeviceProfile
	devprofileid.Id = "bUHYZeqRQXafDmuZeaw"
	x.Data = &FreeformSystemData{
		Type:            SystemTypeInternal,
		Label:           "test_generic_system",
		Hostname:        "systemFoo",
		DeviceProfileId: devprofileid,
	}
	rawjson, err := json.Marshal(&x)
	require.NoError(t, err)
	log.Println(string(rawjson))
}
func TestGSb(t *testing.T) {
	var y FreeformSystem
	rawjson := []byte(`{"id":"foo","label":"test_generic_system","hostname":"systemFoo","device_profile":"bUHYZeqRQXafDmuZeaw"}`)
	err := json.Unmarshal(rawjson, &y)
	require.NoError(t, err)
	require.Equal(t, ObjectId("foo"), y.Id, "id mismatch")
}
