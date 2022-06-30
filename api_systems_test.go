package goapstra

import (
	"context"
	"crypto/tls"
	"log"
	"testing"
)

func systemsTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestListSystems(t *testing.T) {
	client, err := systemsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	systems, err := client.listSystems(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range systems {
		log.Println(s)
	}
}

func TestGetSystems(t *testing.T) {
	client, err := systemsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	systems, err := client.listSystems(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range systems {
		system, err := client.getSystemInfo(context.TODO(), s)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(system.Facts.HwModel)
	}
}
