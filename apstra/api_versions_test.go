package apstra

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

func clientTestVersionsCfg1() (*ClientCfg, error) {
	user, foundUser := os.LookupEnv(EnvApstraUser)
	pass, foundPass := os.LookupEnv(EnvApstraPass)
	scheme, foundScheme := os.LookupEnv(EnvApstraScheme)
	host, foundHost := os.LookupEnv(EnvApstraHost)
	portstr, foundPort := os.LookupEnv(EnvApstraPort)

	switch {
	case !foundUser:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraUser)
	case !foundPass:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraPass)
	case !foundScheme:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraScheme)
	case !foundHost:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraHost)
	case !foundPort:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraPort)
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, fmt.Errorf("error converting '%s' to integer - %w", portstr, err)
	}

	return &ClientCfg{
		Scheme:    scheme,
		Host:      host,
		Port:      uint16(port),
		User:      user,
		Pass:      pass,
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	}, nil
}

func TestGetVersionsServer(t *testing.T) {
	cfg, err := clientTestVersionsCfg1()
	if err != nil {
		t.Fatal(err)
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Logout(context.TODO())

	response, err := client.getVersionsServer(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(string(body))
}
