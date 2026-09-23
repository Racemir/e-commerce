package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Access Token Blacklist
//
// Logout olan kullanıcıların Access Token'ları burada tutulur.
// TTL = token'ın kalan exp süresi → Redis otomatik temizler, bellek sızıntısı olmaz.

// blacklistKeyPrefix, Access Token blacklist key'lerinin öneki.
// Örnek: "blacklist:<tokenString>"
const blacklistKeyPrefix = "blacklist:"

// BlacklistToken, verilen JWT Access Token'ını Redis'e ekler.
// ttl parametresi token'ın kalan geçerlilik süresidir (exp - now).
// Bu token'ı kara listeye al ve süresi dolduğunda otomatik çöpe at.
func BlacklistToken(ctx context.Context, rdb *redis.Client, tokenString string, ttl time.Duration) error {
	key := blacklistKeyPrefix + tokenString
	// SET key "1" EX <ttl_seconds>
	return rdb.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted, token'ın Redis blacklist'inde olup olmadığını kontrol eder.
// true dönerse token iptal edilmiş demektir — istek reddedilmeli.
func IsBlacklisted(ctx context.Context, rdb *redis.Client, tokenString string) (bool, error) {
	key := blacklistKeyPrefix + tokenString
	// Redis'e git, bu key'in değerini getir. Bulursan val içine koy, bulamazsan veya hata çıkarsa err içine yaz.
	val, getError := rdb.Get(ctx, key).Result()
	if getError == redis.Nil {
		// Key yok → blacklist'te değil
		return false, nil
	}
	if getError != nil {
		// Redis bağlantı hatası vb.
		return false, getError
	}
	return val == "1", nil
}

// Refresh Token Store
//
// Refresh Token'lar Redis'te saklanır.
// Her kullanıcının her cihazı için ayrı bir tokenID (UUID) bulunur.
// Rotation sırasında eski tokenID silinip yenisi yazılır.
//
// Key formatı: "refresh:<userID>:<tokenID>"
// Örnek:       "refresh:42:a3f9..."

func refreshKey(userID int, tokenID string) string {
	return fmt.Sprintf("refresh:%d:%s", userID, tokenID)
}

// StoreRefreshToken, yeni bir Refresh Token'ı Redis'e kaydeder.
// ttl genellikle 7 gündür.
func StoreRefreshToken(ctx context.Context, rdb *redis.Client, userID int, tokenID string, tokenString string, ttl time.Duration) error {
	key := refreshKey(userID, tokenID)
	return rdb.Set(ctx, key, tokenString, ttl).Err()
}

// ValidateRefreshToken, Redis'teki Refresh Token ile gelen token'ın eşleşip eşleşmediğini kontrol eder.
// false + nil → geçersiz token (key yok veya değer farklı)
// false + err → Redis hatası
func ValidateRefreshToken(ctx context.Context, rdb *redis.Client, userID int, tokenID string, tokenString string) (bool, error) {
	key := refreshKey(userID, tokenID)
	// Redis'e git, bu key'in değerini getir. Bulursan val içine koy, bulamazsan veya hata çıkarsa err içine yaz.
	stored, getError := rdb.Get(ctx, key).Result()
	if getError == redis.Nil {
		// Key yok → token daha önce rotate edilmiş veya logout olunmuş
		return false, nil
	}
	if getError != nil {
		// Redis bağlantı hatası vb.
		return false, getError
	}
	return stored == tokenString, nil
}

// DeleteRefreshToken, Refresh Token'ı Redis'ten siler.
// Logout veya Rotation sırasında çağrılır.
func DeleteRefreshToken(ctx context.Context, rdb *redis.Client, userID int, tokenID string) error {
	key := refreshKey(userID, tokenID)
	return rdb.Del(ctx, key).Err()
}
