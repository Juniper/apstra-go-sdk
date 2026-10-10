package redundancycache

func Drop(id string) {
	mainMutex.Lock()
	delete(bpToCache, id)
	mainMutex.Unlock()
}

func CountSystems(bpID string) int {
	mainMutex.RLock()
	c := bpToCache[bpID]
	mainMutex.RUnlock()

	c.mu.RLock()
	result := len(c.systemToGroup)
	c.mu.RUnlock()

	return result
}

func CountGroups(bpID string) int {
	mainMutex.RLock()
	c := bpToCache[bpID]
	mainMutex.RUnlock()

	c.mu.RLock()
	result := len(c.groupToSystems)
	c.mu.RUnlock()

	return result
}
