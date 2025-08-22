// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package testclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
)

const EnvApstraExperimental = "APSTRA_EXPERIMENTAL"

type testClientConfig interface {
	clientConfig() apstra.ClientCfg
	clientType() ClientType
	id() string
}

// getTestClientCfgs returns []testClientConfig
func getTestClientCfgs(t testing.TB, ctx context.Context, testConfig TestConfig) []testClientConfig {
	t.Helper()

	var result []testClientConfig

	result = append(result, getAPIOpsClientCfgs(t, ctx, testConfig)...)
	result = append(result, getAWSClientCfgs(t, ctx, testConfig)...)
	result = append(result, getCloudLabsClientCfgs(t, ctx, testConfig)...)
	result = append(result, getSlicerClientCfgs(t, ctx, testConfig)...)

	return result
}

type TestClient struct {
	Client *apstra.Client
	config testClientConfig
}

func (t TestClient) Name() string {
	return fmt.Sprintf("%s/%s/%s", t.config.clientType(), t.config.id(), t.Client.ApiVersion())
}

var (
	testClients      []TestClient
	testClientsMutex sync.Mutex
)

func GetTestClients(t testing.TB, ctx context.Context) []TestClient {
	t.Helper()

	testClientsMutex.Lock()
	defer testClientsMutex.Unlock()

	if testClients != nil {
		return testClients
	}

	experimental, err := strconv.ParseBool(os.Getenv(EnvApstraExperimental))
	require.NoError(t, err)

	// create logfile
	fileName := fmt.Sprintf("test_%s.log", time.Now().Format("20060102-15:04:05"))
	fileFlag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	logFile, err := os.OpenFile(fileName, fileFlag, 0o644)
	require.NoError(t, err)

	testConfig := getTestConfig(t)
	clientCfgs := getTestClientCfgs(t, ctx, testConfig)
	require.NotZerof(t, len(clientCfgs), "There seem to be no clients. Check the environment variables and/or config file: %q", testConfig.path)

	testClients = make([]TestClient, len(clientCfgs))
	for i, testClientCfg := range clientCfgs {
		clientCfg := testClientCfg.clientConfig()
		clientCfg.Experimental = experimental
		clientCfg.Logger = log.New(logFile, "", log.LstdFlags)

		client, err := testClientCfg.clientConfig().NewClient(ctx)
		require.NoError(t, err)

		testClients[i] = TestClient{
			Client: client,
			config: testClientCfg,
		}
	}

	return testClients
}
