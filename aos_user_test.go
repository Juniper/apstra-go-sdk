package apstraTelemetry

import (
	"errors"
	"testing"
)

var testCfg1 = AosClientCfg{
	Host:   "66.129.234.206",
	Port:   uint16(37000),
	Scheme: "hxxps",
	User:   "admin",
	Pass:   "admin",
}

func TestNewAosClient(t *testing.T) {
	c, err := NewAosClient(testCfg1)
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal(errors.New("NewAosClient returned nil client"))
	}
}

func TestAosLogin(t *testing.T) {
	c, err := NewAosClient(testCfg1)
	if err != nil {
		t.Fatal(err)
	}

	err = c.UserLogin()
	if err != nil {
		t.Fatal(err)
	}
}
