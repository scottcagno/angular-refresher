package webapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// A Context carries http response and request along with other values
// across API boundaries. It is the main component involved in the
// Handler type in this package.
type Context interface {
	Response() *response
}

// context is the concrete implementation of a Context. The exported
// Context is here to preserve backward compatibility and to allow the
// library to be extendable.
type context struct {
	res     *response
	req     *http.Request
	path    string
	query   url.Values
	handler RouteHandler
}

func (c *context) Response() *response {
	return c.res
}

func (c *context) Reset(w http.ResponseWriter, r *http.Request) {
	c.res.reset(w)
	c.req = r
	c.path = ""
	c.query = nil
	c.handler = RouteHandlerFunc(NotFound)
}

func (c *context) writeContentType(value string) {
	header := c.res.Header()
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", value)
	}
}

func (c *context) Raw(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)
	c.res.WriteHeader(code)
	_, err = c.res.Write(b)
	return
}

func (c *context) String(code int, msg string) error {
	return c.Raw(code, "text/plain", []byte(msg))
}

func (c *context) JSON(code int, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		code := http.StatusExpectationFailed
		format := fmt.Sprintf("{%q:%d,%q:%s}", "code", code, "message", http.StatusText(code))
		return c.Raw(http.StatusInternalServerError, "application/json", []byte(format))
	}
	return c.Raw(code, "application/json", b)
}

func (c *context) NoContent(code int) error {
	c.res.WriteHeader(code)
	return nil
}

func (c *context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return http.ErrAbortHandler
	}
	c.res.Header().Set("Content-Location", url)
	c.res.WriteHeader(code)
	return nil
}

func RedirectHandler(code int, url string) RouteHandler {
	return RouteHandlerFunc(func(c Context) error {
		if code < 300 || code > 308 {
			return http.ErrBodyNotAllowed
		}
		c.Response().Header().Set("Content-Location", url)
		c.Response().WriteHeader(code)
		return nil
	})
}

func (c *context) Error(err error) {
	ErrorHandler(err, c)
}

func ErrorHandler(err error, c Context) {
	w := c.Response()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(CodeFromError(err))
	fmt.Fprintln(w, err)
}

var NotFound = func(c Context) error {
	return ErrorFromCode(http.StatusNotFound)
}

var MethodNotAllowed = func(c Context) error {
	return ErrorFromCode(http.StatusMethodNotAllowed)
}

var ExpectationFailed = func(c Context) error {
	return ErrorFromCode(http.StatusExpectationFailed)
}

var InternalServerError = func(c Context) error {
	return ErrorFromCode(http.StatusInternalServerError)
}
