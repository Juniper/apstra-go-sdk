// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"bytes"
	"context"
	"log"
	"testing"
)

func TestGetVirtualInfraMgrs(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing getVirtualInfraMgrs() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
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
