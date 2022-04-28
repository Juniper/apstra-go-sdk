package apstraTelemetry

import (
	"errors"
	"testing"
)

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

func TestAosLogout(t *testing.T) {
	c, err := NewAosClient(testCfg1)
	if err != nil {
		t.Fatal(err)
	}

	err = c.UserLogin()
	if err != nil {
		t.Fatal(err)
	}

	err = c.UserLogout()
	if err != nil {
		t.Fatal(err)
	}
}
