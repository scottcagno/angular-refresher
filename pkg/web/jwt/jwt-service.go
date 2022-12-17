package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"sync"
	"time"
)

const (
	expirationTime = 12 * time.Hour
	mySecretKey    = `eyJl2eHAiOjE12NzEyMzI30MjY3InJ3vbGU14iOiJST0xFX330FE5TUlO5I3iwid5X6Nlci7I6ImF`
)

var jwtServiceOnce sync.Once

type JWTService struct {
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	expirationTime time.Time
	signingMethod  SigningMethod
}

var JWTServiceInstance *JWTService

func NewJWTService() *JWTService {
	jwtServiceOnce.Do(
		func() {
			JWTServiceInstance = initJWTServiceInstance()
		},
	)
	return JWTServiceInstance
}

func initJWTServiceInstance() *JWTService {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	err = WriteRSAPrivateKeyAsPEM(privateKey, "cmd/roombooking/private_key.pem")
	if err != nil {
		panic(err)
	}
	err = WriteRSAPublicKeyAsPEM(&privateKey.PublicKey, "cmd/roombooking/public_key.pem")
	if err != nil {
		panic(err)
	}
	return &JWTService{
		privateKey:     privateKey,
		publicKey:      &privateKey.PublicKey,
		expirationTime: time.Now().Add(expirationTime),
		signingMethod:  SigningMethodRS256,
	}
}

func (s *JWTService) GenerateToken(username, role string) *Token {
	return NewTokenWithClaims(
		s.signingMethod, MapClaims{
			"user": username,
			"role": role,
			"exp":  NewNumericDate(s.expirationTime),
		},
	)
}

func (s *JWTService) GenerateSignedToken(username, role string) string {
	token := s.GenerateToken(username, role)
	str, err := token.SignedString(s.publicKey)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *JWTService) ValidateToken(tokenString string) (*Token, error) {
	parser := NewParser([]string{s.signingMethod.Alg()}, true, false)
	token, err := parser.Parse(
		tokenString, func(t *Token) (any, error) {
			// Validate the signing method
			if t.Method.Alg() != s.signingMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return s.publicKey, nil
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
