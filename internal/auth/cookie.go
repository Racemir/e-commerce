package auth

import (
	"net/http"
	"os"
)

// cookieName, auth token cookie'sinin adını döner.
// SESSION_COOKIE_NAME veya JWT_COOKIE_NAME ortam değişkeninden okunur, yoksa varsayılan "token" kullanılır.
func cookieName() string {
	name := os.Getenv("JWT_COOKIE_NAME")
	if name == "" {
		name = os.Getenv("SESSION_COOKIE_NAME")
	}
	if name == "" {
		name = "token"
	}
	return name
}

// isProduction, uygulamanın production ortamında çalışıp çalışmadığını kontrol eder.
// APP_ENV ortam değişkeni "production" ise true döner.
func isProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

// SetSessionCookie, tarayıcıya güvenli bir JWT auth cookie'si gönderir.
// Bu cookie, her HTTP isteğinde otomatik olarak sunucuya geri gönderilir.
//
// Bayraklar:
//   - HttpOnly: JavaScript cookie'ye erişemez (XSS koruması)
//   - SameSite=Strict: Başka sitelerden gelen isteklerde cookie gönderilmez (CSRF koruması)
//   - Path="/": Cookie tüm yollarda geçerlidir
//   - MaxAge=86400: 24 saat sonra tarayıcı cookie'yi siler
func SetSessionCookie(w http.ResponseWriter, tokenString string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName(),
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(), // Production'da true (HTTPS zorunlu), development'ta false
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 saat (saniye cinsinden)
	})
}

// ClearSessionCookie, tarayıcıdaki auth cookie'sini temizler.
// MaxAge=-1 ayarlayarak tarayıcıya "bu cookie'yi hemen sil" der.
// Logout işleminde kullanılır.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName(),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Tarayıcıya cookie'yi hemen silmesini söyler
	})
}

// GetSessionCookie, gelen HTTP isteğindeki auth cookie'sini okur.
// Cookie yoksa veya boşsa hata döner.
func GetSessionCookie(r *http.Request) (string, error) {
	cookie, cookieError := r.Cookie(cookieName())
	if cookieError != nil {
		return "", cookieError
	}
	return cookie.Value, nil
}
