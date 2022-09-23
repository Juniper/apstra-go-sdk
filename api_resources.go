package goapstra

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrlResources               = "/api/resources"
	apiUrlResourcesAsnPools       = apiUrlResources + "/asn-pools"
	apiUrlResourcesAsnPoolsPrefix = apiUrlResourcesAsnPools + apiUrlPathDelim
	apiUrlResourcesAsnPoolById    = apiUrlResourcesAsnPoolsPrefix + "%s"
	apiUrlResourcesVniPools       = apiUrlResources + "/vni-pools"
	apiUrlResourcesVniPoolsPrefix = apiUrlResourcesVniPools + apiUrlPathDelim
	apiUrlResourcesVniPoolById    = apiUrlResourcesVniPoolsPrefix + "%s"
	apiUrlResourcesIpPools        = apiUrlResources + "/ip-pools"
	apiUrlResourcesIpPoolsPrefix  = apiUrlResourcesIpPools + apiUrlPathDelim
	apiUrlResourcesIpPoolById     = apiUrlResourcesIpPoolsPrefix + "%s"
)

// AsnPoolRequest is the public structure used to create/update an ASN pool.
type AsnPoolRequest IntPoolRequest

// AsnPool is the public structure used to convey query responses about ASN
// pools.
type AsnPool IntPool

// polish turns a rawAsnPool from the API into AsnPool for caller consumption
func (o *rawIntPool) makeAsnPool() (*AsnPool, error) {
	r, err := o.polish()
	return (*AsnPool)(r), err
}

func (o *Client) createAsnPool(ctx context.Context, in *AsnPoolRequest) (ObjectId, error) {
	id, err := o.createIntPool(ctx, (*IntPoolRequest)(in), apiUrlResourcesAsnPools)
	return id, err
}

func (o *Client) listAsnPoolIds(ctx context.Context) ([]ObjectId, error) {
	r, err := o.listIntPoolIds(ctx, apiUrlResourcesAsnPools)
	return r, err
}

func (o *Client) getAsnPools(ctx context.Context) ([]AsnPool, error) {
	r, err := o.getIntPools(ctx, apiUrlResourcesAsnPools)
	var r1 []AsnPool
	if err != nil {
		return r1, err
	}
	for _, i := range r {
		a, err := i.makeAsnPool()
		if err != nil {
			return r1, err
		}
		r1 = append(r1, *a)
	}
	return r1, err
}

func (o *Client) getAsnPool(ctx context.Context, poolId ObjectId) (*AsnPool, error) {
	r, err := o.getIntPool(ctx, apiUrlResourcesAsnPoolById, poolId)
	if err != nil {
		return nil, err
	}
	return r.makeAsnPool()
}

func (o *Client) getAsnPoolByName(ctx context.Context, desired string) (*AsnPool, error) {
	pools, err := o.getAsnPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching all ASN pools - %w", err)
	}
	found := -1
	for i, pool := range pools {
		if pool.DisplayName == desired {
			if found >= 0 {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("name '%s' does not uniquely identify an ASN pool", desired),
				}
			}
			found = i
		}
	}
	if found < 0 {
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("pool named '%s' not found", desired),
		}
	}
	return &pools[found], nil
}

func (o *Client) deleteAsnPool(ctx context.Context, poolId ObjectId) error {
	return o.deleteIntPool(ctx, apiUrlResourcesAsnPoolById, poolId)
}

func (o *Client) updateAsnPool(ctx context.Context, poolId ObjectId, pool *AsnPoolRequest) error {
	return o.updateIntPool(ctx, apiUrlResourcesAsnPoolById, poolId, (*IntPoolRequest)(pool))
}

func (o *Client) createAsnPoolRange(ctx context.Context, poolId ObjectId, newRange IntfIntRange) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)
	return o.createIntPoolRange(ctx, apiUrlResourcesAsnPoolById, poolId, newRange)
}

func (o *Client) asnPoolRangeExists(ctx context.Context, poolId ObjectId, asnRange IntfIntRange) (bool, error) {
	return o.IntPoolRangeExists(ctx, apiUrlResourcesAsnPoolById, poolId, asnRange)
}

func (o *Client) deleteAsnPoolRange(ctx context.Context, poolId ObjectId, deleteMe IntfIntRange) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)
	return o.deleteIntPoolRange(ctx, apiUrlResourcesAsnPoolById, poolId, deleteMe)
}

// This code will take care of vni
// VniPool is the public structure used to convey query responses about Vni
// pools.
// AsnPoolRequest is the public structure used to create/update an ASN pool.
type VniPoolRequest IntPoolRequest

// raw() converts an AsnPoolRequest to rawAsnPoolRequest for consumption by the
// Apstra API.
func (o *VniPoolRequest) raw() *rawIntPoolRequest {
	return o.raw()
}

type VniPool IntPool

// polish turns a rawVniPool from the API into VniPool for caller consumption
func (o *rawIntPool) makeVniPool() (*VniPool, error) {
	r, err := o.polish()
	return (*VniPool)(r), err
}

