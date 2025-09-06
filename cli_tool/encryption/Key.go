package encryption

import "fmt"

type Key struct {
	LocalPart []byte  // 16 bytes
	RemotePart []byte // 16 bytes
}

func (k *Key) GetBytes() ([]byte, error) {
	if len(k.LocalPart) != 16 || len(k.RemotePart) != 16 {
		return nil, fmt.Errorf("local and remote parts must be 16 bytes")	
	}

	result := make([]byte, len(k.LocalPart) + len(k.RemotePart))
	result = append(result, k.LocalPart...)
	result = append(result, k.RemotePart...)

	return result, nil
}