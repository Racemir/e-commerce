package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// şifreyi sıfırlamak isteyen müşterinin e-postası al.(şifresi değişekecek müşteri)
// e-posta onaylanıyorsa (verify_email) şifre sıfırlama linki e-postaya gönderilir. (e-posta da hem link hem token olur)
// token hem veritabanında hem redis de tutulur.(72 saat geçerli olur)
// e-posta gidip linke tıklayan müşteri yeni şifresini oluşturur.
// şifre değiştirildikten sonra token veritabanından ve redis'den silinir.

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ForgotPassword, şifre sıfırlama isteği başlatır.
// Kullanıcının e-postasını alır, password_reset_tokens tablosuna token kaydeder
// ve sıfırlama linkini e-posta ile gönderir.
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	var req ForgotPasswordRequest

	decodeError := json.NewDecoder(r.Body).Decode(&req)
	if decodeError != nil {
		http.Error(w, "Invalid JSON format/Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "The email field cannot be left blank/E-posta alanı boş bırakılamaz", http.StatusBadRequest)
		return
	}

	// e-posta veritabanında var mı?
	user, getUserError := GetUserByEmail(r.Context(), h.DB, req.Email)
	if getUserError != nil {
		// Güvenlik: E-posta bulunamasa bile aynı mesajı döndür (e-posta enumeration koruması)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "If such an email is registered in our system, a reset link has been sent./Eğer sistemimizde böyle bir e-posta kayıtlıysa, sıfırlama bağlantısı gönderilmiştir.",
		})
		return
	}

	// Yeni token üret
	token, generateSecureTokenError := GenerateSecureToken()
	if generateSecureTokenError != nil {
		http.Error(w, "Token Could Not Be Generated/Token Üretilemedi", http.StatusInternalServerError)
		return
	}

	// 15 dakika geçerli token
	expiresAt := time.Now().Add(15 * time.Minute)

	// Token'ı password_reset_tokens tablosuna kaydet
	createPasswordResetTokenError := CreatePasswordResetToken(r.Context(), h.DB, user.ID, token, expiresAt)
	if createPasswordResetTokenError != nil {
		http.Error(w, "Failed to create reset token./Sıfırlama jetonu oluşturulamadı.", http.StatusInternalServerError)
		return
	}

	// Redis queue'a ekle (e-posta gönderimi için)
	queuePasswordResetEmailError := QueuePasswordResetEmail(r.Context(), h.RDB, req.Email, token)
	if queuePasswordResetEmailError != nil {
		http.Error(w, "Failed to queue reset email./Sıfırlama e-postası kuyruğa alınamadı.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If such an email is registered in our system, a reset link has been sent./Eğer sistemimizde böyle bir e-posta kayıtlıysa, sıfırlama bağlantısı gönderilmiştir.",
	})
}

// ResetPassword, şifre sıfırlama token'ını doğrular ve yeni şifreyi kaydeder.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	var req ResetPasswordRequest

	decodeError := json.NewDecoder(r.Body).Decode(&req)
	if decodeError != nil {
		http.Error(w, "Invalid JSON format/Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.Password == "" {
		http.Error(w, "The password and token fields cannot be left blank/Şifre ve token alanları boş bırakılamaz", http.StatusBadRequest)
		return
	}

	// Token'ı password_reset_tokens tablosunda doğrula ve kullanıcı ID'sini al
	userID, verifyError := VerifyPasswordResetToken(r.Context(), h.DB, req.Token)
	if verifyError != nil {
		http.Error(w, "Invalid or expired token/Geçersiz veya süresi dolmuş token", http.StatusBadRequest)
		return
	}

	// Yeni şifreyi hashle
	hashedPassword, hashError := HashPassword(req.Password)
	if hashError != nil {
		http.Error(w, "An error occurred during encryption/Şifreleme sırasında bir hata oluştu", http.StatusInternalServerError)
		return
	}

	// Kullanıcının şifresini güncelle
	updateError := UpdateUserPassword(r.Context(), h.DB, userID, hashedPassword)
	if updateError != nil {
		http.Error(w, "Failed to update password./Şifre güncellenemedi.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been reset successfully./Şifre başarıyla sıfırlandı.",
	})
}
