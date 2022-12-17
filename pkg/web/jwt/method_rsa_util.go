package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	return priv, &priv.PublicKey
}

func WriteRSAPrivateKey(key *rsa.PrivateKey, file string) error {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, keyBytes, 0655)
	if err != nil {
		return err
	}
	return nil
}

func WriteRSAPublicKey(key *rsa.PublicKey, file string) error {
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, keyBytes, 0655)
	if err != nil {
		return err
	}
	return nil
}

func ReadRSAPrivateKey(file string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key type is not RSA")
	}
	return rsaKey, nil
}

func ReadRSAPublicKey(file string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key type is not RSA")
	}
	return rsaKey, nil
}

func WriteRSAPrivateKeyAsPEM(key *rsa.PrivateKey, file string) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}
	err := WritePEM(block, file)
	if err != nil {
		return err
	}
	return nil
}

func ReadRSAPrivateKeyFromPEM(file string) (*rsa.PrivateKey, error) {
	block, err := ReadPEM(file)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key type is not RSA")
	}
	return rsaKey, nil
}

func ReadPEM(file string) (*pem.Block, error) {
	pemBytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing keys")
	}
	return block, nil
}

func WritePEM(block *pem.Block, file string) error {
	pemBytes := pem.EncodeToMemory(block)
	err := os.WriteFile(file, pemBytes, 0655)
	if err != nil {
		return err
	}
	return nil
}

func WriteRSAPublicKeyAsPEM(key *rsa.PublicKey, file string) error {
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	}
	err = WritePEM(block, file)
	if err != nil {
		return err
	}
	return nil
}

func ReadRSAPublicKeyFromPEM(file string) (*rsa.PublicKey, error) {
	block, err := ReadPEM(file)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("key type is not RSA")
	}
	return rsaKey, nil
}
