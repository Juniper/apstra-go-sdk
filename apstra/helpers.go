package apstra

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

// keyLogWriter takes an environment variable which might name a logfile for
// exporting TLS session keys. If so, it returns an io.Writer to be used for
// that purpose.
func keyLogWriter(keyLogEnv string) (io.Writer, error) {
	keyLogFile, foundKeyLogFile := os.LookupEnv(keyLogEnv)
	if !foundKeyLogFile {
		return nil, nil
	}

	err := os.MkdirAll(filepath.Dir(keyLogFile), os.FileMode(0600))
	if err != nil {
		return nil, err
	}
	return os.OpenFile(keyLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func intSliceContains(in []int, t int) bool {
	for _, i := range in {
		if i == t {
			return true
		}
	}
	return false
}

func intSliceToStringSlice(in []int) []string {
	var result []string
	for _, i := range in {
		result = append(result, strconv.Itoa(i))
	}
	return result
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
