package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// email ve şifre al
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	// Post methodu kullanılmalı.
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed/İzin Verilmeyen Yöntem", http.StatusMethodNotAllowed)
		return
	}

	// Kullanıcının giriş bilgilerini al. (e-posta ve şifre)
	var req LoginRequest
	newDecoderError := json.NewDecoder(r.Body).Decode(&req)
	if newDecoderError != nil {
		http.Error(w, "Invalid JSON format/Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}
	// Email veya şifre boş olamaz.
	if req.Email == "" || req.Password == "" {
		http.Error(w, "The email and password fields cannot be left blank/E-posta ve şifre alanları boş bırakılamaz", http.StatusBadRequest)
		return
	}

	// Email veritabanında kontrol et. (var ise devam et)
	isEmailTaken, checkDuplicateEmailError := CheckDuplicateEmail(r.Context(), h.DB, req.Email)
	if checkDuplicateEmailError != nil {
		http.Error(w, "Email check error/E-posta kontrol hatası", http.StatusInternalServerError)
		return
	}
	if !isEmailTaken {
		http.Error(w, "Email not found/E-posta bulunamadı", http.StatusUnauthorized)
		return
	}

	// Verilen email'in kullanıcı bilgilerini al (id, email, role, vb.)
	user, getUserError := GetUserByEmail(r.Context(), h.DB, req.Email)
	if getUserError != nil {
		http.Error(w, "User not found/Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	// Girilen şifreyi veritabanındaki hash ile karşılaştır
	match, checkError := CheckPasswordHash(req.Password, user.PasswordHash)
	if checkError != nil {
		http.Error(w, "Password check error/Şifre kontrol hatası", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "Invalid password/Geçersiz şifre", http.StatusUnauthorized)
		return
	}

	// Access Token oluştur (15 dakika)
	accessToken, accessTokenError := JwtCreateAccessToken(user.ID, user.Email, user.Role)
	if accessTokenError != nil {
		http.Error(w, "Token creation error/Token oluşturma hatası", http.StatusInternalServerError)
		return
	}

	// Her refresh token benzersiz bir ID taşır (Rotation için)
	tokenID, tokenIDError := GenerateSecureToken()
	if tokenIDError != nil {
		http.Error(w, "Token ID generation error/Token ID oluşturma hatası", http.StatusInternalServerError)
		return
	}

	// Refresh Token oluştur (7 gün)
	refreshToken, refreshTokenError := JwtCreateRefreshToken(user.ID, tokenID)
	if refreshTokenError != nil {
		http.Error(w, "Refresh token creation error/Refresh token oluşturma hatası", http.StatusInternalServerError)
		return
	}

	// Refresh Token'ı Redis'e kaydet (7 gün TTL)
	storeError := StoreRefreshToken(r.Context(), h.RDB, user.ID, tokenID, refreshToken, 7*24*time.Hour)
	if storeError != nil {
		http.Error(w, "Session store error/Oturum kayıt hatası", http.StatusInternalServerError)
		return
	}

	// Cookie'leri tarayıcıya gönder
	SetAccessTokenCookie(w, accessToken)
	SetRefreshTokenCookie(w, refreshToken)

	// Başarılı giriş yanıtı
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful/Giriş başarılı",
	})
}
