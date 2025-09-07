// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package testclient

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Juniper/apstra-go-sdk/apstra"
	testutils "github.com/Juniper/apstra-go-sdk/internal/test_utils"
	"github.com/stretchr/testify/require"
)

const (
	envAPIOpsTopologyURLList = "API_OPS_URLS"
	envAPIOpsURLSep          = ";"
)

var _ Config = (*APIOpsConfig)(nil)

type APIOpsConfig struct {
	dcID   string
	config apstra.ClientCfg
}

func (a APIOpsConfig) clientConfig() apstra.ClientCfg {
	return a.config
}

func (a APIOpsConfig) clientType() ClientType {
	return ClientTypeAPIOps
}

func (a APIOpsConfig) id() string {
	return a.dcID
}

func getAPIOpsClientCfg(t testing.TB, _ context.Context, APIOpsURL string) Config {
	t.Helper()

	u, err := url.Parse(APIOpsURL)
	require.NoErrorf(t, err, "parsing APIOpsURL: %q", APIOpsURL)

	id := path.Base(u.Path)

	return APIOpsConfig{
		dcID: id,
		config: apstra.ClientCfg{
			APIOpsDCID: &id,
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

func getAPIOpsClientCfgs(t testing.TB, ctx context.Context, testConfig TestConfig) []Config {
	t.Helper()

	topologyIDs := testConfig.ApiOpsProxyUrls

	if len(topologyIDs) == 0 {
		list := os.Getenv(envAPIOpsTopologyURLList)
		if len(list) != 0 {
			topologyIDs = strings.Split(list, envAPIOpsTopologyURLList)
		}
	}

	result := make([]Config, len(topologyIDs))
	for i, id := range topologyIDs {
		result[i] = getAPIOpsClientCfg(t, ctx, id)
	}

	return result
}
