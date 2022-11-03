package api

import (
	"errors"
	"net/http"
)

var (
	// repository errors
	ErrNone         = errors.New("there is nothing to return")
	ErrNoMatchFound = errors.New("nothing matched the provided key")
	ErrBadType      = errors.New("type provided was incorrect")
)

// type Repository interface {
// 	FindAll() (any, error)
// 	FindOne(id string) (any, error)
// 	Insert(v any) error
// 	Update(id string, t any) error
// 	Delete(id string) error
// 	Size() int
// 	GetRepositorySet() any
// 	Init(data map[string]any)
// }

//
// type Repository[T any] interface {
// 	FindAll() (T, error)
// 	FindOne(id string) (T, error)
// 	Insert(v any) error
// 	Update(id string, t T) error
// 	Delete(id string) error
// 	Size() int
// }

// type Service interface {
// 	InitService()
// 	AddRepository(key string, val Repository)
// 	GetRepository(key string) Repository
// }

// Resource is here to represent a rest resource. It can be used to
// represent a single item or a collection of items and is purposely
// created to easily cater to many situations.
type Resource interface {
	// Get implements http.Handler and is responsible for locating
	// and returns all the implementing resource items, or one if a
	// matching identifier is found. Note: the user implementing this
	// interface is responsible for obtaining the identifier from the
	// request.
	Get(w http.ResponseWriter, r *http.Request)
	// Add implements http.Handler and is responsible for locating
	// the provided serialized resource item (written to the request
	// body) and adding the serialized item to the resource set.
	Add(w http.ResponseWriter, r *http.Request)
	// Set implements http.Handler and is responsible for locating
	// an identifier along with a serialized resource item (written to
	// the request body) and updating the resource item that has a
	// matching identifier.
	Set(w http.ResponseWriter, r *http.Request)
	// Del implements http.Handler and is responsible for locating
	// an identifier and deleting the resource item with the matching
	// identifier.
	Del(w http.ResponseWriter, r *http.Request)
}
