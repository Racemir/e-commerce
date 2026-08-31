package auth

import (
	"encoding/json"
	"net/http"
)

// (h *Handler) kullanarak handler.go'daki DB bağlantısına erişim sağlıyoruz
func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {

	// Kullanıcı e-postadaki linke tıkladığında tarayıcı GET isteği gönderir
	if http.MethodGet != r.Method {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	// Token'ı URL'den al: /api/auth/verify-email?token=abc123
	// r.URL.Query().Get("token") → URL'deki "token" parametresinin değerini döndürür
	token := r.URL.Query().Get("token")

	// Token'ın boş olup olmadığını kontrol edicez (validation)
	if token == "" {
		http.Error(w, "The token cannot be empty/Token boş olamaz.", http.StatusBadRequest)
		return
	}
	// Token veritabanı kontrolü yap
	verifyEmailTokenError := VerifyEmailToken(r.Context(), h.DB, token)
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
