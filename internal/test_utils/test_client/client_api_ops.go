// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package testclient

import (
	"context"
	"crypto/tls"
	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
)

const (
	envAPIopsTopologyURLList = "API_OPS_URLS"
	envAPIopsURLSep          = ";"
)

var _ testClientConfig = (*APIopsConfig)(nil)

type APIopsConfig struct {
	dcID   string
	config apstra.ClientCfg
}

func (a APIopsConfig) clientConfig() apstra.ClientCfg {
	return a.config
}

func (a APIopsConfig) clientType() ClientType {
	return ClientTypeAPIops
}

func (a APIopsConfig) id() string {
	return a.dcID
}

func getAPIopsClientCfg(t *testing.T, ctx context.Context, APIopsURL string) testClientConfig {
	t.Helper()

	u, err := url.Parse(APIopsURL)
	require.NoErrorf(t, err, "parsing APIopsURL: %q", APIopsURL)

	id := path.Base(u.Path)

	return APIopsConfig{
		dcID: id,
		config: apstra.ClientCfg{
			APIopsDCID: &id,
			Url:        strings.TrimSuffix(u.String(), u.Path),
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

func getAPIopsClientCfgs(t *testing.T, ctx context.Context, testConfig TestConfig) []testClientConfig {
	topologyIDs := testConfig.ApiOpsProxyUrls

	if len(topologyIDs) == 0 {
		list := os.Getenv(envAPIopsTopologyURLList)
		if len(list) != 0 {
			topologyIDs = strings.Split(list, envAPIopsTopologyURLList)
		}
	}

	result := make([]testClientConfig, len(topologyIDs))
	for i, id := range topologyIDs {
		result[i] = getAWSClientCfg(t, ctx, id)
	}

	return result
}
