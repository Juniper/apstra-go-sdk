// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration

package apstra

import (
	"context"
	"fmt"
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
	return strings.Split(ids, envApiOpsTopologyUrlList)
}

func getApiOpsTestClientCfgs(_ context.Context) (map[string]testClientCfg, error) {
	cfg, err := GetTestConfig()
	if err != nil {
		return nil, err
	}

	var topologyIds []string
	if len(cfg.ApiOpsProxyUrls) > 0 {
		topologyIds = make([]string, len(cfg.ApiOpsProxyUrls))
		for i, s := range cfg.ApiOpsProxyUrls {
			topologyIds[i] = s
		}
	}

	if len(topologyIds) == 0 {
		topologyIds = apiOpsTopologyIdsFromEnv()
	}

	result := make(map[string]testClientCfg, len(topologyIds))
	for _, id := range topologyIds {
		u, err := url.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("api ops proxy url parse error: %w", err)
		}

		p := u.Path
		x := strings.TrimSuffix(u.String(), p)

		//tlsConfig := &tls.Config{InsecureSkipVerify: true}
		//
		//klw, err := keyLogWriterFromEnv(envApstraApiKeyLogFile)
		//if err != nil {
		//	return nil, err
		//}
		//if klw != nil {
		//	tlsConfig.KeyLogWriter = klw
		//}
		//
		//httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

		result[id] = testClientCfg{
			cfgType: "api-ops",
			cfg: &ClientCfg{
				Url:        x,
				apiOpsDcId: toPtr(path.Base(p)),
				// HttpClient: httpClient,
			},
		}
	}
	return result, nil
}
