// Copyright (c) Juniper Networks, Inc., 2026-2026.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package redundancycache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/Juniper/apstra-go-sdk/internal/pointer"
)

const (
	SystemNotFoundError = "system not found in cache"
	GroupNotFoundError  = "group not found in cache"
	NotFoundError       = "ID not found in cache"
)

func newCache(bp *apstra.TwoStageL3ClosClient) *cache {
	return &cache{
		bp:             bp,
		mu:             new(sync.RWMutex),
		groupToSystems: make(map[string][2]string),
		systemToGroup:  make(map[string]*string),
	}
}

func getCache(bp *apstra.TwoStageL3ClosClient) *cache {
	bpID := bp.Id().String()

	mainMutex.RLock()
	if c, ok := bpToCache[bpID]; ok {
		mainMutex.RUnlock()
		return c
	}
	mainMutex.RUnlock()

	mainMutex.Lock()
	defer mainMutex.Unlock()

	// Another goroutine may have created the cache while we waited.
	if c, ok := bpToCache[bpID]; ok {
		return c
	}

	c := newCache(bp)
	bpToCache[bpID] = c
	return c
}

type cache struct {
	bp             *apstra.TwoStageL3ClosClient
	mu             *sync.RWMutex        // protects both maps below
	groupToSystems map[string][2]string // map of redundancy group ID to system ID pair
	systemToGroup  map[string]*string   // map of system ID to redundancy group ID (nil if not part of a group)
}

// lookupGroup returns the Redundancy Group ID and the peer system ID for the given switch system ID in the given Blueprint, if any.
//
// Possible results:
// - System exists and is part of a redundancy group     : returns non-nil pointers to the RG ID and the peer system ID and nil error
// - System exists and is not part of a redundancy group : returns nil, nil, nil
// - System does not exist, or failure during lookup     : returns nil, nil, error
func (c *cache) lookupGroup(ctx context.Context, systemID string) (*string, *string, error) {
	getPeer := func(groupID *string) *string {
		if groupID == nil {
			return nil
		}

		systemIDs := c.groupToSystems[*groupID]

		var peer string
		if systemID == systemIDs[0] {
			peer = systemIDs[1]
		} else {
			peer = systemIDs[0]
		}

		return &peer
	}

	c.mu.RLock() // lock for read
	if groupID, ok := c.systemToGroup[systemID]; ok {
		defer c.mu.RUnlock()                                         // Release the lock for read on exit.
		return pointer.ToCopyOfValue(groupID), getPeer(groupID), nil // Cache hit - Success!
	}

	// Cache miss - We may need to refresh the cache. Replace our reader lock with a writer lock.
	c.mu.RUnlock()      // Release lock for read.
	c.mu.Lock()         // Acquire a lock for write.
	defer c.mu.Unlock() // Unlock on return.

	// Check the cache one more time, in case another goroutine refreshed the cache while we were waiting for lock.
	if groupID, ok := c.systemToGroup[systemID]; ok {
		return pointer.ToCopyOfValue(groupID), getPeer(groupID), nil // Cache hit - Success!
	}

	// Nope. Refresh the cache.
	err := c.refresh(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("refreshing cache: %w", err)
	}

	// Now that we've refreshed the cache, check it one last time.
	if groupID, ok := c.systemToGroup[systemID]; ok {
		return pointer.ToCopyOfValue(groupID), getPeer(groupID), nil // Cache hit - Success!
	}

	// Probably a bogus system ID.
	return nil, nil, errors.New(SystemNotFoundError + ": " + systemID)
}

