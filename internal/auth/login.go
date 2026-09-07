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
	if isEmailTaken == false {
		http.Error(w, "Email not found/E-posta bulunamadı", http.StatusUnauthorized)
		return
	}

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
	if match == false {
		http.Error(w, "Invalid password/Geçersiz şifre", http.StatusUnauthorized)
		return
	}

	// Şifre eşleşirse jwt token oluştur.
	// SHA256 , user_id role , iss
	token, jwtError := JwtCreateToken(user.ID, user.Role)
	if jwtError != nil {
		http.Error(w, "Token creation error/Token oluşturma hatası", http.StatusInternalServerError)
		return
	}

	// Oturum (Session) oluştur ve Redis'e kaydet (PDF Madde 34 & 35 gereksinimi)
	sessionID, generateSessionError := GenerateSecureToken()
	if generateSessionError != nil {
		http.Error(w, "Session generation error/Oturum anahtarı üretilemedi", http.StatusInternalServerError)
		return
	}

	sessionData := SessionData{
		UserID: user.ID,
		Role:   user.Role,
	}

	createSessionError := CreateSession(r.Context(), h.RDB, sessionID, sessionData, 24*time.Hour)
	if createSessionError != nil {
		http.Error(w, "Session could not be created/Oturum oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// Tarayıcıya güvenli session cookie'sini gönder
	SetSessionCookie(w, sessionID)

	// Tokenı response body'e koy.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
