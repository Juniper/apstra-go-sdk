// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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
	"net/http"
	"os"
	"strings"
	"testing"
)

const (
	cloudlabsTopologyUrlById   = "https://cloudlabs.apstra.com/api/v1.0/topologies/%s"
	envCloudlabsTopologyIdList = "CLOUDLABS_TOPOLOGIES"

	vSwitchTypeArista = "veos"
	vSwitchTypeNexus  = "nxosv"
	vSwitchTypeSonic  = "sonic-vs"
	vSwitchTypeQfx    = "vqfx"
	vSwitchTypeVmx    = "vmx"
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
	case vSwitchTypeVmx:
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

func (o *cloudlabsTopology) selfUpdate(ctx context.Context) error {
	current, err := getCloudlabsTopology(ctx, o.Uuid)
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

func (o *cloudlabsTopology) getClientCfg(ctx context.Context) (*ClientCfg, error) {
	if o.State != "up" {
		err := o.selfUpdate(ctx)
		if err != nil {
			return nil, err
		}
		if o.State != "up" {
			return nil, fmt.Errorf("topology '%s' state is '%s'", o.Uuid, o.State)
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

func (o *cloudlabsTopology) getClient(ctx context.Context) (*Client, error) {
	cfg, err := o.getClientCfg(ctx)
	if err != nil {
		return nil, err
	}
	return cfg.NewClient(ctx)
}

func getCloudlabsTopology(ctx context.Context, id string) (*cloudlabsTopology, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(cloudlabsTopologyUrlById, id), nil)
	if err != nil {
		return nil, fmt.Errorf("error preparing http request for cloudlabs topology - %w", err)
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting cloudlabs topology info - %w", err)
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

func cloudLabsTopologyIdsFromEnv() []string {
	ids := os.Getenv(envCloudlabsTopologyIdList)
	if ids == "" {
		return nil
	}
	return strings.Split(ids, envCloudlabsTopologyIdSep)
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
		case vSwitchTypeVmx:
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

// getCloudlabsTestClientCfgs returns map[string]testClientCfg keyed by
// cloudlab topology ID
func getCloudlabsTestClientCfgs(ctx context.Context) (map[string]testClientCfg, error) {
	cfg, err := GetTestConfig()
	if err != nil {
		return nil, err
	}

	var topologyIds []string
	if len(cfg.CloudlabsTopologyIds) > 0 {
		topologyIds = make([]string, len(cfg.CloudlabsTopologyIds))
		for i, s := range cfg.CloudlabsTopologyIds {
			topologyIds[i] = s
		}
	}

	if len(topologyIds) == 0 {
		topologyIds = cloudLabsTopologyIdsFromEnv()
	}

	result := make(map[string]testClientCfg, len(topologyIds))
	for _, id := range topologyIds {
		topology, err := getCloudlabsTopology(ctx, id)
		if err != nil {
			return nil, err
		}
		cfg, err := topology.getClientCfg(ctx)
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
	ctx := context.Background()
	topologyIds := cloudLabsTopologyIdsFromEnv()

	for _, id := range topologyIds {
		topology, err := getCloudlabsTopology(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("topology '%s' state: '%s'", id, topology.State)
	}

	bogusId := "bogus"
	_, err := getCloudlabsTopology(ctx, bogusId)
	if err == nil {
		t.Fatalf("topology id '%s' did not produce an error", bogusId)
	}
}

func TestGetCloudlabsClients(t *testing.T) {
	ctx := context.Background()
	topologyIds := cloudLabsTopologyIdsFromEnv()

	topologies := make([]*cloudlabsTopology, len(topologyIds))
	for i, id := range topologyIds {
		topology, err := getCloudlabsTopology(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		topologies[i] = topology
		log.Printf("topology '%s' state: '%s'", id, topology.State)
	}

	for _, topology := range topologies {
		client, err := topology.getClient(ctx)
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
