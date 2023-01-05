package jwt

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoTokenInRequest = errors.New("no token present in request")
var ErrNoCookieFound = errors.New("no cookie present with specified name in request")

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	tokenHeader := r.Header.Get("Authorization")
	// The usual convention is for "Bearer" to be title-cased. However, there's no
	// strict rule around this, and it's best to follow the robustness principle here.
	if tokenHeader == "" || !strings.HasPrefix(strings.ToLower(tokenHeader), "bearer ") {
		return "", ErrNoTokenInRequest
	}
	return tokenHeader[7:], nil
}

func ExtractTokenFromCookie(name string, r *http.Request) (string, error) {
	c, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		return "", err
	}
	return c.Value, nil
}
