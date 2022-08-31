//go:build integration
// +build integration

package goapstra

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

const (
	cloudlabsTopologyUrlById   = "https://cloudlabs.apstra.com/api/v1.0/topologies/%s"
	envCloudlabsTopologyIdList = "CLOUDLABS_TOPOLOGIES"
	envCloudlabsTopologyIdSep  = ":"

	vSwitchTypeArista = "veos"
	vSwitchTypeNexus  = "nxosv"
	vSwitchTypeSonic  = "sonic-vs"
	vSwitchTypeQfx    = "vqfx"
)

type cloudlabsDeviceType string

func (o cloudlabsDeviceType) platform() AgentPlatform {
	switch o {
	case vSwitchTypeArista:
		return AgentPlatformEOS
	case vSwitchTypeNexus:
		return AgentPlatformNXOS
	case vSwitchTypeQfx:
		return AgentPlatformJunos
	default:
		return AgentPlatformNull
	}
}

type cloudlabsAccess struct {
	Host        string `json:"host"`
	Password    string `json:"password"`
	Port        int    `json:"port"`
	PrivateIp   string `json:"privateIp"`
	PrivatePort int    `json:"privatePort"`
	Protocol    string `json:"protocol"`
	Username    string `json:"username"`
}

type cloudlabsTopology struct {
	Uuid                      string      `json:"uuid"`
	Name                      string      `json:"name"`
	Description               string      `json:"description"`
	CreationTime              string      `json:"creationTime"`
	ExpirationTime            string      `json:"expirationTime"`
	State                     string      `json:"state"`
	Alterations               interface{} `json:"alterations"`
	Owner                     int         `json:"owner"`
	Department                string      `json:"department"`
	EmailNotificationsEnabled bool        `json:"emailNotificationsEnabled"`
	HasBastionHost            bool        `json:"hasBastionHost"`
	Region                    string      `json:"region"`
	StatusMessage             interface{} `json:"statusMessage"`
	Tags                      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"tags"`
	Template int `json:"template"`
	Vms      []struct {
		Access        []cloudlabsAccess   `json:"access"`
		DeviceType    cloudlabsDeviceType `json:"deviceType"`
		Interfaces    []string            `json:"interfaces"`
		Name          string              `json:"name"`
		Role          string              `json:"role"`
		State         string              `json:"state"`
		StatusMessage interface{}         `json:"statusMessage"`
	} `json:"vms"`
}

func (o *cloudlabsTopology) getVmAccessInfo(name string, proto string) (*cloudlabsAccess, error) {
	for _, vm := range o.Vms {
		if vm.Name != name {
			continue
		}
		for _, a := range vm.Access {
			if a.Protocol == proto {
				return &a, nil
			}
		}
		break
	}
	return nil, fmt.Errorf("vm named %s with access protocol %s not found in topology %s", name, proto, o.Uuid)
}

func (o *cloudlabsTopology) selfUpdate() error {
	current, err := getCloudlabsTopology(o.Uuid)
	if err != nil {
		return err
	}

	o.Uuid = current.Uuid
	o.Name = current.Name
	o.Description = current.Description
	o.CreationTime = current.CreationTime
	o.ExpirationTime = current.ExpirationTime
	o.State = current.State
	o.Alterations = current.Alterations
	o.Owner = current.Owner
	o.Department = current.Department
	o.EmailNotificationsEnabled = current.EmailNotificationsEnabled
	o.HasBastionHost = current.HasBastionHost
	o.Region = current.Region
	o.StatusMessage = current.StatusMessage
	o.Tags = current.Tags
	o.Template = current.Template
	o.Vms = current.Vms

	return nil
}

func (o *cloudlabsTopology) getGoapstraClientCfg() (*ClientCfg, error) {
	if o.State != "up" {
		err := o.selfUpdate()
		if err != nil {
			return nil, err
		}
		if o.State != "up" {
			return nil, fmt.Errorf("topology '%s' state is '%s", o.Uuid, o.State)
		}
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	klw, err := keyLogWriterFromEnv(EnvApstraApiKeyLogFile)
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
		Scheme:     access.Protocol,
		Host:       access.Host,
		Port:       uint16(access.Port),
		User:       access.Username,
		Pass:       access.Password,
		HttpClient: httpClient,
	}, nil
}

func (o *cloudlabsTopology) getGoapstraClient() (*Client, error) {
	cfg, err := o.getGoapstraClientCfg()
	if err != nil {
		return nil, err
	}
	return cfg.NewClient()
}

func getCloudlabsTopology(id string) (*cloudlabsTopology, error) {
	topologyUrl, err := url.Parse(fmt.Sprintf(cloudlabsTopologyUrlById, id))
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	klw, err := keyLogWriterFromEnv(EnvApstraApiKeyLogFile)
	if err != nil {
		return nil, err
	}
	if klw != nil {
		tlsConfig.KeyLogWriter = klw
	}

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    topologyUrl,
		Header: http.Header{"Accept": {"application/json"}},
	}
	httpResp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("http %d (%s) - '%s'", httpResp.StatusCode, httpResp.Status, string(body))
	}

	topology := &cloudlabsTopology{}
	return topology, json.NewDecoder(httpResp.Body).Decode(topology)
}

func topologyIdsFromEnv() ([]string, error) {
	env, found := os.LookupEnv(envCloudlabsTopologyIdList)
	if !found {
		return nil, fmt.Errorf("env var '%s' not set", envCloudlabsTopologyIdList)
	}

	return strings.Split(env, envCloudlabsTopologyIdSep), nil
}

type cloudlabsSwitchInfo struct {
	name       string
	role       string
	state      string
	deviceType cloudlabsDeviceType
	sshIp      string
	sshUser    string
	sshPass    string
}

func (o *cloudlabsTopology) getSwitchInfo() ([]cloudlabsSwitchInfo, error) {
	var result []cloudlabsSwitchInfo
	for _, vm := range o.Vms {
		switch vm.DeviceType {
		case vSwitchTypeArista:
		case vSwitchTypeNexus:
		case vSwitchTypeSonic:
		case vSwitchTypeQfx:
		default:
			continue
		}
		access, err := o.getVmAccessInfo(vm.Name, "ssh")
		if err != nil {
			return nil, err
		}
		result = append(result, cloudlabsSwitchInfo{
			name:       vm.Name,
			role:       vm.Role,
			state:      vm.State,
			deviceType: vm.DeviceType,
			sshIp:      access.PrivateIp,
			sshUser:    access.Username,
			sshPass:    access.Password,
		})
	}
	return result, nil
}

type testClient struct {
	clientType string
	client     *Client
}

type testClientCfg struct {
	cfgType string
	cfg     *ClientCfg
}

// getCloudlabsTestClientCfgs returns map[string]testClientCfg keyed by
// cloudlab topology ID
func getCloudlabsTestClientCfgs() (map[string]testClientCfg, error) {
	topologyIds, err := topologyIdsFromEnv()
	if err != nil {
		return nil, err
	}

	result := make(map[string]testClientCfg, len(topologyIds))
	for _, id := range topologyIds {
		topology, err := getCloudlabsTopology(id)
		if err != nil {
			return nil, err
		}
		cfg, err := topology.getGoapstraClientCfg()
		if err != nil {
			return nil, err
		}
		result[id] = testClientCfg{
			cfgType: "cloudlabs",
			cfg:     cfg,
		}
	}
	return result, nil
}

func TestGetCloudlabsTopologies(t *testing.T) {
	topologyIds, err := topologyIdsFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range topologyIds {
		topology, err := getCloudlabsTopology(id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("topology '%s' state: '%s'", id, topology.State)
	}

	bogusId := "bogus"
	_, err = getCloudlabsTopology(bogusId)
	if err == nil {
		t.Fatalf("topology id '%s' did not produce an error", bogusId)
	}
}

func TestGetCloudlabsClients(t *testing.T) {
	topologyIds, err := topologyIdsFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	topologies := make([]*cloudlabsTopology, len(topologyIds))
	for i, id := range topologyIds {
		topology, err := getCloudlabsTopology(id)
		if err != nil {
			t.Fatal(err)
		}
		topologies[i] = topology
		log.Printf("topology '%s' state: '%s'", id, topology.State)
	}

	for _, topology := range topologies {
		client, err := topology.getGoapstraClient()
		if err != nil {
			t.Fatal(err)
		}

		switchInfo, err := topology.getSwitchInfo()
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("Apstra at '%s://%s:%s@%s:%d' has the following switches:\n",
			client.cfg.Scheme, client.cfg.User, client.cfg.Pass, client.cfg.Host, client.cfg.Port)
		for _, si := range switchInfo {
			log.Printf("  %s  \t%s\t %s:%s\t offbox platform: '%s'", si.name, si.sshIp, si.sshUser, si.sshPass, si.deviceType.platform())
		}
	}
}
