package goapstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"testing"
)

func apstraVersionTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestGetVersion(t *testing.T) {
	client, err := apstraVersionTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	ver, err := client.getVersion(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	result, err := json.Marshal(ver)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(string(result))

	err = client.Logout(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}
