// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package testclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

const (
	envCloudlabsTopologyIDList = "CLOUDLABS_TOPOLOGIES"
	envCloudlabsTopologyIDSep  = ":"

	cloudlabsTopologyUrlByID = "https://cloudlabs.apstra.com/api/v1.0/topologies/%s"
)

var _ Config = (*CloudLabsConfig)(nil)

type CloudLabsConfig struct {
	topologyID string
	config     apstra.ClientCfg
	switches   []SwitchInfo
}

func (c CloudLabsConfig) Switches() []SwitchInfo {
	return c.switches
}

func (c CloudLabsConfig) clientConfig() apstra.ClientCfg {
	return c.config
}

func (c CloudLabsConfig) clientType() ClientType {
	return ClientTypeCloudLabs
}

func (c CloudLabsConfig) id() string {
	return c.topologyID
}

func getCloudLabsClientCfg(t testing.TB, ctx context.Context, id string) Config {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(cloudlabsTopologyUrlByID, id), nil)
	require.NoError(t, err, "preparing http request for cloudlabs topology")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "requesting cloudlabs topology info")

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("http %d (%s): '%s'", resp.StatusCode, resp.Status, string(body))
	}

	var topology struct {
		Vms []struct {
			Access []struct {
				Username  string `json:"username"`
				Password  string `json:"password"`
				Protocol  string `json:"protocol"`
				Host      string `json:"host"`
				PrivateIp string `json:"privateIp"`
				Port      int    `json:"port"`
			} `json:"access"`
			DeviceType string `json:"deviceType"`
			Role       string `json:"role"`
			Name       string `json:"name"`
		} `json:"vms"`
	}
	err = json.NewDecoder(resp.Body).Decode(&topology)
	require.NoError(t, err, "decoding topology info")

	var user, pass, url string
	var switches []SwitchInfo
	for _, vm := range topology.Vms {
		switch vm.Role {
		case "aos":
			for _, a := range vm.Access {
				if a.Protocol == "https" {
					user = a.Username
					pass = a.Password
					url = fmt.Sprintf("%s://%s:%d", a.Protocol, a.Host, a.Port)
					break
				}
			}
		case "spine", "leaf":
			var platform apstra.AgentPlatform
			switch vm.DeviceType {
			case "veos":
				platform = apstra.AgentPlatformEOS
			case "nxosv":
				platform = apstra.AgentPlatformNXOS
			case "vqfx", "vmx":
				platform = apstra.AgentPlatformJunos
			case "sonic-vs":
				platform = apstra.AgentPlatformNull
			default:
				t.Fatalf("unhandled platform %q for %q device", vm.DeviceType, vm.Role)
			}
			for _, a := range vm.Access {
				if a.Protocol == "ssh" {
					switches = append(switches, SwitchInfo{
						ManagementIP: netip.MustParseAddr(a.PrivateIp),
						Username:     a.Username,
						Password:     a.Password,
						Platform:     platform,
					})
					break
				}
			}
		}
	}

	if user == "" {
		t.Fatalf("cloudlabs topology %q user is empty", id)
	}
	if pass == "" {
		t.Fatalf("cloudlabs topology %q password is empty", id)
	}
	if url == "" {
		t.Fatalf("cloudlabs topology %q does not have an access url", id)
	}

	return CloudLabsConfig{
		topologyID: id,
		switches:   switches,
		config: apstra.ClientCfg{
			Url:  url,
			User: user,
			Pass: pass,
			HttpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						KeyLogWriter:       testutils.KeyLogWriterFromEnv(t),
					},
				},
			},
		},
	}
}

func getCloudLabsClientCfgs(t testing.TB, ctx context.Context, testConfig TestConfig) []Config {
	t.Helper()

	topologyIDs := testConfig.CloudlabsTopologyIds

	if len(topologyIDs) == 0 {
		list := os.Getenv(envCloudlabsTopologyIDList)
		if len(list) != 0 {
			topologyIDs = strings.Split(list, envCloudlabsTopologyIDSep)
		}
	}

	result := make([]Config, len(topologyIDs))
	for i, id := range topologyIDs {
		result[i] = getCloudLabsClientCfg(t, ctx, id)
	}

	return result
}
