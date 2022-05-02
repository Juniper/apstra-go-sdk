package aosSdk

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func userTestClient1() (*Client, error) {
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
		Scheme: scheme,
		Host:   host,
		Port:   uint16(port),
		User:   user,
		Pass:   pass,
	}), nil
}

func TestLogin(t *testing.T) {
	c, err := userTestClient1()
	if err != nil {
		t.Fatal(err)
	}

	err = c.userLogin()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogout(t *testing.T) {
	c, err := userTestClient1()
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
