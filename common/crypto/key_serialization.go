package crypto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

const version1Signature = "v1"
const separator = ":"

func SerializePublicKey(publicKey *PublicKey) (string, error) {
	publicKeyBytes, err := json.Marshal(publicKey)
	if err != nil {
		return "", err
	}

	return addVersionSignature(base64.StdEncoding.EncodeToString(publicKeyBytes)), nil
}

func DeserializePublicKey(serialized string) (*PublicKey, error) {
	version, serialized := consumeVersionSignature(serialized)
	if version != version1Signature {
		return nil, fmt.Errorf("invalid public key: unsupported version signature '%s'", version)
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}

	var publicKey PublicKey
	err = json.Unmarshal(publicKeyBytes, &publicKey)
	if err != nil {
		return nil, err
	}

	if publicKey.N == nil || publicKey.E == nil {
		return nil, fmt.Errorf("invalid public key: serialized JSON is missing required fields")
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

	if privateKey.N == nil || privateKey.D == nil {
		return nil, fmt.Errorf("invalid private key: serialized JSON is missing required fields")
	}

	return &privateKey, nil
}

func addVersionSignature(serialized string) string {
	return version1Signature + separator + serialized
}

func consumeVersionSignature(serialized string) (string, string) {
	first, second, found := strings.Cut(serialized, separator)
	if !found {
		return "", first
	}

	return first, second
}
