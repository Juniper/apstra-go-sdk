package aosSdk

import (
	"encoding/json"
	"io"
	"net"
	"strconv"
)

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
func ourIpForPeer(us *net.IP) (*net.IP, error) {
	c, err := net.Dial("udp4", us.String()+":1")
	if err != nil {
		return nil, err
	}

	return &c.LocalAddr().(*net.UDPAddr).IP, c.Close()
}
