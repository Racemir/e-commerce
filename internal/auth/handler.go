package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Veritabanı bağlantısı taşıyan struct
type Handler struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

// Server içinde çalışan veritabanı bağlantısını "db" parametre olarak alır
// ve Handler struct'ının içindeki DB alanına yerleştirir hazır olan struct'ı geri döndürür
func NewHandler(db *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{
		DB:  db,
		RDB: rdb,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// fonksiyonun içine parametre olarak db sokamıyorum, veritabanına ulaşmak için
// Handler struct'ına bir metot olarak bağladım
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	// http.Request yapısının içinde bulunan bir string (metin) değişkenidir.
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Json verisini tutacak boş nesne.
	var req RegisterRequest

	newDecoderErr := json.NewDecoder(r.Body).Decode(&req)
	if newDecoderErr != nil {
		http.Error(w, "Invalid JSON format/Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	// Basit Veri Doğrulama (Validation)
	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "The name, email, and password fields cannot be left blank/İsim, email ve şifre alanları boş bırakılamaz", http.StatusBadRequest)
		return
	}

	// Chack duplicate email
	exists, checkDuplicateEmailError := CheckDuplicateEmail(r.Context(), h.DB, req.Email)
	if checkDuplicateEmailError != nil {
		http.Error(w, "Email check error/E-posta kontrol hatası", http.StatusInternalServerError)
		return
	}
	if exists == true {
		http.Error(w, "This email address is already in use/Bu e-posta adresi zaten kullanılıyo", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, hashedPasswordError := HashPassword(req.Password)
	if hashedPasswordError != nil {
		http.Error(w, "An error occurred during encryption/Şifreleme sırasında bir hata oluştu", http.StatusInternalServerError)
		return
	}

	// Create user
	req.Password = hashedPassword
	newUserID, createUserError := CreateUser(r.Context(), h.DB, req.Name, req.Email, hashedPassword)
	if createUserError != nil {
		http.Error(w, "User could not be saved/Kullanıcı kaydedilemedi", http.StatusInternalServerError)
		return
	}

	// Create verification token
	// Token Son kullanam tarihi
	token, generateSecureTokenError := GenerateSecureToken()
	if generateSecureTokenError != nil {
		http.Error(w, "Token Could Not Be Generated/Token Üretilemedi", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	createVerificationTokenError := CreateVerificationToken(r.Context(), h.DB, newUserID, token, expiresAt)
	if createVerificationTokenError != nil {
		http.Error(w, "Verification Token Could Not Be Generated/Doğrulama Jetonu Oluşturulamadı", http.StatusInternalServerError)
		return
	}

	// Queue verification email
	// h.RDB (Redis Client) Handler yapısından geliyor
	queueVerificationEmailError := QueueVerificationEmail(r.Context(), h.RDB, req.Email, token)
	if queueVerificationEmailError != nil {
		http.Error(w, "Could not queue verification email/Doğrulama e-postası kuyruğa eklenemedi", http.StatusInternalServerError)
		return
	}
	// Retrun response
	w.Header().Set("Content-Type", "application/json")
	// Başlıklar (Header) ayarlandıktan sonra yanıt gövdesi (Body) yazılmadan önce çağırılır.
	w.WriteHeader(http.StatusCreated)
	response := map[string]any{"message": "Kullanıcı başarıyla oluşturuldu!", "user_id": newUserID}
	json.NewEncoder(w).Encode(response)
}
