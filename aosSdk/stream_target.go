package aosSdk

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/chrismarget-j/apstraTelemetry/aosStreaming"
)

const (
	sizeOfAosMsgLenHdr = 2
	network            = "tcp4"
	errConnClosed      = "use of closed network connection"
)

// AosStreamTargetCfg is used when initializing an instance of
// AosStreamTarget with NewStreamTarget. If Cert or Key are nil, the
// AosStreamTarget will use bare TCP rather than TLS.
type AosStreamTargetCfg struct {
	Certificate    *x509.Certificate
	Key            *rsa.PrivateKey
	SequencingMode StreamingConfigSequencingMode
	StreamingType  AosApiStreamingConfigStreamingType
	Protocol       AosApiStreamingConfigProtocol
	Port           uint16
}

// NewStreamTarget creates a AosStreamTarget (socket listener) either with TLS
// support (when both x509Cert and privkey are supplied) or using bare TCP
// (when either x509Cert or privkey are nil)
func NewStreamTarget(cfg *AosStreamTargetCfg) (*AosStreamTarget, error) {
	var tlsConfig *tls.Config

	if cfg.Certificate != nil && cfg.Key != nil {
		keyLog, err := keyLogWriter()
		if err != nil {
			return nil, err
		}

		certBlock := bytes.NewBuffer(nil)
		err = pem.Encode(certBlock, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cfg.Certificate.Raw,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to pem encode certificate block - %v", err)
		}

		privateKeyBlock := bytes.NewBuffer(nil)
		err = pem.Encode(privateKeyBlock, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(cfg.Key),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to pem encode private key block - %v", err)
		}

		tlsCert, err := tls.X509KeyPair(certBlock.Bytes(), privateKeyBlock.Bytes())
		if err != nil {
			// todo: wrap error
			return nil, err
		}

		tlsConfig = &tls.Config{
			KeyLogWriter: keyLog,
			Rand:         rand.Reader,
			Certificates: []tls.Certificate{tlsCert},
		}
	}

	return &AosStreamTarget{
		errChan:   make(chan error),
		stopChan:  make(chan struct{}),
		msgChan:   make(chan *aosStreaming.AosMessage),
		tlsConfig: tlsConfig,
	}, nil
}

// AosStreamTarget is a listener for AOS streaming objects
// todo: this thing should have a waitgroup to keep track of each accept() func
type AosStreamTarget struct {
	tlsConfig *tls.Config
	nl        net.Listener
	stopChan  chan struct{}
	errChan   chan error
	msgChan   chan *aosStreaming.AosMessage
	port      uint16
	//sessions    map[int]*Session
	//sessChMap   map[chan *Session]struct{}
	//sessChMutex *sync.Mutex
}

// Start loops forever handling new connections from the AOS streaming service
// as they arrive. Messages generated by socket clients are sent to msgChan.
// Receive errors are sent to errChan. An error is returned immediately if
// there's a problem starting the client handling loop.
func (o AosStreamTarget) Start() (msgChan <-chan *aosStreaming.AosMessage, errChan <-chan error, err error) {
	var nl net.Listener

	laddr := ":" + strconv.Itoa(int(o.port))
	if o.tlsConfig != nil {
		nl, err = tls.Listen(network, laddr, o.tlsConfig)
	} else {
		nl, err = net.Listen(network, laddr)
	}
	if err != nil {
		return nil, nil, err
	}

	// loop accepting incoming connections
	go o.receive(nl)

	// this should stop everything.
	// todo: find out if this works the way i think it does
	go func() {
		<-o.stopChan
		nl.Close()
		close(o.stopChan)
		close(o.errChan)
	}()

	return o.msgChan, o.errChan, nil
}

// Stop shuts down the receiver
// todo: block on waitgroup while each accept() func shuts down
func (o AosStreamTarget) Stop() {
	close(o.stopChan)
}

// receive loops until the listener shuts down, handing off connections from the
// AOS server to instances of handleMessages().
func (o *AosStreamTarget) receive(nl net.Listener) {
	// loop accepting new connections
	for {
		conn, err := nl.Accept()
		if err != nil {
			if strings.HasSuffix(err.Error(), errConnClosed) {
				o.errChan <- err
				return
			}
			o.errChan <- err
			continue
		}

		go handleMessages(conn, o.msgChan, o.errChan)
	}
}

func getBytesFromConn(i int, conn net.Conn) ([]byte, error) {
	data := make([]byte, i)
	n, err := io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	if n != i {
		return nil, fmt.Errorf("expected %d bytes, got %d", i, n)
	}
	return data, nil
}

func msgLenFromConn(conn net.Conn) (uint16, error) {
	msgLenHdr, err := getBytesFromConn(sizeOfAosMsgLenHdr, conn)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(msgLenHdr), nil
}

func handleClientConn(conn net.Conn, msgChan chan<- *aosStreaming.AosMessage, errChan chan<- error) {
	// todo: defer waitgroup.done() here
	defer conn.Close()

	for {
		msgLen, err := msgLenFromConn(conn)
		if err != nil {
			errChan <- err
			if err == io.EOF {
				return
			}
		}

		payload, err := getBytesFromConn(int(msgLen), conn)
		if err != nil {
			errChan <- err
			if err == io.EOF {
				return
			}
		}

		msg, err := msgFromBytes(payload)
		if err != nil {
			errChan <- err
		} else {
			msgChan <- msg
		}
	}
}

func msgFromBytes(in []byte) (*aosStreaming.AosMessage, error) {
	msg := &aosStreaming.AosMessage{}
	err := proto.Unmarshal(in, msg)
	return msg, err
}
