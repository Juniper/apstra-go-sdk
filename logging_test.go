package goapstra

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	topologyIds, err := topologyIdsFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	if len(topologyIds) == 0 {
		t.Skip("can't test client logging without any topologies")
	}

	topologies := make([]cloudlabsTopology, len(topologyIds))
	for i, id := range topologyIds {
		topology, err := getCloudlabsTopology(id)
		if err != nil {
			t.Fatal(err)
		}
		topologies[i] = *topology
	}

	for i, topology := range topologies {
		access, err := topology.getVmAccessInfo("aos-vm1", "https")
		if err != nil {
			t.Fatal(err)
		}

		loggers := make([]*log.Logger, LogVerbosityLevels)
		logFiles := make([]*os.File, LogVerbosityLevels)
		for j := 0; j < LogVerbosityLevels; j++ {
			logFiles[j], err = ioutil.TempFile("/tmp", fmt.Sprintf("client-%d-level-%d-log-", i, j))
			if err != nil {
				t.Fatal(err)
			}
			loggers[j] = log.New(logFiles[j], "", 0)
		}

		clientCfg := &ClientCfg{
			Scheme:    access.Protocol,
			Host:      access.Host,
			Port:      uint16(access.Port),
			User:      access.Username,
			Pass:      access.Password,
			Loggers:   loggers,
			TlsConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client, err := NewClient(clientCfg)
		if err != nil {
			t.Fatal(err)
		}

		msgs := make([]string, LogVerbosityLevels)
		for j := 0; j < LogVerbosityLevels; j++ {
			msgs[j] = randString(10, "hex")
			log.Printf("logging '%s' with level %d.", msgs[j], j)
			client.logStr(j, msgs[j])
		}

	severity:
		for j := 0; j < LogVerbosityLevels; j++ {
			err = logFiles[j].Close()
			if err != nil {
				log.Fatal(err)
			}
			file, err := os.Open(logFiles[j].Name())
			if err != nil {
				t.Fatal(err)
			}
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				if strings.TrimSuffix(scanner.Text(), "\n") == msgs[j] {
					log.Printf("message '%s' found in file '%s'", msgs[j], file.Name())
					err = file.Close()
					if err != nil {
						t.Fatal(err)
					}
					err = os.Remove(file.Name())
					if err != nil {
						t.Fatal(err)
					}
					continue severity
				}
			}
			log.Fatalf("logged string '%s' not found in file '%s'", msgs[j], file.Name())
		}
	}
}
