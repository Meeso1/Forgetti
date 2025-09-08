package encryption

import (
	"Forgetti/models"
	"fmt"
	"forgetti-common/logging"
)

type Key struct {
	LocalPart  []byte // 16 bytes
	RemotePart []byte // 16 bytes
}

func CreateKey(localPart string, remotePart string, version models.AlgVersion) (*Key, error) {
	logger := logging.MakeLogger("encryption.CreateKey")
	logger.Verbose("Creating symmetric key with algorithm version: %s", version.String())

	logger.Verbose("Hashing local part with algorithm: %s", version.LocalHash)
	localPartBytes, err := HashLocalPart(localPart, version.LocalHash)
	if err != nil {
		logger.Error("Failed to hash local part: %v", err)
		return nil, err
	}
	logger.Verbose("Local part hashed successfully (%d bytes)", len(localPartBytes))

	logger.Verbose("Hashing remote part with algorithm: %s", version.PostRemoteHash)
	remotePartBytes, err := HashEncryptedRemotePart(remotePart, version.PostRemoteHash)
	if err != nil {
		logger.Error("Failed to hash encrypted remote part: %v", err)
		return nil, err
	}
	logger.Verbose("Remote part hashed successfully (%d bytes)", len(remotePartBytes))

	key := &Key{
		LocalPart:  localPartBytes,
		RemotePart: remotePartBytes,
	}
	logger.Info("Successfully created symmetric key")
	return key, nil
}

func (k *Key) GetBytes() ([]byte, error) {
	logger := logging.MakeLogger("encryption.Key.GetBytes")
	logger.Verbose("Converting key to bytes (local: %d bytes, remote: %d bytes)", len(k.LocalPart), len(k.RemotePart))

	if len(k.LocalPart) != 16 || len(k.RemotePart) != 16 {
		logger.Error("Invalid key part sizes - local: %d bytes, remote: %d bytes (expected 16 each)", len(k.LocalPart), len(k.RemotePart))
		return nil, fmt.Errorf("local and remote parts must be 16 bytes")
	}

	result := make([]byte, len(k.LocalPart)+len(k.RemotePart))
	copy(result[:len(k.LocalPart)], k.LocalPart)
	copy(result[len(k.LocalPart):], k.RemotePart)

	logger.Verbose("Successfully converted key to %d bytes", len(result))
	return result, nil
}
