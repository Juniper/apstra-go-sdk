package goapstra

import (
	"context"
	"crypto/tls"
	"log"
	"testing"
)

func userTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestLogin(t *testing.T) {
	c, err := userTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Login(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogout(t *testing.T) {
	c, err := userTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Login(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	err = c.Logout(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserLogin(t *testing.T) {
	DebugLevel = 2
	clients, _, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		client.cfg.Timeout = -1
		log.Printf("testing Login() and Logout() with %s client", clientName)
		err = client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		err = client.Logout(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		err = client.Logout(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
	}
}
