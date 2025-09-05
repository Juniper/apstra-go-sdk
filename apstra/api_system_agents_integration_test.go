// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	testclient "github.com/Juniper/apstra-go-sdk/internal/test_utils/test_client"
	"github.com/stretchr/testify/require"
)

const (
	envSkipSwitchAgentTest  = "SPEEDY_SKIP_SWITCH_AGENT_TEST"
	fileSkipSwitchAgentTest = "/tmp/SPEEDY_SKIP_SWITCH_AGENT_TEST"
)

func TestGetSystemAgent(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, t.Context())
	clients := testclient.GetTestClients(t, ctx)

	skipMsg := make(map[int]string)
	for i, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, t.Context())

			list, err := client.Client.ListSystemAgents(ctx)
			require.NoError(t, err)

			if len(list) <= 0 {
				skipMsg[i] = fmt.Sprintf("cannot get system agents because none exist on '%s'", client.Name())
				return
			}

			for i := range len(list) {
				info, err := client.Client.GetSystemAgent(ctx, list[i])
				require.NoError(t, err)
				t.Log(info.Id, info.DeviceFacts.DeviceOsFamily, info.Config.ManagementIp, info.Config.AgentTypeOffBox)
			}
		})
	}

	if len(skipMsg) > 0 {
		sb := strings.Builder{}
		for _, m := range skipMsg {
			sb.WriteString(m + ";")
		}
		t.Skip(sb.String())
	}
}

func TestCreateDeleteSwitchAgent(t *testing.T) {
	ctx := testutils.WrapCtxWithTestId(t, t.Context())
	clients := testclient.GetTestClients(t, ctx)

	if s, ok := os.LookupEnv(envSkipSwitchAgentTest); ok {
		b, err := strconv.ParseBool(s)
		require.NoError(t, err)
		if b {
			t.Skipf("skipping switch agent tests because %q has value %q", envSkipSwitchAgentTest, s)
		}
	}

	if _, err := os.Stat(fileSkipSwitchAgentTest); !errors.Is(err, os.ErrNotExist) {
		// don't check this error - we only care if it's "file not found"
		t.Skipf("skipping switch agent tests because %q exists", fileSkipSwitchAgentTest)
	}

	skipMsg := make(map[int]string)
	for _, client := range clients {
		t.Run(client.Name(), func(t *testing.T) {
			t.Parallel()
			ctx := testutils.WrapCtxWithTestId(t, t.Context())

			var switches []testclient.SwitchInfo
			switch client.Type() {
			case testclient.ClientTypeCloudLabs:
				// get the switch info
				switches = client.Config().(testclient.CloudLabsConfig).Switches()
			}

			agentIds := make([]apstra.ObjectId, 0, len(switches))
			labels := make([]string, 0, len(switches))
			for _, testSwitch := range switches {
				if strings.HasSuffix(testSwitch.ManagementIP.String(), ".15") {
					continue
				}
				label := testutils.RandString(5, "hex")
				agentId, err := client.Client.CreateSystemAgent(ctx, &apstra.SystemAgentRequest{
					ManagementIp:    testSwitch.ManagementIP.String(),
					Username:        testSwitch.Username,
					Password:        testSwitch.Password,
					Platform:        testSwitch.Platform,
					Label:           label,
					AgentTypeOffbox: testSwitch.Platform == apstra.AgentPlatformJunos,
					OperationMode:   apstra.SystemManagementLevelFullControl,
				})
				if err != nil {
					var ace apstra.ClientErr
					if errors.As(err, &ace) && ace.Type() == apstra.ErrConflict {
						log.Printf("skipping switch at %q because: %q", testSwitch.ManagementIP, err.Error())
						continue
					} else {
						t.Fatal(err)
					}
				}
				agentIds = append(agentIds, agentId)
				labels = append(labels, label)
			}
			if len(agentIds) == 0 {
				t.Skipf("no switches available in '%s' for testing", client.Name())
			}

			// run these jobs in parallel
			type runJobResult struct {
				jobStatus *apstra.AgentJobStatus
				err       error
			}
			installRunJobResultChan := make(chan runJobResult, len(switches))
			for _, id := range agentIds {
				go func(agentId apstra.ObjectId, resultChan chan<- runJobResult) {
					status, err := client.Client.SystemAgentRunJob(ctx, agentId, apstra.AgentJobTypeInstall)
					resultChan <- runJobResult{jobStatus: status, err: err}
				}(id, installRunJobResultChan)
			}

			installJobStatus := make([]apstra.AgentJobStatus, len(agentIds))
			for i := range agentIds {
				log.Printf("%s waiting for SystemAgentRunJob(install) result %d of %d", client.Name(), i+1, len(agentIds))
				result := <-installRunJobResultChan
				require.NoError(t, result.err)
				installJobStatus[i] = *result.jobStatus
			}

			agents := make([]apstra.SystemAgent, len(agentIds))
			for i, agentId := range agentIds {
				agent, err := client.Client.GetSystemAgent(ctx, agentId)
				require.NoError(t, err)
				require.NotNil(t, agent)
				require.Equal(t, labels[i], agent.Config.Label)
				agents[i] = *agent

				jsonAgent, err := json.Marshal(agent)
				require.NoError(t, err)
				log.Println(string(jsonAgent))

				systemInfo, err := client.Client.GetSystemInfo(ctx, agent.Status.SystemId)
				require.NoError(t, err)

				err = client.Client.UpdateSystemByAgentId(ctx, agentId, &apstra.SystemUserConfig{
					AdminState:  apstra.SystemAdminStateNormal,
					AosHclModel: systemInfo.Facts.AosHclModel,
					Location:    testutils.RandString(10, "hex"),
				})
				require.NoError(t, err)
				log.Println("acknowledged!")
			}
			if len(skipMsg) > 0 {
				sb := strings.Builder{}
				for _, msg := range skipMsg {
					sb.WriteString(msg + ";")
				}
				t.Skip(sb.String())
			}

			log.Printf("%s uninstalling agents...", client.Name())
			// run these jobs in parallel
			type unInstallAgentResult struct {
				jobStatus *apstra.AgentJobStatus
				err       error
			}
			unInstallAgentResultChan := make(chan unInstallAgentResult, len(agentIds))
			for _, agentId := range agentIds {
				go func(id apstra.ObjectId, resultChan chan<- unInstallAgentResult) {
					status, err := client.Client.SystemAgentRunJob(ctx, id, apstra.AgentJobTypeUninstall)
					resultChan <- unInstallAgentResult{jobStatus: status, err: err}
				}(agentId, unInstallAgentResultChan)
			}
			unInstallJobStatus := make([]apstra.AgentJobStatus, len(agentIds))
			for i := range agentIds {
				log.Printf("%s waiting for SystemAgentRunJob(unInstall) result %d of %d", client.Name(), i+1, len(agentIds))
				result := <-unInstallAgentResultChan
				require.NoError(t, result.err)
				require.NotNil(t, result.jobStatus)
				unInstallJobStatus[i] = *result.jobStatus
			}

			for _, agent := range agents {
				err := client.Client.DeleteSystemAgent(ctx, agent.Id)
				require.NoError(t, err)

				err = client.Client.DeleteSystem(ctx, agent.Status.SystemId)
				require.NoError(t, err)
			}
		})
	}
}
