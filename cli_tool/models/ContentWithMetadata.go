package models

import (
	"forgetti-common/dto"
	"time"
)

type Metadata struct {
	KeyId 	 		string 	  `json:"key_id"`
	Expiration 		time.Time `json:"expiration"`
	VerificationKey string 	  `json:"verification_key"`
	ServerAddress   string 	  `json:"server_address"`
}

type FileContentWithMetadata struct {
	FileContent []byte `json:"file_content"`
	Metadata  Metadata `json:"metadata"`
}

func ToFileMetadata(metadata dto.Metadata, serverAddress string) *Metadata {
	return &Metadata{
		KeyId: metadata.KeyId,
		Expiration: metadata.Expiration,
		VerificationKey: metadata.VerificationKey,
		ServerAddress: serverAddress,
	}
}
