package goapstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// IntfIntRange allows both IntRangeRequest (sparse type created by the caller)
// and IntRange (detailed type sent by API) to be used in create/update methods
type IntfIntRange interface {
	first() uint32
	last() uint32
}

// IntRanges is used in IntPool responses. It exists as a standalone type to
// facilitate checks with IndexOf() and Overlaps() methods.
type IntRanges []IntRange

// IntRangeRequest is the public structure found within an IntPoolRequest.
type intRangeRequest struct {
	First uint32
	Last  uint32
}

func (o intRangeRequest) first() uint32 {
	return o.First
}

func (o intRangeRequest) last() uint32 {
	return o.Last
}

// rawIntRangeRequest is the API-friendly structure sent within a
// rawIntPoolRequest.
type rawIntRangeRequest struct {
	First uint32 `json:"first"`
	Last  uint32 `json:"last"`
}

// IntPoolRequest is the public structure used to create/update an Int pool.
type IntPoolRequest struct {
	DisplayName string
	Ranges      []IntfIntRange
	Tags        []string
}

// raw() converts an IntPoolRequest to rawIntPoolRequest for consumption by the
// Apstra API.
func (o *IntPoolRequest) raw() *rawIntPoolRequest {
	ranges := make([]rawIntRangeRequest, len(o.Ranges))
	for i, r := range o.Ranges {
		ranges[i] = rawIntRangeRequest{
			First: r.first(),
			Last:  r.last(),
		}
	}
	return &rawIntPoolRequest{
		DisplayName: o.DisplayName,
		Ranges:      ranges,
		Tags:        o.Tags,
	}
}

// rawIntPoolRequest is formatted for the Apstra API, and is used to create or
// update an Int pool.
type rawIntPoolRequest struct {
	DisplayName string               `json:"display_name"`
	Ranges      []rawIntRangeRequest `json:"ranges"`
	Tags        []string             `json:"tags,omitempty"`
}

// IndexOf returns index of 'b'. If not found, it returns -1
func (o IntRanges) IndexOf(b IntfIntRange) int {
	for i, a := range o {
		if a.first() == b.first() && a.last() == b.last() {
			return i
		}
	}
	return -1
}

func (o IntRanges) Overlaps(b IntfIntRange) bool {
	for _, a := range o {
		if IntOverlap(a, b) {
			return true
		}
	}
	return false // no overlap
}

// IntPool is the public structure used to convey query responses about Int
// pools.
type IntPool struct {
	Id             ObjectId
	DisplayName    string
	Ranges         IntRanges // use the named slice type so we can call IndexOf()
	Tags           []string
	Status         string
	CreatedAt      time.Time
	LastModifiedAt time.Time
	Total          uint32
	Used           uint32
	UsedPercentage float32
}

// rawIntPool contains some clunky types (integers as strings, etc.), is
// cleaned up into an IntPool before being presented to callers
type rawIntPool struct {
	Status         string        `json:"status"`
	Used           string        `json:"used"`
	DisplayName    string        `json:"display_name"`
	Tags           []string      `json:"tags"`
	CreatedAt      time.Time     `json:"created_at"`
	LastModifiedAt time.Time     `json:"last_modified_at"`
	Ranges         []rawIntRange `json:"ranges"`
	UsedPercentage float32       `json:"used_percentage"`
	Total          string        `json:"total"`
	Id             ObjectId      `json:"id"`
}

