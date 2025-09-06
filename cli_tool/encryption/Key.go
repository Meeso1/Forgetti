package encryption

import "fmt"

type Key struct {
	LocalPart  []byte // 16 bytes
	RemotePart []byte // 16 bytes
}

func CreateKey(localPart string, remotePart string) (*Key, error) {
	localPartBytes, err := HashLocalPart(localPart)
	if err != nil {
		return nil, err
	}

	remotePartBytes, err := HashEncryptedRemotePart(remotePart)
	if err != nil {
		return nil, err
	}

	return &Key{
		LocalPart:  localPartBytes,
		RemotePart: remotePartBytes,
	}, nil
}

func (k *Key) GetBytes() ([]byte, error) {
	if len(k.LocalPart) != 16 || len(k.RemotePart) != 16 {
		return nil, fmt.Errorf("local and remote parts must be 16 bytes")
	}

	result := make([]byte, len(k.LocalPart) + len(k.RemotePart))
	copy(result[:len(k.LocalPart)], k.LocalPart)
	copy(result[len(k.LocalPart):], k.RemotePart)

	return result, nil
}
