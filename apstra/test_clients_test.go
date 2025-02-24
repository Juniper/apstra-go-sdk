// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration
// +build integration

package apstra

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	clientTimeoutSeconds = 30

	clientTypeCloudlabs = "cloudlabs"
	clientTypeAws       = "aws"

	envCloudlabsTopologyIdSep = ":"
	envApstraApiKeyLogFile    = "SSLKEYLOGFILE"
	envApstraExperimental     = "APSTRA_EXPERIMENTAL"
)

type testClientCfg struct {
	cfgType string
	cfg     *ClientCfg
}

type testClient struct {
	clientType string
	client     *Client
	id         string
}

func (o testClient) name() string {
	return fmt.Sprintf("%s/%s/%s", o.clientType, o.id, o.client.apiVersion)
}

var testClients map[string]testClient

func getTestClients(ctx context.Context, t *testing.T) (map[string]testClient, error) {
	t.Helper()

	if testClients != nil {
		return testClients, nil
	}

	clientCfgs, err := getTestClientCfgs(ctx)
	require.NoError(t, err)

	if _, ok := os.LookupEnv(envApstraExperimental); ok {
		for k := range clientCfgs {
			clientCfgs[k].cfg.Experimental = true
		}
	}

	for k := range clientCfgs {
		clientCfgs[k].cfg.Timeout = clientTimeoutSeconds * time.Second
	}

	testClients = make(map[string]testClient, len(clientCfgs))
	for k, cfg := range clientCfgs {
		client, err := cfg.cfg.NewClient(ctx)
		require.NoError(t, err)

		testClients[k] = testClient{
			clientType: cfg.cfgType,
			client:     client,
			id:         k,
		}
	}

	// set logfile
	fileName := fmt.Sprintf("test_%s.log", time.Now().Format("20060102-15:04:05"))
	fileFlag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(fileName, fileFlag, 0o644)
	require.NoError(t, err)

	// There are no test clients. Might be worth logging
	if len(testClients) == 0 {
		t.Fatal("Error : There seem to be no clients. Check the environment variables and/or config file.")
	}

	for k := range testClients {
		testClients[k].client.logger = log.New(f, "", log.LstdFlags)
		testClients[k].client.cfg.LogLevel = 1
	}

	return testClients, nil
}

// getTestClientCfgs returns map[string]testClientCfg keyed by
// the test environment name (e.g. cloudlabs topology ID).
func getTestClientCfgs(ctx context.Context) (map[string]testClientCfg, error) {
	testClientCfgs := make(map[string]testClientCfg)

	// add cloudlabs clients to testClients slice
	clTestClientCfgs, err := getCloudlabsTestClientCfgs(ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range clTestClientCfgs {
		testClientCfgs[k] = v
	}

	// add aws clients to testClients slice
	awsTestClientCfgs, err := getAwsTestClientCfgs(ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range awsTestClientCfgs {
		testClientCfgs[k] = v
	}

	// add slicer to testClients slice here
	slicerTestClientCfgs, err := getSlicerTestClientCfgs(ctx)
	if err != nil {
		return nil, err
	}

	for k, v := range slicerTestClientCfgs {
		testClientCfgs[k] = v
	}

	return testClientCfgs, nil
}
