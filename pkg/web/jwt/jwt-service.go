package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"sync"
	"time"
)

const expirationTime = 30 * time.Minute

var jwtServiceOnce sync.Once

type JWTService struct {
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	expirationTime time.Time
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
	return &JWTService{
		privateKey:     privateKey,
		publicKey:      &privateKey.PublicKey,
		expirationTime: time.Now().Add(expirationTime),
	}
}

func (s *JWTService) GenerateToken(username, role string) *Token {
	return NewTokenWithClaims(
		SigningMethodRS256, MapClaims{
			"user": username,
			"role": role,
			"exp":  NewNumericDate(s.expirationTime),
		},
	)
}

func (s *JWTService) GenerateSignedToken(username, role string) string {
	str, err := s.GenerateToken(username, role).Sign()
	if err != nil {
		panic(err)
	}
	return str
}
