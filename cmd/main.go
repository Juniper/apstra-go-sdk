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
	envApstraStreamHost     = "APSTRA_STREAM_HOST"
	envApstraStreamBasePort = "APSTRA_STREAM_BASE_PORT"
)

// getConfigIn includes details for the clientCfg (AOS API
// location+credentials) the required streamingConfig (AOS API configuration to
// point events at our collector) and streamTarget (our listener for AOS
// protobuf messages)
type getConfigIn struct {
	clientCfg       *aosSdk.ClientCfg           // AOS API client config
	streamTarget    []aosSdk.StreamTargetCfg    // Our protobuf stream listener
	streamingConfig []aosSdk.StreamingConfigCfg // Tell AOS API about our stream listener
}

func getConfig(in getConfigIn) error {
	aosScheme, foundAosScheme := os.LookupEnv(aosSdk.EnvApstraScheme)
	aosUser, foundAosUser := os.LookupEnv(aosSdk.EnvApstraUser)
	aosPass, foundAosPass := os.LookupEnv(aosSdk.EnvApstraPass)
	aosHost, foundAosHost := os.LookupEnv(aosSdk.EnvApstraHost)
	aosPort, foundAosPort := os.LookupEnv(aosSdk.EnvApstraPort)
	recHost, foundRecHost := os.LookupEnv(envApstraStreamHost)
	recPort, foundRecPort := os.LookupEnv(envApstraStreamBasePort)

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
	case !foundRecHost:
		return fmt.Errorf("environment variable '%s' not found", envApstraStreamHost)
	case !foundRecPort:
		return fmt.Errorf("environment variable '%s' not found", envApstraStreamBasePort)
	}

	aosPortInt, err := strconv.Atoi(aosPort)
	if err != nil {
		return fmt.Errorf("error converting '%s' to integer - %v", aosPort, err)
	}

	recPortInt, err := strconv.Atoi(recPort)
	if err != nil {
		return fmt.Errorf("error converting '%s' to integer - %v", recPort, err)
	}

	in.clientCfg.Scheme = aosScheme
	in.clientCfg.Host = aosHost
	in.clientCfg.Port = uint16(aosPortInt)
	in.clientCfg.User = aosUser
	in.clientCfg.Pass = aosPass
	in.clientCfg.TlsConfig = tls.Config{InsecureSkipVerify: true} // todo: something less shameful

	for i := range in.streamTarget {
		in.streamTarget[i].StreamingType = aosSdk.StreamingConfigStreamingType(aosSdk.StreamingConfigStreamingTypeUnknown + 1)
		in.streamTarget[i].Protocol = aosSdk.StreamingConfigProtocolProtoBufOverTcp
		in.streamTarget[i].SequencingMode = aosSdk.StreamingConfigSequencingModeUnsequenced
		in.streamTarget[i].Port = uint16(recPortInt + i)
	}

	for i := range in.streamingConfig {
		in.streamingConfig[i].StreamingType = aosSdk.StreamingConfigStreamingType(aosSdk.StreamingConfigStreamingTypeUnknown + 1)
		in.streamingConfig[i].Protocol = aosSdk.StreamingConfigProtocolProtoBufOverTcp
		in.streamingConfig[i].SequencingMode = aosSdk.StreamingConfigSequencingModeUnsequenced
		in.streamingConfig[i].Hostname = recHost
		in.streamingConfig[i].Port = uint16(recPortInt + i)
	}

	return nil
}

