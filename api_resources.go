package goapstra

import (
	"context"
	"errors"
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
	apiUrlResourcesIpPools        = apiUrlResources + "/ip-pools"
	apiUrlResourcesIpPoolsPrefix  = apiUrlResourcesIpPools + apiUrlPathDelim
	apiUrlResourcesIpPoolById     = apiUrlResourcesIpPoolsPrefix + "%s"
)

// AsnRangeRequest is the public structure found within an AsnPoolRequest.
type AsnRangeRequest struct {
	First uint32
	Last  uint32
}

func (o AsnRangeRequest) first() uint32 {
	return o.First
}

func (o AsnRangeRequest) last() uint32 {
	return o.Last
}

type IntfAsnRange interface {
	first() uint32
	last() uint32
}

// rawAsnRangeRequest is the API-friendly structure sent within a
// rawAsnPoolRequest.
type rawAsnRangeRequest struct {
	First uint32 `json:"first"`
	Last  uint32 `json:"last"`
}

// AsnPoolRequest is the public structure used to create/update an ASN pool.
type AsnPoolRequest struct {
	DisplayName string
	Ranges      []IntfAsnRange
	Tags        []string
}

// raw() converts an AsnPoolRequest to rawAsnPoolRequest for consumption by the
// Apstra API.
func (o *AsnPoolRequest) raw() *rawAsnPoolRequest {
	ranges := make([]rawAsnRangeRequest, len(o.Ranges))
	for i, r := range o.Ranges {
		ranges[i] = rawAsnRangeRequest{
			First: r.first(),
			Last:  r.last(),
		}
	}
	return &rawAsnPoolRequest{
		DisplayName: o.DisplayName,
		Ranges:      ranges,
		Tags:        o.Tags,
	}
}

// rawAsnPoolRequest is formatted for the Apstra API, and is used to create or
// update an ASN pool.
type rawAsnPoolRequest struct {
	DisplayName string               `json:"display_name"`
	Ranges      []rawAsnRangeRequest `json:"ranges"`
	Tags        []string             `json:"tags,omitempty"'`
}

type AsnRanges []AsnRange

// indexOf returns index of 'b'. If not found, it returns -1
func (o AsnRanges) indexOf(b IntfAsnRange) int {
	for i, a := range o {
		if a.first() == b.first() && a.last() == b.last() {
			return i
		}
	}
	return -1
}

func (o AsnRanges) overlaps(b IntfAsnRange) bool {
	for _, a := range o {
		if AsnOverlap(a, b) {
			return true
		}
	}
	return false // no overlap
}

// AsnPool is the public structure used to convey query responses about ASN
// pools.
type AsnPool struct {
	Id             ObjectId
	DisplayName    string
	Ranges         AsnRanges
	Tags           []string
	Status         string
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Total          uint32
	Used           uint32
	UsedPercentage float32
}

// rawAsnPool contains some clunky types (integers as strings, etc.), is
// cleaned up into an AsnPool before being presented to callers
type rawAsnPool struct {
	Status         string        `json:"status"`
	Used           string        `json:"used"`
	DisplayName    string        `json:"display_name"`
	Tags           []string      `json:"tags"`
	CreatedAt      time.Time     `json:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at"`
	Ranges         []rawAsnRange `json:"ranges"`
	UsedPercentage float32       `json:"used_percentage"`
	Total          string        `json:"total"`
	Id             ObjectId      `json:"id"`
}

// polish turns a rawAsnPool from the API into AsnPool for caller consumption
func (o *rawAsnPool) polish() (*AsnPool, error) {
	ranges := make([]AsnRange, len(o.Ranges))
	for i, r := range o.Ranges {
		p, err := r.polish()
		if err != nil {
			return nil, err
		}
		ranges[i] = *p
	}

	var err error
	var used uint64
	if o.Used == "" {
		used = 0
	} else {
		used, err = strconv.ParseUint(o.Used, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'used' element of ASN Pool '%s' - %w", o.Id, err)
		}
	}

	var total uint64
	if o.Total == "" {
		total = 0
	} else {
		total, err = strconv.ParseUint(o.Total, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'total' element of ASN Pool '%s' - %w", o.Id, err)
		}
	}
	return &AsnPool{
		Id:             o.Id,
		DisplayName:    o.DisplayName,
		Ranges:         ranges,
		Tags:           o.Tags,
		Status:         o.Status,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Total:          uint32(total),
		Used:           uint32(used),
		UsedPercentage: o.UsedPercentage,
	}, nil
}

// AsnRange is the public structure found within AsnPool
type AsnRange struct {
	Status         string
	First          uint32
	Last           uint32
	Total          uint32
	Used           uint32
	UsedPercentage float32
}

func (o AsnRange) first() uint32 {
	return o.First
}

func (o AsnRange) last() uint32 {
	return o.Last
}

// rawAsnRange contains some clunky types (integers as strings, etc.), is
// cleaned up into an AsnRange before being presented to callers
type rawAsnRange struct {
	Status         string  `json:"status"`
	First          uint32  `json:"first"`
	Last           uint32  `json:"last"`
	Total          string  `json:"total"`
	Used           string  `json:"used"`
	UsedPercentage float32 `json:"used_percentage"`
}

