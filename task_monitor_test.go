package goapstra

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"testing"
)

func taskMonitorConfigTestClient1() (*Client, error) {
	return NewClient(&ClientCfg{
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
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

func TestBlueprintIdFromUrl(t *testing.T) {
	testBpId := ObjectId("lkasdlfaj")
	test := "https://host:443" + fmt.Sprintf(apiUrlBlueprintById, testBpId)
	url, err := url.Parse(string(test))
	if err != nil {
		log.Fatal(err)
	}
	resultBpId, err := blueprintIdFromUrl(url)
	if err != nil {
		log.Fatal(err)
	}
	if testBpId != resultBpId {
		log.Fatalf("expected '%s', got '%s'", testBpId, resultBpId)
	}
}
