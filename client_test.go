package goapstra

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestLoginLogoutAuthFail(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing empty password Login() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		client.client.cfg.Pass = ""
		err := client.client.Login(context.TODO())
		if err == nil {
			log.Fatal(fmt.Errorf("tried logging in with empty password, did not get errror"))
		}

		log.Printf("testing bad password Login() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		client.client.cfg.Pass = randString(10, "hex")
		err = client.client.Login(context.TODO())
		if err == nil {
			log.Fatal(fmt.Errorf("tried logging in with bad password, did not get errror"))
		}

		log.Printf("testing failed Logout() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		client.client.httpHeaders[apstraAuthHeader] = randJwt()
		err = client.client.Logout(context.TODO())
		if err == nil {
			log.Fatal(fmt.Errorf("tried logging out with bad token, did not get errror"))
		}
	}
}
