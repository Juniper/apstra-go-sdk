package goapstra

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"testing"
)

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
	parsed, err := url.Parse(test)
	if err != nil {
		log.Fatal(err)
	}
	resultBpId := blueprintIdFromUrl(parsed)
	if testBpId != resultBpId {
		log.Fatalf("expected '%s', got '%s'", testBpId, resultBpId)
	}
}
