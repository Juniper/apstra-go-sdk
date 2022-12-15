package goapstraw

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	apiUrlResourcesIp4Pools       = apiUrlResources + "/ip-pools"
	apiUrlResourcesIp4PoolsPrefix = apiUrlResourcesIp4Pools + apiUrlPathDelim
	apiUrlResourcesIp4PoolById    = apiUrlResourcesIp4PoolsPrefix + "%s"

	apiUrlResourcesIp6Pools       = apiUrlResources + "/ipv6-pools"
	apiUrlResourcesIp6PoolsPrefix = apiUrlResourcesIp6Pools + apiUrlPathDelim
	apiUrlResourcesIp6PoolById    = apiUrlResourcesIp6PoolsPrefix + "%s"
)

type IpPool struct {
	Id             ObjectId   `json:"id"`
	DisplayName    string     `json:"display_name"`
	Status         string     `json:"status"`
	Tags           []string   `json:"tags"`
	Used           big.Int    `json:"used"`
	Total          big.Int    `json:"total"`
	UsedPercentage float32    `json:"used_percentage"`
	CreatedAt      time.Time  `json:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at"`
	Subnets        []IpSubnet `json:"subnets"`
}

type rawIpPool struct {
	Id             ObjectId      `json:"id"`
	DisplayName    string        `json:"display_name"`
	Status         string        `json:"status"`
	Tags           []string      `json:"tags"`
	Used           string        `json:"used"`
	Total          string        `json:"total"`
	UsedPercentage float32       `json:"used_percentage"`
	CreatedAt      time.Time     `json:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at"`
	Subnets        []rawIpSubnet `json:"subnets"`
}

func (o *rawIpPool) polish() (*IpPool, error) {
	var used, total big.Int
	_, ok := used.SetString(o.Used, 10)
	if !ok {
		return nil, fmt.Errorf("failed parsing IP Pool field 'used' ('%s')", o.Used)
	}
	_, ok = total.SetString(o.Total, 10)
	if !ok {
		return nil, fmt.Errorf("failed parsing IP Pool field 'used' ('%s')", o.Total)
	}

	subnets := make([]IpSubnet, len(o.Subnets))
	for i, rs := range o.Subnets {
		ps, err := rs.polish()
		if err != nil {
			return nil, err
		}
		subnets[i] = *ps
	}

	return &IpPool{
		Id:             o.Id,
		DisplayName:    o.DisplayName,
		Status:         o.Status,
		Tags:           o.Tags,
		Used:           used,
		Total:          total,
		UsedPercentage: o.UsedPercentage,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Subnets:        subnets,
	}, nil
}

type IpSubnet struct {
	Network        *net.IPNet
	Status         string
	Used           big.Int
	Total          big.Int
	UsedPercentage float32
}

type rawIpSubnet struct {
	Network        string  `json:"network,omitempty"`
	Status         string  `json:"status,omitempty"`
	Used           string  `json:"used,omitempty"`
	Total          string  `json:"total,omitempty"`
	UsedPercentage float32 `json:"used_percentage"`
}

func (o *rawIpSubnet) polish() (*IpSubnet, error) {
	var used, total big.Int
	_, ok := used.SetString(o.Used, 10)
	if !ok {
		return nil, fmt.Errorf("failed parsing IP Pool field 'used' ('%s')", o.Used)
	}
	_, ok = total.SetString(o.Total, 10)
	if !ok {
		return nil, fmt.Errorf("failed parsing IP Pool field 'used' ('%s')", o.Total)
	}

	_, parsed, err := net.ParseCIDR(o.Network)
	if err != nil {
		return nil, fmt.Errorf("error parsing subnet string from apstra '%s' - %w", o.Network, err)
	}
	return &IpSubnet{
		Network:        parsed,
		Status:         o.Status,
		Used:           used,
		Total:          total,
		UsedPercentage: o.UsedPercentage,
	}, nil
}

type NewIpPoolRequest struct {
	DisplayName string        `json:"display_name"`
	Tags        []string      `json:"tags"`
	Subnets     []NewIpSubnet `json:"subnets"`
}

type NewIpSubnet struct {
	Network string `json:"network"`
}

func (o *Client) listIpPoolIds(ctx context.Context, urlStr string) ([]ObjectId, error) {
	var response struct {
		Items []ObjectId `json:"items"`
	}
	return response.Items, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      urlStr,
		apiResponse: &response,
	})
}

