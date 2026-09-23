package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Racemir/e-commerce/internal/auth"
	"github.com/Racemir/e-commerce/internal/orders"
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
	ordersHandler := orders.NewHandler(db)

	// Normal Rotalar
	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/api/health", HealthHandler)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/verify-email", authHandler.VerifyEmail)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("/api/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/api/auth/refresh", authHandler.Refresh)


	// Korumalı Rotalar
	mux.Handle("/api/auth/me", authHandler.RequireAuth(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/auth/logout", authHandler.RequireAuth(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/orders/{id}", authHandler.RequireAuth(http.HandlerFunc(ordersHandler.GetOrder)))

	// Sipariş Rotaları (Resource Ownership korumalı)
	mux.Handle("/api/orders/{id}", authHandler.RequireAuth(http.HandlerFunc(ordersHandler.GetOrder)))

	// Admin Rotaları
	adminMux := http.NewServeMux()

	mux.Handle("/api/admin/", authHandler.RequireAuth(auth.AdminOnly(adminMux)))
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
