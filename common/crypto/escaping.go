package crypto

const EscapeByte byte = 92 // '\'

func EscapeBytes(bytes []byte, bytesToEscape []byte) []byte {
	bytesToEscape = prepareBytesToEscape(bytesToEscape)
	for _, bToEscape := range bytesToEscape {
		bytes = escapeSingleByte(bytes, bToEscape)
	}

	return bytes
}

func prepareBytesToEscape(bytes []byte) []byte {
	result := make(map[byte]bool)
	for _, b := range bytes {
		if b == EscapeByte {
			continue
		}
		result[b] = true
	}

	resultBytes := make([]byte, 0, len(result))
	// Escape byte must be escaped first
	resultBytes = append(resultBytes, EscapeByte)
	for b := range result {
		resultBytes = append(resultBytes, b)
	}

	return resultBytes
}

func escapeSingleByte(bytes []byte, byteToEscape byte) []byte {
	count := 0
	for _, b := range bytes {
		if b == byteToEscape {
			count++
		}
	}

	result := make([]byte, 0, len(bytes) + count)
	for _, b := range bytes {
		if b == byteToEscape {
			result = append(result, EscapeByte)
		}
		result = append(result, b)
	}
	return result
}

func UnescapeBytes(bytes []byte) []byte {
	// Count escape bytes, except when the escape byte is escaped (they would be counted twice)
	count := 0
	for i, b := range bytes[:len(bytes)-1] {
		if b == EscapeByte && bytes[i+1] != EscapeByte {
			count++
		}
	}

	result := make([]byte, 0, len(bytes) - count)
	for i, b := range bytes {
		// Skip escape bytes
		if b == EscapeByte {
			// Except if the next byte is also an escape byte
			// The last byte should never be a single escape byte, but if it is, we just keep it
			if i == len(bytes) - 1 || bytes[i+1] == EscapeByte {
				result = append(result, EscapeByte)
			}
			continue
		}

		result = append(result, b)
	}
	
	return result
}
