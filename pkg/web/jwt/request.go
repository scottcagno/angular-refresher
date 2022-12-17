package jwt

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoTokenInRequest = errors.New("no token present in request")

func ExtractToken(r *http.Request) (string, error) {
	tokenHeader := r.Header.Get("Authorization")
	// The usual convention is for "Bearer" to be title-cased. However, there's no
	// strict rule around this, and it's best to follow the robustness principle here.
	if tokenHeader == "" || !strings.HasPrefix(strings.ToLower(tokenHeader), "bearer ") {
		return "", ErrNoTokenInRequest
	}
	return tokenHeader[7:], nil
}
