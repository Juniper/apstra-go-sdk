package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetSystemAgent(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		log.Printf("testing listAgents() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		list, err := client.client.listAgents(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(list) <= 0 {
			t.Fatalf("cannot get system agent - %d agents exist on this apstra", len(list))
		}

		for i := 0; i < len(list); i++ {
			log.Printf("testing getSystemAgent() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
			info, err := client.client.getSystemAgent(context.TODO(), list[i])
			if err != nil {
				t.Fatal(err)
			}
			log.Println(info.Id, info.DeviceFacts.DeviceOsFamily, info.Config.ManagementIp, info.Config.AgentTypeOffBox)
		}
	}
}

type testSwitchInfo struct {
	ip       string
	user     string
	pass     string
	platform AgentPlatform
	offbox   AgentTypeOffbox
}

func TestCreateOffboxAgent(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		var switchInfo []testSwitchInfo // collect topology-specific switch info here

		switch client.clientType {
		case clientTypeCloudlabs:
			// get the topology by name
			clTopology, err := getCloudlabsTopology(client.clientName)
			if err != nil {
				t.Fatal(err)
			}

			// get the switch info
			clSwitchInfo, err := clTopology.getSwitchInfo()
			if err != nil {
				t.Fatal(err)
			}

			// save the switch info to the topology-independent slice
			switchInfo = make([]testSwitchInfo, len(clSwitchInfo))
			for i, si := range clSwitchInfo {
				switchInfo[i] = testSwitchInfo{ip: si.sshIp, user: si.sshUser, pass: si.sshPass, platform: si.deviceType.platform(), offbox: si.deviceType.platform().offbox()}
			}
		}

		type createAgentResult struct {
			agentId ObjectId
			label   string
			err     error
		}

		createAgentResultChan := make(chan createAgentResult, len(switchInfo))
		for _, testSwitch := range switchInfo {
			go func(si testSwitchInfo, result chan createAgentResult) {
				log.Printf("testing CreateAgent() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
				label := randString(5, "hex")
				id, err := client.client.CreateAgent(context.TODO(), &SystemAgentRequest{
					ManagementIp:    si.ip,
					Username:        si.user,
					Password:        si.pass,
					Platform:        si.platform,
					Label:           label,
					AgentTypeOffbox: si.platform.offbox(),
				})
				createAgentResultChan <- createAgentResult{agentId: id, label: label, err: err}
			}(testSwitch, createAgentResultChan)
		}

		agentIds := make([]ObjectId, len(switchInfo))
		labels := make([]string, len(switchInfo))
		for i := 0; i < len(switchInfo); i++ {
			result := <-createAgentResultChan
			if result.err != nil {
				t.Fatal(err)
			}
			agentIds[i] = result.agentId
			labels[i] = result.label
		}

		//log.Printf("testing SystemAgentRunJob() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		//jobStatus, err := client.client.SystemAgentRunJob(context.TODO(), agentId, AgentJobTypeInstall)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//jsonJobStatus, err := json.Marshal(jobStatus)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//log.Printf("jobstatus: %s", string(jsonJobStatus))
		//
		//log.Printf("testing GetSystemAgent() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		//agent, err := client.client.GetSystemAgent(context.TODO(), agentId)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//
		//if agent.Config.Label != label {
		//	t.Fatalf("label mismatch: expected '%s', got '%s'", label, agent.Config.Label)
		//}
		//
		//jsonAgent, err := json.Marshal(agent)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//log.Println(string(jsonAgent))
		//
		//log.Printf("testing GetSystemInfo() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		//systemInfo, err := client.client.GetSystemInfo(context.TODO(), agent.Status.SystemId)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//
		//log.Printf("testing updateSystemByAgentId() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		//err = client.client.updateSystemByAgentId(context.TODO(), agentId, &SystemUserConfig{
		//	AdminState:  SystemAdminStateNormal,
		//	AosHclModel: systemInfo.Facts.AosHclModel,
		//	Location:    randString(10, "hex"),
		//})
		//if err != nil {
		//	t.Fatal(err)
		//}
		//
		//log.Println("acknowledged!")
		//log.Println("deleting...")
		//
		//log.Printf("testing SystemAgentRunJob() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		//jobStatus, err = client.client.SystemAgentRunJob(context.TODO(), agentId, AgentJobTypeUninstall)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//jsonJobStatus, err = json.Marshal(jobStatus)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//log.Println(string(jsonJobStatus))
		//
		//log.Printf("testing DeleteSystemAgent() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		//err = client.client.DeleteSystemAgent(context.TODO(), agentId)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//
		//log.Printf("testing deleteSystem() against %s %s (%s)", client.clientType, client.clientName, client.client.ApiVersion())
		//err = client.client.deleteSystem(context.TODO(), agent.Status.SystemId)
		//if err != nil {
		//	log.Println(err)
		//}
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