func main() {
	// handle interrupts
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt, os.Kill)

	// configuration objects
	clientCfg := aosSdk.ClientCfg{}                          // config for interacting with AOS API
	streamingConfigs := make([]aosSdk.StreamingConfigCfg, 3) // config for pointing event stream at our target
	streamTargetConfigs := make([]aosSdk.StreamTargetCfg, 3) // config for our event stream target

	// populate configuration objects using local function
	err := getConfig(getConfigIn{
		clientCfg:       &clientCfg,
		streamingConfig: streamingConfigs,
		streamTarget:    streamTargetConfigs,
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

	// create aggregator channels where we'll get messages from all target services
	msgChan := make(chan *aosSdk.StreamingMessage)
	errChan := make(chan error)

	var streamTargets []aosSdk.StreamTarget
	//var msgChans []<-chan *aosStreaming.AosMessage
	//var errChans []<-chan error
	for i := range streamTargetConfigs {
		// create each AOS stream target service
		st, err := aosSdk.NewStreamTarget(streamTargetConfigs[i])
		if err != nil {
			log.Fatal(err)
		}
		streamTargets = append(streamTargets,

			*st)

		// start each AOS stream target service
		mc, ec, err := st.Start()
		if err != nil {
			log.Fatal(err)
		}

		// copy messages from this target's message channel to aggregated message channel
		go func(in <-chan *aosSdk.StreamingMessage, out chan<- *aosSdk.StreamingMessage) {
			for {
				out <- <-in
			}
		}(mc, msgChan)

		// copy errors from this target's error channel to aggregated error channel
		go func(in <-chan error, out chan<- error) {
			for {
				out <- <-in
			}
		}(ec, errChan)
		//msgChans = append(msgChans, msgChan)
		//errChans = append(errChans, errChan)
	}

	// tell AOS about our stream targets (create StreamingConfig objects)
	var streamingConfigIds []aosSdk.StreamingConfigId
	for i := range streamingConfigs {
		id, err := c.NewStreamingConfig(&streamingConfigs[i])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Stream target registered with AOS API - ID: %s", string(id))
		streamingConfigIds = append(streamingConfigIds, id)
	}

	//streamId1, err := c.NewStreamingConfig(&aosSdk.StreamingConfigCfg{
	//	StreamingType:  aosSdk.StreamingConfigStreamingTypePerfmon,
	//	SequencingMode: aosSdk.StreamingConfigSequencingModeSequenced,
	//	Protocol:       aosSdk.StreamingConfigProtocolProtoBufOverTcp,
	//	Hostname:       "one.one.one.one",
	//	Port:           1111,
	//})
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(streamId1)
	//}

	//streamId2, err := c.NewStreamingConfig(&aosSdk.StreamingConfigCfg{
	//	StreamingType:  aosSdk.StreamingConfigStreamingTypePerfmon,
	//	SequencingMode: aosSdk.StreamingConfigSequencingModeSequenced,
	//	Protocol:       aosSdk.StreamingConfigProtocolProtoBufOverTcp,
	//	Hostname:       "1.1.1.1",
	//	Port:           1111,
	//})
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(streamId2)
	//}

	//streamId3, err := c.NewStreamingConfig(&aosSdk.StreamingConfigCfg{
	//	StreamingType:  aosSdk.StreamingConfigStreamingTypePerfmon,
	//	SequencingMode: aosSdk.StreamingConfigSequencingModeSequenced,
	//	Protocol:       aosSdk.StreamingConfigProtocolProtoBufOverTcp,
	//	Hostname:       "1.0.0.1",
	//	Port:           1111,
	//})
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(streamId3)
	//}

MAINLOOP:
	for {
		select {
		// interrupt (ctrl-c or whatever)
		case <-quitChan:
			break MAINLOOP
		// aosStreamTarget has a message
		case msg := <-msgChan:
			var seqNumStr string
			switch msg.SequencingMode {
			case aosSdk.StreamingConfigSequencingModeSequenced:
				seqNumStr = strconv.Itoa(int(*msg.SequenceNum))
			case aosSdk.StreamingConfigSequencingModeUnsequenced:
				seqNumStr = "N/A"
			}
			log.Printf("%s / %s / message number %s / %s\n", msg.StreamingType, msg.SequencingMode, seqNumStr, msg.Message)
		// aosStreamTarget has an error
		case err := <-errChan:
			log.Fatal(err)
		}
	}

	for _, id := range streamingConfigIds {
		log.Printf("deleting stream id %s\n", id)
		err := c.DeleteStreamingConfig(id)
		if err != nil {
			log.Println(err)
		}
	}

	//err = c.DeleteStreamingConfig(streamId1)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//err = c.DeleteStreamingConfig(streamId2)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//err = c.DeleteStreamingConfig(streamId3)
	//if err != nil {
	//	log.Println(err)
	//}

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
