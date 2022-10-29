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

type Repository interface {
	Init()
	FindAll() (any, error)
	FindOne(key string) (any, error)
	Insert(v any) error
	Update(key string, v any) error
	Delete(key string) error
	Size() int
}

type Service interface {
	InitService()
	GetRepository(key string) any
}

// Resource is here to represent a rest resource. It can be used to
// represent a single item or a collection of items and is purposely
// created to easily cater to many situations.
type Resource interface {
	Inject(s Service)
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
