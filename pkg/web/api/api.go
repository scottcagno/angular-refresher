package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/scottcagno/angular-refresher/pkg/web/api/middleware"
)

type M = map[string]any

type APIConfig struct {
	CORS   *middleware.CORSConfig
	Muxer  *http.ServeMux
	Logger *log.Logger
	Auth   map[string]string
}

var defaultAPIConfig = &APIConfig{
	CORS:   middleware.DefaultCORSConfig,
	Muxer:  http.NewServeMux(),
	Logger: log.New(os.Stderr, "[DEFAULT] ", log.LstdFlags),
	Auth:   nil,
}

func checkConf(c *APIConfig) *APIConfig {
	if c == nil {
		c = defaultAPIConfig
	}
	if c.Muxer == nil {
		c.Muxer = http.NewServeMux()
	}
	if c.Logger == nil {
		c.Logger = log.New(os.Stderr, "[DEFAULT] ", log.LstdFlags)
	}
	if c.Auth == nil {
		c.Auth = make(map[string]string)
	}
	return c
}

type API struct {
	base     string
	conf     *APIConfig
	cors     http.Handler
	logger   *log.Logger
	mux      *http.ServeMux
	handlers []handler
}

func NewAPI(base string, conf *APIConfig) *API {
	conf = checkConf(conf)
	api := new(API)
	api.base = base
	api.conf = conf
	if conf.CORS != nil {
		api.cors = middleware.CORSHandler(conf.CORS)
	}
	api.logger = conf.Logger
	api.mux = conf.Muxer
	api.mux.Handle("/", http.RedirectHandler(api.base, http.StatusSeeOther))
	api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "stats")), api.StatsHandler())
	if conf.Auth != nil && len(conf.Auth) > 0 {
		api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "validate")), api.BasicAuthHandler())
	}
	api.handlers = make([]handler, 0)
	return api
}

func _NewAPI(base string, cors http.Handler, logger *log.Logger, mux *http.ServeMux) *API {
	if logger == nil {
		logger = log.New(os.Stderr, "[DEFAULT] ", log.LstdFlags)
	}
	if mux == nil {
		mux = http.NewServeMux()
	}
	api := &API{
		base:     base,
		cors:     cors,
		logger:   logger,
		mux:      mux,
		handlers: make([]handler, 0),
	}
	api.mux.Handle("/", http.RedirectHandler(api.base, http.StatusSeeOther))
	api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "stats")), api.StatsHandler())
	api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "validate")), api.BasicAuthHandler())
	return api
}

func (api *API) StatsHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, map[string]any{"routes": api.handlers})
		},
	)
}

func (api *API) BasicAuthHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			log.Printf("running basic auth handler: u=%q, p=%q, ok=%v\n", username, password, ok)
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			storedPass, hasPass := api.conf.Auth[username]
			if !hasPass || storedPass != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/", http.StatusOK)
		},
	)
}

func (api *API) Register(name string, re Resource) {
	h := &handler{
		name: name,
		path: filepath.ToSlash(filepath.Join(api.base, name)),
		reso: re,
	}
	api.handlers = append(api.handlers, *h)
	api.mux.Handle(h.path, middleware.WithLogging(api.logger, h))
}

func (api *API) RegisterCustom(name string, re CustomResource) {
	h := &customHandler{
		path: filepath.ToSlash(filepath.Join(api.base, name)),
		fn:   re.Custom(),
	}
	api.mux.Handle(h.path, middleware.WithLogging(api.logger, h))
}

// func (api *API) _RegisterSecure(name string, re SecureResource) {
// 	h := &customHandler{
// 		path: filepath.ToSlash(filepath.Join(api.base, name)),
// 		fn:   api.Secure(re),
// 	}
// 	api.mux.Handle(h.path, middleware.WithLogging(api.logger, h))
// }

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// apply cors handler if we have one
	if api.cors != nil {
		api.cors.ServeHTTP(w, r)
	}
	// lookup resource handler
	rh, pat := api.mux.Handler(r)
	// do something with the pattern if we need to
	_ = pat
	// if strings.HasSuffix(pat, "/validate") {
	// 	rh.ServeHTTP(w, r)
	// }
	// call the resource handler
	fmt.Println(">>>", pat)
	rh.ServeHTTP(w, r)
}

// func (api *API) HandleRequestMapping(mapping RequestMapping) {
// 	fmt.Printf("type=%T, value=%#v\n", mapping, mapping)
// }
//
// func (api *API) HandleRequestMappingFunc(reso ResourceV2, handler http.Handler) {
// 	fmt.Printf("type=%T, value=%#v\n", reso, reso)
// 	fmt.Printf("type=%T, value=%#v\n", handler, handler)
// }

type customHandler struct {
	path string
	fn   http.HandlerFunc
}

func (ch customHandler) IsNil() bool {
	return ch.path == "" && ch.fn == nil
}

type handler struct {
	name string
	path string
	reso Resource
}

func (h handler) String() string {
	return fmt.Sprintf("name=%q, path=%q", h.name, h.path)
}

// func writeReqCtx(r *http.Request, key, val any) *http.Request {
// 	// create a new context from the parent context in the incoming request
// 	ctx := context.WithValue(r.Context(), key, val)
// 	// create and return a new request using the new context
// 	return r.WithContext(ctx)
// }
//
// func readReqCtx(r *http.Request, key any) any {
// 	// get the value from the request context
// 	return r.Context().Value(key)
// }

func HasParam(r *http.Request, key string) bool {
	params := r.URL.Query()
	if params.Has(key) {
		return true
	}
	return false
}

func GetParam(r *http.Request, key string) (string, bool) {
	params := r.URL.Query()
	if params.Has(key) {
		return params.Get(key), true
	}
	return "", false
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.reso.Get(w, r)
	case http.MethodPost:
		h.reso.Add(w, r)
	case http.MethodPut:
		h.reso.Set(w, r)
	case http.MethodDelete:
		h.reso.Del(w, r)
	case http.MethodOptions:
		middleware.Options(w, r)
	default:
		middleware.NotFound(w, r)
	}
}

func (h *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check for a custom handler request
	if h.path == r.URL.Path {
		h.fn(w, r)
		return
	}
}