func (o *Client) createVniPool(ctx context.Context, in *VniPoolRequest) (ObjectId, error) {
	id, err := o.createIntPool(ctx, (*IntPoolRequest)(in), apiUrlResourcesVniPools)
	return id, err
}

func (o *Client) listVniPoolIds(ctx context.Context) ([]ObjectId, error) {
	r, err := o.listIntPoolIds(ctx, apiUrlResourcesVniPools)
	return r, err
}

func (o *Client) getVniPools(ctx context.Context) ([]VniPool, error) {
	r, err := o.getIntPools(ctx, apiUrlResourcesVniPools)
	var r1 []VniPool
	if err != nil {
		return r1, err
	}
	for _, i := range r {
		a, err := i.makeVniPool()
		if err != nil {
			return r1, err
		}
		r1 = append(r1, *a)
	}
	return r1, err
}

func (o *Client) getVniPool(ctx context.Context, poolId ObjectId) (*VniPool, error) {
	r, err := o.getIntPool(ctx, apiUrlResourcesVniPoolById, poolId)
	if err != nil {
		return nil, err
	}
	return r.makeVniPool()
}

func (o *Client) deleteVniPool(ctx context.Context, poolId ObjectId) error {
	return o.deleteIntPool(ctx, apiUrlResourcesVniPoolById, poolId)
}

func (o *Client) updateVniPool(ctx context.Context, poolId ObjectId, pool *VniPoolRequest) error {
	return o.updateIntPool(ctx, apiUrlResourcesVniPoolById, poolId, (*IntPoolRequest)(pool))
}

func (o *Client) createVniPoolRange(ctx context.Context, poolId ObjectId, newRange IntfIntRange) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceVniPoolRangeMutex)
	defer o.unlock(clientApiResourceVniPoolRangeMutex)
	return o.createIntPoolRange(ctx, apiUrlResourcesVniPoolById, poolId, newRange)
}

func (o *Client) vniPoolRangeExists(ctx context.Context, poolId ObjectId, VniRange IntfIntRange) (bool, error) {
	return o.IntPoolRangeExists(ctx, apiUrlResourcesVniPoolById, poolId, VniRange)
}

func (o *Client) deleteVniPoolRange(ctx context.Context, poolId ObjectId, deleteMe IntfIntRange) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceVniPoolRangeMutex)
	defer o.unlock(clientApiResourceVniPoolRangeMutex)
	return o.deleteIntPoolRange(ctx, apiUrlResourcesVniPoolById, poolId, deleteMe)
}

type Ip4Pool struct {
	Id             ObjectId    `json:"id"`
	DisplayName    string      `json:"display_name"`
	Status         string      `json:"status"`
	Tags           []string    `json:"tags"`
	Used           int64       `json:"used"`
	Total          int64       `json:"total"`
	UsedPercentage float32     `json:"used_percentage"`
	CreatedAt      time.Time   `json:"created_at"`
	LastModifiedAt time.Time   `json:"last_modified_at"`
	Subnets        []Ip4Subnet `json:"subnets"`
}

type rawIp4Pool struct {
	Id             ObjectId       `json:"id"`
	DisplayName    string         `json:"display_name"`
	Status         string         `json:"status"`
	Tags           []string       `json:"tags"`
	Used           string         `json:"used"`
	Total          string         `json:"total"`
	UsedPercentage float32        `json:"used_percentage"`
	CreatedAt      time.Time      `json:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at"`
	Subnets        []rawIp4Subnet `json:"subnets"`
}

func (o *rawIp4Pool) polish() (*Ip4Pool, error) {
	used, err := strconv.ParseInt(o.Used, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing IP Pool field 'used' ('%s') - %w", o.Used, err)
	}

	total, err := strconv.ParseInt(o.Total, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing IP Pool field 'total' ('%s') - %w", o.Total, err)
	}

	var subnets []Ip4Subnet
	for _, rs := range o.Subnets {
		ps, err := rs.polish()
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, *ps)
	}

	return &Ip4Pool{
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

type Ip4Subnet struct {
	Network        *net.IPNet
	Status         string
	Used           int64
	Total          int64
	UsedPercentage float32
}

type rawIp4Subnet struct {
	Network        string  `json:"network,omitempty"`
	Status         string  `json:"status,omitempty"`
	Used           string  `json:"used,omitempty"`
	Total          string  `json:"total,omitempty"`
	UsedPercentage float32 `json:"used_percentage"`
}

func (o *rawIp4Subnet) polish() (*Ip4Subnet, error) {
	used, err := strconv.ParseInt(o.Used, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing subnet field 'used' ('%s') - %w", o.Used, err)
	}

	total, err := strconv.ParseInt(o.Total, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing subnet field 'total' ('%s') - %w", o.Total, err)
	}

	_, parsed, err := net.ParseCIDR(o.Network)
	if err != nil {
		return nil, fmt.Errorf("error parsing subnet string from apstra '%s' - %w", o.Network, err)
	}
	return &Ip4Subnet{
		Network:        parsed,
		Status:         o.Status,
		Used:           used,
		Total:          total,
		UsedPercentage: o.UsedPercentage,
	}, nil
}

type NewIp4PoolRequest struct {
	DisplayName string         `json:"display_name"`
	Tags        []string       `json:"tags"`
	Subnets     []NewIp4Subnet `json:"subnets"`
}

type NewIp4Subnet struct {
	Network string `json:"network"`
}

func (o *Client) listIp4PoolIds(ctx context.Context) ([]ObjectId, error) {
	var response struct {
		Items []ObjectId `json:"items"`
	}
	return response.Items, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlResourcesIpPools,
		apiResponse: &response,
	})
}

