package goapstra

import (
	"bytes"
	"context"
	"log"
	"testing"
)

func TestGetVirtualInfraMgrs(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing getVirtualInfraMgrs() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		vim, err := client.client.getVirtualInfraMgrs(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		buf := bytes.NewBuffer([]byte{})
		err = pp(vim, buf)
		if err != nil {
			t.Fatal(err)
		}
	}
}
