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

var _ testClientConfig = (*cloudLabsConfig)(nil)

type cloudLabsConfig struct {
	topologyID string
	config     apstra.ClientCfg
}

func (c cloudLabsConfig) clientConfig() apstra.ClientCfg {
	return c.config
}

func (c cloudLabsConfig) clientType() ClientType {
	return ClientTypeCloudLabs
}

func (c cloudLabsConfig) id() string {
	return c.topologyID
}

func getCloudLabsClientCfg(t *testing.T, ctx context.Context, id string) testClientConfig {
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
				Username string `json:"username"`
				Password string `json:"password"`
				Protocol string `json:"protocol"`
				Host     string `json:"host"`
				Port     int    `json:"port"`
			} `json:"access"`
			Name string `json:"name"`
		} `json:"vms"`
	}
	err = json.NewDecoder(resp.Body).Decode(&topology)
	require.NoError(t, err, "decoding topology info")

	var user, pass, url string
VMLOOP:
	for _, vm := range topology.Vms {
		if vm.Name != "aos-vm1" {
			continue
		}
		for _, a := range vm.Access {
			if a.Protocol == "https" {
				user = a.Username
				pass = a.Password
				url = fmt.Sprintf("%s://%s:%d", a.Protocol, a.Host, a.Port)
				break VMLOOP
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

	return cloudLabsConfig{
		topologyID: id,
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

func getCloudLabsClientCfgs(t *testing.T, ctx context.Context, testConfig TestConfig) []testClientConfig {
	topologyIDs := testConfig.CloudlabsTopologyIds

	if len(topologyIDs) == 0 {
		list := os.Getenv(envCloudlabsTopologyIDList)
		if len(list) != 0 {
			topologyIDs = strings.Split(list, envCloudlabsTopologyIDSep)
		}
	}

	result := make([]testClientConfig, len(topologyIDs))
	for i, id := range topologyIDs {
		result[i] = getCloudLabsClientCfg(t, ctx, id)
	}

	return result
}
