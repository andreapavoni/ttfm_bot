package collections

import (
	"sync"
)

// Map type that can be safely shared between
// goroutines that require read/write access to a map

// type SmartList[T comparable] struct {
// 	list []T
// 	lock sync.RWMutex
// }

type SmartMap[T any] struct {
	sync.RWMutex
	items map[string]T
}

// Concurrent map item
type SmartMapItem[T any] struct {
	Key   string
	Value T
}

func NewSmartMap[T any]() *SmartMap[T] {
	return &SmartMap[T]{items: map[string]T{}}
}

// Set sets a key in a concurrent map
func (sm *SmartMap[T]) Set(key string, value T) {
	sm.Lock()
	defer sm.Unlock()

	sm.items[key] = value
}

// Delete removes a key in a concurrent map
func (sm *SmartMap[T]) Delete(key string) {
	sm.Lock()
	defer sm.Unlock()

	delete(sm.items, key)
}

// HasKey checks if map has the given key
func (sm *SmartMap[T]) HasKey(key string) bool {
	sm.Lock()
	defer sm.Unlock()

	_, ok := sm.items[key]
	return ok
}

// Keys returns a list of string keys in the map
func (sm *SmartMap[T]) Keys() (keys []string) {
	sm.Lock()
	defer sm.Unlock()

	for k := range sm.items {
		keys = append(keys, k)
	}

	return keys
}

// Get gets a key from a concurrent map
func (sm *SmartMap[T]) Get(key string) (T, bool) {
	sm.Lock()
	defer sm.Unlock()

	value, ok := sm.items[key]

	return value, ok
}

// Length gets the number of keys from a concurrent map
func (sm *SmartMap[T]) Length() int {
	sm.Lock()
	defer sm.Unlock()

	return len(sm.items)
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *SmartMap[T]) Iter() <-chan SmartMapItem[T] {
	c := make(chan SmartMapItem[T])

	f := func() {
		cm.Lock()
		defer cm.Unlock()

		for k, v := range cm.items {
			c <- SmartMapItem[T]{k, v}
		}
		close(c)
	}
	go f()

	return c
}
