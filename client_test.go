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

func getTestClientsAndMockAPIs() (map[string]*Client, map[string]*mockApstraApi, error) {
	clientResult := make(map[string]*Client)
	apiResult := make(map[string]*mockApstraApi)

	if useLiveClient() {
		log.Println("generating a live client")
		c, err := newLiveTestClient()
		if err != nil {
			return nil, nil, err
		}
		clientResult["live"] = c
		apiResult["live"] = nil
	}

	if useMockClient() {
		log.Println("generating a mock client")
		c, api, err := newMockTestClient()
		if err != nil {
			return nil, nil, err
		}
		clientResult["mock"] = c
		apiResult["mock"] = api
	}

	return clientResult, apiResult, nil
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

func newMockTestClient() (*Client, *mockApstraApi, error) {
	c, err := NewClient(&ClientCfg{
		Scheme: "mock",
		Host:   "mock",
		Port:   uint16(0),
		User:   mockApstraUser,
		Pass:   mockApstraPass,
	})
	if err != nil {
		return nil, nil, err
	}

	mockApi, err := newMockApstraApi(mockApstraPass)
	if err != nil {
		return nil, nil, err
	}

	c.httpClient = mockApi

	return c, mockApi, nil
}

func TestLoginLogout(t *testing.T) {
	clients, _, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("testing with %d clients", len(clients))

	for ctype, c := range clients {
		log.Printf("testing Login() with %s client", ctype)
		err = c.Login(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("testing Logout() with %s client", ctype)
		err = c.Logout(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestLoginLogoutAuthFail(t *testing.T) {
	clients, _, err := getTestClientsAndMockAPIs()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("testing with %d clients", len(clients))

	for clientName, client := range clients {
		log.Printf("testing for Login() fail with %s client", clientName)
		client.cfg.Pass = randString(10, "hex")
		err := client.Login(context.TODO())
		if err == nil {
			log.Fatal(fmt.Errorf("tried logging in with bad password, did not get errror"))
		}

		log.Printf("testing for Logout() fail with %s client", clientName)
		client.httpHeaders[apstraAuthHeader] = randJwt()
		err = client.Logout(context.TODO())
		if err == nil {
			log.Fatal(fmt.Errorf("tried logging in with bad password, did not get errror"))
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
