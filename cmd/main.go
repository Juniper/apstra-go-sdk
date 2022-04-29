package main

import (
	"fmt"
	aosSdk "github.com/chrismarget-j/apstraTelemetry/aosSdk"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
)

const (
	envApstraStreamHost = "APSTRA_STREAM_HOST"
	envApstraStreamPort = "APSTRA_STREAM_PORT"
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

func client() (*aosSdk.AosClient, error) {
	cfg, err := aosClientClientCfg()
	if err != nil {
		return nil, fmt.Errorf("error getting client config - %v", err)
	}

	aosClient, err := aosSdk.NewAosClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating client - %v", err)
	}

	err = aosClient.Login()
	if err != nil {
		return nil, fmt.Errorf("error logging in AOS client - %v", err)
	}

	return aosClient, nil
}

func listener() (net.Listener, error) {
	streamPort, found := os.LookupEnv(envApstraStreamPort)
	if !found {
		return nil, fmt.Errorf("environment variable '%s' not found", envApstraStreamPort)
	}
	return net.Listen("tcp", ":"+streamPort)
}

//func connHandler(conn net.Conn, errChan chan<- error) {
//	msgLenBuf := make([]byte, 2)
//	buf := bytes.Buffer{}
//	conn.Read(msgLenBuf)
//}

func main() {
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt, os.Kill)

	client, err := client()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := listener()
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	//go func() {
	//	for {
	//		cxn, err := listener.Accept()
	//		if err != nil {
	//			log.Fatalf("socket accept error - %v", err)
	//		}
	//	}
	//}()

	<-quitChan
	os.Exit(0)

	// get streaming configs
	//sc, err := aosClient.GetStreamingConfigs()
	//if err != nil {
	//	log.Fatalf("error getting all streaming configs - %v", err)
	//}
	//err = enc.Encode(sc)
	//if err != nil {
	//	log.Fatalf("error encoding data to JSON - %v", err)
	//}

	err = client.Logout()
	if err != nil {
		log.Fatalf("error logging out in AOS client - %v", err)
	}

	// print the buffer
	//log.Println(buf.String())

	//s := grpc.NewServer()
	//aosST.
}
