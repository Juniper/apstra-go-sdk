package goapstra

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
)

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

func TestLoginLogout(t *testing.T) {
	client, err := newLiveTestClient()
	//client, err := newMockTestClient()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Login(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Logout(context.TODO())
	if err != nil {
		log.Fatal(err)
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
