package goapstra

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
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

type AsnPool struct {
	Id             ObjectId   `json:"id"`
	DisplayName    string     `json:"display_name"`
	Ranges         []AsnRange `json:"ranges"`
	Tags           []string   `json:"tags"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at"`
	Total          uint32     `json:"total"`
	Used           uint32     `json:"used"`
	UsedPercentage float32    `json:"used_percentage"`
}

type AsnRange struct {
	Status         string  `json:"status"`
	First          uint32  `json:"first"`
	Last           uint32  `json:"last"`
	Total          uint32  `json:"total"`
	Used           uint32  `json:"used"`
	UsedPercentage float32 `json:"used_percentage"`
}

// newAsnPool is minimal version of AsnPool which omits statistical elements.
// It is used with create/update commands upstream towards Apstra.
type newAsnPool struct {
	DisplayName string        `json:"display_name"`
	Ranges      []newAsnRange `json:"ranges"`
	Tags        []string      `json:"tags"`
}

// newAsnRange is minimal version of AsnRange which omits statisitcal elements.
// It is used with create/update commands upstream towards Apstra.
type newAsnRange struct {
	First uint32 `json:"first"`
	Last  uint32 `json:"last"`
}

// rawAsnPool is contains some clunky types (integers as strings, etc), is
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

// rawAsnRange is contains some clunky types (integers as strings, etc), is
// cleaned up into an AsnPool before being presented to callers
type rawAsnRange struct {
	Status         string  `json:"status"`
	First          uint32  `json:"first"`
	Last           uint32  `json:"last"`
	Total          string  `json:"total"`
	Used           string  `json:"used"`
	UsedPercentage float32 `json:"used_percentage"`
}

type getAsnPoolsResponse struct {
	Items []rawAsnPool `json:"items"`
}

type getIp4PoolsResponse struct {
	Items []rawIp4Pool `json:"items"`
}

type optionsResourcePoolResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

func (o *Client) createAsnPool(ctx context.Context, in *AsnPool) (*objectIdResponse, error) {
	method := http.MethodPost
	urlStr := apiUrlResourcesAsnPools
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	response := &objectIdResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiInput:    asnPoolToNewAsnPool(in),
		apiResponse: response,
	})
}

func (o *Client) listAsnPoolIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsResourcePoolResponse{}
	return response.Items, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlResourcesAsnPools,
		apiResponse: response,
	})
}

func (o *Client) getAsnPools(ctx context.Context) ([]AsnPool, error) {
	method := http.MethodGet
	urlStr := apiUrlResourcesAsnPools
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	response := &getAsnPoolsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s' - %w", method, urlStr, err)
	}

	var pools []AsnPool
	for _, rawPool := range response.Items {
		p, err := rawAsnPoolToAsnPool(&rawPool)
		if err != nil {
			return nil, err
		}
		pools = append(pools, *p)
	}
	return pools, nil
}

func (o *Client) getAsnPool(ctx context.Context, poolId ObjectId) (*AsnPool, error) {
	if poolId == "" {
		return nil, errors.New("attempt to get ASN Pool info with empty pool ID")
	}
	method := http.MethodGet
	urlStr := fmt.Sprintf(apiUrlResourcesAsnPoolById, poolId)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	raw := &rawAsnPool{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      method,
		url:         apstraUrl,
		apiResponse: raw,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling '%s' at '%s' - %w", method, urlStr, convertTtaeToAceWherePossible(err))
	}
	return rawAsnPoolToAsnPool(raw)

}

func (o *Client) deleteAsnPool(ctx context.Context, poolId ObjectId) error {
	if poolId == "" {
		return errors.New("attempt to delete ASN Pool with empty pool ID")
	}
	method := http.MethodDelete
	urlStr := fmt.Sprintf(apiUrlResourcesAsnPoolById, poolId)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		url:    apstraUrl,
	})
	if err != nil {
		return fmt.Errorf("error calling '%s' at '%s - %w", method, urlStr, convertTtaeToAceWherePossible(err))
	}
	return nil
}

