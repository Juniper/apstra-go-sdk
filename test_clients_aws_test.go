//go:build integration
// +build integration

package goapstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	//cloudlabsTopologyUrlById   = "https://cloudlabs.apstra.com/api/v1.0/topologies/%s"
	//envApstraApiKeyLogFile     = "APSTRA_API_TLS_LOGFILE"
	envAwsInstanceIdList   = "APSTRA_AWS_INSTANCE_IDS"
	awsToplogySecretPrefix = "apstra-info-"
	//envCloudlabsTopologyIdSep  = ":"
	//
	//vSwitchTypeArista = "veos"
	//vSwitchTypeNexus  = "nxosv"
	//vSwitchTypeSonic  = "sonic-vs"
	//vSwitchTypeQfx    = "vqfx"
)

//type cloudlabsDeviceType string

//func (o cloudlabsDeviceType) platform() AgentPlatform {
//	switch o {
//	case vSwitchTypeArista:
//		return AgentPlatformEOS
//	case vSwitchTypeNexus:
//		return AgentPlatformNXOS
//	case vSwitchTypeQfx:
//		return AgentPlatformJunos
//	default:
//		return AgentPlatformNull
//	}
//}

//type cloudlabsAccess struct {
//	Host        string `json:"host"`
//	Password    string `json:"password"`
//	Port        int    `json:"port"`
//	PrivateIp   string `json:"privateIp"`
//	PrivatePort int    `json:"privatePort"`
//	Protocol    string `json:"protocol"`
//	Username    string `json:"username"`
//}

type awsTopology struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//func (o *cloudlabsTopology) getVmAccessInfo(name string, proto string) (*cloudlabsAccess, error) {
//	for _, vm := range o.Vms {
//		if vm.Name != name {
//			continue
//		}
//		for _, a := range vm.Access {
//			if a.Protocol == proto {
//				return &a, nil
//			}
//		}
//		break
//	}
//	return nil, fmt.Errorf("vm named %s with access protocol %s not found in topology %s", name, proto, o.Uuid)
//}

//func (o *cloudlabsTopology) selfUpdate() error {
//	current, err := getCloudlabsTopology(o.Uuid)
//	if err != nil {
//		return err
//	}
//
//	o.Uuid = current.Uuid
//	o.Name = current.Name
//	o.Description = current.Description
//	o.CreationTime = current.CreationTime
//	o.ExpirationTime = current.ExpirationTime
//	o.State = current.State
//	o.Alterations = current.Alterations
//	o.Owner = current.Owner
//	o.Department = current.Department
//	o.EmailNotificationsEnabled = current.EmailNotificationsEnabled
//	o.HasBastionHost = current.HasBastionHost
//	o.Region = current.Region
//	o.StatusMessage = current.StatusMessage
//	o.Tags = current.Tags
//	o.Template = current.Template
//	o.Vms = current.Vms
//
//	return nil
//}

func (o *awsTopology) getGoapstraClientCfg() (*ClientCfg, error) {
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
		Url:        o.Url,
		User:       o.Username,
		Pass:       o.Password,
		HttpClient: httpClient,
	}, nil
}

//func (o *cloudlabsTopology) getGoapstraClient() (*Client, error) {
//	cfg, err := o.getGoapstraClientCfg()
//	if err != nil {
//		return nil, err
//	}
//	return cfg.NewClient()
//}

func getAwsTopology(ctx context.Context, id string) (*awsTopology, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error loading default AWS config - %w", err)
	}

	sm := secretsmanager.NewFromConfig(awsCfg)
	secretId := fmt.Sprintf("%s%s", awsToplogySecretPrefix, id)
	gsvi := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}
	gsvo, err := sm.GetSecretValue(ctx, gsvi)
	if err != nil {
		return nil, fmt.Errorf("error getting secret '%s' value - %w", secretId, err)
	}
	if gsvo.SecretString == nil {
		return nil, fmt.Errorf("secret '%s' value is nil - %w", secretId, err)
	}

	result := &awsTopology{}
	return result, json.NewDecoder(strings.NewReader(*gsvo.SecretString)).Decode(result)
}

func awsInstanceIdsFromEnv() []string {
	ids := os.Getenv(envAwsInstanceIdList)
	if ids == "" {
		return nil
	}
	return strings.Split(ids, envCloudlabsTopologyIdSep)
}

//type cloudlabsSwitchInfo struct {
//	name       string
//	role       string
//	state      string
//	deviceType cloudlabsDeviceType
//	sshIp      string
//	sshUser    string
//	sshPass    string
//}

//func (o *cloudlabsTopology) getSwitchInfo() ([]cloudlabsSwitchInfo, error) {
//	var result []cloudlabsSwitchInfo
//	for _, vm := range o.Vms {
//		switch vm.DeviceType {
//		case vSwitchTypeArista:
//		case vSwitchTypeNexus:
//		case vSwitchTypeSonic:
//		case vSwitchTypeQfx:
//		default:
//			continue
//		}
//		access, err := o.getVmAccessInfo(vm.Name, "ssh")
//		if err != nil {
//			return nil, err
//		}
//		result = append(result, cloudlabsSwitchInfo{
//			name:       vm.Name,
//			role:       vm.Role,
//			state:      vm.State,
//			deviceType: vm.DeviceType,
//			sshIp:      access.PrivateIp,
//			sshUser:    access.Username,
//			sshPass:    access.Password,
//		})
//	}
//	return result, nil
//}

// getAwsTestClientCfgs returns map[string]testClientCfg keyed by
// AWS instance ID
func getAwsTestClientCfgs(ctx context.Context) (map[string]testClientCfg, error) {
	topologyIds := awsInstanceIdsFromEnv()

	result := make(map[string]testClientCfg, len(topologyIds))
	for _, id := range topologyIds {
		topology, err := getAwsTopology(ctx, id)
		if err != nil {
			return nil, err
		}
		cfg, err := topology.getGoapstraClientCfg()
		if err != nil {
			return nil, err
		}
		result[id] = testClientCfg{
			cfgType: "aws",
			cfg:     cfg,
		}
	}
	return result, nil
}

//func TestGetCloudlabsTopologies(t *testing.T) {
//	topologyIds, err := cloudLabsTopologyIdsFromEnv()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for _, id := range topologyIds {
//		topology, err := getCloudlabsTopology(id)
//		if err != nil {
//			t.Fatal(err)
//		}
//		log.Printf("topology '%s' state: '%s'", id, topology.State)
//	}
//
//	bogusId := "bogus"
//	_, err = getCloudlabsTopology(bogusId)
//	if err == nil {
//		t.Fatalf("topology id '%s' did not produce an error", bogusId)
//	}
//}
//
//func TestGetCloudlabsClients(t *testing.T) {
//	topologyIds, err := cloudLabsTopologyIdsFromEnv()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	topologies := make([]*cloudlabsTopology, len(topologyIds))
//	for i, id := range topologyIds {
//		topology, err := getCloudlabsTopology(id)
//		if err != nil {
//			t.Fatal(err)
//		}
//		topologies[i] = topology
//		log.Printf("topology '%s' state: '%s'", id, topology.State)
//	}
//
//	for _, topology := range topologies {
//		client, err := topology.getGoapstraClient()
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		switchInfo, err := topology.getSwitchInfo()
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		log.Printf("Apstra at '%s' has the following switches:\n", client.cfg.Url)
//		for _, si := range switchInfo {
//			log.Printf("  %s  \t%s\t %s:%s\t offbox platform: '%s'", si.name, si.sshIp, si.sshUser, si.sshPass, si.deviceType.platform())
//		}
//	}
//}
