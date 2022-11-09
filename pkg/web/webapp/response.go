package webapp

import (
	"bufio"
	"net"
	"net/http"
)

// response wraps a http.ResponseWriter and implements its interface to be
// used by an HTTP handler to construct and HTTP response.
// See [http.ResponseWriter](https://pkg.go.dev/net/http#ResponseWriter)
type response struct {
	writer    http.ResponseWriter
	status    int
	size      int64
	committed bool
}

// newResponse creates and returns a new instance of a *response
func newResponse(w http.ResponseWriter) *response {
	return &response{writer: w}
}

// Header implements the http.ResponseWriter interface to allow an HTTP
// handler to return the header map for the writer that will be sent by
// WriteHeader. Changing the header after a call to WriteHeader (or Write)
// has no effect unless the modified headers were declared as trailers by
// setting the "Trailer" header before the call to WriteHeader. To suppress
// implicit response headers, set their value to nil.
func (r *response) Header() http.Header {
	return r.writer.Header()
}

// Write implements the http.ResponseWriter interface to allow an HTTP handler
// to write the data to the connection as part of an HTTP reply. If WriteHeader
// has not yet been called, Write calls WriteHeader(http.StatusOK) before
// writing the data. If the Header does not contain a Content-Type line, Write
// adds a Content-Type set to the result of passing the initial 512 bytes of
// written data to DetectContentType. Additionally, if the total size of all
// written data is under a few KB and there are no Flush calls, the Content-Length
// header is added automatically.
func (r *response) Write(bytes []byte) (int, error) {
	if !r.committed {
		if r.status == 0 {
			r.status = http.StatusOK
		}
		r.WriteHeader(r.status)
	}
	n, err := r.writer.Write(bytes)
	r.size += int64(n)
	return n, err
}

// WriteHeader implements the http.ResponseWriter interface to allow an HTTP
// handler to write a status code to the response. WriteHeader sends an HTTP
// response header with status code. If WriteHeader is not called explicitly,
// the first call to Write will trigger an implicit WriteHeader(http.StatusOK).
// Thus, explicit calls to WriteHeader are mainly use to send error codes.
func (r *response) WriteHeader(statusCode int) {
	if r.committed {
		return
	}
	r.status = statusCode
	r.writer.WriteHeader(r.status)
	r.committed = true
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://pkg.go.dev/net/http#Flusher)
func (r *response) Flush() {
	r.writer.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://pkg.go.dev/net/http#Hijacker)
func (r *response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.writer.(http.Hijacker).Hijack()
}

// reset is here to reset a response
func (r *response) reset(w http.ResponseWriter) {
	r.writer = w
	r.status = http.StatusOK
	r.size = 0
	r.committed = false
}
