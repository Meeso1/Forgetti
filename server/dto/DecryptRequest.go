package dto

type DecryptRequest struct {
	EncryptedContent string `json:"encrypted_content" binding:"required"`
	VerificationKey string `json:"verification_key" binding:"required"`
}
