//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestClientLog(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		client.client.Logf(1, "log test - client '%s'", clientName)

	}
}

func TestLoginEmptyPassword(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing empty password Login() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		client.client.cfg.Pass = ""
		err := client.client.Login(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging in with empty password, did not get errror"))
		}
	}
}

func TestLoginBadPassword(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing bad password Login() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		client.client.cfg.Pass = randString(10, "hex")
		err = client.client.Login(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging in with bad password, did not get errror"))
		}
	}
}

func TestLogoutAuthFail(t *testing.T) {
	ctx := context.Background()
	clientCfgs, err := getTestClientCfgs(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for name, cfg := range clientCfgs {
		client, err := cfg.cfg.NewClient(ctx)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("testing Login() against %s %s (%s)", cfg.cfgType, name, client.ApiVersion())
		err = client.Login(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("client has this authtoken: '%s'", client.httpHeaders[apstraAuthHeader])
		client.httpHeaders[apstraAuthHeader] = randJwt()
		log.Printf("client authtoken changed to: '%s'", client.httpHeaders[apstraAuthHeader])
		log.Printf("testing failed Logout() against %s %s (%s)", cfg.cfgType, name, client.ApiVersion())
		err = client.Logout(context.TODO())
		if err == nil {
			t.Fatal(fmt.Errorf("tried logging out with bad token, did not get errror"))
		}
	}
}

func TestGetBlueprintOverlayControlProtocol(t *testing.T) {
	ctx := context.Background()

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		bpFunc      func(context.Context, *testing.T, *Client) (*TwoStageL3ClosClient, func(context.Context) error)
		expectedOcp OverlayControlProtocol
	}

	testCases := []testCase{
		{bpFunc: testBlueprintA, expectedOcp: OverlayControlProtocolEvpn},
		{bpFunc: testBlueprintB, expectedOcp: OverlayControlProtocolNone},
	}

	for clientName, client := range clients {
		for i := range testCases {
			bpClient, bpDel := testCases[i].bpFunc(ctx, t, client.client)
			defer func() {
				err := bpDel(ctx)
				if err != nil {
					t.Fatal(err)
				}
			}()

			log.Printf("testing BlueprintOverlayControlProtocol() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			ocp, err := bpClient.client.BlueprintOverlayControlProtocol(ctx, bpClient.blueprintId)
			if err != nil {
				t.Fatal(err)
			}

			if ocp != testCases[i].expectedOcp {
				t.Fatalf("expected overlay control protocol %q, got %q", testCases[i].expectedOcp.String(), ocp.String())
			}
			log.Printf("blueprint %q has overlay control protocol %q", bpClient.blueprintId, ocp.String())
		}
	}
}
