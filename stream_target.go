package goapstra

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	EnvApstraStreamHost     = "APSTRA_STREAM_HOST"
	EnvApstraStreamBasePort = "APSTRA_STREAM_BASE_PORT"

	sizeOfAosMessageLenHdr = 2 // Apstra 'protoBufOverTcp' streaming includes a 16-bit length w/each protobuf
	network                = "tcp4"
	errConnClosed          = "use of closed network connection"
)

// StreamTargetCfg is used when initializing an instance of
// StreamTarget with NewStreamTarget. If Cert or Key are nil, the
// StreamTarget will use bare TCP rather than TLS.
type StreamTargetCfg struct {
	Certificate       *x509.Certificate
	Key               *rsa.PrivateKey
	SequencingMode    string
	StreamingType     string
	Protocol          string
	Port              uint16
	AosTargetHostname string
}

// StreamingMessage is a wrapper structure for messages delivered by both
// StreamingConfigSequencingModeSequenced channels and by
// StreamingConfigSequencingModeUnsequenced channels. In the latter case, 'seq'
// will always be nil.
type StreamingMessage struct {
	SequencingMode string
	StreamingType  string
	Message        *AosMessage
	SequenceNum    *uint64
}

// NewStreamTarget creates a StreamTarget (socket listener) either with TLS
// support (when both x509Cert and privkey are supplied) or using bare TCP
// (when either x509Cert or privkey are nil)
func NewStreamTarget(cfg *StreamTargetCfg) (*StreamTarget, error) {
	var tlsConfig *tls.Config

	keyLog, err := keyLogWriterFromEnv(EnvApstraStreamKeyLogFile)
	if err != nil {
		return nil, err
	}

	if cfg.Certificate != nil && cfg.Key != nil {
		certBlock := bytes.NewBuffer(nil)
		err = pem.Encode(certBlock, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cfg.Certificate.Raw,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to pem encode certificate block - %w", err)
		}

		privateKeyBlock := bytes.NewBuffer(nil)
		err = pem.Encode(privateKeyBlock, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(cfg.Key),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to pem encode private key block - %w", err)
		}

		tlsCert, err := tls.X509KeyPair(certBlock.Bytes(), privateKeyBlock.Bytes())
		if err != nil {
			return nil, fmt.Errorf("error parsing tls.Certificate object from cert and key - %w", err)
		}

		tlsConfig = &tls.Config{
			KeyLogWriter: keyLog,
			Rand:         rand.Reader,
			Certificates: []tls.Certificate{tlsCert},
		}
	}

	return &StreamTarget{
		cfg:       cfg,
		errChan:   make(chan error),
		stopChan:  make(chan struct{}),
		msgChan:   make(chan *StreamingMessage),
		tlsConfig: tlsConfig,
	}, nil
}

// StreamTarget is a listener for AOS streaming objects
type StreamTarget struct {
	tlsConfig *tls.Config            // if we're a TLS listener
	nl        net.Listener           // clients (apstra server) connect here
	stopChan  chan struct{}          // close to rip everythign down
	errChan   chan error             // client handlers pass errors here
	msgChan   chan *StreamingMessage // client handlers pass messages here
	clientWG  sync.WaitGroup         // keeps track of client handlers
	cfg       *StreamTargetCfg       // submitted by caller
	apstraIP  *net.IP                // for filtering incoming connections
	strmCfgId ObjectId               // AOS streaming ID, populated by Register
	client    *Client                // populated by Register, we hang onto it for Unregister
}

// Start loops forever handling new connections from the AOS streaming service
// as they arrive. Messages generated by socket clients are sent to msgChan.
// Receive errors are sent to errChan. An error is returned immediately if
// there's a problem starting the client handling loop.
func (o *StreamTarget) Start() (msgChan <-chan *StreamingMessage, errChan <-chan error, err error) {
	var nl net.Listener

	laddr := ":" + strconv.Itoa(int(o.cfg.Port)) // something like ":6000" (a port number)
	if o.tlsConfig != nil {
		nl, err = tls.Listen(network, laddr, o.tlsConfig) // if we're doing TLS (tls.listener)
	} else {
		nl, err = net.Listen(network, laddr) // if we're doing raw TCP (net.TCPListener)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error starting listener - %w", err)
	}

	// loop accepting incoming connections
	// this will stop when nl.Close() gets called
	go o.receive(nl)

	// anonymous shutdown go func kicks in when stopChan is closed
	// todo: find out if this works the way i think it does
	go func() {
		<-o.stopChan      // wait for Stop() to close stopChan
		err := nl.Close() // close socket listener
		if err != nil {
			o.errChan <- err
		}
		o.clientWG.Wait() // wait for client conn handlers to exit
		close(o.errChan)  // close errChan to signal to Stop() that we're done
	}()

	return o.msgChan, o.errChan, nil
}

