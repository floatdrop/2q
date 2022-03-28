package twoqueue

import (
	"github.com/floatdrop/lru"
)

const (
	// Default2QRecentRatio is the ratio of the 2Q cache dedicated
	// to recently added entries that have only been accessed once.
	Default2QRecentRatio = 0.25

	// Default2QGhostEntries is the default ratio of ghost
	// entries kept to track entries recently evicted
	Default2QGhostEntries = 0.50
)

// TwoQueue is a thread-safe fixed size 2Q cache.
// 2Q is an enhancement over the standard LRU cache
// in that it tracks both frequently and recently used
// entries separately. This avoids a burst in access to new
// entries from evicting frequently used entries. It adds some
// additional tracking overhead to the standard LRU cache, and is
// computationally about 2x the cost, and adds some metadata over
// head.
type TwoQueue[K comparable, V any] struct {
	recent      *lru.LRU[K, V]        // A1in in paper (should have type FIFO[entry[K, V]])
	recentEvict *lru.LRU[K, struct{}] // A1out in paper (should be FIFO[k])
	frequent    *lru.LRU[K, V]        // Am in paper
}

// Evicted holds key/value pair that was evicted from cache.
type Evicted[K comparable, V any] struct {
	Key K
	Value V
}

func fromLruEvicted[K comparable, V any](e *lru.Evicted[K, V]) *Evicted[K, V] {
	if e == nil {
		return nil
	}

	return &Evicted[K, V]{
		e.Key,
		e.Value,
	}
}

// Get probes frequent and recent cached items and returns pointer to value (or nil if it was not found).
func (L *TwoQueue[K, V]) Get(key K) *V {
	if e := L.frequent.Get(key); e != nil {
		return e
	}

	if e := L.recent.Peek(key); e != nil {
		return e
	}

	return nil
}

// Set stores key/value pair in 2Q cache following 2Q Full Version promotion algorytm.
func (L *TwoQueue[K, V]) Set(key K, value V) *Evicted[K, V] {
	if e := L.frequent.Peek(key); e != nil {
		return fromLruEvicted(L.frequent.Set(key, value))
	}

	if L.recentEvict.Peek(key) != nil {
		L.recentEvict.Remove(key)
		return fromLruEvicted(L.frequent.Set(key, value))
	}

	if e := L.recent.Peek(key); e != nil {
		return nil
	} else if re := L.recent.Set(key, value); re != nil {
		L.recentEvict.Set(re.Key, struct{}{})
		return fromLruEvicted(re)
	}

	return nil
}

// Len returns size of cache (frequent + recent items)
func (L *TwoQueue[K, V]) Len() int {
	return L.frequent.Len() + L.recent.Len()
}

// Peek returns value for key (if key was in cache), but does not modify its recency.
func (L *TwoQueue[K, V]) Peek(key K) *V {
	if e := L.frequent.Get(key); e != nil {
		return e
	}

	return L.recent.Peek(key)
}

// Remove method removes entry associated with key and returns pointer to removed value (or nil if entry was not in cache).
func (L *TwoQueue[K, V]) Remove(key K) *V {
	if e := L.frequent.Remove(key); e != nil {
		return e
	}

	return L.recent.Remove(key)
}

// New creates 2Q cache with specified capacities:
//
// - Kin defines A1in FIFO size for key/value pairs
// - Kout defines A1out FIFO size for keys
// - size defines frequent LRU size for key/value pairs
//
// It's recommended to hold 25% of available memory in Kin. Kout size should correspond
// to 50% memory for values. And size should consume rest of memory. You can refer to
// original paper (http://www.vldb.org/conf/1994/P439.PDF) for computing sizes.
//
// For example, if you can store around 10000 items in cache:
// - Kin should hold around 2500 items.
// - Kout should hold 5000 items.
// - And size should take the rest 7500 items.
//
// Cache will preallocate size count of internal structures to avoid allocation in process.
func NewParams[K comparable, V any](Kin int, Kout int, size int) *TwoQueue[K, V] {
	return &TwoQueue[K, V]{
		recent:      lru.New[K, V](Kin),
		recentEvict: lru.New[K, struct{}](Kout),
		frequent:    lru.New[K, V](size),
	}
}

// New creates 2Q cache with predefined size splits. 25% of size goes to Kin, 50% to KOut and rest to Am size.
func New[K comparable, V any](size int) *TwoQueue[K, V] {
	return NewParams[K, V](
		int(Default2QRecentRatio * float64(size)),
		int(Default2QGhostEntries * float64(size)),
		int((1 - Default2QRecentRatio) * float64(size)),
	)
}