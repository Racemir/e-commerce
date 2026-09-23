package auth

import (
	"net/http"
	"os"
)

// isProduction, uygulamanın production ortamında çalışıp çalışmadığını kontrol eder.
// APP_ENV ortam değişkeni "production" ise true döner.
func isProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

// Access Token Cookie

// SetAccessTokenCookie, tarayıcıya kısa ömürlü (15 dk) Access Token cookie'si gönderir.
// Bu cookie her HTTP isteğinde sunucuya otomatik iletilir.
//
// Bayraklar:
//   - HttpOnly: JavaScript erişemez (XSS koruması)
//   - SameSite=Strict: CSRF koruması
//   - Path="/": Tüm yollarda geçerli
//   - MaxAge=900: 15 dakika
func SetAccessTokenCookie(w http.ResponseWriter, tokenString string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 dakika
	})
}

// ClearAccessTokenCookie, tarayıcıdaki access_token cookie'sini temizler.
func ClearAccessTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

// GetAccessTokenCookie, gelen istekteki access_token cookie'sini okur.
func GetAccessTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Refresh Token Cookie

// SetRefreshTokenCookie, tarayıcıya uzun ömürlü (7 gün) Refresh Token cookie'si gönderir.
//
// Kritik güvenlik özelliği:
//   - Path="/api/auth/refresh": Cookie YALNIZCA bu endpoint'e gider.
//     Diğer API çağrılarında tarayıcı bu cookie'yi göndermez — saldırı yüzeyi minimize edilir.
//   - MaxAge=604800: 7 gün
func SetRefreshTokenCookie(w http.ResponseWriter, tokenString string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenString,
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 gün
	})
}

// ClearRefreshTokenCookie, tarayıcıdaki refresh_token cookie'sini temizler.
func ClearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/refresh",
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

// GetRefreshTokenCookie, gelen istekteki refresh_token cookie'sini okur.
func GetRefreshTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Yardımcı
// ClearAllAuthCookies, logout sırasında hem access hem refresh cookie'sini temizler.
func ClearAllAuthCookies(w http.ResponseWriter) {
	ClearAccessTokenCookie(w)
	ClearRefreshTokenCookie(w)
}
