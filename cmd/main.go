package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	aosSdk "github.com/chrismarget-j/apstraTelemetry/aosSdk"
	"log"
	"os"
	"strconv"
)

func aosClientClientCfg() (*aosSdk.AosClientCfg, error) {
	user, foundUser := os.LookupEnv(aosSdk.EnvApstraUser)
	pass, foundPass := os.LookupEnv(aosSdk.EnvApstraPass)
	scheme, foundScheme := os.LookupEnv(aosSdk.EnvApstraScheme)
	host, foundHost := os.LookupEnv(aosSdk.EnvApstraHost)
	portstr, foundPort := os.LookupEnv(aosSdk.EnvApstraPort)

	switch {
	case !foundUser:
		return nil, fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraUser)
	case !foundPass:
		return nil, fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraPass)
	case !foundScheme:
		return nil, fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraScheme)
	case !foundHost:
		return nil, fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraHost)
	case !foundPort:
		return nil, fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraPort)
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, fmt.Errorf("error converting '%s' to integer - %v", portstr, err)
	}

	return &aosSdk.AosClientCfg{
		Scheme: scheme,
		Host:   host,
		Port:   uint16(port),
		User:   user,
		Pass:   pass,
	}, nil
}

func main() {
	cfg, err := aosClientClientCfg()
	if err != nil {
		log.Fatal(err)
	}

	aosClient, err := aosSdk.NewAosClient(cfg)
	if err != nil {
		log.Fatalf("error creating new AOS client - %v", err)
	}

	err = aosClient.Login()
	if err != nil {
		log.Fatalf("error logging in AOS client - %v", err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "    ")

	// get version
	ver, err := aosClient.GetVersion()
	if err != nil {
		log.Fatalf("error getting AOS version - %v", err)
	}
	err = enc.Encode(ver)
	if err != nil {
		log.Fatalf("error encoding data to JSON - %v", err)
	}

	// get streaming configs
	sc, err := aosClient.GetStreamingConfigs()
	if err != nil {
		log.Fatalf("error getting all streaming configs - %v", err)
	}
	err = enc.Encode(sc)
	if err != nil {
		log.Fatalf("error encoding data to JSON - %v", err)
	}

	err = aosClient.Logout()
	if err != nil {
		log.Fatalf("error logging out in AOS client - %v", err)
	}

	// print the buffer
	log.Println(buf.String())

	//s := grpc.NewServer()
	//aosST.
}
