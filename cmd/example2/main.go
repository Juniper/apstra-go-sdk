package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chrismarget-j/apstraTelemetry/apstra"
)

func keyLogWriter() (io.Writer, error) {
	keyLogDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	keyLogFile := filepath.Join(keyLogDir, ".aosSdk.log")

	err = os.MkdirAll(filepath.Dir(keyLogFile), os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	return os.OpenFile(keyLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func getConfig() (*apstra.ClientCfg, error) {
	aosScheme, foundAosScheme := os.LookupEnv(apstra.EnvApstraScheme)
	aosUser, foundAosUser := os.LookupEnv(apstra.EnvApstraUser)
	aosPass, foundAosPass := os.LookupEnv(apstra.EnvApstraPass)
	aosHost, foundAosHost := os.LookupEnv(apstra.EnvApstraHost)
	aosPort, foundAosPort := os.LookupEnv(apstra.EnvApstraPort)

	switch {
	case !foundAosScheme:
		return nil, fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraScheme)
	case !foundAosUser:
		return nil, fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraUser)
	case !foundAosPass:
		return nil, fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraPass)
	case !foundAosHost:
		return nil, fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraHost)
	case !foundAosPort:
		return nil, fmt.Errorf("environment variable '%s' not found", apstra.EnvApstraPort)
	}

	aosPortInt, err := strconv.Atoi(aosPort)
	if err != nil {
		return nil, fmt.Errorf("error converting '%s' to integer - %w", aosPort, err)
	}

	klw, err := keyLogWriter()
	if err != nil {
		return nil, err
	}

	var result apstra.ClientCfg
	result.Scheme = aosScheme
	result.Host = aosHost
	result.Port = uint16(aosPortInt)
	result.User = aosUser
	result.Pass = aosPass
	result.TlsConfig = &tls.Config{
		InsecureSkipVerify: true, // todo: something less shameful
		KeyLogWriter:       klw,
	}

	return &result, nil
}

func randString(n int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// handle interrupts
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt, os.Kill)

	rand.Seed(time.Now().UnixNano())

	clientCfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	// create AOS client
	c, err := apstra.NewClient(clientCfg)
	if err != nil {
		log.Fatal(err)
	}

	// login
	err = c.Login(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// noinspection GoUnhandledErrorResult
	defer c.Logout(context.TODO())

	blueprints, err := c.GetAllBlueprintIds(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(blueprints)
	bpId := blueprints[0]

	name := randString(10)
	rzid, err := c.CreateRoutingZone(context.TODO(), &apstra.CreateRoutingZoneCfg{
		SzType:      "evpn",
		VrfName:     name,
		Label:       "label-" + name,
		BlueprintId: bpId,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("VRF: " + name + " ID: " + string(rzid))

	log.Println("waiting")
	<-quitChan
	log.Println("got ctrl-c in main")

	vrfs, err := c.GetRoutingZones(context.TODO(), bpId)
	if err != nil {
		log.Fatal(err)
	}

	for _, vrf := range vrfs {
		if vrf.VrfName != "default" {
			log.Println("deleting ", vrf.Id, vrf.VrfName)
			err = c.DeleteRoutingZone(context.TODO(), bpId, vrf.Id)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
