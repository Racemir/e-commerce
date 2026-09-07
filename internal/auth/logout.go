package auth

import (
	"encoding/json"
	"net/http"
)

// Logout, kullanıcının oturumunu sonlandırır.
//
// Akış:
//  1. Cookie'den session ID'yi oku
//  2. Redis'ten session kaydını sil
//  3. Tarayıcıdaki cookie'yi temizle
//
// Bu işlem sonrası aynı session ID ile gelen istekler reddedilecektir.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	// Sadece POST methodu kabul edilir
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	// Cookie'den session ID'yi oku
	sessionID, cookieError := GetSessionCookie(r)
	if cookieError != nil {
		http.Error(w, "No active session/Aktif oturum bulunamadı", http.StatusUnauthorized)
		return
	}

	// Redis'ten session kaydını sil
	deleteError := DeleteSession(r.Context(), h.RDB, sessionID)
	if deleteError != nil {
		http.Error(w, "Session deletion error/Oturum silme hatası", http.StatusInternalServerError)
		return
	}

	// Tarayıcıdaki cookie'yi temizle
	ClearSessionCookie(w)

	// Başarılı çıkış yanıtı
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful/Çıkış başarılı",
	})
}
