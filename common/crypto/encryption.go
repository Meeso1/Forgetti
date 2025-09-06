package crypto

import (
	"encoding/base64"
	"fmt"
	"math/big"
)

const paddingByte byte = 7
const paddingSize int = 11

func Encrypt(content string, key *PublicKey) (string, error) {
	if err := ValidatePublicKey(key); err != nil {
		return "", fmt.Errorf("invalid public key: %w", err)
	}

	contentBytes := []byte(content)
	escaped := EscapeBytes(contentBytes, []byte{paddingByte})

	// Chunk into parts that are smaller (not equal) than N
	chunkSize := key.N.BitLen()/8 - paddingSize - 1
	result := make([]byte, 0, len(escaped)/chunkSize+1)
	for i := 0; i < len(escaped); i += chunkSize {
		end := min(i + chunkSize, len(escaped))

		encryptedChunk, err := encryptChunk(escaped[i:end], key)
		if err != nil {
			return "", fmt.Errorf("failed to encrypt chunk: %w", err)
		}

		result = append(result, encryptedChunk...)
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

func encryptChunk(chunk []byte, key *PublicKey) ([]byte, error) {
	chunkWithPadding := make([]byte, key.N.BitLen()/8-1) // We need to ensure that the chunk is smaller than N
	// Fill the chunk with padding bytes
	for i := 0; i < len(chunkWithPadding); i++ {
		chunkWithPadding[i] = paddingByte
	}
	
	// Copy the chunk at the start
	copy(chunkWithPadding, chunk)

	chunkAsInt := new(big.Int).SetBytes(chunkWithPadding)
	encryptedChunkAsInt := new(big.Int).Exp(chunkAsInt, key.E, key.N)

	encryptedBytes := make([]byte, key.N.BitLen()/8)
	if encryptedChunkAsInt.BitLen()/8 > key.N.BitLen()/8 {
		return nil, fmt.Errorf("encrypted chunk is too large: %d > %d", encryptedChunkAsInt.BitLen()/8, key.N.BitLen()/8)
	}

	encryptedChunkAsInt.FillBytes(encryptedBytes)
	return encryptedBytes, nil
}

// TODO: Implement some "signing" mechanism, and return error if the signature is not present in decrypted content
// This will prevent someone from guessing the content to decrypt, and subsequently guessing the exponent from public key
func Decrypt(encryptedContent string, key *PrivateKey) (string, error) {
	if err := ValidatePrivateKey(key); err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedContent)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted content: %w", err)
	}

	chunkSize := key.N.BitLen() / 8
	if chunkSize <= paddingSize {
		return "", fmt.Errorf("'N' is too small: %d bits (N) <= %d bytes (padding size)", 
			key.N.BitLen(), paddingSize * 8)
	}
	
	result := make([]byte, 0, len(encryptedBytes)/chunkSize+1)
	for i := 0; i < len(encryptedBytes); i += chunkSize {
		end := min(i + chunkSize, len(encryptedBytes))

		encryptedChunk, err := decryptChunk(encryptedBytes[i:end], key)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt chunk: %w", err)
		}

		result = append(result, encryptedChunk...)
	}

	result = removePadding(result)

	return string(result), nil
}

func decryptChunk(chunk []byte, key *PrivateKey) ([]byte, error) {
	chunkAsInt := new(big.Int).SetBytes(chunk)
	decryptedChunkAsInt := new(big.Int).Exp(chunkAsInt, key.D, key.N)

	if decryptedChunkAsInt.BitLen()/8 > key.N.BitLen()/8 {
		return nil, fmt.Errorf("decrypted chunk is too large: %d bits (decrypted chunk) > %d bits (N)", 
			decryptedChunkAsInt.BitLen(), key.N.BitLen())
	}

	decryptedBytes := make([]byte, key.N.BitLen()/8)
	decryptedChunkAsInt.FillBytes(decryptedBytes)

	return decryptedBytes[1:], nil
}

func removePadding(decryptedBytes []byte) []byte {
	result := make([]byte, 0, len(decryptedBytes))
	for i, b := range decryptedBytes {
		// Skip padding bytes
		if b == paddingByte {
			// Except if the previous byte is an escape byte
			if i != 0 && decryptedBytes[i-1] == EscapeByte {
				result = append(result, b)
			}
			continue
		}

		result = append(result, b)
	}

	return result
}
