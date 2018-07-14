package main

import (
	"hash/fnv"
	"sort"
	"sync"
)

type StripedLock interface {
	Get(id string) sync.Locker
	Lock(id string)
	Unlock(id string)
	BatchLock(ids []string)
}

type stripedLock struct {
	size  uint32
	locks []sync.Mutex
}

// NewStripedLock creates an array of locks that are striped by the hash code
// of an id string.
//
// !!! Warning: never initialize with size 0.
func NewStripedLock(size uint32) StripedLock {
	sl := &stripedLock{
		size:  size,
		locks: make([]sync.Mutex, size),
	}
	return sl
}

// Get returns the lock associated with a given id based on the id's hash code
func (sl *stripedLock) Get(id string) sync.Locker {
	return &sl.locks[sl.idToIndex(id)]
}

// Lock acquires the lock associated with a given id based on the id's hash code
func (sl *stripedLock) Lock(id string) {
	sl.locks[sl.idToIndex(id)].Lock()
}

// Unlock releases the lock associated with a given id based on the id's hash code
func (sl *stripedLock) Unlock(id string) {
	sl.locks[sl.idToIndex(id)].Unlock()
}

func (sl *stripedLock) BatchLock(ids []string) {
	sort.Strings(ids)
	for _, id := range ids {
		sl.Lock(id)
	}
}

// idToIndex is an id hashing function
func (sl *stripedLock) idToIndex(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32() % sl.size
}
