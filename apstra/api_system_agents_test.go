// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import "testing"

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
