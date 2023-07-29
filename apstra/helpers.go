package apstra

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"sync"
	"time"
)

func pp(in interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	err := enc.Encode(in)
	return err
}

// ourIpForPeer returns a *net.IP representing the local interface selected by
// the system for talking to the passed *net.IP. The returned value might also
// be the best choice for that peer to reach us.
func ourIpForPeer(them net.IP) (*net.IP, error) {
	c, err := net.Dial("udp4", them.String()+":1")
	if err != nil {
		return nil, err
	}

	return &c.LocalAddr().(*net.UDPAddr).IP, c.Close()
}

func immediateTicker(interval time.Duration) *time.Ticker {
	nc := make(chan time.Time, 1)
	ticker := time.NewTicker(interval)
	oc := ticker.C
	go func() {
		nc <- time.Now()
		for tm := range oc {
			nc <- tm
		}
	}()
	ticker.C = nc
	return ticker
}

func itemInSlice[A comparable](item A, slice []A) bool {
	for i := range slice {
		if item == slice[i] {
			return true
		}
	}
	return false
}

var uuidInit bool
var uuidInitMutex sync.Mutex

// initUUID sets the "hardware address" used for generating UUIDv1 strings to "apstra"
func initUUID() {
	uuidInitMutex.Lock()
	defer uuidInitMutex.Unlock()
	if !uuidInit {
		uuid.SetNodeID([]byte("apstra"))
		uuidInit = true
	}
}

// uuid1AsObjectId returns a new UUIDv1 string as an ObjectId
func uuid1AsObjectId() (ObjectId, error) {
	initUUID()
	uuid1, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("failed while invoking uuid>NewUUID() - %w", err)
	}
	return ObjectId(uuid1.String()), nil
}
