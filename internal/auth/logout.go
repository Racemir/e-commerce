package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// Logout, kullanıcının oturumunu sonlandırır.
// Access Token blacklist'e eklenir; Refresh Token Redis'ten silinir.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	// Sadece POST methodu kabul edilir
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	// Access Token'ı blacklist'e ekle
	accessTokenString, accessCookieErr := GetAccessTokenCookie(r)
	if accessCookieErr == nil && accessTokenString != "" {
		claims, jwtErr := JwtVerifyToken(accessTokenString)
		if jwtErr == nil {
			// Token'ın kalan süresini hesapla (Redis TTL için)
			if expRaw, hasExp := claims["exp"]; hasExp {
				// Token son kullanım tarihini int64'e çevir
				expUnix := int64(expRaw.(float64))
				// Token bitiş tarihine şu an itibariyle ne kadar kaldıysa o süreyi al
				ttl := time.Until(time.Unix(expUnix, 0))
				if ttl > 0 {
					// Token henüz geçerliyken blacklist'e ekle
					_ = BlacklistToken(r.Context(), h.RDB, accessTokenString, ttl)
				}
			}
		}
	}

	// Refresh Token'ı Redis'ten sil
	refreshTokenString, refreshCookieErr := GetRefreshTokenCookie(r)
	if refreshCookieErr == nil && refreshTokenString != "" {
		claims, jwtErr := JwtVerifyToken(refreshTokenString)
		if jwtErr == nil {
			userIDRaw, hasUserID := claims["user_id"]
			tokenIDRaw, hasTokenID := claims["token_id"]
			if hasUserID && hasTokenID {
				userID := int(userIDRaw.(float64))
				tokenID, _ := tokenIDRaw.(string)
				_ = DeleteRefreshToken(r.Context(), h.RDB, userID, tokenID)
			}
		}
	}

	// Her iki cookie'yi temizle
	ClearAllAuthCookies(w)

	// Başarılı çıkış yanıtı
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful/Çıkış başarılı",
	})
}
