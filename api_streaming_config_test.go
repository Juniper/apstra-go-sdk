package goapstra

import (
	"context"
	"crypto/tls"
	"testing"
)

func streamingConfigTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestClient_GetAllStreamingConfigs(t *testing.T) {
	client, err := streamingConfigTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Login(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}
