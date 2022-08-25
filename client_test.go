package goapstra

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestLoginEmptyPassword(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing empty password Login() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		client.client.cfg.Pass = ""
		err := client.client.Login(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging in with empty password, did not get errror"))
		}
	}
}

func TestLoginBadPassword(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing bad password Login() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		client.client.cfg.Pass = randString(10, "hex")
		err = client.client.Login(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging in with bad password, did not get errror"))
		}
	}
}

func TestLogoutAuthFail(t *testing.T) {
	clientCfgs, err := getTestClientCfgs()
	if err != nil {
		t.Fatal(err)
	}

	for _, cfg := range clientCfgs {
		client, err := cfg.cfg.NewClient()
		log.Printf("testing Login() against %s %s (%s)", cfg.cfgType, cfg.cfgName, client.ApiVersion())
		err = client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("client has this authtoken: '%s'", client.httpHeaders[apstraAuthHeader])
		client.httpHeaders[apstraAuthHeader] = randJwt()
		log.Printf("client authtoken changed to: '%s'", client.httpHeaders[apstraAuthHeader])
		log.Printf("testing failed Logout() against %s %s (%s)", cfg.cfgType, cfg.cfgName, client.ApiVersion())
		err = client.Logout(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging out with bad token, did not get errror"))
		}
	}
}