// lookupNodeType expects either the ID of a system node with system_type == switch,
// or the ID of a redundancy group node. It returns the type from cache, updating the
// system redundancy cache if required.
func (c *cache) lookupNodeType(ctx context.Context, nodeID string) (apstra.NodeType, error) {
	typeFromCache := func() apstra.NodeType {
		if _, ok := c.groupToSystems[nodeID]; ok {
			return apstra.NodeTypeRedundancyGroup
		}
		if _, ok := c.systemToGroup[nodeID]; ok {
			return apstra.NodeTypeSystem
		}
		return apstra.NodeTypeNone
	}

	c.mu.RLock() // lock for read
	if t := typeFromCache(); t != apstra.NodeTypeNone {
		c.mu.RUnlock()
		return t, nil
	}

	// Try refreshing the cache. Release the read lock before acquiring the write lock.
	c.mu.RUnlock()      // Release the lock for read.
	c.mu.Lock()         // Acquire a lock for write.
	defer c.mu.Unlock() // Release the lock for write on return.

	// Check the cache one more time after acquiring the write lock, in case another thread refreshed it while we were waiting.
	if t := typeFromCache(); t != apstra.NodeTypeNone {
		return t, nil
	}

	// Another cache miss - refresh the cache.
	err := c.refresh(ctx)
	if err != nil {
		return apstra.NodeTypeNone, fmt.Errorf("refreshing cache: %w", err)
	}

	// Now that we've refreshed the cache, check it one last time.
	if t := typeFromCache(); t != apstra.NodeTypeNone {
		return t, nil
	}

	return apstra.NodeTypeNone, errors.New(NotFoundError + ": " + nodeID)
}

// lookupSystems returns a pair of System IDs representing the given redundancy group ID in the given Blueprint.
//
// Possible results:
// - Redundancy Group exists                                 : returns the member system IDs, nil
// - Redundancy Group does not exist or failure during lookup: returns a zero-value array, error
func (c *cache) lookupSystems(ctx context.Context, groupID string) ([2]string, error) {
	c.mu.RLock() // lock for read
	if sysIDs, ok := c.groupToSystems[groupID]; ok {
		c.mu.RUnlock()     // Release the lock for read.
		return sysIDs, nil // Cache hit - Success!
	}

	// Cache miss - We may need to refresh the cache. Replace our reader lock with a writer lock.
	c.mu.RUnlock()      // Release the lock for read.
	c.mu.Lock()         // Acquire a lock for write.
	defer c.mu.Unlock() // Release the lock for write on return.

	// Check the cache one more time, in case another goroutine refreshed the cache while we were waiting for lock.
	if sysIDs, ok := c.groupToSystems[groupID]; ok {
		return sysIDs, nil // Cache hit - Success!
	}

	// Nope. Refresh the cache.
	err := c.refresh(ctx)
	if err != nil {
		return [2]string{}, fmt.Errorf("refreshing cache: %w", err)
	}

	// Now that we've refreshed the cache, check it one last time.
	if sysIDs, ok := c.groupToSystems[groupID]; ok {
		return sysIDs, nil // Cache hit - Success!
	}

	// Probably a bogus group ID.
	return [2]string{}, errors.New(GroupNotFoundError + ": " + groupID)
}

// refresh queries the blueprint for all switches and their redundancy groups,
// if any, and updates the redundancy group membership maps.
// The caller must hold the write lock for the cache when calling this function.
func (c *cache) refresh(ctx context.Context) error {
	systemToGroup, err := getSystemToGroup(ctx, c.bp)
	if err != nil {
		return fmt.Errorf("refreshing system redundancy group cache: %w", err)
	}

	groupToSystems, err := buildGroupToSystemsMap(systemToGroup, c.bp)
	if err != nil {
		return fmt.Errorf("building redundancy group to systems cache: %w", err)
	}

	c.systemToGroup = systemToGroup
	c.groupToSystems = groupToSystems

	return nil
}

