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
		LinkType        linkType               `json:"link_type"`
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
		Tags []ObjectId `json:"tags"`
	}
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	if len(raw.Endpoints) != 2 {
		return fmt.Errorf("got %d endpoints, expected 2", len(raw.Endpoints))
	}

	linkType, err := raw.LinkType.parse()
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformLinkData)
	o.Data.Speed = raw.Speed
	o.Data.Label = raw.Label
	o.Data.Type = LinkType(linkType)
	o.Data.AggregateLinkId = raw.AggregateLinkId
	o.Data.Endpoints[0] = FreeformEndpoint{
		SystemId:  raw.Endpoints[0].System.Id,
		Interface: *raw.Endpoints[0].Interface.Data,
	}
	o.Data.Endpoints[1] = FreeformEndpoint{
		SystemId:  raw.Endpoints[1].System.Id,
		Interface: *raw.Endpoints[1].Interface.Data,
	}
	o.Data.Tags = raw.Tags

	return nil
}

type FreeformLinkData struct {
	Type            LinkType
	AggregateLinkId *ObjectId
	Label           string
	Speed           LogicalDevicePortSpeed
	Tags            []ObjectId
	Endpoints       [2]FreeformEndpoint
}

var _ json.Marshaler = new(FreeformInterfaceData)

type FreeformInterfaceData struct {
	IfName           string
	TransformationId int
	Ipv4Address      *net.IPNet
	Ipv6Address      *net.IPNet
	Tags             []ObjectId
}

func (o FreeformInterfaceData) MarshalJSON() ([]byte, error) {
	var raw struct {
		IfName           string     `json:"if_name"`
		TransformationId int        `json:"transformation_id"`
		Ipv4Addr         string     `json:"ipv4_addr,omitempty"`
		Ipv6Addr         string     `json:"ipv6_addr,omitempty"`
		Tags             []ObjectId `json:"tags"`
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

var _ json.Unmarshaler = new(FreeformInterface)

type FreeformInterface struct {
	Id   ObjectId `json:"id"`
	Data *FreeformInterfaceData
}

func (o *FreeformInterface) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id               ObjectId   `json:"id"`
		IfName           string     `json:"if_name"`
		TransformationId int        `json:"transformation_id"`
		Ipv4Addr         *string    `json:"ipv4_addr"`
		Ipv6Addr         *string    `json:"ipv6_addr"`
		Tags             []ObjectId `json:"tags"`
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
	Interface FreeformInterfaceData
}

func (o *FreeformEndpoint) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		System struct {
			Id ObjectId `json:"id"`
		} `json:"system"`
		Interface FreeformInterfaceData `json:"interface"`
	}

	o.SystemId = raw.System.Id
	o.Interface = raw.Interface

	return nil
}

func (o FreeformEndpoint) MarshalJSON() ([]byte, error) {
	var raw struct {
		System struct {
			Id ObjectId `json:"id"`
		} `json:"system"`
		Interface FreeformInterfaceData `json:"interface"`
	}
	raw.System.Id = o.SystemId
	raw.Interface = o.Interface
	return json.Marshal(&raw)
}

type FreeformLinkRequest struct {
	Label     string              `json:"label"`
	Tags      []ObjectId          `json:"tags"`
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

func (o *FreeformClient) UpdateLink(ctx context.Context, id ObjectId, in *FreeformLinkData) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfLinkById, o.blueprintId, id),
		apiInput: in,
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
