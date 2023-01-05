package api

import (
	"errors"
	"sync"
)

var (
	ErrExists   = errors.New("already exists")
	ErrNotFound = errors.New("not found")
)

type MemoryStore[K comparable, V any] struct {
	store   sync.Map
	zeroVal V
}

func NewMemoryStore[K comparable, V any]() *MemoryStore[K, V] {
	return &MemoryStore[K, V]{
		store: sync.Map{},
	}
}

func (ms *MemoryStore[K, V]) Add(k K, v V) error {
	_, added := ms.store.LoadOrStore(k, v)
	if !added {
		return ErrExists
	}
	return nil
}

func (ms *MemoryStore[K, V]) Set(k K, v V) {
	ms.store.Store(k, v)
}

func (ms *MemoryStore[K, V]) Get(k K) (V, error) {
	v, found := ms.store.Load(k)
	if !found {
		return ms.zeroVal, ErrNotFound
	}
	return v.(V), nil
}

func (ms *MemoryStore[K, V]) Del(k K) (V, error) {
	prev, present := ms.store.LoadAndDelete(k)
	if !present {
		return ms.zeroVal, ErrNotFound
	}
	return prev.(V), nil
}
