package api

import (
	"net/http"
)

// Resource is here to represent a rest resource. It can be used to
// represent a single item or a collection of items and is purposely
// created to easily cater to many situations.
type Resource interface {
	// Init is will automatically be called when the Resource is
	// registered by an API server
	Init()
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
