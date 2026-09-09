package auth

import (
	"encoding/json"
	"net/http"
)

// Logout, kullanıcının oturumunu sonlandırır.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	// Sadece POST methodu kabul edilir
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	// Cookie'den auth token'ını oku
	_, cookieError := GetSessionCookie(r)
	if cookieError != nil {
		http.Error(w, "No active session/Aktif oturum bulunamadı", http.StatusUnauthorized)
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
