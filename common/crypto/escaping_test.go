package crypto

import (
	"fmt"
	"testing"
)

func TestEscapeBytes(t *testing.T) {
	tests := []struct {
		name          string
		bytes         []byte
		bytesToEscape []byte
		expected      []byte
	}{
		{
			name:          "empty bytes",
			bytes:         []byte{},
			bytesToEscape: []byte{1, 2, 3},
			expected:      []byte{},
		},
		{
			name:          "no bytes to escape",
			bytes:         []byte{1, 2, 3, 4},
			bytesToEscape: []byte{5, 6},
			expected:      []byte{1, 2, 3, 4},
		},
		{
			name:          "single byte to escape",
			bytes:         []byte{1, 2, 1, 3},
			bytesToEscape: []byte{1},
			expected:      []byte{EscapeByte, 1, 2, EscapeByte, 1, 3},
		},
		{
			name:          "multiple bytes to escape",
			bytes:         []byte{1, 2, 3, 1, 2},
			bytesToEscape: []byte{1, 2},
			expected:      []byte{EscapeByte, 1, EscapeByte, 2, 3, EscapeByte, 1, EscapeByte, 2},
		},
		{
			name:          "escape byte itself needs escaping",
			bytes:         []byte{EscapeByte, 1, EscapeByte, 2},
			bytesToEscape: []byte{1, EscapeByte},
			expected:      []byte{EscapeByte, EscapeByte, EscapeByte, 1, EscapeByte, EscapeByte, 2},
		},
		{
			name:          "escape byte in bytesToEscape is filtered out",
			bytes:         []byte{1, 2, 3},
			bytesToEscape: []byte{1, EscapeByte, 2},
			expected:      []byte{EscapeByte, 1, EscapeByte, 2, 3},
		},
		{
			name:          "duplicate bytes in bytesToEscape are deduplicated",
			bytes:         []byte{1, 2, 1, 2},
			bytesToEscape: []byte{1, 1, 2, 2},
			expected:      []byte{EscapeByte, 1, EscapeByte, 2, EscapeByte, 1, EscapeByte, 2},
		},
		{
			name:          "all bytes need escaping",
			bytes:         []byte{1, 1, 1},
			bytesToEscape: []byte{1},
			expected:      []byte{EscapeByte, 1, EscapeByte, 1, EscapeByte, 1},
		},
		{
			name:          "only escape byte in input",
			bytes:         []byte{EscapeByte},
			bytesToEscape: []byte{EscapeByte},
			expected:      []byte{EscapeByte, EscapeByte},
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("EscapeBytes() panicked: %v", r)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeBytes(tt.bytes, tt.bytesToEscape)
			shouldEqual(t, result, tt.expected)
		})
	}
}

func TestUnescapeBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected []byte
	}{
		{
			name:     "empty bytes",
			bytes:    []byte{},
			expected: []byte{},
		},
		{
			name:     "no escape bytes",
			bytes:    []byte{1, 2, 3, 4},
			expected: []byte{1, 2, 3, 4},
		},
		{
			name:     "single escaped byte",
			bytes:    []byte{EscapeByte, 1, 2, 3},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "multiple escaped bytes",
			bytes:    []byte{EscapeByte, 1, EscapeByte, 2, EscapeByte, 3},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "escaped escape byte (double escape)",
			bytes:    []byte{EscapeByte, EscapeByte, 1, EscapeByte, EscapeByte, 2},
			expected: []byte{EscapeByte, 1, EscapeByte, 2},
		},
		{
			name:     "mixed escaped and unescaped",
			bytes:    []byte{1, EscapeByte, 2, 3, EscapeByte, 4},
			expected: []byte{1, 2, 3, 4},
		},
		{
			name:     "escape byte at start",
			bytes:    []byte{EscapeByte, 1, 2, 3},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "escape byte at end",
			bytes:    []byte{1, 2, 3, EscapeByte, 4},
			expected: []byte{1, 2, 3, 4},
		},
		{
			name:     "escape byte as last byte (should be kept)",
			bytes:    []byte{1, 2, EscapeByte},
			expected: []byte{1, 2, EscapeByte},
		},
		{
			name:     "double escape at end",
			bytes:    []byte{1, 2, EscapeByte, EscapeByte},
			expected: []byte{1, 2, EscapeByte},
		},
		{
			name:     "triple escape sequence",
			bytes:    []byte{EscapeByte, EscapeByte, EscapeByte, 1},
			expected: []byte{EscapeByte, 1},
		},
		{
			name:     "only escape bytes",
			bytes:    []byte{EscapeByte, EscapeByte, EscapeByte, EscapeByte},
			expected: []byte{EscapeByte, EscapeByte},
		},
		{
			name:     "complex mixed sequence",
			bytes:    []byte{1, EscapeByte, 2, EscapeByte, EscapeByte, 3, EscapeByte, 4, EscapeByte, EscapeByte},
			expected: []byte{1, 2, EscapeByte, 3, 4, EscapeByte},
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("UnescapeBytes() panicked: %v", r)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnescapeBytes(tt.bytes)
			shouldEqual(t, result, tt.expected)
		})
	}
}

func shouldEqual(t *testing.T, actual, expected []byte) {
	for i, b := range actual {
		if i >= len(expected) {
			t.Errorf("Too many bytes: %v > %v", len(actual), len(expected))
			return
		}

		if b != expected[i] {
			t.Errorf("Bytes at index %v differ: %v != %v", i, getContext(actual, i), getContext(expected, i))
			return
		}
	}
}

func getContext(bytes []byte, position int) string {
	if position < 0 || position >= len(bytes) {
		return ""
	}

	preceding := bytes[:position]
	atStart := len(preceding) <= 3
	if !atStart {
		preceding = preceding[len(preceding)-3:]
	}

	following := bytes[position+1:]
	atEnd := len(following) <= 3
	if !atEnd {
		following = following[:3]
	}

	result := "["
	if !atStart {
		result += "..."
	}
	for _, b := range preceding {
		result += fmt.Sprintf(" %v", b)
	}
	result += fmt.Sprintf(" (%v)", bytes[position])
	for _, b := range following {
		result += fmt.Sprintf(" %v", b)
	}
	if !atEnd {
		result += " ..."
	}
	return result + " ]"
}
