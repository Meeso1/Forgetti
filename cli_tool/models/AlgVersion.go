package models

import (
	"fmt"
	"strings"
)

const separator = ":"

type AlgVersion struct {
	Symmetric string
	LocalHash string
	PreRemoteHash string
	PostRemoteHash string
}

func CurrentAlgVersion() AlgVersion {
	return AlgVersion{Symmetric: "1", LocalHash: "1", PreRemoteHash: "1", PostRemoteHash: "1"}
}

func (v AlgVersion) String() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", v.Symmetric, separator, v.LocalHash, separator, v.PreRemoteHash, separator, v.PostRemoteHash)
}

func ParseAlgVersion(s string) AlgVersion {
	parts := strings.Split(s, separator)

	result := AlgVersion{Symmetric: parts[0]}
	if len(parts) > 1 {
		result.LocalHash = parts[1]
	} else {
		result.LocalHash = ""
	}

	if len(parts) > 2 {
		result.PreRemoteHash = parts[2]
	} else {
		result.PreRemoteHash = ""
	}

	if len(parts) > 3 {
		result.PostRemoteHash = parts[3]
	} else {
		result.PostRemoteHash = ""
	}

	return result
}
