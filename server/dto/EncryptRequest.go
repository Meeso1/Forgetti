package dto

type EncryptRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
	KeyId string `json:"key_id" binding:"required,uuid"`
}