// asnPoolToNewAsnPool copies relevant fields from an AsnPool to a newAsnPool.
// The latter type does not include statistical info, is suitable for sending
// to Apstra in create/update methods.
func asnPoolToNewAsnPool(in *AsnPool) *newAsnPool {
	var nars []newAsnRange
	for _, r := range in.Ranges {
		nars = append(nars, newAsnRange{First: r.First, Last: r.Last})
	}

	// Apstra wants '"ranges": []' rather than '"ranges": null'
	if nars == nil {
		nars = []newAsnRange{}
	}

	// Apstra wants '"tags": []' rather than '"tags": null'
	if in.Tags == nil {
		in.Tags = []string{}
	}

	return &newAsnPool{
		DisplayName: in.DisplayName,
		Ranges:      nars,
		Tags:        in.Tags,
	}
}

// rawAsnPoolToAsnPool cleans up a rawAsnPool object (ints as strings) into
// and AsnPool.
func rawAsnPoolToAsnPool(in *rawAsnPool) (*AsnPool, error) {
	var total, used uint64
	var err error

	if in.Used == "" {
		used = 0
	} else {
		used, err = strconv.ParseUint(in.Used, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'used' element of ASN Pool '%s' - %w", in.Id, err)
		}
	}

	if in.Total == "" {
		total = 0
	} else {
		total, err = strconv.ParseUint(in.Total, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'total' element of ASN Pool '%s' - %w", in.Id, err)
		}
	}

	result := AsnPool{
		Status:         in.Status,
		Used:           uint32(used),
		Total:          uint32(total),
		DisplayName:    in.DisplayName,
		Tags:           in.Tags,
		CreatedAt:      in.CreatedAt,
		LastModifiedAt: in.LastModifiedAt,
		UsedPercentage: in.UsedPercentage,
		Id:             in.Id,
	}

	for i, r := range in.Ranges {
		used, err := strconv.ParseUint(in.Used, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing ASN Pool '%s', 'ranges[%d]', 'used' element - %w", in.Id, i, err)
		}

		total, err := strconv.ParseUint(in.Total, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing ASN Pool '%s', 'ranges[%d]', 'total' element - %w", in.Id, i, err)
		}

		result.Ranges = append(result.Ranges, AsnRange{
			Status:         r.Status,
			First:          r.First,
			Last:           r.Last,
			Total:          uint32(total),
			Used:           uint32(used),
			UsedPercentage: r.UsedPercentage,
		})

	}

	return &result, nil
}

func (o *Client) updateAsnPool(ctx context.Context, poolId ObjectId, poolInfo *AsnPool) error {
	if poolId == "" {
		return errors.New("attempt to update ASN Pool with empty pool ID")
	}
	method := http.MethodPut
	urlStr := fmt.Sprintf(apiUrlResourcesAsnPoolById, poolId)
	apstraUrl, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", urlStr, err)
	}

	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   method,
		url:      apstraUrl,
		apiInput: asnPoolToNewAsnPool(poolInfo),
	})
}

func hashAsnPoolRange(in *AsnRange) string {
	first := make([]byte, 4)
	last := make([]byte, 4)

	binary.BigEndian.PutUint32(first, in.First)
	binary.BigEndian.PutUint32(last, in.Last)

	hash := sha256.Sum256(append(first, last...))
	printable := hex.EncodeToString(hash[0:len(hash)])
	return printable
}

func (o *Client) hashAsnPoolRanges(ctx context.Context, poolId ObjectId) (map[string]AsnRange, error) {
	result := make(map[string]AsnRange)
	pool, err := o.getAsnPool(ctx, poolId)
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) {
			if ttae.Response.StatusCode == http.StatusNotFound {
				return nil, ApstraClientErr{
					errType: ErrNotfound,
					err:     err,
				}
			}
		} else {
			return nil, fmt.Errorf("error getting ASN pool info for pool '%s' - %w", poolId, err)
		}
	}

	for _, r := range pool.Ranges {
		rid := hashAsnPoolRange(&r)
		result[rid] = r
	}

	return result, nil
}