func (o *Client) getIpPools(ctx context.Context, urlStr string) ([]rawIpPool, error) {
	var response struct {
		Items []rawIpPool `json:"items"`
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      urlStr,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *Client) getIpPool(ctx context.Context, urlStr string) (*rawIpPool, error) {
	response := &rawIpPool{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      urlStr,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response, nil
}

func (o *Client) getIpPoolsByName(ctx context.Context, urlStr string, desired string) ([]rawIpPool, error) {
	pools, err := o.getIpPools(ctx, urlStr)
	if err != nil {
		return nil, err
	}

	i := len(pools) - 1
	for i >= 0 {
		if pools[i].DisplayName != desired { // undesired. delete element.
			pools[i] = pools[len(pools)-1] // swap last to current
			pools = pools[:len(pools)-1]   // delete last
		}
		i--
	}
	return pools, nil
}

func (o *Client) getIpPoolByName(ctx context.Context, urlStr string, desired string) (*rawIpPool, error) {
	pools, err := o.getIpPoolsByName(ctx, urlStr, desired)
	if err != nil {
		return nil, err
	}
	switch len(pools) {
	case 0:
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("no pool named '%s' found", desired),
		}
	case 1:
		return &pools[0], nil
	default:
		return nil, ApstraClientErr{
			errType: ErrMultipleMatch,
			err:     fmt.Errorf("name '%s' does not uniquely identify a single IPv4 pool", desired),
		}
	}
}

func (o *Client) createIpPool(ctx context.Context, ipv6 bool, request *NewIpPoolRequest) (ObjectId, error) {
	if request.Subnets == nil {
		request.Subnets = []NewIpSubnet{}
	}
	for _, s := range request.Subnets {
		if ipv6 && !strings.Contains(s.Network, ":") {
			return "", fmt.Errorf("network '%s' not compatible with IPv6 pool", s.Network)
		}
		if !ipv6 && strings.Contains(s.Network, ":") {
			return "", fmt.Errorf("network '%s' not compatible with IPv4 pool", s.Network)
		}
	}

	if request.Tags == nil {
		request.Tags = []string{}
	}

	var urlStr string
	if ipv6 {
		urlStr = apiUrlResourcesIp6Pools
	} else {
		urlStr = apiUrlResourcesIp4Pools
	}
	response := &objectIdResponse{}
	return response.Id, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      urlStr,
		apiInput:    request,
		apiResponse: response,
	})
}

func (o *Client) deleteIpPool(ctx context.Context, urlStr string, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodDelete,
		urlStr:   fmt.Sprintf(urlStr, id),
		apiInput: nil,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) updateIpPool(ctx context.Context, urlStr string, request *NewIpPoolRequest) error {
	if request.Subnets == nil {
		request.Subnets = []NewIpSubnet{}
	}
	if request.Tags == nil {
		request.Tags = []string{}
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   urlStr,
		apiInput: request,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) addSubnetToIpPool(ctx context.Context, poolId ObjectId, new *net.IPNet) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceIpPoolRangeMutex)
	defer o.unlock(clientApiResourceIpPoolRangeMutex)

	var urlStr string // url we'll use to request changes
	var pool *IpPool  // existing pool
	var err error
	if strings.Contains(new.String(), ":") {
		urlStr = fmt.Sprintf(apiUrlResourcesIp6PoolById, poolId)
		pool, err = o.GetIp6Pool(ctx, poolId)
		if err != nil {
			return fmt.Errorf("cannot fetch ip pool - %w", err)
		}
	} else {
		urlStr = fmt.Sprintf(apiUrlResourcesIp4PoolById, poolId)
		pool, err = o.GetIp4Pool(ctx, poolId)
		if err != nil {
			return fmt.Errorf("cannot fetch ip pool - %w", err)
		}
	}

	// check for subnet collisions while copying existing subnets to new request object
	subnets := []NewIpSubnet{{Network: new.String()}} // start the list with the new one
	for _, s := range pool.Subnets {
		old := s.Network
		if old.Contains(new.IP) || new.Contains(old.IP) {
			return fmt.Errorf("new subnet '%s' overlaps existing subnet %s'", new.String(), s.Network)
		}
		subnets = append(subnets, NewIpSubnet{Network: s.Network.String()})
	}

	err = o.updateIpPool(ctx, urlStr, &NewIpPoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
		Subnets:     subnets,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteSubnetFromIpPool(ctx context.Context, poolId ObjectId, target *net.IPNet) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceIpPoolRangeMutex)
	defer o.unlock(clientApiResourceIpPoolRangeMutex)

	var urlStr string // url we'll use to request changes
	var pool *IpPool  // existing pool
	var err error
	if strings.Contains(target.String(), ":") {
		urlStr = fmt.Sprintf(apiUrlResourcesIp6PoolById, poolId)
		pool, err = o.GetIp6Pool(ctx, poolId)
		if err != nil {
			return fmt.Errorf("cannot fetch ip pool - %w", err)
		}
	} else {
		urlStr = fmt.Sprintf(apiUrlResourcesIp4PoolById, poolId)
		pool, err = o.GetIp4Pool(ctx, poolId)
		if err != nil {
			return fmt.Errorf("cannot fetch ip pool - %w", err)
		}
	}

	// prep new request structure
	newRequest := &NewIpPoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
		Subnets:     []NewIpSubnet{}, // empty slice, but not nil for apstra
	}

	// work through the list copy non-target subnets to the new request
	var targetFound bool
	for _, s := range pool.Subnets {
		// copy old subnets which don't match deletion target to new request slice
		if s.Network.String() != target.String() {
			newRequest.Subnets = append(newRequest.Subnets, NewIpSubnet{Network: s.Network.String()})
		} else {
			targetFound = true
		}
	}

	if !targetFound {
		// nothing to do
		return ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("target '%s' not found in pool '%s'", target.String(), poolId),
		}
	}

	err = o.updateIpPool(ctx, urlStr, newRequest)
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
