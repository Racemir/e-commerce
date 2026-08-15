package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Sağlık kontrolü için json kalıbı
type Healthresponse struct {
	Status string `json:"status"`
}

// Mux oluşturup geri dönene ana fonksiyon
func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/api/health", HealthHandler)

	return mux
}

// Anasayfa fonksiyonu
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "E-commerce Application")
}

// Health endpoint fonksiyonu
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	healthData := Healthresponse{Status: "ok"}
	json.NewEncoder(w).Encode(healthData)
}
