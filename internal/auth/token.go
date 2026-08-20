package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken, kriptografik olarak güvenli, rastgele bir metin (token) üretir.
func GenerateSecureToken() (string, error) {
	emptyMemory := make([]byte, 32)

	// Boş belleği işletim sisteminin kripto motorundan gelen rastgele verilerle doldur
	_, readError := rand.Read(emptyMemory)
	if readError != nil {
		return "", nil
	}
	// Byte (sayı) verisini, Hex (metin) formatına çevir
	return hex.EncodeToString(emptyMemory), nil
}
