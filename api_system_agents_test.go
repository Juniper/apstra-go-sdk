package goapstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func systemAgentsTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		Timeout:   -1,
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
}

func TestListSystemAgents(t *testing.T) {
	client, err := systemAgentsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	list, err := client.listAgents(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range list {
		log.Println(a)
	}
}

func TestGetSystemAgent(t *testing.T) {
	client, err := systemAgentsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	list, err := client.listAgents(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(list) <= 0 {
		t.Fatalf("cannot get system agent - %d agents exist on this apstra", len(list))
	}

	for i := 0; i < len(list); i++ {
		info, err := client.getSystemAgent(context.TODO(), list[i])
		if err != nil {
			t.Fatal(err)
		}
		log.Println(info.Id, info.DeviceFacts.DeviceOsFamily, info.Config.ManagementIp, info.Config.AgentTypeOffBox)
	}
}

func TestCreateOffboxAgent(t *testing.T) {
	client, err := systemAgentsTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	qfxIp, found := os.LookupEnv("QFX_IP")
	if !found {
		t.Fatal("env QFX_IP not found, cannot create system agent")
	}
	qfxUser, found := os.LookupEnv("QFX_USER")
	if !found {
		t.Fatal("env QFX_USER not found, cannot create system agent")
	}
	qfxPass, found := os.LookupEnv("QFX_PASS")
	if !found {
		t.Fatal("env QFX_PASS not found, cannot create system agent")
	}

	label := randString(5, "hex")

	agentId, err := client.CreateAgent(context.TODO(), &SystemAgentRequest{
		ManagementIp: qfxIp,
		Username:     qfxUser,
		Password:     qfxPass,
		//Platform:     AgentPlatformJunos,
		Label:           label,
		AgentTypeOffbox: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(agentId)

	jobStatus, err := client.SystemAgentRunJob(context.TODO(), agentId, AgentJobTypeInstall)
	if err != nil {
		t.Fatal(err)
	}
	jsonJobStatus, err := json.Marshal(jobStatus)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("jobstatus: %s", string(jsonJobStatus))

	agent, err := client.GetSystemAgent(context.TODO(), agentId)
	if err != nil {
		t.Fatal(err)
	}

	if agent.Config.Label != label {
		t.Fatalf("label mismatch: expected '%s', got '%s'", label, agent.Config.Label)
	}

	jsonAgent, err := json.Marshal(agent)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(jsonAgent))

	systemInfo, err := client.GetSystemInfo(context.TODO(), agent.Status.SystemId)
	if err != nil {
		t.Fatal(err)
	}

	err = client.updateSystemByAgentId(context.TODO(), agentId, &SystemUserConfig{
		AdminState:  SystemAdminStateNormal,
		AosHclModel: systemInfo.Facts.AosHclModel,
		Location:    randString(10, "hex"),
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Println("acknowledged!")
	log.Println("deleting...")

	jobStatus, err = client.SystemAgentRunJob(context.TODO(), agentId, AgentJobTypeUninstall)
	if err != nil {
		t.Fatal(err)
	}
	jsonJobStatus, err = json.Marshal(jobStatus)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(jsonJobStatus))

	err = client.DeleteSystemAgent(context.TODO(), agentId)
	if err != nil {
		t.Fatal(err)
	}

	err = client.deleteSystem(context.TODO(), agent.Status.SystemId)
	if err != nil {
		log.Println(err)
	}
}

func TestSystemAgentsStrings(t *testing.T) {
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
		{stringVal: "connected", intType: AgentCxnStateConnected, stringType: agentCxnStateConnected},
		{stringVal: "disconnected", intType: AgentCxnStateDisconnected, stringType: agentCxnStateDisconnected},
		{stringVal: "auth_failed", intType: AgentCxnStateAuthFail, stringType: agentCxnStateAuthFail},
		{stringVal: "full_control", intType: AgentModeFull, stringType: agentModeFull},
		{stringVal: "telemetry_only", intType: AgentModeTelemetry, stringType: agentModeTelemetry},

		{stringVal: "", intType: AgentJobTypeNull, stringType: agentJobTypeNull},
		{stringVal: "none", intType: AgentJobTypeNone, stringType: agentJobTypeNone},
		{stringVal: "check", intType: AgentJobTypeCheck, stringType: agentJobTypeCheck},
		{stringVal: "install", intType: AgentJobTypeInstall, stringType: agentJobTypeInstall},
		{stringVal: "revertToPristine", intType: AgentJobTypeRevertToPristine, stringType: agentJobTypeRevertToPristine},
		{stringVal: "upgrade", intType: AgentJobTypeUpgrade, stringType: agentJobTypeUpgrade},
		{stringVal: "uninstall", intType: AgentJobTypeUninstall, stringType: agentJobTypeUninstall},

		{stringVal: "", intType: AgentPlatformNull, stringType: agentPlatformNull},
		{stringVal: "junos", intType: AgentPlatformJunos, stringType: agentPlatformJunos},
		{stringVal: "eos", intType: AgentPlatformEOS, stringType: agentPlatformEOS},
		{stringVal: "nxos", intType: AgentPlatformNXOS, stringType: agentPlatformNXOS},

		{stringVal: "", intType: AgentJobStateNull, stringType: agentJobStateNull},
		{stringVal: "init", intType: AgentJobStateInit, stringType: agentJobStateInit},
		{stringVal: "inprogress", intType: AgentJobStateInProgress, stringType: agentJobStateInProgress},
		{stringVal: "success", intType: AgentJobStateSuccess, stringType: agentJobStateSuccess},
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

func TestAgentTypeOffbox(t *testing.T) {
	t1 := AgentTypeOffbox(true)
	e1 := rawAgentType("offbox")
	r1 := t1.raw()
	if r1 != e1 {
		t.Fatalf("expected '%s', got '%s'", e1, r1)
	}

	t2 := AgentTypeOffbox(false)
	e2 := rawAgentType("onbox")
	r2 := t2.raw()
	if r2 != e2 {
		t.Fatalf("expected '%s', got '%s'", e2, r2)
	}

	t3 := rawAgentType("offbox")
	e3 := AgentTypeOffbox(true)
	p3 := t3.parse()
	if p3 != e3 {
		t.Fatalf("expected '%t', got '%t'", e3, p3)
	}

	t4 := rawAgentType("onbox")
	e4 := AgentTypeOffbox(false)
	p4 := t4.parse()
	if p4 != e4 {
		t.Fatalf("expected '%t', got '%t'", e4, p4)
	}
}
