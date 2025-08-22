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
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/require"
)

const (
	envAWSInstanceIDList = "APSTRA_AWS_INSTANCE_IDS"
	envAWSTopologyIDSep  = ":"

	awsToplogySecretPrefix = "apstra-info-"
)

var _ testClientConfig = (*awsConfig)(nil)

type awsConfig struct {
	instanceID string
	config     apstra.ClientCfg
}

func (a awsConfig) clientConfig() apstra.ClientCfg {
	return a.config
}

func (a awsConfig) clientType() ClientType {
	return ClientTypeAWS
}

func (a awsConfig) id() string {
	return a.instanceID
}

func getAWSClientCfg(t testing.TB, ctx context.Context, id string) testClientConfig {
	t.Helper()

	awsCfg, err := config.LoadDefaultConfig(ctx)
	require.NoError(t, err, "loading default AWS config")

	sm := secretsmanager.NewFromConfig(awsCfg)
	secretId := fmt.Sprintf("%s%s", awsToplogySecretPrefix, id)
	gsvi := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}
	gsvo, err := sm.GetSecretValue(ctx, gsvi)
	require.NoErrorf(t, err, "getting secret '%s' value", secretId)
	require.NotNil(t, gsvo, "secret '%s' value is nil", secretId)

	var topologyInfo struct {
		Url      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = json.NewDecoder(strings.NewReader(*gsvo.SecretString)).Decode(&topologyInfo)
	require.NoError(t, err, "decoding topology info")

	return awsConfig{
		instanceID: id,
		config: apstra.ClientCfg{
			Url:  topologyInfo.Url,
			User: topologyInfo.Username,
			Pass: topologyInfo.Password,
			HttpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						KeyLogWriter: testutils.KeyLogWriterFromEnv(t),
					},
				},
			},
		},
	}
}

func getAWSClientCfgs(t testing.TB, ctx context.Context, testConfig TestConfig) []testClientConfig {
	t.Helper()

	topologyIDs := testConfig.AwsTopologyIds

	if len(topologyIDs) == 0 {
		list := os.Getenv(envAWSInstanceIDList)
		if len(list) != 0 {
			topologyIDs = strings.Split(list, envAWSTopologyIDSep)
		}
	}

	result := make([]testClientConfig, len(topologyIDs))
	for i, id := range topologyIDs {
		result[i] = getAWSClientCfg(t, ctx, id)
	}

	return result
}
