//go:build integration
// +build integration

package apstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	envSlicerInstanceIdList = "SLICER_INSTANCE_IDS"
	slicerUserName          = "admin"
	slicerPassword          = "admin"
	slicerURL               = "https://%s"
	slicerTopologyUrlById   = "http://slicer-topology-management-ui.k8s-ci.dc1.apstra.com/v1_1/systest/%s"
)

type slicerTopology struct {
	AosIp string `json:"aos-vm1"`
}

func (o *slicerTopology) getGoapstraClientCfg() (*ClientCfg, error) {
	tlsConfig := &tls.Config{}

	klw, err := keyLogWriterFromEnv(envApstraApiKeyLogFile)
	if err != nil {
		return nil, err
	}
	if klw != nil {
		tlsConfig.KeyLogWriter = klw
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	return &ClientCfg{
		Url:        fmt.Sprintf(slicerURL, o.AosIp),
		User:       slicerUserName,
		Pass:       slicerPassword,
		HttpClient: httpClient,
	}, nil
}

func getSlicerTopology(ctx context.Context, id string) (*slicerTopology, error) {
	requestURL := fmt.Sprintf(slicerTopologyUrlById, id)
	response, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error getting slicerTopology '%s' - %w", id, err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading slicer topology '%s' - %w", id, err)
	}

	var st slicerTopology
	err = json.Unmarshal(body, &st)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling slicer topology '%s' - %w", id, err)
	}

	return &st, nil
}

func slicerInstanceIdsFromEnv() []string {
	ids := os.Getenv(envSlicerInstanceIdList)
	if ids == "" {
		return nil
	}
	return strings.Split(ids, envCloudlabsTopologyIdSep)
}

// getSlicerTestClientCfgs returns map[string]testClientCfg keyed by
// Slicer instance ID
func getSlicerTestClientCfgs(ctx context.Context) (map[string]testClientCfg, error) {
	cfg, err := GetTestConfig()
	if err != nil {
		return nil, err
	}

	var topologyIds []string
	if len(cfg.SlicerTopologyIds) > 0 {
		topologyIds = make([]string, len(cfg.SlicerTopologyIds))
		for i, s := range cfg.SlicerTopologyIds {
			topologyIds[i] = s
			i++
		}
	}

	if len(topologyIds) == 0 {
		topologyIds = slicerInstanceIdsFromEnv()
	}

	result := make(map[string]testClientCfg, len(topologyIds))
	for _, id := range topologyIds {
		topology, err := getSlicerTopology(ctx, id)
		if err != nil {
			return nil, err
		}
		cfg, err := topology.getGoapstraClientCfg()
		if err != nil {
			return nil, err
		}
		result[id] = testClientCfg{
			cfgType: "slicer",
			cfg:     cfg,
		}
	}
	return result, nil
}
