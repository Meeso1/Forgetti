package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func SerializePublicKey(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(publicKeyBytes), nil
}

func DeserializePublicKey(serialized string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}

	deserialized, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	if typed, ok := deserialized.(*rsa.PublicKey); ok {
		return typed, nil
	}

	return nil, fmt.Errorf("deserialized public key is not an RSA public key")
}

func SerializePrivateKey(privateKey *rsa.PrivateKey) (string, error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(privateKeyBytes), nil
}

func DeserializePrivateKey(serialized string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	if typed, ok := privateKey.(*rsa.PrivateKey); ok {
		return typed, nil
	}

	return nil, fmt.Errorf("deserialized private key is not an RSA private key")
}
