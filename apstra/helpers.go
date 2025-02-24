// Copyright (c) Juniper Networks, Inc., 2022-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
	"github.com/google/uuid"
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

var (
	uuidInit      bool
	uuidInitMutex sync.Mutex
)

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

func toPtr[A any](a A) *A {
	return &a
}

func stringerPtrToStringPtr(in fmt.Stringer) *string {
	if in == nil {
		return nil
	}
	// this interesting thing checks to make sure this thing is really nil...
	if reflect.ValueOf(in).Kind() == reflect.Ptr && reflect.ValueOf(in).IsNil() {
		return nil
	}
	return toPtr(in.String())
}

func featureSwitchEnumFromStringPtr(in *string) *enum.FeatureSwitch {
	if in == nil {
		return nil
	}
	return enum.FeatureSwitches.Parse(*in)
}

func isv4(ip net.IP) bool {
	return 4 == len(ip.To4())
}

func isv6(ip net.IP) bool {
	if ip.To4() != nil {
		return false
	}
	return 16 == len(ip.To16())
}

// ipFromString is an improvement on calling net.ParseIP() directly because it
// handles empty strings gracefully (returns nil net.IP) and because it returns
// errors in case of un-parseable input strings.
func ipFromString(s string) (net.IP, error) {
	if s == "" {
		return nil, nil
	}

	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("cannot parse IP %q", s)
	}

	return ip, nil
}

// ipNetFromString is an improvement on calling net.ParseCIDR() directly because
// it handles empty strings gracefully (returns nil pointer) and because it
// returns a net.IPNet with the actual IP address, rather than the base address.
func ipNetFromString(s string) (*net.IPNet, error) {
	if s == "" {
		return nil, nil
	}

	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, fmt.Errorf("while parsing CIDR string %q - %w", s, err)
	}
	ipNet.IP = ip

	return ipNet, nil
}

// macFromString is an improvement on calling net.ParseMAC() directly because it
// handles empty strings gracefully (returns nil net.HardwareAddr)
func macFromString(s string) (net.HardwareAddr, error) {
	if s == "" {
		return nil, nil
	}

	mac, err := net.ParseMAC(s)
	if err != nil {
		return nil, fmt.Errorf("cannot parse hardware address %q", s)
	}

	return mac, nil
}
