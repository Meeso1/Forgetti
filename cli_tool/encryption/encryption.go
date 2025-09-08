package encryption

import "forgetti-common/crypto"

func Encrypt(content []byte, key *Key) ([]byte, error) {
	keyBytes, err := key.GetBytes()
	if err != nil {
		return nil, err
	}

	return crypto.EncryptAes256(content, keyBytes)
}


func Decrypt(content []byte, key *Key) ([]byte, error) {
	keyBytes, err := key.GetBytes()
	if err != nil {
		return nil, err
	}

	return crypto.DecryptAes256(content, keyBytes)
}
