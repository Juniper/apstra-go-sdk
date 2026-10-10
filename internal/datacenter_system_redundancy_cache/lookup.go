package redundancycache

import (
	"context"
	"sync"

	"github.com/Juniper/apstra-go-sdk/apstra"
)

var (
	mainMutex = new(sync.RWMutex)
	bpToCache = make(map[string]*cache)
)

// LookupGroup returns the Redundancy Group ID and the peer system ID for the given switch system ID in the given Blueprint, if any.
//
// Possible results:
// - System exists and is part of a redundancy group     : returns non-nil pointers to the RG ID and the peer system ID and nil error
// - System exists and is not part of a redundancy group : returns nil, nil, nil
// - System does not exist, or failure during lookup     : returns nil, nil, error
func LookupGroup(ctx context.Context, bp *apstra.TwoStageL3ClosClient, systemID string) (*string, *string, error) {
	return getCache(bp).lookupGroup(ctx, systemID)
}

// LookupSystems returns an unordered pair of System IDs representing members of the given redundancy group ID in the given Blueprint.
//
// Possible results:
// - Redundancy Group exists                                 : returns the member system IDs, nil
// - Redundancy Group does not exist or failure during lookup: returns a zero-value array, error
func LookupSystems(ctx context.Context, bp *apstra.TwoStageL3ClosClient, groupID string) ([2]string, error) {
	return getCache(bp).lookupSystems(ctx, groupID)
}

// LookupNodeType expects either the ID of a system node with system_type == switch,
// or the ID of a redundancy group node. It returns the type from cache, updating the
// system redundancy cache if required. If the node ID is not found, it returns an error.
func LookupNodeType(ctx context.Context, bp *apstra.TwoStageL3ClosClient, nodeID string) (apstra.NodeType, error) {
	return getCache(bp).lookupNodeType(ctx, nodeID)
}
