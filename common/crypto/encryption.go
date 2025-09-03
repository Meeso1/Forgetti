package crypto

import (
	"encoding/base64"
	"fmt"
	"math/big"
)

const paddingByte byte = 7
const paddingSize int = 11

func Encrypt(content string, key *PublicKey) (string, error) {
	contentBytes := []byte(content)

	// Chunk into parts that are significantly smaller than N
	chunkSize := key.N.BitLen()/8 - paddingSize - 1
	result := make([]byte, 0, len(contentBytes)/chunkSize+1)
	for i := 0; i < len(contentBytes); i += chunkSize {
		end := i + chunkSize
		if end > len(contentBytes) {
			end = len(contentBytes)
		}

		encryptedChunk, err := EncryptChunk(contentBytes[i:end], key)
		if err != nil {
			return "", fmt.Errorf("failed to encrypt chunk: %w", err)
		}

		result = append(result, encryptedChunk...)
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

func EncryptChunk(chunk []byte, key *PublicKey) ([]byte, error) {
	chunkWithPadding := make([]byte, key.N.BitLen()/8-1) // We need to ensure that the chunk is smaller than N
	// Copy the chunk at the start
	copy(chunkWithPadding, chunk)
	// Fill the rest with 0s
	for i := len(chunk); i < len(chunkWithPadding); i++ {
		chunkWithPadding[i] = 0
	}
	// Fill the reserved space with padding bytes
	for i := 0; i < paddingSize; i++ {
		chunkWithPadding[len(chunkWithPadding)-i-1] = paddingByte
	}

	chunkAsInt := new(big.Int).SetBytes(chunkWithPadding)
	encryptedChunkAsInt := new(big.Int).Exp(chunkAsInt, key.E, key.N)

	encryptedBytes := make([]byte, key.N.BitLen()/8)
	if encryptedChunkAsInt.BitLen()/8 > key.N.BitLen()/8 {
		return nil, fmt.Errorf("encrypted chunk is too large: %d > %d", encryptedChunkAsInt.BitLen()/8, key.N.BitLen()/8)
	}

	encryptedChunkAsInt.FillBytes(encryptedBytes)
	return encryptedBytes, nil
}

func Decrypt(encryptedContent string, key *PrivateKey) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedContent)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted content: %w", err)
	}

	chunkSize := key.N.BitLen() / 8
	result := make([]byte, 0, len(encryptedBytes)/chunkSize+1)
	for i := 0; i < len(encryptedBytes); i += chunkSize {
		end := i + chunkSize
		if end > len(encryptedBytes) {
			end = len(encryptedBytes)
		}

		encryptedChunk, err := DecryptChunk(encryptedBytes[i:end], key)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt chunk: %w", err)
		}

		result = append(result, encryptedChunk...)
	}

	return string(result), nil
}

func DecryptChunk(chunk []byte, key *PrivateKey) ([]byte, error) {
	chunkAsInt := new(big.Int).SetBytes(chunk)
	decryptedChunkAsInt := new(big.Int).Exp(chunkAsInt, key.D, key.N)

	decryptedBytes := make([]byte, key.N.BitLen()/8 - 1)
	decryptedChunkAsInt.FillBytes(decryptedBytes)

	// Strip padding
	endOfData := len(decryptedBytes) - paddingSize

	//TODO: Use escaping instead of stripping trailing 0s
	if decryptedBytes[endOfData-1] == byte(0) {
		for endOfData >= 0 && decryptedBytes[endOfData-1] == byte(0) {
			endOfData--
		}
	}

	return decryptedBytes[:endOfData], nil
}
