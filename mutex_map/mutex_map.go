// Copyright (c) Juniper Networks, Inc., 2025-2025.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package mutexmap

import "sync"

type MutexMap struct {
	l *sync.Mutex
	m map[string]*sync.Mutex
}

// Lock locks the sync.Mutex specified by id, and then locks it. The specified
// sync.Mutex will be created if it does not exist.
func (mm *MutexMap) Lock(id string) {
	mm.l.Lock() // lock the map of locks - no defer unlock here, we explicitly unlock in both cases below
	if mu, found := mm.m[id]; found {
		// the sync.Mutex specified by `id` is found to exist
		mm.l.Unlock() // unlock the map because we will not be writing to it
		mu.Lock()     // lock the specified sync.Mutex
	} else {
		// specified sync.Mutex does not exist. Create it, lock it, add it to the map
		mu = new(sync.Mutex)
		mu.Lock()
		mm.m[id] = mu
		mm.l.Unlock()
	}
}

// Unlock releases the sync.Mutex specified by id. It is a run-time error if
// the specified sync.Mutex does not exist or is not locked
func (mm *MutexMap) Unlock(id string) {
	mm.l.Lock()
	defer mm.l.Unlock()

	if mu, ok := mm.m[id]; ok {
		mu.Unlock()
	} else {
		panic("attempt to unlock unknown mutex: '" + id + "'")
	}
}

func NewMutexMap() MutexMap {
	return MutexMap{
		l: new(sync.Mutex),
		m: make(map[string]*sync.Mutex),
	}
}
