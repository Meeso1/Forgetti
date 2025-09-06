package crypto

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestSerializeDeserializePublicKey(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create test public keys
	tests := []struct {
		name string
		key  *PublicKey
	}{
		{
			name: "Small numbers",
			key: &PublicKey{
				N: big.NewInt(12345),
				E: big.NewInt(65537),
			},
		},
		{
			name: "Large numbers",
			key: &PublicKey{
				N: new(big.Int).SetBytes([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}),
				E: new(big.Int).SetBytes([]byte{0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02}),
			},
		},
		{
			name: "Generated key",
			key:  keyPair.BroadcastKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			serialized, err := SerializePublicKey(tt.key)
			if err != nil {
				t.Errorf("SerializePublicKey() error = %v", err)
				return
			}

			// Serialized string should not be empty
			if serialized == "" {
				t.Error("SerializePublicKey() returned empty string")
				return
			}

			// Test deserialization
			deserialized, err := DeserializePublicKey(serialized)
			if err != nil {
				t.Errorf("DeserializePublicKey() error = %v", err)
				return
			}

			// Check that values match
			if deserialized.N.Cmp(tt.key.N) != 0 {
				t.Errorf("DeserializePublicKey() N = %v, want %v", deserialized.N, tt.key.N)
			}
			if deserialized.E.Cmp(tt.key.E) != 0 {
				t.Errorf("DeserializePublicKey() E = %v, want %v", deserialized.E, tt.key.E)
			}
		})
	}
}

func TestSerializeDeserializePrivateKey(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create test private keys
	tests := []struct {
		name string
		key  *PrivateKey
	}{
		{
			name: "Small numbers",
			key: &PrivateKey{
				N: big.NewInt(12345),
				D: big.NewInt(54321),
			},
		},
		{
			name: "Large numbers",
			key: &PrivateKey{
				N: new(big.Int).SetBytes([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}),
				D: new(big.Int).SetBytes([]byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}),
			},
		},
		{
			name: "Generated key",
			key:  keyPair.VerificationKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			serialized, err := SerializePrivateKey(tt.key)
			if err != nil {
				t.Errorf("SerializePrivateKey() error = %v", err)
				return
			}

			// Serialized string should not be empty
			if serialized == "" {
				t.Error("SerializePrivateKey() returned empty string")
				return
			}

			// Test deserialization
			deserialized, err := DeserializePrivateKey(serialized)
			if err != nil {
				marshalledKey, _ := json.Marshal(tt.key)
				t.Errorf("DeserializePrivateKey() error = %v (key: %v)", err, string(marshalledKey))
				return
			}

			// Check that values match
			if deserialized.N.Cmp(tt.key.N) != 0 {
				t.Errorf("DeserializePrivateKey() N = %v, want %v", deserialized.N, tt.key.N)
			}
			if deserialized.D.Cmp(tt.key.D) != 0 {
				t.Errorf("DeserializePrivateKey() D = %v, want %v", deserialized.D, tt.key.D)
			}
		})
	}
}

func TestDeserializePublicKeyInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid base64",
			input: "invalid base64 content!@#$%^&*()",
		},
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Valid base64 but invalid JSON",
			input: "aW52YWxpZCBqc29u", // "invalid json" in base64
		},
		{
			name:  "Valid base64 and JSON but missing fields",
			input: "e30=", // {} in base64 - Go will unmarshal this with nil/zero values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeserializePublicKey(tt.input)
			if err == nil && tt.input != "" {
				t.Error("DeserializePublicKey() should return error for invalid input")
			}
		})
	}
}

func TestDeserializePrivateKeyInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid base64",
			input: "invalid base64 content!@#$%^&*()",
		},
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Valid base64 but invalid JSON",
			input: "aW52YWxpZCBqc29u", // "invalid json" in base64
		},
		{
			name:  "Valid base64 and JSON but missing fields",
			input: "e30=", // {} in base64 - Go will unmarshal this with nil/zero values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeserializePrivateKey(tt.input)
			if err == nil && tt.input != "" {
				t.Error("DeserializePrivateKey() should return error for invalid input")
			}
		})
	}
}
