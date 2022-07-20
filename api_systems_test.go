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

func TestGetAllSystems(t *testing.T) {
	client, err := systemsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	systemIds, err := client.listSystems(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	systems, err := client.getAllSystemsInfo(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(systemIds) != len(systems) {
		t.Fatalf("system count discrepancy: %d vs. %d", len(systemIds), len(systems))
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

func TestSystemsStrings(t *testing.T) {
	type apiStringIota interface {
		String() string
		Int() int
	}

	type apiIotaString interface {
		parse() int
		string() string
	}

	type stringTestData struct {
		stringVal  string
		intType    apiStringIota
		stringType apiIotaString
	}
	testData := []stringTestData{
		{stringVal: "maint", intType: SystemAdminStateMaint, stringType: systemAdminStateMaint},
		{stringVal: "normal", intType: SystemAdminStateNormal, stringType: systemAdminStateNormal},
		{stringVal: "decomm", intType: SystemAdminStateDecomm, stringType: systemAdminStateDecomm},
	}

	for i, td := range testData {
		ii := td.intType.Int()
		is := td.intType.String()
		sp := td.stringType.parse()
		ss := td.stringType.string()
		if td.intType.String() != td.stringType.string() ||
			td.intType.Int() != td.stringType.parse() ||
			td.stringType.string() != td.stringVal {
			t.Fatalf("test index %d mismatch: %d %d '%s' '%s' '%s'",
				i, ii, sp, is, ss, td.stringVal)
		}
	}
}
