//go:build integration
// +build integration

package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

const envSkipSwitchAgentTest = "GOAPSTRA_SPEEDY_SKIP_SWITCH_AGENT_TEST"
const fileSkipSwitchAgentTest = "/tmp/GOAPSTRA_SPEEDY_SKIP_SWITCH_AGENT_TEST"

func TestGetSystemAgent(t *testing.T) {
	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	skipMsg := make(map[string]string)
	for clientName, client := range clients {
		log.Printf("testing listSystemAgents() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		list, err := client.client.listSystemAgents(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(list) <= 0 {
			skipMsg[clientName] = fmt.Sprintf("cannot get system agents because none exist on '%s'", clientName)
			continue
		}

		for i := 0; i < len(list); i++ {
			log.Printf("testing getSystemAgent(%s) against %s %s (%s)", list[i], client.clientType, clientName, client.client.ApiVersion())
			info, err := client.client.getSystemAgent(context.TODO(), list[i])
			if err != nil {
				t.Fatal(err)
			}
			log.Println(info.Id, info.DeviceFacts.DeviceOsFamily, info.Config.ManagementIp, info.Config.AgentTypeOffBox)
		}
	}
	if len(skipMsg) > 0 {
		sb := strings.Builder{}
		for _, m := range skipMsg {
			sb.WriteString(m + ";")
		}
		t.Skip(sb.String())
	}
}

type testSwitchInfo struct {
	ip       string
	user     string
	pass     string
	platform AgentPlatform
	offbox   AgentTypeOffbox
}

func TestCreateDeleteSwitchAgent(t *testing.T) {
	if os.Getenv(envSkipSwitchAgentTest) == "true" {
		t.Skipf("skipping switch agent tests because '%s' == '%s'", envSkipSwitchAgentTest, "true")
	}
	if _, err := os.Stat(fileSkipSwitchAgentTest); !errors.Is(err, os.ErrNotExist) {
		t.Skipf("skipping switch agent tests because '%s' exists", fileSkipSwitchAgentTest)
	}

	clients, err := getTestClients(context.Background(), t)
	if err != nil {
		t.Fatal(err)
	}

	skipMsg := make(map[string]string)
	for clientName, client := range clients {
		var switchInfo []testSwitchInfo // collect topology-specific switch info here

		switch client.clientType {
		case clientTypeCloudlabs:
			// get the topology by name
			clTopology, err := getCloudlabsTopology(context.Background(), clientName)
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

		var agentIds []ObjectId // do not make fixed length -- some switches might not be available
		var labels []string     // do not make fixed length -- some switches might not be available
		for _, testSwitch := range switchInfo {
			log.Printf("testing CreateSystemAgent() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			label := randString(5, "hex")
			id, err := client.client.CreateSystemAgent(context.TODO(), &SystemAgentRequest{
				ManagementIp:    testSwitch.ip,
				Username:        testSwitch.user,
				Password:        testSwitch.pass,
				Platform:        testSwitch.platform,
				Label:           label,
				AgentTypeOffbox: testSwitch.platform.offbox(),
			})
			if err != nil {
				ace := &ApstraClientErr{}
				if errors.As(err, ace) && ace.Type() == ErrConflict {
					log.Printf("skipping switch '%s' because: '%s'", testSwitch.ip, err.Error())
					continue
				} else {
					t.Fatal(err)
				}
			}
			agentIds = append(agentIds, id)
			labels = append(labels, label)
		}
		if len(agentIds) == 0 {
			skipMsg[clientName] = fmt.Sprintf("no switches available in '%s' for testing", clientName)
			continue
		}

		// run these jobs in parallel
		type runJobResult struct {
			jobStatus *AgentJobStatus
			err       error
		}
		installRunJobResultChan := make(chan runJobResult, len(switchInfo))
		for _, id := range agentIds {
			go func(agentId ObjectId, resultChan chan<- runJobResult) {
				log.Printf("testing SystemAgentRunJob(install) against agent %s %s %s (%s)", agentId, client.clientType, clientName, client.client.ApiVersion())
				status, err := client.client.SystemAgentRunJob(context.TODO(), agentId, AgentJobTypeInstall)
				resultChan <- runJobResult{jobStatus: status, err: err}
			}(id, installRunJobResultChan)
		}
		installJobStatus := make([]AgentJobStatus, len(agentIds))
		for i := range agentIds {
			log.Printf("waiting for SystemAgentRunJob(install) result %d of %d", i+1, len(agentIds))
			result := <-installRunJobResultChan
			if result.err != nil {
				t.Fatal(result.err)
			}
			installJobStatus[i] = *result.jobStatus
		}

		agents := make([]SystemAgent, len(agentIds))
		for i, agentId := range agentIds {
			log.Printf("testing GetSystemAgent() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			agent, err := client.client.GetSystemAgent(context.TODO(), agentId)
			if err != nil {
				t.Fatal(err)
			}
			agents[i] = *agent

			if agent.Config.Label != labels[i] {
				t.Fatalf("label mismatch: expected '%s', got '%s'", labels[i], agent.Config.Label)
			}

			jsonAgent, err := json.Marshal(agent)
			if err != nil {
				t.Fatal(err)
			}
			log.Println(string(jsonAgent))

			log.Printf("testing GetSystemInfo() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			systemInfo, err := client.client.GetSystemInfo(context.TODO(), agent.Status.SystemId)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing updateSystemByAgentId() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.updateSystemByAgentId(context.TODO(), agentId, &SystemUserConfig{
				AdminState:  SystemAdminStateNormal,
				AosHclModel: systemInfo.Facts.AosHclModel,
				Location:    randString(10, "hex"),
			})
			if err != nil {
				t.Fatal(err)
			}
			log.Println("acknowledged!")
		}
		if len(skipMsg) > 0 {
			sb := strings.Builder{}
			for _, msg := range skipMsg {
				sb.WriteString(msg + ";")
			}
			t.Skip(sb.String())
		}

		log.Println("uninstalling agents...")
		// run these jobs in parallel
		type unInstallAgentResult struct {
			jobStatus *AgentJobStatus
			err       error
		}
		unInstallAgentResultChan := make(chan unInstallAgentResult, len(agentIds))
		for _, agentId := range agentIds {
			go func(id ObjectId, resultChan chan<- unInstallAgentResult) {
				log.Printf("testing SystemAgentRunJob(unInstall) against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
				status, err := client.client.SystemAgentRunJob(context.TODO(), id, AgentJobTypeUninstall)
				resultChan <- unInstallAgentResult{jobStatus: status, err: err}
			}(agentId, unInstallAgentResultChan)
		}
		unInstallJobStatus := make([]AgentJobStatus, len(agentIds))
		for i := range agentIds {
			log.Printf("waiting for SystemAgentRunJob(unInstall) result %d of %d", i+1, len(agentIds))
			result := <-unInstallAgentResultChan
			if result.err != nil {
				t.Fatal(result.err)
			}
			unInstallJobStatus[i] = *result.jobStatus
		}

		for _, agent := range agents {
			log.Printf("testing DeleteSystemAgent() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.DeleteSystemAgent(context.TODO(), agent.Id)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("testing deleteSystem() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			err = client.client.deleteSystem(context.TODO(), agent.Status.SystemId)
			if err != nil {
				log.Println(err)
			}
		}
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
		{stringVal: "failed", intType: AgentJobStateFailed, stringType: agentJobStateFailed},
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
