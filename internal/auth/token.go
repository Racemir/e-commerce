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
// E-posta doğrulama ve şifre sıfırlama tokenları için kullanılır.
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

// JwtCreateAccessToken, kısa ömürlü (15 dakika) bir Access Token üretir.
// Her API isteğinde "access_token" cookie'si ile taşınır.
func JwtCreateAccessToken(userID int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// JwtCreateRefreshToken, uzun ömürlü (7 gün) bir Refresh Token üretir.
// İçine benzersiz bir tokenID gömülür — Rotation sırasında bu ID ile Redis'teki kayıt eşleştirilir.
// Yalnızca /api/auth/refresh endpoint'ine "refresh_token" cookie'si ile gider.
func JwtCreateRefreshToken(userID int, tokenID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"token_id": tokenID,
		"type":     "refresh",
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// JwtVerifyToken, gelen token stringini doğrular ve içindeki verileri (claims) döndürür.
// Base64 formatından çözer (decode) ve onu bellekteki jwt.Token yapısına yükler.
// Bu fonksiyon, bir token'ın sadece şifreleme imzası (signature) üzerinden doğruluğunu kontrol eder.
// Token'ın "kullanılıp kullanılmadığını" (blacklist) veya "kullanıcı tarafından iptal edilip edilmediğini" kontrol etmez.
// Bu tür kontroler middleware'de (AuthRequired) ayrıca yapılır.
func JwtVerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, parseError := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Type Assertion: Algoritmanın HMAC tipinde bir yapı (struct) olup olmadığını kontrol eder.
		// Eğer token'ın imzalama algoritması HMAC (Hash-based Message Authentication Code) değilse,
		// token geçersiz kabul edilir ve jwt.ErrSignatureInvalid hatası döndürülür.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})
	// Kütüphane, exp değerindeki zaman damgasını (timestamp) sunucunun o anki saati ile karşılaştırır.
	// Eğer şu anki zaman, exp zamanını geçmişse,
	if parseError != nil {
		return nil, parseError
	}

	// Claims Type Assertion: Token'ın içindeki verilerin jwt.MapClaims tipinde olup olmadığını kontrol eder.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
