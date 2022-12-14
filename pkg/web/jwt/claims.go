package jwt

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"time"
)

type Claims interface {
	// VerifyAudience Compares the aud claim against cmp.
	// If req is false, it will return true, if aut is unset.
	VerifyAudience(cmp string, req bool) bool

	// VerifyExpiresAt compares the exp claim against cmp (cmp <= exp).
	// If req is false, it will return true, if exp is unset.
	VerifyExpiresAt(cmp int64, req bool) bool

	// VerifyIssuedAt compares the exp claim against cmp (cmp >= iat).
	// If req is false, it will return true, if iat is unset.
	VerifyIssuedAt(cmp int64, req bool) bool

	// VerifyNotBefore compares the nbf claim against cmp (cmp >= nbf).
	// If req is false, it will return true, if nbf is unset.
	VerifyNotBefore(cmp int64, req bool) bool

	// VerifyIssuer compares the iss claim against cmp.
	// If req is false, it will return true, if iss is unset.
	VerifyIssuer(cmp string, req bool) bool

	// Valid validates time based claims "exp, iat, nbf".
	// There is no accounting for clock skew.
	// As well, if any of the above claims are not in the token, it will
	// still be considered a valid claim.
	Valid() error
}

// RegisteredClaims are a structured version of the JWT Claims set which
// is restricted to Registered Claim Names as referenced in the documentation
// found at https://datatracker.ietf.org/doc/html/rfc7519#section-4.1. The
// RegisteredClaims type can be used on its own, but any additional private
// or public claims embedded in the Token will not be parsed. The default use
// case is to embed the RegisteredClaims type in a user-defined claim type.
type RegisteredClaims struct {
	// Issuer is the "iss" claim.
	Issuer string `json:"iss,omitempty"`

	// Subject is the "sub"  claim.
	Subject string `json:"sub,omitempty"`

	// Audience is the "aud"  claim.
	Audience ClaimStrings `json:"aud,omitempty"`

	// ExpiresAt is the "exp"  claim.
	ExpiresAt *NumericDate `json:"exp,omitempty"`

	// NotBefore is the "nbf" claim.
	NotBefore *NumericDate `json:"nbf,omitempty"`

	// IssuedAt is the "iat" claim.
	IssuedAt *NumericDate `json:"iat,omitempty"`

	// ID is the "jti" claim.
	ID string `json:"jti,omitempty"`
}

func (r RegisteredClaims) VerifyAudience(cmp string, req bool) bool {
	return verifyAud(r.Audience, cmp, req)
}

func (r RegisteredClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	if r.ExpiresAt == nil {
		return verifyExp(nil, time.Unix(cmp, 0), req)
	}
	return verifyExp(&r.ExpiresAt.Time, time.Unix(cmp, 0), req)
}

func (r RegisteredClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	if r.IssuedAt == nil {
		return verifyIat(nil, time.Unix(cmp, 0), req)
	}
	return verifyIat(&r.IssuedAt.Time, time.Unix(cmp, 0), req)
}

func (r RegisteredClaims) VerifyNotBefore(cmp int64, req bool) bool {
	if r.NotBefore == nil {
		return verifyNbf(nil, time.Unix(cmp, 0), req)
	}
	return verifyNbf(&r.NotBefore.Time, time.Unix(cmp, 0), req)
}

func (r RegisteredClaims) VerifyIssuer(cmp string, req bool) bool {
	return verifyIss(r.Issuer, cmp, req)
}

func (r RegisteredClaims) Valid() error {
	vErr := new(ValidationError)
	now := time.Now()
	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !r.VerifyExpiresAt(now.Unix(), false) {
		delta := now.Sub(r.ExpiresAt.Time)
		vErr.Inner = fmt.Errorf("%s by %s", ErrTokenExpired, delta)
		vErr.Errors |= ValidationErrorExpired
	}
	if !r.VerifyIssuedAt(now.Unix(), false) {
		vErr.Inner = ErrTokenUsedBeforeIssued
		vErr.Errors |= ValidationErrorIssuedAt
	}
	if !r.VerifyNotBefore(now.Unix(), false) {
		vErr.Inner = ErrTokenNotValidYet
		vErr.Errors |= ValidationErrorNotValidYet
	}
	if vErr.valid() {
		return nil
	}
	return vErr
}

