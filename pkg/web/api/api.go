package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/scottcagno/angular-refresher/pkg/web"
	"github.com/scottcagno/angular-refresher/pkg/web/api/middleware"
)

type M = map[string]any

type APIConfig struct {
	CORS   *middleware.CORSConfig
	Muxer  *http.ServeMux
	Logger *log.Logger
	//Auth   *jwt.JWTService
}

var defaultAPIConfig = &APIConfig{
	CORS:   middleware.DefaultCORSConfig,
	Muxer:  http.NewServeMux(),
	Logger: log.New(os.Stderr, "[DEFAULT] ", log.LstdFlags),
	//Auth:   nil,
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
	// if c.Auth == nil {
	// 	c.Auth = jwt.NewJWTService()
	// }
	return c
}

type API struct {
	base        string
	conf        *APIConfig
	cors        http.Handler
	logger      *log.Logger
	mux         *http.ServeMux
	sessions    *web.SessionStore
	handlers    []handler
	authService *AuthService
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
	//api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "auth/basic")), api.BasicAuthHandler())
	//api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "auth/token")), api.AuthTokenHandler())
	// if conf.Auth != nil {
	// 	api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "validate")), api.AuthHandler())
	// }
	// api.sessions = web.NewSessionStore(&web.SessionStoreConfig{
	// 	SessionID: "go_sess_id",
	// 	Domain:    "localhost",
	// 	Timeout:   time.Duration(30) * time.Minute,
	// })
	api.sessions = web.NewSessionStore(nil)
	api.handlers = make([]handler, 0)
	// api.logger.Println(api.conf.Auth.Keys())
	return api
}

// func _NewAPI(base string, cors http.Handler, logger *log.Logger, mux *http.ServeMux) *API {
// 	if logger == nil {
// 		logger = log.New(os.Stderr, "[DEFAULT] ", log.LstdFlags)
// 	}
// 	if mux == nil {
// 		mux = http.NewServeMux()
// 	}
// 	api := &API{
// 		base:     base,
// 		cors:     cors,
// 		logger:   logger,
// 		mux:      mux,
// 		handlers: make([]handler, 0),
// 	}
// 	api.mux.Handle("/", http.RedirectHandler(api.base, http.StatusSeeOther))
// 	api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "stats")), api.StatsHandler())
// 	//api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "auth/basic")), api.BasicAuthHandler())
// 	//api.mux.Handle(filepath.ToSlash(filepath.Join(api.base, "auth/token")), api.AuthTokenHandler())
// 	return api
// }

func (api *API) StatsHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, map[string]any{"routes": api.handlers})
		},
	)
}

// func (api *API) AuthTokenHandler() http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		// get the user session
// 		sess, exists := api.sessions.Get(r)
// 		if !exists {
// 			// do something here
// 			WriteJSON(w, http.StatusUnauthorized, nil)
// 			return
// 		}
// 		user, found := sess.GetUser()
// 		if !found {
// 			// do something here
// 			WriteJSON(w, http.StatusUnauthorized, nil)
// 			return
// 		}
// 		// generate JWT token
// 		token := api.conf.Auth.GenerateSignedToken(user.Username, user.Password, user.Role)
// 		WriteJSON(w, http.StatusOK, map[string]string{"results": token})
// 		return
// 	}
// 	return http.HandlerFunc(fn)
// }

// func (api *API) BasicAuthHandler() http.Handler {
// 	users := []struct {
// 		Username string
// 		Password string
// 		Role     string
// 	}{
// 		{"admin", "secret", "ROLE_ADMIN"},
// 		{"matt", "secret", "ROLE_ADMIN"},
// 		{"jane", "secret", "ROLE_USER"},
// 	}
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			username, password, hasBasicAuth := r.BasicAuth()
// 			log.Printf("running basic auth handler: u=%q, p=%q, ok=%v\n", username, password, hasBasicAuth)
// 			if !hasBasicAuth {
// 				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}
// 			var foundMatch bool
// 			for _, user := range users {
// 				if user.Username == username && user.Password == password {
// 					foundMatch = true
// 					break
// 				}
// 			}
// 			if !foundMatch {
// 				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}
// 			http.Redirect(w, r, "/", http.StatusOK)
// 		},
// 	)
// }

