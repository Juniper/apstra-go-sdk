package main

import (
	apstratelemetry "github.com/chrismarget-j/apstraTelemetry"
	"log"
)

func main() {
	cfg := apstratelemetry.AosClientCfg{
		Host:   "66.129.234.206",
		Port:   uint16(37000),
		Scheme: "hxxps",
		User:   "admin",
		Pass:   "admin",
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