type MapClaims map[string]any

func (m MapClaims) VerifyAudience(cmp string, req bool) bool {
	var aud []string
	switch v := m["aud"].(type) {
	case string:
		aud = append(aud, v)
	case []string:
		aud = v
	case []interface{}:
		for _, a := range v {
			vs, ok := a.(string)
			if !ok {
				return false
			}
			aud = append(aud, vs)
		}
	}
	return verifyAud(aud, cmp, req)
}

func (m MapClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	v, ok := m["exp"]
	if !ok {
		return !req
	}
	cmpTime := time.Unix(cmp, 0)
	switch exp := v.(type) {
	case float64:
		if exp == 0 {
			return verifyExp(nil, cmpTime, req)
		}
		return verifyExp(&newNumericDateFromSeconds(exp).Time, cmpTime, req)
	case json.Number:
		v, _ := exp.Float64()
		return verifyExp(&newNumericDateFromSeconds(v).Time, cmpTime, req)
	}
	return false
}

func (m MapClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	v, ok := m["iat"]
	if !ok {
		return !req
	}
	cmpTime := time.Unix(cmp, 0)
	switch iat := v.(type) {
	case float64:
		if iat == 0 {
			return verifyIat(nil, cmpTime, req)
		}
		return verifyIat(&newNumericDateFromSeconds(iat).Time, cmpTime, req)
	case json.Number:
		v, _ := iat.Float64()
		return verifyIat(&newNumericDateFromSeconds(v).Time, cmpTime, req)
	}
	return false
}

func (m MapClaims) VerifyNotBefore(cmp int64, req bool) bool {
	v, ok := m["nbf"]
	if !ok {
		return !req
	}
	cmpTime := time.Unix(cmp, 0)
	switch nbf := v.(type) {
	case float64:
		if nbf == 0 {
			return verifyNbf(nil, cmpTime, req)
		}
		return verifyNbf(&newNumericDateFromSeconds(nbf).Time, cmpTime, req)
	case json.Number:
		v, _ := nbf.Float64()

		return verifyNbf(&newNumericDateFromSeconds(v).Time, cmpTime, req)
	}

	return false
}

func (m MapClaims) VerifyIssuer(cmp string, req bool) bool {
	iss, _ := m["iss"].(string)
	return verifyIss(iss, cmp, req)
}

func (m MapClaims) Valid() error {
	vErr := new(ValidationError)
	now := time.Now().Unix()
	if !m.VerifyExpiresAt(now, false) {
		vErr.Inner = ErrTokenExpired
		vErr.Errors |= ValidationErrorExpired
	}
	if !m.VerifyIssuedAt(now, false) {
		vErr.Inner = ErrTokenUsedBeforeIssued
		vErr.Errors |= ValidationErrorIssuedAt
	}
	if !m.VerifyNotBefore(now, false) {
		vErr.Inner = ErrTokenNotValidYet
		vErr.Errors |= ValidationErrorNotValidYet
	}
	if vErr.valid() {
		return nil
	}

	return vErr
}

// The code found below are helpers for the different claims types
//

func verifyAud(aud []string, cmp string, required bool) bool {
	if len(aud) == 0 {
		return !required
	}
	// use a var here to keep constant time compare when looping over a number of claims
	result := false

	var stringClaims string
	for _, a := range aud {
		if subtle.ConstantTimeCompare([]byte(a), []byte(cmp)) != 0 {
			result = true
		}
		stringClaims = stringClaims + a
	}

	// case where "" is sent in one or many aud claims
	if len(stringClaims) == 0 {
		return !required
	}

	return result
}

func verifyExp(exp *time.Time, now time.Time, required bool) bool {
	if exp == nil {
		return !required
	}
	return now.Before(*exp)
}

func verifyIat(iat *time.Time, now time.Time, required bool) bool {
	if iat == nil {
		return !required
	}
	return now.After(*iat) || now.Equal(*iat)
}

func verifyNbf(nbf *time.Time, now time.Time, required bool) bool {
	if nbf == nil {
		return !required
	}
	return now.After(*nbf) || now.Equal(*nbf)
}

func verifyIss(iss string, cmp string, required bool) bool {
	if iss == "" {
		return !required
	}
	return subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0
}