func (o *Client) getIp4Pools(ctx context.Context) ([]Ip4Pool, error) {
	var response struct {
		Items []rawIp4Pool `json:"items"`
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlResourcesIpPools,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	var polishedPools []Ip4Pool
	for _, rp := range response.Items {
		pp, err := rp.polish()
		if err != nil {
			return nil, fmt.Errorf("error parsing raw pool content - %w", err)
		}
		polishedPools = append(polishedPools, *pp)
	}
	return polishedPools, nil
}

func (o *Client) getIp4Pool(ctx context.Context, poolId ObjectId) (*Ip4Pool, error) {
	response := &rawIp4Pool{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlResourcesIpPoolById, poolId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.polish()
}

func (o *Client) getIp4PoolByName(ctx context.Context, desired string) (*Ip4Pool, error) {
	pools, err := o.getIp4Pools(ctx)
	if err != nil {
		return nil, err
	}
	found := -1
	for i, pool := range pools {
		if pool.DisplayName == desired {
			if found >= 0 {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("name '%s' does not uniquely identify an IPv4 pool", desired),
				}
			}
			found = i
		}
	}
	if found < 0 {
		return nil, ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("pool named '%s' not found", desired),
		}
	}
	return &pools[found], nil
}

func (o *Client) createIp4Pool(ctx context.Context, request *NewIp4PoolRequest) (ObjectId, error) {
	if request.Subnets == nil {
		request.Subnets = []NewIp4Subnet{}
	}
	if request.Tags == nil {
		request.Tags = []string{}
	}
	response := &objectIdResponse{}
	return response.Id, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlResourcesIpPools,
		apiInput:    request,
		apiResponse: response,
	})
}

func (o *Client) deleteIp4Pool(ctx context.Context, id ObjectId) error {
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodDelete,
		urlStr:   fmt.Sprintf(apiUrlResourcesIpPoolById, id),
		apiInput: nil,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) updateIp4Pool(ctx context.Context, poolId ObjectId, request *NewIp4PoolRequest) error {
	if request.Subnets == nil {
		request.Subnets = []NewIp4Subnet{}
	}
	if request.Tags == nil {
		request.Tags = []string{}
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlResourcesIpPoolById, poolId),
		apiInput: request,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) addSubnetToIp4Pool(ctx context.Context, poolId ObjectId, new *net.IPNet) error {
	// IPv4 only, buddy
	if strings.Contains(new.String(), ":") {
		return fmt.Errorf("error attmempt to add '%s' to IPv4 address pool", new.String())
	}

	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceIp4PoolRangeMutex)
	defer o.unlock(clientApiResourceIp4PoolRangeMutex)

	// grab the existing pool
	pool, err := o.getIp4Pool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("cannot fetch ip pool - %w", err)
	}

	// check for subnet collisions while copying existing subnets to new request object
	subnets := []NewIp4Subnet{{Network: new.String()}} // start the list with the new one
	for _, s := range pool.Subnets {
		old := s.Network
		if err != nil {
			return fmt.Errorf("error parsing subnet string returned by apstra %s - %w", s.Network, err)
		}
		if old.Contains(new.IP) || new.Contains(old.IP) {
			return fmt.Errorf("new subnet '%s' overlaps existing subnet %s'", new.String(), s.Network)
		}
		subnets = append(subnets, NewIp4Subnet{Network: s.Network.String()})
	}

	err = o.updateIp4Pool(ctx, poolId, &NewIp4PoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
		Subnets:     subnets,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) deleteSubnetFromIp4Pool(ctx context.Context, poolId ObjectId, target *net.IPNet) error {
	// IPv4 only, buddy
	if strings.Contains(target.String(), ":") {
		return fmt.Errorf("error attmempt to add '%s' to IPv4 address pool", target.String())
	}

	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceIp4PoolRangeMutex)
	defer o.unlock(clientApiResourceIp4PoolRangeMutex)

	// grab the existing pool
	pool, err := o.getIp4Pool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("cannot fetch ip pool - %w", err)
	}

	// prep new request structure
	newRequest := &NewIp4PoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
		Subnets:     []NewIp4Subnet{}, // empty slice, but not nil for apstra
	}

	// work through the list copy non-target subnets to the new request
	var targetFound bool
	for _, s := range pool.Subnets {
		if err != nil {
			return fmt.Errorf("error parsing subnet string returned by apstra %s - %w", s.Network, err)
		}

		// copy old subnets which don't match deletion target to new request slice
		if s.Network.String() != target.String() {
			newRequest.Subnets = append(newRequest.Subnets, NewIp4Subnet{Network: s.Network.String()})
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

	err = o.updateIp4Pool(ctx, poolId, newRequest)
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}
