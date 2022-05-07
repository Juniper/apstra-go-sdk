package aosSdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func blueprintsTestClient1() (*Client, error) {
	user, foundUser := os.LookupEnv(EnvApstraUser)
	pass, foundPass := os.LookupEnv(EnvApstraPass)
	scheme, foundScheme := os.LookupEnv(EnvApstraScheme)
	host, foundHost := os.LookupEnv(EnvApstraHost)
	portstr, foundPort := os.LookupEnv(EnvApstraPort)
	keyLogFile, foundKeyLogFile := os.LookupEnv(EnvApstraApiKeyLogFile)

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

	var kl io.Writer
	var err error
	if foundKeyLogFile {
		kl, err = keyLogWriter(keyLogFile)
		if err != nil {
			return nil, fmt.Errorf("error creating keyLogWriter - %w", err)
		}
	} else {
		kl = nil
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, fmt.Errorf("error converting '%s' to integer - %w", portstr, err)
	}

	return NewClient(&ClientCfg{
		Scheme:    scheme,
		Host:      host,
		Port:      uint16(port),
		User:      user,
		Pass:      pass,
		TlsConfig: tls.Config{InsecureSkipVerify: true, KeyLogWriter: kl},
		Timeout:   5 * time.Minute,
	}), nil
}

func TestGetBlueprints(t *testing.T) {
	client, err := blueprintsTestClient1()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Logout()

	blueprints, err := client.GetBlueprints()
	if err != nil {
		t.Fatal(err)
	}

	result, err := json.Marshal(blueprints)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(result))
}

func TestCreateRoutingZone(t *testing.T) {
	client, err := blueprintsTestClient1()
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.createRoutingZone("db10754a-610e-475b-9baa-4c85f82282e8", &CreateRoutingZoneCfg{
		SzType:  "evpn",
		VrfName: "test",
		Label:   "test_label",
	})

	buf := bytes.Buffer{}
	pp(result, &buf)
	log.Print(buf.String())
}
