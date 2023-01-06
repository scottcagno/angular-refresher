package jwt

import (
	"fmt"
	"testing"
	"time"
)

func TestJWTService_GenerateSignedToken(t *testing.T) {

	signingMethod := SigningMethodRS256

	token := NewTokenWithClaims(
		signingMethod, MapClaims{
			"user": "admin",
			"pass": "secret",
			"role": "ROLE_ADMIN",
			"exp":  NewNumericDate(time.Now().Add(12 * time.Hour)),
		},
	)

	key, err := RSAPrivateKeyFromString(privateKey)
	if err != nil {
		t.Error(err)
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("tokenStr: %q\n", tokenString)
}

func TestJWTService_ValidateTokenString(t *testing.T) {

	signingMethod := SigningMethodRS256

	parser := NewParser([]string{signingMethod.Alg()}, true, false)

	keyFn := func(t *Token) (any, error) {
		// Validate the signing method
		if t.Method.Alg() != signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		key, err := RSAPublicKeyFromString(publicKey)
		if err != nil {
			panic(err)
		}
		return key, nil
	}

	token, err := parser.Parse(jwtTokenString, keyFn)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("token: %+v\nvalid: %v\n", token, token.Valid)

}

const jwtTokenString = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiO` +
	`jE2NzMwODI1NTksInBhc3MiOiJzZWNyZXQiLCJyb2xlIjoiUk9MRV9BRE1JTiIsInVzZXIiOiJhZG1pbiJ9.` +
	`RXQfaLSGE01pU5EwG1VyQpUWZz5uVKdL9kOR8S06lRQ5KZAmZQBHsnrwy--yCiHvUH8F8YlTYPDXa8fXe37t` +
	`5d2MviVU1Mk09oL8qU0ABO8XMuYITXU7kBkV4gfTIeZpC44JbMpHU_NNui9fXTAEfriPQT2TUiFIcDc0N-Wf` +
	`yfcs72vVNAQkr4bmo26BZOMuGAdDmGhvsyRLKyxRvBh89e0NFeHze_uzErGGY3lc2nzp7QoPVWn3DFBqk_qG` +
	`aY0ahNDL_MwLBCisYa8vOQUOKT3EDt-5nxXA9dWVGFHz012kKRMYXcRrkYd1l7AAUEMTQa0jpPrKTMnxXDV5qbJZHw`

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAqN3Qme44/bArcJTrYHr61Fyxy7+ksfuB4QuQOw1KxXDJb3Y/
lx/UygJKdGwTa4PJKh/0yQ4AfXlqZRQBD9pZI5KL54Ho/HlHmiKuM3OfzrfFna96
tUiBCj4GddwOCpC8zWBix1/8e5/tehyouvNRXFZv/J2TTEq0KpjuIcL0YPonVGqt
FnmUbrAcx2dBDPOyMuI6rINxS2Tl7guLoNNTVK9/HjaEnx/EPWsoyaCplQvirQp+
W0p0zIXj5ulCbi2kP9K1d0wJdYbBgy+H62WXGPYeWbG6zN7I/FqxN1sA+JPc5yJ5
UjpLgPoLdB24WbHm11iMpXDPvQcGEJADEvbG4QIDAQABAoIBAGqh5x/MtmA75rJo
11lNTybagctPcQiS8SjSzHj9o8GZvxmLagxcJVqKp8lslbbGuTjIhSKQnO6exPie
8Sy5jKBR4daGykDjtLs4OxhyYu9+TGOOc8YVyqZVFG9ITfWOACsQOk/75MuL6cG6
ZzHmg/HzAzRZrLH4WlyrxXO8T+UMCsiHAfyodCnQiwjPGCAXAI0j0yI4mIEl5SzW
Eb1rgxt9WXZWfb4LVh5edWJLYwjro4b6ilXrCLiBzCr2lXPkW6hD1mA/xpIFXG7d
cGgiIueYBSL75IXPy3bgVXdQsOsQV0dqGI6uDYaN98UAx3+LUwuXBBYA2OZ5K1kJ
Z1g83gECgYEA3OBlCqupg0isZVekmkIaWEkUcjHq8yWURZrlVj7WmtYIlKU7wXSH
Xw00G3PC1DASqnNFqdEYaiZFHt0BTagyt790007k9QPvbfSAWAQRqGSrWg0ZZJjT
mmYk20TsDcizH4rgUBsKYtglO0p5qpTpiuJGNhPHXDRRuzzGkjfzXdkCgYEAw7gn
0J24dM2AWJwWk7H2TC34ifBBIUxG3OGBX6xVskP07WtFkAR7qqiXIlBFCDERsERF
sCnb8P4tV3GEL7KrM6+IGthxIhdZxfJj+GS2ykVF40KK7bli6nNMRrtQwtB1MDkQ
bQ72VB5LJg4SoYl6mRTUewKpT0OC87capdjCpEkCgYEAkeRF2T55wRWHiYjSWHHB
JP9gWe1O2zu/LBqb0NPAvJUTJdveFHH72HTILjnQPodiTOPG59wM3FBa53/jFIA8
v9HeQJSj9pKa223cOEa3wxp7dAei9idb3WgKgCqOIKyoY/U/JKo3ugI61Wbj5iBm
Ai5jYeS+kdCdC6ehIYODZEECgYEAwXAQgeJwfZjiQjG7KqyYNoC1BXgclxFxdDu5
B1sns7HwsHr2XLnhlDFedn6JS+hbiDBiBBPLGqvNGoDKWe4nwUS6q3XCkyQrCTZh
Ug4Qj2faBBwvfXdd6USdXccisfkf6dJshq1kDo2GTo1YIqnjLstkmlNuDTY3hjMx
tjq/XWkCgYEAkOl8IQSy51i5/gnMPTfUne0Cbxa8IL6lT87LWzEYYyqLzTqZQzdd
W9BDV+dXv5elzZZk/aE8ZqKwBDaoO7FDasQ5qZDeYhi2nvld4u2EWlnAK41hmm28
quRRckhCfiLNMKIs+Ghb/XjYMDQDnmYR83C9no0bR7Px4bHxhDIvcoU=
-----END RSA PRIVATE KEY-----`

const publicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqN3Qme44/bArcJTrYHr6
1Fyxy7+ksfuB4QuQOw1KxXDJb3Y/lx/UygJKdGwTa4PJKh/0yQ4AfXlqZRQBD9pZ
I5KL54Ho/HlHmiKuM3OfzrfFna96tUiBCj4GddwOCpC8zWBix1/8e5/tehyouvNR
XFZv/J2TTEq0KpjuIcL0YPonVGqtFnmUbrAcx2dBDPOyMuI6rINxS2Tl7guLoNNT
VK9/HjaEnx/EPWsoyaCplQvirQp+W0p0zIXj5ulCbi2kP9K1d0wJdYbBgy+H62WX
GPYeWbG6zN7I/FqxN1sA+JPc5yJ5UjpLgPoLdB24WbHm11iMpXDPvQcGEJADEvbG
4QIDAQAB
-----END RSA PUBLIC KEY-----`
