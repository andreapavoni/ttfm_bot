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

// Sets a key in a concurrent map
func (sm *SmartMap[T]) Set(key string, value T) {
	sm.Lock()
	defer sm.Unlock()

	sm.items[key] = value
}

// Removes a key in a concurrent map
func (sm *SmartMap[T]) Delete(key string) {
	sm.Lock()
	defer sm.Unlock()

	delete(sm.items, key)
}

// Gets a key from a concurrent map
func (sm *SmartMap[T]) Get(key string) (T, bool) {
	sm.Lock()
	defer sm.Unlock()

	value, ok := sm.items[key]

	return value, ok
}

// Gets the number of keys from a concurrent map
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
