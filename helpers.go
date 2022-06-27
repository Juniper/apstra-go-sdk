package goapstra

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// keyLogWriter takes an environment variable which might name a logfile for
// exporting TLS session keys. If so, it returns an io.Writer to be used for
// that purpose.
func keyLogWriter(keyLogEnv string) (io.WriteCloser, error) {
	fileName, foundKeyLogFile := os.LookupEnv(keyLogEnv)
	if !foundKeyLogFile {
		return nil, nil
	}

	// expand ~ style home directory
	if strings.HasPrefix(fileName, "~/") {
		dirname, _ := os.UserHomeDir()
		fileName = filepath.Join(dirname, fileName[2:])
	}

	err := os.MkdirAll(filepath.Dir(fileName), os.FileMode(0600))
	if err != nil {
		return nil, err
	}
	return os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func pp(in interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(in); err != nil {
		return err
	}
	return nil
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

// AsnOverlap returns a bool indicating whether two AsnRange objects overlap
func AsnOverlap(a, b AsnRange) bool {
	if a.First >= b.First && a.First <= b.Last { // begin 'a' falls within 'b'
		return true
	}
	if a.Last <= b.Last && a.Last >= b.First { // end 'a' falls within 'b'
		return true
	}
	if b.First >= a.First && b.First <= a.Last { // begin 'b' falls within 'a'
		return true
	}
	if b.Last <= a.Last && b.Last >= a.First { // end 'b' falls within 'a'
		return true
	}
	return false // no overlap
}
