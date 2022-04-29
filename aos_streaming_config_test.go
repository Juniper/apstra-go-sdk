package apstraTelemetry

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

func aosStreamingConfigTestClient1() (*AosClient, error) {
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
		return nil, fmt.Errorf("error converting '%s' to integer - %v", portstr, err)
	}

	return NewAosClient(&AosClientCfg{
		Scheme: scheme,
		Host:   host,
		Port:   uint16(port),
		User:   user,
		Pass:   pass,
	})
}

func TestAosClient_GetAllStreamingConfigs(t *testing.T) {
	client, err := aosStreamingConfigTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = client.UserLogin()
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.GetAllStreamingConfigs()
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	err = pp(response, &out)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(out.String())

}
