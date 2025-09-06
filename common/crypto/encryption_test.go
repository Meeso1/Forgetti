package crypto

import (
	"math/big"
	"strings"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}	

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "Simple text",
			content: "Hello, World!",
		},
		{
			name:    "Empty string",
			content: "",
		},
		{
			name:    "Single character",
			content: "A",
		},
		{
			name:    "Numbers and special characters",
			content: "1234567890!@#$%^&*()",
		},
		{
			name:    "Long text",
			content: strings.Repeat("This is a longer text to test encryption and decryption with multiple chunks. ", 10),
		},
		{
			name:    "Unicode characters",
			content: "Hello 世界! 🌍 Привет мир!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encryption
			encrypted, err := Encrypt(tt.content, keyPair.BroadcastKey)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
				return
			}

			// Encrypted string should not be empty (unless original content was empty)
			if tt.content != "" && encrypted == "" {
				t.Error("Encrypt() returned empty string for non-empty input")
				return
			}

			// Encrypted string should be different from original
			if tt.content != "" && encrypted == tt.content {
				t.Error("Encrypt() returned same string as input")
				return
			}

			// Test decryption
			decrypted, err := Decrypt(encrypted, keyPair.VerificationKey)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}

			// Decrypted content should match original
			if decrypted != tt.content {
				t.Errorf("Decrypt() = [%d]'%v', want [%d]'%v'", len(decrypted), decrypted, len(tt.content), tt.content)
				t.Errorf("Decrypted end: '%d %d %d %d %d %d %d %d %d %d'", decrypted[len(decrypted)-10], decrypted[len(decrypted)-9], decrypted[len(decrypted)-8], decrypted[len(decrypted)-7], decrypted[len(decrypted)-6], decrypted[len(decrypted)-5], decrypted[len(decrypted)-4], decrypted[len(decrypted)-3], decrypted[len(decrypted)-2], decrypted[len(decrypted)-1])
			}
		})
	}
}

func TestEncryptInvalidKey(t *testing.T) {
	// Test with invalid public key (nil values)
	validKeyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	tests := []struct {
		name string
		key *PublicKey
	}{
		{
			name: "Nil N",
			key: &PublicKey{
				N: nil,
				E: validKeyPair.BroadcastKey.E,
			},
		},
		{
			name: "Nil E",
			key: &PublicKey{
				N: validKeyPair.BroadcastKey.N,
				E: nil,
			},
		},
		{
			name: "Negative N",
			key: &PublicKey{
				N: new(big.Int).Neg(validKeyPair.BroadcastKey.N),
				E: validKeyPair.BroadcastKey.E,
			},
		},
		{
			name: "Negative E",
			key: &PublicKey{
				N: validKeyPair.BroadcastKey.N,
				E: new(big.Int).Neg(validKeyPair.BroadcastKey.E),
			},
		},
		{
			name: "Small N",
			key: &PublicKey{
				N: big.NewInt(1024),
				E: validKeyPair.BroadcastKey.E,
			},
		},
		{
			name: "Small E",
			key: &PublicKey{
				N: validKeyPair.BroadcastKey.N,
				E: big.NewInt(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Encrypt("test content", tt.key)
			if err == nil {
				t.Error("Encrypt() should return error for invalid key")
			}
		})
	}
}

func TestDecryptInvalidInput(t *testing.T) {
	// Generate a valid key pair
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "Invalid base64",
			content: "invalid base64 content!@#$%^&*()",
		},
		{
			name:    "Empty string",
			content: "",
		},
		{
			name:    "Random base64",
			content: "dGVzdCBjb250ZW50", // "test content" but not encrypted
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Decrypt() panicked for input '%v': %v", tt.content, r)
				}
			}()

			// We only care that it does not panic, error or not is fine
			_, _ = Decrypt(tt.content, keyPair.VerificationKey)
		})
	}
}

func TestDecryptInvalidKey(t *testing.T) {
	// Generate valid encrypted content
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	encrypted, err := Encrypt("test content", keyPair.BroadcastKey)
	if err != nil {
		t.Fatalf("Failed to encrypt test content: %v", err)
	}

	tests := []struct {
		name string
		key *PrivateKey
	}{
		{
			name: "Nil N",
			key: &PrivateKey{
				N: nil,
				D: keyPair.VerificationKey.D,
			},
		},
		{
			name: "Nil D",
			key: &PrivateKey{
				N: keyPair.VerificationKey.N,
				D: nil,
			},
		},
		{
			name: "Negative N",
			key: &PrivateKey{
				N: new(big.Int).Neg(keyPair.VerificationKey.N),
				D: keyPair.VerificationKey.D,
			},
		},
		{
			name: "Negative D",
			key: &PrivateKey{
				N: keyPair.VerificationKey.N,
				D: new(big.Int).Neg(keyPair.VerificationKey.D),
			},
		},
		{
			name: "Small N",
			key: &PrivateKey{
				N: big.NewInt(1024),
				D: keyPair.VerificationKey.D,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(encrypted, tt.key)
			if err == nil {
				t.Error("Decrypt() should return error for invalid key")
			}
		})
	}
}

func TestEncryptConsistency(t *testing.T) {
	// Generate a key pair
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	content := "Consistency test content"

	// Encrypt multiple times and ensure results are consistent
	encrypted1, err1 := Encrypt(content, keyPair.BroadcastKey)
	encrypted2, err2 := Encrypt(content, keyPair.BroadcastKey)

	if err1 != nil || err2 != nil {
		t.Fatalf("Encryption failed: err1=%v, err2=%v", err1, err2)
	}

	if encrypted1 != encrypted2 {
		t.Errorf("Encryption inconsistent: got %v and %v", encrypted1, encrypted2)
	}
}
