package main

import (
	"crypto/tls"
	"fmt"
	aosSdk "github.com/chrismarget-j/apstraTelemetry/aosSdk"
	"log"
	"os"
	"os/signal"
	"strconv"
)

const (
	envApstraStreamHost = "APSTRA_STREAM_HOST"
	envApstraStreamPort = "APSTRA_STREAM_PORT"
)

// getConfigIn includes details for the clientCfg (AOS API
// location+credentials) the required streamingConfig (AOS API configuration to
// point events at our collector) and streamTarget (our listener for AOS
// protobuf messages)
type getConfigIn struct {
	clientCfg       *aosSdk.ClientCfg          // AOS API client
	streamingConfig *aosSdk.StreamingConfigCfg // Specifies target to AOS API
	streamTarget    *aosSdk.AosStreamTargetCfg // Our protobuf stream listener
}

func getConfig(in getConfigIn) error {
	aosScheme, foundAosScheme := os.LookupEnv(aosSdk.EnvApstraScheme)
	aosUser, foundAosUser := os.LookupEnv(aosSdk.EnvApstraUser)
	aosPass, foundAosPass := os.LookupEnv(aosSdk.EnvApstraPass)
	aosHost, foundAosHost := os.LookupEnv(aosSdk.EnvApstraHost)
	aosPort, foundAosPort := os.LookupEnv(aosSdk.EnvApstraPort)
	//recHost, foundRecHost := os.LookupEnv(envApstraStreamHost)
	//recPort, foundRecPort := os.LookupEnv(envApstraStreamPort)

	switch {
	case !foundAosScheme:
		return fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraScheme)
	case !foundAosUser:
		return fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraUser)
	case !foundAosPass:
		return fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraPass)
	case !foundAosHost:
		return fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraHost)
	case !foundAosPort:
		return fmt.Errorf("environment variable '%s' not found", aosSdk.EnvApstraPort)
		//case !foundRecHost:
		//	return fmt.Errorf("environment variable '%s' not found", envApstraStreamHost)
		//case !foundRecPort:
		//	return fmt.Errorf("environment variable '%s' not found", envApstraStreamPort)
	}

	aosPortInt, err := strconv.Atoi(aosPort)
	if err != nil {
		return fmt.Errorf("error converting '%s' to integer - %v", aosPort, err)
	}

	//recPortInt, err := strconv.Atoi(recPort)
	//if err != nil {
	//	return fmt.Errorf("error converting '%s' to integer - %v", recPort, err)
	//}

	in.clientCfg.Scheme = aosScheme
	in.clientCfg.Host = aosHost
	in.clientCfg.Port = uint16(aosPortInt)
	in.clientCfg.User = aosUser
	in.clientCfg.Pass = aosPass
	in.clientCfg.TlsConfig = &tls.Config{InsecureSkipVerify: true}

	in.streamingConfig.StreamingType = aosSdk.StreamingConfigStreamingTypeAlerts
	in.streamingConfig.Protocol = aosSdk.StreamingConfigProtocolProtoBufOverTcp
	in.streamingConfig.SequencingMode = aosSdk.StreamingConfigSequencingModeSequenced
	//in.streamingConfig.Hostname = recHost
	//in.streamingConfig.Port = uint16(recPortInt)

	in.streamTarget.StreamingType = aosSdk.StreamingConfigStreamingTypeAlerts
	in.streamTarget.Protocol = aosSdk.StreamingConfigProtocolProtoBufOverTcp
	//in.streamTarget.SequencingMode = aosSdk.StreamingConfigSequencingModeSequenced
	//in.streamTarget.Port = uint16(recPortInt)

	return nil
}

func main() {
	// handle interrupts
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt, os.Kill)

	// configuration objects
	clientCfg := aosSdk.ClientCfg{}                   // config for interacting with AOS API
	streamingConfig := aosSdk.StreamingConfigCfg{}    // config for pointing event stream at our target
	streamTargetConfig := aosSdk.AosStreamTargetCfg{} // config for our event stream target

	// populate configuration objects using local function
	err := getConfig(getConfigIn{
		clientCfg:       &clientCfg,
		streamingConfig: &streamingConfig,
		streamTarget:    &streamTargetConfig,
	})
	if err != nil {
		log.Fatal(err)
	}

	// create AOS client
	c, err := aosSdk.NewClient(clientCfg)
	if err != nil {
		log.Fatal(err)
	}

	// fetch AOS auth token
	err = c.Login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	//// create AOS stream target service
	//streamTarget, err := aosSdk.NewStreamTarget(&streamTargetConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// start AOS stream target service
	//msgChan, errChan, err := streamTarget.Start()

	//// tell AOS about our stream target (create a StreamingConfig)
	//streamConfigId, err := c.NewStreamingConfig(&streamingConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("Stream target registered with AOS API - ID: %s", string(*streamConfigId))

	streamId1, err := c.NewStreamingConfig(&aosSdk.StreamingConfigCfg{
		StreamingType:  aosSdk.StreamingConfigStreamingTypePerfmon,
		SequencingMode: aosSdk.StreamingConfigSequencingModeSequenced,
		Protocol:       aosSdk.StreamingConfigProtocolProtoBufOverTcp,
		Hostname:       "one.one.one.one",
		Port:           1111,
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(streamId1)
	}

	streamId2, err := c.NewStreamingConfig(&aosSdk.StreamingConfigCfg{
		StreamingType:  aosSdk.StreamingConfigStreamingTypePerfmon,
		SequencingMode: aosSdk.StreamingConfigSequencingModeSequenced,
		Protocol:       aosSdk.StreamingConfigProtocolProtoBufOverTcp,
		Hostname:       "1.1.1.1",
		Port:           1111,
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(streamId2)
	}

	streamId3, err := c.NewStreamingConfig(&aosSdk.StreamingConfigCfg{
		StreamingType:  aosSdk.StreamingConfigStreamingTypePerfmon,
		SequencingMode: aosSdk.StreamingConfigSequencingModeSequenced,
		Protocol:       aosSdk.StreamingConfigProtocolProtoBufOverTcp,
		Hostname:       "1.0.0.1",
		Port:           1111,
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(streamId3)
	}

MAINLOOP:
	for {
		select {
		// interrupt (ctrl-c or whatever)
		case <-quitChan:
			break MAINLOOP
			//// aosStreamTarget has a message
			//case msg := <-msgChan:
			//	log.Println(msg)
			//// aosStreamTarget has an error
			//case err := <-errChan:
			//	log.Fatal(err)
		}
	}

	err = c.DeleteStreamingConfig(streamId1)
	if err != nil {
		log.Println(err)
	}

	err = c.DeleteStreamingConfig(streamId2)
	if err != nil {
		log.Println(err)
	}

	err = c.DeleteStreamingConfig(streamId3)
	if err != nil {
		log.Println(err)
	}

	//streamTarget.Stop()
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

	// print the buffer
	//log.Println(buf.String())

	//s := grpc.NewServer()
	//aosST.
}
