// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"log"
	"testing"
)

func TestListSystems(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}
	for clientName, client := range clients {
		log.Printf("testing listSystems() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systems, err := client.client.listSystems(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		for _, s := range systems {
			log.Println(s)
		}
	}
}

func TestGetAllSystems(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listSystems() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systemIds, err := client.client.ListSystems(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing getAllSystemsInfo() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systems, err := client.client.GetAllSystemsInfo(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(systemIds) != len(systems) {
			t.Fatalf("system count discrepancy: %d vs. %d", len(systemIds), len(systems))
		}
	}
}

func TestGetSystems(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		log.Printf("testing listSystems() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		systems, err := client.client.listSystems(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		for _, s := range systems {
			log.Printf("testing getSystemInfo() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			system, err := client.client.GetSystemInfo(context.TODO(), s)
			if err != nil {
				t.Fatal(err)
			}
			log.Println(system.Facts.HwModel)
		}
	}
}

func TestSystemsStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() (int, error)
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "", intType: SystemAdminStateNone, stringType: systemAdminStateNone},
		{stringVal: "normal", intType: SystemAdminStateNormal, stringType: systemAdminStateNormal},
		{stringVal: "decomm", intType: SystemAdminStateDecomm, stringType: systemAdminStateDecomm},
		{stringVal: "maint", intType: SystemAdminStateMaint, stringType: systemAdminStateMaint},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp, err := td.stringType.parse()
		if err != nil {
			t.Fatalf("index %d error: %q", i, err)
		}
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != sp ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}
