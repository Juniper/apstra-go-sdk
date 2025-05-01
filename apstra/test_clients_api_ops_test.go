// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	envApiOpsTopologyUrlList = "API_OPS_URLS"
	envApiOpsUrlSep          = ";"
)

func apiOpsTopologyIdsFromEnv() []string {
	ids := os.Getenv(envApiOpsTopologyUrlList)
	if ids == "" {
		return nil
	}
	return strings.Split(ids, envApiOpsUrlSep)
}

func getApiOpsTestClientCfgs(_ context.Context) (map[string]testClientCfg, error) {
	cfg, err := GetTestConfig()
	if err != nil {
		return nil, err
	}

	var proxyUrls []string
	if len(cfg.ApiOpsProxyUrls) > 0 {
		proxyUrls = make([]string, len(cfg.ApiOpsProxyUrls))
		for i, proxyUrl := range cfg.ApiOpsProxyUrls {
			proxyUrls[i] = proxyUrl
		}
	}

	if len(proxyUrls) == 0 {
		proxyUrls = apiOpsTopologyIdsFromEnv()
	}

	result := make(map[string]testClientCfg, len(proxyUrls))
	for _, id := range proxyUrls {
		u, err := url.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("api ops proxy url parse error: %w", err)
		}

		klw, err := keyLogWriterFromEnv(envApstraApiKeyLogFile)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			KeyLogWriter:       klw,
		}

		result[id] = testClientCfg{
			cfgType: "api-ops",
			cfg: &ClientCfg{
				Url:        strings.TrimSuffix(u.String(), u.Path),
				apiOpsDcId: toPtr(path.Base(u.Path)),
				HttpClient: &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}},
			},
		}
	}

	return result, nil
}