func (o *Client) createAsnPoolRange(ctx context.Context, poolId ObjectId, newRange *AsnRange) error {
	if newRange.First <= 0 || newRange.Last > math.MaxUint32 {
		return ApstraClientErr{
			errType: ErrAsnOutOfRange,
			err:     fmt.Errorf("error invalid ASN Range %d-%d", newRange.First, newRange.Last),
		}
	}

	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)

	// read the ASN pool info (that's where the configured ranges are found)
	poolInfo, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("error getting ASN pool ranges - %w", err)
	}

	// we don't expect to find the "new" range in there already
	if AsnPoolRangeInSlice(newRange, poolInfo.Ranges) {
		return ApstraClientErr{
			errType: ErrExists,
			err:     fmt.Errorf("ASN range %d-%d in ASN pool '%s' already exists, cannot create", newRange.First, newRange.Last, poolInfo.Id),
		}
	}

	// sanity check: the new range shouldn't overlap any existing range (the API will reject it)
	for _, r := range poolInfo.Ranges {
		if AsnOverlap(r, *newRange) {
			return ApstraClientErr{
				errType: ErrAsnRangeOverlap,
				err: fmt.Errorf("new ASN range %d-%d overlaps with existing ASN range %d-%d in ASN Pool '%s'",
					newRange.First, newRange.Last, r.First, r.Last, poolId),
			}
		}
	}

	poolInfo.Ranges = append(poolInfo.Ranges, *newRange)
	return o.updateAsnPool(ctx, poolId, poolInfo)
}

func (o *Client) asnPoolRangeExists(ctx context.Context, poolId ObjectId, asnRange *AsnRange) (bool, error) {
	poolInfo, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return false, fmt.Errorf("error getting ASN ranges from pool '%s' - %w", poolId, err)
	}

	return AsnPoolRangeInSlice(asnRange, poolInfo.Ranges), nil
}

func (o *Client) deleteAsnPoolRange(ctx context.Context, poolId ObjectId, deleteMe *AsnRange) error {
	// we read, then replace the pool range. this is not concurrency safe.
	o.lock(clientApiResourceAsnPoolRangeMutex)
	defer o.unlock(clientApiResourceAsnPoolRangeMutex)

	poolInfo, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("error getting ASN ranges from pool '%s' - %w", poolId, err)
	}

	initialRangeCount := len(poolInfo.Ranges)

	targetHash := hashAsnPoolRange(deleteMe)
	for i := initialRangeCount - 1; i >= 0; i-- {
		if targetHash == hashAsnPoolRange(&poolInfo.Ranges[i]) {
			poolInfo.Ranges = append(poolInfo.Ranges[:i], poolInfo.Ranges[i+1:]...)
		}
	}

	if initialRangeCount == len(poolInfo.Ranges) {
		return ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("ASN range '%d-%d' not found in ASN Pool '%s'", deleteMe.First, deleteMe.Last, poolId),
		}
	}

	return o.updateAsnPool(ctx, poolId, poolInfo)
}

// ip4 pool stuff below here

// minimal create subnet request:
// {
//  "subnets": [], // mandatory
//  "tags": ["t1", "t2", "t3"],
//  "display_name": "bang"
//}

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

func (o *Ip4Pool) ToNew() *NewIp4PoolRequest {
	var subnets []NewIp4Subnet
	for _, s := range o.Subnets {
		subnets = append(subnets, *s.ToNew())
	}
	return &NewIp4PoolRequest{
		DisplayName: o.DisplayName,
		Tags:        o.Tags,
		Subnets:     subnets,
	}
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

func (o *Ip4Subnet) ToNew() *NewIp4Subnet {
	return &NewIp4Subnet{Network: o.Network.String()}
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
	response := &optionsResourcePoolResponse{}
	return response.Items, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlResourcesIpPools,
		apiResponse: response,
	})
}

func (o *Client) getIp4Pools(ctx context.Context) ([]Ip4Pool, error) {
	response := &getIp4PoolsResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlResourcesIpPools,
		apiResponse: response,
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

	polishedPool, err := response.polish()
	if err != nil {
		return nil, fmt.Errorf("error parsing raw pool content - %w", err)
	}

	return polishedPool, nil
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
			if found == true {
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
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodDelete,
		urlStr:   fmt.Sprintf(apiUrlResourcesIpPoolById, id),
		apiInput: nil,
	})
}

func (o *Client) updateIp4Pool(ctx context.Context, poolId ObjectId, request *NewIp4PoolRequest) error {
	if request.Subnets == nil {
		request.Subnets = []NewIp4Subnet{}
	}
	if request.Tags == nil {
		request.Tags = []string{}
	}
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlResourcesIpPoolById, poolId),
		apiInput: request,
	})
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
	return err
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

	return o.updateIp4Pool(ctx, poolId, newRequest)
}