// getSystemToGroup returns map[string]*string keyed by system node ID where the values are
// pointers to the redundancy group associated with the system. If the system is not a member
// of a redundancy group, the value will be a nil pointer so that we clearly know the system
// is known and is *not* a member of a group (cached negative result).
func getSystemToGroup(ctx context.Context, bp *apstra.TwoStageL3ClosClient) (map[string]*string, error) {
	query := new(apstra.MatchQuery).
		SetBlueprintId(bp.Id()).
		SetClient(bp.Client()).
		Match(new(apstra.PathQuery).
			Node([]apstra.QEEAttribute{
				apstra.NodeTypeSystem.QEEAttribute(),
				// We filter on system_type='switch' to reduce the cache size. Generic Systems are also "systems"
				// in the DC refdesign graph, but there's lots of them and they're not interesting to us.
				// But not all switches can be part of a redundancy group. DC refdesign has a validation called
				// RG_SUPPORTS_LEAF_ACCESS which ensures that only those with role=is_in(['leaf', 'access']) can
				// be part of a redundancy group, so we *could* disregard spines and superspines as well.
				// We're not doing that because the count of spines and superspines will be low/insignificant
				// and this simplifies the error diagnostic returned by the lookup functions in case of an unknown
				// system ID. Rather than saying "no such switch of type leaf or access with that ID" (we won't
				// know which type we're looking for), we can return "no switch with that ID" because we'll have
				// cache entries for all switches, regardless of their role.
				{Key: "system_type", Value: apstra.QEStringVal(apstra.SystemTypeSwitch.String())},
				{Key: "name", Value: apstra.QEStringVal("n_sys")},
			}),
		).
		Optional(new(apstra.PathQuery).
			Node([]apstra.QEEAttribute{{Key: "name", Value: apstra.QEStringVal("n_sys")}}).
			Out([]apstra.QEEAttribute{apstra.RelationshipTypePartOfRedundancyGroup.QEEAttribute()}).
			Node([]apstra.QEEAttribute{
				apstra.NodeTypeRedundancyGroup.QEEAttribute(),
				{Key: "name", Value: apstra.QEStringVal("n_grp")},
			}),
		)

	// target collects only the system ID and group ID (if any) for each system in the blueprint. The group
	// ID is a pointer to ensure we notice when a system is not part of a redundancy group (nil pointer)
	var target struct {
		Items []struct {
			System struct {
				ID string `json:"id"`
			} `json:"n_sys"`
			Group struct {
				ID *string `json:"id"`
			} `json:"n_grp"`
		} `json:"items"`
	}

	// Run the query.
	err := query.Do(ctx, &target)
	if err != nil {
		return nil, fmt.Errorf("refreshing system redundancy group cache: %w", err)
	}

	// Populate a new system ID -> group ID map.
	result := make(map[string]*string, len(target.Items))
	for _, item := range target.Items {
		result[item.System.ID] = item.Group.ID // nil if system is not part of a redundancy group
	}

	return result, nil
}

// buildGroupToSystemsMap takes a SystemToGroup map and reverses to facilitate lookup by group ID.
func buildGroupToSystemsMap(systemToGroup map[string]*string, bp *apstra.TwoStageL3ClosClient) (map[string][2]string, error) {
	// groupToSystems is a temporary slice that collects system IDs for each redundancy group ID.
	// The value associated with each redundancy group ID is a slice of system IDs.
	groupToSystems := make(map[string][]string)
	for sysID, rgID := range systemToGroup {
		if rgID == nil {
			continue // Skip systems which are not part of a redundancy group.
		}

		// Add the system ID to the slice for this redundancy group.
		groupToSystems[*rgID] = append(groupToSystems[*rgID], sysID)
	}

	// Having created a slice of system IDs for each redundancy group, we can now
	// create the final map of redundancy group ID to system ID pair ([2]string).
	// If a redundancy group does not have exactly 2 members, return an error.
	result := make(map[string][2]string)
	for rgID, sysIDs := range groupToSystems {
		if len(sysIDs) != 2 {
			return nil, fmt.Errorf("%s %q in Blueprint %q does not have exactly 2 members", apstra.NodeTypeRedundancyGroup, rgID, bp.Id())
		}

		// Convert the slice of system IDs to an [2]string system ID pair.
		result[rgID] = [2]string{sysIDs[0], sysIDs[1]} // Direct assignment since we know there are exactly 2 members.
	}

	return result, nil
}
