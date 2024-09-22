// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const testConfigFile = "../.testconfig.hcl"

type TestConfig struct {
	CloudlabsTopologyIds []string `hcl:"cloudlabs_topology_ids,optional"`
	AwsTopologyIds       []string `hcl:"aws_topology_ids,optional"`
	SlicerTopologyIds    []string `hcl:"slicer_topology_ids,optional"`
}

var (
	testCfg      *TestConfig
	testCfgMutex sync.Mutex
)

func GetTestConfig() (TestConfig, error) {
	testCfgMutex.Lock()
	defer testCfgMutex.Unlock()

	if testCfg != nil {
		return *testCfg, nil
	}

	absPath, err := filepath.Abs(testConfigFile)
	if err != nil {
		return TestConfig{}, fmt.Errorf("error expanding config file path %s - %w", testConfigFile, err)
	}

	if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
		return TestConfig{}, nil
	}

	testCfg = new(TestConfig)
	err = hclsimple.DecodeFile(absPath, nil, testCfg)
	if err != nil {
		return *testCfg, fmt.Errorf("failed to parse configuration from %q - %w", absPath, err)
	}

	return *testCfg, nil
}
