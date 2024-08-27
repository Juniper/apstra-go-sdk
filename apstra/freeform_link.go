package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	apiUrlFfLinks    = apiUrlBlueprintById + apiUrlPathDelim + "links"
	apiUrlFfLinkById = apiUrlFfLinks + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformLink)

type FreeformLink struct {
	Id   ObjectId
	Data *FreeformLinkData
}

func (o *FreeformLink) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id              ObjectId               `json:"id"`
		Speed           LogicalDevicePortSpeed `json:"speed"`
		LinkType        string                 `json:"link_type"`
		Label           string                 `json:"label"`
		AggregateLinkId *ObjectId              `json:"aggregate_link_id"`
		Endpoints       []struct {
			System struct {
				Id         ObjectId   `json:"id"`
				Label      string     `json:"label"`
				SystemType systemType `json:"system_type"`
			} `json:"system"`
			Interface FreeformInterface `json:"interface"`
		} `json:"endpoints"`
		Tags []string `json:"tags"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	if len(raw.Endpoints) != 2 {
		return fmt.Errorf("got %d endpoints, expected 2", len(raw.Endpoints))
	}

	o.Id = raw.Id
	o.Data = new(FreeformLinkData)
	o.Data.Speed = raw.Speed
	o.Data.Label = raw.Label
	err = o.Data.Type.FromString(raw.LinkType)
	if err != nil {
		return err
	}
	o.Data.AggregateLinkId = raw.AggregateLinkId
	o.Data.Endpoints[0] = FreeformEndpoint{
		SystemId: raw.Endpoints[0].System.Id,
		Interface: FreeformInterface{
			Id:   raw.Endpoints[0].Interface.Id,
			Data: raw.Endpoints[0].Interface.Data,
		},
	}
	o.Data.Endpoints[1] = FreeformEndpoint{
		SystemId: raw.Endpoints[1].System.Id,
		Interface: FreeformInterface{
			Id:   raw.Endpoints[1].Interface.Id,
			Data: raw.Endpoints[1].Interface.Data,
		},
	}
	o.Data.Tags = raw.Tags

	return nil
}

var _ json.Marshaler = new(FreeformLinkData)

type FreeformLinkData struct {
	Type            FFLinkType
	AggregateLinkId *ObjectId
	Label           string
	Speed           LogicalDevicePortSpeed
	Tags            []string
	Endpoints       [2]FreeformEndpoint
}

func (o FreeformLinkData) MarshalJSON() ([]byte, error) {
	type rawEndpointInterface struct {
		IfName           string   `json:"if_name"`
		TransformationId int      `json:"transformation_id"`
		Ipv4Addr         *string  `json:"ipv4_addr"`
		Ipv6Addr         *string  `json:"ipv6_addr"`
		Tags             []string `json:"tags"`
	}
	type rawEndpoint struct {
		Interface rawEndpointInterface `json:"interface"`
	}
	var raw struct {
		Label     string         `json:"label"`
		Tags      []string       `json:"tags"`
		Endpoints [2]rawEndpoint `json:"endpoints"`
	}

	raw.Label = o.Label
	raw.Tags = o.Tags
	raw.Endpoints[0] = rawEndpoint{
		Interface: rawEndpointInterface{
			IfName:           o.Endpoints[0].Interface.Data.IfName,
			TransformationId: o.Endpoints[0].Interface.Data.TransformationId,
			Tags:             o.Endpoints[0].Interface.Data.Tags,
		},
	}
	if o.Endpoints[0].Interface.Data.Ipv4Address != nil {
		raw.Endpoints[0].Interface.Ipv4Addr = toPtr(o.Endpoints[0].Interface.Data.Ipv4Address.String())
	}
	if o.Endpoints[0].Interface.Data.Ipv6Address != nil {
		raw.Endpoints[0].Interface.Ipv6Addr = toPtr(o.Endpoints[0].Interface.Data.Ipv6Address.String())
	}

	raw.Endpoints[1] = rawEndpoint{
		Interface: rawEndpointInterface{
			IfName:           o.Endpoints[1].Interface.Data.IfName,
			TransformationId: o.Endpoints[1].Interface.Data.TransformationId,
			Tags:             o.Endpoints[1].Interface.Data.Tags,
		},
	}
	if o.Endpoints[1].Interface.Data.Ipv4Address != nil {
		raw.Endpoints[1].Interface.Ipv4Addr = toPtr(o.Endpoints[1].Interface.Data.Ipv4Address.String())
	}
	if o.Endpoints[1].Interface.Data.Ipv6Address != nil {
		raw.Endpoints[1].Interface.Ipv6Addr = toPtr(o.Endpoints[1].Interface.Data.Ipv6Address.String())
	}

	return json.Marshal(raw)
}

var _ json.Marshaler = new(FreeformInterfaceData)

type FreeformInterfaceData struct {
	IfName           string
	TransformationId int
	Ipv4Address      *net.IPNet
	Ipv6Address      *net.IPNet
	Tags             []string
}

func (o FreeformInterfaceData) MarshalJSON() ([]byte, error) {
	var raw struct {
		IfName           string   `json:"if_name"`
		TransformationId int      `json:"transformation_id"`
		Ipv4Addr         string   `json:"ipv4_addr,omitempty"`
		Ipv6Addr         string   `json:"ipv6_addr,omitempty"`
		Tags             []string `json:"tags"`
	}

	raw.IfName = o.IfName
	raw.TransformationId = o.TransformationId
	if o.Ipv4Address != nil {
		raw.Ipv4Addr = o.Ipv4Address.String()
		if strings.Contains(raw.Ipv4Addr, "<nil>") {
			return nil, fmt.Errorf("cannot marshall ipv4 address - %s", raw.Ipv4Addr)
		}
	}
	if o.Ipv6Address != nil {
		raw.Ipv6Addr = o.Ipv6Address.String()
		if strings.Contains(raw.Ipv6Addr, "<nil>") {
			return nil, fmt.Errorf("cannot marshall ipv6 address - %s", raw.Ipv6Addr)
		}
	}
	raw.Tags = o.Tags

	return json.Marshal(&raw)
}

var (
	_ json.Unmarshaler = new(FreeformInterface)
	_ json.Marshaler   = new(FreeformInterface)
)

type FreeformInterface struct {
	Id   *ObjectId
	Data *FreeformInterfaceData
}

func (o *FreeformInterface) MarshalJSON() ([]byte, error) {
	var raw struct {
		Id               *ObjectId `json:"id"`
		IfName           string    `json:"if_name"`
		TransformationId int       `json:"transformation_id"`
		Ipv4Addr         *string   `json:"ipv4_addr"`
		Ipv6Addr         *string   `json:"ipv6_addr"`
		Tags             []string  `json:"tags"`
	}
	raw.Id = o.Id
	raw.IfName = o.Data.IfName
	raw.TransformationId = o.Data.TransformationId
	if o.Data.Ipv4Address != nil {
		raw.Ipv4Addr = toPtr(o.Data.Ipv4Address.String())
	}
	if o.Data.Ipv6Address != nil {
		raw.Ipv6Addr = toPtr(o.Data.Ipv6Address.String())
	}
	raw.Tags = o.Data.Tags
	return json.Marshal(raw)
}

func (o *FreeformInterface) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id               *ObjectId `json:"id"`
		IfName           string    `json:"if_name"`
		TransformationId int       `json:"transformation_id"`
		Ipv4Addr         *string   `json:"ipv4_addr"`
		Ipv6Addr         *string   `json:"ipv6_addr"`
		Tags             []string  `json:"tags"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformInterfaceData)
	if raw.Ipv4Addr != nil {
		ip, net4, err := net.ParseCIDR(*raw.Ipv4Addr)
		if err != nil {
			return fmt.Errorf("failed parsing IPv4 API response - %w", err)
		}
		net4.IP = ip
		o.Data.Ipv4Address = net4
	}

	if raw.Ipv6Addr != nil {
		ip, net6, err := net.ParseCIDR(*raw.Ipv6Addr)
		if err != nil {
			return fmt.Errorf("failed parsing IPv6 API response - %w", err)
		}
		net6.IP = ip
		o.Data.Ipv6Address = net6
	}
	o.Data.IfName = raw.IfName
	o.Data.TransformationId = raw.TransformationId
	o.Data.Tags = raw.Tags

	return nil
}

