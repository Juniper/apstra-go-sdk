package apstraTelemetry

import (
	"encoding/json"
	"log"
	"testing"
)

func TestGetVersion(t *testing.T) {
	c, err := NewAosClient(testCfg1)
	if err != nil {
		t.Fatal(err)
	}

	err = c.UserLogin()
	if err != nil {
		t.Fatal(err)
	}

	ver, err := c.GetVersion()
	if err != nil {
		t.Fatal(err)
	}

	result, err := json.Marshal(ver)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(string(result))

	err = c.UserLogout()
	if err != nil {
		t.Fatal(err)
	}
}
