package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/chrismarget-j/apstraTelemetry/aosSdk"
	"github.com/chrismarget-j/apstraTelemetry/aosStreamTarget"
)

// getConfigIn includes details for the clientCfg (AOS API
// location+credentials) the required streamingConfigParams (AOS API configuration to
// point events at our collector) and streamTargetCfg (our listener for AOS
// protobuf messages)
type getConfigIn struct {
	clientCfg             *aosSdk.ClientCfg                 // AOS API client config
	streamTargetCfg       []aosStreamTarget.StreamTargetCfg // Our protobuf stream listener
	streamingConfigParams []aosSdk.StreamingConfigParams    // Tell AOS API about our stream listener
}

func keyLogWriter() (io.Writer, error) {
	keyLogDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	keyLogFile := filepath.Join(keyLogDir, ".example1.log")

	err = os.MkdirAll(filepath.Dir(keyLogFile), os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	return os.OpenFile(keyLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func getConfig(in *getConfigIn) error {
	aosScheme, foundAosScheme := os.LookupEnv(aosSdk.EnvApstraScheme)
	aosUser, foundAosUser := os.LookupEnv(aosSdk.EnvApstraUser)
	aosPass, foundAosPass := os.LookupEnv(aosSdk.EnvApstraPass)
	aosHost, foundAosHost := os.LookupEnv(aosSdk.EnvApstraHost)
	aosPort, foundAosPort := os.LookupEnv(aosSdk.EnvApstraPort)
	recHost, foundRecHost := os.LookupEnv(aosStreamTarget.EnvApstraStreamHost)
	recPort, foundRecPort := os.LookupEnv(aosStreamTarget.EnvApstraStreamBasePort)

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
		return fmt.Errorf("environment variable '%s' not found", aosStreamTarget.EnvApstraStreamHost)
	case !foundRecPort:
		return fmt.Errorf("environment variable '%s' not found", aosStreamTarget.EnvApstraStreamBasePort)
	}

	aosPortInt, err := strconv.Atoi(aosPort)
	if err != nil {
		return fmt.Errorf("error converting '%s' to integer - %w", aosPort, err)
	}

	recPortInt, err := strconv.Atoi(recPort)
	if err != nil {
		return fmt.Errorf("error converting '%s' to integer - %w", recPort, err)
	}

	klw, err := keyLogWriter()
	if err != nil {
		return err
	}

	in.clientCfg.Scheme = aosScheme
	in.clientCfg.Host = aosHost
	in.clientCfg.Port = uint16(aosPortInt)
	in.clientCfg.User = aosUser
	in.clientCfg.Pass = aosPass
	in.clientCfg.TlsConfig = tls.Config{
		InsecureSkipVerify: true, // todo: something less shameful
		KeyLogWriter:       klw,
	}

	for i, streamType := range []string{
		aosSdk.StreamingConfigStreamingTypeAlerts,
		aosSdk.StreamingConfigStreamingTypeEvents,
		aosSdk.StreamingConfigStreamingTypePerfmon,
	} {
		stc := aosStreamTarget.StreamTargetCfg{
			StreamingType:     streamType,
			SequencingMode:    aosSdk.StreamingConfigSequencingModeSequenced,
			Protocol:          aosSdk.StreamingConfigProtocolProtoBufOverTcp,
			AosTargetHostname: recHost,
			Port:              uint16(i + recPortInt),
		}
		in.streamTargetCfg = append(in.streamTargetCfg, stc)

		scp := aosSdk.StreamingConfigParams{
			StreamingType:  streamType,
			SequencingMode: aosSdk.StreamingConfigProtocolProtoBufOverTcp,
			Protocol:       aosSdk.StreamingConfigSequencingModeSequenced,
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
	//var streamTargetConfigs []aosStreamTarget.StreamTargetCfg // config for our event stream target
	cfg := getConfigIn{
		clientCfg:             &aosSdk.ClientCfg{},
		streamingConfigParams: []aosSdk.StreamingConfigParams{},
		streamTargetCfg:       []aosStreamTarget.StreamTargetCfg{},
	}

	// populate configuration objects using local function
	err := getConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// create AOS client
	// noinspection GoVetCopyLock
	c := aosSdk.NewClient(cfg.clientCfg)

	// noinspection GoUnhandledErrorResult
	defer c.Logout()

	// create aggregator channels where we'll get messages from all target services
	msgChan := make(chan *aosStreamTarget.StreamingMessage)
	errChan := make(chan error)

	var streamTargets []*aosStreamTarget.StreamTarget
	for i := range cfg.streamTargetCfg {
		// create each AOS stream target service
		st, err := aosStreamTarget.NewStreamTarget(&cfg.streamTargetCfg[i])
		if err != nil {
			log.Fatal(err)
		}

		// start this AOS stream target service
		mc, ec, err := st.Start()
		if err != nil {
			log.Fatal(err)
		}

		// register this AOS stream target as a streaming config / receiver
		err = st.Register(c)
		if err != nil {
			log.Fatal(err)
		}

		// copy messages from this target's message channel to aggregated message channel
		go func(in <-chan *aosStreamTarget.StreamingMessage, out chan<- *aosStreamTarget.StreamingMessage) {
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

	err = c.Login()
	if err != nil {
		log.Fatal(err)
	}

MainLoop:
	for {
		select {
		// interrupt (ctrl-c or whatever)
		case <-quitChan:
			break MainLoop
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

	for _, st := range streamTargets {
		err = st.Unregister()
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}
