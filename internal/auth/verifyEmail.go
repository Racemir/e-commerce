package auth

import (
	"encoding/json"
	"net/http"
)

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

// (h *Handler) kullanarak handler.go'daki DB bağlantısına erişim sağlıyoruz
func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {

	// isteğin Post metodu olmalı
	if http.MethodPost != r.Method {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}
	// Gelen json verisini VerifyEmailRequest yapısında çöz
	var verify VerifyEmailRequest
	decodeEror := json.NewDecoder(r.Body).Decode(&verify)
	if decodeEror != nil {
		http.Error(w, "Invalid JSON format/Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}
	// Token'ın boş olup olmadığını konrol edicez (validation)
	if verify.Token == "" {
		http.Error(w, "The token cannot be empty/Token boş olamaz.", http.StatusBadRequest)
		return
	}
	// Token veritabanı kontrolü yap
	verifyEmailTokenError := VerifyEmailToken(r.Context(), h.DB, verify.Token)
	// Geçersiz token
	if verifyEmailTokenError != nil {
		http.Error(w, "Invalid Token/Geçersiz Token", http.StatusUnauthorized)
		return
	}
	// Geçerli token E-posta onaylandı statusok
	// Yanıt başlığı json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Durum kodu
	json.NewEncoder(w).Encode(map[string]string{"message": "Email successfully verified/E-Posta başarıyla doğrulandı"})
}