func (o *rawAsnRange) polish() (*AsnRange, error) {
	var err error
	var used uint64
	if o.Used == "" {
		used = 0
	} else {
		used, err = strconv.ParseUint(o.Used, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'used' element of ASN Pool Range - %w", err)
		}
	}

	var total uint64
	if o.Total == "" {
		total = 0
	} else {
		total, err = strconv.ParseUint(o.Total, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'total' element of ASN Pool Range - %w", err)
		}
	}

	return &AsnRange{
		Status:         o.Status,
		First:          o.First,
		Last:           o.Last,
		Total:          uint32(total),
		Used:           uint32(used),
		UsedPercentage: o.UsedPercentage,
	}, nil
}

func (o *Client) createAsnPool(ctx context.Context, in *AsnPoolRequest) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlResourcesAsnPools,
		apiInput:    in.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) listAsnPoolIds(ctx context.Context) ([]ObjectId, error) {
	var response struct {
		Items []ObjectId `json:"items"`
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlResourcesAsnPools,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getAsnPools(ctx context.Context) ([]rawAsnPool, error) {
	var response struct {
		Items []rawAsnPool `json:"items"`
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlResourcesAsnPools,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *Client) getAsnPool(ctx context.Context, poolId ObjectId) (*rawAsnPool, error) {
	if poolId == "" {
		return nil, errors.New("attempt to get ASN Pool info with empty pool ID")
	}
	response := &rawAsnPool{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlResourcesAsnPoolById, poolId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) deleteAsnPool(ctx context.Context, poolId ObjectId) error {
	if poolId == "" {
		return errors.New("attempt to delete ASN Pool with empty pool ID")
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlResourcesAsnPoolById, poolId),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) updateAsnPool(ctx context.Context, poolId ObjectId, pool *AsnPoolRequest) error {
	if poolId == "" {
		return errors.New("attempt to update ASN Pool with empty pool ID")
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlResourcesAsnPoolById, poolId),
		apiInput: pool.raw(),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) createAsnPoolRange(ctx context.Context, poolId ObjectId, newRange *AsnRangeRequest) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)

	// read the ASN pool info (that's where the configured ranges are found)
	pool, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("error getting ASN pool ranges - %w", err)
	}

	// we don't expect to find the "new" range in there already
	if pool.Ranges.indexOf(newRange) >= 0 {
		return ApstraClientErr{
			errType: ErrExists,
			err:     fmt.Errorf("ASN range %d-%d in ASN pool '%s' already exists, cannot create", newRange.First, newRange.Last, pool.Id),
		}
	}

	// sanity check: the new range shouldn't overlap any existing range (the API will reject it)
	if pool.Ranges.overlaps(newRange) {
		return ApstraClientErr{
			errType: ErrAsnRangeOverlap,
			err: fmt.Errorf("new ASN range %d-%d overlaps with existing range in ASN Pool '%s'",
				newRange.First, newRange.Last, poolId),
		}
	}

	req := &AsnPoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
	}

	// make one extra slice element for the new range element
	req.Ranges = make([]IntfAsnRange, len(pool.Ranges)+1)

	// fill the first elements with the retrieved data
	for i, r := range pool.Ranges {
		req.Ranges[i] = AsnRangeRequest{
			First: r.First,
			Last:  r.Last,
		}
	}

	// populate the final element (index matches length of retrieved data) with
	// the new range element
	req.Ranges[len(pool.Ranges)] = *newRange

	return o.updateAsnPool(ctx, poolId, req)
}

func (o *Client) asnPoolRangeExists(ctx context.Context, poolId ObjectId, asnRange IntfAsnRange) (bool, error) {
	poolInfo, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return false, fmt.Errorf("error getting ASN ranges from pool '%s' - %w", poolId, err)
	}

	if poolInfo.Ranges.indexOf(asnRange) >= 0 {
		return true, nil
	}
	return false, nil
}

func (o *Client) deleteAsnPoolRange(ctx context.Context, poolId ObjectId, deleteMe *AsnRangeRequest) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)

	pool, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("error getting ASN ranges from pool '%s' - %w", poolId, err)
	}

	initialRangeCount := len(pool.Ranges)

	indexOf := pool.Ranges.indexOf(deleteMe)
	if indexOf < 0 {
		if initialRangeCount == len(pool.Ranges) {
			return ApstraClientErr{
				errType: ErrNotfound,
				err:     fmt.Errorf("ASN range '%d-%d' not found in ASN Pool '%s'", deleteMe.First, deleteMe.Last, poolId),
			}
		}
	}

	pool.Ranges = append(pool.Ranges[:indexOf], pool.Ranges[indexOf+1:]...)

	req := &AsnPoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
	}
	req.Ranges = make([]IntfAsnRange, len(pool.Ranges))
	for i, r := range pool.Ranges {
		req.Ranges[i] = AsnRangeRequest{
			First: r.First,
			Last:  r.Last,
		}
	}

	return o.updateAsnPool(ctx, poolId, req)
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

	total, err := strconv.ParseInt(o.Used, 10, 64)
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

func (o *Client) getIp4PoolByName(ctx context.Context, desiredName string) (*Ip4Pool, error) {
	pools, err := o.getIp4Pools(ctx)
	if err != nil {
		return nil, err
	}

	var pool Ip4Pool
	var found bool

	for _, p := range pools {
		if p.DisplayName == desiredName {
			if found {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("multiple matches for IP Pool with name '%s'", desiredName),
				}
			}
			pool = p
			found = true
		}
	}
	return &pool, nil
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