// Stop shuts down the receiver
func (o *StreamTarget) Stop() {
	close(o.stopChan) // signal exit to client conn handlers
	o.clientWG.Wait() // wait for client conn handlers to exit
	for range o.errChan {
	} // We're done when err channel gets closed
}

// Id returns the StreamTarget ID returned by Apstra during registration or ""
// if unregistered.
func (o *StreamTarget) Id() ObjectId {
	return o.strmCfgId
}

// receive loops until the listener gets closed, handing off connections from the
// AOS server to instances of handleClientConn().
func (o *StreamTarget) receive(nl net.Listener) {
	// loop accepting new connections
	for {
		conn, err := nl.Accept() // block here waiting for inbound client
		if err != nil {
			// nl got closed (graceful shutdown) or we've encountered an error
			if strings.HasSuffix(err.Error(), errConnClosed) {
				o.errChan <- err // this is a graceful close, but send the error along anyway
				return           // that's all folks
			} else {
				o.errChan <- err // real error. fire into the channel
			}
			continue // go collect the next client
		}

		o.clientWG.Add(1)                                 // defered close in handleClientConn
		go o.handleClientConn(conn, o.msgChan, o.errChan) // read messages from the socket

		// close this Conn when triggered by stopChan
		go func() {
			<-o.stopChan
			// noinspection GoUnhandledErrorResult
			conn.Close()
		}()
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
	msgLenHdr, err := getBytesFromConn(sizeOfAosMessageLenHdr, conn)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(msgLenHdr), nil
}

func (o *StreamTarget) handleClientConn(conn net.Conn, msgChan chan<- *StreamingMessage, errChan chan<- error) {
	// noinspection GoUnhandledErrorResult
	defer conn.Close()
	defer o.clientWG.Done()

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

		msg, err := o.msgFromBytes(payload)
		if err != nil {
			errChan <- err
		} else {
			msgChan <- msg
		}
	}
}

func (o *StreamTarget) msgFromBytes(in []byte) (*StreamingMessage, error) {
	var msgOut AosMessage
	var seqPtr *uint64

	// extract AosMessage from AosSequencedMessage wrapper if configured for StreamingConfigSequencingModeSequenced
	if o.cfg.SequencingMode == StreamingConfigSequencingModeSequenced {
		var seqMsg AosSequencedMessage      // outer wrapper structure
		err := proto.Unmarshal(in, &seqMsg) // unwrap inner message
		if err != nil {
			return nil, err
		}
		in = seqMsg.AosProto         // redefine 'in' to payload of sequencing wrapper
		seqNum := seqMsg.GetSeqNum() // extract sequence number from sequencing wrapper
		seqPtr = &seqNum             // record address of sequence number in pointer we'll return
	}

	err := proto.Unmarshal(in, &msgOut) // extract inner message
	return &StreamingMessage{
		StreamingType:  o.cfg.StreamingType,
		SequencingMode: o.cfg.SequencingMode,
		Message:        &msgOut, // pointer to inner message
		SequenceNum:    seqPtr,  // pointer to sequence number (nil if unsequenced)
	}, err
}

// Register registers this StreamTarget as a streaming config / receiver on the
// AOS server. If o.cfg.AosTargetHostname is non-empty, that value will be told
// to AOS when configuring the streaming config / receiver. If it's empty, we
// attempt to determine the local IP nearest to the AOS server, use that value
// (as a string)
func (o *StreamTarget) Register(ctx context.Context, client *Client) error {
	// figure out what the AOS server should call us (string: IP or DNS name)
	var apstraTargetHostname string
	switch o.cfg.AosTargetHostname {
	case "": // no value supplied - find a local IP
		ourIp, err := ourIpForPeer(net.ParseIP(client.ServerName()))
		if err != nil {
			return fmt.Errorf("error determinging local IP for AOS '%s' streaming config - %w", client.ServerName(), err)
		}
		apstraTargetHostname = ourIp.String()
	default: // use whatever is in our configuration
		apstraTargetHostname = o.cfg.AosTargetHostname
	}

	// Register this target with Apstra
	id, err := client.NewStreamingConfig(ctx, &StreamingConfigParams{
		StreamingType:  o.cfg.StreamingType,
		SequencingMode: o.cfg.SequencingMode,
		Protocol:       o.cfg.Protocol,
		Hostname:       apstraTargetHostname,
		Port:           o.cfg.Port,
	})
	if err != nil {
		return fmt.Errorf("error in Register() - %w", err)
	}
	o.strmCfgId = id  // save the streamingConfig ID returned by Apstra
	o.client = client // hang onto the client pointer for use in Unregister()
	return nil
}

// Unregister deletes the streaming config / receiver associated with this
// StreamTarget from the AOS server.
func (o *StreamTarget) Unregister(ctx context.Context) error {
	if o.strmCfgId == "" {
		return errors.New("no stream id for this StreamTarget, cannot UnRegister")
	}

	err := o.client.DeleteStreamingConfig(ctx, o.strmCfgId)
	if err != nil {
		return err
	}

	o.strmCfgId = ""

	return nil
}
