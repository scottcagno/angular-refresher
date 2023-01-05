package jwt

import (
	"crypto/rsa"
	"fmt"
)

func ValidateTokenString(signingMethod SigningMethod, public *rsa.PublicKey, tokenString string) (*Token, error) {
	parser := NewParser([]string{signingMethod.Alg()}, true, false)
	token, err := parser.Parse(
		tokenString, func(t *Token) (any, error) {
			// Validate the signing method
			if t.Method.Alg() != signingMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return public, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(MapClaims); ok && token.Valid {
		// you could check some claims in here if you wanted...
		_ = claims
	}
	return token, nil
}
