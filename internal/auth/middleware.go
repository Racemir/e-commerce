package auth

import (
	"context"
	"net/http"
)

// RequireAuth, korumalı endpoint'lere erişim için kimlik doğrulama middleware'i.
//
// Akış:
//  1. "access_token" cookie'yi oku
//  2. JWT doğrula (imza + süre)
//  3. Redis blacklist kontrolü (logout edilmiş token'ları reddet)
//  4. Claims'den User oluştur → context'e yaz
//  5. Sonraki handler'ı çağır
func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cookie'den Access Token'ı oku
		tokenString, cookieError := GetAccessTokenCookie(r)
		if cookieError != nil || tokenString == "" {
			http.Error(w, "Unauthorized/Giriş yapılmamış", http.StatusUnauthorized)
			return
		}

		// JWT token'ını doğrula
		claims, err := JwtVerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized/Oturum geçersiz veya süresi dolmuş", http.StatusUnauthorized)
			return
		}

		// Redis blacklist kontrolü — logout edilmiş token'ları reddet
		blacklisted, blacklistErr := IsBlacklisted(r.Context(), h.RDB, tokenString)
		if blacklistErr != nil {
			http.Error(w, "Internal Server Error/Oturum kontrol edilemedi", http.StatusInternalServerError)
			return
		}
		if blacklisted {
			http.Error(w, "Unauthorized/Oturum sonlandırılmış", http.StatusUnauthorized)
			return
		}

		// JWT içindeki verilerden User objesini oluştur (Veritabanı sorgusunu atlıyoruz bu bilgileri doğrudan JWT'nin içinden okuyacağız.)
		userID := int(claims["user_id"].(float64))
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		// Claims'den kullanıcı oluştur
		user := &User{
			ID:    userID,
			Email: email,
			Role:  role,
		}

		// Context'e yaz, bir sonraki handler'a ilet
		r = r.WithContext(context.WithValue(r.Context(), userContextKey, user))
		next.ServeHTTP(w, r)
	})
}
