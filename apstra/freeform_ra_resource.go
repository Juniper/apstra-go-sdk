package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

const (
	apiUrlFfRaResources    = apiUrlBlueprintById + apiUrlPathDelim + "ra-resources"
	apiUrlFfRaResourceById = apiUrlFfRaResources + apiUrlPathDelim + "%s"
)

var _ json.Unmarshaler = new(FreeformRaResource)

type FreeformRaResource struct {
	Id   ObjectId
	Data *FreeformRaResourceData
}

func (o *FreeformRaResource) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		Id              ObjectId  `json:"id"`
		ResourceType    string    `json:"resource_type"`
		Label           string    `json:"label"`
		Value           *string   `json:"value"`
		AllocatedFrom   *ObjectId `json:"allocated_from"`
		GroupId         ObjectId  `json:"group_id"`
		SubnetPrefixLen *int      `json:"subnet_prefix_len"`
		GeneratorId     *ObjectId `json:"generator_id"`
	}

	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	o.Id = raw.Id
	o.Data = new(FreeformRaResourceData)
	o.Data.Label = raw.Label
	o.Data.Value = raw.Value
	o.Data.AllocatedFrom = raw.AllocatedFrom
	o.Data.GroupId = raw.GroupId
	o.Data.SubnetPrefixLen = raw.SubnetPrefixLen
	o.Data.GeneratorId = raw.GeneratorId
	err = o.Data.ResourceType.FromString(raw.ResourceType)
	if err != nil {
		return err
	}

	return nil
}

var _ json.Marshaler = new(FreeformRaResourceData)

type FreeformRaResourceData struct {
	ResourceType    FFResourceType
	Label           string
	Value           *string
	AllocatedFrom   *ObjectId
	GroupId         ObjectId
	SubnetPrefixLen *int
	GeneratorId     *ObjectId
}

func (o FreeformRaResourceData) MarshalJSON() ([]byte, error) {
	var raw struct {
		ResourceType    string    `json:"resource_type"`
		Label           string    `json:"label"`
		Value           *string   `json:"value"`
		AllocatedFrom   *ObjectId `json:"allocated_from"`
		GroupId         ObjectId  `json:"group_id"`
		SubnetPrefixLen *int      `json:"subnet_prefix_len"`
		GeneratorId     *ObjectId `json:"generator_id"`
	}

	raw.ResourceType = o.ResourceType.String()
	raw.Label = o.Label
	raw.Value = o.Value
	raw.AllocatedFrom = o.AllocatedFrom
	raw.GroupId = o.GroupId
	raw.SubnetPrefixLen = o.SubnetPrefixLen
	raw.GeneratorId = o.GeneratorId

	return json.Marshal(&raw)
}

func (o FreeformRaResourceData) validate() error {
	switch o.ResourceType.String() {
	case FFResourceTypeAsn.String(), FFResourceTypeInt.String(), FFResourceTypeVni.String(), FFResourceTypeVlan.String():
		if o.Value != nil {
			_, err := strconv.Atoi(*o.Value)
			if err != nil {
				return fmt.Errorf("value cannot be %q when resource type is %q - %w", *o.Value, o.ResourceType, err)
			}
		}
		if o.SubnetPrefixLen != nil {
			return fmt.Errorf("subnet prefix len must not be specified when resource type is %q", o.ResourceType)
		}
	case FFResourceTypeHostIpv4.String():
		if o.Value != nil {
			ip, _, err := net.ParseCIDR(*o.Value)
			if err != nil {
				return fmt.Errorf("value cannot be %q when resource type is %q - %w", *o.Value, o.ResourceType, err)
			}
			ip = ip.To4()
			if ip == nil {
				return fmt.Errorf("value %q must be an ipv4 address in cidr notation when resource type is %q", *o.Value, o.ResourceType)
			}
		}
		if o.SubnetPrefixLen != nil {
			return fmt.Errorf("subnet prefix len must not be specified when resource type is %q", o.ResourceType)
		}
	case FFResourceTypeHostIpv6.String():
		if o.Value != nil {
			ip, _, err := net.ParseCIDR(*o.Value)
			if err != nil {
				return fmt.Errorf("value cannot be %q when resource type is %q - %w", *o.Value, o.ResourceType, err)
			}
			shouldBeNil := ip.To4()
			if shouldBeNil != nil {
				return fmt.Errorf("value %q must be an ipv6 address in cidr notation when resource type is %q", *o.Value, o.ResourceType)
			}
		}
		if o.SubnetPrefixLen != nil {
			return fmt.Errorf("subnet prefix len must not be specified when resource type is %q", o.ResourceType)
		}
	case FFResourceTypeIpv4.String():
		var ip net.IP
		var ipNet *net.IPNet
		var err error
		if o.Value != nil {
			ip, ipNet, err = net.ParseCIDR(*o.Value)
			if err != nil {
				return fmt.Errorf("value cannot be %q when resource type is %q - %w", *o.Value, o.ResourceType, err)
			}
			ip = ip.To4()
			if ip == nil {
				return fmt.Errorf("value %q must be an ipv4 address in cidr notation when resource type is %q", *o.Value, o.ResourceType)
			}
			if ipNet.IP.String() != ip.String() {
				return errors.New("value must be the base subnet address")
			}
			if o.SubnetPrefixLen != nil {
				ones, _ := ipNet.Mask.Size()
				if ones != *o.SubnetPrefixLen {
					return errors.New("subnetPrefixLen must be the same as the specified netmask")
				}
			}
		}
	case FFResourceTypeIpv6.String():
		var ip net.IP
		var ipNet *net.IPNet
		var err error
		if o.Value != nil {
			ip, ipNet, err = net.ParseCIDR(*o.Value)
			if err != nil {
				return fmt.Errorf("value cannot be %q when resource type is %q - %w", *o.Value, o.ResourceType, err)
			}
			shouldBeNil := ip.To4()
			if shouldBeNil != nil {
				return fmt.Errorf("value %q must be an ipv6 address in cidr notation when resource type is %q", *o.Value, o.ResourceType)
			}
			if ipNet.IP.String() != ip.String() {
				return errors.New("value must be the base subnet address")
			}
			if o.SubnetPrefixLen != nil {
				ones, _ := ipNet.Mask.Size()
				if ones != *o.SubnetPrefixLen {
					return errors.New("subnetPrefixLen must be the same as the specified netmask")
				}
			}
		}
	}
	return nil
}

func (o *FreeformClient) CreateRaResource(ctx context.Context, in *FreeformRaResourceData) (ObjectId, error) {
	err := in.validate()
	if err != nil {
		return "", err
	}

	var response objectIdResponse

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlFfRaResources, o.blueprintId),
		apiInput:    in,
		apiResponse: &response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}

	return response.Id, nil
}

func (o *FreeformClient) GetAllRaResources(ctx context.Context) ([]FreeformRaResource, error) {
	var response struct {
		Items []FreeformRaResource `json:"items"`
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaResources, o.blueprintId),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *FreeformClient) GetRaResource(ctx context.Context, id ObjectId) (*FreeformRaResource, error) {
	var response FreeformRaResource

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlFfRaResourceById, o.blueprintId, id),
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return &response, nil
}

func (o *FreeformClient) UpdateRaResource(ctx context.Context, id ObjectId, in *FreeformRaResourceData) error {
	err := in.validate()
	if err != nil {
		return err
	}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlFfRaResourceById, o.blueprintId, id),
		apiInput: in,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}

func (o *FreeformClient) DeleteRaResource(ctx context.Context, id ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlFfRaResourceById, o.blueprintId, id),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
