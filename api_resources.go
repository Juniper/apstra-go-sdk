package goapstra

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiUrlResources               = "/api/resources"
	apiUrlResourcesAsnPools       = apiUrlResources + "/asn-pools"
	apiUrlResourcesAsnPoolsPrefix = apiUrlResourcesAsnPools + apiUrlPathDelim
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

type optionsAsnPoolsResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

func (o *Client) createAsnPool(ctx context.Context, in *AsnPool) (*objectIdResponse, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPools)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}

	response := &objectIdResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         apstraUrl,
		apiInput:    asnPoolToNewAsnPool(in),
		apiResponse: response,
	})
}

func (o *Client) getAsnPools(ctx context.Context) ([]AsnPool, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPools)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPools, err)
	}
	response := &getAsnPoolsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching ASN pools - %w", err)
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
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(poolId))
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPoolsPrefix+string(poolId), err)
	}
	raw := &rawAsnPool{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: raw,
	})
	if err != nil {
		var ttae TalkToApstraErr
		if errors.As(err, &ttae) && ttae.Response.StatusCode == http.StatusNotFound {
			return nil, ApstraClientErr{
				errType: ErrNotfound,
				err:     err,
			}
		}
		return nil, fmt.Errorf("error fetching ASN pool '%s' - %w", poolId, err)
	}
	return rawAsnPoolToAsnPool(raw)

}

func (o *Client) deleteAsnPool(ctx context.Context, poolId ObjectId) error {
	if poolId == "" {
		return errors.New("attempt to delete ASN Pool with empty pool ID")
	}
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(poolId))
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPoolsPrefix+string(poolId), err)
	}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		url:    apstraUrl,
	})
	if err != nil {
		return fmt.Errorf("error fetching ASN pools - %w", err)
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
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(poolId))
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}

	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
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
	poolInfo, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("error getting ASN pool ranges - %w", err)
	}

	for _, r := range poolInfo.Ranges {
		if asnOverlap(r, *newRange) {
			return ApstraClientErr{
				errType: ErrAsnRangeOverlap,
				err: fmt.Errorf("new ASN range %d-%d overlaps with existing ASN range %d-%d in ASN Pool '%s'",
					newRange.First, newRange.Last, r.First, r.Last, poolId),
			}
		}
	}

	poolInfo.Ranges = append(poolInfo.Ranges, *newRange)
	return o.UpdateAsnPool(ctx, poolId, poolInfo)
}

func (o *Client) deleteAsnPoolRange(ctx context.Context, poolId ObjectId, deleteMe *AsnRange) error {
	poolInfo, err := o.GetAsnPool(ctx, poolId)
	if err != nil {
		return fmt.Errorf("error getting ASN pool ranges - %w", err)
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

	return o.UpdateAsnPool(ctx, poolId, poolInfo)
}

func (o *Client) listAsnPoolIds(ctx context.Context) ([]ObjectId, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPools)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}

	response := &optionsAsnPoolsResponse{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPut,
		url:         apstraUrl,
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}
