package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

const (
	apiUrlSubinterfaces = apiUrlBlueprintByIdPrefix + "subinterfaces"
)

var _ json.Marshaler = (*TwoStageL3ClosSubinterface)(nil)
var _ json.Unmarshaler = (*TwoStageL3ClosSubinterface)(nil)

type TwoStageL3ClosSubinterface struct {
	Ipv4AddrType *InterfaceNumberingIpv4Type
	Ipv6AddrType *InterfaceNumberingIpv6Type
	Ipv4Addr     *net.IPNet
	Ipv6Addr     *net.IPNet
}

func (o TwoStageL3ClosSubinterface) MarshalJSON() ([]byte, error) {
	var raw struct {
		Ipv4AddrType *string `json:"ipv4_addr_type,omitempty"`
		Ipv6AddrType *string `json:"ipv6_addr_type,omitempty"`
		Ipv4Addr     *string `json:"ipv4_addr,omitempty"`
		Ipv6Addr     *string `json:"ipv6_addr,omitempty"`
	}

	if o.Ipv4AddrType != nil {
		raw.Ipv4AddrType = toPtr(o.Ipv4AddrType.String())
	}

	if o.Ipv6AddrType != nil {
		raw.Ipv6AddrType = toPtr(o.Ipv6AddrType.String())
	}

	if o.Ipv4Addr != nil {
		if len(o.Ipv4Addr.IP) > 0 && len(o.Ipv4Addr.Mask) > 0 {
			raw.Ipv4Addr = toPtr(o.Ipv4Addr.String()) // send the string representation
		} else {
			raw.Ipv4Addr = toPtr("") // send an empty string to clear the value on the server
		}
	}

	if o.Ipv6Addr != nil {
		if len(o.Ipv6Addr.IP) > 0 && len(o.Ipv6Addr.Mask) > 0 {
			raw.Ipv6Addr = toPtr(o.Ipv6Addr.String()) // send the string representation
		} else {
			raw.Ipv6Addr = toPtr("") // send an empty string to clear the value on the server
		}
	}

	return json.Marshal(raw)
}

func (o *TwoStageL3ClosSubinterface) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Ipv4AddrType string  `json:"ipv4_addr_type,omitempty"`
		Ipv6AddrType string  `json:"ipv6_addr_type,omitempty"`
		Ipv4Addr     *string `json:"ipv4_addr,omitempty"`
		Ipv6Addr     *string `json:"ipv6_addr,omitempty"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Ipv4AddrType = new(InterfaceNumberingIpv4Type)
	err = o.Ipv4AddrType.FromString(raw.Ipv4AddrType)
	if err != nil {
		return fmt.Errorf("failed parsing ipv4_addr_type %q while unmarshaling TwoStageL3ClosSubinterface", raw.Ipv4AddrType)
	}

	o.Ipv6AddrType = new(InterfaceNumberingIpv6Type)
	err = o.Ipv6AddrType.FromString(raw.Ipv6AddrType)
	if err != nil {
		return fmt.Errorf("failed parsing ipv6_addr_type %q while unmarshaling TwoStageL3ClosSubinterface", raw.Ipv6AddrType)
	}

	if raw.Ipv4Addr != nil {
		ip, net, err := net.ParseCIDR(*raw.Ipv4Addr)
		if err != nil {
			return fmt.Errorf("failed parsing ipv4_addr while unmarshaling subinterface - %w", err)
		}
		net.IP = ip
		o.Ipv4Addr = net
	}

	if raw.Ipv6Addr != nil {
		ip, net, err := net.ParseCIDR(*raw.Ipv6Addr)
		if err != nil {
			return fmt.Errorf("failed parsing ipv6_addr while unmarshaling subinterface - %w", err)
		}
		net.IP = ip
		o.Ipv6Addr = net
	}

	return nil
}

func (o *TwoStageL3ClosClient) UpdateSubinterfaces(ctx context.Context, in map[ObjectId]TwoStageL3ClosSubinterface) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPatch,
		urlStr: fmt.Sprintf(apiUrlSubinterfaces, o.Id()),
		apiInput: &struct {
			Subinterfaces map[ObjectId]TwoStageL3ClosSubinterface `json:"subinterfaces"`
		}{
			Subinterfaces: in,
		},
	})

	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
