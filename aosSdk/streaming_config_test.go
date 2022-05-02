package aosSdk

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

func streamingConfigTestClient1() (*Client, error) {
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

	return NewClient(ClientCfg{
		Scheme:    scheme,
		Host:      host,
		Port:      uint16(port),
		User:      user,
		Pass:      pass,
		TlsConfig: tls.Config{InsecureSkipVerify: true},
	}), nil
}

func TestThing(t *testing.T) {
	var mySCST StreamingConfigStreamingType
	log.Println(mySCST)

}

func TestClient_GetAllStreamingConfigs(t *testing.T) {
	client, err := streamingConfigTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Login()
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.getAllStreamingConfigs()
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
