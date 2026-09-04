package auth

import (
	"encoding/json"
	"net/http"
)

// Me, oturum açmış olan mevcut kullanıcının bilgilerini döner. (PDF Madde 34 - Current User)
//
// Akış:
//  İstekten session cookie'sini oku
//  Redis'ten oturum bilgilerini doğrula
//  Veritabanından güncel kullanıcı bilgilerini çek
//  Kullanıcı bilgilerini JSON olarak döndür
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {

	// Sadece GET methodu kabul edilir
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	// Cookie'den session ID oku
	sessionID, cookieError := GetSessionCookie(r)
	if cookieError != nil || sessionID == "" {
		http.Error(w, "Unauthorized/Giriş yapılmamış", http.StatusUnauthorized)
		return
	}

	// Redis'ten session verisini kontrol et
	sessionData, sessionError := GetSession(r.Context(), h.RDB, sessionID)
	if sessionError != nil {
		http.Error(w, "Unauthorized/Oturum geçersiz veya süresi dolmuş", http.StatusUnauthorized)
		return
	}

	// Veritabanından kullanıcı bilgilerini al
	user, getUserError := GetUserByID(r.Context(), h.DB, sessionData.UserID)
	if getUserError != nil {
		http.Error(w, "User not found/Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	// JSON olarak döndür
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
