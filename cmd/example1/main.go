package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/chrismarget-j/apstraTelemetry/apstra"
)

// getConfigIn includes details for the clientCfg (AOS API
// location+credentials) the required streamingConfigParams (AOS API configuration to
// point events at our collector) and streamTargetCfg (our listener for AOS
// protobuf messages)
type getConfigIn struct {
	clientCfg             *apstra.ClientCfg              // AOS API client config
	streamTargetCfg       []apstra.StreamTargetCfg       // Our protobuf stream listener
	streamingConfigParams []apstra.StreamingConfigParams // Tell AOS API about our stream listener
}

func getConfig(in *getConfigIn) error {
	aosScheme, foundAosScheme := os.LookupEnv(apstra.EnvApstraScheme)
	aosUser, foundAosUser := os.LookupEnv(apstra.EnvApstraUser)
	aosPass, foundAosPass := os.LookupEnv(apstra.EnvApstraPass)
	aosHost, foundAosHost := os.LookupEnv(apstra.EnvApstraHost)
	aosPort, foundAosPort := os.LookupEnv(apstra.EnvApstraPort)
	recHost, foundRecHost := os.LookupEnv(apstra.EnvApstraStreamHost)
	recPort, foundRecPort := os.LookupEnv(apstra.EnvApstraStreamBasePort)

	switch {
	case !foundAosScheme:
		return fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraScheme)
	case !foundAosUser:
		return fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraUser)
	case !foundAosPass:
		return fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraPass)
	case !foundAosHost:
		return fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraHost)
	case !foundAosPort:
		return fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraPort)
	case !foundRecHost:
		return fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraStreamHost)
	case !foundRecPort:
		return fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraStreamBasePort)
	}

	aosPortInt, err := strconv.Atoi(aosPort)
	if err != nil {
		return fmt.Errorf("error converting '%s' to integer - %w", aosPort, err)
	}

	recPortInt, err := strconv.Atoi(recPort)
	if err != nil {
		return fmt.Errorf("error converting '%s' to integer - %w", recPort, err)
	}

	in.clientCfg.Scheme = aosScheme
	in.clientCfg.Host = aosHost
	in.clientCfg.Port = uint16(aosPortInt)
	in.clientCfg.User = aosUser
	in.clientCfg.Pass = aosPass
	in.clientCfg.TlsConfig = &tls.Config{
		InsecureSkipVerify: true, // todo: something less shameful
	}

	for i, streamType := range []string{
		apstra.StreamingConfigStreamingTypeAlerts,
		apstra.StreamingConfigStreamingTypeEvents,
		apstra.StreamingConfigStreamingTypePerfmon,
	} {
		stc := apstra.StreamTargetCfg{
			StreamingType:     streamType,
			SequencingMode:    apstra.StreamingConfigSequencingModeSequenced,
			Protocol:          apstra.StreamingConfigProtocolProtoBufOverTcp,
			AosTargetHostname: recHost,
			Port:              uint16(i + recPortInt),
		}
		in.streamTargetCfg = append(in.streamTargetCfg, stc)

		scp := apstra.StreamingConfigParams{
			StreamingType:  streamType,
			SequencingMode: apstra.StreamingConfigProtocolProtoBufOverTcp,
			Protocol:       apstra.StreamingConfigSequencingModeSequenced,
			Port:           uint16(i + recPortInt),
		}
		in.streamingConfigParams = append(in.streamingConfigParams, scp)
	}

	return nil
}

func main() {
	// handle interrupts
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt, os.Kill)

	// configuration objects
	//clientCfg := aosSdk.ClientCfg{}                           // config for interacting with AOS API
	//var streamingConfigs []aosSdk.StreamingConfigParams       // config for pointing event stream at our target
	//var streamTargetConfigs []apstraStreamTarget.StreamTargetCfg // config for our event stream target
	cfg := getConfigIn{
		clientCfg:             &apstra.ClientCfg{},
		streamingConfigParams: []apstra.StreamingConfigParams{},
		streamTargetCfg:       []apstra.StreamTargetCfg{},
	}

	// populate configuration objects using local function
	err := getConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// create AOS client
	// noinspection GoVetCopyLock
	c, err := apstra.NewClient(cfg.clientCfg)
	if err != nil {
		log.Fatal(err)
	}

	// noinspection GoUnhandledErrorResult
	defer c.Logout(context.TODO())

	// create aggregator channels where we'll get messages from all target services
	msgChan := make(chan *apstra.StreamingMessage)
	errChan := make(chan error)

	var streamTargets []*apstra.StreamTarget
	for i := range cfg.streamTargetCfg {
		// create each AOS stream target service
		st, err := apstra.NewStreamTarget(&cfg.streamTargetCfg[i])
		if err != nil {
			log.Fatal(err)
		}

		// start this AOS stream target service
		mc, ec, err := st.Start()
		if err != nil {
			log.Fatal(err)
		}

		// register this AOS stream target as a streaming config / receiver
		err = st.Register(context.TODO(), c)
		if err != nil {
			log.Fatal(err)
		}

		// copy messages from this target's message channel to aggregated message channel
		go func(in <-chan *apstra.StreamingMessage, out chan<- *apstra.StreamingMessage) {
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

		streamTargets = append(streamTargets, st)
	}

MainLoop:
	for {
		select {
		// interrupt (ctrl-c or whatever)
		case <-quitChan:
			break MainLoop
		// apstraStreamTarget has a message
		case msg := <-msgChan:
			var seqNumStr string
			switch msg.SequencingMode {
			case apstra.StreamingConfigSequencingModeSequenced:
				seqNumStr = strconv.Itoa(int(*msg.SequenceNum))
			case apstra.StreamingConfigSequencingModeUnsequenced:
				seqNumStr = "N/A"
			}
			log.Printf("%s / %s / message number %s / %s\n", msg.StreamingType, msg.SequencingMode, seqNumStr, msg.Message)
		// apstraStreamTarget has an error
		case err := <-errChan:
			log.Fatal(err)
		}
	}

	for _, st := range streamTargets {
		err = st.Unregister(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}
