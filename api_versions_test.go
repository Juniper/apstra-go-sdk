package goapstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"testing"
)

func clientTestVersionsCfg1() (*ClientCfg, error) {
	return &ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	}, nil
}

func TestGetVersionsServer(t *testing.T) {
	cfg, err := clientTestVersionsCfg1()
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Logout(context.TODO())

	response, err := client.getVersionsServer(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(string(body))
}
