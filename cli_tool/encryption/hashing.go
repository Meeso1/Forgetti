package encryption

import (
	"crypto/sha256"
	"encoding/base64"
)

// TODO: Things like this should be versioned
const remoteSalt string = "remote"
const localSalt string = "local"

func HashForRemotePart(key string) string {
	return hashWithSha256(key, remoteSalt)
}

func HashForLocalPart(key string) string {
	return hashWithSha256(key, localSalt)
}

func hashWithSha256(key string, salt string) (string) {
	hash := sha256.New()
	hash.Write([]byte(key))
	hash.Write([]byte(salt))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}