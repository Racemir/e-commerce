package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Racemir/e-commerce/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Sağlık kontrolü için json kalıbı
type Healthresponse struct {
	Status string `json:"status"`
}

// Mux oluşturup geri dönene ana fonksiyon
func SetupRoutes(db *pgxpool.Pool, rdb *redis.Client) *http.ServeMux {
	// Gelen HTTP isteklerinin URL yolarına bakarak fonksiyonlara yönlendirir.
	mux := http.NewServeMux()

	authHandler := auth.NewHandler(db, rdb)
	// Tarayıcıdan girilen adres "/api/health" bu isteği özel fonksiyona (handler) gönderir.
	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/api/health", HealthHandler)
	mux.HandleFunc("/api/auth/register", authHandler.Register)

	return mux
}

// Anasayfa fonksiyonu
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "E-commerce Application")
}

// Health endpoint fonksiyonu
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Header content-type verinin fomratı nedir
	// application: binary veya metin verimiz json formatında
	w.Header().Set("Content-Type", "application/json")
	healthData := Healthresponse{Status: "ok"}

	// NewEncoder veriyi doğrudan HTTP yanıtına (w) yazacak bir dönüştürücü oluşturur.
	// Encode healthData isimli Go yapısını alır, JSON'a çevirir ve (w) üzerinden gönderir.
	// Veriyi bellekte bir değişkene atamadık RAM'i yormadık doğrudan stream (akış) olarak gönderdik.
	json.NewEncoder(w).Encode(healthData)
}