// func (api *API) JWTAuthorizationFilter(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		// Check for a header key of Authorization and a value starting
// 		// with Bearer
// 		bearer := r.Header.Get("Authorization")
// 		if bearer == "" || !strings.HasPrefix(bearer, "Bearer") {
// 			// We did not find an Authorization Bearer <token>, so we return a 401 status
// 			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 			return
// 		}
// 		// Pass the Bearer token (jwtToken) to our authentication method to
// 		// validate the token and extract the user info
// 		user, err := getAuthentication(bearer, api.conf.Auth)
// 		if err != nil {
// 			// We found an Authorization Bearer <token>, but the token could
// 			// not be validated, so we return the error and a 400 status
// 			msg := fmt.Sprintf("(%T) %s", err, err)
// 			http.Error(w, msg, http.StatusBadRequest)
// 			return
// 		}
// 		// We successfully got an Authorization Bearer <token> and validated
// 		// the token and have obtained the payload from the token. We will
// 		// request a session and add the payload to our session (make sure
// 		// we remember to persist the session)
// 		sess, _ := api.sessions.MustGet(r)
// 		sess.Set("user", user)
// 		api.sessions.Save(w, r, sess)
// 		// And finally, we will call the next handler in our chain
// 		next.ServeHTTP(w, r)
// 	}
// 	return http.HandlerFunc(fn)
// }
//
// func getAuthentication(bearerToken string, jwtService *jwt.JWTService) (*web.SystemUser, error) {
// 	jwtToken := strings.Split(bearerToken, " ")[1]
// 	token, err := jwtService.ValidateTokenString(jwtToken)
// 	if err != nil {
// 		if e, ok := err.(*jwt.ValidationError); ok {
// 			return nil, e
// 		}
// 		return nil, err
// 	}
// 	return &web.SystemUser{
// 		Username: token.Claims.(jwt.MapClaims)["user"].(string),
// 		Password: "",
// 		Role:     token.Claims.(jwt.MapClaims)["role"].(string),
// 	}, nil
// }

func (api *API) Register(name string, re Resource, secure bool) {
	h := &handler{
		name: name,
		path: filepath.ToSlash(filepath.Join(api.base, name)),
		reso: re,
	}
	api.handlers = append(api.handlers, *h)
	var hand http.Handler
	if secure {
		hand = api.authService.Secure(middleware.WithLogging(api.logger, h))
	} else {
		hand = middleware.WithLogging(api.logger, h)
	}
	api.mux.Handle(h.path, hand)
}

func (api *API) RegisterAuthService(base string, as *AuthService) {
	h1 := &customHandler{
		path: filepath.ToSlash(filepath.Join(base, as.RegisterPath)),
		fn:   as.Authenticator.Register,
	}
	h2 := &customHandler{
		path: filepath.ToSlash(filepath.Join(base, as.ValidatePath)),
		fn:   as.Authenticator.Validate,
	}

	api.mux.Handle(h1.path, middleware.WithLogging(api.logger, h1))
	api.mux.Handle(h2.path, middleware.WithLogging(api.logger, h2))

	api.logger.Printf("::Registered authentication service...\n")
	api.logger.Printf("::Register a token at %q\n", h1.path)
	api.logger.Printf("::Validate a token at %q\n", h2.path)

	api.authService = as
}

func (api *API) RegisterCustom(name string, re CustomResource, secure bool) {
	h := &customHandler{
		path: filepath.ToSlash(filepath.Join(api.base, name)),
		fn:   re.Custom(),
	}
	var hand http.Handler
	if secure {
		hand = api.authService.Secure(middleware.WithLogging(api.logger, h))
	} else {
		hand = middleware.WithLogging(api.logger, h)
	}
	api.mux.Handle(h.path, hand)
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
	path   string
	fn     http.HandlerFunc
	secure bool
}

func (ch customHandler) IsNil() bool {
	return ch.path == "" && ch.fn == nil
}

type handler struct {
	name   string
	path   string
	reso   Resource
	secure bool
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
