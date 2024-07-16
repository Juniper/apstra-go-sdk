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

var (
	_ json.Marshaler   = (*TwoStageL3ClosSubinterface)(nil)
	_ json.Unmarshaler = (*TwoStageL3ClosSubinterface)(nil)
)

type TwoStageL3ClosSubinterface struct {
	Ipv4AddrType InterfaceNumberingIpv4Type
	Ipv6AddrType InterfaceNumberingIpv6Type
	Ipv4Addr     *net.IPNet
	Ipv6Addr     *net.IPNet
}

func (o TwoStageL3ClosSubinterface) MarshalJSON() ([]byte, error) {
	var raw struct {
		Ipv4AddrType *string `json:"ipv4_addr_type"`
		Ipv6AddrType *string `json:"ipv6_addr_type"`
		Ipv4Addr     *string `json:"ipv4_addr"`
		Ipv6Addr     *string `json:"ipv6_addr"`
	}

	if o.Ipv4AddrType != InterfaceNumberingIpv4TypeNone {
		raw.Ipv4AddrType = toPtr(o.Ipv4AddrType.String())
	}

	if o.Ipv6AddrType != InterfaceNumberingIpv6TypeNone {
		raw.Ipv6AddrType = toPtr(o.Ipv6AddrType.String())
	}

	if o.Ipv4Addr != nil {
		raw.Ipv4Addr = toPtr(o.Ipv4Addr.String()) // send the string representation
	}

	if o.Ipv6Addr != nil {
		raw.Ipv6Addr = toPtr(o.Ipv6Addr.String()) // send the string representation
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

	err = o.Ipv4AddrType.FromString(raw.Ipv4AddrType)
	if err != nil {
		return fmt.Errorf("failed parsing ipv4_addr_type %q while unmarshaling TwoStageL3ClosSubinterface", raw.Ipv4AddrType)
	}

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

func (o *TwoStageL3ClosClient) GetSubinterface(ctx context.Context, id ObjectId) (*TwoStageL3ClosSubinterface, error) {
	var node struct {
		Type     nodeType `json:"type"`
		IfType   string   `json:"if_type"`
		Ipv4Addr *string  `json:"ipv4_addr"`
		Ipv4Type *string  `json:"ipv4_addr_type"`
		Ipv6Addr *string  `json:"ipv6_addr"`
		Ipv6Type *string  `json:"ipv6_addr_type"`
	}

	err := o.Client().GetNode(ctx, o.Id(), id, &node)
	if err != nil {
		return nil, err
	}

	if node.Type != nodeTypeInterface {
		return nil, fmt.Errorf("node %q exists but has type %q - expected %q", id, node.Type, nodeTypeInterface)
	}

	if node.IfType != "subinterface" {
		return nil, fmt.Errorf("interface node %q has if_type %q - expected \"subinterface\"", id, node.IfType)
	}

	result := TwoStageL3ClosSubinterface{
		Ipv4AddrType: InterfaceNumberingIpv4TypeNone,
		Ipv6AddrType: InterfaceNumberingIpv6TypeNone,
	}

	if node.Ipv4Type != nil {
		err = result.Ipv4AddrType.FromString(*node.Ipv4Type)
		if err != nil {
			return nil, fmt.Errorf("cannot parse node %q ipv4_addr_type value %q - %w", id, *node.Ipv4Type, err)
		}
	}

	if node.Ipv6Type != nil {
		err = result.Ipv6AddrType.FromString(*node.Ipv6Type)
		if err != nil {
			return nil, fmt.Errorf("cannot parse node %q ipv6_addr_type value %q - %w", id, *node.Ipv6Type, err)
		}
	}

	if node.Ipv4Addr != nil {
		var ip net.IP
		ip, result.Ipv4Addr, err = net.ParseCIDR(*node.Ipv4Addr)
		if err != nil {
			return nil, fmt.Errorf("cannot parse node %q ipv4_addr value %q - %w", id, *node.Ipv4Addr, err)
		}
		result.Ipv4Addr.IP = ip
	}

	if node.Ipv6Addr != nil {
		var ip net.IP
		ip, result.Ipv6Addr, err = net.ParseCIDR(*node.Ipv6Addr)
		if err != nil {
			return nil, fmt.Errorf("cannot parse node %q ipv6_addr value %q - %w", id, *node.Ipv6Addr, err)
		}
		result.Ipv6Addr.IP = ip
	}

	return &result, nil
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
