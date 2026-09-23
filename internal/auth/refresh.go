package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// Refresh, süresi dolmuş Access Token'ı Refresh Token kullanarak yeniler.
//
// Akış:
//  1. "refresh_token" cookie'yi oku
//  2. JWT doğrula → user_id ve token_id çek
//  3. Redis'te refresh:<userID>:<tokenID> key'ini doğrula
//  4. Yeni Access Token üret → cookie set et
//  5. Yeni Refresh Token üret (yeni token_id ile) → cookie set et   ← Rotation
//  6. Eski Refresh Token'ı Redis'ten sil, yeniyi kaydet             ← Rotation
//
// Rotation garantisi: Çalınan bir Refresh Token yalnızca bir kez kullanılabilir.
// İkinci kullanımda Redis'te key bulunamaz → 401.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {

	// Sadece POST methodu kabul edilir
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	// 1. Refresh Token cookie'yi oku 
	refreshTokenString, cookieErr := GetRefreshTokenCookie(r)
	if cookieErr != nil || refreshTokenString == "" {
		http.Error(w, "Unauthorized/Refresh token bulunamadı", http.StatusUnauthorized)
		return
	}

	// 2. JWT doğrula
	claims, jwtErr := JwtVerifyToken(refreshTokenString)
	if jwtErr != nil {
		http.Error(w, "Unauthorized/Refresh token geçersiz veya süresi dolmuş", http.StatusUnauthorized)
		return
	}

	// Token tipini kontrol et — sadece refresh token kabul edilir
	if tokenType, _ := claims["type"].(string); tokenType != "refresh" {
		http.Error(w, "Unauthorized/Geçersiz token tipi", http.StatusUnauthorized)
		return
	}

	// Claims'den user_id ve token_id çek
	userIDRaw, hasUserID := claims["user_id"]
	tokenIDRaw, hasTokenID := claims["token_id"]
	if !hasUserID || !hasTokenID {
		http.Error(w, "Unauthorized/Token içeriği eksik", http.StatusUnauthorized)
		return
	}
	userID := int(userIDRaw.(float64))
	oldTokenID, _ := tokenIDRaw.(string)

	// 3. Redis'te Refresh Token'ı doğrula
	valid, validateErr := ValidateRefreshToken(r.Context(), h.RDB, userID, oldTokenID, refreshTokenString)
	if validateErr != nil {
		http.Error(w, "Internal Server Error/Oturum doğrulanamadı", http.StatusInternalServerError)
		return
	}
	if !valid {
		// Token rotate edilmiş veya logout olunmuş — muhtemel token hırsızlığı
		http.Error(w, "Unauthorized/Oturum sonlandırılmış veya token zaten kullanılmış", http.StatusUnauthorized)
		return
	}

	// 4. Kullanıcı bilgilerini DB'den al (güncel rol için)
	user, userErr := GetUserByID(r.Context(), h.DB, userID)
	if userErr != nil {
		http.Error(w, "Unauthorized/Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	// 5. Yeni Access Token üret
	newAccessToken, accessErr := JwtCreateAccessToken(user.ID, user.Email, user.Role)
	if accessErr != nil {
		http.Error(w, "Internal Server Error/Access token oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// 6. Refresh Token Rotation
	// Yeni benzersiz token_id üret
	newTokenID, tokenIDErr := GenerateSecureToken()
	if tokenIDErr != nil {
		http.Error(w, "Internal Server Error/Token ID oluşturulamadı", http.StatusInternalServerError)
		return
	}

	newRefreshToken, refreshErr := JwtCreateRefreshToken(user.ID, newTokenID)
	if refreshErr != nil {
		http.Error(w, "Internal Server Error/Refresh token oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// Eski Refresh Token'ı Redis'ten sil
	_ = DeleteRefreshToken(r.Context(), h.RDB, userID, oldTokenID)

	// Yeni Refresh Token'ı Redis'e kaydet (7 gün TTL)
	storeErr := StoreRefreshToken(r.Context(), h.RDB, user.ID, newTokenID, newRefreshToken, 7*24*time.Hour)
	if storeErr != nil {
		http.Error(w, "Internal Server Error/Oturum güncellenemedi", http.StatusInternalServerError)
		return
	}

	// 7. Cookie'leri güncelle
	SetAccessTokenCookie(w, newAccessToken)
	SetRefreshTokenCookie(w, newRefreshToken)

	// Başarılı yanıt
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"refreshed": true,
	})
}
