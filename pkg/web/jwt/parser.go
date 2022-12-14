package jwt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ParserOptions struct {
	ValidMethods         []string
	UseJSONNumber        bool
	SkipClaimsValidation bool
}

type Parser struct {
	validMethods         []string
	useJSONNumber        bool
	skipClaimsValidation bool
}

func NewParser(validMethods []string, useJSONNumber, skipClaimsValidation bool) *Parser {
	p := new(Parser)
	if validMethods != nil && len(validMethods) > 0 {
		p.validMethods = validMethods
	}
	if useJSONNumber {
		p.useJSONNumber = true
	}
	if skipClaimsValidation {
		p.skipClaimsValidation = true
	}
	return p
}

func (p *Parser) Parse(token string, keyFn KeyFunc) (*Token, error) {
	return p.ParseWithClaims(token, make(MapClaims), keyFn)
}

func (p *Parser) ParseWithClaims(token string, claims Claims, keyFn KeyFunc) (*Token, error) {

	jwt, parts, err := p.ParseUnverified(token, claims)
	if err != nil {
		return jwt, err
	}

	// Verify signing method is in the required set
	if p.validMethods != nil {
		var signingMethodValid = false
		var alg = jwt.Method.Alg()
		for _, m := range p.validMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			// signing method is not in the listed set
			return jwt, NewValidationError(
				fmt.Sprintf("signing method %v is invalid", alg),
				ValidationErrorSignatureInvalid,
			)
		}
	}

	// Lookup key
	var key any
	if keyFn == nil {
		// keyFunc was not provided, short-circuiting validation
		return jwt, NewValidationError("no KeyFunc was provided.", ValidationErrorUnverifiable)
	}
	key, err = keyFn(jwt)
	if err != nil {
		// keyFunc returned an error
		ve, ok := err.(*ValidationError)
		if ok {
			return jwt, ve
		}
		return jwt, &ValidationError{Inner: err, Errors: ValidationErrorUnverifiable}
	}

	vErr := new(ValidationError)

	// Validate Claims
	if !p.skipClaimsValidation {
		err = jwt.Claims.Valid()
		if err != nil {
			// If the Claims Valid returned an error, check if it is a validation error,
			// If it was another error type, create a ValidationError with a generic ClaimsInvalid flag set
			e, ok := err.(*ValidationError)
			if !ok {
				vErr = &ValidationError{Inner: err, Errors: ValidationErrorClaimsInvalid}
			} else {
				vErr = e
			}
		}
	}

	// Perform validation
	jwt.Signature = parts[2]
	err = jwt.Method.Verify(strings.Join(parts[0:2], "."), jwt.Signature, key)
	if err != nil {
		vErr.Inner = err
		vErr.Errors |= ValidationErrorSignatureInvalid
	}
	if vErr.valid() {
		jwt.Valid = true
		return jwt, nil
	}
	return jwt, vErr
}

func (p *Parser) ParseUnverified(token string, claims Claims) (*Token, []string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, parts, NewValidationError(
			"token contains an invalid number of segments",
			ValidationErrorMalformed,
		)
	}
	jwt := &Token{
		Raw: token,
	}
	// parse Header
	var headerBytes []byte
	headerBytes, err := decodeSegment(parts[0])
	if err != nil {
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			return jwt, parts, NewValidationError(
				"tokenstring should not contain 'bearer '",
				ValidationErrorMalformed,
			)
		}
		return jwt, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}
	err = json.Unmarshal(headerBytes, &jwt.Header)
	if err != nil {
		return jwt, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}

	// parse Claims
	var claimBytes []byte
	jwt.Claims = claims

	claimBytes, err = decodeSegment(parts[1])
	if err != nil {
		return jwt, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	if p.useJSONNumber {
		dec.UseNumber()
	}
	// JSON Decode.  Special case for map type to avoid weird pointer behavior
	c, ok := jwt.Claims.(MapClaims)
	if ok {
		err = dec.Decode(&c)
	} else {
		err = dec.Decode(&claims)
	}
	// Handle decode error
	if err != nil {
		return jwt, parts, &ValidationError{Inner: err, Errors: ValidationErrorMalformed}
	}

	// Lookup signature method
	method, ok := jwt.Header["alg"].(string)
	if ok {
		jwt.Method = GetSigningMethod(method)
		if jwt.Method == nil {
			return jwt, parts, NewValidationError(
				"signing method (alg) is unavailable.",
				ValidationErrorUnverifiable,
			)
		}
	} else {
		return jwt, parts, NewValidationError(
			"signing method (alg) is unspecified.",
			ValidationErrorUnverifiable,
		)
	}
	return jwt, parts, nil
}
