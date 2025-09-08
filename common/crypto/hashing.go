package crypto

import (
	"crypto/sha256"
	"fmt"
)

func HashToSize(key string, salt string, size int) ([]byte, error) {
	if size > 32 {
		return nil, fmt.Errorf("size must be less than or equal to 32: got %d", size)
	}

	hash := sha256.New()
	hash.Write([]byte(key))
	hash.Write([]byte(salt))
	return hash.Sum(nil)[:size], nil
}
