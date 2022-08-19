package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestUserLogin(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing Login() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing Logout() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.Logout(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing redundant Logout() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		err = client.client.Logout(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
	}
}
