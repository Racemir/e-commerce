package auth

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT imzalama anahtarı – ortam değişkeninden okunur.
var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateSecureToken, kriptografik olarak güvenli, rastgele bir metin (token) üretir.
func GenerateSecureToken() (string, error) {
	emptyMemory := make([]byte, 32)

	// Boş belleği işletim sisteminin kripto motorundan gelen rastgele verilerle doldur
	_, readError := rand.Read(emptyMemory)
	if readError != nil {
		return "", readError
	}
	// Byte (sayı) verisini, Hex (metin) formatına çevir
	return hex.EncodeToString(emptyMemory), nil
}

// Jwt token oluşturma
func JwtCreateToken(userID int, email string, role string) (string, error) {
	// Jwt Token içine koyacağımız verileri tutacak yapı.
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	// HS256 algoritması ile token oluşturuyorum. (HMAC)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// JWT Token imzalama
	tokenString, signStringError := token.SignedString(secretKey)
	if signStringError != nil {
		return "", signStringError
	}
	return tokenString, nil

}

// JwtVerifyToken, gelen token stringini doğrular ve içindeki verileri (claims) döndürür.
func JwtVerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Algoritmanın HMAC olduğunu doğrula
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
