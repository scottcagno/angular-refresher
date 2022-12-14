package api

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const secretkey = "foobar"

func GenerateJWT(username, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretkey))

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

// func makeToken() {
// 	tok := jwt.NewWithClaims(
// 		jwt.SigningMethodHS256, jwt.MapClaims{
// 			"foo": "bar",
// 			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
// 		},
// 	)
// 	tokenString, err := tok.SignedString([]byte("foo"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	_ = tokenString
// }

// type JWT struct {
// 	RSAPrivateKey  *rsa.PrivateKey
// 	RSAPublicKey   *rsa.PublicKey
// 	ExpirationTime int
// }
//
// func (jwt *JWT) initKeys() {
// 	buf := new(bytes.Buffer)
// 	key, err := rsa.GenerateKey(buf, 2048)
// 	if err != nil {
// 		panic(err)
// 	}
// 	jwt.RSAPrivateKey = key
// 	jwt.RSAPublicKey = &key.PublicKey
// }
//
// func (jwt *JWT) generateToken(name, role string) string {
// 	return ""
// }
