package goapstra

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

func getTestClients() (map[string]*Client, error) {
	result := make(map[string]*Client)

	if useLiveClient() {
		log.Println("generating a live client")
		c, err := newLiveTestClient()
		if err != nil {
			return nil, err
		}
		result["live"] = c
	}

	if useMockClient() {
		log.Println("generating a mock client")
		c, err := newMockTestClient()
		if err != nil {
			return nil, err
		}
		result["mock"] = c
	}

	return result, nil
}

func useLiveClient() bool {
	_, err := os.Stat("/tmp/live")
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func useMockClient() bool {
	_, err := os.Stat("/tmp/nomock")
	if errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}

func newLiveTestClient() (*Client, error) {
	user, foundUser := os.LookupEnv(EnvApstraUser)
	pass, foundPass := os.LookupEnv(EnvApstraPass)
	scheme, foundScheme := os.LookupEnv(EnvApstraScheme)
	host, foundHost := os.LookupEnv(EnvApstraHost)
	portStr, foundPort := os.LookupEnv(EnvApstraPort)

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

	kl, err := keyLogWriter(EnvApstraApiKeyLogFile)
	if err != nil {
		return nil, fmt.Errorf("error creating keyLogWriter - %w", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("error converting '%s' to integer - %w", portStr, err)
	}

	return NewClient(&ClientCfg{
		Scheme:    scheme,
		Host:      host,
		Port:      uint16(port),
		User:      user,
		Pass:      pass,
		TlsConfig: &tls.Config{InsecureSkipVerify: true, KeyLogWriter: kl},
	})
}

func newMockTestClient() (*Client, error) {
	c, err := NewClient(&ClientCfg{
		Scheme: "mock",
		Host:   "mock",
		Port:   uint16(0),
		User:   "mockUser",
		Pass:   "mockPass",
	})
	if err != nil {
		return nil, err
	}

	c.httpClient = &mockApstraApi{
		username: "mockUser",
		password: "mockPass",
	}

	return c, err
}

func TestLoginLogout(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("testing with %d clients", len(clients))

	for t, c := range clients {
		log.Printf("testing Login() with %s client", t)
		err = c.Login(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("testing Logout() with %s client", t)
		err = c.Logout(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
}

//func TestParseBytesAsTaskId(t *testing.T) {
//	var testData [][]byte
//	var expected []bool
//
//	testData = append(testData, []byte(""))
//	expected = append(expected, false)
//
//	testData = append(testData, []byte("{}"))
//	expected = append(expected, false)
//
//	testData = append(testData, []byte("[]"))
//	expected = append(expected, false)
//
//	if len(testData) != len(expected) {
//		t.Fatalf("test setup error - have %d tests, but expect %d results", len(testData), len(expected))
//	}
//
//	for i, td := range testData {
//		result := &taskIdResponse{}
//		ok, err := peekParseResponseBodyAsTaskId(td, result)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if ok != expected[i] {
//			t.Fatalf("test data '%s' produced '%t', expected '%t'", string(td), ok, expected[i])
//		}
//	}
//}