var (
	_ json.Marshaler   = new(FreeformEndpoint)
	_ json.Unmarshaler = new(FreeformEndpoint)
)

type FreeformEndpoint struct {
	SystemId  ObjectId
	Interface FreeformInterface
}

func (o *FreeformEndpoint) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		System struct {
			Id ObjectId `json:"id"`
		} `json:"system"`
		Interface *FreeformInterfaceData `json:"interface"`
	}

	o.SystemId = raw.System.Id
	o.Interface.Data = raw.Interface

	return json.Unmarshal(bytes, &raw)
}

func (o FreeformEndpoint) MarshalJSON() ([]byte, error) {
	var raw struct {
		System *struct {
			Id ObjectId `json:"id"`
		} `json:"system,omitempty"`
		Interface struct {
			Id          *ObjectId `json:"id,omitempty"`
			IfName      string    `json:"if_name"`
			TransformId int       `json:"transformation_id"`
			Ipv4Addr    *string   `json:"ipv4_addr"`
			Ipv6Addr    *string   `json:"ipv6_addr"`
			Tags        []string  `json:"tags"`
		} `json:"interface"`
	}
	if o.SystemId != "" {
		raw.System = new(struct {
			Id ObjectId `json:"id"`
		})
		raw.System.Id = o.SystemId
	}
	raw.Interface.Id = o.Interface.Id
	raw.Interface.IfName = o.Interface.Data.IfName
	if o.Interface.Data.Ipv4Address != nil {
		raw.Interface.Ipv4Addr = toPtr(o.Interface.Data.Ipv4Address.String())
	}
	if o.Interface.Data.Ipv6Address != nil {
		raw.Interface.Ipv6Addr = toPtr(o.Interface.Data.Ipv6Address.String())
	}
	raw.Interface.TransformId = o.Interface.Data.TransformationId
	raw.Interface.Tags = o.Interface.Data.Tags

	return json.Marshal(raw)
}

