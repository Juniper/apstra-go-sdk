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
	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

const (
	envSlicerTopologyIDList = "SLICER_TOPOLOGIES"
	envSlicerTopologyIDSep  = ":"

	slicerTopologyUrlByID = "http://slicer-topology-management-ui.k8s-autobuild.dc1.apstra.com/v1_1/systest/%s"
)

var _ testClientConfig = (*slicerConfig)(nil)

type slicerConfig struct {
	topologyID string
	config     apstra.ClientCfg
}

func (s slicerConfig) clientConfig() apstra.ClientCfg {
	return s.config
}

func (s slicerConfig) clientType() ClientType {
	return ClientTypeSlicer
}

func (s slicerConfig) id() string {
	return s.topologyID
}

func getSlicerClientCfg(t *testing.T, ctx context.Context, id string) testClientConfig {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(slicerTopologyUrlByID, id), nil)
	require.NoError(t, err, "preparing http request for slicer topology")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "requesting slicer topology info")

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("http %d (%s): '%s'", resp.StatusCode, resp.Status, string(body))
	}

	var topology struct {
		DeployStatus string `json:"deploy_status"`
		DeployModel  struct {
			DutmgmtConnectivity map[string]string `json:"dutmgmt_connectivity"`
		} `json:"deploy_model"`
	}
	err = json.NewDecoder(resp.Body).Decode(&topology)
	require.NoError(t, err, "decoding topology info")
	require.Contains(t, topology.DeployModel.DutmgmtConnectivity, "aos-vm1")

	return slicerConfig{
		topologyID: id,
		config: apstra.ClientCfg{
			Url:  fmt.Sprintf("https://%s", topology.DeployModel.DutmgmtConnectivity["aos-vm1"]),
			User: "admin",
			Pass: "admin",
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

func getSlicerClientCfgs(t *testing.T, ctx context.Context, testConfig TestConfig) []testClientConfig {
	topologyIDs := testConfig.SlicerTopologyIds

	if len(topologyIDs) == 0 {
		list := os.Getenv(envSlicerTopologyIDList)
		if len(list) != 0 {
			topologyIDs = strings.Split(list, envSlicerTopologyIDSep)
		}
	}

	result := make([]testClientConfig, len(topologyIDs))
	for i, id := range topologyIDs {
		result[i] = getSlicerClientCfg(t, ctx, id)
	}

	return result
}
