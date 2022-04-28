package main

import (
	"flag"
	apstratelemetry "github.com/chrismarget-j/apstraTelemetry"
	"log"
	"net/url"
	"os"
	"strconv"
)

const (
	ENV_APSTRA_PW = "APSTRA_PW"
)

func main() {
	pw, found := os.LookupEnv(ENV_APSTRA_PW)
	if !found {
		log.Fatalf("error environment '%s' is required", ENV_APSTRA_PW)
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
	aosClient, err := apstratelemetry.NewAosClient(cfg)
	if err != nil {
		log.Fatalf("error creating new AOS client - %v", err)
	}

	err = aosClient.Login()
	if err != nil {
		log.Fatalf("error logging in AOS client - %v", err)
	}

	err = aosClient.Logout()
	if err != nil {
		log.Fatalf("error logging out in AOS client - %v", err)
	}
}
