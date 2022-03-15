package collections

import (
	"errors"
	"sync"

	"github.com/andreapavoni/ttfm_bot/utils"
)

type SmartList[T comparable] struct {
	list []T
	lock sync.RWMutex
}

func NewSmartList[T comparable]() *SmartList[T] {
	return &SmartList[T]{list: make([]T, 0)}
}

func NewSmartListFromSlice[T comparable](slice []T) *SmartList[T] {
	list := NewSmartList[T]()

	for _, item := range slice {
		list.Push(item)
	}

	return list
}

func (l *SmartList[T]) Push(value T) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.list = append(l.list, value)
}

func (l *SmartList[T]) Shift() (T, error) {
	var elem T
	if l.Size() > 0 {
		l.lock.Lock()
		defer l.lock.Unlock()

		elem = l.list[0]
		l.list = l.list[1:]
		return elem, nil
	}
	return elem, errors.New("List is empty")
}

func (l *SmartList[T]) Remove(value T) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if idx := utils.IndexOf(value, l.list); idx >= 0 {
		l.list = append(l.list[:idx], l.list[idx+1:]...)
		return nil
	}

	return errors.New("Element not found")
}

func (l *SmartList[T]) HasElement(value T) bool {
	return utils.IndexOf(value, l.list) >= 0
}

func (l *SmartList[T]) List() []T {
	return l.list
}

func (l *SmartList[T]) IndexOf(value T) int {
	return utils.IndexOf(value, l.list)
}

func (l *SmartList[T]) Find(f func(*T) bool) (*T, int) {
	for idx, item := range l.list {
		if f(&item) {
			return &item, idx
		}
	}
	return (new(T)), -1
}

func (l *SmartList[T]) Size() int {
	return len(l.list)
}

func (l *SmartList[T]) Empty() {
	l.list = make([]T, 0)
}