// polish turns a rawIntPool from the API into IntPool for caller consumption
func (o *rawIntPool) polish() (*IntPool, error) {
	ranges := make(IntRanges, len(o.Ranges))
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
			return nil, fmt.Errorf("error parsing 'used' element of Int Pool '%s' - %w", o.Id, err)
		}
	}

	var total uint64
	if o.Total == "" {
		total = 0
	} else {
		total, err = strconv.ParseUint(o.Total, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'total' element of Int Pool '%s' - %w", o.Id, err)
		}
	}
	return &IntPool{
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

// IntRange is the public structure found within IntPool
type IntRange struct {
	Status         string
	First          uint32
	Last           uint32
	Total          uint32
	Used           uint32
	UsedPercentage float32
}

func (o IntRange) first() uint32 {
	return o.First
}

func (o IntRange) last() uint32 {
	return o.Last
}

// rawIntRange contains some clunky types (integers as strings, etc.), is
// cleaned up into an IntRange before being presented to callers
type rawIntRange struct {
	Status         string  `json:"status"`
	First          uint32  `json:"first"`
	Last           uint32  `json:"last"`
	Total          string  `json:"total"`
	Used           string  `json:"used"`
	UsedPercentage float32 `json:"used_percentage"`
}

func (o rawIntRange) first() uint32 {
	//TODO implement me
	return o.First
}

func (o rawIntRange) last() uint32 {
	//TODO implement me
	return o.Last
}

func (o *rawIntRange) polish() (*IntRange, error) {
	var err error
	var used uint64
	if o.Used == "" {
		used = 0
	} else {
		used, err = strconv.ParseUint(o.Used, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'used' element of Int Pool Range - %w", err)
		}
	}

	var total uint64
	if o.Total == "" {
		total = 0
	} else {
		total, err = strconv.ParseUint(o.Total, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'total' element of Int Pool Range - %w", err)
		}
	}

	return &IntRange{
		Status:         o.Status,
		First:          o.First,
		Last:           o.Last,
		Total:          uint32(total),
		Used:           uint32(used),
		UsedPercentage: o.UsedPercentage,
	}, nil
}

func (o *Client) createIntPool(ctx context.Context, in *IntPoolRequest, apiUrlResourcePool string) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlResourcePool, //Will be apiUrlResourcesAsnPool or apiUrlResourcesVniPool
		apiInput:    in.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) listIntPoolIds(ctx context.Context, apiUrlResourcePool string) ([]ObjectId, error) {
	var response struct {
		Items []ObjectId `json:"items"`
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlResourcePool, //Will be apiUrlResourcesAsnPool or apiUrlResourcesVniPool
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getIntPools(ctx context.Context, apiUrlResourcePool string) ([]rawIntPool, error) {
	var response struct {
		Items []rawIntPool `json:"items"`
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlResourcePool,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	return response.Items, nil
}

func (o *Client) getIntPool(ctx context.Context, apiUrlResourcePoolById string, poolId ObjectId) (*rawIntPool, error) {
	if poolId == "" {
		return nil, errors.New("attempt to get Int Pool info with empty pool ID")
	}
	response := &rawIntPool{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlResourcePoolById, poolId), //ApiUrl
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) deleteIntPool(ctx context.Context, apiUrlResourcePoolById string, poolId ObjectId) error {
	if poolId == "" {
		return errors.New("attempt to delete Int Pool with empty pool ID")
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlResourcePoolById, poolId),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) updateIntPool(ctx context.Context, apiUrlResourcePoolById string, poolId ObjectId, pool *IntPoolRequest) error {
	if poolId == "" {
		return errors.New("attempt to update Int Pool with empty pool ID")
	}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlResourcePoolById, poolId),
		apiInput: pool.raw(),
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}
	return nil
}

func (o *Client) createIntPoolRange(ctx context.Context, apiResourcesPoolById string, poolId ObjectId, newRange *intRangeRequest) error {
	// we read, then replace the pool range. this is not concurrency safe.
	// read the Int pool info (that's where the configured ranges are found)
	p, err := o.getIntPool(ctx, apiResourcesPoolById, poolId)
	if err != nil {
		return fmt.Errorf("error getting Int pool ranges - %w", err)
	}
	pool, err := p.polish()
	if err != nil {
		return fmt.Errorf("error getting Int pool ranges - %w", err)
	}
	// we don't expect to find the "new" range in there already
	if pool.Ranges.IndexOf(newRange) >= 0 {
		return ApstraClientErr{
			errType: ErrExists,
			err:     fmt.Errorf(" range %d-%d in  pool '%s' already exists, cannot create", newRange.First, newRange.Last, pool.Id),
		}
	}

	// sanity check: the new range shouldn't overlap any existing range (the API will reject it)
	if pool.Ranges.Overlaps(newRange) {
		return ApstraClientErr{
			errType: ErrRangeOverlap,
			err: fmt.Errorf("new range %d-%d overlaps with existing range in  Pool '%s'",
				newRange.First, newRange.Last, poolId),
		}
	}

	req := &IntPoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
	}

	// make one extra slice element for the new range element
	req.Ranges = make([]IntfIntRange, len(pool.Ranges)+1)

	// fill the first elements with the retrieved data
	for i, r := range pool.Ranges {
		req.Ranges[i] = r
	}

	// populate the final element (index matches length of retrieved data) with
	// the new range element
	req.Ranges[len(pool.Ranges)] = *newRange

	return o.updateIntPool(ctx, apiResourcesPoolById, poolId, req)
}

func (o *Client) IntPoolRangeExists(ctx context.Context, apiResourcesPoolById string, poolId ObjectId, IntRange IntfIntRange) (bool, error) {
	poolInfo, err := o.getIntPool(ctx, apiResourcesPoolById, poolId)
	if err != nil {
		return false, fmt.Errorf("error getting Int ranges from pool '%s' - %w", poolId, err)
	}
	p, err := poolInfo.polish()
	if err != nil {
		return false, fmt.Errorf("error polishing Int ranges from pool '%s' - %w", poolId, err)
	}
	if p.Ranges.IndexOf(IntRange) >= 0 {
		return true, nil
	}
	return false, nil
}

func (o *Client) deleteIntPoolRange(ctx context.Context, apiResourcesPoolById string, poolId ObjectId, deleteMe IntfIntRange) error {
	// we read, then replace the pool range. this is not concurrency safe.
	// Caller must take a lock

	pool, err := o.getIntPool(ctx, apiResourcesPoolById, poolId)
	if err != nil {
		return fmt.Errorf("error getting ranges from pool '%s' - %w", poolId, err)
	}
	p, err := pool.polish()
	if err != nil {
		return fmt.Errorf("error getting ranges from pool '%s' - %w", poolId, err)
	}
	deleteIdx := p.Ranges.IndexOf(deleteMe)
	if deleteIdx < 0 {
		return ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("range '%d-%d' not found in Int Pool '%s'", deleteMe.first(), deleteMe.last(), poolId),
		}
	}

	req := &IntPoolRequest{
		DisplayName: pool.DisplayName,
		Tags:        pool.Tags,
	}
	for i := range pool.Ranges {
		if i == deleteIdx {
			continue
		}
		req.Ranges = append(req.Ranges, &pool.Ranges[i])
	}
	return o.updateIntPool(ctx, apiResourcesPoolById, poolId, req)
}
func IntOverlap(a, b IntfIntRange) bool {
	if a.first() >= b.first() && a.first() <= b.last() { // begin 'a' falls within 'b'
		return true
	}
	if a.last() <= b.last() && a.last() >= b.first() { // end 'a' falls within 'b'
		return true
	}
	if b.first() >= a.first() && b.first() <= a.last() { // begin 'b' falls within 'a'
		return true
	}
	if b.last() <= a.last() && b.last() >= a.first() { // end 'b' falls within 'a'
		return true
	}
	return false // no overlap
}
