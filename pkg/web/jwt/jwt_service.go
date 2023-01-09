package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
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

func NewJWTService(privateKeyFile, publicKeyFile string) *JWTService {
	jwtServiceOnce.Do(
		func() {
			service, err := initJWTServiceInstance(privateKeyFile, publicKeyFile)
			if err != nil {
				panic(err)
			}
			JWTServiceInstance = service
		},
	)
	return JWTServiceInstance
}

func initJWTServiceInstance(privateKeyFile, publicKeyFile string) (*JWTService, error) {
	// Initialize our variables
	var key *rsa.PrivateKey
	// Check to see if the specified private and public key files exist
	_, err := os.Stat(privateKeyFile)
	if os.IsNotExist(err) {
		// The private key file does not exist, so the public should not
		// exist either, but if it does, we should remove it.
		_, err = os.Stat(publicKeyFile)
		if os.IsExist(err) {
			// The public key file does indeed exist; let us blow it away
			err = os.RemoveAll(publicKeyFile)
			if err != nil {
				return nil, err
			}
		}
		// Now, since we do not have any private and public keys we must create them.
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		// We should write our *rsa.PrivateKey to a file
		err = WriteRSAPrivateKeyAsPEM(key, privateKeyFile)
		if err != nil {
			return nil, err
		}
		// And write our *rsa.PublicKey to a file
		err = WriteRSAPublicKeyAsPEM(&key.PublicKey, publicKeyFile)
		if err != nil {
			return nil, err
		}
	} else {
		// We already have private and public key files, so we can read them
		key, err = ReadRSAPrivateKeyFromPEM(privateKeyFile)
		if err != nil {
			return nil, err
		}
		// Next, read our public key file
		var pub *rsa.PublicKey
		pub, err = ReadRSAPublicKeyFromPEM(publicKeyFile)
		if err != nil {
			return nil, err
		}
		key.PublicKey = *pub
	}
	// Finally, we can assemble and return our JWTService
	return &JWTService{
		privateKey:     key,
		publicKey:      &key.PublicKey,
		expirationTime: time.Now().Add(expirationTime),
		signingMethod:  SigningMethodRS256,
	}, nil
}

func (s *JWTService) generateToken(username, password, role string) *Token {
	return NewTokenWithClaims(
		s.signingMethod, MapClaims{
			"user": username,
			"pass": password,
			"role": role,
			"exp":  NewNumericDate(s.expirationTime),
		},
	)
}

func (s *JWTService) GenerateSignedToken(username, password, role string) string {
	token := s.generateToken(username, password, role)
	// SignedString must take a *PrivateKey
	str, err := token.SignedString(s.privateKey)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *JWTService) ValidateTokenString(tokenString string) (*Token, error) {
	parser := NewParser([]string{s.signingMethod.Alg()}, true, false)

	// keyFn must return a *PublicKey type when parsing
	keyFn := func(t *Token) (any, error) {
		// Validate the signing method
		if t.Method.Alg() != s.signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.publicKey, nil
	}

	token, err := parser.Parse(tokenString, keyFn)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(MapClaims); ok && token.Valid {
		// you could check some claims in here if you wanted...
		_ = claims
	}
	return token, nil
}

func (s *JWTService) ExpireTime() time.Time {
	return s.expirationTime
}
