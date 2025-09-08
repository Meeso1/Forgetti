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
	AlgVersion      string 	  `json:"alg_version"`
}

type FileContentWithMetadata struct {
	FileContent []byte `json:"file_content"`
	Metadata  Metadata `json:"metadata"`
}

func ToFileMetadata(metadata dto.Metadata, serverAddress string) Metadata {
	return Metadata{
		KeyId: metadata.KeyId,
		Expiration: metadata.Expiration,
		VerificationKey: metadata.VerificationKey,
		ServerAddress: serverAddress,
		AlgVersion: CurrentAlgVersion().String(),
	}
}

func (f *FileContentWithMetadata) String() string {
	roundedDuration := time.Until(f.Metadata.Expiration).Round(time.Second)
	return fmt.Sprintf("Encrypted content length: %d bytes\n", len(f.FileContent)) +
		   fmt.Sprintf("Key ID:                   %s\n", f.Metadata.KeyId) +
		   fmt.Sprintf("Expires at:               %s (in %s)\n", f.Metadata.Expiration.String(), roundedDuration.String()) +
		   fmt.Sprintf("Server Address:           %s\n", f.Metadata.ServerAddress) +
		   fmt.Sprintf("Algorithm Version:        %s\n", f.Metadata.AlgVersion)
}
