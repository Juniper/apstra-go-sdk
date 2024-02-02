//go:build integration
// +build integration

package apstra

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
	envAwsInstanceIdList   = "APSTRA_AWS_INSTANCE_IDS"
	awsToplogySecretPrefix = "apstra-info-"
)

type awsTopology struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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

// getAwsTestClientCfgs returns map[string]testClientCfg keyed by
// AWS instance ID
func getAwsTestClientCfgs(ctx context.Context) (map[string]testClientCfg, error) {
	cfg, err := GetTestConfig()
	if err != nil {
		return nil, err
	}

	var topologyIds []string
	if len(cfg.AwsTopologyIds) > 0 {
		topologyIds = make([]string, len(cfg.AwsTopologyIds))
		for i, s := range cfg.AwsTopologyIds {
			topologyIds[i] = s
			i++
		}
	}

	if len(topologyIds) == 0 {
		topologyIds = awsInstanceIdsFromEnv()
	}

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
