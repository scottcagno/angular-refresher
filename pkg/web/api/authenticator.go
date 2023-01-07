package api

import (
	"errors"
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
	jwts  *jwt.JWTService
	users *UserStore
}

func NewJWTAuthService(defaultUser *web.SystemUser, privateKeyFile, publicKeyFile string) *JWTAuthService {
	jwtAuthService := &JWTAuthService{
		jwts:  jwt.NewJWTService(privateKeyFile, publicKeyFile),
		users: NewUserStore(),
	}
	if defaultUser != nil {
		jwtAuthService.users.AddUser(defaultUser.Username, defaultUser.Password, defaultUser.Role)
	}
	return jwtAuthService
}

func (js *JWTAuthService) Register(w http.ResponseWriter, r *http.Request) {
	// Check for HTTP basic authentication
	username, password, hasBasicAuth := r.BasicAuth()
	if !hasBasicAuth {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var user *web.SystemUser
	// Received valid HTTP basic authentication token; check store
	user = js.users.GetUser(username)
	if user == nil {
		// User was not found, we must fill out a new one based on the provided
		// basic authentication we received.
		user.Username = username
		user.Password = password
		user.Role = "ROLE_USER"
	}
	// Found valid user in store, generate and sign new token
	tokenString := js.jwts.GenerateSignedToken(user.Username, user.Password, user.Role)
	// Create a new cookie, and add the JWT token string to the cookie
	chocoChip := &http.Cookie{
		Name:       "token",
		Value:      tokenString,
		Path:       "",
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
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
	_, err = js.jwts.ValidateTokenString(tokenString)
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
