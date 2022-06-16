package main

import (
	"context"
	"crypto/tls"
	"github.com/chrismarget-j/goapstra"
	"log"
	"net"
	"os"
	"os/signal"
)

// ourIpForPeer returns a *net.IP representing the local interface selected by
// the system for talking to the passed *net.IP. The returned value might also
// be the best choice for that peer to reach us.
func ourIpForPeer(them net.IP) (*net.IP, error) {
	c, err := net.Dial("udp4", them.String()+":1")
	if err != nil {
		return nil, err
	}

	return &c.LocalAddr().(*net.UDPAddr).IP, c.Close()
}

func main() {
	// create an apstra client object
	clientCfg := goapstra.ClientCfg{
		//Host:      "", // omit to use env var 'APSTRA_HOST'
		//Port:      "", // omit to use env var 'APSTRA_PORT'
		//User:      "", // omit to use env var 'APSTRA_USER'
		//Pass:      "", // omit to use env var 'APSTRA_PASS'
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client, err := goapstra.NewClient(&clientCfg)
	if err != nil {
		log.Fatal(err)
	}

	client.GetAllBlueprintIds(context.TODO())

	// figure out our IP address -- we'll tell Apstra to fire protobuf structures at this address
	ourIp, err := ourIpForPeer(net.ParseIP(os.Getenv("APSTRA_HOST")))
	if err != nil {
		log.Fatal(err)
	}

	// create local stream target object
	streamTargetConfig := goapstra.StreamTargetCfg{
		Certificate:       nil, // apstra doesn't support TLS (?!?)
		Key:               nil, // apstra doesn't support TLS (?!?)
		SequencingMode:    goapstra.StreamingConfigSequencingModeSequenced,
		StreamingType:     goapstra.StreamingConfigStreamingTypeAlerts,
		Protocol:          goapstra.StreamingConfigProtocolProtoBufOverTcp,
		Port:              9999,
		AosTargetHostname: ourIp.String(),
	}
	streamTarget, err := goapstra.NewStreamTarget(&streamTargetConfig)
	if err != nil {
		log.Fatal(err)
	}

	// start the stream target service listening for incoming messages
	streamMsgChan, streamErrChan, err := streamTarget.Start()
	if err != nil {
		log.Fatal(err)
	}

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt)

	// tell apstra to send protobuf messages to us
	err = streamTarget.Register(context.TODO(), client)
	if err != nil {
		log.Fatal(err)
	}

	// check the API knows about our receiver
	ids, err := client.GetAllStreamingConfigIds(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Our Streaming Receiver has ID '%s'. The complete set of receivers is %s", streamTarget.Id(), ids)

	// loop until ctrl-c, print messages+errors as they arrive
	for {
		select {
		case <-quitChan:
			err = streamTarget.Unregister(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			return
		case msg := <-streamMsgChan:
			log.Printf(msg.Message.String())
		case err := <-streamErrChan:
			log.Printf(err.Error())
		}
	}
}
