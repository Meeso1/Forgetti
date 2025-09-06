package models

import (
	"forgetti-common/dto"
	"time"
	"fmt"
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

func (f *FileContentWithMetadata) String() string {
	return fmt.Sprintf("Encrypted content length: %d\n", len(f.FileContent)) +
		   fmt.Sprintf("Key ID: 				  %s\n", f.Metadata.KeyId) +
		   fmt.Sprintf("Expires at: 			  %s (in %s)\n", f.Metadata.Expiration.String(), time.Until(f.Metadata.Expiration).String()) +
		   fmt.Sprintf("Server Address: 		  %s\n", f.Metadata.ServerAddress)
}
