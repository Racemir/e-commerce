package auth

import (
	"context"
	"net/http"
)

func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cookie'den token'ı oku
		tokenString, cookieError := GetSessionCookie(r)
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

		// JWT içindeki verilerden User objesini oluştur (Veritabanı sorgusunu atlıyoruz)
		userID := int(claims["user_id"].(float64))
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		user := &User{
			ID:    userID,
			Email: email,
			Role:  role,
		}
		
		r = r.WithContext(context.WithValue(r.Context(), userContextKey, user))
		next.ServeHTTP(w, r)
	})

}
