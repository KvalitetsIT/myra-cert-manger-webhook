package testutil

import (
	"sync"
)

type AtomicMap[K comparable, T any] struct {
	entries map[K]T
	sync.RWMutex
	name string
}

func NewAtomicMap[K comparable, T any](name string) *AtomicMap[K, T] {
	return &AtomicMap[K, T]{
		entries: make(map[K]T),
		name:    name,
	}
}

func (s *AtomicMap[K, T]) GetAll() []T {
	s.RLock()
	defer s.RUnlock()

	all := make([]T, 0, len(s.entries))
	for _, entry := range s.entries {
		all = append(all, entry)
	}
	return all
}

func (s *AtomicMap[K, T]) Get(key K) (T, bool) {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	value, found := s.entries[key]

	return value, found
}

func (s *AtomicMap[K, T]) Set(key K, value T) T {
	s.Lock()
	defer s.Unlock()
	s.entries[key] = value
	return value
}

func (s *AtomicMap[K, T]) Delete(key K) (T, bool) {
	s.Lock()
	defer s.Unlock()

	record, found := s.entries[key]
	if !found {
		var zero T
		return zero, false
	}

	delete(s.entries, key)
	return record, true
}
