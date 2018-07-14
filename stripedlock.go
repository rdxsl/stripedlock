package main

import (
	"fmt"
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

// BatchLock will try to acquire locks on a batch of ids.
// The ids are deduped to avoid hash collisions within the set and the locks are acquired based on
// the sorted order of ids, so that concurrent batch updates cannot deadlock.
func (sl *stripedLock) BatchLock(ids []string) {
	hashcodes := sl.getHashcodes(ids)
	for _, hashcode := range hashcodes {
		sl.locks[hashcode].Lock()
	}
}

// BatchUnlock will try to release locks on a batch of ids.
func (sl *stripedLock) BatchUnlock(ids []string) {
	hashcodes := sl.getHashcodes(ids)
	for _, i := range hashcodes {
		sl.locks[i].Unlock()
	}
}

// getHashcodes handles id sorting and de duplication of an hashcodes for a batch of ids
func (sl *stripedLock) getHashcodes(ids []string) []uint32 {

	// sort the ids so that locks are acquired in a stable order
	sort.Strings(ids)

	var hashcodes []uint32
	for _, id := range ids {
		hashcodes = append(hashcodes, sl.idToIndex(id))
	}

	// init a map to de dupe the entries
	hashcodesMap := make(map[uint32]bool)
	for _, i := range hashcodes {
		hashcodesMap[i] = true
	}

	// init new array to return filtered values
	var ret []uint32
	for i, _ := range hashcodesMap {
		ret = append(ret, i)
	}

	fmt.Println(ret)

	return ret
}

// idToIndex is an id hashing function
func (sl *stripedLock) idToIndex(id string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(id))
	return h.Sum32() % sl.size
}
