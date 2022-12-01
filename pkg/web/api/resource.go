package api

import (
	"net/http"
)

// RequestMappingFunc is a function that is used to validate if a given
// Handler should run or not. It is used as part of a RequestMapping. It
// Takes a method, a path and any (if any) required query param keys and
// should return an appropriate and valid HTTP status code.
// type RequestMappingFunc = func(method, path string, requiredParams ...string) int

// RequestMapping is here to fulfill custom API requests. It takes a
// RequestMappingFunc and a Handler. The provided handler is not called
// if the RequestMappingFunc returns any HTTP code other than HTTP 200 OK.
// type RequestMapping func(mapping RequestMappingFunc, next http.Handler) http.Handler
type RequestMapping struct {
	Method            string
	Path              string
	RequiredParamKeys []string
}

type Resource = ResourceV1

type ResourceV2 interface {
	// RequestMapping(mapping RequestMappingFunc, next http.Handler) http.Handler
	RequestMappingFunc(mapping RequestMapping) int
}

type CustomResource interface {
	// Custom is a custom defined handler
	Custom() http.HandlerFunc
}

// ResourceV1 is here to represent a rest resource. It can be used to
// represent a single item or a collection of items and is purposely
// created to easily cater to many situations.
type ResourceV1 interface {
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

// SecureResource is here to represent a secure rest resource. It can be
// used to represent a single item or a collection of items and is purposely
// created to easily cater to many situations. It differs from Resource
// by automatically integrating with the API in a secure fashion, expecting
// secure tokens to be exchanged before being able to process the request.
type SecureResource interface {
}
