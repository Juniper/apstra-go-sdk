package goapstra

import (
	"bytes"
	"context"
	"crypto/tls"
	"log"
	"testing"
)

func telemetryServicesTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestGetTelemetryServicesDeviceMapping(t *testing.T) {
	client, err := telemetryServicesTestClient1()
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Login(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Logout(context.TODO())

	result, err := client.GetTelemetryServicesDeviceMapping(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = pp(result, buf)
	if err != nil {
		t.Fatal(err)
	}
}