type FreeformLinkRequest struct {
	Label     string              `json:"label"`
	Tags      []string            `json:"tags"`
	Endpoints [2]FreeformEndpoint `json:"endpoints"`
}

func (o *FreeformClient) CreateLink(ctx context.Context, in *FreeformLinkRequest) (ObjectId, error) {
	var response objectIdResponse

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfLinks, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetLink(ctx context.Context, id ObjectId) (*FreeformLink, error) {
	var response FreeformLink

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfLinkById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) GetLinkByName(ctx context.Context, name string) (*FreeformLink, error) {
	all, err := o.GetAllLinks(ctx)
	if err != nil {
		return nil, err
	}

	var result *FreeformLink
	for _, link := range all {
		link := link
		if link.Data.Label == name {
			if result != nil {
				return nil, ClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple links in blueprint %q have name %q", o.client.id, name),
				}
			}

			result = &link
		}
	}

	if result == nil {
		return nil, ClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no link in blueprint %q has name %q", o.client.id, name),
		}
	}

	return result, nil
}

func (o *FreeformClient) GetAllLinks(ctx context.Context) ([]FreeformLink, error) {
	var response struct {
		Items []FreeformLink `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfLinks, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) UpdateLink(ctx context.Context, id ObjectId, in *FreeformLinkRequest) error {
	// make a copy to be certain to clear the systemId without damaging the caller's struct.
	copy := *in
	copy.Endpoints[0].SystemId = ""
	copy.Endpoints[1].SystemId = ""

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfLinkById, o.blueprintId, id),
		apiInput: &copy,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteLink(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfLinkById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
