package ttfm

import (
	"errors"
	"sync"
)

type SmartList struct {
	list []interface{}
	lock sync.RWMutex
}

func NewSmartList() *SmartList {
	return &SmartList{list: make([]interface{}, 0)}
}

func NewSmartListFromSlice[T any](slice []T) *SmartList {
	list := NewSmartList()

	for _, item := range slice {
		list.Push(item)
	}

	return list
}

func (l *SmartList) Push(value any) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.list = append(l.list, value)
}

func (l *SmartList) Shift() (any, error) {
	if l.Size() > 0 {
		l.lock.Lock()
		defer l.lock.Unlock()

		elem := l.list[0]
		l.list = l.list[1:]
		return elem, nil
	}
	return nil, errors.New("List is empty")
}

func (l *SmartList) Remove(value any) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if idx := l.findElement(value); idx >= 0 {
		l.list = append(l.list[:idx], l.list[idx+1:]...)
		return nil
	}

	return errors.New("Element not found")
}

func (l *SmartList) List() []interface{} {
	return l.list
}

func ListElements[T any](l *SmartList, list []T) []T {
	l.lock.Lock()
	defer l.lock.Unlock()

	for _, elem := range l.list {
		list = append(list, elem.(T))
	}

	return list
}

func (l *SmartList) HasElement(value any) bool {
	return l.findElement(value) >= 0
}

func (l *SmartList) findElement(value any) int {
	for idx, elem := range l.list {
		if elem == value {
			return idx
		}
	}
	return -1
}

func (l *SmartList) Size() int {
	return len(l.list)
}

func (l *SmartList) Empty() bool {
	return l.Size() == 0
}

// Map type that can be safely shared between
// goroutines that require read/write access to a map
type SmartMap struct {
	sync.RWMutex
	items map[string]any
}

// Concurrent map item
type SmartMapItem struct {
	Key   string
	Value any
}

func NewSmartMap() *SmartMap {
	return &SmartMap{items: map[string]any{}}
}

// Sets a key in a concurrent map
func (sm *SmartMap) Set(key string, value any) {
	sm.Lock()
	defer sm.Unlock()

	sm.items[key] = value
}

// Removes a key in a concurrent map
func (sm *SmartMap) Delete(key string) {
	sm.Lock()
	defer sm.Unlock()

	delete(sm.items, key)
}

// Gets a key from a concurrent map
func (sm *SmartMap) Get(key string) (interface{}, bool) {
	sm.Lock()
	defer sm.Unlock()

	value, ok := sm.items[key]

	return value, ok
}

// Gets the number of keys from a concurrent map
func (sm *SmartMap) Length() int {
	sm.Lock()
	defer sm.Unlock()

	return len(sm.items)
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *SmartMap) Iter() <-chan SmartMapItem {
	c := make(chan SmartMapItem)

	f := func() {
		cm.Lock()
		defer cm.Unlock()

		for k, v := range cm.items {
			c <- SmartMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}
