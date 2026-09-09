package auth

import (
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

// AdminOnly, sadece admin rolüne sahip kullanıcıların geçişine izin verir.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Contextten kullanıcı verisini çeker.
		user := r.Context().Value(userContextKey)
		// Eğer kullanıcı yoksa Unauthorized hatası döner
		if user == nil {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}
		// Eğer yetkili yoksa Forbidden hatası döner.
		if user.(*User).Role != "admin" {
			http.Error(w, "Admin role required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
