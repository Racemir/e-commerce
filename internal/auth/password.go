package auth

import "github.com/alexedwards/argon2id"

// HashPassword, düz metin şifreyi güvenli Argon2id formatına çevirir (Register işlemi için)
func HashPassword(password string) (string, error) {
	hash, createHashError := argon2id.CreateHash(password, argon2id.DefaultParams)
	if createHashError != nil {
		return "", createHashError
	}
	return hash, nil
}

// CheckPasswordHash, girilen şifre ile veritabanındaki hash'i karşılaştırır (Login işlemi için)
func CheckPasswordHash(password, hash string) (bool, error) {
	match, comparePasswordAndHashError := argon2id.ComparePasswordAndHash(password, hash)
	if comparePasswordAndHashError != nil {
		return false, comparePasswordAndHashError
	}

	return match, nil
}
