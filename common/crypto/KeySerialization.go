package crypto

import (
	"encoding/base64"
	"encoding/json"
)

func SerializePublicKey(publicKey *PublicKey) (string, error) {
	publicKeyBytes, err := json.Marshal(publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(publicKeyBytes), nil
}

func DeserializePublicKey(serialized string) (*PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}

	var publicKey PublicKey
	err = json.Unmarshal(publicKeyBytes, &publicKey)
	if err != nil {
		return nil, err
	}

	return &publicKey, nil
}

func SerializePrivateKey(privateKey *PrivateKey) (string, error) {
	privateKeyBytes, err := json.Marshal(privateKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(privateKeyBytes), nil
}

func DeserializePrivateKey(serialized string) (*PrivateKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}

	var privateKey PrivateKey
	err = json.Unmarshal(privateKeyBytes, &privateKey)
	if err != nil {
		return nil, err
	}

	return &privateKey, nil
}
