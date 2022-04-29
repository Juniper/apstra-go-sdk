package apstraTelemetry

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func aosUserTestClient1() (*AosClient, error) {
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

func TestAosLogin(t *testing.T) {
	c, err := aosUserTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = c.userLogin()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAosLogout(t *testing.T) {
	c, err := aosUserTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = c.userLogin()
	if err != nil {
		t.Fatal(err)
	}

	err = c.userLogout()
	if err != nil {
		t.Fatal(err)
	}
}
