//go:build integration
// +build integration

package apstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

const (
	slicerTopologyUrlById   = "http://slicer-topology-management-ui.k8s-ci.dc1.apstra.com/v1_1/systest/%s"
	envSlicerTopologyIdList = "SLICER_TOPOLOGIES"
)

type slicerDeviceType string

func (o slicerDeviceType) platform() AgentPlatform {
	switch o {
	case vSwitchTypeArista:
		return AgentPlatformEOS
	case vSwitchTypeNexus:
		return AgentPlatformNXOS
	case vSwitchTypeQfx:
		return AgentPlatformJunos
	case vSwitchTypeVmx:
		return AgentPlatformJunos
	default:
		return AgentPlatformNull
	}
}

type slicerAccess struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Username string `json:"username"`
}

type slicerTopology struct {
	id           string
	DeployStatus string `json:"deploy_status"`
	DeployModel  struct {
		DutmgmtConnectivity map[string]string `json:"dutmgmt_connectivity"`
	} `json:"deploy_model"`
}

func (o *slicerTopology) getVmAccessInfo(name string, proto string) (*slicerAccess, error) {
	port, err := net.LookupPort("tcp", proto)
	if err != nil {
		return nil, err
	}

	for dutName, dutIpStr := range o.DeployModel.DutmgmtConnectivity {
		if dutName == name {
			return &slicerAccess{
				Host:     dutIpStr,
				Password: "admin",
				Port:     port,
				Protocol: proto,
				Username: "admin",
			}, nil
		}
	}
	return nil, fmt.Errorf("vm named %s not found in topology %s", name, o.id)
}

func (o *slicerTopology) selfUpdate(ctx context.Context) error {
	current, err := getSlicerTopology(ctx, o.id)
	if err != nil {
		return err
	}

	o.DeployStatus = current.DeployStatus
	o.DeployModel = current.DeployModel

	return nil
}

func (o *slicerTopology) getClientCfg(ctx context.Context) (*ClientCfg, error) {
	deployedRegex := regexp.MustCompile("^(.*_)?[_]?deployed[_]?(_.*)?$")
	if !deployedRegex.MatchString(o.DeployStatus) {
		err := o.selfUpdate(ctx)
		if err != nil {
			return nil, err
		}
		if !deployedRegex.MatchString(o.DeployStatus) {
			return nil, fmt.Errorf("topology '%s' deploy status is '%s'", o.id, o.DeployStatus)
		}
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	klw, err := keyLogWriterFromEnv(envApstraApiKeyLogFile)
	if err != nil {
		return nil, err
	}
	if klw != nil {
		tlsConfig.KeyLogWriter = klw
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	access, err := o.getVmAccessInfo("aos-vm1", "https")
	if err != nil {
		return nil, err
	}

	return &ClientCfg{
		Url:        fmt.Sprintf("%s://%s:%d", access.Protocol, access.Host, access.Port),
		User:       access.Username,
		Pass:       access.Password,
		HttpClient: httpClient,
	}, nil
}

func (o *slicerTopology) getGoapstraClient(ctx context.Context) (*Client, error) {
	cfg, err := o.getClientCfg(ctx)
	if err != nil {
		return nil, err
	}
	return cfg.NewClient(ctx)
}

func getSlicerTopology(ctx context.Context, id string) (*slicerTopology, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(slicerTopologyUrlById, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error preparing http request for slicer topology - %w", err)
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting slicer topology info - %w", err)
	}

	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("http %d (%s) - '%s'", httpResp.StatusCode, httpResp.Status, string(body))
	}

	topology := &slicerTopology{}
	return topology, json.NewDecoder(httpResp.Body).Decode(topology)
}

func slicerTopologyIdsFromEnv() []string {
	ids := os.Getenv(envSlicerTopologyIdList)
	if ids == "" {
		return nil
	}
	return strings.Split(ids, envCloudlabsTopologyIdSep)
}

type slicerSwitchInfo struct {
	name       string
	role       string
	state      string
	deviceType slicerDeviceType
	sshIp      string
	sshUser    string
	sshPass    string
}

func (o *slicerTopology) getSwitchInfo() ([]slicerSwitchInfo, error) {
	// todo get the switch info from slicer somehow
	return nil, nil
}

// getSlicerTestClientCfgs returns map[string]testClientCfg keyed by
// slicer topology ID
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
		}
	}

	if len(topologyIds) == 0 {
		topologyIds = slicerTopologyIdsFromEnv()
	}

	result := make(map[string]testClientCfg, len(topologyIds))
	for _, id := range topologyIds {
		topology, err := getSlicerTopology(ctx, id)
		if err != nil {
			return nil, err
		}
		cfg, err := topology.getClientCfg(ctx)
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

func TestGetSlicerTopologies(t *testing.T) {
	ctx := context.Background()
	topologyIds := slicerTopologyIdsFromEnv()

	for _, id := range topologyIds {
		topology, err := getSlicerTopology(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("topology '%s' deploy status: '%s'", id, topology.DeployStatus)
	}

	bogusId := "bogus"
	_, err := getSlicerTopology(ctx, bogusId)
	if err == nil {
		t.Fatalf("topology id '%s' did not produce an error", bogusId)
	}
}

func TestGetSlicerClients(t *testing.T) {
	ctx := context.Background()
	topologyIds := slicerTopologyIdsFromEnv()

	topologies := make([]*slicerTopology, len(topologyIds))
	for i, id := range topologyIds {
		topology, err := getSlicerTopology(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		topologies[i] = topology
		log.Printf("topology '%s' deploy status: '%s'", id, topology.DeployStatus)
	}

	for _, topology := range topologies {
		client, err := topology.getGoapstraClient(ctx)
		if err != nil {
			t.Fatal(err)
		}

		switchInfo, err := topology.getSwitchInfo()
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("Apstra at '%s' has the following switches:\n", client.cfg.Url)
		for _, si := range switchInfo {
			log.Printf("  %s  \t%s\t %s:%s\t offbox platform: '%s'", si.name, si.sshIp, si.sshUser, si.sshPass, si.deviceType.platform())
		}
	}
}
