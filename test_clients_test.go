//go:build integration
// +build integration

package goapstra

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	clientTypeCloudlabs = "cloudlabs"
	clientTypeAws       = "aws"

	envCloudlabsTopologyIdSep = ":"
	envApstraApiKeyLogFile    = "APSTRA_API_TLS_LOGFILE"
)

type testClientCfg struct {
	cfgType string
	cfg     *ClientCfg
}

type testClient struct {
	clientType string
	client     *Client
}

var testClients map[string]testClient

func getTestClients(ctx context.Context) (map[string]testClient, error) {
	if testClients != nil {
		return testClients, nil
	}

	clientCfgs, err := getTestClientCfgs(ctx)
	if err != nil {
		return nil, err
	}

	testClients = make(map[string]testClient, len(clientCfgs))
	for k, cfg := range clientCfgs {
		client, err := cfg.cfg.NewClient()
		if err != nil {
			return nil, err
		}
		testClients[k] = testClient{
			clientType: cfg.cfgType,
			client:     client,
		}
	}

	//set logfile
	fileName := fmt.Sprintf("test_%s.log", time.Now().Format("20060102-15:04:05"))
	fileFlag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(fileName, fileFlag, 0644)
	if err != nil {
		return nil, err
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

	// add future type clients (slicer?) to testClients slice here

	return testClientCfgs, nil
}
