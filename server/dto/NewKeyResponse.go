package dto

import "time"

type Metadata struct {
	KeyId string `json:"key_id"`
	Expiration time.Time `json:"expiration"`
	VerificationKey string `json:"verification_key"`
}

type NewKeyResponse struct {
	EncryptedContent string `json:"encrypted_content"`
	Metadata Metadata `json:"metadata"`
}
