package rest

import (
	"fmt"
	"mime"
	"net/http"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.data.status = statusCode
}

func logStr(code int, r *http.Request) (string, []interface{}) {
	return "# %s - - [%s] \"%s %s %s\" %d %d\n", []interface{}{
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.URL.EscapedPath(),
		r.Proto,
		code,
		r.ContentLength,
	}
}

func HandleWithLogging(logger *Logger, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("err: %v, trace: %s\n", err, debug.Stack())
			}
		}()
		lrw := loggingResponseWriter{
			ResponseWriter: w,
			data: &responseData{
				status: 200,
				size:   0,
			},
		}
		next.ServeHTTP(&lrw, r)
		if 400 <= lrw.data.status && lrw.data.status <= 599 {
			str, args := logStr(lrw.data.status, r)
			logger.Error(str, args...)
			return
		}
		str, args := logStr(lrw.data.status, r)
		logger.Info(str, args...)
		return
	}
	return http.HandlerFunc(fn)
}

func HandleMetrics(title string, ss []string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var data []string
		data = append(data, fmt.Sprintf("<h3>%s</h3>", title))
		sort.Strings(ss)
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		_, err := fmt.Fprintf(w, strings.Join(data, "<br>"))
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		return
	}
	return http.HandlerFunc(fn)
}

func HandleStatic(prefix, path string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
}
