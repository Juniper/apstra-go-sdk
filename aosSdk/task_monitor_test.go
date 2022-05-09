package aosSdk

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

func taskMonitorConfigTestClient1() (*Client, error) {
	user, foundUser := os.LookupEnv(EnvApstraUser)
	pass, foundPass := os.LookupEnv(EnvApstraPass)
	scheme, foundScheme := os.LookupEnv(EnvApstraScheme)
	host, foundHost := os.LookupEnv(EnvApstraHost)
	portstr, foundPort := os.LookupEnv(EnvApstraPort)

	switch {
	case !foundUser:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraUser)
	case !foundPass:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraPass)
	case !foundScheme:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraScheme)
	case !foundHost:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraHost)
	case !foundPort:
		return nil, fmt.Errorf("environment variable '%s' not found", EnvApstraPort)
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return nil, fmt.Errorf("error converting '%s' to integer - %w", portstr, err)
	}

	return NewClient(&ClientCfg{
		Scheme:    scheme,
		Host:      host,
		Port:      uint16(port),
		User:      user,
		Pass:      pass,
		TlsConfig: tls.Config{InsecureSkipVerify: true},
	})
}

//func TestGetTaskByBlueprintIdAndTaskId(t *testing.T) {
//	client, err := taskMonitorConfigTestClient1()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	result, err := client.GetTaskByBlueprintIdAndTaskId("db10754a-610e-475b-9baa-4c85f82282e8", "46dbddde-8d2d-410d-ac70-8bf6c110afe2")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	buf := bytes.Buffer{}
//	err = pp(result, &buf)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Print(buf.String())
//}

func TestBufIoReaderStuff(t *testing.T) {
	peekSize := 5
	unBufReader := strings.NewReader("01234567890123456789") //20
	bufReader := bufio.NewReader(unBufReader)                // 20
	peek, err := bufReader.Peek(peekSize)
	if err != nil && !errors.Is(err, io.EOF) { // error other than EOF?
		log.Fatal(err)
	}
	log.Println("peek", string(peek))
	buf := bytes.Buffer{}

	i, err := buf.ReadFrom(bufReader)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got %d bytes from buffered reader: %s\n", i, buf.String())

	j, err := buf.ReadFrom(bufReader)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got %d bytes from unbuffered reader: %s\n", j, buf.String())
}

func TestChanClose(t *testing.T) {
	testChan := make(chan struct{})
	select {
	case <-testChan:
		log.Println("read from testChan")
	default:
	}
	log.Println("closing testChan")
	close(testChan)
	log.Println("closing testChan again")
	close(testChan)
	log.Println("closed testChan twice")
}
