package apstra

import (
	"context"
	"fmt"
)

const (
	apiUrlResources               = "/api/resources"
	apiUrlResourcesAsnPools       = apiUrlResources + "/asn-pools"
	apiUrlResourcesAsnPoolsPrefix = apiUrlResourcesAsnPools + apiUrlPathDelim
	apiUrlResourcesAsnPoolById    = apiUrlResourcesAsnPoolsPrefix + "%s"
	apiUrlResourcesVniPools       = apiUrlResources + "/vni-pools"
	apiUrlResourcesVniPoolsPrefix = apiUrlResourcesVniPools + apiUrlPathDelim
	apiUrlResourcesVniPoolById    = apiUrlResourcesVniPoolsPrefix + "%s"

	apiUrlResourcesIntegerPools       = apiUrlResources + "/integer-pools"
	apiUrlResourcesIntegerPoolsPrefix = apiUrlResourcesIntegerPools + apiUrlPathDelim
	apiUrlResourcesIntegerPoolById    = apiUrlResourcesIntegerPoolsPrefix + "%s"
)

// Following code will take care of ASN Pools

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

// Following code will take care of Integer Pools

func (o *Client) getIntPoolByName(ctx context.Context, desired string) (*rawIntPool, error) {
	pools, err := o.getIntPools(ctx, apiUrlResourcesIntegerPools)
	if err != nil {
		return nil, fmt.Errorf("error fetching all Integer Pools - %w", err)
	}
	found := -1
	for i, pool := range pools {
		if pool.DisplayName == desired {
			if found >= 0 {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("name '%s' does not uniquely identify an Integer Pool", desired),
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

// Following code will take care of VNI Pools

// VniPoolRequest is the public structure used to create/update an ASN pool.
type VniPoolRequest IntPoolRequest

// VniPool is the public structure used to convey query responses about Vni
// pools.
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

func (o *Client) getVniPoolByName(ctx context.Context, desired string) (*VniPool, error) {
	pools, err := o.getVniPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching all VNI pools - %w", err)
	}
	found := -1
	for i, pool := range pools {
		if pool.DisplayName == desired {
			if found >= 0 {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("name '%s' does not uniquely identify an VNI pool", desired),
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
