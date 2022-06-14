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

	asnPoolRangeIdPrefix = "pool-"
)

type NewAsnRange struct {
	B uint32 `json:"first"`
	E uint32 `json:"last"`
}

func (o NewAsnRange) String() string {
	return fmt.Sprintf("%d-%d", o.B, o.E)
}

type NewAsnPoolCfg struct {
	DisplayName string        `json:"display_name"`
	Ranges      []NewAsnRange `json:"ranges"`
	Tags        []string      `json:"tags"`
}

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

type rawAsnRange struct {
	Status         string  `json:"status"`
	First          uint32  `json:"first"`
	Last           uint32  `json:"last"`
	Total          string  `json:"total"`
	Used           string  `json:"used"`
	UsedPercentage float32 `json:"used_percentage"`
}

type AsnPool struct {
	Status         string     `json:"status"`
	Used           uint32     `json:"used"`
	DisplayName    string     `json:"display_name"`
	Tags           []string   `json:"tags"`
	CreatedAt      time.Time  `json:"created_at"`
	LastModifiedAt time.Time  `json:"last_modified_at"`
	Ranges         []AsnRange `json:"ranges"`
	UsedPercentage float32    `json:"used_percentage"`
	Total          uint32     `json:"total"`
	Id             ObjectId   `json:"id"`
}

type AsnRange struct {
	Status         string  `json:"status"`
	First          uint32  `json:"first"`
	Last           uint32  `json:"last"`
	Total          uint32  `json:"total"`
	Used           uint32  `json:"used"`
	UsedPercentage float32 `json:"used_percentage"`
}

type getAsnPoolsResponse struct {
	Items []rawAsnPool `json:"items"`
}

func (o *Client) createAsnPool(ctx context.Context, in *NewAsnPoolCfg) (*objectIdResponse, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPools)
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}
	if in.Ranges == nil {
		in.Ranges = []NewAsnRange{} // avoid '"ranges": null' in JSON output
	}
	if in.Tags == nil {
		in.Tags = []string{} // avoid '"tags": null' in JSON output
	}
	response := &objectIdResponse{}
	return response, o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		url:         apstraUrl,
		apiInput:    in,
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

func (o *Client) getAsnPool(ctx context.Context, in ObjectId) (*AsnPool, error) {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(in))
	if err != nil {
		return nil, fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPoolsPrefix+string(in), err)
	}
	raw := &rawAsnPool{}
	err = o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: raw,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching ASN pool '%s' - %w", in, err)
	}
	return rawAsnPoolToAsnPool(raw)

}

func (o *Client) deleteAsnPool(ctx context.Context, in ObjectId) error {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(in))
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlResourcesAsnPoolsPrefix+string(in), err)
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

func rawAsnPoolToAsnPool(in *rawAsnPool) (*AsnPool, error) {
	used, err := strconv.ParseUint(in.Used, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'used' element of ASN Pool '%s' - %w", in.Id, err)
	}

	total, err := strconv.ParseUint(in.Total, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'total' element of ASN Pool '%s' - %w", in.Id, err)
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

func (o *Client) updateAsnPool(ctx context.Context, poolId ObjectId, poolInfo *NewAsnPoolCfg) error {
	apstraUrl, err := url.Parse(apiUrlResourcesAsnPoolsPrefix + string(poolId))
	if err != nil {
		return fmt.Errorf("error parsing url '%s' - %w", apiUrlVersion, err)
	}
	if poolInfo.Ranges == nil {
		poolInfo.Ranges = []NewAsnRange{}
	}
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		url:      apstraUrl,
		apiInput: poolInfo,
	})
}

func getAsnPoolRangeId(in *AsnRange) string {
	first := make([]byte, 4)
	last := make([]byte, 4)

	binary.BigEndian.PutUint32(first, in.First)
	binary.BigEndian.PutUint32(last, in.Last)

	hash := sha256.Sum256(append(first, last...))
	printable := hex.EncodeToString(hash[0:len(hash)])
	return asnPoolRangeIdPrefix + printable[:8]
}

func (o *Client) getAsnPoolRanges(ctx context.Context, poolId ObjectId) (map[string]AsnRange, error) {
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
		rid := getAsnPoolRangeId(&r)
		result[rid] = r
	}

	return result, nil
}
