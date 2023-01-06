package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// KeyFunc will be used by the Parse functions as a callback function in order
// to supply the key for verification. The function receives the parsed but
// un-verified Token. This allows you to use properties in the Header of the
// token (like "kid") in order to identify which key to use.
// *** MAKE SURE YOU RETURN A PUBLIC KEY IN HERE ***
type KeyFunc func(*Token) (any, error)

// Token represents a JWT. Different fields will be used depending on if you
// are creating, parsing or verifying a token.
type Token struct {
	// Raw is the raw token. It is populated when you parse a token.
	Raw string

	// Method is the signing method that is used to sign the token.
	Method SigningMethod

	// Header is the first segment of the token and holds the token
	// type and the algorithm of the signing method.
	Header map[string]any

	// Claims is the second segment of a token and holds the payload.
	Claims Claims

	// Signature is the third segment of a token and is populated when you
	// parse a token.
	Signature string

	// Valid is populated when you parse or verify a token.
	Valid bool
}

// NewToken creates and returns a new Token with the specified signing method
// and an empty map of claims.
func NewToken(method SigningMethod) *Token {
	return NewTokenWithClaims(method, make(MapClaims))
}

// NewTokenWithClaims creates and returns a new Token with the specified signing
// method and provided claims.
func NewTokenWithClaims(method SigningMethod, claims Claims) *Token {
	return &Token{
		Method: method,
		Header: map[string]any{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
	}
}

// Sign generates the signing string. This is the most expensive call. Try to use
// SignedString when possible, and if the signature already exists, it will be
// returned, otherwise SignedString will call Sign to generate it.
func (t *Token) Sign() (string, error) {
	b, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	header := encodeSegment(b)
	b, err = json.Marshal(t.Claims)
	if err != nil {
		return "", err
	}
	claim := encodeSegment(b)
	return strings.Join([]string{header, claim}, "."), nil
}

// SignedString returns the token Signature. If the signature has not been generated
// yet, Sign will be called.
func (t *Token) SignedString(key any) (string, error) {
	str, err := t.Sign()
	if err != nil {
		return "", err
	}
	sig, err := t.Method.Sign(str, key)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{str, sig}, "."), nil
}

// encodeSegment encodes a JWT specific base64url encoding with padding stripped.
func encodeSegment(seg []byte) string {
	return base64.RawURLEncoding.EncodeToString(seg)
}

// decodeSegment decodes a JWT specific base64url encoding with padding stripped.
func decodeSegment(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(seg)
}

// Parse parses, validates, verifies the signature and returns the parsed token.
// keyFn will receive the parsed token and should return the cryptographic key
// for verifying the signature.
// The caller is strongly encouraged to set the WithValidMethods option to
// validate the 'alg' claim in the token matches the expected algorithm.
func Parse(token string, keyFn KeyFunc, options ParserOptions) (*Token, error) {
	return NewParser(options.ValidMethods, options.UseJSONNumber, options.SkipClaimsValidation).Parse(token, keyFn)
}

// ParseWithClaims is a shortcut for NewParser().ParseWithClaims().
//
// Note: If you provide a custom claim implementation that embeds one of the standard claims
// (such as RegisteredClaims), make sure that a: you either embed a non-pointer version of
// the claims or b: if you are using a pointer, allocate the proper memory for it before passing
// in the overall claims, otherwise you might run into a panic.
func ParseWithClaims(token string, claims Claims, keyFn KeyFunc, options ParserOptions) (*Token, error) {
	return NewParser(options.ValidMethods, options.UseJSONNumber, options.SkipClaimsValidation).ParseWithClaims(token, claims, keyFn)
}
