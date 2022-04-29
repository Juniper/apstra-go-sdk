package main

import (
	"bytes"
	"encoding/json"
	"flag"
	apstratelemetry "github.com/chrismarget-j/apstraTelemetry"
	"log"
	"net/url"
	"os"
	"strconv"
)

func main() {
	pw, found := os.LookupEnv(apstratelemetry.EnvApstraPass)
	if !found {
		log.Fatalf("error environment '%s' is required", apstratelemetry.EnvApstraPass)
	}

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("You need to specify Apstra URL")
	}

	aosUrl, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatalf("error parsing url from command line - %v", err)
	}
	port, err := strconv.Atoi(aosUrl.Port())
	if err != nil {
		log.Fatalf("error parsing port from URL - %v", err)
	}

	cfg := apstratelemetry.AosClientCfg{
		Host:   aosUrl.Hostname(),
		Port:   uint16(port),
		Scheme: aosUrl.Scheme,
		User:   aosUrl.User.Username(),
		Pass:   pw,
	}
	aosClient, err := apstratelemetry.NewAosClient(&cfg)
	if err != nil {
		log.Fatalf("error creating new AOS client - %v", err)
	}

	err = aosClient.UserLogin()
	if err != nil {
		log.Fatalf("error logging in AOS client - %v", err)
	}

	ver, err := aosClient.GetVersion()
	if err != nil {
		log.Fatalf("error getting AOS version - %v", err)
	}

	sc, err := aosClient.GetAllStreamingConfigs()
	if err != nil {
		log.Fatalf("error getting all streaming configs - %v", err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "    ")

	err = enc.Encode(ver)
	if err != nil {
		log.Fatalf("error encoding data to JSON - %v", err)
	}
	log.Println(buf.String())

	err = enc.Encode(sc)
	if err != nil {
		log.Fatalf("error encoding data to JSON - %v", err)
	}
	log.Println(buf.String())

	err = aosClient.UserLogout()
	if err != nil {
		log.Fatalf("error logging out in AOS client - %v", err)
	}
}
