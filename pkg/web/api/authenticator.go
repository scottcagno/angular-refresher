package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/scottcagno/angular-refresher/pkg/web"
	"github.com/scottcagno/angular-refresher/pkg/web/jwt"
)

const (
	registerEndpoint = "/register"
	validateEndpoint = "/validate"
)

type Authenticator interface {
	Register(w http.ResponseWriter, r *http.Request)
	Validate(w http.ResponseWriter, r *http.Request)
}

type AuthService struct {
	RegisterPath string
	ValidatePath string
	Authenticator
}

func MakeAuthService(authenticator Authenticator) *AuthService {
	return &AuthService{
		RegisterPath:  registerEndpoint,
		ValidatePath:  validateEndpoint,
		Authenticator: authenticator,
	}
}

func (s *AuthService) Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.Authenticator.Validate(w, r)
			next.ServeHTTP(w, r)
		},
	)
}

type JWTAuthService struct {
	Service *jwt.JWTService
	Users   *UserStore
}

func NewJWTAuthService(privateKeyFile, publicKeyFile string, defaultUsers ...*web.SystemUser) *JWTAuthService {
	jwtAuthService := &JWTAuthService{
		Service: jwt.NewJWTService(privateKeyFile, publicKeyFile),
		Users:   NewUserStore(),
	}
	if defaultUsers != nil {
		for _, user := range defaultUsers {
			jwtAuthService.Users.AddUser(user.Username, user.Password, user.Role)
		}
	}
	return jwtAuthService
}

func (js *JWTAuthService) Register(w http.ResponseWriter, r *http.Request) {
	// Check for HTTP basic authentication
	username, password, hasBasicAuth := r.BasicAuth()
	log.Println(username, password, hasBasicAuth)
	if !hasBasicAuth {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var user *web.SystemUser
	// Received valid HTTP basic authentication token; check store
	user = js.Users.GetUser(username)
	if user == nil {
		// User was not found, we must fill out a new one based on the provided
		// basic authentication we received.
		user.Username = username
		user.Password = password
		user.Role = "ROLE_USER"
	}
	// Found valid user in store, generate and sign new token
	tokenString := js.Service.GenerateSignedToken(user.Username, user.Password, user.Role)
	// Create a new cookie, and add the JWT token string to the cookie
	chocoChip := &http.Cookie{
		Name:       "token",
		Value:      tokenString,
		Path:       "/",
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     js.Service.ExpireTime().Second(),
		Secure:     true, // set true, when in production
		HttpOnly:   true,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	// Add the cookie to the request
	http.SetCookie(w, chocoChip)
}

func (js *JWTAuthService) Validate(w http.ResponseWriter, r *http.Request) {
	var err error
	// Check for a cookie containing our JWT string
	chocoChip, err := r.Cookie("token")
	if err != nil && errors.Is(err, http.ErrNoCookie) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	tokenString := chocoChip.Value

	// // Check the header for an Authorization Bearer
	// bearer := r.Header.Get("Authorization")
	// if bearer == "" || !strings.HasPrefix(bearer, "Bearer") {
	// 	// We did not find an Authorization Bearer <token>, so we return a 401 status
	// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	// 	return
	// }
	// tokenString := bearer

	// Attempt to validate the token string we found in the cookie
	_, err = js.Service.ValidateTokenString(tokenString)
	if err != nil {
		if _, is := err.(*jwt.ValidationError); is {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// If we get here, then our token was valid, return 200 OK
	w.WriteHeader(http.StatusOK)
	return
}

func VerifyJWT(service *jwt.JWTService, next http.HandlerFunc) http.HandlerFunc {
	const errCode = http.StatusUnauthorized
	return func(w http.ResponseWriter, r *http.Request) {
		//
		// Check for a cookie containing a JWT token
		//
		cook, err := r.Cookie("token")
		if err != nil && errors.Is(err, http.ErrNoCookie) {
			// No cookie found, return status unauthorized
			http.Error(w, http.StatusText(errCode), errCode)
			return
		}
		tokenString := cook.Value
		// //
		// // Check the header for an Authorization Bearer
		// //
		// bearer := r.Header.Get("Authorization")
		// if bearer == "" || !strings.HasPrefix(bearer, "Bearer") {
		// 	// No Bearer token found, return status unauthorized
		// 	http.Error(w, http.StatusText(errCode), errCode)
		// 	return
		// }
		// tokenString := bearer

		//
		// Attempt to validate the token string
		//
		_, err = service.ValidateTokenString(tokenString)
		if err != nil {
			if _, is := err.(*jwt.ValidationError); is {
				http.Error(w, http.StatusText(errCode), errCode)
				return
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}
}
