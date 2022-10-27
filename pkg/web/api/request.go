package api

import (
	"net/http"
	"net/url"
)

type Request struct {
	Method  string
	Path    string
	Params  url.Values
	Handler http.Handler
}

func NewRequest(method, path string, h http.Handler) *Request {
	uri, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return &Request{
		Method:  method,
		Path:    uri.Path,
		Params:  uri.Query(),
		Handler: h,
	}
}

func match(r *http.Request, method, path string, params url.Values) bool {
	if r.Method != method {
		return false
	}
	if r.URL.Path != path {
		return false
	}
	if params != nil && len(params) > 0 {
		reqParams := r.URL.Query()
		for k, _ := range params {
			if !reqParams.Has(k) {
				return false
			}
		}
	}
	return true
}

func (r *Request) Matches(req *http.Request) bool {
	if r.Method != req.Method {
		return false
	}
	if r.Path != req.URL.Path {
		return false
	}
	if r.Params != nil && len(r.Params) > 0 {
		reqParams := req.URL.Query()
		for key, _ := range r.Params {
			if !reqParams.Has(key) {
				return false
			}
		}
	}
	return true
}
