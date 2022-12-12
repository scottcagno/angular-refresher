package api

import (
	"errors"
	"sync"
)

// QueryFunc is a function that is injected with a type
// and returns a boolean
type QueryFunc[T any] func(t T) bool
type FindFunc[T any] func(query QueryFunc[T]) ([]T, error)
type ExecFunc[T any] func(query QueryFunc[T], exec QueryFunc[T]) (int, error)
type InsertFunc[K comparable, T any] func(newK K, newT T) error
type UpdateFunc[K comparable, T any] func(oldK K, newT T) error
type DeleteFunc[K comparable] func(oldK K) error

// Repository is an interface for representing a generic
// data repository
type Repository[T any, K comparable] interface {

	// Find provides the user with an interface for
	// locating, retrieving or viewing one or more
	// entries.
	Find(query QueryFunc[T]) ([]T, error)

	// FindOne provides the user with an interface for
	// locating, retrieving or viewing one entity. It
	// guarantees that only one entry is returned. If
	// more than one entry is found, it returns the first
	// one matched.
	FindOne(query QueryFunc[T]) (T, error)

	// Exec provides the user with an interface for
	// locating and executing a function on the items matching
	// the query criteria.
	Exec(query QueryFunc[T], exec QueryFunc[T]) (int, error)

	// Insert provides the user with an interface for
	// creating and adding new entries.
	Insert(newK K, newT T) error

	// Update provides the user with an interface for
	// editing or updating an existing entry.
	Update(oldK K, newT T) error

	// Delete provides the user with an interface for
	// removing or invalidating an existing entry.
	Delete(oldK K) error

	// Type returns the data type that is used with the Repository
	Type() T

	// KeyType returns the data type that is used as a primary key with the Repository
	KeyType() K
}

type MemoryRepository[T any, K comparable] struct {
	lock     sync.Mutex
	isLocked bool
	data     map[K]T
}

func NewMemoryRepository[T any, K comparable]() *MemoryRepository[T, K] {
	return &MemoryRepository[T, K]{
		data: make(map[K]T),
	}
}

func (repo *MemoryRepository[T, K]) Find(query QueryFunc[T]) ([]T, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	if len(repo.data) < 1 {
		return nil, errors.New("error: cannot find anything because there is no data")
	}
	var res []T
	for k, t := range repo.data {
		if query(t) {
			res = append(res, repo.data[k])
		}
	}
	if len(res) == 0 {
		return nil, errors.New("error: query did not match anything")
	}
	return res, nil
}

func (repo *MemoryRepository[T, K]) FindOne(query QueryFunc[T]) (T, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	var res T
	if len(repo.data) < 1 {
		return res, errors.New("error: cannot find anything because there is no data")
	}
	var foundMatch bool
	for k, t := range repo.data {
		if query(t) {
			res = repo.data[k]
			foundMatch = true
			break
		}
	}
	if !foundMatch {
		return res, errors.New("error: query did not match anything")
	}
	return res, nil
}

func (repo *MemoryRepository[T, K]) Exec(query QueryFunc[T], exec QueryFunc[T]) (int, error) {
	var res []T
	f1 := func() (int, error) {
		repo.lock.Lock()
		defer repo.lock.Unlock()
		if len(repo.data) < 1 {
			return -1, errors.New("error: cannot find anything because there is no data")
		}

		for k, t := range repo.data {
			if query(t) {
				res = append(res, repo.data[k])
			}
		}
		if len(res) == 0 {
			return 0, errors.New("error: query did not match anything")
		}
		return len(res), nil
	}
	f2 := func() (int, error) {
		var ops int
		for i := range res {
			if exec(res[i]) {
				ops++
			}
		}
		return ops, nil
	}
	n, err := f1()
	if err != nil {
		return n, err
	}
	return f2()
}

func (repo *MemoryRepository[T, K]) Insert(newK K, newT T) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	_, exists := repo.data[newK]
	if exists {
		return errors.New("error: cannot insert, item already exists")
	}
	repo.data[newK] = newT
	return nil
}

func (repo *MemoryRepository[T, K]) Update(oldK K, newT T) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	_, exists := repo.data[oldK]
	if !exists {
		return errors.New("error: cannot update, item does not exist")
	}
	repo.data[oldK] = newT
	return nil
}

func (repo *MemoryRepository[T, K]) Delete(oldK K) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	_, exists := repo.data[oldK]
	if !exists {
		return errors.New("error: cannot remove, item is not present")
	}
	delete(repo.data, oldK)
	return nil
}

func (repo *MemoryRepository[T, K]) Type() (t T) {
	return t
}

func (repo *MemoryRepository[T, K]) KeyType() (k K) {
	return k
}
