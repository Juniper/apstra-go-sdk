package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/Juniper/apstra-go-sdk/apstra"
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
	clientCfg := apstra.ClientCfg{
		Url:        "https://apstra.example.com",
		User:       "admin",
		Pass:       "password",
		HttpClient: &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
	}
	client, err := clientCfg.NewClient(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// figure out our IP address -- we'll tell Apstra to fire protobuf structures at this address
	ourIp, err := ourIpForPeer(net.ParseIP(os.Getenv("APSTRA_HOST")))
	if err != nil {
		log.Fatal(err)
	}

	// create local stream target object
	streamTargetConfig := apstra.StreamTargetCfg{
		Certificate:       nil, // apstra doesn't support TLS (?!?)
		Key:               nil, // apstra doesn't support TLS (?!?)
		SequencingMode:    apstra.StreamingConfigSequencingModeSequenced,
		StreamingType:     apstra.StreamingConfigStreamingTypeAlerts,
		Protocol:          apstra.StreamingConfigProtocolProtoBufOverTcp,
		Port:              9999,
		AosTargetHostname: ourIp.String(),
	}
	streamTarget, err := apstra.NewStreamTarget(&streamTargetConfig)
	if err != nil {
		log.Fatal(err)
	}

	// start the stream target service listening for incoming messages
	streamMsgChan, streamErrChan, err := streamTarget.Start()
	if err != nil {
		log.Fatal(err)
	}

	quitChan := make(chan os.Signal, 1)
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
			log.Println(msg.Message.String())
		case err := <-streamErrChan:
			log.Println(err.Error())
		}
	}
}
