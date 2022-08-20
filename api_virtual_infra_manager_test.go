package goapstra

import (
	"bytes"
	"context"
	"log"
	"testing"
)

const (
	mockVirtualInfraMgrInfo = `{
      "items": [
        {
          "connection_state": "connected",
          "last_successful_collection_time": "2022-05-17T20:36:00.712641Z",
          "service_enabled": true,
          "management_ip": "100.123.91.106",
          "system_id": "382f335c-13cc-47b1-8591-7f90d92a20bc",
          "agent_id": "6ec9ac7d-bde2-48b6-9d01-71fd815531a6",
          "virtual_infra_type": "vcenter"
        }
      ]
    }`
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
