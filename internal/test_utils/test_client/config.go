// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build integration && requiretestutils

package testclient

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/stretchr/testify/require"
)

const testConfigFileName = ".testconfig.hcl"

func testConfigFilePath(t testing.TB) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, err = os.Stat(filepath.Join(dir, testConfigFileName)); err == nil {
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find test config file %s", testConfigFileName)
		}
		dir = parent
	}

	return path.Join(dir, testConfigFileName)
}

type TestConfig struct {
	path                 string
	CloudlabsTopologyIds []string `hcl:"cloudlabs_topology_ids,optional"`
	AwsTopologyIds       []string `hcl:"aws_topology_ids,optional"`
	SlicerTopologyIds    []string `hcl:"slicer_topology_ids,optional"`
	ApiOpsProxyUrls      []string `hcl:"api_ops_proxy_urls,optional"`
}

func getTestConfig(t testing.TB) TestConfig {
	t.Helper()

	fileName := testConfigFilePath(t)
	absPath, err := filepath.Abs(fileName)
	require.NoErrorf(t, err, "expanding config file path %s", fileName)

	var testCfg TestConfig
	err = hclsimple.DecodeFile(absPath, nil, &testCfg)
	require.NoErrorf(t, err, "parsing configuraiton file %s", absPath)

	testCfg.path = absPath

	return testCfg
}
