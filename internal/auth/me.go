package auth

import (
	"encoding/json"
	"net/http"
)

// Me, oturum açmış olan mevcut kullanıcının bilgilerini döner. (Current User)
//
// Akış:
//
//	İstekten session cookie'sini oku
//	Redis'ten oturum bilgilerini doğrula
//	Veritabanından güncel kullanıcı bilgilerini çek
//	Kullanıcı bilgilerini JSON olarak döndür
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {

	// Sadece GET methodu kabul edilir
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	user := r.Context().Value(userContextKey)
	// JSON olarak döndür
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
