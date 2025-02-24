// Copyright (c) Juniper Networks, Inc., 2022-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package apstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDesignInterfaceMapDigests        = apiUrlDesignPrefix + "interface-map-digests"
	apiUrlDesignInterfaceMapsDigestsPrefix = apiUrlDesignInterfaceMapDigests + apiUrlPathDelim
	apiUrlDesignInterfaceMapDigestById     = apiUrlDesignInterfaceMapsDigestsPrefix + "%s"
)

type InterfaceMapDigest struct {
	Id            ObjectId `json:"id"`
	Label         string   `json:"label"`
	LogicalDevice struct {
		Id    ObjectId `json:"id"`
		Label string   `json:"label"`
	} `json:"logical_device"`
	DeviceProfile struct {
		Id    ObjectId `json:"id"`
		Label string   `json:"label"`
	} `json:"device_profile"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

type InterfaceMapDigests []InterfaceMapDigest

// SupportsDeviceProfile returns bool indicating whether any InterfaceMapDigest
// in the slice indicates support for the given DeviceProfile ID
func (o *InterfaceMapDigests) SupportsDeviceProfile(id ObjectId) bool {
	for _, imd := range *o {
		if imd.DeviceProfile.Id == id {
			return true
		}
	}
	return false
}

// SupportsLogicalDevice returns bool indicating whether any InterfaceMapDigest
// in the slice indicates support for the given LogicalDevice ID
func (o *InterfaceMapDigests) SupportsLogicalDevice(id ObjectId) bool {
	for _, imd := range *o {
		if imd.LogicalDevice.Id == id {
			return true
		}
	}
	return false
}

func (o *Client) getInterfaceMapDigest(ctx context.Context, id ObjectId) (*InterfaceMapDigest, error) {
	response := &InterfaceMapDigest{}
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignInterfaceMapDigestById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getInterfaceMapDigests(ctx context.Context) (InterfaceMapDigests, error) {
	response := &struct {
		Items InterfaceMapDigests `json:"items"`
	}{}
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignInterfaceMapDigests,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (o *Client) getInterfaceMapDigestsByDeviceProfile(ctx context.Context, desired ObjectId) (InterfaceMapDigests, error) {
	response := &struct {
		Items InterfaceMapDigests `json:"items"`
	}{}
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignInterfaceMapDigests,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	for i := len(response.Items) - 1; i >= 0; i-- {
		if response.Items[i].DeviceProfile.Id != desired {
			response.Items[i] = response.Items[len(response.Items)-1] // move last item to position [i]
			response.Items = response.Items[:len(response.Items)-1]   // trim last element
		}
	}
	return response.Items, nil
}

func (o *Client) getInterfaceMapDigestsByLogicalDevice(ctx context.Context, desired ObjectId) (InterfaceMapDigests, error) {
	response := &struct {
		Items InterfaceMapDigests `json:"items"`
	}{}
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignInterfaceMapDigests,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	for i := len(response.Items) - 1; i >= 0; i-- {
		if response.Items[i].LogicalDevice.Id != desired {
			response.Items[i] = response.Items[len(response.Items)-1] // move last item to position [i]
			response.Items = response.Items[:len(response.Items)-1]   // trim last element
		}
	}
	return response.Items, nil
}

func (o *Client) getInterfaceMapDigestsLogicalDeviceAndDeviceProfile(ctx context.Context, ldId ObjectId, dpId ObjectId) (InterfaceMapDigests, error) {
	response := &struct {
		Items InterfaceMapDigests `json:"items"`
	}{}
	err := o.talkToApstra(ctx, talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignInterfaceMapDigests,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	for i := len(response.Items) - 1; i >= 0; i-- {
		if response.Items[i].LogicalDevice.Id != ldId || response.Items[i].DeviceProfile.Id != dpId {
			response.Items[i] = response.Items[len(response.Items)-1] // move last item to position [i]
			response.Items = response.Items[:len(response.Items)-1]   // trim last element
		}
	}
	return response.Items, nil
}
